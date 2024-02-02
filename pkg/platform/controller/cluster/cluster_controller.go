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

	"golang.org/x/time/rate"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/rand"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	clusterconfig "tkestack.io/tke/pkg/platform/controller/cluster/config"
	"tkestack.io/tke/pkg/platform/controller/cluster/deletion"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	vendor "tkestack.io/tke/pkg/platform/util/kubevendor"
	workqueue_extension "tkestack.io/tke/pkg/platform/util/workqueue"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

type ContextKey int

const (
	KeyLister                ContextKey = iota
	conditionTypeHealthCheck            = "HealthCheck"
	failedHealthCheckReason             = "FailedHealthCheck"
)

// Controller is responsible for performing actions dependent upon a cluster phase.
type Controller struct {
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.ClusterLister
	listerSynced cache.InformerSynced

	log                                        log.Logger
	platformClient                             platformversionedclient.PlatformV1Interface
	deleter                                    deletion.ClusterDeleterInterface
	healthCheckPeriod                          time.Duration
	randomeRangeLowerLimitForHealthCheckPeriod time.Duration
	randomeRangeUpperLimitForHealthCheckPeriod time.Duration
	isCRDMode                                  bool
}

// NewController creates a new Controller object.
func NewController(
	platformClient platformversionedclient.PlatformV1Interface,
	clusterInformer platformv1informer.ClusterInformer,
	configuration clusterconfig.ClusterControllerConfiguration,
	finalizerToken platformv1.FinalizerName) *Controller {
	rand.Seed(time.Now().Unix())

	c := &Controller{
		log:            log.WithName("ClusterController"),
		platformClient: platformClient,
		deleter: deletion.NewClusterDeleter(platformClient.Clusters(),
			platformClient,
			finalizerToken,
			true),
		isCRDMode: configuration.IsCRDMode,
	}
	rateLimit := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(configuration.BucketRateLimiterLimit), configuration.BucketRateLimiterBurst)},
	)
	c.queue = workqueue_extension.NewNamedRateLimitingWithCustomQueue(rateLimit,
		workqueue_extension.NewNamed("platform", 12, c.getPriority),
		"cluster")

	if platformClient != nil && platformClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("cluster_controller", platformClient.RESTClient().GetRateLimiter())
	}

	clusterInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    c.addCluster,
				UpdateFunc: c.updateCluster,
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

	c.lister = clusterInformer.Lister()
	c.listerSynced = clusterInformer.Informer().HasSynced
	c.healthCheckPeriod = configuration.HealthCheckPeriod
	c.randomeRangeLowerLimitForHealthCheckPeriod = configuration.RandomeRangeLowerLimitForHealthCheckPeriod
	c.randomeRangeUpperLimitForHealthCheckPeriod = configuration.RandomeRangeUpperLimitForHealthCheckPeriod

	return c
}

// The higher the priority value, the higher the priority, such as priorityInitializing(10) > priorityTerminating(8)
const (
	priorityIdling       int = 2
	priorityFailed       int = 4
	priorityRunning      int = 6
	priorityTerminating  int = 8
	priorityInitializing int = 10
)

func (c *Controller) getPriority(item interface{}) int {
	var key string
	var ok bool
	if key, ok = item.(string); !ok {
		return priorityRunning
	}

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return priorityRunning
	}

	cluster, err := c.lister.Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Infof("getPriority item is not found, item: %v", item)
		}
		return priorityRunning
	}
	if cluster == nil {
		return priorityRunning
	}

	switch {
	case cluster.Status.Phase == platformv1.ClusterPhase("Idling"):
		return priorityIdling
	case cluster.Status.Phase == platformv1.ClusterFailed:
		return priorityFailed
	case cluster.Status.Phase == platformv1.ClusterRunning:
		return priorityRunning
	case cluster.Status.Phase == platformv1.ClusterTerminating:
		return priorityTerminating
	case cluster.Status.Phase == platformv1.ClusterInitializing:
		return priorityInitializing
	}

	return priorityRunning
}

