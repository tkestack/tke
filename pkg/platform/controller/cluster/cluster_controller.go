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
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/rand"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/controller/cluster/deletion"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
	"tkestack.io/tke/pkg/util/strategicpatch"
)

type ContextKey int

const (
	KeyLister                ContextKey = iota
	conditionTypeHealthCheck            = "HealthCheck"
	failedHealthCheckReason             = "FailedHealthCheck"

	resyncInternal = 5 * time.Minute
)

// Controller is responsible for performing actions dependent upon a cluster phase.
type Controller struct {
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.ClusterLister
	listerSynced cache.InformerSynced

	log            log.Logger
	platformClient platformversionedclient.PlatformV1Interface
	deleter        deletion.ClusterDeleterInterface
}

// NewController creates a new Controller object.
func NewController(
	platformClient platformversionedclient.PlatformV1Interface,
	clusterInformer platformv1informer.ClusterInformer,
	resyncPeriod time.Duration,
	finalizerToken platformv1.FinalizerName) *Controller {

	rand.Seed(time.Now().Unix())

	c := &Controller{
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cluster"),

		log:            log.WithName("ClusterController"),
		platformClient: platformClient,
		deleter: deletion.NewClusterDeleter(platformClient.Clusters(),
			platformClient,
			finalizerToken,
			true),
	}

	if platformClient != nil && platformClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("cluster_controller", platformClient.RESTClient().GetRateLimiter())
	}

	clusterInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addCluster,
			UpdateFunc: c.updateCluster,
		},
		resyncPeriod,
	)
	c.lister = clusterInformer.Lister()
	c.listerSynced = clusterInformer.Informer().HasSynced

	return c
}

func (c *Controller) addCluster(obj interface{}) {
	cluster := obj.(*platformv1.Cluster)
	c.log.Info("Adding cluster", "clusterName", cluster.Name)
	c.enqueue(cluster)
}

func (c *Controller) updateCluster(old, obj interface{}) {
	oldCluster := old.(*platformv1.Cluster)
	cluster := obj.(*platformv1.Cluster)
	if !c.needsUpdate(oldCluster, cluster) {
		return
	}
	c.log.Info("Updating cluster", "clusterName", cluster.Name)
	c.enqueue(cluster)
}

