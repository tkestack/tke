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

package meshmanager

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	meshv1informer "tkestack.io/tke/api/client/informers/externalversions/mesh/v1"
	meshv1lister "tkestack.io/tke/api/client/listers/mesh/v1"
	v1 "tkestack.io/tke/api/mesh/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/mesh/controller/meshmanager/images"
	"tkestack.io/tke/pkg/mesh/util"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"

	//"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	namespace      = "tke"
	controllerName = "mesh-controller"

	crbName        = "mesh-manager-role-binding"
	svcName        = "mesh-manager"
	svcAccountName = "mesh-manager"
	deploymentName = "mesh-manager"

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	timeOut       = 5 * time.Minute
	maxRetryCount = 5

	upgradePatchTemplate = `[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`

	storageTypeElasticsearch = "elasticsearch"
	storageTypeThanos        = "thanos"
)

// Controller is responsible for performing actions dependent upon a MeshManager phase.
type Controller struct {
	platformClient platformversionedclient.PlatformV1Interface
	client         clientset.Interface
	cache          *meshManagerCache
	health         sync.Map
	checking       sync.Map
	upgrading      sync.Map
	queue          workqueue.RateLimitingInterface
	lister         meshv1lister.MeshManagerLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
}

// NewController creates a new MeshManager Controller object.
func NewController(client clientset.Interface, platformClient platformversionedclient.PlatformV1Interface,
	informer meshv1informer.MeshManagerInformer,
	resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		platformClient: platformClient,
		client:         client,
		cache:          &meshManagerCache{lcMap: make(map[string]*cachedMeshManager)},
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.MeshV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.MeshV1().RESTClient().GetRateLimiter())
	}

	// configure the informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueMeshManager,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldMeshManager, ok1 := oldObj.(*v1.MeshManager)
				curMeshManager, ok2 := newObj.(*v1.MeshManager)
				if ok1 && ok2 && controller.needsUpdate(oldMeshManager, curMeshManager) {
					controller.enqueueMeshManager(newObj)
				}
			},
			DeleteFunc: controller.enqueueMeshManager,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

func (c *Controller) enqueueMeshManager(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for MeshManager object",
			log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
	log.Infof("enqueue MeshManager with key %v current queue size is %v", key, c.queue.Len())
}

func (c *Controller) needsUpdate(old *v1.MeshManager, new *v1.MeshManager) bool {
	return !reflect.DeepEqual(old, new)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting Mesh controller")
	defer log.Info("Shutting down Mesh controller")

	if !cache.WaitForCacheSync(stopCh, c.listerSynced) {
		return fmt.Errorf("failed to wait for MeshManager cache to sync")
	}

	c.stopCh = stopCh

	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}

	<-stopCh
	return nil
}

// worker processes the queue of namespace objects.
// Each key can be in the queue at most once.
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

	err := c.syncMeshManager(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing MeshManager %s (will retry): %v", key, err))
	c.queue.AddRateLimited(key)

	return true
}

// syncMeshManager will sync the MeshManager with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncMeshManager(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing MeshManager", log.String("name", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// MeshManager holds the latest MeshManager info from apiserver.
	MeshManager, err := c.lister.Get(name)
	switch {
	case k8serrors.IsNotFound(err):
		log.Info("MeshManager has been deleted. Attempting to cleanup resources", log.String("name", key))
		err = c.processMeshManagerDeletion(context.Background(), key)
	case err != nil:
		log.Warn("Unable to retrieve MeshManager from store", log.String("name", key), log.Err(err))
	default:
		cachedMeshManager := c.cache.getOrCreate(key)
		err = c.processMeshManagerUpdate(context.Background(), cachedMeshManager, MeshManager, key)
	}

	return err
}

