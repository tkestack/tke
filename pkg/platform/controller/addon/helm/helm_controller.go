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

package helm

import (
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	helmClientRetryCount = 5

	helmMaxRetryCount = 5
	helmTimeOut       = 5 * time.Minute
)

// Controller is responsible for performing actions dependent upon a helm phase.
type Controller struct {
	client       clientset.Interface
	cache        *helmCache
	prober       Prober
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.HelmLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, helmInformer platformv1informer.HelmInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client: client,
		cache:  &helmCache{helmMap: make(map[string]*cachedHelm)},
		prober: NewHealthProber(helmInformer.Lister(), client),
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "helm"),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("helm_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the helm informer event handlers
	helmInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldHelm, ok1 := oldObj.(*v1.Helm)
				curHelm, ok2 := newObj.(*v1.Helm)
				if ok1 && ok2 && controller.needsUpdate(oldHelm, curHelm) {
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.lister = helmInformer.Lister()
	controller.listerSynced = helmInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.Helm, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldHelm *v1.Helm, newHelm *v1.Helm) bool {
	return !reflect.DeepEqual(oldHelm, newHelm)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting helm controller")
	defer log.Info("Shutting down helm controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for cluster caches to sync")
	}

	c.prober.Run(stopCh)
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

	err := c.sync(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing helm %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// sync will sync the Helm with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) sync(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing helm", log.String("helmName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// helm holds the latest helm info from apiserver
	helm, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Helm has been deleted. Attempting to cleanup resources", log.String("helm", key))
		err = c.processHelmDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve helm from store", log.String("helm", key), log.Err(err))
	default:
		err = c.processCreateOrUpdate(helm, key)
	}
	return err
}

func (c *Controller) processHelmDeletion(key string) error {
	cachedHelm, ok := c.cache.get(key)
	if !ok {
		log.Error("Helm not in cache even though the watcher thought it was. Ignoring the deletion", log.String("helmName", key))
		return nil
	}
	return c.processDelete(cachedHelm, key)
}

func (c *Controller) processDelete(cachedHelm *cachedHelm, key string) error {
	log.Info("helm will be dropped", log.String("helm", key))
	if c.cache.Exist(key) {
		c.cache.delete(key)
	}
	if c.prober.Exist(key) {
		c.prober.Del(key)
	}
	helm := cachedHelm.state

	var provisioner Provisioner
	var err error
	if provisioner, err = createProvisioner(helm, c.client); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return provisioner.Uninstall()
}

func (c *Controller) processCreateOrUpdate(holder *v1.Helm, key string) error {
	cached := c.cache.getOrCreate(key)
	if cached.state != nil && cached.state.UID != holder.UID {
		// if the same key holder is created, delete it
		if err := c.processDelete(cached, key); err != nil {
			return err
		}
	}

	err := c.handlePhase(key, cached, holder)
	if err != nil {
		return err
	}

	cached.state = holder
	// Always update the cache upon success.
	c.cache.set(key, cached)
	return nil
}

func (c *Controller) handlePhase(key string, cachedHelm *cachedHelm, holder *v1.Helm) error {
	phase := holder.Status.Phase
	log.Info(fmt.Sprintf("Helm is %s", string(phase)), log.String("helm", key))
	switch phase {
	case v1.AddonPhaseInitializing:
		return c.doInitializing(key, holder)
	case v1.AddonPhaseReinitializing:
		c.doReinitializing(key, holder)
	case v1.AddonPhaseChecking, v1.AddonPhaseRunning, v1.AddonPhaseFailed, v1.AddonPhaseUnhealthy:
		if !c.prober.ExistByPhase(key, phase) {
			c.prober.Set(key, phase)
		}
	}
	return nil
}

