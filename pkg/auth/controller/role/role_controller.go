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

package role

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
	v1 "tkestack.io/tke/api/auth/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/role/deletion"
	authutil "tkestack.io/tke/pkg/auth/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// roleDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	roleDeletionGracePeriod = 5 * time.Second

	controllerName = "role-controller"
)

// Controller is responsible for performing actions dependent upon a role phase.
type Controller struct {
	client           clientset.Interface
	queue            workqueue.RateLimitingInterface
	roleLister       authv1lister.RoleLister
	roleListerSynced cache.InformerSynced
	ruleLister       authv1lister.RuleLister
	ruleListerSynced cache.InformerSynced
	// helper to delete all resources in the role when the role is deleted.
	roleedResourcesDeleter deletion.RoledResourcesDeleterInterface
	enforcer               *casbin.SyncedEnforcer
}

// NewController creates a new role object.
func NewController(client clientset.Interface, roleInformer authv1informer.RoleInformer, ruleInformer authv1informer.RuleInformer, enforcer *casbin.SyncedEnforcer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                 client,
		queue:                  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		enforcer:               enforcer,
		roleedResourcesDeleter: deletion.NewRoleedResourcesDeleter(client.AuthV1().Roles(), client.AuthV1(), enforcer, finalizerToken, true),
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("role_controller", client.AuthV1().RESTClient().GetRateLimiter())
	}

	roleInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.Role)
				cur, ok2 := newObj.(*v1.Role)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					log.Info("Update enqueue")
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.roleLister = roleInformer.Lister()
	controller.roleListerSynced = roleInformer.Informer().HasSynced

	controller.ruleLister = ruleInformer.Lister()
	controller.ruleListerSynced = ruleInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.role, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, roleDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.Role, new *v1.Role) bool {
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
	log.Info("Starting role controller")
	defer log.Info("Shutting down role controller")

	if ok := cache.WaitForCacheSync(stopCh, c.roleListerSynced, c.ruleListerSynced); !ok {
		log.Error("Failed to wait for role caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of role objects.
// Each role can be in the queue at most once.
// The system ensures that no two workers can process
// the same role at the same time.
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

		// rather than wait for a full resync, re-add the role to the queue to be processed
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

// syncItem will sync the role with the given key if it has had
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing role", log.String("role", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	role, err := c.roleLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("role has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve role from store", log.String("role name", key), log.Err(err))
	default:
		if role.Status.Phase == v1.RoleTerminating {
			err = c.roleedResourcesDeleter.Delete(key)
		} else {
			err = c.processUpdate(role, key)
		}

		log.Debug("Handle role", log.Any("role", role))
	}
	return err
}

func (c *Controller) processUpdate(role *v1.Role, key string) error {

	// start update role if needed
	err := c.handlePhase(key, role)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) handlePhase(key string, role *v1.Role) error {

	var errs []error
	err := c.handleSpec(key, role)
	if err != nil {
		errs = append(errs, err)
	}

	err = c.handleSubjects(key, role)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) handleSpec(key string, role *v1.Role) error {
	rules, err := c.enforcer.GetRolesForUser(role.Name)
	if err != nil {
		return err
	}

	log.Debugf("Get roles for role: %s, %v", role.Name, rules)
	var existsPolicies []string
	for _, rule := range rules {
		if strings.HasPrefix(rule, "pol-") {
			existsPolicies = append(existsPolicies, rule)
		}
	}

	expectedPolicies := role.Spec.Policies

	added, removed := util.DiffStringSlice(existsPolicies, expectedPolicies)

	if len(added) != 0 || len(removed) != 0 {
		log.Info("Handle role added and removed", log.String("role", key), log.Any("added", added), log.Any("removed", removed))
	}

	var errs []error
	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUser(role.Name, add); err != nil {
				log.Errorf("Bind policy to role failed", log.String("role", role.Name), log.String("user", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUser(role.Name, remove); err != nil {
				log.Errorf("Unbind policy to role failed", log.String("group", role.Name), log.String("user", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) handleSubjects(key string, role *v1.Role) error {
	rules := c.enforcer.GetFilteredGroupingPolicy(1, role.Name)
	log.Debugf("Get grouping rules for role: %s, %v", role.Name, rules)
	var existUsers []string
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], authutil.UserPrefix(role.Spec.TenantID)) {
			existUsers = append(existUsers, strings.TrimPrefix(rule[0], authutil.UserPrefix(role.Spec.TenantID)))
		}
	}

	var expectedUsers []string
	for _, subj := range role.Status.Users {
		expectedUsers = append(expectedUsers, subj.Name)
	}

	var errs []error
	added, removed := util.DiffStringSlice(existUsers, expectedUsers)
	if len(added) != 0 || len(removed) != 0 {
		log.Info("Handle role users changed", log.String("role", key), log.Strings("added", added), log.Strings("removed", removed))
	}
	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUser(authutil.UserKey(role.Spec.TenantID, add), role.Name); err != nil {
				log.Errorf("Bind user to role failed", log.String("role", role.Name), log.String("user", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUser(authutil.UserKey(role.Spec.TenantID, remove), role.Name); err != nil {
				log.Errorf("Unbind user to role failed", log.String("role", role.Name), log.String("user", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	var existGroups []string
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], authutil.GroupPrefix(role.Spec.TenantID)) {
			existGroups = append(existGroups, strings.TrimPrefix(rule[0], authutil.GroupPrefix(role.Spec.TenantID)))
		}
	}

	var expectedGroups []string
	for _, subj := range role.Status.Groups {
		expectedGroups = append(expectedGroups, subj.ID)
	}

	added, removed = util.DiffStringSlice(existGroups, expectedGroups)
	if len(added) != 0 || len(removed) != 0 {
		log.Info("Handle role groups changed", log.String("role", key), log.Strings("added", added), log.Strings("removed", removed))
	}
	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUser(authutil.GroupKey(role.Spec.TenantID, add), role.Name); err != nil {
				log.Errorf("Bind groups to role failed", log.String("role", role.Name), log.String("group", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUser(authutil.GroupKey(role.Spec.TenantID, remove), role.Name); err != nil {
				log.Errorf("Unbind group to role failed", log.String("role", role.Name), log.String("group", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}
