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

package channel

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	notifyinformers "tkestack.io/tke/api/client/informers/externalversions/notify/v1"
	notifylisters "tkestack.io/tke/api/client/listers/notify/v1"
	"tkestack.io/tke/api/notify/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/notify/controller/channel/deletion"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// channelDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	channelDeletionGracePeriod = 5 * time.Second
)

// Controller is responsible for performing actions dependent upon a channel phase
type Controller struct {
	// lister that can list channels from a shared cache
	lister notifylisters.ChannelLister
	// returns true when the channel cache is ready
	listerSynced cache.InformerSynced
	// channels that have been queued up for processing by workers
	queue workqueue.RateLimitingInterface
	// helper to delete all resources in the channel when the channel is deleted.
	channeldResourcesDeleter deletion.ChanneledResourcesDeleterInterface
}

// NewController creates a new Controller
func NewController(
	client clientset.Interface,
	channelInformer notifyinformers.ChannelInformer,
	resyncPeriod time.Duration,
	finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	channelController := &Controller{
		queue:                    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "channel"),
		channeldResourcesDeleter: deletion.NewChanneledResourcesDeleter(client.NotifyV1().Channels(), client.NotifyV1(), finalizerToken, true),
	}

	if client != nil && client.NotifyV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("channel_controller", client.NotifyV1().RESTClient().GetRateLimiter())
	}

	// configure the channel informer event handlers
	channelInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				channel := obj.(*v1.Channel)
				channelController.enqueueChannel(channel)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				channel := newObj.(*v1.Channel)
				channelController.enqueueChannel(channel)
			},
		},
		resyncPeriod,
	)
	channelController.lister = channelInformer.Lister()
	channelController.listerSynced = channelInformer.Informer().HasSynced

	return channelController
}

// enqueueChannel adds an object to the controller work queue
// obj could be an *v1.Channel, or a DeletionFinalStateUnknown item.
func (nm *Controller) enqueueChannel(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}

	channel := obj.(*v1.Channel)
	// don't queue if we aren't deleted
	if channel.DeletionTimestamp == nil || channel.DeletionTimestamp.IsZero() {
		return
	}

	// delay processing channel events to allow HA api servers to observe channel deletion,
	// and HA etcd servers to observe last minute object creations inside the channel
	nm.queue.AddAfter(key, channelDeletionGracePeriod)
}

// worker processes the queue of channel objects.
// Each channel can be in the queue at most once.
// The system ensures that no two workers can process
// the same channel at the same time.
func (nm *Controller) worker() {
	workFunc := func() bool {
		key, quit := nm.queue.Get()
		if quit {
			return true
		}
		defer nm.queue.Done(key)

		err := nm.syncChannelFromKey(key.(string))
		if err == nil {
			// no error, forget this entry and return
			nm.queue.Forget(key)
			return false
		}

		// rather than wait for a full resync, re-add the channel to the queue to be processed
		nm.queue.AddRateLimited(key)
		utilruntime.HandleError(err)
		return false
	}

	for {
		quit := workFunc()

		if quit {
			return
		}
	}
}

// syncChannelFromKey looks for a channel with the specified key in its store and synchronizes it
func (nm *Controller) syncChannelFromKey(key string) (err error) {
	startTime := time.Now()
	defer func() {
		log.Infof("Finished syncing channel %q (%v)", key, time.Since(startTime))
	}()

	channel, err := nm.lister.Get(key)
	if errors.IsNotFound(err) {
		log.Infof("Channel has been deleted %v", key)
		return nil
	}
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("unable to retrieve channel %v from store: %v", key, err))
		return err
	}
	return nm.channeldResourcesDeleter.Delete(channel.Name)
}

// Run starts observing the system with the specified number of workers.
func (nm *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer nm.queue.ShutDown()

	log.Infof("Starting channel controller")
	defer log.Infof("Shutting down channel controller")

	if !controllerutil.WaitForCacheSync("channel", stopCh, nm.listerSynced) {
		return
	}

	log.Info("Starting workers of channel controller")
	for i := 0; i < workers; i++ {
		go wait.Until(nm.worker, time.Second, stopCh)
	}
	<-stopCh
}