func (c *Controller) addCluster(obj interface{}) {
	cluster := obj.(*platformv1.Cluster)
	c.log.Info("Adding cluster", "clusterName", cluster.Name)
	c.enqueue(cluster)
}

func (c *Controller) updateCluster(old, obj interface{}) {
	oldCluster := old.(*platformv1.Cluster)
	cluster := obj.(*platformv1.Cluster)

	controllerNeedUpddateResult := c.needsUpdate(oldCluster, cluster)
	var providerNeedUpddateResult bool
	provider, _ := clusterprovider.GetProvider(cluster.Spec.Type)
	if provider != nil {
		providerNeedUpddateResult = provider.NeedUpdate(oldCluster, cluster)
	}
	if !(controllerNeedUpddateResult || providerNeedUpddateResult) {
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
	healthCondition := new.GetCondition(conditionTypeHealthCheck)
	if !reflect.DeepEqual(old.Spec, new.Spec) {
		return true

	}
	if !reflect.DeepEqual(old.ObjectMeta.Labels, new.ObjectMeta.Labels) {
		return true
	}
	if !reflect.DeepEqual(old.ObjectMeta.Annotations, new.ObjectMeta.Annotations) {
		return true
	}
	if old.Status.Phase != new.Status.Phase {
		return true
	}
	if new.Status.Phase == platformv1.ClusterInitializing {
		// if ResourceVersion is equal, it's an resync envent, should return true.
		if old.ResourceVersion == new.ResourceVersion {
			return true
		}
		if len(new.Status.Conditions) == 0 {
			return true
		}
		if new.Status.Conditions[len(new.Status.Conditions)-1].Status == platformv1.ConditionUnknown {
			return true
		}
		// if user set last condition false block procesee until resync envent
		if new.Status.Conditions[len(new.Status.Conditions)-1].Status == platformv1.ConditionFalse {
			return false
		}
	}
	// if last health check is not long enough， return false
	if healthCondition != nil &&
		time.Since(healthCondition.LastProbeTime.Time) < c.healthCheckPeriod {
		return false
	}
	return true
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
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.FromContext(ctx).Info("cluster has been deleted")
			return nil
		}

		utilruntime.HandleError(fmt.Errorf("unable to retrieve cluster %v from store: %v", key, err))
		return err
	}

	valueCtx := context.WithValue(ctx, KeyLister, &c.lister)
	return c.reconcile(valueCtx, key, cluster)
}

