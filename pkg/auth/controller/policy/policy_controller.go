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

package policy

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"tkestack.io/tke/api/auth"
	v1 "tkestack.io/tke/api/auth/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/policy/deletion"
	authutil "tkestack.io/tke/pkg/auth/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// policyDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	policyDeletionGracePeriod = 5 * time.Second

	controllerName = "policy-controller"
)

// Controller is responsible for performing actions dependent upon a policy phase.
type Controller struct {
	client             clientset.Interface
	queue              workqueue.RateLimitingInterface
	policyLister       authv1lister.PolicyLister
	policyListerSynced cache.InformerSynced
	ruleLister         authv1lister.RuleLister
	ruleListerSynced   cache.InformerSynced
	// helper to delete all resources in the policy when the policy is deleted.
	policyedResourcesDeleter deletion.PoliciedResourcesDeleterInterface
	enforcer                 *casbin.SyncedEnforcer
}

// NewController creates a new policy object.
func NewController(client clientset.Interface, policyInformer authv1informer.PolicyInformer, ruleInformer authv1informer.RuleInformer, enforcer *casbin.SyncedEnforcer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                   client,
		queue:                    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		enforcer:                 enforcer,
		policyedResourcesDeleter: deletion.NewPoliciedResourcesDeleter(client.AuthV1().Policies(), client.AuthV1(), enforcer, finalizerToken, true),
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("policy_controller", client.AuthV1().RESTClient().GetRateLimiter())
	}

	policyInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.Policy)
				cur, ok2 := newObj.(*v1.Policy)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					log.Info("Update enqueue")
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
	c.queue.AddAfter(key, policyDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.Policy, new *v1.Policy) bool {
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
		log.Info("Finished syncing policy", log.String("policy", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	policy, err := c.policyLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("policy has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve policy from store", log.String("policy name", key), log.Err(err))
	default:
		if policy.Status.Phase == v1.PolicyActive {
			err = c.processUpdate(policy, key)
		} else if policy.Status.Phase == v1.PolicyTerminating {
			err = c.policyedResourcesDeleter.Delete(key)
		}

		//log.Info("Handle policy", log.Any("policy", policy))
	}
	return err
}

func (c *Controller) processUpdate(policy *v1.Policy, key string) error {

	// start update policy if needed
	err := c.handlePhase(key, policy)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) handlePhase(key string, policy *v1.Policy) error {

	var errs []error
	err := c.handleSpec(key, policy)
	if err != nil {
		errs = append(errs, err)
	}

	err = c.handleSubjects(key, policy)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) handleSpec(key string, policy *v1.Policy) error {
	existedRule := c.enforcer.GetFilteredPolicy(0, key)

	var outPolicy = &auth.Policy{}
	err := v1.Convert_v1_Policy_To_auth_Policy(policy, outPolicy, nil)
	if err != nil {
		log.Error("unable to convert policy object: %v", log.Err(err))
		return err
	}

	expectedRule := authutil.ConvertPolicyToRuleArray(outPolicy)
	added, removed := util.Diff2DStringSlice(existedRule, expectedRule)

	log.Info("Handle policy added and removed", log.String("policy", key), log.Any("added", added), log.Any("removed", removed))
	var errs []error
	if len(added) != 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddPolicy(add); err != nil {
				log.Errorf("Add policy failed", log.Strings("rule", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) != 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.RemovePolicy(remove); err != nil {
				log.Errorf("Remove policy failed", log.Strings("rule", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) handleSubjects(key string, policy *v1.Policy) error {
	rules := c.enforcer.GetFilteredGroupingPolicy(1, policy.Name)
	log.Debugf("Get grouping rules for policy: %s, %v", policy.Name, rules)
	var existSubj []string
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], authutil.UserPrefix(policy.Spec.TenantID)) {
			existSubj = append(existSubj, strings.TrimPrefix(rule[0], authutil.UserPrefix(policy.Spec.TenantID)))
		}
	}

	var expectedSubj []string
	for _, subj := range policy.Status.Subjects {
		expectedSubj = append(expectedSubj, subj.Name)
	}

	var errs []error
	added, removed := util.DiffStringSlice(existSubj, expectedSubj)
	log.Info("Handle policy subjects changed", log.String("policy", key), log.Strings("added", added), log.Strings("removed", removed))
	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUser(authutil.UserKey(policy.Spec.TenantID, add), policy.Name); err != nil {
				log.Errorf("Bind policy to user failed", log.String("policy", policy.Name), log.String("user", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUser(authutil.UserKey(policy.Spec.TenantID, remove), policy.Name); err != nil {
				log.Errorf("Unbind policy to user failed", log.String("policy", policy.Name), log.String("user", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}
