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

package namespace

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
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	businessv1informer "tkestack.io/tke/api/client/informers/externalversions/business/v1"
	businessv1lister "tkestack.io/tke/api/client/listers/business/v1"
	"tkestack.io/tke/pkg/business/controller/namespace/deletion"
	businessUtil "tkestack.io/tke/pkg/business/util"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	// namespaceDeletionGracePeriod is the time period to wait before processing a received channel event.
	// This allows time for the following to occur:
	// * lifecycle admission plugins on HA apiservers to also observe a channel
	//   deletion and prevent new objects from being created in the terminating channel
	// * non-leader etcd servers to observe last-minute object creations in a channel
	//   so this controller's cleanup can actually clean up all objects
	namespaceDeletionGracePeriod = 5 * time.Second
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second
)
const (
	controllerName = "namespace-controller"
)

// Controller is responsible for performing actions dependent upon a namespace phase.
type Controller struct {
	client         clientset.Interface
	platformClient platformversionedclient.PlatformV1Interface
	cache          *namespaceCache
	health         *namespaceHealth
	queue          workqueue.RateLimitingInterface
	lister         businessv1lister.NamespaceLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// helper to delete all resources in the namespace when the namespace is deleted.
	namespacedResourcesDeleter deletion.NamespacedResourcesDeleterInterface
}

// NewController creates a new Controller object.
func NewController(platformClient platformversionedclient.PlatformV1Interface, client clientset.Interface, namespaceInformer businessv1informer.NamespaceInformer, resyncPeriod time.Duration, finalizerToken v1.FinalizerName) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:                     client,
		platformClient:             platformClient,
		cache:                      &namespaceCache{m: make(map[string]*cachedNamespace)},
		health:                     &namespaceHealth{namespaces: sets.NewString()},
		queue:                      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		namespacedResourcesDeleter: deletion.NewNamespacedResourcesDeleter(platformClient, client.BusinessV1(), finalizerToken, true),
	}

	if client != nil && client.BusinessV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("namespace_controller", client.BusinessV1().RESTClient().GetRateLimiter())
	}

	namespaceInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueue,
			UpdateFunc: func(oldObj, newObj interface{}) {
				old, ok1 := oldObj.(*v1.Namespace)
				cur, ok2 := newObj.(*v1.Namespace)
				if ok1 && ok2 && controller.needsUpdate(old, cur) {
					controller.enqueue(newObj)
				}
			},
			DeleteFunc: controller.enqueue,
		},
		resyncPeriod,
	)
	controller.lister = namespaceInformer.Lister()
	controller.listerSynced = namespaceInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.Namespace, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueue(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.AddAfter(key, namespaceDeletionGracePeriod)
}

func (c *Controller) needsUpdate(old *v1.Namespace, new *v1.Namespace) bool {
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
	log.Info("Starting namespace controller")
	defer log.Info("Shutting down namespace controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		log.Error("Failed to wait for namespace caches to sync")
		return
	}

	c.stopCh = stopCh
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
}

// worker processes the queue of namespace objects.
// Each namespace can be in the queue at most once.
// The system ensures that no two workers can process
// the same namespace at the same time.
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

		// rather than wait for a full resync, re-add the namespace to the queue to be processed
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

// syncItem will sync the Namespace with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncItem(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing namespace", log.String("namespace", key), log.Duration("processTime", time.Since(startTime)))
	}()

	projectName, namespaceName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	var cachedNamespace *cachedNamespace
	// Namespace holds the latest Namespace info from apiserver
	namespace, err := c.lister.Namespaces(projectName).Get(namespaceName)
	switch {
	case errors.IsNotFound(err):
		log.Info("Namespace has been deleted. Attempting to cleanup resources", log.String("projectName", projectName), log.String("namespaceName", namespaceName))
		err = c.processDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve namespace from store", log.String("projectName", projectName), log.String("namespaceName", namespaceName), log.Err(err))
	default:
		if namespace.Status.Phase == v1.NamespacePending || namespace.Status.Phase == v1.NamespaceAvailable || namespace.Status.Phase == v1.NamespaceFailed {
			cachedNamespace = c.cache.getOrCreate(key)
			err = c.processUpdate(cachedNamespace, namespace, key)
		} else if namespace.Status.Phase == v1.NamespaceTerminating {
			log.Info("Namespace has been terminated. Attempting to cleanup resources", log.String("projectName", projectName), log.String("namespaceName", namespaceName))
			_ = c.processDeletion(key)
			err = c.namespacedResourcesDeleter.Delete(projectName, namespaceName)
		} else {
			log.Debug(fmt.Sprintf("Namespace %s status is %s, not to process", key, namespace.Status.Phase))
		}
	}
	return err
}