func (c *Controller) processMeshManagerDeletion(ctx context.Context, key string) error {
	cachedMeshManager, ok := c.cache.get(key)
	if !ok {
		log.Error("MeshManager not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processMeshManagerDelete(ctx, cachedMeshManager, key)
}

func (c *Controller) processMeshManagerDelete(ctx context.Context, cachedMeshManager *cachedMeshManager, key string) error {
	log.Info("MeshManager will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the MeshManager cache", log.String("name", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the MeshManager health cache", log.String("name", key))
		c.health.Delete(key)
	}

	MeshManager := cachedMeshManager.state
	return c.uninstallMeshManager(ctx, MeshManager)
}

func (c *Controller) uninstallMeshManager(ctx context.Context, MeshManager *v1.MeshManager) error {
	log.Info("Start to uninstall MeshManager",
		log.String("name", MeshManager.Name),
		log.String("clusterName", MeshManager.Spec.ClusterName))

	kubeClient, err := util.GetClusterClient(ctx, MeshManager.Spec.ClusterName, c.platformClient)
	if err != nil {
		log.Errorf("unable to get cluster client")
		return err
	}

	// Delete the mesh-manager deployment.
	clearDeploymentErr := kubeClient.AppsV1().
		Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{})
	// Delete the ClusterRoleBinding.
	clearCRBErr := kubeClient.RbacV1().
		ClusterRoleBindings().Delete(ctx, crbName, metav1.DeleteOptions{})
	// Delete the ServiceAccount.
	clearSVCErr := kubeClient.CoreV1().ServiceAccounts(namespace).
		Delete(ctx, svcAccountName, metav1.DeleteOptions{})

	failed := false

	if clearDeploymentErr != nil && !k8serrors.IsNotFound(clearDeploymentErr) {
		failed = true
		log.Error("delete deployment for MeshManager failed",
			log.String("name", MeshManager.Name),
			log.String("clusterName", MeshManager.Spec.ClusterName),
			log.Err(clearDeploymentErr))
	}

	if clearCRBErr != nil && !k8serrors.IsNotFound(clearCRBErr) {
		failed = true
		log.Error("delete crb for MeshManager failed",
			log.String("name", MeshManager.Name),
			log.String("clusterName", MeshManager.Spec.ClusterName),
			log.Err(clearCRBErr))
	}

	if clearSVCErr != nil && !k8serrors.IsNotFound(clearSVCErr) {
		failed = true
		log.Error("delete service account for MeshManager failed",
			log.String("name", MeshManager.Name),
			log.String("clusterName", MeshManager.Spec.ClusterName),
			log.Err(clearSVCErr))
	}

	if failed {
		return errors.New("delete MeshManager failed")
	}

	return nil
}

func (c *Controller) processMeshManagerUpdate(ctx context.Context, cachedMeshManager *cachedMeshManager, MeshManager *v1.MeshManager, key string) error {
	if cachedMeshManager.state != nil {
		// exist and the cluster name changed
		if cachedMeshManager.state.UID != MeshManager.UID {
			if err := c.processMeshManagerDelete(ctx, cachedMeshManager, key); err != nil {
				return err
			}
		}
	}
	err := c.createMeshManagerIfNeeded(ctx, key, cachedMeshManager, MeshManager)
	if err != nil {
		return err
	}

	cachedMeshManager.state = MeshManager
	// Always update the cache upon success.
	c.cache.set(key, cachedMeshManager)
	return nil
}

func (c *Controller) createMeshManagerIfNeeded(
	ctx context.Context,
	key string,
	cachedMeshManager *cachedMeshManager,
	MeshManager *v1.MeshManager) error {

	switch MeshManager.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info("MeshManager will be created", log.String("name", key))
		err := c.installMeshManager(ctx, MeshManager)
		if err == nil {
			MeshManager = MeshManager.DeepCopy()
			fillOperatorStatus(MeshManager)
			MeshManager.Status.Phase = v1.AddonPhaseChecking
			MeshManager.Status.Reason = ""
			MeshManager.Status.RetryCount = 0
			return c.persistUpdate(ctx, MeshManager)
		}
		log.Errorf("MeshManager install error %v %v", err, log.String("name", key))

		// Install MeshManager failed.
		MeshManager = MeshManager.DeepCopy()
		fillOperatorStatus(MeshManager)
		MeshManager.Status.Phase = v1.AddonPhaseReinitializing
		MeshManager.Status.Reason = err.Error()
		MeshManager.Status.RetryCount = 1
		MeshManager.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(ctx, MeshManager)
	case v1.AddonPhaseReinitializing:
		log.Info("MeshManager will be reinitialized", log.String("name", key))
		var interval = time.Since(MeshManager.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go func() {
			reInitialErr := wait.Poll(waitTime, timeOut,
				c.meshManagerReinitialize(ctx, key, cachedMeshManager, MeshManager))
			if reInitialErr != nil {
				log.Error("Reinitialize MeshManager failed",
					log.String("name", MeshManager.Name),
					log.String("clusterName", MeshManager.Spec.ClusterName),
					log.Err(reInitialErr))
			}
		}()
	case v1.AddonPhaseChecking:
		log.Info("MeshManager will be checked", log.String("name", key))
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, true)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				checkStatusErr := wait.PollImmediate(5*time.Second, 5*time.Minute+10*time.Second,
					c.checkMeshManagerStatus(ctx, MeshManager, key, initDelay))
				if checkStatusErr != nil {
					log.Error("Check status of MeshManager failed",
						log.String("name", MeshManager.Name),
						log.String("clusterName", MeshManager.Spec.ClusterName),
						log.Err(checkStatusErr))
				}
			}()
		}
	case v1.AddonPhaseRunning:
		log.Info("MeshManager will be running", log.String("name", key))
		if c.needUpgrade(MeshManager) {
			c.health.Delete(key)
			MeshManager = MeshManager.DeepCopy()
			MeshManager.Status.Phase = v1.AddonPhaseUpgrading
			MeshManager.Status.Reason = ""
			MeshManager.Status.RetryCount = 0
			return c.persistUpdate(ctx, MeshManager)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go func() {
				defer c.health.Delete(key)
				healthErr := wait.PollImmediateUntil(5*time.Minute,
					c.watchMeshManagerHealth(ctx, key), c.stopCh)
				if healthErr != nil {
					log.Error("Watch health of MeshManager failed",
						log.String("name", MeshManager.Name),
						log.String("clusterName", MeshManager.Spec.ClusterName),
						log.Err(healthErr))
				}
			}()
		}
	case v1.AddonPhaseUpgrading:
		log.Info("MeshManager will be upgraded", log.String("name", key))
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, true)
			upgradeDelay := time.Now().Add(timeOut)
			go func() {
				defer c.upgrading.Delete(key)
				upgradeErr := wait.PollImmediate(5*time.Second, timeOut,
					c.upgradeMeshManager(ctx, MeshManager, key, upgradeDelay))
				if upgradeErr != nil {
					log.Error("Upgrade MeshManager failed",
						log.String("name", MeshManager.Name),
						log.String("clusterName", MeshManager.Spec.ClusterName),
						log.Err(upgradeErr))
				}
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("MeshManager failed", log.String("name", key))
		c.health.Delete(key)
		c.checking.Delete(key)
		c.upgrading.Delete(key)
	}
	return nil
}

