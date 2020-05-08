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

package platform

import (
	"context"
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
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessv1informer "tkestack.io/tke/api/client/informers/externalversions/business/v1"
	businessv1lister "tkestack.io/tke/api/client/listers/business/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	platformEnqueueGracePeriod = 1 * time.Second
)

const (
	controllerName = "platform-controller"
)

// Controller is responsible for performing actions dependent upon a platform phase.
type Controller struct {
	client       clientset.Interface
	queue        workqueue.RateLimitingInterface
	lister       businessv1lister.PlatformLister
	listerSynced cache.InformerSynced
	authClient   authversionedclient.AuthV1Interface
}

// NewController creates a new platform object.
func NewController(client clientset.Interface, authClient authversionedclient.AuthV1Interface, platformInformer businessv1informer.PlatformInformer,
	resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:     client,
		authClient: authClient,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("platform_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	platformInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.Platform)
				cur, ok2 := newObj.(*v1.Platform)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
		},
		resyncPeriod,
	)
	controller.lister = platformInformer.Lister()
	controller.listerSynced = platformInformer.Informer().HasSynced
	return controller
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, platformEnqueueGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.Platform, new *v1.Platform) bool {
	if old.UID != new.UID {
		return true
	}

	if !reflect.DeepEqual(old.Spec, new.Spec) {
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
	log.Info("Starting platform controller")
	defer log.Info("Shutting down platform controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for platform caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of platform objects.
// Each platform can be in the queue at most once.
// The system ensures that no two workers can process
// the same platform at the same time.
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

		// rather than wait for a full resync, re-add the platform to the queue to be processed
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

// syncItem will sync the platform with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// platform created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing platform", log.String("platformName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	platform, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Platform has been deleted.", log.String("platformName", key))
	case err != nil:
		log.Warn("Unable to retrieve platform from store", log.String("platformName", key), log.Err(err))
	default:
		log.Info("Platform has been updated. Attempting to update identity provider", log.String("platformName", key))
		err = c.process(platform)
	}
	return err
}

func (c *Controller) process(platform *v1.Platform) error {
	tenantID := platform.Spec.TenantID
	if len(tenantID) == 0 {
		log.Warn("TenantID is empty", log.String("platform", platform.Name))
		return nil
	}

	idp, err := c.authClient.IdentityProviders().Get(context.Background(), tenantID, metav1.GetOptions{})
	if err != nil {
		log.Error("Get identity provider for tenant failed", log.String("tennantID", tenantID), log.Err(err))
		return err
	}

	if !reflect.DeepEqual(platform.Spec.Administrators, idp.Spec.Administrators) {
		log.Info("Attempting update identity provider for tenant with new administrators", log.String("tenant", tenantID),
			log.Strings("administrators", platform.Spec.Administrators))
		idp.Spec.Administrators = platform.Spec.Administrators
		_, err = c.authClient.IdentityProviders().Update(context.Background(), idp, metav1.UpdateOptions{})
	}

	return err
}
