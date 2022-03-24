/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package role

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
	apiauthzv1 "tkestack.io/tke/api/authz/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authzv1informer "tkestack.io/tke/api/client/informers/externalversions/authz/v1"
	authzv1 "tkestack.io/tke/api/client/listers/authz/v1"
	"tkestack.io/tke/pkg/authz/constant"
	"tkestack.io/tke/pkg/authz/controller/policyrolecache"
	authzprovider "tkestack.io/tke/pkg/authz/provider"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	appDeletionGracePeriod = 5 * time.Second
	controllerName         = "role-controller"
)

type Controller struct {
	client     clientset.Interface
	queue      workqueue.RateLimitingInterface
	roleLister authzv1.RoleLister
	roleSynced cache.InformerSynced
	mcrbLister authzv1.MultiClusterRoleBindingLister
	mcrbSynced cache.InformerSynced
	stopCh     <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	roleInformer authzv1informer.RoleInformer,
	mcrbInformer authzv1informer.MultiClusterRoleBindingInformer,
	resyncPeriod time.Duration) *Controller {
	controller := &Controller{
		client: client,
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}
	if client != nil &&
		client.AuthzV1().RESTClient() != nil &&
		!reflect.ValueOf(client.AuthzV1().RESTClient()).IsNil() &&
		client.AuthzV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.AuthzV1().RESTClient().GetRateLimiter())
	}

	roleInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					role := obj.(*apiauthzv1.Role)
					controller.enqueue(obj)
					policyrolecache.Cache.PutByRole(role)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					oldRole := oldObj.(*apiauthzv1.Role)
					newRole := newObj.(*apiauthzv1.Role)
					controller.enqueue(newObj)
					policyrolecache.Cache.UpdateByRole(oldRole, newRole)
				},
				DeleteFunc: func(obj interface{}) {
					role, _ := obj.(*apiauthzv1.Role)
					controller.enqueue(obj)
					policyrolecache.Cache.DeleteRole(role)
				},
			},
			FilterFunc: func(obj interface{}) bool {
				role, ok := obj.(*apiauthzv1.Role)
				if !ok || role.Scope != apiauthzv1.MultiClusterScope {
					return false
				}
				provider, err := authzprovider.GetProvider(role.Annotations)
				if err != nil {
					return true
				}
				return provider.OnFilter(context.TODO(), role.Annotations)
			},
		},
		resyncPeriod,
	)

	controller.roleLister = roleInformer.Lister()
	controller.roleSynced = roleInformer.Informer().HasSynced
	controller.mcrbLister = mcrbInformer.Lister()
	controller.mcrbSynced = mcrbInformer.Informer().HasSynced
	return controller
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, appDeletionGracePeriod)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting role controller")
	defer log.Info("Shutting down role controller")

	if ok := cache.WaitForCacheSync(stopCh, c.roleSynced, c.mcrbSynced); !ok {
		log.Error("Failed to wait for role caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of app objects.
// Each app can be in the queue at most once.
// The system ensures that no two workers can process
// the same app at the same time.
func (c *Controller) worker() {
	workFunc := func() bool {
		key, quit := c.queue.Get()
		if quit {
			return true
		}
		defer c.queue.Done(key)

		err := c.syncItem(key.(string))
		if err == nil {
			// no error, forget this entry and return
			c.queue.Forget(key)
			return false
		}

		// rather than wait for a full resync, re-add the app to the queue to be processed
		c.queue.AddRateLimited(key)
		runtime.HandleError(err)
		return false
	}

	for {
		quit := workFunc()
		if quit {
			return
		}
	}
}

func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing role", log.String("role", key), log.Duration("processTime", time.Since(startTime)))
	}()
	var roleDeleted bool
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	role, err := c.roleLister.Roles(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Role has been deleted. Attempting to cleanup resources",
				log.String("namespace", ns),
				log.String("name", name))
			// 删除mcrb资源
			roleDeleted = true
			mcrbs, err := c.getMultiClusterRoleBindings(ns, name)
			if err != nil {
				log.Warn("Unable to retrieve MultiClusterRoleBindings from store",
					log.String("roleNs", ns),
					log.String("roleName", name), log.Err(err))
				return err
			}
			err = c.updateOrDeleteMultiClusterRoleBindings(mcrbs, roleDeleted)
			if err != nil {
				log.Warn("Unable to update MultiClusterRoleBindings",
					log.String("roleNs", ns),
					log.String("roleName", name), log.Err(err))
			}
			return err
		}
		log.Warn("Unable to retrieve role from store",
			log.String("namespace", ns),
			log.String("name", name), log.Err(err))
		return err
	}
	roleDeleted = role.DeletionTimestamp != nil
	role = role.DeepCopy()

	mcrbs, err := c.getMultiClusterRoleBindings(ns, name)
	if err != nil {
		log.Warn("Unable to retrieve MultiClusterRoleBindings from store",
			log.String("roleNs", ns),
			log.String("roleName", name), log.Err(err))
		return err
	}

	err = c.updateOrDeleteMultiClusterRoleBindings(mcrbs, roleDeleted)
	if err != nil {
		log.Warn("Unable to update MultiClusterRoleBindings",
			log.String("roleNs", ns),
			log.String("roleName", name), log.Err(err))
		return err
	}

	if roleDeleted {
		roleFinalize := apiauthzv1.Role{}
		roleFinalize.ObjectMeta = role.ObjectMeta
		roleFinalize.Finalizers = []string{}
		if err := c.client.AuthzV1().RESTClient().Put().Resource("roles").
			Namespace(ns).
			Name(name).
			SubResource("finalize").
			Body(&roleFinalize).
			Do(context.Background()).
			Into(&roleFinalize); err != nil {
			log.Warnf("Unable to finalize role '%s/%s', err: %v", ns, name, err)
			return err
		}
		if err = c.client.AuthzV1().Roles(ns).Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				log.Warnf("Unable to delete role '%s/%s', err: %v", ns, name, err)
				return err
			}
		}
	}
	return nil
}

