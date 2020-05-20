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

package imagenamespace

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
	"tkestack.io/tke/pkg/business/controller/imagenamespace/deletion"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// imageNamespaceDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	imageNamespaceDeletionGracePeriod = 5 * time.Second
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)
const (
	controllerName = "imageNamespace-controller"
)

// Controller is responsible for performing actions dependent upon an imageNamespace phase.
type Controller struct {
	client         clientset.Interface
	registryClient registryversionedclient.RegistryV1Interface
	cache          *imageNamespaceCache
	health         *imageNamespaceHealth
	queue          workqueue.RateLimitingInterface
	lister         businessv1lister.ImageNamespaceLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// helper to delete all resources in the imageNamespace when the imageNamespace is deleted.
	imageNamespaceResourcesDeleter deletion.ImageNamespaceResourcesDeleterInterface
}

// NewController creates a new Controller object.
func NewController(registryClient registryversionedclient.RegistryV1Interface,
	client clientset.Interface, imageNamespaceInformer businessv1informer.ImageNamespaceInformer,
	resyncPeriod time.Duration, finalizerToken businessv1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                         client,
		registryClient:                 registryClient,
		cache:                          &imageNamespaceCache{m: make(map[string]*cachedImageNamespace)},
		health:                         &imageNamespaceHealth{imageNamespaces: sets.NewString()},
		queue:                          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		imageNamespaceResourcesDeleter: deletion.NewImageNamespaceResourcesDeleter(registryClient, client.BusinessV1(), finalizerToken, true),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("imageNamespace_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	imageNamespaceInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*businessv1.ImageNamespace)
				cur, ok2 := newObj.(*businessv1.ImageNamespace)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.lister = imageNamespaceInformer.Lister()
	controller.listerSynced = imageNamespaceInformer.Informer().HasSynced

	return controller
}

// obj could be an *businessv1.ImageNamespace, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, imageNamespaceDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *businessv1.ImageNamespace, new *businessv1.ImageNamespace) bool {
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
	log.Info("Starting imageNamespace controller")
	defer log.Info("Shutting down imageNamespace controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for imageNamespace caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of imageNamespace objects.
// Each imageNamespace can be in the queue at most once.
// The system ensures that no two workers can process
// the same imageNamespace at the same time.
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

		// rather than wait for a full resync, re-add the imageNamespace to the queue to be processed
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

// syncItem will sync the ImageNamespace with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// imageNamespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing imageNamespace", log.String("imageNamespace", key), log.Duration("processTime", time.Since(startTime)))
	}()

	projectName, imageNamespaceName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// imageNamespace holds the latest ImageNamespace info from apiserver
	imageNamespace, err := c.lister.ImageNamespaces(projectName).Get(imageNamespaceName)
	switch {
	case errors.IsNotFound(err):
		log.Info("ImageNamespace has been deleted. Attempting to cleanup resources",
			log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName))
		err = c.processDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve imageNamespace from store",
			log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName), log.Err(err))
	default:
		if imageNamespace.Status.Phase == businessv1.ImageNamespacePending ||
			imageNamespace.Status.Phase == businessv1.ImageNamespaceAvailable ||
			imageNamespace.Status.Phase == businessv1.ImageNamespaceLocked {
			cachedImageNamespace := c.cache.getOrCreate(key)
			err = c.processUpdate(context.Background(), cachedImageNamespace, imageNamespace, key)
		} else if imageNamespace.Status.Phase == businessv1.ImageNamespaceTerminating {
			log.Info("ImageNamespace has been terminated. Attempting to cleanup resources",
				log.String("projectName", projectName), log.String("imageNamespaceName", imageNamespaceName))
			_ = c.processDeletion(key)
			err = c.imageNamespaceResourcesDeleter.Delete(projectName, imageNamespaceName)
		} else {
			log.Debug(fmt.Sprintf("ImageNamespace %s status is %s, not to process", key, imageNamespace.Status.Phase))
		}
	}
	return err
}

func (c *Controller) processDeletion(key string) error {
	cachedImageNamespace, ok := c.cache.get(key)
	if !ok {
		log.Debug("ImageNamespace not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processDelete(cachedImageNamespace, key)
}

func (c *Controller) processDelete(cachedImageNamespace *cachedImageNamespace, key string) error {
	log.Info("ImageNamespace will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the imageNamespace cache", log.String("name", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the imageNamespace health cache", log.String("name", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) processUpdate(ctx context.Context, cachedImageNamespace *cachedImageNamespace, imageNamespace *businessv1.ImageNamespace, key string) error {
	if cachedImageNamespace.state != nil {
		// exist and the imageNamespace name changed
		if cachedImageNamespace.state.UID != imageNamespace.UID {
			if err := c.processDelete(cachedImageNamespace, key); err != nil {
				return err
			}
		}
	}
	// start update machine if needed
	err := c.handlePhase(ctx, key, cachedImageNamespace, imageNamespace)
	if err != nil {
		return err
	}
	cachedImageNamespace.state = imageNamespace
	// Always update the cache upon success.
	c.cache.set(key, cachedImageNamespace)
	return nil
}

func (c *Controller) handlePhase(ctx context.Context, key string, cachedImageNamespace *cachedImageNamespace, imageNamespace *businessv1.ImageNamespace) error {
	switch imageNamespace.Status.Phase {
	case businessv1.ImageNamespacePending:
		err := c.createRegistryNamespace(ctx, imageNamespace)
		if err != nil {
			imageNamespace.Status.Phase = businessv1.ImageNamespaceFailed
			imageNamespace.Status.Message = "CreateRegistryNamespace failed"
			imageNamespace.Status.Reason = err.Error()
			imageNamespace.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, imageNamespace)
		}
		imageNamespace.Status.Phase = businessv1.ImageNamespaceAvailable
		imageNamespace.Status.Message = ""
		imageNamespace.Status.Reason = ""
		imageNamespace.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdate(ctx, imageNamespace)
	case businessv1.ImageNamespaceAvailable, businessv1.ImageNamespaceLocked:
		c.startImageNamespaceHealthCheck(key)
	}
	return nil
}

func (c *Controller) createRegistryNamespace(ctx context.Context, imageNamespace *businessv1.ImageNamespace) error {
	_, err := c.registryClient.Namespaces().Create(ctx, &registryv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"projectName": imageNamespace.Namespace,
			},
		},
		Spec: registryv1.NamespaceSpec{
			Name:        imageNamespace.Name,
			DisplayName: imageNamespace.Spec.DisplayName,
			TenantID:    imageNamespace.Spec.TenantID,
		}}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) persistUpdate(ctx context.Context, imageNamespace *businessv1.ImageNamespace) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.BusinessV1().ImageNamespaces(imageNamespace.Namespace).UpdateStatus(ctx, imageNamespace, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to imageNamespace that no longer exists",
				log.String("projectName", imageNamespace.Namespace),
				log.String("imageNamespaceName", imageNamespace.Name),
				log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to imageNamespace '%s/%s' that has been changed since we received it: %v",
				imageNamespace.Namespace, imageNamespace.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of imageNamespace '%s/%s/%s'",
			imageNamespace.Namespace, imageNamespace.Name, imageNamespace.Status.Phase),
			log.String("imageNamespaceName", imageNamespace.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}
