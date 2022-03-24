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

package multiclusterrolebinding

import (
	"context"
	"encoding/json"
	"fmt"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
	apiauthzv1 "tkestack.io/tke/api/authz/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	authzv1informer "tkestack.io/tke/api/client/informers/externalversions/authz/v1"
	authzv1 "tkestack.io/tke/api/client/listers/authz/v1"
	apiplatformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/authz/constant"
	"tkestack.io/tke/pkg/authz/controller/multiclusterrolebinding/deletion"
	authzprovider "tkestack.io/tke/pkg/authz/provider"
	controllerutil "tkestack.io/tke/pkg/controller"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// appDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	appDeletionGracePeriod = 5 * time.Second
)

const (
	controllerName = "multiclusterrolebinding-controller"
)

type Controller struct {
	client         clientset.Interface
	platformClient platformversionedclient.PlatformV1Interface
	queue          workqueue.RateLimitingInterface
	policyLister   authzv1.PolicyLister
	policySynced   cache.InformerSynced
	roleLister     authzv1.RoleLister
	roleSynced     cache.InformerSynced
	mcrbLister     authzv1.MultiClusterRoleBindingLister
	mcrbSynced     cache.InformerSynced
	mcrbDeleter    deletion.MultiClusterRoleBindingDeleter
	stopCh         <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	policyInformer authzv1informer.PolicyInformer,
	roleInformer authzv1informer.RoleInformer,
	mcrbInformer authzv1informer.MultiClusterRoleBindingInformer,
	resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:         client,
		platformClient: platformClient,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		mcrbDeleter:    deletion.New(client, platformClient),
	}
	if client != nil &&
		client.AuthzV1().RESTClient() != nil &&
		!reflect.ValueOf(client.AuthzV1().RESTClient()).IsNil() &&
		client.AuthzV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.AuthzV1().RESTClient().GetRateLimiter())
	}

	mcrbInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: controller.enqueue,
				UpdateFunc: func(oldObj, newObj interface{}) {
					old, ok1 := oldObj.(*apiauthzv1.MultiClusterRoleBinding)
					cur, ok2 := newObj.(*apiauthzv1.MultiClusterRoleBinding)
					if ok1 && ok2 {
						if cur.Labels[constant.DispatchAllClusters] == "true" {
							controller.enqueue(newObj)
						} else if controller.needsUpdate(old, cur) {
							controller.enqueue(newObj)
						}
					}
				},
				DeleteFunc: controller.enqueue,
			},
			FilterFunc: func(obj interface{}) bool {
				mcrb, ok := obj.(*apiauthzv1.MultiClusterRoleBinding)
				if !ok {
					return false
				}
				provider, err := authzprovider.GetProvider(mcrb.Annotations)
				if err != nil {
					return true
				}
				return provider.OnFilter(context.TODO(), mcrb.Annotations)
			},
		},
		resyncPeriod,
	)

	controller.policyLister = policyInformer.Lister()
	controller.policySynced = policyInformer.Informer().HasSynced
	controller.roleLister = roleInformer.Lister()
	controller.roleSynced = roleInformer.Informer().HasSynced
	controller.mcrbLister = mcrbInformer.Lister()
	controller.mcrbSynced = mcrbInformer.Informer().HasSynced
	return controller
}

func (c *Controller) enqueueCluster(obj interface{}) {
	cluster := obj.(*apiplatformv1.Cluster)
	selector, _ := labels.Parse(fmt.Sprintf("%s=%s", constant.DispatchAllClusters, "true"))
	list, err := c.mcrbLister.MultiClusterRoleBindings(cluster.Spec.TenantID).List(selector)
	if err != nil {
		log.Warnf("failed to list mcrbs for tenant '%s'", cluster.Spec.TenantID)
		return
	}
	for _, item := range list {
		c.enqueue(item)
	}
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, appDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *apiauthzv1.MultiClusterRoleBinding, new *apiauthzv1.MultiClusterRoleBinding) bool {
	if old.UID != new.UID {
		return true
	}
	if !reflect.DeepEqual(old.Annotations, new.Annotations) {
		return true
	}
	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true
	}
	if !reflect.DeepEqual(old.Status, new.Status) {
		return true
	}
	if !reflect.DeepEqual(old.DeletionTimestamp, new.DeletionTimestamp) {
		return true
	}
	return false
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting app controller")
	defer log.Info("Shutting down app controller")

	if ok := cache.WaitForCacheSync(stopCh, c.mcrbSynced, c.policySynced); !ok {
		log.Error("Failed to wait for app caches to sync")
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
		log.Info("Finished syncing MultiClusterRoleBinding", log.String("MultiClusterRoleBinding", key), log.Duration("processTime", time.Since(startTime)))
	}()
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	mcrb, err := c.mcrbLister.MultiClusterRoleBindings(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("MultiClusterRoleBinding has been deleted. Attempting to cleanup resources",
				log.String("namespace", ns),
				log.String("name", name))
			return nil
		}
		log.Warn("Unable to retrieve MultiClusterRoleBinding from store",
			log.String("namespace", ns),
			log.String("name", name), log.Err(err))
		return err
	}
	mcrb = mcrb.DeepCopy()
	provider, err := authzprovider.GetProvider(mcrb.Annotations)
	if err != nil {
		log.Warn("Unable to retrieve provider",
			log.String("namespace", ns),
			log.String("name", name), log.Err(err))
		return err
	}
	ctx := provider.InitContext(mcrb)
	if mcrb.Labels[constant.DispatchAllClusters] == "true" {
		mcrb.Spec.Clusters, err = provider.GetTenantClusters(ctx, c.platformClient, mcrb.Namespace)
		if err != nil {
			log.Warnf("failed to get tenant clusters, err '%v'", err)
			return err
		}
	}

	switch mcrb.Status.Phase {
	case apiauthzv1.BindingActive:
		return c.handleActive(ctx, mcrb, provider)
	case apiauthzv1.BindingTerminating:
		return c.mcrbDeleter.Delete(ctx, mcrb, provider)
	default:
		return fmt.Errorf("unknown MultiClusterRoleBinding phase '%s'", mcrb.Status.Phase)
	}
}

