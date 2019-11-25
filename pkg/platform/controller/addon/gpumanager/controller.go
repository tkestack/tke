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

package gpumanager

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/health"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/operator"
	"tkestack.io/tke/pkg/platform/controller/addon/gpumanager/utils"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

// Controller is responsible for performing actions dependent upon a gpu manager phase.
type Controller struct {
	operator     operator.ObjectOperator
	cache        *gmCache
	health       health.Prober
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.GPUManagerLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

const (
	maxInstallationRetryTime = 5
)

// NewController creates a new Controller object.
func NewController(tkeClient tkeclientset.Interface, gmInformer platformv1informer.GPUManagerInformer, resyncPeriod time.Duration) *Controller {
	limiter := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 10*time.Second),
		// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
	)

	// create the controller so we can inject the enqueue function
	controller := &Controller{
		operator: operator.NewObjectOperator(tkeClient),
		cache:    &gmCache{store: make(map[string]*gmCachedItem)},
		queue:    workqueue.NewNamedRateLimitingQueue(limiter, "gpumanager"),
	}

	if tkeClient != nil && tkeClient.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("gm_controller", tkeClient.PlatformV1().RESTClient().GetRateLimiter())
	}

	gmInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueItem,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldItem, ok1 := oldObj.(*v1.GPUManager)
				curItem, ok2 := newObj.(*v1.GPUManager)
				if ok1 && ok2 && controller.needsUpdate(oldItem, curItem) {
					controller.enqueueItem(newObj)
				}
			},
			DeleteFunc: controller.enqueueItem,
		},
		resyncPeriod,
	)
	controller.lister = gmInformer.Lister()
	controller.listerSynced = gmInformer.Informer().HasSynced
	controller.health = health.NewHealthProber(gmInformer.Lister(), controller.operator)

	return controller
}

// obj could be an *v1.GPUManager, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueItem(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("obj", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldItem *v1.GPUManager, newItem *v1.GPUManager) bool {
	return !reflect.DeepEqual(oldItem, newItem)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info(fmt.Sprintf("Starting gpu manager controller"))
	defer log.Info(fmt.Sprintf("Shutting down gpu manager controller"))

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for cluster caches to sync")
	}

	c.health.Run(stopCh)
	c.stopCh = stopCh

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

// worker processes the queue of namespace objects.
// Each namespace can be in the queue at most once.
// The system ensures that no two workers can process
// the same namespace at the same time.
func (c *Controller) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncItem(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing gpu manager %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncItem will sync the GPUManager with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Finished syncing GPUManager %s, elapsed %f secs", key, time.Since(startTime).Seconds()))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	var cachedItem *gmCachedItem
	gm, err := c.lister.Get(name)
	switch {
	case apierror.IsNotFound(err):
		log.Info(fmt.Sprintf("GPUManager %s has been deleted. Attempting to cleanup resources", key))
		err = c.processDeletion(key)
	case err != nil:
		log.Warn(fmt.Sprintf("Unable to retrieve GPUManager %s from store, err %+v", key, err))
	default:
		log.Info(fmt.Sprintf("Start syncing GPUManager %s, status %s", gm.Name, string(gm.Status.Phase)))
		cachedItem = c.cache.getOrCreate(key)
		err = c.processUpdate(cachedItem, gm, key)
	}
	return err
}

func (c *Controller) processDeletion(key string) error {
	cachedItem, ok := c.cache.get(key)
	if !ok {
		log.Warn(fmt.Sprintf("GPUManager %s not in store. Ignoring the deletion", key))
		return nil
	}
	return c.delete(cachedItem, key)
}

func (c *Controller) processUpdate(cachedItem *gmCachedItem, holder *v1.GPUManager, key string) error {
	if cachedItem.holder != nil {
		// exist and the cluster name changed
		if cachedItem.holder.UID != holder.UID {
			if err := c.delete(cachedItem, key); err != nil {
				return err
			}
		}
	}

	err := c.createIfNeeded(key, cachedItem, holder)
	if err != nil {
		return err
	}

	cachedItem.holder = holder
	// Always updateStatus the store upon success.
	c.cache.set(key, cachedItem)
	return nil
}