func (c *Controller) getMultiClusterRoleBindings(roleNs, roleName string) ([]*apiauthzv1.MultiClusterRoleBinding, error) {
	var mcrbs []*apiauthzv1.MultiClusterRoleBinding
	var err error
	selector := labels.SelectorFromSet(map[string]string{
		constant.RoleNamespace: roleNs,
		constant.RoleName:      roleName,
	})
	if roleNs == "" {
		mcrbs, err = c.mcrbLister.List(selector)
	} else {
		mcrbs, err = c.mcrbLister.MultiClusterRoleBindings(roleNs).List(selector)
	}
	return mcrbs, err
}

func (c *Controller) updateOrDeleteMultiClusterRoleBindings(mcrbs []*apiauthzv1.MultiClusterRoleBinding, roleDeleted bool) error {
	for _, mcrb := range mcrbs {
		if mcrb.DeletionTimestamp != nil {
			continue
		}
		if roleDeleted {
			deleteOpt := metav1.DeletePropagationBackground
			if err := c.client.AuthzV1().MultiClusterRoleBindings(mcrb.Namespace).Delete(context.Background(), mcrb.Name, metav1.DeleteOptions{PropagationPolicy: &deleteOpt}); err != nil {
				if errors.IsNotFound(err) {
					continue
				} else {
					log.Warnf("Unable to delete MultiClusterRoleBinding '%s/%s', err: '%v'", mcrb.Namespace, mcrb.Name, err)
					return err
				}
			}
		} else {
			// 触发更新mcrb
			deepCopy := mcrb.DeepCopy()
			annotations := deepCopy.Annotations
			if annotations == nil {
				annotations = map[string]string{}
			}
			annotations[constant.UpdatedByRoleController] = time.Now().Format("2006-01-02T15:04:05")
			deepCopy.Annotations = annotations
			if _, err := c.client.AuthzV1().MultiClusterRoleBindings(mcrb.Namespace).Update(context.Background(), deepCopy, metav1.UpdateOptions{}); err != nil {
				if errors.IsNotFound(err) {
					continue
				} else {
					log.Warnf("Unable to delete MultiClusterRoleBinding '%s/%s', err: '%v'", mcrb.Namespace, mcrb.Name, err)
					return err
				}
			}
		}
	}
	return nil
}