func (c *Controller) handleActive(ctx context.Context, mcrb *apiauthzv1.MultiClusterRoleBinding, provider authzprovider.Provider) error {
	roleNs, roleName, err := cache.SplitMetaNamespaceKey(mcrb.Spec.RoleName)
	if err != nil {
		log.Warnf("failed to parse Role namespace/name '%s'", mcrb.Spec.RoleName)
		return err
	}

	role, err := c.roleLister.Roles(roleNs).Get(roleName)
	if err != nil {
		log.Warn("Unable to retrieve Role from store",
			log.String("namespace", roleNs),
			log.String("name", roleName), log.Err(err))
		return err
	}

	// 将Role关联的多个Policy合并
	policies, err := c.combineRolePolicies(role)
	if err != nil {
		log.Warn("Unable to combine role policies",
			log.String("namespace", roleNs),
			log.String("name", roleName), log.Err(err))
		return err
	}

	// 获取user在各个cluster内的subject
	clusterSubjects := map[string]*rbacv1.Subject{}
	for _, cls := range mcrb.Spec.Clusters {
		cluster, err := clusterprovider.GetV1ClusterByName(ctx, c.platformClient, cls, mcrb.Spec.Username)
		if err != nil {
			log.Warnf("GetV1ClusterByName failed, cluster: '%s', user: '%s', err: '%#v'", cls, mcrb.Spec.Username, err)
			return err
		}
		subject, err := provider.GetSubject(ctx, mcrb.Spec.Username, cluster)
		if err != nil {
			log.Warnf("GetSubject failed, cluster: '%s',  user: '%s', err: '%#v'", cls, mcrb.Spec.Username, err)
			return err
		}
		clusterSubjects[cls] = subject
	}

	// 执行权限分发
	if err = provider.DispatchMultiClusterRoleBinding(ctx, c.platformClient, mcrb, policies, clusterSubjects); err != nil {
		log.Warnf("DispatchMultiClusterRoleBinding failed, MultiClusterRoleBinding: '%s', err: '%#v'", mcrb.Name, err)
		return err
	}

	// 删除已经解绑的资源
	lastDispatchedClusters := []string{}
	if lastStr, ok := mcrb.Annotations[constant.LastDispatchedClusters]; ok {
		err = json.Unmarshal([]byte(lastStr), &lastDispatchedClusters)
		if err != nil {
			log.Warnf("Unmarshal lastDispatchedClusters failed', err: '%#v'", err)
			return err
		}
	}
	oldSet := sets.NewString(lastDispatchedClusters...)
	newSet := sets.NewString(mcrb.Spec.Clusters...)
	difference := oldSet.Difference(newSet)
	if len(difference) != 0 {
		if err = provider.DeleteUnbindingResources(ctx, c.platformClient, mcrb, difference.List()); err != nil {
			log.Warnf("DeleteUnbindingResources '%s/%s' failed', err: '%#v'", mcrb.Namespace, mcrb.Name, err)
			return err
		}
		clsBytes, _ := json.Marshal(mcrb.Spec.Clusters)
		mcrb.Annotations[constant.LastDispatchedClusters] = string(clsBytes)
		if mcrb.Labels[constant.DispatchAllClusters] == "true" {
			mcrb.Spec.Clusters = []string{"*"}
		}
		if _, err = c.client.AuthzV1().MultiClusterRoleBindings(mcrb.Namespace).Update(context.Background(), mcrb, metav1.UpdateOptions{}); err != nil {
			log.Warnf("Update MultiClusterRoleBindings '%s/%s' failed', err: '%#v'", mcrb.Namespace, mcrb.Name, err)
			return err
		}
	}
	return nil
}

func (c *Controller) combineRolePolicies(role *apiauthzv1.Role) ([]rbacv1.PolicyRule, error) {
	var policyRules []rbacv1.PolicyRule
	for _, policy := range role.Policies {
		policyNamespace, policyName, _ := cache.SplitMetaNamespaceKey(policy)
		pol, err := c.policyLister.Policies(policyNamespace).Get(policyName)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Warnf("Policy '%s/%s' is not exist", policyNamespace, policyName)
				continue
			}
			log.Warnf("Unable get policy '%s/%s', err: '%v'", policyNamespace, policyName, err)
			return nil, err
		}
		policyRules = append(policyRules, pol.Rules...)
	}
	return policyRules, nil
}
