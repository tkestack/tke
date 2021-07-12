/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package clusterapp

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"golang.org/x/time/rate"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	applicationv1informer "tkestack.io/tke/api/client/informers/externalversions/application/v1"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	applicationv1lister "tkestack.io/tke/api/client/listers/application/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	clusterconfig "tkestack.io/tke/pkg/platform/controller/cluster/config"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

// SyncAppToClsController is responsible for performing actions dependent upon a cluster phase.
type SyncAppToClsController struct {
	queue     workqueue.RateLimitingInterface
	clsLister platformv1lister.ClusterLister
	clsSynced cache.InformerSynced
	appLister applicationv1lister.AppLister
	appSynced cache.InformerSynced

	log               log.Logger
	platformClient    platformversionedclient.PlatformV1Interface
	applicationClient applicationversionedclient.ApplicationV1Interface
}

// NewSyncAppToClsController creates a new Controller object.
func NewSyncAppToClsController(
	platformClient platformversionedclient.PlatformV1Interface,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	clsInformer platformv1informer.ClusterInformer,
	appInformer applicationv1informer.AppInformer,
	configuration clusterconfig.ClusterControllerConfiguration,
	finalizerToken platformv1.FinalizerName) *SyncAppToClsController {
	rand.Seed(time.Now().Unix())
	rateLimit := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(configuration.BucketRateLimiterLimit), configuration.BucketRateLimiterBurst)},
	)
	c := &SyncAppToClsController{
		queue: workqueue.NewNamedRateLimitingQueue(rateLimit, "application"),

		log:               log.WithName("SyncAppToClsController"),
		platformClient:    platformClient,
		applicationClient: applicationClient,
	}

	if applicationClient != nil && applicationClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("syncAppToCls_controller", platformClient.RESTClient().GetRateLimiter())
	}

	appInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*applicationv1.App)
				cur, ok2 := newObj.(*applicationv1.App)
				if ok1 && ok2 && c.needsUpdate(old, cur) {
					c.enqueue(newObj)
				}
			},
		},
		configuration.ClusterSyncPeriod,
	)

	c.clsLister = clsInformer.Lister()
	c.clsSynced = clsInformer.Informer().HasSynced
	c.appLister = appInformer.Lister()
	c.appSynced = appInformer.Informer().HasSynced

	return c
}

func (c *SyncAppToClsController) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *SyncAppToClsController) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting syncAppToCls controller")
	defer log.Info("Shutting down syncAppToCls controller")

	if ok := cache.WaitForCacheSync(stopCh, c.appSynced); !ok {
		return fmt.Errorf("failed to wait for cluster caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh

	return nil
}

// worker processes the queue of persistent event objects.
// Each cluster can be in the queue at most once.
// The system ensures that no two workers can process
// the same namespace at the same time.
func (c *SyncAppToClsController) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *SyncAppToClsController) processNextWorkItem() bool {
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

	utilruntime.HandleError(fmt.Errorf("error processing sync app to cluster %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// sync will sync the Cluster with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *SyncAppToClsController) sync(key string) error {
	ctx := c.log.WithValues("appToCluster", key).WithContext(context.TODO())

	startTime := time.Now()
	defer func() {
		log.FromContext(ctx).Info("Finished syncing app to cluster", "processTime", time.Since(startTime).String())
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	app, err := c.appLister.Apps(namespace).Get(name)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("unable to retrieve cluster %v from store: %v", key, err))
			return err
		}
		log.FromContext(ctx).Info("app has been deleted")
	}

	valueCtx := context.WithValue(ctx, KeyLister, c.appLister)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return c.syncAppToCluster(valueCtx, app)
	})
}

func (c *SyncAppToClsController) syncAppToCluster(ctx context.Context, obj interface{}) error {
	app := obj.(*applicationv1.App)
	err := c.addOrUpdateAppToCluster(ctx, app)
	if err != nil {
		return fmt.Errorf("add/update app to cluster %s failed: %v", app.Spec.TargetCluster, err)
	}
	return nil
}

func (c *SyncAppToClsController) addOrUpdateAppToCluster(ctx context.Context, app *applicationv1.App) error {
	cls, err := c.platformClient.Clusters().Get(ctx, app.Spec.TargetCluster, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	clsApp := platformv1.ClusterApp{
		AppNamespace: app.Namespace,
		App: platformv1.App{
			ObjectMeta: app.ObjectMeta,
			Spec:       app.Spec,
			Status:     app.Status,
		},
	}
	for i, item := range cls.Spec.ClusterApps {
		if item.AppNamespace == app.Namespace && item.App.Spec.Name == app.Spec.Name {
			if !reflect.DeepEqual(item.App.Spec, app.Spec) || !reflect.DeepEqual(item.App.Status, app.Status) {
				cls.Spec.ClusterApps[i] = clsApp
				_, err = c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
				return err
			}
			return nil
		}
	}

	// if app is terminating, no need add it to cluster apps
	if app.Status.Phase == applicationv1.AppPhaseTerminating {
		return nil
	}

	cls.Spec.ClusterApps = append(cls.Spec.ClusterApps, clsApp)

	_, err = c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
	return err
}

func (c *SyncAppToClsController) needsUpdate(old *applicationv1.App, new *applicationv1.App) bool {
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
