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

package chartgroup

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	businessv1 "tkestack.io/tke/api/business/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	registryv1informer "tkestack.io/tke/api/client/informers/externalversions/registry/v1"
	registryv1lister "tkestack.io/tke/api/client/listers/registry/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/registry/controller/chartgroup/deletion"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// chartGroupDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	chartGroupDeletionGracePeriod = 5 * time.Second
)

const (
	controllerName = "chartGroup-controller"
)

// Controller is responsible for performing actions dependent upon an chartGroup phase.
type Controller struct {
	client         clientset.Interface
	businessClient businessversionedclient.BusinessV1Interface
	cache          *chartGroupCache
	health         *chartGroupHealth
	queue          workqueue.RateLimitingInterface
	lister         registryv1lister.ChartGroupLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// helper to delete all resources in the chartGroup when the chartGroup is deleted.
	chartGroupResourcesDeleter deletion.ChartGroupResourcesDeleterInterface
}

// NewController creates a new Controller object.
func NewController(businessClient businessversionedclient.BusinessV1Interface,
	client clientset.Interface, chartGroupInformer registryv1informer.ChartGroupInformer,
	resyncPeriod time.Duration, finalizerToken registryv1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                     client,
		businessClient:             businessClient,
		cache:                      &chartGroupCache{m: make(map[string]*cachedChartGroup)},
		health:                     &chartGroupHealth{chartGroups: sets.NewString()},
		queue:                      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		chartGroupResourcesDeleter: deletion.NewChartGroupResourcesDeleter(businessClient, client.RegistryV1(), finalizerToken, true),
	}

	if client != nil &&
		client.RegistryV1().RESTClient() != nil &&
		!reflect.ValueOf(client.RegistryV1().RESTClient()).IsNil() &&
		client.RegistryV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("chartGroup_controller", client.RegistryV1().RESTClient().GetRateLimiter())
	}

	chartGroupInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*registryv1.ChartGroup)
				cur, ok2 := newObj.(*registryv1.ChartGroup)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.lister = chartGroupInformer.Lister()
	controller.listerSynced = chartGroupInformer.Informer().HasSynced

	return controller
}

// obj could be an *registryv1.ChartGroup, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, chartGroupDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *registryv1.ChartGroup, new *registryv1.ChartGroup) bool {
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
	log.Info("Starting chartGroup controller")
	defer log.Info("Shutting down chartGroup controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for chartGroup caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of chartGroup objects.
// Each chartGroup can be in the queue at most once.
// The system ensures that no two workers can process
// the same chartGroup at the same time.
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

		// rather than wait for a full resync, re-add the chartGroup to the queue to be processed
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

// syncItem will sync the ChartGroup with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// chartGroups created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing chartGroup", log.String("chartGroup", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, chartGroupName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// chartGroup holds the latest ChartGroup info from apiserver
	chartGroup, err := c.lister.Get(chartGroupName)
	switch {
	case errors.IsNotFound(err):
		log.Info("ChartGroup has been deleted. Attempting to cleanup resources",
			log.String("chartGroupName", chartGroupName))
		_ = c.processDeletion(key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve chartGroup from store",
			log.String("chartGroupName", chartGroupName), log.Err(err))
		return err
	default:
		cachedChartGroup := c.cache.getOrCreate(key)
		if c.needCompatibleUpgrade(context.Background(), chartGroup) {
			return c.compatibleUpgrade(context.Background(), key, cachedChartGroup, chartGroup)
		}

		if chartGroup.Status.Phase == registryv1.ChartGroupPending ||
			chartGroup.Status.Phase == registryv1.ChartGroupAvailable {
			err = c.processUpdate(context.Background(), key, cachedChartGroup, chartGroup)
		} else if chartGroup.Status.Phase == registryv1.ChartGroupTerminating {
			log.Info("ChartGroup has been terminated. Attempting to cleanup resources",
				log.String("chartGroupName", chartGroupName))
			_ = c.processDeletion(key)
			err = c.chartGroupResourcesDeleter.Delete(context.Background(), chartGroupName)
			// If err is not nil, do not update object status when phase is Terminating.
			// DeletionTimestamp is not empty and object will be deleted when you request updateStatus
		} else {
			log.Debug(fmt.Sprintf("ChartGroup %s status is %s, not to process", key, chartGroup.Status.Phase))
		}
	}
	return err
}

