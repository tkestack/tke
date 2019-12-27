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

package localidentity

import (
	"fmt"
	"reflect"
	"time"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	v1 "tkestack.io/tke/api/auth/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/localidentity/deletion"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// localIdentityDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	localIdentityDeletionGracePeriod = 5 * time.Second

	controllerName = "localidentity-controller"
)

// Controller is responsible for performing actions dependent upon a policy phase.
type Controller struct {
	client                    clientset.Interface
	queue                     workqueue.RateLimitingInterface
	localIdentityLister       authv1lister.LocalIdentityLister
	localIdentityListerSynced cache.InformerSynced
	ruleLister                authv1lister.RuleLister
	ruleListerSynced          cache.InformerSynced
	// helper to delete all resources in the localIdentity when the localIdentity is deleted.
	localIdentityedResourcesDeleter deletion.LocalIdentitiedResourcesDeleterInterface
	enforcer                        *casbin.SyncedEnforcer
}

// NewController creates a new localIdentity object.
func NewController(client clientset.Interface, localIdentityInformer authv1informer.LocalIdentityInformer, ruleInformer authv1informer.RuleInformer, enforcer *casbin.SyncedEnforcer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                          client,
		queue:                           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		enforcer:                        enforcer,
		localIdentityedResourcesDeleter: deletion.NewLocalIdentitiedResourcesDeleter(client.AuthV1().LocalIdentities(), client.AuthV1(), enforcer, finalizerToken, true),
	}

	if client != nil && client.AuthV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("localidentity_controller", client.AuthV1().RESTClient().GetRateLimiter())
	}

	localIdentityInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.LocalIdentity)
				cur, ok2 := newObj.(*v1.LocalIdentity)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					log.Info("Update enqueue")
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.localIdentityLister = localIdentityInformer.Lister()
	controller.localIdentityListerSynced = localIdentityInformer.Informer().HasSynced

	controller.ruleLister = ruleInformer.Lister()
	controller.ruleListerSynced = ruleInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.localIdentity, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, localIdentityDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.LocalIdentity, new *v1.LocalIdentity) bool {
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
	log.Info("Starting localIdentity controller")
	defer log.Info("Shutting down localIdentity controller")

	if ok := cache.WaitForCacheSync(stopCh, c.localIdentityListerSynced, c.ruleListerSynced); !ok {
		log.Error("Failed to wait for localIdentity caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of localIdentity objects.
// Each localIdentity can be in the queue at most once.
// The system ensures that no two workers can process
// the same localIdentity at the same time.
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

		// rather than wait for a full resync, re-add the localIdentity to the queue to be processed
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

// syncItem will sync the localIdentity with the given key if it has had
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing localIdentity", log.String("localIdentity", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	localIdentity, err := c.localIdentityLister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Infof("LocalIdentity has been deleted %v", key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve localIdentity from store", log.String("localIdentity name", key), log.Err(err))
	default:
		// Only check deleting localIdentity for now
		if localIdentity.Status.Phase == v1.LocalIdentityDeleting {
			err = c.localIdentityedResourcesDeleter.Delete(key)
		}

		log.Debug("Handle localIdentity", log.Any("localIdentity", localIdentity))
	}
	return err
}