func (c *Controller) installMeshManager(ctx context.Context, MeshManager *v1.MeshManager) error {
	kubeClient, err := util.GetClusterClient(ctx, MeshManager.Spec.ClusterName, c.platformClient)

	if err != nil {
		log.Infof("unable to get cluster client %v", err)
		return err
	}

	// Create ServiceAccount.
	if err := c.installServiceAccount(ctx, MeshManager, kubeClient); err != nil {
		return err
	}

	// Create ClusterRoleBinding.
	if err := c.installClusterRoleBinding(ctx, MeshManager, kubeClient); err != nil {
		return err
	}

	// Create ConfigMap
	//if err := c.installClusterRoleBinding(ctx, MeshManager, kubeClient); err != nil {
	//	return err
	//}

	// Create Service
	if err := c.installService(ctx, MeshManager, kubeClient); err != nil {
		return err
	}

	// Create Deployment.
	return c.installDeployment(ctx, MeshManager, kubeClient)
}

func (c *Controller) installServiceAccount(ctx context.Context, MeshManager *v1.MeshManager, kubeClient kubernetes.Interface) error {
	svc := c.genServiceAccount()
	svcClient := kubeClient.CoreV1().ServiceAccounts(namespace)

	_, err := svcClient.Get(ctx, svc.Name, metav1.GetOptions{})
	if err == nil {
		log.Info("ServiceAccount of MeshManager is already created",
			log.String("name", MeshManager.Name))
		return nil
	}

	if k8serrors.IsNotFound(err) {
		_, err = svcClient.Create(ctx, svc, metav1.CreateOptions{})
		return err
	}

	return fmt.Errorf("get svc failed: %v", err)
}

