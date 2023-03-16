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

package app

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
	applicationv1 "tkestack.io/tke/api/application/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	applicationv1informer "tkestack.io/tke/api/client/informers/externalversions/application/v1"
	applicationv1lister "tkestack.io/tke/api/client/listers/application/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	"tkestack.io/tke/pkg/application/controller/app/action"
	"tkestack.io/tke/pkg/application/controller/app/deletion"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// appDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	appDeletionGracePeriod = 5 * time.Second
)

const (
	controllerName = "app-controller"
)

// Controller is responsible for performing actions dependent upon an app phase.
type Controller struct {
	client         clientset.Interface
	platformClient platformversionedclient.PlatformV1Interface
	repo           appconfig.RepoConfiguration
	cache          *applicationCache
	health         *applicationHealth
	queue          workqueue.RateLimitingInterface
	lister         applicationv1lister.AppLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// helper to delete all resources in the app when the app is deleted.
	appResourcesDeleter deletion.AppResourcesDeleterInterface
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	repo appconfig.RepoConfiguration,
	applicationInformer applicationv1informer.AppInformer,
	resyncPeriod time.Duration, finalizerToken applicationv1.FinalizerName,
) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:              client,
		platformClient:      platformClient,
		repo:                repo,
		cache:               &applicationCache{m: make(map[string]*cachedApp)},
		health:              &applicationHealth{applications: sets.NewString()},
		queue:               workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		appResourcesDeleter: deletion.NewAppResourcesDeleter(client.ApplicationV1(), platformClient, repo, finalizerToken, true),
	}

	if client != nil &&
		client.ApplicationV1().RESTClient() != nil &&
		!reflect.ValueOf(client.ApplicationV1().RESTClient()).IsNil() &&
		client.ApplicationV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("app_controller", client.ApplicationV1().RESTClient().GetRateLimiter())
	}

	applicationInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: controller.enqueue,
				UpdateFunc: func(oldObj, newObj interface{}) {
					old, ok1 := oldObj.(*applicationv1.App)
					cur, ok2 := newObj.(*applicationv1.App)
					if ok1 && ok2 && controller.needsUpdate(old, cur) {
						controller.enqueue(newObj)
					}
				},
				DeleteFunc: controller.enqueue,
			},
			FilterFunc: func(obj interface{}) bool {
				app, ok := obj.(*applicationv1.App)
				if !ok {
					return false
				}
				provider, err := applicationprovider.GetProvider(app)
				if err != nil {
					return true
				}
				return provider.OnFilter(context.TODO(), app)
			},
		},
		resyncPeriod,
	)
	controller.lister = applicationInformer.Lister()
	controller.listerSynced = applicationInformer.Informer().HasSynced

	return controller
}

// obj could be an *applicationv1.App, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, appDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *applicationv1.App, new *applicationv1.App) bool {
	if old.UID != new.UID {
		return true
	}

	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true
	}

	if !reflect.DeepEqual(old.Status, new.Status) {
		return true
	}

	if new.Status.Phase == applicationv1.AppPhaseSyncFailed ||
		new.Status.Phase == applicationv1.AppPhaseInstallFailed ||
		new.Status.Phase == applicationv1.AppPhaseUpgradFailed ||
		new.Status.Phase == applicationv1.AppPhaseSucceeded ||
		new.Status.Phase == applicationv1.AppPhaseTerminating {
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
	log.Info("Starting app controller")
	defer log.Info("Shutting down app controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for app caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of app objects.
// Each app can be in the queue at most once.
// The system ensures that no two workers can process
// the same app at the same time.
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

		// rather than wait for a full resync, re-add the app to the queue to be processed
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

// syncItem will sync the App with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// applications created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing app", log.String("app", key), log.Duration("processTime", time.Since(startTime)))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-c.stopCh:
			log.Info("stop ch", log.String("namespace", namespace), log.String("name", name))
			cancel()
			return
		case <-ctx.Done():
			log.Info("success done", log.String("namespace", namespace), log.String("name", name))
			return
		}
	}()
	// app holds the latest App info from apiserver
	app, err := c.lister.Apps(namespace).Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("App has been deleted. Attempting to cleanup resources",
			log.String("namespace", namespace),
			log.String("name", name))
		_ = c.processDeletion(key)
		return nil
	case err != nil:
		log.Warn("Unable to retrieve app from store",
			log.String("namespace", namespace),
			log.String("name", name), log.Err(err))
		return err
	default:
		if app.Status.Phase == applicationv1.AppPhaseTerminating {
			log.Info("App has been terminated. Attempting to cleanup resources",
				log.String("namespace", namespace),
				log.String("name", name))
			_ = c.processDeletion(key)
			err = c.appResourcesDeleter.Delete(ctx, namespace, name)
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationSyncFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			// If err is not nil, do not update object status when phase is Terminating.
			// DeletionTimestamp is not empty and object will be deleted when you request updateStatus
		} else if app.Status.Phase == "Isolated" {
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationSyncFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		} else {
			cachedApp := c.cache.getOrCreate(key)
			err = c.processUpdate(ctx, cachedApp, app, key)
		}
	}
	return err
}

