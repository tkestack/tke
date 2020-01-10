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

package group

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
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
	"tkestack.io/tke/pkg/auth/controller/group/deletion"
	authutil "tkestack.io/tke/pkg/auth/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// groupDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	groupDeletionGracePeriod = 5 * time.Second

	controllerName = "group-controller"

	groupSyncedPeriod = 1 * time.Minute
)

// Controller is responsible for performing actions dependent upon a group phase.
type Controller struct {
	client            clientset.Interface
	queue             workqueue.RateLimitingInterface
	groupLister       authv1lister.LocalGroupLister
	groupListerSynced cache.InformerSynced
	ruleLister        authv1lister.RuleLister
	ruleListerSynced  cache.InformerSynced
	// helper to delete all resources in the group when the group is deleted.
	groupedResourcesDeleter deletion.GroupedResourcesDeleterInterface
	enforcer                *casbin.SyncedEnforcer
}

// NewController creates a new group object.
func NewController(client clientset.Interface, groupInformer authv1informer.LocalGroupInformer, ruleInformer authv1informer.RuleInformer, enforcer *casbin.SyncedEnforcer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                  client,
		queue:                   workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		enforcer:                enforcer,
		groupedResourcesDeleter: deletion.NewGroupedResourcesDeleter(client.AuthV1().LocalGroups(), client.AuthV1(), enforcer, finalizerToken, true),
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("group_controller", client.AuthV1().RESTClient().GetRateLimiter())
	}

	groupInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.LocalGroup)
				cur, ok2 := newObj.(*v1.LocalGroup)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					log.Info("Update enqueue")
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.groupLister = groupInformer.Lister()
	controller.groupListerSynced = groupInformer.Informer().HasSynced

	controller.ruleLister = ruleInformer.Lister()
	controller.ruleListerSynced = ruleInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.group, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, groupDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.LocalGroup, new *v1.LocalGroup) bool {
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
	log.Info("Starting group controller")
	defer log.Info("Shutting down group controller")

	if ok := cache.WaitForCacheSync(stopCh, c.groupListerSynced, c.ruleListerSynced); !ok {
		log.Error("Failed to wait for group caches to sync")
	}

	// sync groups(include 3rd party) for identity providers into casbin
	go c.pollThirdPartyGroup(stopCh)

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of group objects.
// Each group can be in the queue at most once.
// The system ensures that no two workers can process
// the same group at the same time.
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

		// rather than wait for a full resync, re-add the group to the queue to be processed
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

// syncItem will sync the group with the given key if it has had
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing group", log.String("group", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	group, err := c.groupLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("group has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve group from store", log.String("group name", key), log.Err(err))
	default:
		if group.Status.Phase == v1.GroupTerminating {
			err = c.groupedResourcesDeleter.Delete(key)
		} else {
			err = c.processUpdate(group, key)
		}

		log.Debug("Handle group", log.Any("group", group))
	}
	return err
}

func (c *Controller) processUpdate(group *v1.LocalGroup, key string) error {
	return c.handleSubjects(key, convertToGroup(group))
}

func convertToGroup(localGroup *v1.LocalGroup) *v1.Group {
	return &v1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: localGroup.ObjectMeta.Name,
		},
		Spec: v1.GroupSpec{
			ID:          localGroup.ObjectMeta.Name,
			DisplayName: localGroup.Spec.DisplayName,
			TenantID:    localGroup.Spec.TenantID,
			Description: localGroup.Spec.TenantID,
		},
		Status: v1.GroupStatus{
			Users: localGroup.Status.Users,
		},
	}
}

func (c *Controller) handleSubjects(key string, group *v1.Group) error {
	rules := c.enforcer.GetFilteredGroupingPolicy(1, authutil.GroupKey(group.Spec.TenantID, key))
	log.Debugf("Get grouping rules for group: %s, %v", group.Name, rules)
	var existMembers []string
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], authutil.UserPrefix(group.Spec.TenantID)) {
			existMembers = append(existMembers, strings.TrimPrefix(rule[0], authutil.UserPrefix(group.Spec.TenantID)))
		}
	}

	var expectedMembers []string
	for _, subj := range group.Status.Users {
		expectedMembers = append(expectedMembers, subj.Name)
	}

	var errs []error
	added, removed := util.DiffStringSlice(existMembers, expectedMembers)
	if len(added) != 0 || len(removed) != 0 {
		log.Info("Handle group subjects changed", log.String("group", key), log.Strings("added", added), log.Strings("removed", removed))
	}
	if len(added) > 0 {
		for _, add := range added {
			if _, err := c.enforcer.AddRoleForUser(authutil.UserKey(group.Spec.TenantID, add), authutil.GroupKey(group.Spec.TenantID, group.Name)); err != nil {
				log.Errorf("Bind group to user failed", log.String("group", group.Name), log.String("user", add), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	if len(removed) > 0 {
		for _, remove := range removed {
			if _, err := c.enforcer.DeleteRoleForUser(authutil.UserKey(group.Spec.TenantID, remove), authutil.GroupKey(group.Spec.TenantID, group.Name)); err != nil {
				log.Errorf("Unbind group to user failed", log.String("group", group.Name), log.String("user", remove), log.Err(err))
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

// pollThirdPartyGroup syncs groups with members into storage
func (c *Controller) pollThirdPartyGroup(stopCh <-chan struct{}) {
	timerC := time.NewTicker(groupSyncedPeriod)
	for {
		select {
		case <-timerC.C:
			c.resyncGroups()
		case <-stopCh:
			timerC.Stop()
			return
		}
	}
}

func (c *Controller) resyncGroups() {
	defer log.Info("Finished syncing groups with users")

	idpList, err := c.client.AuthV1().IdentityProviders().List(metav1.ListOptions{})
	if err != nil {
		log.Error("List all identity providers failed", log.Err(err))
		return
	}

	for _, idp := range idpList.Items {
		tenantSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("spec.tenantID", idp.Name),
			fields.OneTermEqualSelector(auth.QueryLimitTag, "0"),
		)

		groups, err := c.client.AuthV1().Groups().List(metav1.ListOptions{FieldSelector: tenantSelector.String()})
		if err != nil {
			log.Error("List groups for tenant failed", log.String("tenant", idp.Name), log.Err(err))
			continue
		}
		log.Debug("syncing groups for tenantID", log.String("tenant", idp.Name), log.Any("groups", groups))
		for _, grp := range groups.Items {
			_ = c.handleSubjects(grp.Name, grp.DeepCopy())
		}
	}

}
