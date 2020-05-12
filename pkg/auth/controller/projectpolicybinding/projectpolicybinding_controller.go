/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package projectpolicybinding

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	v1 "tkestack.io/tke/api/auth/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/projectpolicybinding/deletion"
	authutil "tkestack.io/tke/pkg/auth/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// projectPolicyDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	projectPolicyDeletionGracePeriod = 5 * time.Second

	controllerName = "projectpolicybinding-controller"
)

// Controller is responsible for performing actions dependent upon a policy phase.
type Controller struct {
	client             clientset.Interface
	queue              workqueue.RateLimitingInterface
	policyLister       authv1lister.ProjectPolicyBindingLister
	policyListerSynced cache.InformerSynced
	ruleLister         authv1lister.RuleLister
	ruleListerSynced   cache.InformerSynced
	// helper to delete all resources in the policy when the policy is deleted.
	projectpolicyedResourcesDeleter deletion.ProjectPolicyBindingResourcesDeleterInterface
	enforcer                        *casbin.SyncedEnforcer
}

// NewController creates a new projectpolicy controller object.
func NewController(client clientset.Interface, policyInformer authv1informer.ProjectPolicyBindingInformer, ruleInformer authv1informer.RuleInformer, enforcer *casbin.SyncedEnforcer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                          client,
		queue:                           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		enforcer:                        enforcer,
		projectpolicyedResourcesDeleter: deletion.NewProjectPolicyBindingResourcesDeleter(client.AuthV1().ProjectPolicyBindings(), client.AuthV1(), enforcer, finalizerToken, true),
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.AuthV1().RESTClient().GetRateLimiter())
	}

	policyInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.ProjectPolicyBinding)
				cur, ok2 := newObj.(*v1.ProjectPolicyBinding)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					log.Info("Update enqueue", log.String("project policy binding", cur.Name))
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.policyLister = policyInformer.Lister()
	controller.policyListerSynced = policyInformer.Informer().HasSynced

	controller.ruleLister = ruleInformer.Lister()
	controller.ruleListerSynced = ruleInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.policy, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, projectPolicyDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.ProjectPolicyBinding, new *v1.ProjectPolicyBinding) bool {
	if old.UID != new.UID {
		return true
	}

	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true
	}

	if !reflect.DeepEqual(old.Status, new.Status) {
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
	log.Info("Starting policy controller")
	defer log.Info("Shutting down policy controller")

	if ok := cache.WaitForCacheSync(stopCh, c.policyListerSynced, c.ruleListerSynced); !ok {
		log.Error("Failed to wait for policy caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of policy objects.
// Each policy can be in the queue at most once.
// The system ensures that no two workers can process
// the same policy at the same time.
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

		// rather than wait for a full resync, re-add the policy to the queue to be processed
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

// syncItem will sync the policy with the given key if it has had
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing projectPolicy", log.String("projectPolicy", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	projectPolicy, err := c.policyLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("projectPolicy has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve projectPolicy from store", log.String("projectPolicy name", key), log.Err(err))
	default:
		if projectPolicy.Status.Phase == v1.BindingTerminating {
			log.Info("Delete project policy", log.String("key", key))
			err = c.projectpolicyedResourcesDeleter.Delete(key)
		} else {
			err = c.processUpdate(projectPolicy, key)
		}

		log.Debug("Handle projectPolicy", log.Any("projectPolicy", projectPolicy))
	}
	return err
}

func (c *Controller) processUpdate(policy *v1.ProjectPolicyBinding, key string) error {

	// start update policy if needed
	err := c.handlePhase(key, policy)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) handlePhase(key string, policy *v1.ProjectPolicyBinding) error {

	var errs []error

	err := c.handleSubjects(key, policy)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) handleSubjects(key string, policy *v1.ProjectPolicyBinding) error {
	if policy.Spec.PolicyID == "" || policy.Spec.ProjectID == "" {
		log.Info("PolicyID or projectID is empty for projectPolicy", log.String("projectpolicy", policy.Name))
		return nil
	}

	expectProjectSubj := sets.String{}
	for _, subj := range policy.Spec.Users {
		expectProjectSubj.Insert(authutil.UserKey(policy.Spec.TenantID, subj.Name))
	}

	for _, subj := range policy.Spec.Groups {
		expectProjectSubj.Insert(authutil.GroupKey(policy.Spec.TenantID, subj.Name))
	}

	return c.handleRules(policy.Spec.TenantID, policy.Spec.PolicyID, policy.Spec.ProjectID, expectProjectSubj.UnsortedList())
}

func (c *Controller) handleRules(tenantID, policyID, projectID string, expectSubj []string) error {
	rules := c.enforcer.GetFilteredGroupingPolicy(1, policyID, projectID)

	log.Debugf("Get grouping rules for policy: %s, %v", policyID, rules)

	var existSubj []string
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], authutil.UserPrefix(tenantID)) || strings.HasPrefix(rule[0], authutil.GroupPrefix(tenantID)) {
			existSubj = append(existSubj, rule[0])
		}
	}

	var errs []error
	added, removed := util.DiffStringSlice(existSubj, expectSubj)
	if len(added) != 0 || len(removed) != 0 {
		log.Info("Handle project policy subj changed", log.String("policy", policyID), log.String("project", projectID), log.Strings("added", added), log.Strings("removed", removed))
	}

	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUserInDomain(add, policyID, projectID); err != nil {
				log.Errorf("Bind subj to policy failed", log.String("policy", policyID), log.String("project", projectID), log.String("subj", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUserInDomain(remove, policyID, projectID); err != nil {
				log.Errorf("Unbind subj to policy failed", log.String("policy", policyID), log.String("project", projectID), log.String("subj", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}
