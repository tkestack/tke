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

package project

import (
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
	v1 "tkestack.io/tke/api/business/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	businessv1informer "tkestack.io/tke/api/client/informers/externalversions/business/v1"
	businessv1lister "tkestack.io/tke/api/client/listers/business/v1"
	"tkestack.io/tke/pkg/business/controller/project/deletion"
	businessUtil "tkestack.io/tke/pkg/business/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// projectDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	projectDeletionGracePeriod = 5 * time.Second
)

const (
	controllerName = "project-controller"
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)

// Controller is responsible for performing actions dependent upon a project phase.
type Controller struct {
	client       clientset.Interface
	cache        *projectCache
	queue        workqueue.RateLimitingInterface
	lister       businessv1lister.ProjectLister
	listerSynced cache.InformerSynced
	// helper to delete all resources in the project when the project is deleted.
	projectedResourcesDeleter deletion.ProjectedResourcesDeleterInterface
}

// NewController creates a new Project object.
func NewController(client clientset.Interface, projectInformer businessv1informer.ProjectInformer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                    client,
		cache:                     &projectCache{m: make(map[string]*cachedProject)},
		queue:                     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		projectedResourcesDeleter: deletion.NewProjectedResourcesDeleter(client.BusinessV1().Projects(), client.BusinessV1(), finalizerToken, true),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("project_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	projectInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.Project)
				cur, ok2 := newObj.(*v1.Project)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
		},
		resyncPeriod,
	)
	controller.lister = projectInformer.Lister()
	controller.listerSynced = projectInformer.Informer().HasSynced
	return controller
}

// obj could be an *v1.Project, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.AddAfter(key, projectDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.Project, new *v1.Project) bool {
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
	log.Info("Starting project controller")
	defer log.Info("Shutting down project controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for project caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of project objects.
// Each project can be in the queue at most once.
// The system ensures that no two workers can process
// the same project at the same time.
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

		// rather than wait for a full resync, re-add the project to the queue to be processed
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

// syncItem will sync the Project with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// project created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()

	defer func() {
		log.Info("Finished syncing project", log.String("projectName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	var cachedProject *cachedProject
	// Project holds the latest Project info from apiserver
	project, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Project has been deleted. Attempting to cleanup resources", log.String("projectName", key))
		err = c.processDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve project from store", log.String("projectName", key), log.Err(err))
	default:
		if project.Status.Phase == v1.ProjectActive {
			cachedProject = c.cache.getOrCreate(key)
			err = c.processUpdate(cachedProject, project, key)
		} else if project.Status.Phase == v1.ProjectTerminating {
			log.Info("Project has been terminated. Attempting to cleanup resources", log.String("projectName", key))
			_ = c.processDeletion(key)
			err = c.projectedResourcesDeleter.Delete(key)
		} else {
			log.Debug(fmt.Sprintf("Project %s status is %s, not to process", key, project.Status.Phase))
		}
	}
	return err
}

func (c *Controller) processUpdate(cachedProject *cachedProject, project *v1.Project, key string) error {
	if cachedProject.state != nil {
		// exist and the project name changed
		if cachedProject.state.UID != project.UID {
			if err := c.processDelete(cachedProject, key); err != nil {
				return err
			}
		}
	}
	// start update machine if needed
	err := c.handlePhase(key, cachedProject, project)
	if err != nil {
		return err
	}
	cachedProject.state = project
	// Always update the cache upon success.
	c.cache.set(key, cachedProject)
	return nil
}

func (c *Controller) processDeletion(key string) error {
	cachedProject, ok := c.cache.get(key)
	if !ok {
		log.Debug("Project not in cache even though the watcher thought it was. Ignoring the deletion", log.String("projectName", key))
		return nil
	}
	return c.processDelete(cachedProject, key)
}

func (c *Controller) processDelete(cachedProject *cachedProject, key string) error {
	log.Info("Project will be dropped", log.String("projectName", key))

	if c.cache.Exist(key) {
		log.Info("Delete the project cache", log.String("projectName", key))
		c.cache.delete(key)
	}

	return nil
}

func (c *Controller) handlePhase(key string, cachedProject *cachedProject, project *v1.Project) error {
	if project.Spec.ParentProjectName != "" {
		parentProject, err := c.client.BusinessV1().Projects().Get(project.Spec.ParentProjectName, metav1.GetOptions{})
		if err != nil {
			log.Error("Failed to get the parent project", log.String("projectName", key), log.Err(err))
			return err
		}
		calculatedChildProjectNames := sets.NewString(parentProject.Status.CalculatedChildProjects...)
		if !calculatedChildProjectNames.Has(project.ObjectMeta.Name) {
			parentProject.Status.CalculatedChildProjects = append(parentProject.Status.CalculatedChildProjects, project.ObjectMeta.Name)
			if parentProject.Status.Clusters == nil {
				parentProject.Status.Clusters = make(v1.ClusterUsed)
			}
			businessUtil.AddClusterHardToUsed(&parentProject.Status.Clusters, project.Spec.Clusters)
			return c.persistUpdate(parentProject)
		}
		if cachedProject.state != nil && !reflect.DeepEqual(cachedProject.state.Spec.Clusters, project.Spec.Clusters) {
			if parentProject.Status.Clusters == nil {
				parentProject.Status.Clusters = make(v1.ClusterUsed)
			}
			// sub old
			businessUtil.SubClusterHardFromUsed(&parentProject.Status.Clusters, cachedProject.state.Spec.Clusters)
			// add new
			businessUtil.AddClusterHardToUsed(&parentProject.Status.Clusters, project.Spec.Clusters)
			return c.persistUpdate(parentProject)
		}
	}
	return nil
}

func (c *Controller) persistUpdate(project *v1.Project) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.BusinessV1().Projects().UpdateStatus(project)
		if err == nil {
			return nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to projects that no longer exists", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to projects '%s' that has been changed since we received it: %v", project.ObjectMeta.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of project %s", project.ObjectMeta.Name), log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}