func (c *Controller) processDeletion(key string) error {
	cachedNamespace, ok := c.cache.get(key)
	if !ok {
		log.Debug("Namespace not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processDelete(cachedNamespace, key)
}

func (c *Controller) processDelete(cachedNamespace *cachedNamespace, key string) error {
	log.Info("Namespace will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the namespace cache", log.String("name", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the namespace health cache", log.String("name", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) processUpdate(cachedNamespace *cachedNamespace, namespace *v1.Namespace, key string) error {
	if cachedNamespace.state != nil {
		// exist and the namespace name changed
		if cachedNamespace.state.UID != namespace.UID {
			if err := c.processDelete(cachedNamespace, key); err != nil {
				return err
			}
		}
	}
	// start update machine if needed
	err := c.handlePhase(key, cachedNamespace, namespace)
	if err != nil {
		return err
	}
	cachedNamespace.state = namespace
	// Always update the cache upon success.
	c.cache.set(key, cachedNamespace)
	return nil
}

func (c *Controller) handlePhase(key string, cachedNamespace *cachedNamespace, namespace *v1.Namespace) error {
	switch namespace.Status.Phase {
	case v1.NamespacePending:
		if err := c.calculateProjectUsed(cachedNamespace, namespace); err != nil {
			return err
		}
		if err := c.createNamespaceOnCluster(namespace); err != nil {
			namespace.Status.Phase = v1.NamespaceFailed
			namespace.Status.Message = "CreateNamespaceOnClusterFailed"
			namespace.Status.Reason = err.Error()
			namespace.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(namespace)
		}
		namespace.Status.Phase = v1.NamespaceAvailable
		namespace.Status.Message = ""
		namespace.Status.Reason = ""
		namespace.Status.LastTransitionTime = metav1.Now()
		return c.persistUpdate(namespace)
	case v1.NamespaceAvailable, v1.NamespaceFailed:
		c.startNamespaceHealthCheck(key)
	}
	return nil
}

func (c *Controller) calculateProjectUsed(cachedNamespace *cachedNamespace, namespace *v1.Namespace) error {
	project, err := c.client.BusinessV1().Projects().Get(namespace.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get the project", log.String("projectName", namespace.ObjectMeta.Namespace), log.Err(err))
		return err
	}
	calculatedNamespaceNames := sets.NewString(project.Status.CalculatedNamespaces...)
	if !calculatedNamespaceNames.Has(project.ObjectMeta.Name) {
		project.Status.CalculatedNamespaces = append(project.Status.CalculatedNamespaces, namespace.ObjectMeta.Name)
		if project.Status.Clusters == nil {
			project.Status.Clusters = make(v1.ClusterUsed)
		}
		businessUtil.AddClusterHardToUsed(&project.Status.Clusters,
			v1.ClusterHard{
				namespace.Spec.ClusterName: v1.HardQuantity{
					Hard: namespace.Spec.Hard,
				},
			})
		return c.persistUpdateProject(project)
	}
	if cachedNamespace.state != nil && !reflect.DeepEqual(cachedNamespace.state.Spec.Hard, namespace.Spec.Hard) {
		if project.Status.Clusters == nil {
			project.Status.Clusters = make(v1.ClusterUsed)
		}
		// sub old
		businessUtil.SubClusterHardFromUsed(&project.Status.Clusters,
			v1.ClusterHard{
				namespace.Spec.ClusterName: v1.HardQuantity{
					Hard: cachedNamespace.state.Spec.Hard,
				},
			})
		// add new
		businessUtil.AddClusterHardToUsed(&project.Status.Clusters,
			v1.ClusterHard{
				namespace.Spec.ClusterName: v1.HardQuantity{
					Hard: namespace.Spec.Hard,
				},
			})
		return c.persistUpdateProject(project)
	}
	return nil
}

func (c *Controller) createNamespaceOnCluster(namespace *v1.Namespace) error {
	kubeClient, err := util.BuildExternalClientSetWithName(c.platformClient, namespace.Spec.ClusterName)
	if err != nil {
		log.Error("Failed to create the kubernetes client", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
		return err
	}
	if err := createNamespaceOnCluster(kubeClient, namespace); err != nil {
		return err
	}
	return createResourceQuotaOnCluster(kubeClient, namespace)
}

func (c *Controller) persistUpdate(namespace *v1.Namespace) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.BusinessV1().Namespaces(namespace.ObjectMeta.Namespace).UpdateStatus(namespace)
		if err == nil {
			return nil
		}
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to namespace that no longer exists", log.String("projectName", namespace.ObjectMeta.Namespace), log.String("namespaceName", namespace.ObjectMeta.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to namespace '%s/%s' that has been changed since we received it: %v", namespace.ObjectMeta.Namespace, namespace.ObjectMeta.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of namespace '%s/%s/%s'", namespace.ObjectMeta.Namespace, namespace.ObjectMeta.Name, namespace.Status.Phase), log.String("namespaceName", namespace.ObjectMeta.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}

func (c *Controller) persistUpdateProject(project *v1.Project) error {
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
		log.Warn(fmt.Sprintf("Failed to persist updated status of projects '%s/%s'", project.ObjectMeta.Name, project.Status.Phase), log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		time.Sleep(clientRetryInterval)
	}
	return err
}