func (c *Controller) installClusterRoleBinding(ctx context.Context, MeshManager *v1.MeshManager, kubeClient kubernetes.Interface) error {
	crb := c.genClusterRoleBinding()
	crbClient := kubeClient.RbacV1().ClusterRoleBindings()

	oldCRB, err := crbClient.Get(ctx, crb.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = crbClient.Create(ctx, crb, metav1.CreateOptions{})
			return err
		}
		return fmt.Errorf("get crb failed: %v", err)
	}

	if equality.Semantic.DeepEqual(oldCRB.RoleRef, crb.RoleRef) &&
		equality.Semantic.DeepEqual(oldCRB.Subjects, crb.Subjects) {
		log.Info("ClusterRoleBinding of MeshManager created",
			log.String("name", MeshManager.Name))
		return nil
	}

	newCRB := oldCRB.DeepCopy()
	newCRB.RoleRef = crb.RoleRef
	newCRB.Subjects = crb.Subjects
	_, err = crbClient.Update(ctx, newCRB, metav1.UpdateOptions{})

	return err
}

func (c *Controller) installService(ctx context.Context, MeshManager *v1.MeshManager, kubeClient kubernetes.Interface) error {
	svc := c.genService()
	svcClient := kubeClient.CoreV1().Services(namespace)

	_, err := svcClient.Get(ctx, svc.Name, metav1.GetOptions{})
	if err == nil {
		log.Info("ServiceAccount of MeshManager is already created",
			log.String("name", MeshManager.Name))
		return nil
	}

	if k8serrors.IsNotFound(err) {
		_, err = svcClient.Create(ctx, svc, metav1.CreateOptions{})
		return err
	}

	return fmt.Errorf("get svc failed: %v", err)
}

func (c *Controller) installDeployment(ctx context.Context, MeshManager *v1.MeshManager, kubeClient kubernetes.Interface) error {
	deploy := c.genDeployment(MeshManager)
	deploymentClient := kubeClient.AppsV1().Deployments(namespace)

	oldDeploy, err := deploymentClient.Get(ctx, deploy.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = deploymentClient.Create(ctx, deploy, metav1.CreateOptions{})
			return err
		}
		return fmt.Errorf("get deployment failed: %v", err)
	}

	newDeploy := oldDeploy.DeepCopy()
	newDeploy.Labels = deploy.Labels
	newDeploy.Spec = deploy.Spec

	if len(oldDeploy.Spec.Template.Spec.Containers) == 1 {
		newDeploy.Spec.Template.Spec.Containers[0].Resources = oldDeploy.Spec.Template.Spec.Containers[0].Resources
	}

	_, err = deploymentClient.Update(ctx, newDeploy, metav1.UpdateOptions{})

	return err
}

func (c *Controller) genServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountName,
			Namespace: namespace,
		},
	}
}

func (c *Controller) genClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountName,
				Namespace: namespace,
			},
		},
	}
}

func (c *Controller) genService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: namespace,
			Labels:    map[string]string{"app": svcName},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       10072,
					TargetPort: intstr.FromInt(10072),
					Protocol:   corev1.ProtocolTCP,
					Name:       "http",
				},
			},
			Selector: map[string]string{"app": deploymentName},
		},
	}
}