func (c *Controller) reconcile(ctx context.Context, key string, cluster *platformv1.Cluster) error {
	var err error

	switch cluster.Status.Phase {
	// empty string is for crd without mutating webhook
	case "":
		cluster.Status.Phase = platformv1.ClusterInitializing
		err = c.onCreate(ctx, cluster)
	case platformv1.ClusterInitializing, platformv1.ClusterWaiting:
		err = c.onCreate(ctx, cluster)
	case platformv1.ClusterRunning, platformv1.ClusterFailed:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterUpgrading:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterUpscaling, platformv1.ClusterDownscaling:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterIdling, platformv1.ClusterConfined, platformv1.ClusterRecovering:
		err = c.onUpdate(ctx, cluster)
	case platformv1.ClusterTerminating:
		log.FromContext(ctx).Info("Cluster has been terminated. Attempting to cleanup resources")
		err = c.deleter.Delete(ctx, key)
		if err == nil {
			log.FromContext(ctx).Info("Cluster has been successfully deleted")
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
	clusterWrapper, err := clusterprovider.GetV1Cluster(ctx, c.platformClient, cluster, clusterprovider.AdminUsername)
	if err != nil {
		return err
	}

	for clusterWrapper.Status.Phase == platformv1.ClusterInitializing {
		err = provider.OnCreate(ctx, clusterWrapper)
		if err != nil {
			// Update status, ignore failure
			_, _ = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			updatedCls, _ := c.platformClient.Clusters().Update(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			// if using crd, cluster status cannot be updated through update cluster
			if c.isCRDMode {
				var clsStatus *platformv1.Cluster
				if updatedCls == nil {
					clsStatus = clusterWrapper.Cluster
				} else {
					clsStatus = updatedCls
					clsStatus.Status = clusterWrapper.Cluster.Status
				}
				_, _ = c.platformClient.Clusters().UpdateStatus(ctx, clsStatus, metav1.UpdateOptions{})
			}
			return err
		}
		clusterWrapper.ClusterCredential, err = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		clusterWrapper.RegisterRestConfig(clusterWrapper.ClusterCredential.RESTConfig(cluster))
		cls, err := c.platformClient.Clusters().Update(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		// if using crd, cluster status cannot be updated through update cluster
		if c.isCRDMode {
			cls.Status = clusterWrapper.Cluster.Status
			clusterWrapper.Cluster, err = c.platformClient.Clusters().UpdateStatus(ctx, cls, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		} else {
			clusterWrapper.Cluster = cls
		}
	}

	return nil
}

func (c *Controller) onUpdate(ctx context.Context, cluster *platformv1.Cluster) error {
	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return err
	}
	clusterWrapper, err := clusterprovider.GetV1Cluster(ctx, c.platformClient, cluster, clusterprovider.AdminUsername)
	if err != nil {
		return err
	}
	if clusterWrapper.Status.Phase == platformv1.ClusterRunning ||
		clusterWrapper.Status.Phase == platformv1.ClusterFailed ||
		clusterWrapper.Status.Phase == platformv1.ClusterIdling ||
		clusterWrapper.Status.Phase == platformv1.ClusterConfined ||
		clusterWrapper.Status.Phase == platformv1.ClusterRecovering {
		err = provider.OnUpdate(ctx, clusterWrapper)
		clusterWrapper = c.checkHealth(ctx, clusterWrapper)
		if err != nil {
			// Update status, ignore failure
			if clusterWrapper.IsCredentialChanged {
				_, _ = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			}

			updatedCluster, _ := c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			if c.isCRDMode {
				if !reflect.DeepEqual(updatedCluster.ObjectMeta.Annotations, clusterWrapper.Cluster.ObjectMeta.Annotations) {
					updatedCluster.Annotations = clusterWrapper.Cluster.Annotations
					_, _ = c.platformClient.Clusters().Update(ctx, updatedCluster, metav1.UpdateOptions{})
				}
			}
			return err
		}
		if clusterWrapper.IsCredentialChanged {
			clusterWrapper.ClusterCredential, err = c.platformClient.ClusterCredentials().Update(ctx, clusterWrapper.ClusterCredential, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
			clusterWrapper.RegisterRestConfig(clusterWrapper.ClusterCredential.RESTConfig(cluster))
		}
		cls, err := c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
		if err != nil {
			clusterWrapper.Cluster = cls
			return err
		}
		if c.isCRDMode {
			if !reflect.DeepEqual(cls.ObjectMeta.Annotations, clusterWrapper.Cluster.ObjectMeta.Annotations) {
				cls.Annotations = clusterWrapper.Cluster.Annotations
				clusterWrapper.Cluster, err = c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
				if err != nil {
					return err
				}
			}
		} else {
			clusterWrapper.Cluster = cls
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
				clusterWrapper.RegisterRestConfig(clusterWrapper.ClusterCredential.RESTConfig(cluster))
			}
			cls, err := c.platformClient.Clusters().UpdateStatus(ctx, clusterWrapper.Cluster, metav1.UpdateOptions{})
			if err != nil {
				clusterWrapper.Cluster = cls
				return err
			}
			if c.isCRDMode {
				if !reflect.DeepEqual(cls.ObjectMeta.Annotations, clusterWrapper.Cluster.ObjectMeta.Annotations) {
					cls.Annotations = clusterWrapper.Cluster.Annotations
					clusterWrapper.Cluster, err = c.platformClient.Clusters().Update(ctx, cls, metav1.UpdateOptions{})
					if err != nil {
						return err
					}
				}
			} else {
				clusterWrapper.Cluster = cls
			}
		}
	}
	return nil
}

// ensureCreateClusterCredential creates ClusterCredential for cluster if ClusterCredentialRef is nil.
func (c *Controller) ensureCreateClusterCredential(ctx context.Context, cluster *platformv1.Cluster) (*platformv1.Cluster, error) {
	if cluster.Spec.ClusterCredentialRef != nil {
		// Set OwnerReferences for imported cluster credentials
		cc, err := c.platformClient.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		cc.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(cluster, platformv1.SchemeGroupVersion.WithKind("Cluster"))}
		_, err = c.platformClient.ClusterCredentials().Update(ctx, cc, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}
		return cluster, nil
	}

	// TODO use informer search by labels.
	var clustercredentials *platformv1.ClusterCredentialList
	var err error
	if c.isCRDMode {
		labelSelector := fields.OneTermEqualSelector(platformv1.ClusterNameLable, cluster.Name).String()
		clustercredentials, err = c.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, err
		}
	} else {
		fieldSelector := fields.OneTermEqualSelector("clusterName", cluster.Name).String()
		clustercredentials, err = c.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return nil, err
		}
	}

	// [Idempotent] if not found cluster credentials, create one for next logic
	var credential *platformv1.ClusterCredential
	if len(clustercredentials.Items) == 0 {
		credential = &platformv1.ClusterCredential{
			TenantID:    cluster.Spec.TenantID,
			ClusterName: cluster.Name,
			ObjectMeta: metav1.ObjectMeta{
				Labels:       map[string]string{platformv1.ClusterNameLable: cluster.Name},
				GenerateName: "cc-",
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(cluster, platformv1.SchemeGroupVersion.WithKind("Cluster"))},
			},
		}

		credential, err = c.platformClient.ClusterCredentials().Create(ctx, credential, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		if len(clustercredentials.Items) > 1 {
			log.Warnf("cluster %s has more than one credentials, need attention!")
		}

		credential = &clustercredentials.Items[0]
	}

	cluster.Spec.ClusterCredentialRef = &corev1.LocalObjectReference{Name: credential.Name}

	return cluster, nil
}

func (c *Controller) checkHealth(ctx context.Context, cluster *typesv1.Cluster) *typesv1.Cluster {
	if !(cluster.Status.Phase == platformv1.ClusterRunning ||
		cluster.Status.Phase == platformv1.ClusterFailed) {
		return cluster
	}

	pseudo := time.Now().Add(time.Second * time.Duration(rand.Int63nRange(
		int64(c.randomeRangeLowerLimitForHealthCheckPeriod.Seconds()),
		int64(c.randomeRangeUpperLimitForHealthCheckPeriod.Seconds()))))

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
		version, err := controllerutil.CheckClusterHealthStatus(client)
		if err != nil {
			cluster.Status.Phase = platformv1.ClusterFailed

			healthCheckCondition.Reason = failedHealthCheckReason
			healthCheckCondition.Message = err.Error()
		} else {
			cluster.Status.Phase = platformv1.ClusterRunning
			cluster.Status.Version = strings.TrimPrefix(version.String(), "v")
			cluster.Status.KubeVendor = vendor.GetKubeVendor(cluster.Status.Version)

			healthCheckCondition.Status = platformv1.ConditionTrue
		}
	}

	cluster.SetCondition(healthCheckCondition, false)

	log.FromContext(ctx).Info("Update cluster health status",
		"version", cluster.Status.Version,
		"kubevendor", cluster.Status.KubeVendor,
		"phase", cluster.Status.Phase)

	return cluster
}
