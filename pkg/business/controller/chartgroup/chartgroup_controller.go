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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	businessv1 "tkestack.io/tke/api/business/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	businessv1informer "tkestack.io/tke/api/client/informers/externalversions/business/v1"
	businessv1lister "tkestack.io/tke/api/client/listers/business/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/business/controller/chartgroup/deletion"
	controllerutil "tkestack.io/tke/pkg/controller"
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
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)
const (
	controllerName = "chartGroup-controller"
)

// Controller is responsible for performing actions dependent upon an chartGroup phase.
type Controller struct {
	client         clientset.Interface
	registryClient registryversionedclient.RegistryV1Interface
	cache          *chartGroupCache
	health         *chartGroupHealth
	queue          workqueue.RateLimitingInterface
	lister         businessv1lister.ChartGroupLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// helper to delete all resources in the chartGroup when the chartGroup is deleted.
	chartGroupResourcesDeleter deletion.ChartGroupResourcesDeleterInterface
}

// NewController creates a new Controller object.
func NewController(registryClient registryversionedclient.RegistryV1Interface,
	client clientset.Interface, chartGroupInformer businessv1informer.ChartGroupInformer,
	resyncPeriod time.Duration, finalizerToken businessv1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                     client,
		registryClient:             registryClient,
		cache:                      &chartGroupCache{m: make(map[string]*cachedChartGroup)},
		health:                     &chartGroupHealth{chartGroups: sets.NewString()},
		queue:                      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		chartGroupResourcesDeleter: deletion.NewChartGroupResourcesDeleter(registryClient, client.BusinessV1(), finalizerToken, true),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("chartGroup_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	chartGroupInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*businessv1.ChartGroup)
				cur, ok2 := newObj.(*businessv1.ChartGroup)
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

// obj could be an *businessv1.ChartGroup, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, chartGroupDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *businessv1.ChartGroup, new *businessv1.ChartGroup) bool {
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

	projectName, chartGroupName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// chartGroup holds the latest ChartGroup info from apiserver
	chartGroup, err := c.lister.ChartGroups(projectName).Get(chartGroupName)
	switch {
	case errors.IsNotFound(err):
		log.Info("ChartGroup has been deleted. Attempting to cleanup resources",
			log.String("projectName", projectName), log.String("chartGroupName", chartGroupName))
		err = c.processDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve chartGroup from store",
			log.String("projectName", projectName), log.String("chartGroupName", chartGroupName), log.Err(err))
	default:
		if chartGroup.Status.Phase == businessv1.ChartGroupPending ||
			chartGroup.Status.Phase == businessv1.ChartGroupAvailable ||
			chartGroup.Status.Phase == businessv1.ChartGroupLocked {
			cachedChartGroup := c.cache.getOrCreate(key)
			err = c.processUpdate(context.Background(), cachedChartGroup, chartGroup, key)
		} else if chartGroup.Status.Phase == businessv1.ChartGroupTerminating {
			log.Info("ChartGroup has been terminated. Attempting to cleanup resources",
				log.String("projectName", projectName), log.String("chartGroupName", chartGroupName))
			_ = c.processDeletion(key)
			err = c.chartGroupResourcesDeleter.Delete(projectName, chartGroupName)
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

func (c *Controller) processUpdate(ctx context.Context, cachedChartGroup *cachedChartGroup, chartGroup *businessv1.ChartGroup, key string) error {
	if cachedChartGroup.state != nil {
		// exist and the chartGroup name changed
		if cachedChartGroup.state.UID != chartGroup.UID {
			if err := c.processDelete(cachedChartGroup, key); err != nil {
				return err
			}
		}
	}
	// start update machine if needed
	err := c.handlePhase(ctx, key, cachedChartGroup, chartGroup)
	if err != nil {
		return err
	}
	cachedChartGroup.state = chartGroup
	// Always update the cache upon success.
	c.cache.set(key, cachedChartGroup)
	return nil
}

func (c *Controller) handlePhase(ctx context.Context, key string, cachedChartGroup *cachedChartGroup, chartGroup *businessv1.ChartGroup) error {
	switch chartGroup.Status.Phase {
	case businessv1.ChartGroupPending:
		err := c.createChartGroup(ctx, chartGroup)
		if err != nil {
			chartGroup.Status.Phase = businessv1.ChartGroupFailed
			chartGroup.Status.Message = "createChartGroup failed"
			chartGroup.Status.Reason = err.Error()
			chartGroup.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, chartGroup)
		}
		chartGroup.Status.Phase = businessv1.ChartGroupAvailable
		chartGroup.Status.Message = ""
		chartGroup.Status.Reason = ""
		chartGroup.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdate(ctx, chartGroup)
	case businessv1.ChartGroupAvailable, businessv1.ChartGroupLocked:
		c.startChartGroupHealthCheck(ctx, key)
	}
	return nil
}

func (c *Controller) createChartGroup(ctx context.Context, chartGroup *businessv1.ChartGroup) error {
	_, err := c.registryClient.ChartGroups().Create(ctx, &registryv1.ChartGroup{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"projectName": chartGroup.Namespace,
			},
		},
		Spec: registryv1.ChartGroupSpec{
			Name:        chartGroup.Name,
			DisplayName: chartGroup.Spec.DisplayName,
			TenantID:    chartGroup.Spec.TenantID,
		}}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) persistUpdate(ctx context.Context, chartGroup *businessv1.ChartGroup) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.BusinessV1().ChartGroups(chartGroup.Namespace).UpdateStatus(ctx, chartGroup, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to chartGroup that no longer exists",
				log.String("projectName", chartGroup.Namespace),
				log.String("chartGroupName", chartGroup.Name),
				log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to chartGroup '%s/%s' that has been changed since we received it: %v",
				chartGroup.Namespace, chartGroup.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of chartGroup '%s/%s/%s'",
			chartGroup.Namespace, chartGroup.Name, chartGroup.Status.Phase),
			log.String("chartGroupName", chartGroup.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}