func (c *Controller) processDeletion(key string) error {
	cachedApp, ok := c.cache.get(key)
	if !ok {
		log.Debug("App not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processDelete(cachedApp, key)
}

func (c *Controller) processDelete(cachedApp *cachedApp, key string) error {
	log.Info("App will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the app cache", log.String("name", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the app health cache", log.String("name", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) processUpdate(ctx context.Context, cachedApp *cachedApp, app *applicationv1.App, key string) error {
	if cachedApp.state != nil {
		// exist and the app name changed
		if cachedApp.state.UID != app.UID {
			if err := c.processDelete(cachedApp, key); err != nil {
				return err
			}
		}
	}
	// start update app if needed
	updated, err := c.handlePhase(ctx, key, cachedApp, app)
	if err != nil {
		log.Error("processUpdate failed",
			log.String("namespace", app.Namespace),
			log.String("name", app.Name),
			log.Err(err))
		return err
	}
	cachedApp.state = updated
	// Always update the cache upon success.
	c.cache.set(key, cachedApp)
	return nil
}

func (c *Controller) handlePhase(ctx context.Context, key string, cachedApp *cachedApp, app *applicationv1.App) (*applicationv1.App, error) {
	switch app.Status.Phase {
	case applicationv1.AppPhaseInstalling:
		return action.Install(ctx, c.client.ApplicationV1(), c.platformClient, app, c.repo, c.updateStatus)
	case applicationv1.AppPhaseUpgrading:
		// if only update status, generation won't change, and we won't upgrade
		if hasSynced(app) {
			newStatus := app.Status.DeepCopy()
			newStatus.Phase = applicationv1.AppPhaseSucceeded
			newStatus.Message = ""
			newStatus.Reason = ""
			newStatus.LastTransitionTime = metav1.Now()
			return c.updateStatus(ctx, app, &app.Status, newStatus)
		}
		return action.Upgrade(ctx, c.client.ApplicationV1(), c.platformClient, app, c.repo, c.updateStatus)
	case applicationv1.AppPhaseInstallFailed:
		return action.Install(ctx, c.client.ApplicationV1(), c.platformClient, app, c.repo, c.updateStatus)
	case applicationv1.AppPhaseSucceeded:
		c.startAppHealthCheck(ctx, key)
		// sync release status
		return c.syncAppFromRelease(ctx, cachedApp, app)
	case applicationv1.AppPhaseUpgradFailed:
		return action.Upgrade(ctx, c.client.ApplicationV1(), c.platformClient, app, c.repo, c.updateStatus)
	case applicationv1.AppPhaseRollingBack:
		if app.Status.RollbackRevision > 0 {
			return action.Rollback(ctx, c.client.ApplicationV1(), c.platformClient, app, c.repo, c.updateStatus)
		}
	case applicationv1.AppPhaseRolledBack:
		// sync release status
		return c.syncAppFromRelease(ctx, cachedApp, app)
	case applicationv1.AppPhaseRollbackFailed:
		break
	case applicationv1.AppPhaseSyncFailed:
		return c.syncAppFromRelease(ctx, cachedApp, app)
	default:
		break
	}
	return app, nil
}

func (c *Controller) syncAppFromRelease(ctx context.Context, cachedApp *cachedApp, app *applicationv1.App) (*applicationv1.App, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("syncAppFromRelease panic")
		}
	}()
	newStatus := app.Status.DeepCopy()
	rels, err := action.List(ctx, c.client.ApplicationV1(), c.platformClient, app)
	if err != nil {
		newStatus.Phase = applicationv1.AppPhaseSyncFailed
		newStatus.Message = "sync app failed"
		newStatus.Reason = err.Error()
		newStatus.LastTransitionTime = metav1.Now()
		metrics.GaugeApplicationSyncFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		return c.updateStatus(ctx, app, &app.Status, newStatus)
	}
	rel, found := helmutil.Filter(rels, app.Spec.TargetNamespace, app.Spec.Name)
	if !found {
		// release not found, reinstall for reconcile
		newStatus.Phase = applicationv1.AppPhaseInstalling
		newStatus.Message = "sync app failed"
		newStatus.Reason = fmt.Sprintf("release not found: %s/%s", app.Spec.TargetNamespace, app.Spec.Name)
		newStatus.LastTransitionTime = metav1.Now()
		metrics.GaugeApplicationSyncFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		return c.updateStatus(ctx, app, &app.Status, newStatus)
	}

	newStatus.Phase = applicationv1.AppPhaseSucceeded
	newStatus.Message = ""
	newStatus.Reason = ""
	newStatus.LastTransitionTime = metav1.Now()
	metrics.GaugeApplicationSyncFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)

	newStatus.ReleaseStatus = string(rel.Info.Status)
	newStatus.Revision = int64(rel.Version)
	if tspb := rel.Info.LastDeployed; !tspb.IsZero() {
		newStatus.ReleaseLastUpdated = metav1.NewTime(tspb.Time)
	}
	// updates the observed generation status of the App to the given generation.
	newStatus.ObservedGeneration = app.Generation
	// clean revision
	newStatus.RollbackRevision = 0
	if app.Status.Phase == applicationv1.AppPhaseRolledBack && app.Spec.Chart.ChartVersion != rel.Chart.Metadata.Version {
		newObj := app.DeepCopy()
		newObj.Spec.Chart.ChartVersion = rel.Chart.Metadata.Version
		newObj.Status = *newStatus
		_, err = c.client.ApplicationV1().Apps(app.Namespace).Update(ctx, newObj, metav1.UpdateOptions{})
		if err != nil {
			return app, fmt.Errorf("update chart version failed %v", err)
		}
		return app, err
	}
	if app.Status.Phase == applicationv1.AppPhaseSucceeded && hasSynced(app) {
		return app, nil
	}
	return c.updateStatus(ctx, app, &app.Status, newStatus)
}

func (c *Controller) updateStatus(ctx context.Context, app *applicationv1.App, previousStatus, newStatus *applicationv1.AppStatus) (*applicationv1.App, error) {
	if reflect.DeepEqual(previousStatus, newStatus) {
		return nil, nil
	}
	// Make a copy so we don't mutate the shared informer cache.
	newObj := app.DeepCopy()
	newObj.Status = *newStatus

	updated, err := c.client.ApplicationV1().Apps(newObj.Namespace).UpdateStatus(ctx, newObj, metav1.UpdateOptions{})
	if err == nil {
		return updated, nil
	}
	if errors.IsNotFound(err) {
		log.Info("Not persisting update to app that no longer exists",
			log.String("namespace", app.Namespace),
			log.String("name", app.Name),
			log.Err(err))
		return updated, nil
	}
	if errors.IsConflict(err) {
		return nil, fmt.Errorf("not persisting update to app '%s' that has been changed since we received it: %v",
			app.Name, err)
	}
	log.Warn(fmt.Sprintf("Failed to persist updated status of app '%s/%s'",
		app.Name, app.Status.Phase),
		log.String("namespace", app.Namespace),
		log.String("name", app.Name), log.Err(err))
	return nil, err
}

// hasSynced returns if the app has been processed by the controller.
func hasSynced(app *applicationv1.App) bool {
	return app.Status.ObservedGeneration >= app.Generation
}
