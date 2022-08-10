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

package policy

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	deletionGracePeriod = 5 * time.Second
)

const (
	controllerName = "policy-controller"
)

type Controller struct {
	client       clientset.Interface
	queue        workqueue.RateLimitingInterface
	policyLister authzv1.PolicyLister
	policySynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	policyInformer authzv1informer.PolicyInformer,
	resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
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

	policyInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					controller.enqueue(obj)
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					old, ok1 := oldObj.(*apiauthzv1.Policy)
					cur, ok2 := newObj.(*apiauthzv1.Policy)
					if ok1 && ok2 && controller.needsUpdate(old, cur) {
						controller.enqueue(cur)
					}
				},
				DeleteFunc: func(obj interface{}) {
					controller.enqueue(obj)
				},
			},
			FilterFunc: func(obj interface{}) bool {
				policy, ok := obj.(*apiauthzv1.Policy)
				if !ok || policy.Scope != apiauthzv1.MultiClusterScope {
					return false
				}
				provider, err := authzprovider.GetProvider(policy.Annotations)
				if err != nil {
					return true
				}
				return provider.OnFilter(context.TODO(), policy.Annotations)
			},
		},
		resyncPeriod,
	)
	controller.policyLister = policyInformer.Lister()
	controller.policySynced = policyInformer.Informer().HasSynced
	return controller
}

func (c *Controller) needsUpdate(old *apiauthzv1.Policy, new *apiauthzv1.Policy) bool {
	if old.UID != new.UID {
		return true
	}
	if !reflect.DeepEqual(old.Rules, new.Rules) {
		return true
	}
	return false
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, deletionGracePeriod)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting policy controller")
	defer log.Info("Shutting down policy controller")

	if ok := cache.WaitForCacheSync(stopCh, c.policySynced); !ok {
		log.Error("Failed to wait for policy caches to sync")
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

		_, err := c.syncItem(key.(string))
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

func (c *Controller) syncItem(key string) (policyDeleted bool, retErr error) {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing policy", log.String("policy", key), log.Duration("processTime", time.Since(startTime)))
	}()
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return policyDeleted, err
	}

	defer func() {
		if policyDeleted && retErr == nil {
			log.Infof("Delete key '%s' from policy role cache", key)
			policyrolecache.Cache.DeletePolicy(key)
		}
	}()

	policy, err := c.policyLister.Policies(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Policy has been deleted. Attempting to cleanup resources",
				log.String("namespace", ns),
				log.String("name", name))
			policyDeleted = true
			return policyDeleted, c.updateRelatedRoles(key, policyDeleted)
		}
		log.Warn("Unable to retrieve policy from store",
			log.String("namespace", ns),
			log.String("name", name), log.Err(err))
		return policyDeleted, err
	}
	policy = policy.DeepCopy()
	policyDeleted = policy.DeletionTimestamp != nil
	return policyDeleted, c.updateRelatedRoles(key, policyDeleted)
}

func (c *Controller) updateRelatedRoles(policyName string, policyDeleted bool) error {
	roles := policyrolecache.Cache.GetRolesByPolicy(policyName)
	for roleName := range roles {
		roleNs, roleName, err := cache.SplitMetaNamespaceKey(roleName)
		if err != nil {
			return err
		}
		role, err := c.client.AuthzV1().Roles(roleNs).Get(context.Background(), roleName, metav1.GetOptions{ResourceVersion: "0"})
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			} else {
				log.Warn("Unable to retrieve role from store",
					log.String("namespace", roleNs),
					log.String("name", roleName), log.Err(err))
				return err
			}
		}
		annotations := role.Annotations
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[constant.UpdatedByPolicyController] = time.Now().Format("2006-01-02T15:04:05")
		role.Annotations = annotations
		if policyDeleted {
			role.Policies = removeItem(role.Policies, policyName)
		}
		_, err = c.client.AuthzV1().Roles(roleNs).Update(context.Background(), role, metav1.UpdateOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			log.Warn("Unable to update role",
				log.String("namespace", roleNs),
				log.String("name", roleName), log.Err(err))
			return err
		}
	}
	return nil
}

func removeItem(raw []string, item string) []string {
	var result []string
	for _, str := range raw {
		if str != item {
			result = append(result, str)
		}
	}
	return result
}