func (c *Controller) enqueue(obj *platformv1.Cluster) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %+v: %v", obj, err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *platformv1.Cluster, new *platformv1.Cluster) bool {
	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true
	}

	if old.Status.Phase == platformv1.ClusterRunning && new.Status.Phase == platformv1.ClusterTerminating {
		return true
	}

	if !reflect.DeepEqual(old.ObjectMeta.Annotations, new.ObjectMeta.Annotations) {
		return true
	}

	if !reflect.DeepEqual(old.ObjectMeta.Labels, new.ObjectMeta.Labels) {
		return true
	}

	// Control the synchronization interval through the health detection interval
	// to avoid version conflicts caused by concurrent modification
	healthCondition := new.GetCondition(conditionTypeHealthCheck)
	if healthCondition == nil {
		return true
	}
	if time.Since(healthCondition.LastProbeTime.Time) > resyncInternal {
		return true
	}

	return false
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting cluster controller")
	defer log.Info("Shutting down cluster controller")

	if err := clusterprovider.Setup(); err != nil {
		return err
	}

	if ok := cache.WaitForCacheSync(stopCh, c.listerSynced); !ok {
		return fmt.Errorf("failed to wait for cluster caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh

	if err := clusterprovider.Teardown(); err != nil {
		return err
	}

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

	utilruntime.HandleError(fmt.Errorf("error processing cluster %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncCluster will sync the Cluster with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncCluster(key string) error {
	ctx := c.log.WithValues("cluster", key).WithContext(context.TODO())

	startTime := time.Now()
	defer func() {
		log.FromContext(ctx).Info("Finished syncing cluster", "processTime", time.Since(startTime).String())
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	cluster, err := c.lister.Get(name)
	if apierrors.IsNotFound(err) {
		log.FromContext(ctx).Info("cluster has been deleted")
	}
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("unable to retrieve cluster %v from store: %v", key, err))
		return err
	}

	valueCtx := context.WithValue(ctx, KeyLister, &c.lister)
	return c.reconcile(valueCtx, key, cluster)
}

func (c *Controller) reconcile(ctx context.Context, key string, cluster *platformv1.Cluster) error {
	var err error

	c.ensureSyncCredentialClusterName(ctx, cluster)
	c.ensureSyncClusterMachineNodeLabel(ctx, cluster)

	switch cluster.Status.Phase {
	case platformv1.ClusterInitializing:
		err = c.onCreate(ctx, cluster)
	case platformv1.ClusterRunning, platformv1.ClusterFailed:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterUpgrading:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterUpscaling, platformv1.ClusterDownscaling:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterTerminating:
		log.FromContext(ctx).Info("Cluster has been terminated. Attempting to cleanup resources")
		err = c.deleter.Delete(ctx, key)
		if err == nil {
			log.FromContext(ctx).Info("Machine has been successfully deleted")
		}
	default:
		log.FromContext(ctx).Info("unknown cluster phase", "status.phase", cluster.Status.Phase)
	}

	return err
}

func (c *Controller) onCreate(ctx context.Context, cluster *platformv1.Cluster) error {
	var err error

	cluster, err = c.ensureCreateClusterCredential(ctx, cluster)
	if err != nil {
		return fmt.Errorf("ensureCreateClusterCredential error: %w", err)
	}

	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return err
	}
	clusterWrapper, err := typesv1.GetCluster(ctx, c.platformClient, cluster)
	if err != nil {
		return err
	}

	for clusterWrapper.Status.Phase == platformv1.ClusterInitializing {
		err = provider.OnCreate(ctx, clusterWrapper)
		if err != nil {
			// Update status, ignore failure
			_, _ = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			_, _ = c.platformClient.Clusters().Update(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			return err
		}
		clusterWrapper.ClusterCredential, err = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		clusterWrapper.Cluster, err = c.platformClient.Clusters().Update(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) onUpdate(ctx context.Context, cluster *platformv1.Cluster) error {
	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return err
	}
	clusterWrapper, err := typesv1.GetCluster(ctx, c.platformClient, cluster)
	if err != nil {
		return err
	}
	if clusterWrapper.Status.Phase == platformv1.ClusterRunning || clusterWrapper.Status.Phase == platformv1.ClusterFailed {
		err = provider.OnUpdate(ctx, clusterWrapper)
		clusterWrapper = c.checkHealth(ctx, clusterWrapper)
		if err != nil {
			// Update status, ignore failure
			if clusterWrapper.IsCredentialChanged {
				_, _ = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			}

			_, _ = c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			return err
		}
		if clusterWrapper.IsCredentialChanged {
			clusterWrapper.ClusterCredential, err = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
		clusterWrapper.Cluster, err = c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	} else {
		for clusterWrapper.Status.Phase != platformv1.ClusterRunning {
			err = provider.OnUpdate(ctx, clusterWrapper)
			if err != nil {
				// Update status, ignore failure
				if clusterWrapper.IsCredentialChanged {
					_, _ = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
				}

				_, _ = c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
				return err
			}
			if clusterWrapper.IsCredentialChanged {
				clusterWrapper.ClusterCredential, err = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
			}
			clusterWrapper.Cluster, err = c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ensureCreateClusterCredential creates ClusterCredential for cluster if ClusterCredentialRef is nil.
// TODO: add gc collector for clean non reference ClusterCredential.
func (c *Controller) ensureCreateClusterCredential(ctx context.Context, cluster *platformv1.Cluster) (*platformv1.Cluster, error) {
	if cluster.Spec.ClusterCredentialRef != nil {
		return cluster, nil
	}

	var err error
	credential := &platformv1.ClusterCredential{
		TenantID:    cluster.Spec.TenantID,
		ClusterName: cluster.Name,
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cluster, platformv1.SchemeGroupVersion.WithKind("Cluster"))},
		},
	}

	credential, err = c.platformClient.ClusterCredentials().Create(ctx, credential, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	cluster.Spec.ClusterCredentialRef = &corev1.LocalObjectReference{Name: credential.Name}
	cluster, err = c.platformClient.Clusters().Update(ctx, cluster, metav1.UpdateOptions{})
	if err != nil {
		// Possible deletion failure will result in dirty data. So need gc collector.
		_ = c.platformClient.ClusterCredentials().Delete(ctx, credential.Name, metav1.DeleteOptions{})
		return nil, err
	}

	return cluster, nil
}

func (c *Controller) ensureSyncCredentialClusterName(ctx context.Context, cluster *platformv1.Cluster) {
	if cluster.Spec.ClusterCredentialRef == nil {
		clusterCredentials, err := c.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("clusterName=%s", cluster.Name)})
		if err != nil {
			return
		}
		if len(clusterCredentials.Items) > 0 {
			credential := &clusterCredentials.Items[0]
			oldCredential := credential.DeepCopy()
			if credential.ObjectMeta.OwnerReferences == nil {
				patchBytes, err := strategicpatch.GetPatchBytes(oldCredential, credential)
				credential.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*metav1.NewControllerRef(cluster, platformv1.SchemeGroupVersion.WithKind("Cluster"))}
				_, err = c.platformClient.ClusterCredentials().Patch(ctx, credential.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
				if err != nil {
					return
				}
			}
		}
		return
	}

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		credential, err := c.platformClient.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		oldCredential := credential.DeepCopy()
		if credential.ClusterName != cluster.Name {
			credential.ClusterName = cluster.Name

			patchBytes, err := strategicpatch.GetPatchBytes(oldCredential, credential)
			if err != nil {
				return fmt.Errorf("GetPatchBytes for credential error: %w", err)
			}
			_, err = c.platformClient.ClusterCredentials().Patch(ctx, credential.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.FromContext(ctx).Error(err, "sync ClusterCredential.ClusterName error")
	}
}

func (c *Controller) checkHealth(ctx context.Context, cluster *typesv1.Cluster) *typesv1.Cluster {
	if !(cluster.Status.Phase == platformv1.ClusterRunning ||
		cluster.Status.Phase == platformv1.ClusterFailed) {
		return cluster
	}

	pseudo := time.Now().Add(time.Minute * time.Duration(rand.Intn(5)))

	log.Infof("next heart beat time. now:%s pesudo:%s cls:%s", time.Now(), pseudo, cluster.Name)

	healthCheckCondition := platformv1.ClusterCondition{
		Type:          conditionTypeHealthCheck,
		Status:        platformv1.ConditionFalse,
		LastProbeTime: metav1.NewTime(pseudo),
	}

	client, err := cluster.Clientset()
	if err != nil {
		cluster.Status.Phase = platformv1.ClusterFailed

		healthCheckCondition.Reason = failedHealthCheckReason
		healthCheckCondition.Message = err.Error()
	} else {
		version, err := client.Discovery().ServerVersion()
		if err != nil {
			cluster.Status.Phase = platformv1.ClusterFailed

			healthCheckCondition.Reason = failedHealthCheckReason
			healthCheckCondition.Message = err.Error()
		} else {
			cluster.Status.Phase = platformv1.ClusterRunning
			cluster.Status.Version = strings.TrimPrefix(version.String(), "v")

			healthCheckCondition.Status = platformv1.ConditionTrue
		}
	}

	cluster.SetCondition(healthCheckCondition, false)

	log.FromContext(ctx).Info("Update cluster health status",
		"version", cluster.Status.Version,
		"phase", cluster.Status.Phase)

	return cluster
}

func (c *Controller) ensureSyncClusterMachineNodeLabel(ctx context.Context, cluster *platformv1.Cluster) {

	clusterWrapper, err := typesv1.GetCluster(ctx, c.platformClient, cluster)
	if err != nil {
		log.FromContext(ctx).Error(err, "Get cluster error")
		return
	}

	client, err := clusterWrapper.Clientset()
	if err != nil {
		log.FromContext(ctx).Error(err, "get client set error")
		return
	}

	for _, machine := range cluster.Spec.Machines {
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			node, err := client.CoreV1().Nodes().Get(ctx, machine.IP, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return nil
				}
				return err
			}

			labels := node.GetLabels()
			_, ok := labels[string(apiclient.LabelMachineIPV4)]
			if ok {
				return nil
			}

			oldNode := node.DeepCopy()
			labels[string(apiclient.LabelMachineIPV4)] = machine.IP
			node.SetLabels(labels)

			patchBytes, err := strategicpatch.GetPatchBytes(oldNode, node)
			if err != nil {
				return fmt.Errorf("GetPatchBytes for node error: %w", err)
			}

			_, err = client.CoreV1().Nodes().Patch(ctx, node.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
			return err
		})

		if err != nil {
			log.FromContext(ctx).Error(err, "sync ClusterMachine node label error")
		}
	}
}