func (c *Controller) processDeletion(key string) error {
	cachedChartGroup, ok := c.cache.get(key)
	if !ok {
		log.Debug("ChartGroup not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processDelete(cachedChartGroup, key)
}

func (c *Controller) processDelete(cachedChartGroup *cachedChartGroup, key string) error {
	log.Info("ChartGroup will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the chartGroup cache", log.String("name", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the chartGroup health cache", log.String("name", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) processUpdate(ctx context.Context, key string, cachedChartGroup *cachedChartGroup, chartGroup *registryv1.ChartGroup) error {
	if cachedChartGroup.state != nil {
		// exist and the chartGroup name changed
		if cachedChartGroup.state.UID != chartGroup.UID {
			if err := c.processDelete(cachedChartGroup, key); err != nil {
				return err
			}
		}
	}
	// start update chartGroup if needed
	updated, err := c.handlePhase(ctx, key, cachedChartGroup, chartGroup)
	if err != nil {
		return err
	}
	cachedChartGroup.state = updated
	// Always update the cache upon success.
	c.cache.set(key, cachedChartGroup)
	return nil
}

func (c *Controller) handlePhase(ctx context.Context, key string, cachedChartGroup *cachedChartGroup, chartGroup *registryv1.ChartGroup) (*registryv1.ChartGroup, error) {
	newStatus := chartGroup.Status.DeepCopy()
	switch chartGroup.Status.Phase {
	case registryv1.ChartGroupPending:
		err := c.createOrDeleteBusinessChartGroup(ctx, cachedChartGroup, chartGroup)
		if err != nil {
			newStatus.Phase = registryv1.ChartGroupFailed
			newStatus.Message = "createOrDeleteBusinessChartGroup failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			return c.updateStatus(ctx, chartGroup, &chartGroup.Status, newStatus)
		}
		newStatus.Phase = registryv1.ChartGroupAvailable
		newStatus.Message = ""
		newStatus.Reason = ""
		newStatus.LastTransitionTime = metav1.Now()
		return c.updateStatus(ctx, chartGroup, &chartGroup.Status, newStatus)
	case registryv1.ChartGroupAvailable:
		c.startChartGroupHealthCheck(ctx, key)

		err := c.createOrDeleteBusinessChartGroup(ctx, cachedChartGroup, chartGroup)
		if err != nil {
			newStatus.Phase = registryv1.ChartGroupFailed
			newStatus.Message = "createOrDeleteBusinessChartGroup failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			return c.updateStatus(ctx, chartGroup, &chartGroup.Status, newStatus)
		}
	}
	return chartGroup, nil
}

func (c *Controller) createOrDeleteBusinessChartGroup(ctx context.Context, cachedChartGroup *cachedChartGroup, chartGroup *registryv1.ChartGroup) error {
	if c.businessClient == nil {
		return nil
	}

	cachedProjects := []string{}
	if cachedChartGroup.state != nil {
		cachedProjects = cachedChartGroup.state.Spec.Projects
	}
	added, removed := util.DiffStringSlice(cachedProjects, chartGroup.Spec.Projects)
	var errs []error
	for _, projectID := range added {
		_, err := c.businessClient.ChartGroups(projectID).Get(ctx, chartGroup.Spec.Name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			_, err = c.businessClient.ChartGroups(projectID).Create(ctx, &businessv1.ChartGroup{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: projectID,
					Name:      chartGroup.Spec.Name,
				},
				Spec: businessv1.ChartGroupSpec{
					Name:     chartGroup.Spec.Name,
					TenantID: chartGroup.Spec.TenantID,
				}}, metav1.CreateOptions{})
			if err != nil {
				log.Warn("ChartGroup controller - addBusinessChartGroup failed",
					log.String("projectID", projectID),
					log.String("chartGroupName", chartGroup.Spec.Name), log.Err(err))
				errs = append(errs, err)
			} else {
				log.Info("ChartGroup controller - addBusinessChartGroup",
					log.String("projectID", projectID),
					log.String("chartGroupName", chartGroup.Spec.Name))
			}
		}
	}
	for _, projectID := range removed {
		err := c.businessClient.ChartGroups(projectID).Delete(ctx, chartGroup.Spec.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Warn("ChartGroup controller - deleteBusinessChartGroup failed",
				log.String("projectID", projectID),
				log.String("chartGroupName", chartGroup.Spec.Name), log.Err(err))
			if !errors.IsNotFound(err) {
				errs = append(errs, err)
			}
		} else {
			log.Info("ChartGroup controller - deleteBusinessChartGroup",
				log.String("projectID", projectID),
				log.String("chartGroupName", chartGroup.Spec.Name))
		}
	}
	return utilerrors.NewAggregate(errs)
}

func (c *Controller) updateStatus(ctx context.Context, chartGroup *registryv1.ChartGroup, previousStatus, newStatus *registryv1.ChartGroupStatus) (*registryv1.ChartGroup, error) {
	if reflect.DeepEqual(previousStatus, newStatus) {
		return nil, nil
	}
	// Make a copy so we don't mutate the shared informer cache.
	newObj := chartGroup.DeepCopy()
	newObj.Status = *newStatus

	updated, err := c.client.RegistryV1().ChartGroups().UpdateStatus(ctx, newObj, metav1.UpdateOptions{})
	if err == nil {
		return updated, nil
	}
	if errors.IsNotFound(err) {
		log.Info("Not persisting update to chartGroup that no longer exists",
			log.String("chartGroupName", chartGroup.Name),
			log.Err(err))
		return updated, nil
	}
	if errors.IsConflict(err) {
		return nil, fmt.Errorf("not persisting update to chartGroup '%s' that has been changed since we received it: %v",
			chartGroup.Name, err)
	}
	log.Warn(fmt.Sprintf("Failed to persist updated status of chartGroup '%s/%s'",
		chartGroup.Name, chartGroup.Status.Phase),
		log.String("chartGroupName", chartGroup.Name), log.Err(err))
	return nil, err
}

// If need to upgrade compatibly
func (c *Controller) needCompatibleUpgrade(ctx context.Context, cg *registryv1.ChartGroup) bool {
	switch cg.Spec.Type {
	// lower case
	case registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypePersonal))),
		registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypeProject))),
		registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypeSystem))):
		{
			return true
		}
	}
	return false
}

