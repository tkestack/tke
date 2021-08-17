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

const (
	conditionTypeHealthCheck         = "HealthCheck"
	conditionTypeEnsureBootstrapApps = "EnsureBootstrapApps"
)

// BootstrapAppsController is responsible for performing actions dependent upon a cluster phase.
type BootstrapAppsController struct {
	queue     workqueue.RateLimitingInterface
	clsLister platformv1lister.ClusterLister
	clsSynced cache.InformerSynced
	appLister applicationv1lister.AppLister
	appSynced cache.InformerSynced

	log               log.Logger
	platformClient    platformversionedclient.PlatformV1Interface
	applicationClient applicationversionedclient.ApplicationV1Interface
}

// NewBootstrapAppsController creates a new Controller object.
func NewBootstrapAppsController(
	platformClient platformversionedclient.PlatformV1Interface,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	clsInformer platformv1informer.ClusterInformer,
	appInformer applicationv1informer.AppInformer,
	configuration clusterconfig.ClusterControllerConfiguration,
	finalizerToken platformv1.FinalizerName) *BootstrapAppsController {
	rand.Seed(time.Now().Unix())
	rateLimit := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(configuration.BucketRateLimiterLimit), configuration.BucketRateLimiterBurst)},
	)
	c := &BootstrapAppsController{
		queue: workqueue.NewNamedRateLimitingQueue(rateLimit, "cluster"),

		log:               log.WithName("BootstrapAppsController"),
		platformClient:    platformClient,
		applicationClient: applicationClient,
	}

	if platformClient != nil && platformClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("bootstrapApps_controller", platformClient.RESTClient().GetRateLimiter())
	}

	clsInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*platformv1.Cluster)
				cur, ok2 := newObj.(*platformv1.Cluster)
				if ok1 && ok2 && c.needsUpdate(old, cur) {
					c.enqueue(cur)
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

func (c *BootstrapAppsController) enqueue(obj *platformv1.Cluster) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

func (c *BootstrapAppsController) needsUpdate(old *platformv1.Cluster, new *platformv1.Cluster) bool {
	healthCheckDone := false
	for _, condition := range new.Status.Conditions {
		if condition.Type == conditionTypeHealthCheck && condition.Status == platformv1.ConditionTrue {
			healthCheckDone = true
		}
		if condition.Type == conditionTypeEnsureBootstrapApps && condition.Status == platformv1.ConditionTrue {
			return false
		}
	}

	return healthCheckDone
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *BootstrapAppsController) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting bootstrapApps controller")
	defer log.Info("Shutting down bootstrapApps controller")

	if ok := cache.WaitForCacheSync(stopCh, c.clsSynced); !ok {
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
func (c *BootstrapAppsController) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *BootstrapAppsController) processNextWorkItem() bool {
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

	utilruntime.HandleError(fmt.Errorf("error processing cluster %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// sync will sync the Cluster with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *BootstrapAppsController) sync(key string) error {
	ctx := c.log.WithValues("bootstrapApps", key).WithContext(context.TODO())

	startTime := time.Now()
	defer func() {
		log.FromContext(ctx).Info("Finished syncing bootstrap apps", "processTime", time.Since(startTime).String())
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	cluster, err := c.clsLister.Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.FromContext(ctx).Info("cluster has been deleted")
			return nil
		}
		utilruntime.HandleError(fmt.Errorf("unable to retrieve cluster %v from store: %v", key, err))
		return err
	}

	return c.syncBootstrapApps(ctx, cluster)
}

func (c *BootstrapAppsController) syncBootstrapApps(ctx context.Context, cls *platformv1.Cluster) error {
	logger := log.FromContext(ctx)
	conditon := platformv1.ClusterCondition{
		Type:   conditionTypeEnsureBootstrapApps,
		Status: platformv1.ConditionTrue,
	}

	for _, clusterApp := range cls.Spec.BootstrapApps {
		clusterApp.App.Spec.TargetCluster = cls.Name
		err := c.installApplication(ctx, clusterApp)
		if err != nil && apierrors.IsAlreadyExists(err) {
			conditon.Status = platformv1.ConditionFalse
			conditon.Reason = clusterApp.App.Name
			conditon.Message = err.Error()
			err := c.updateClsCondition(ctx, cls.Name, conditon)
			if err != nil {
				return fmt.Errorf("update cls conditon failed: %v", err)
			}
			return fmt.Errorf("install application failed. %v, %v", clusterApp.App.Name, err)
		}
		logger.Infof("finish application installation %v", clusterApp.App.Name)
	}
	return c.updateClsCondition(ctx, cls.Name, conditon)
}

func (c *BootstrapAppsController) installApplication(ctx context.Context, clusterApp platformv1.BootstrapApp) error {
	app := applicationv1.App{
		ObjectMeta: clusterApp.App.ObjectMeta,
		Spec:       clusterApp.App.Spec,
	}
	_, err := c.applicationClient.Apps(clusterApp.App.Namespace).Create(ctx, &app, metav1.CreateOptions{})

	return err
}

func (c *BootstrapAppsController) updateClsCondition(ctx context.Context, clsName string, condition platformv1.ClusterCondition) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cls, err := c.clsLister.Get(clsName)
		if err != nil {
			return err
		}
		cls.SetCondition(condition, false)
		_, err = c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
		return err
	})
}