func (c *Controller) createIfNeeded(key string, cachedItem *gmCachedItem, holder *v1.GPUManager) error {
	switch holder.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info(fmt.Sprintf("GPUManager %s will be created", key))
		return c.updateStatus(utils.ForUpdateItem(holder, v1.AddonPhasePending, ""))
	case v1.AddonPhasePending:
		log.Info(fmt.Sprintf("GPUManager %s is prepared to install ", key))
		return c.install(holder)
	case v1.AddonPhaseRunning:
		c.health.Set(holder)

		//If phase transits from pending to running, at the same time, we update spec version.
		//In this situation, the controller will ignore this update which is unexpected. So we
		//retrieve a fresh daemonset, and compare the version.
		needUpgrade, err := c.operator.DiffDaemonSet(holder)
		if err != nil {
			return err
		}

		if needUpgrade {
			log.Info(fmt.Sprintf("GPUManager %s needs upgrading", key))
			return c.updateStatus(utils.ForUpdateItem(holder, v1.AddonPhaseUpgrading, ""))
		}

		log.Warn(fmt.Sprintf("GPUManager %s is changed but without the version field modification", key))
		return nil
	case v1.AddonPhaseUpgrading:
		log.Info(fmt.Sprintf("GPUManager %s is going to upgrade", key))

		return c.upgrade(holder)
	case v1.AddonPhaseFailed:
		log.Info(fmt.Sprintf("GPUManager %s is error", key))
		c.health.Del(key)
	case v1.AddonPhaseUnhealthy:
		c.health.Set(holder)
	default:
		log.Warn(fmt.Sprintf("GPUManager %s got unknown phrase %s", key, holder.Status.Phase))
	}
	return nil
}

func (c *Controller) delete(item *gmCachedItem, key string) error {
	log.Info(fmt.Sprintf("GPUManager %s will be dropped", key))

	if c.cache.Exist(key) {
		log.Info(fmt.Sprintf("Delete the GPUManager %s from store", key))
		c.cache.delete(key)
	}

	log.Info(fmt.Sprintf("Delete the GPUManager %s health store", key))
	c.health.Del(key)

	log.Info(fmt.Sprintf("Delete service of GPUManager %s", key))
	if err := c.operator.DeleteService(item.holder.Spec.ClusterName); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Delete service account and role binding of GPUManager %s", key))
	if err := c.operator.DeleteServiceAccount(item.holder.Spec.ClusterName); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Delete daemonset of GPUManager %s", key))
	if err := c.operator.DeleteDaemonSet(item.holder); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Delete deployment of GPUManager quota admission %s", key))
	return c.operator.DeleteDeployment(item.holder)
}

func (c *Controller) install(holder *v1.GPUManager) error {
	var err error

	log.Info(fmt.Sprintf("Create service account and role binding of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateServiceAccount(holder.Spec.ClusterName)
	if err != nil {
		goto FAILED
	}

	log.Info(fmt.Sprintf("Create configmap of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateConfigmap(holder.Spec.ClusterName)
	if err != nil {
		goto FAILED
	}

	log.Info(fmt.Sprintf("Create daemonset of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateDaemonSet(holder)
	if err != nil {
		goto FAILED
	}

	log.Info(fmt.Sprintf("Create deployment of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateDeployment(holder)
	if err != nil {
		goto FAILED
	}

	log.Info(fmt.Sprintf("Create service of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateService(holder.Spec.ClusterName)
	if err != nil {
		goto FAILED
	}

	log.Info(fmt.Sprintf("Create metric service of GPUManager %s", holder.Spec.ClusterName))
	err = c.operator.CreateServiceMetric(holder.Spec.ClusterName)
	if err != nil {
		goto FAILED
	}

	return c.updateStatus(utils.ForUpdateItem(holder, v1.AddonPhaseRunning, ""))

FAILED:
	log.Error(fmt.Sprintf("installation failed, backoff %s", holder.Spec.ClusterName))
	if holder.Status.RetryCount >= maxInstallationRetryTime {
		_ = c.updateStatus(utils.ForUpdateItem(holder, v1.AddonPhaseFailed,
			errors.Wrapf(err, "reach %d installation retry limit", maxInstallationRetryTime).Error()))
		return err
	}

	newStatus := utils.ForUpdateItem(holder, holder.Status.Phase, "")
	newStatus.Status.RetryCount++
	_ = c.updateStatus(newStatus)

	_ = c.operator.DeleteConfigmap(holder.Spec.ClusterName)
	_ = c.operator.DeleteService(holder.Spec.ClusterName)
	_ = c.operator.DeleteServiceMetric(holder.Spec.ClusterName)
	_ = c.operator.DeleteServiceAccount(holder.Spec.ClusterName)
	_ = c.operator.DeleteDaemonSet(holder)
	_ = c.operator.DeleteDeployment(holder)

	return err
}

func (c *Controller) updateStatus(holder *v1.GPUManager) error {
	return c.operator.UpdateGPUManagerStatus(holder)
}

func (c *Controller) upgrade(holder *v1.GPUManager) error {
	_ = c.operator.UpdateDaemonSet(holder)
	_ = c.operator.UpdateDeployment(holder)
	return nil
}