// If need to upgrade compatibly
func (c *Controller) compatibleUpgrade(ctx context.Context, key string, cachedChartGroup *cachedChartGroup, cg *registryv1.ChartGroup) error {
	newObj := cg.DeepCopy()
	switch newObj.Spec.Type {
	case registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypePersonal))):
		{
			newObj.Spec.Creator = cg.Spec.Name
			newObj.Spec.Type = registryv1.RepoTypeSelfBuilt
			if cg.Spec.Visibility == registryv1.VisibilityPrivate {
				newObj.Spec.Visibility = registryv1.VisibilityUser
				newObj.Spec.Users = []string{cg.Spec.Name}
			}
			break
		}
	case registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypeProject))):
		{
			newObj.Spec.Type = registryv1.RepoTypeSelfBuilt
			if cg.Spec.Visibility == registryv1.VisibilityPrivate {
				newObj.Spec.Visibility = registryv1.VisibilityProject
			}
			break
		}
	case registryv1.RepoType(strings.ToLower(string(registryv1.RepoTypeSystem))):
		{
			newObj.Spec.Type = registryv1.RepoTypeSystem
			break
		}
	default:
		return nil
	}

	updated, err := c.client.RegistryV1().ChartGroups().Update(ctx, newObj, metav1.UpdateOptions{})
	if errors.IsNotFound(err) {
		log.Info("Not persisting upgrade to chartGroup that no longer exists",
			log.String("chartGroupName", cg.Name),
			log.Err(err))
		return nil
	}

	cachedChartGroup.state = updated
	// Always update the cache upon success.
	c.cache.set(key, cachedChartGroup)
	return nil
}