func (c *Controller) doInitializing(key string, holder *v1.Helm) error {
	defer controllerutil.CatchPanic("doInitializing", "Helm")

	if c.prober.Exist(key) {
		c.prober.Del(key)
	}

	var provisioner Provisioner
	var err error
	if provisioner, err = createProvisioner(holder, c.client); err != nil {
		return err
	}
	if err := provisioner.Install(); err != nil {
		// if user install his own tiller, update helm status to fail
		if errors.IsConflict(err) {
			return updateHelmStatus(getUpdateObj(holder, v1.AddonPhaseFailed, err.Error()), c.client)
		}
		return updateHelmStatus(getUpdateObj(holder, v1.AddonPhaseReinitializing, err.Error()), c.client)
	}
	newObj := getUpdateObj(holder, v1.AddonPhaseChecking, "")
	newObj.Status.LastReInitializingTimestamp = metav1.Now()
	return updateHelmStatus(newObj, c.client)
}

func (c *Controller) doReinitializing(key string, helm *v1.Helm) {
	var interval = time.Since(helm.Status.LastReInitializingTimestamp.Time)
	var waitTime time.Duration
	if interval >= helmTimeOut {
		waitTime = time.Duration(1)
	} else {
		waitTime = helmTimeOut - interval
	}
	go func() {
		defer controllerutil.CatchPanic("reinitialize", "Helm")
		if err := wait.Poll(waitTime, helmTimeOut, c.reinitialize(key, helm)); err != nil {
			log.Info(fmt.Sprintf("reinitialize err: %v", err))
		}
	}()
}

func (c *Controller) reinitialize(key string, holder *v1.Helm) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		var provisioner Provisioner
		var err error
		if provisioner, err = createProvisioner(holder, c.client); err == nil {
			_ = provisioner.Uninstall()
			if err = provisioner.Install(); err == nil {
				newObj := getUpdateObj(holder, v1.AddonPhaseChecking, "")
				newObj.Status.LastReInitializingTimestamp = metav1.Now()
				return true, updateHelmStatus(newObj, c.client)
			}
			if errors.IsConflict(err) {
				return true, updateHelmStatus(getUpdateObj(holder, v1.AddonPhaseFailed, err.Error()), c.client)
			}
		}
		if holder.Status.RetryCount >= helmMaxRetryCount {
			err := fmt.Sprintf("Install error and retried max(%d) times already.", helmMaxRetryCount)
			return true, updateHelmStatus(getUpdateObj(holder, v1.AddonPhaseFailed, err), c.client)
		}
		return true, updateHelmStatus(getUpdateObj(holder, v1.AddonPhaseReinitializing, err.Error()), c.client)
	}
}

func createProvisioner(helm *v1.Helm, client clientset.Interface) (Provisioner, error) {
	cluster, err := client.PlatformV1().Clusters().Get(helm.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, client.PlatformV1())
	if err != nil {
		return nil, err
	}
	isExtensionsAPIGroup := controllerutil.IsClusterVersionBefore1_9(kubeClient)

	provisioner := NewProvisioner(kubeClient, &Option{
		version:              helm.Spec.Version,
		isExtensionsAPIGroup: isExtensionsAPIGroup,
	})
	return provisioner, nil
}

// updateHelmStatus means update status to the given object
func updateHelmStatus(obj *v1.Helm, client clientset.Interface) error {
	return wait.PollImmediate(time.Second, helmClientRetryCount*time.Second, func() (done bool, err error) {
		_, err = client.PlatformV1().Helms().UpdateStatus(obj)
		if err == nil {
			return true, nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to helm that no longer exists", log.String("clusterName", obj.Spec.ClusterName), log.Err(err))
			return true, nil
		}
		if errors.IsConflict(err) {
			return false, fmt.Errorf("not persisting update to helm '%s' that has been changed since we received it: %v", obj.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of helm '%s/%s'", obj.Spec.ClusterName, obj.Status.Phase), log.String("clusterName", obj.Spec.ClusterName), log.Err(err))
		return false, err
	})
}

func getUpdateObj(original *v1.Helm, phase v1.AddonPhase, reason string) *v1.Helm {
	cloned := original.DeepCopy()
	cloned.Status.Phase = phase
	cloned.Status.Reason = reason
	switch phase {
	case v1.AddonPhaseReinitializing:
		cloned.Status.RetryCount++
		cloned.Status.LastReInitializingTimestamp = metav1.Now()
	default:
	}
	return cloned
}
