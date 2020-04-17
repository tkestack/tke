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

package emigration

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
	v1 "tkestack.io/tke/api/business/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	businessv1informer "tkestack.io/tke/api/client/informers/externalversions/business/v1"
	businessv1lister "tkestack.io/tke/api/client/listers/business/v1"
	businessns "tkestack.io/tke/pkg/business/controller/namespace"
	cls "tkestack.io/tke/pkg/business/controller/namespace/cluster"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// emigrationDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	emigrationDeletionGracePeriod = 5 * time.Second
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)
const (
	controllerName = "nsemigration-controller"
)

// Controller is responsible for performing actions dependent upon a emigration phase.
type Controller struct {
	client         clientset.Interface
	platformClient platformversionedclient.PlatformV1Interface
	queue          workqueue.RateLimitingInterface
	lister         businessv1lister.NsEmigrationLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(platformClient platformversionedclient.PlatformV1Interface, client clientset.Interface,
	emigrationInformer businessv1informer.NsEmigrationInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:         client,
		platformClient: platformClient,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("emigration_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	emigrationInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.NsEmigration)
				cur, ok2 := newObj.(*v1.NsEmigration)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.lister = emigrationInformer.Lister()
	controller.listerSynced = emigrationInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.NsEmigration, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, emigrationDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.NsEmigration, new *v1.NsEmigration) bool {
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
	log.Info("Starting emigration controller")
	defer log.Info("Shutting down emigration controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for emigration caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of emigration objects.
// Each emigration can be in the queue at most once.
// The system ensures that no two workers can process
// the same emigration at the same time.
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

		// rather than wait for a full resync, re-add the emigration to the queue to be processed
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

// syncItem will sync the NsEmigration with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// emigrations created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	log.Info("Start syncing emigration", log.String("emigration", key))
	defer func() {
		log.Info("Finished syncing emigration", log.String("emigration", key),
			log.Duration("processTime", time.Since(startTime)))
	}()

	projectName, emigrationName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// Emigration holds the latest Emigration info from apiserver
	emigration, err := c.lister.NsEmigrations(projectName).Get(emigrationName)
	switch {
	case errors.IsNotFound(err):
		log.Info("Emigration has been deleted. Attempting to cleanup resources",
			log.String("projectName", projectName), log.String("emigrationName", emigrationName))
		return nil
	case err != nil:
		log.Warn("Unable to retrieve emigration from store", log.String("projectName", projectName),
			log.String("emigrationName", emigrationName), log.Err(err))
	default:
		err = c.processUpdate(emigration)
	}
	return err
}

func (c *Controller) processUpdate(emigration *v1.NsEmigration) error {
	switch emigration.Status.Phase {
	case v1.NsEmigrationPending:
		return c.processPending(emigration)
	case v1.NsEmigrationOldOneLocked:
		return c.processOldOneLocked(emigration)
	case v1.NsEmigrationOldOneDetached:
		return c.processOldOneDetached(emigration)
	case v1.NsEmigrationNewOneCreated:
		return c.processNewOneCreated(emigration)
	case v1.NsEmigrationOldOneTerminating:
		return c.processOldOneTerminating(emigration)
	case v1.NsEmigrationFinished:
		return c.processFinished(emigration)
	default:
		return c.processOthers(emigration)
	}
}

func (c *Controller) getOldNamespace(emigration *v1.NsEmigration) (*v1.Namespace, error) {
	oldNS, err := c.client.BusinessV1().Namespaces(emigration.Namespace).Get(emigration.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		oldNS = nil
		emigration.Status.Message = fmt.Sprintf("%s, failed to get namespace %s/%s",
			emigration.Status.Phase, emigration.Namespace, emigration.Spec.NsShowName)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.Reason = err.Error()
		emigration.Status.LastTransitionTime = metav1.Now()
	}
	return oldNS, err
}

func (c *Controller) processPending(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationPending {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationPending))
	}
	oldNS, err := c.getOldNamespace(emigration)
	if err != nil {
		return c.persistUpdateEmigration(emigration)
	}
	if oldNS.Status.Phase != v1.NamespaceAvailable {
		emigration.Status.Message = fmt.Sprintf("%s, namespace %s/%s is NOT in phase %s",
			emigration.Status.Phase, oldNS.Namespace, oldNS.Spec.Namespace, oldNS.Status.Phase)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	oldNS.Status.Phase = v1.NamespaceLocked
	if err := businessns.PersistUpdateNamesapce(c.client, oldNS); err != nil {
		emigration.Status.Message = fmt.Sprintf("%s, failed to lock namespace %s/%s",
			emigration.Status.Phase, oldNS.Namespace, oldNS.Spec.Namespace)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.Reason = err.Error()
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	emigration.Status.Phase = v1.NsEmigrationOldOneLocked
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) processOldOneLocked(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationOldOneLocked {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationOldOneLocked))
	}
	oldNS, err := c.getOldNamespace(emigration)
	if err != nil {
		return c.persistUpdateEmigration(emigration)
	}
	if err := c.detachFromClusterNamespace(oldNS); err != nil {
		emigration.Status.Message = fmt.Sprintf("%s, failed to detach namespace %s/%s from cluster",
			emigration.Status.Phase, oldNS.Namespace, oldNS.Spec.Namespace)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.Reason = err.Error()
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	emigration.Status.Phase = v1.NsEmigrationOldOneDetached
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) processOldOneDetached(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationOldOneDetached {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationOldOneDetached))
	}
	oldNS, err := c.getOldNamespace(emigration)
	if err != nil {
		return c.persistUpdateEmigration(emigration)
	}
	newNS := v1.Namespace{}
	newNS.Namespace = emigration.Spec.Destination
	newNS.Spec.TenantID = oldNS.Spec.TenantID
	newNS.Spec.ClusterName = oldNS.Spec.ClusterName
	newNS.Spec.Namespace = oldNS.Spec.Namespace
	newNS.Spec.Hard = oldNS.Spec.Hard
	if _, err := c.client.BusinessV1().Namespaces(newNS.Namespace).Create(&newNS); err != nil {
		emigration.Status.Message = fmt.Sprintf("%s, failed to create namespace %s/%s",
			emigration.Status.Phase, emigration.Spec.Destination, emigration.Spec.NsShowName)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.Reason = err.Error()
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	emigration.Status.Phase = v1.NsEmigrationNewOneCreated
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) processNewOneCreated(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationNewOneCreated {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationNewOneCreated))
	}
	newNS, err := c.client.BusinessV1().Namespaces(emigration.Spec.Destination).Get(emigration.Spec.Namespace, metav1.GetOptions{})
	if err != nil {
		emigration.Status.Message = fmt.Sprintf("%s, failed to check status of namespace %s/%s",
			emigration.Status.Phase, emigration.Spec.Destination, emigration.Spec.NsShowName)
		emigration.Status.Reason = err.Error()
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	// waiting for newNS to be NamespaceAvailable
	if newNS.Status.Phase != v1.NamespaceAvailable {
		if newNS.Status.Phase != v1.NamespacePending {
			emigration.Status.Message = fmt.Sprintf("%s, status of namespace %s/%s is %s",
				emigration.Status.Phase, emigration.Spec.Destination, emigration.Spec.NsShowName, newNS.Status.Phase)
			emigration.Status.Phase = v1.NsEmigrationFailed
			emigration.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdateEmigration(emigration)
		}
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	oldNS, err := c.client.BusinessV1().Namespaces(emigration.Namespace).Get(emigration.Spec.Namespace, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		emigration.Status.Message = fmt.Sprintf("%s, failed to check status of namespace %s/%s",
			emigration.Status.Phase, emigration.Namespace, emigration.Spec.NsShowName)
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.Reason = err.Error()
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	if err == nil {
		if oldNS.Status.Phase != v1.NamespaceLocked {
			emigration.Status.Message = fmt.Sprintf("%s, status of namespace %s/%s is %s",
				emigration.Status.Phase, oldNS.Namespace, oldNS.Spec.Namespace, oldNS.Status.Phase)
			emigration.Status.Phase = v1.NsEmigrationFailed
			emigration.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdateEmigration(emigration)
		}
		background := metav1.DeletePropagationBackground
		deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
		if err := c.client.BusinessV1().Namespaces(oldNS.Namespace).Delete(oldNS.Name, deleteOpt); err != nil && !errors.IsNotFound(err) {
			emigration.Status.Message = fmt.Sprintf("%s, failed to delete namespace %s/%s",
				emigration.Status.Phase, oldNS.Namespace, oldNS.Spec.Namespace)
			emigration.Status.Phase = v1.NsEmigrationFailed
			emigration.Status.Reason = err.Error()
			emigration.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdateEmigration(emigration)
		}
	}
	emigration.Status.Phase = v1.NsEmigrationOldOneTerminating
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) processOldOneTerminating(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationOldOneTerminating {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationOldOneTerminating))
	}
	oldNS, err := c.client.BusinessV1().Namespaces(emigration.Namespace).Get(emigration.Spec.Namespace, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		emigration.Status.Message = fmt.Sprintf("%s, failed to check status of namespace %s/%s",
			emigration.Status.Phase, emigration.Namespace, emigration.Spec.NsShowName)
		emigration.Status.Reason = err.Error()
		emigration.Status.Phase = v1.NsEmigrationFailed
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	if err == nil {
		if oldNS.Status.Phase != v1.NamespaceTerminating {
			emigration.Status.Message = fmt.Sprintf("%s, status namespace %s/%s is %s",
				emigration.Status.Phase, emigration.Namespace, emigration.Spec.NsShowName, oldNS.Status.Phase)
			emigration.Status.Phase = v1.NsEmigrationFailed
			emigration.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdateEmigration(emigration)
		}
		emigration.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdateEmigration(emigration)
	}
	emigration.Status.Phase = v1.NsEmigrationFinished
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) processFinished(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase != v1.NsEmigrationFinished {
		panic(fmt.Sprintf("%s != %s", emigration.Status.Phase, v1.NsEmigrationFinished))
	}
	return c.client.BusinessV1().NsEmigrations(emigration.Namespace).Delete(emigration.Name, &metav1.DeleteOptions{})
}

func (c *Controller) processOthers(emigration *v1.NsEmigration) error {
	if emigration.Status.Phase == v1.NsEmigrationFailed {
		return nil
	}
	emigration.Status.Message = fmt.Sprintf("invalid emigration phase %s", emigration.Status.Phase)
	emigration.Status.Phase = v1.NsEmigrationFailed
	emigration.Status.LastTransitionTime = metav1.Now()
	return c.persistUpdateEmigration(emigration)
}

func (c *Controller) detachFromClusterNamespace(namespace *v1.Namespace) error {
	kubeClient, err := util.BuildExternalClientSetWithName(c.platformClient, namespace.Spec.ClusterName)
	if err != nil {
		log.Error("Failed to create the kubernetes client",
			log.String("namespaceName", namespace.ObjectMeta.Name),
			log.String("clusterName", namespace.Spec.ClusterName),
			log.Err(err))
		return err
	}
	return cls.DetachFromClusterNamespace(kubeClient, namespace)
}

func (c *Controller) persistUpdateEmigration(emigration *v1.NsEmigration) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.BusinessV1().NsEmigrations(emigration.ObjectMeta.Namespace).Update(emigration)
		if err == nil {
			return nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to emigration that no longer exists",
				log.String("projectName", emigration.ObjectMeta.Namespace),
				log.String("emigrationName", emigration.ObjectMeta.Name),
				log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to emigration '%s/%s' that has been changed since we received it: %v",
				emigration.ObjectMeta.Namespace, emigration.ObjectMeta.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of emigration '%s/%s/%s'",
			emigration.ObjectMeta.Namespace, emigration.ObjectMeta.Name, emigration.Status.Phase),
			log.String("emigrationName", emigration.ObjectMeta.Name),
			log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}
