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
	"sort"
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
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

type ContextKey int

const (
	KeyLister ContextKey = iota
)

// SyncClsToAppController is responsible for performing actions dependent upon a cluster phase.
type SyncClsToAppController struct {
	queue     workqueue.RateLimitingInterface
	clsLister platformv1lister.ClusterLister
	clsSynced cache.InformerSynced
	appLister applicationv1lister.AppLister
	appSynced cache.InformerSynced

	log               log.Logger
	platformClient    platformversionedclient.PlatformV1Interface
	applicationClient applicationversionedclient.ApplicationV1Interface
}

// NewSyncClsToAppController creates a new Controller object.
func NewSyncClsToAppController(
	platformClient platformversionedclient.PlatformV1Interface,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	clsInformer platformv1informer.ClusterInformer,
	appInformer applicationv1informer.AppInformer,
	configuration clusterconfig.ClusterControllerConfiguration,
	finalizerToken platformv1.FinalizerName) *SyncClsToAppController {
	rand.Seed(time.Now().Unix())
	rateLimit := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(configuration.BucketRateLimiterLimit), configuration.BucketRateLimiterBurst)},
	)
	c := &SyncClsToAppController{
		queue: workqueue.NewNamedRateLimitingQueue(rateLimit, "cluster"),

		log:               log.WithName("SyncClsToAppController"),
		platformClient:    platformClient,
		applicationClient: applicationClient,
	}

	if platformClient != nil && platformClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("syncClsToApp_controller", platformClient.RESTClient().GetRateLimiter())
	}

	clsInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				UpdateFunc: func(oldObj, newObj interface{}) {
					old, ok1 := oldObj.(*platformv1.Cluster)
					cur, ok2 := newObj.(*platformv1.Cluster)
					if ok1 && ok2 && c.needsUpdate(old, cur) {
						go c.cleanupClusterApps(old, cur)
						c.enqueue(cur)
					}
				},
			},
			FilterFunc: func(obj interface{}) bool {
				cluster, ok := obj.(*platformv1.Cluster)
				if !ok {
					return false
				}
				provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
				if err != nil {
					return false
				}
				return provider.OnFilter(context.TODO(), cluster)
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

func (c *SyncClsToAppController) enqueue(obj *platformv1.Cluster) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

func (c *SyncClsToAppController) needsUpdate(old *platformv1.Cluster, new *platformv1.Cluster) bool {
	if !reflect.DeepEqual(old.Spec.ClusterApps, new.Spec.ClusterApps) && new.Status.Phase == platformv1.ClusterRunning {
		return true
	}

	if old.Status.Phase == platformv1.ClusterInitializing && new.Status.Phase == platformv1.ClusterRunning {
		return true
	}

	if old.Status.Phase == platformv1.ClusterFailed && new.Status.Phase == platformv1.ClusterRunning {
		return true
	}

	return false
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *SyncClsToAppController) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting syncClsToApp controller")
	defer log.Info("Shutting down syncClsToApp controller")

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
func (c *SyncClsToAppController) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *SyncClsToAppController) processNextWorkItem() bool {
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
func (c *SyncClsToAppController) sync(key string) error {
	ctx := c.log.WithValues("clusterToApps", key).WithContext(context.TODO())

	startTime := time.Now()
	defer func() {
		log.FromContext(ctx).Info("Finished syncing cluster to apps", "processTime", time.Since(startTime).String())
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	cluster, err := c.clsLister.Get(name)
	if apierrors.IsNotFound(err) {
		log.FromContext(ctx).Info("cluster has been deleted")
	}
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("unable to retrieve cluster %v from store: %v", key, err))
		return err
	}

	valueCtx := context.WithValue(ctx, KeyLister, &c.clsLister)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return c.syncClusterToApp(valueCtx, cluster)
	})
}

func (c *SyncClsToAppController) syncClusterToApp(ctx context.Context, cls *platformv1.Cluster) error {
	logger := log.FromContext(ctx)
	if c.applicationClient == nil {
		logger.Info("application client is nil, skip install apps")
		return nil
	}
	apps, err := c.applicationClient.Apps("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.targetCluster=%s", cls.Name),
	})
	if err != nil {
		return fmt.Errorf("get applications failed %v", err)
	}
	clusterApps := cls.Spec.ClusterApps
	sort.Sort(clusterApps)
	for i, clusterApp := range clusterApps {
		if clusterApp.App.Status.Phase == applicationv1.AppPhaseTerminating {
			cls.Spec.ClusterApps = append(cls.Spec.ClusterApps[:i], cls.Spec.ClusterApps[i+1:]...)
			_, err := c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
			logger.Errorf("cluster is: %v", cls.Spec.ClusterApps)
			return err
		}
		if c.isNeedUpdateApp(clusterApp, apps.Items) {
			if err := c.updateApplication(ctx, clusterApp); err != nil {
				logger.Errorf("update app failed: %v", err)
			}
			continue
		}
		clusterApp.App.Spec.TargetCluster = cls.Name
		err := c.installApplication(ctx, clusterApp)
		if err != nil {
			return fmt.Errorf("install application failed. %v, %v", clusterApp.App.Name, err)
		}
		logger.Infof("finish application installation %v", clusterApp.App.Name)
	}
	return nil
}

func (c *SyncClsToAppController) cleanupClusterApps(oldcls, newcls *platformv1.Cluster) {
	ctx := context.Background()
	for _, oldclsapp := range oldcls.Spec.ClusterApps {
		if oldclsapp.App.Status.Phase == applicationv1.AppPhaseTerminating {
			continue
		}
		if !newcls.Spec.ClusterApps.HasApp(oldclsapp.AppNamespace, oldclsapp.App.Spec.Name) {
			apps, err := c.applicationClient.Apps(oldclsapp.AppNamespace).List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.name=%s", oldclsapp.App.Spec.Name),
			})
			if err != nil {
				c.log.Errorf("list apps failed: %v", err)
			}
			if len(apps.Items) > 0 {
				err := c.applicationClient.Apps(oldclsapp.AppNamespace).Delete(ctx, apps.Items[0].Name, metav1.DeleteOptions{})
				if err != nil {
					c.log.Errorf("delete app %s failed: %v", apps.Items[0].Name, err)

				}

			}
		}
	}
}

func (c *SyncClsToAppController) isNeedUpdateApp(clusterApp platformv1.ClusterApp, installedApps []applicationv1.App) bool {
	for _, installedApp := range installedApps {
		if clusterApp.App.Spec.Name == installedApp.Spec.Name &&
			clusterApp.AppNamespace == installedApp.Namespace &&
			clusterApp.App.Spec.TargetCluster == installedApp.Spec.TargetCluster &&
			clusterApp.App.Spec.Chart.ChartVersion != installedApp.Spec.Chart.ChartVersion {
			return true
		}
	}
	return false
}

func (c *SyncClsToAppController) updateApplication(ctx context.Context, clusterApp platformv1.ClusterApp) error {
	if len(clusterApp.App.Name) == 0 {
		return nil
	}
	newApp := applicationv1.App(clusterApp.App)
	_, err := c.applicationClient.Apps(clusterApp.AppNamespace).Update(ctx, &newApp, metav1.UpdateOptions{})
	return err
}

func (c *SyncClsToAppController) installApplication(ctx context.Context, clusterApp platformv1.ClusterApp) error {
	app := applicationv1.App(clusterApp.App)
	_, err := c.applicationClient.Apps(clusterApp.AppNamespace).Create(ctx, &app, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create application failed %v,%v", clusterApp.App.Spec.Chart.ChartName, err)
	}
	return nil
}