func (c *Controller) genDeployment(meshManager *v1.MeshManager) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Labels:    map[string]string{"app": deploymentName},
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": deploymentName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": deploymentName},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcAccountName,
					Tolerations: []corev1.Toleration{
						{Key: "node-role.kubernetes.io/master", Effect: corev1.TaintEffectNoSchedule},
					},
					Containers: []corev1.Container{
						{
							Name:  deploymentName,
							Image: images.Get(meshManager.Spec.Version).MeshManager.FullName(),
							Resources: corev1.ResourceRequirements{
								// TODO: add support for configuring them
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(500, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewScaledQuantity(512, resource.Mega),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(300, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewScaledQuantity(250, resource.Mega),
								},
							},
							Args: []string{
								"--insecure-bind-address=0.0.0.0",
								"--insecure-port=10072",
								"--log-level=info",
								"--image-hub="+ containerregistryutil.GetPrefix(),
								"--db-host=" + meshManager.Spec.DataBase.Host,
								"--db-port=" + strconv.Itoa((int)(meshManager.Spec.DataBase.Port)),
								"--db-user=" + meshManager.Spec.DataBase.UserName,
								"--db-password=" + meshManager.Spec.DataBase.Password,
								"--db-name=" + meshManager.Spec.DataBase.DbName,
								"--monitor-tracing-storage-type=" + meshManager.Spec.TracingStorageBackend.StorageType,
								"--monitor-tracing-storage-addresses=" + strings.Join(meshManager.Spec.TracingStorageBackend.StorageAddresses, ","),
								"--monitor-tracing-storage-username=" + meshManager.Spec.TracingStorageBackend.UserName,
								"--monitor-tracing-storage-password=" + meshManager.Spec.TracingStorageBackend.Password,
								"--monitor-metrics-storage-type=" + meshManager.Spec.MetricStorageBackend.StorageType,
								"--monitor-metrics-storage-addresses=" + strings.Join(meshManager.Spec.MetricStorageBackend.StorageAddresses, ","),
								"--monitor-metrics-storage-query-address=" + meshManager.Spec.MetricStorageBackend.QueryAddress,
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "localtime", MountPath: "/etc/localtime"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{Name: "localtime", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/localtime"}}},
					},
				},
			},
		},
	}

	return deploy
}

func (c *Controller) persistUpdate(ctx context.Context, MeshManager *v1.MeshManager) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.MeshV1().MeshManagers().UpdateStatus(ctx, MeshManager, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if k8serrors.IsNotFound(err) {
			log.Info("Not persisting update to MeshManager that no longer exists",
				log.String("clusterName", MeshManager.Spec.ClusterName), log.Err(err))
			return nil
		}
		if k8serrors.IsConflict(err) {
			return fmt.Errorf("not persisting update to MeshManager %q that has been changed since we received it: %v", MeshManager.Spec.ClusterName, err)
		}
		log.Warn("Failed to persist updated status of MeshManager",
			log.String("name", MeshManager.Name),
			log.String("clusterName", MeshManager.Spec.ClusterName),
			log.String("phase", string(MeshManager.Status.Phase)), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}

func fillOperatorStatus(MeshManager *v1.MeshManager) { //what's the purpose of this?
	log.Infof("spec.version is %v, status.version is %v", MeshManager.Spec.Version, MeshManager.Status.Version)
	MeshManager.Status.Version = MeshManager.Spec.Version
}

func (c *Controller) meshManagerReinitialize(
	ctx context.Context,
	key string,
	cachedMeshManager *cachedMeshManager,
	MeshManager *v1.MeshManager) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installMeshManager(ctx, MeshManager)
		if err == nil {
			MeshManager = MeshManager.DeepCopy()
			MeshManager.Status.Phase = v1.AddonPhaseChecking
			MeshManager.Status.Reason = ""
			MeshManager.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(ctx, MeshManager)
			if err != nil {
				return true, err
			}
			return true, nil
		}

		// First, rollback the MeshManager.
		log.Info("Rollback MeshManager",
			log.String("name", MeshManager.Name),
			log.String("clusterName", MeshManager.Spec.ClusterName))
		if err := c.uninstallMeshManager(ctx, MeshManager); err != nil {
			log.Error("Uninstall MeshManager failed", log.Err(err))
			return true, err
		}

		if MeshManager.Status.RetryCount == maxRetryCount {
			MeshManager = MeshManager.DeepCopy()
			MeshManager.Status.Phase = v1.AddonPhaseFailed
			MeshManager.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(ctx, MeshManager)
			if err != nil {
				log.Error("Update MeshManager failed", log.Err(err))
				return true, err
			}
			return true, nil
		}

		// Add the retry count will trigger reinitialize function from the persistent controller again.
		MeshManager = MeshManager.DeepCopy()
		MeshManager.Status.Phase = v1.AddonPhaseReinitializing
		MeshManager.Status.Reason = err.Error()
		MeshManager.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		MeshManager.Status.RetryCount++
		return true, c.persistUpdate(ctx, MeshManager)
	}
}

