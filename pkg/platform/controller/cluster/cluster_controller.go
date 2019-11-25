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

package cluster

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/controller/cluster/deletion"
	"tkestack.io/tke/pkg/platform/provider"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	clusterClientRetryCount    = 5
	clusterClientRetryInterval = 5 * time.Second

	reasonFailedInit   = "FailedInit"
	reasonFailedUpdate = "FailedUpdate"
)

// Controller is responsible for performing actions dependent upon a cluster phase.
type Controller struct {
	client           clientset.Interface
	cache            *clusterCache
	health           *clusterHealth
	queue            workqueue.RateLimitingInterface
	lister           platformv1lister.ClusterLister
	listerSynced     cache.InformerSynced
	stopCh           <-chan struct{}
	clusterProviders *sync.Map
	// helper to delete all resources in the cluster when the cluster is deleted.
	clusterDeleter deletion.ClusterDeleterInterface
}

// obj could be an *platformv1.Cluster, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueCluster(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *platformv1.Cluster, new *platformv1.Cluster) bool {
	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true
	}

	if !reflect.DeepEqual(old.Status, new.Status) {
		return true
	}

	return false
}

// NewController creates a new Controller object.
func NewController(
	client clientset.Interface,
	clusterInformer platformv1informer.ClusterInformer,
	resyncPeriod time.Duration,
	finalizerToken platformv1.FinalizerName,
	clusterProviders *sync.Map) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:           client,
		cache:            &clusterCache{clusterMap: make(map[string]*cachedCluster)},
		health:           &clusterHealth{clusterMap: make(map[string]*platformv1.Cluster)},
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cluster"),
		clusterProviders: clusterProviders,
		clusterDeleter: deletion.NewClusterDeleter(client.PlatformV1().Clusters(),
			client.PlatformV1(),
			clusterProviders,
			finalizerToken,
			true),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("cluster_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the namespace informer event handlers
	clusterInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueCluster,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCluster, ok1 := oldObj.(*platformv1.Cluster)
				curCluster, ok2 := newObj.(*platformv1.Cluster)
				if ok1 && ok2 && controller.needsUpdate(oldCluster, curCluster) {
					controller.enqueueCluster(newObj)
				} else {
					log.Debug("Update new cluster not to add", log.String("clusterName", curCluster.Name), log.String("resourceversion", curCluster.ResourceVersion), log.String("old-resourceversion", oldCluster.ResourceVersion), log.String("cur-resourceversion", curCluster.ResourceVersion))
				}
			},
		},
		resyncPeriod,
	)
	controller.lister = clusterInformer.Lister()
	controller.listerSynced = clusterInformer.Informer().HasSynced

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting cluster controller")
	defer log.Info("Shutting down cluster controller")

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for cluster caches to sync")
	}

	c.stopCh = stopCh

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
func (c *Controller) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncCluster(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing cluster %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncCluster will sync the Cluster with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncCluster(key string) error {
	startTime := time.Now()
	var cachedCluster *cachedCluster
	defer func() {
		log.Info("Finished syncing cluster", log.String("clusterName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// cluster holds the latest cluster info from apiserver
	cluster, err := c.lister.Get(name)

	switch {
	case errors.IsNotFound(err):
		log.Info("Cluster has been deleted. Attempting to cleanup resources", log.String("clusterName", key))
		err = c.processClusterDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve cluster from store", log.String("clusterName", key), log.Err(err))
	default:
		if (cluster.Status.Phase == platformv1.ClusterRunning) || (cluster.Status.Phase == platformv1.ClusterFailed) || (cluster.Status.Phase == platformv1.ClusterInitializing) {
			cachedCluster = c.cache.getOrCreate(key)
			err = c.processClusterUpdate(cachedCluster, cluster, key)
		} else if cluster.Status.Phase == platformv1.ClusterTerminating {
			log.Info("Cluster has been terminated. Attempting to cleanup resources", log.String("clusterName", key))
			_ = c.processClusterDeletion(key)
			err = c.clusterDeleter.Delete(key)
		} else {
			log.Debug(fmt.Sprintf("Cluster %s status is %s, not to process", key, cluster.Status.Phase), log.String("clusterName", key))
		}
	}
	return err
}

func (c *Controller) processClusterUpdate(cachedCluster *cachedCluster, cluster *platformv1.Cluster, key string) error {
	if cachedCluster.state != nil {
		if cachedCluster.state.UID != cluster.UID {
			err := c.processClusterDelete(key)
			if err != nil {
				return err
			}
		}
	}

	// start update cluster if needed
	err := c.handlePhase(key, cachedCluster, cluster)
	if err != nil {
		return err
	}

	cachedCluster.state = cluster
	// Always update the cache upon success.
	c.cache.set(key, cachedCluster)

	return nil
}

func (c *Controller) processClusterDeletion(key string) error {
	_, ok := c.cache.get(key)
	if !ok {
		log.Debug("Cluster not in cache even though the watcher thought it was. Ignoring the deletion", log.String("clusterName", key))
		return nil
	}
	return c.processClusterDelete(key)
}

func (c *Controller) processClusterDelete(key string) error {
	log.Info("Cluster will be dropped", log.String("clusterName", key))

	if c.cache.Exist(key) {
		log.Info("Delete the cluster cache", log.String("clusterName", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("Delete the cluster health cache", log.String("clusterName", key))
		c.health.Del(key)
	}

	return nil
}

func (c *Controller) handlePhase(key string, cachedCluster *cachedCluster, cluster *platformv1.Cluster) error {
	var err error
	switch cluster.Status.Phase {
	case platformv1.ClusterInitializing:
		err = c.doInitializing(cluster)
	case platformv1.ClusterRunning:
		c.ensureHealthCheck(key, cluster)
		err = c.doUpdate(cluster)
		log.Info("update cluster", log.String("clusterName", cluster.Name), log.Err(err))
	case platformv1.ClusterFailed:
		c.ensureHealthCheck(key, cluster)
	}

	return err
}

func (c *Controller) addOrUpdateCondition(cluster *platformv1.Cluster, newCondition platformv1.ClusterCondition) {
	var conditions []platformv1.ClusterCondition
	exist := false
	for _, condition := range cluster.Status.Conditions {
		if condition.Type == newCondition.Type {
			exist = true
			if newCondition.Status != condition.Status {
				condition.Status = newCondition.Status
			}
			if newCondition.Message != condition.Message {
				condition.Message = newCondition.Message
			}
			if newCondition.Reason != condition.Reason {
				condition.Reason = newCondition.Reason
			}
			if !newCondition.LastProbeTime.IsZero() && newCondition.LastProbeTime != condition.LastProbeTime {
				condition.LastProbeTime = newCondition.LastProbeTime
			}
			if !newCondition.LastTransitionTime.IsZero() && newCondition.LastTransitionTime != condition.LastTransitionTime {
				condition.LastTransitionTime = newCondition.LastTransitionTime
			}
		}
		conditions = append(conditions, condition)
	}
	if !exist {
		if newCondition.LastProbeTime.IsZero() {
			newCondition.LastProbeTime = metav1.Now()
		}
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = metav1.Now()
		}
		conditions = append(conditions, newCondition)
	}
	cluster.Status.Conditions = conditions
}

func (c *Controller) persistUpdate(cluster *platformv1.Cluster) error {
	var err error
	for i := 0; i < clusterClientRetryCount; i++ {
		_, err = c.client.PlatformV1().Clusters().UpdateStatus(cluster)
		if err == nil {
			return nil
		}
		// if the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to cluster set that no longer exists", log.String("clusterName", cluster.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to cluster '%s' that has been changed since we received it: %v", cluster.ClusterName, err)
		}
		log.Warn("Failed to persist updated status of cluster", log.String("clusterName", cluster.ClusterName), log.Err(err))
		time.Sleep(clusterClientRetryInterval)
	}

	return err
}

func (c *Controller) doInitializing(cluster *platformv1.Cluster) error {
	if cluster.Spec.Type == platformv1.ClusterImported {
		c.ensureHealthCheck(cluster.Name, cluster)
		return nil
	}

	clusterCredential, err := util.ClusterCredentialV1(c.client.PlatformV1(), cluster.Name)
	if err != nil {
		if errors.IsNotFound(err) && cluster.Spec.Type != platformv1.ClusterImported {
			clusterCredential = &platformv1.ClusterCredential{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("cc-%s", cluster.Name),
				},
				TenantID:    cluster.Spec.TenantID,
				ClusterName: cluster.Name,
			}
			clusterCredential, err = c.client.PlatformV1().ClusterCredentials().Create(clusterCredential)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	clusterProvider, err := provider.LoadClusterProvider(c.clusterProviders, string(cluster.Spec.Type))
	if err != nil {
		return err
	}
	args := clusterprovider.Cluster{
		Cluster:           *cluster,
		ClusterCredential: *clusterCredential,
	}
	resp, err := clusterProvider.OnInitialize(args)
	if err != nil {
		cluster.Status.Message = err.Error()
		cluster.Status.Reason = reasonFailedInit
		_, _ = c.client.PlatformV1().Clusters().Update(cluster)
		return err
	}
	condition := resp.Status.Conditions[len(resp.Status.Conditions)-1]
	if condition.Status == platformv1.ConditionFalse { // means current condition run into error
		resp.Status.Message = condition.Message
		resp.Status.Reason = condition.Reason
		_, _ = c.client.PlatformV1().Clusters().Update(&resp.Cluster)
		return fmt.Errorf("OnInitialize.%s [Failed] reason: %s message: %s",
			condition.Type, condition.Reason, condition.Message)
	}
	resp.Status.Message = ""
	resp.Status.Reason = ""
	_, err = c.client.PlatformV1().Clusters().Update(&resp.Cluster)
	if err != nil {
		return err
	}
	_, err = c.client.PlatformV1().ClusterCredentials().Update(&resp.ClusterCredential)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) doUpdate(cluster *platformv1.Cluster) error {
	clusterCredential, err := util.ClusterCredentialV1(c.client.PlatformV1(), cluster.Name)
	if err != nil {
		return err
	}
	switch cluster.Spec.Type {
	case platformv1.ClusterImported:
		log.Warn("Imported cluster don's support update", log.String("clusterName", cluster.Name))
		return nil
	default:
		clusterProvider, err := provider.LoadClusterProvider(c.clusterProviders, string(cluster.Spec.Type))
		if err != nil {
			return err
		}
		args := clusterprovider.Cluster{
			Cluster:           *cluster,
			ClusterCredential: *clusterCredential,
		}
		resp, err := clusterProvider.OnUpdate(args)
		if err != nil {
			cluster.Status.Message = err.Error()
			cluster.Status.Reason = reasonFailedUpdate
			_, _ = c.client.PlatformV1().Clusters().Update(cluster)
			return err
		}
		resp.Status.Message = ""
		resp.Status.Reason = ""
		_, err = c.client.PlatformV1().Clusters().Update(&resp.Cluster)
		if err != nil {
			return err
		}
	}

	return nil
}