func (c *Controller) checkMeshManagerStatus(
	ctx context.Context,
	MeshManager *v1.MeshManager,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check MeshManager health", log.String("name", MeshManager.Name))

		kubeClient, err := util.GetClusterClient(ctx, MeshManager.Spec.ClusterName, c.platformClient)

		if err != nil { //what if cluster does not exists?
			return false, err
		}

		if _, ok := c.checking.Load(key); !ok {
			log.Debug("Checking over MeshManager addon status")
			return true, nil
		}

		MeshManager, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		deploy, err := kubeClient.AppsV1().Deployments(namespace).
			Get(ctx, deploymentName, metav1.GetOptions{})

		if err != nil || deploy.Status.Replicas == 0 ||
			deploy.Status.AvailableReplicas < 1 {
			if time.Now().After(initDelay) {
				if err != nil {
					log.Errorf("check mesh-manager %v status failed %v", MeshManager.Name, err)
				} else {
					log.Warnf("meshManager %v not healthy", MeshManager.Name)
				}
				MeshManager = MeshManager.DeepCopy()
				MeshManager.Status.Phase = v1.AddonPhaseChecking
				MeshManager.Status.Reason = "mesh-manager is not healthy in status check"
				if err = c.persistUpdate(ctx, MeshManager); err != nil {
					return false, err
				}
			}
			return false, nil
		}

		MeshManager = MeshManager.DeepCopy()
		MeshManager.Status.Phase = v1.AddonPhaseRunning
		MeshManager.Status.Reason = ""
		if err = c.persistUpdate(ctx, MeshManager); err != nil {
			return false, err
		}

		return true, nil
	}
}

func (c *Controller) needUpgrade(MeshManager *v1.MeshManager) bool {
	return MeshManager.Spec.Version != MeshManager.Status.Version
}

func (c *Controller) watchMeshManagerHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		MeshManager, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		log.Info("Start check health of MeshManager", log.String("name", MeshManager.Name))

		kubeClient, err := util.GetClusterClient(ctx, MeshManager.Spec.ClusterName, c.platformClient)
		if err != nil {
			return false, err
		}

		if _, ok := c.health.Load(key); !ok {
			log.Info("Health check over.")
			return true, nil
		}

		_, err = kubeClient.AppsV1().Deployments(namespace).
			Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("watch MeshManager %v status failed %v", MeshManager.Name, err)
			MeshManager = MeshManager.DeepCopy()
			MeshManager.Status.Phase = v1.AddonPhaseChecking
			MeshManager.Status.Reason = "MeshManager is not healthy in watch."
			if err = c.persistUpdate(ctx, MeshManager); err != nil {
				return false, err
			}
			return true, nil
		}

		return false, nil
	}
}

func (c *Controller) upgradeMeshManager(
	ctx context.Context,
	MeshManager *v1.MeshManager,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade MeshManager", log.String("name", MeshManager.Name))
		kubeClient, err := util.GetClusterClient(ctx, MeshManager.Spec.ClusterName, c.platformClient)
		if err != nil {
			return false, err
		}
		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading MeshManager", log.String("name", MeshManager.Name))
			return true, nil
		}

		MeshManager, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		patch := fmt.Sprintf(upgradePatchTemplate, images.Get(MeshManager.Spec.Version).MeshManager.FullName())

		_, err = kubeClient.AppsV1().Deployments(namespace).
			Patch(ctx, deploymentName, types.JSONPatchType, []byte(patch), metav1.PatchOptions{})
		if err != nil {
			if time.Now().After(initDelay) {
				MeshManager = MeshManager.DeepCopy()
				MeshManager.Status.Phase = v1.AddonPhaseFailed
				MeshManager.Status.Reason = "Failed to upgrade MeshManager."
				if err = c.persistUpdate(ctx, MeshManager); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		MeshManager = MeshManager.DeepCopy()
		fillOperatorStatus(MeshManager)
		MeshManager.Status.Phase = v1.AddonPhaseChecking
		MeshManager.Status.Reason = ""
		if err = c.persistUpdate(ctx, MeshManager); err != nil {
			return false, err
		}
		return true, nil
	}
}
