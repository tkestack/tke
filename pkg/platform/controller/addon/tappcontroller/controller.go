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

package tappcontroller

import (
	normalerrors "errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/tappcontroller/images"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	maxRetryCount = 5
	timeOut       = 5 * time.Minute
)

const (
	controllerName               = "tapp-controller"
	deployTappControllerName     = "tapp-controller"
	svcTappControllerName        = "tapp-controller"
	svcAccountTappControllerName = "tapp-controller"
	crbTappControllerName        = "tapp-controller"
)

// Controller is responsible for performing actions dependent upon a tapp controller phase.
type Controller struct {
	client       clientset.Interface
	cache        *tappControllerCache
	health       sync.Map
	checking     sync.Map
	upgrading    sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.TappControllerLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, informer platformv1informer.TappControllerInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{

		client: client,
		cache:  &tappControllerCache{tappControllerMap: make(map[string]*cachedTappController)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the tapp controller informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueTappController,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldTappController, ok1 := oldObj.(*v1.TappController)
				curTappController, ok2 := newObj.(*v1.TappController)
				if ok1 && ok2 && controller.needsUpdate(oldTappController, curTappController) {
					controller.enqueueTappController(newObj)
				}
			},
			DeleteFunc: controller.enqueueTappController,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.TappController, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueTappController(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *v1.TappController, new *v1.TappController) bool {
	return !reflect.DeepEqual(old, new)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting tappcontroller controller")
	defer log.Info("Shutting down tappcontroller controller")

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

// worker processes the queue of namespace objects.
// Each namespace can be in the queue at most once.
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

	err := c.syncTappController(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing tapp controller %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncTappController will sync the tapp controller with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncTappController(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing tappController", log.String("tappControllerName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// tappController holds the latest tappController info from apiserver
	tappController, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Tapp controller has been deleted. Attempting to cleanup resources", log.String("tappControllerName", key))
		err = c.processTappControllerDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve tapp controller from store", log.String("tappControllerName", key), log.Err(err))
	default:
		cachedTappController := c.cache.getOrCreate(key)
		err = c.processTappControllerUpdate(cachedTappController, tappController, key)
	}
	return err
}

func (c *Controller) processTappControllerDeletion(key string) error {
	cachedTappController, ok := c.cache.get(key)
	if !ok {
		log.Error("Tapp controller not in cache even though the watcher thought it was. Ignoring the deletion", log.String("tappControllerName", key))
		return nil
	}
	return c.processTappControllerDelete(cachedTappController, key)
}

func (c *Controller) processTappControllerDelete(cachedTappController *cachedTappController, key string) error {
	log.Info("Tapp controller will be dropped", log.String("tappControllerName", key))

	if c.cache.Exist(key) {
		log.Info("Delete the tapp controller cache", log.String("tappControllerName", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the tapp controller health cache", log.String("tappControllerName", key))
		c.health.Delete(key)
	}

	tappController := cachedTappController.state
	return c.uninstallTappController(tappController)
}

func (c *Controller) processTappControllerUpdate(cachedTappController *cachedTappController, tappController *v1.TappController, key string) error {
	if cachedTappController.state != nil {
		// exist and the cluster name changed
		if cachedTappController.state.UID != tappController.UID {
			if err := c.processTappControllerDelete(cachedTappController, key); err != nil {
				return err
			}
		}
	}
	err := c.createTappControllerIfNeeded(key, cachedTappController, tappController)
	if err != nil {
		return err
	}

	cachedTappController.state = tappController
	// Always update the cache upon success.
	c.cache.set(key, cachedTappController)
	return nil
}

func (c *Controller) tappControllerReinitialize(key string, cachedTappController *cachedTappController, tappController *v1.TappController) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installTappController(tappController)
		if err == nil {
			tappController = tappController.DeepCopy()
			tappController.Status.Phase = v1.AddonPhaseChecking
			tappController.Status.Reason = ""
			tappController.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(tappController)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		// First, rollback the tappController
		if err := c.uninstallTappController(tappController); err != nil {
			log.Error("Uninstall tapp controller error.")
			return true, err
		}
		if tappController.Status.RetryCount == maxRetryCount {
			tappController = tappController.DeepCopy()
			tappController.Status.Phase = v1.AddonPhaseFailed
			tappController.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(tappController)
			if err != nil {
				log.Error("Update tapp controller error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		tappController = tappController.DeepCopy()
		tappController.Status.Phase = v1.AddonPhaseReinitializing
		tappController.Status.Reason = err.Error()
		tappController.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		tappController.Status.RetryCount++
		err = c.persistUpdate(tappController)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createTappControllerIfNeeded(key string, cachedTappController *cachedTappController, tappController *v1.TappController) error {
	switch tappController.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("Tapp controller will be created", log.String("tappControllerName", key))
		err := c.installTappController(tappController)
		if err == nil {
			tappController = tappController.DeepCopy()
			tappController.Status.Version = tappController.Spec.Version
			tappController.Status.Phase = v1.AddonPhaseChecking
			tappController.Status.Reason = ""
			tappController.Status.RetryCount = 0
			return c.persistUpdate(tappController)
		}
		tappController = tappController.DeepCopy()
		tappController.Status.Version = tappController.Spec.Version
		tappController.Status.Phase = v1.AddonPhaseReinitializing
		tappController.Status.Reason = err.Error()
		tappController.Status.RetryCount = 1
		tappController.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(tappController)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(tappController.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go wait.Poll(waitTime, timeOut, c.tappControllerReinitialize(key, cachedTappController, tappController))
	case v1.AddonPhaseChecking:
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, true)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkTappControllerStatus(tappController, key, initDelay))
			}()
		}
	case v1.AddonPhaseRunning:
		if needUpgrade(tappController) {
			c.health.Delete(key)
			tappController = tappController.DeepCopy()
			tappController.Status.Phase = v1.AddonPhaseUpgrading
			tappController.Status.Reason = ""
			tappController.Status.RetryCount = 0
			return c.persistUpdate(tappController)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go wait.PollImmediateUntil(5*time.Minute, c.watchTappControllerHealth(key), c.stopCh)
		}
	case v1.AddonPhaseUpgrading:
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, true)
			upgradeDelay := time.Now().Add(timeOut)
			go func() {
				defer c.upgrading.Delete(key)
				wait.PollImmediate(5*time.Second, timeOut, c.upgradeTappController(tappController, key, upgradeDelay))
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("Tapp controller is error", log.String("tappControllerName", key))
		c.health.Delete(key)
		c.checking.Delete(key)
		c.upgrading.Delete(key)
	}
	return nil
}

func needUpgrade(tappController *v1.TappController) bool {
	return tappController.Spec.Version != tappController.Status.Version
}

func (c *Controller) installTappController(tappController *v1.TappController) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(tappController.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	// ServiceAccount TappController
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountTappController()); err != nil {
		return err
	}
	// ClusterRoleBinding TappController
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(crbTappController()); err != nil {
		return err
	}
	// Deployment TappController
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deploymentTappController(images.Get(tappController.Spec.Version))); err != nil {
		return err
	}
	// Service TappController
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(serviceTappController()); err != nil {
		return err
	}

	return nil
}

func serviceAccountTappController() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountTappControllerName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func crbTappController() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbTappControllerName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountTappControllerName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func deploymentTappController(components images.Components) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployTappControllerName,
			Labels:    map[string]string{"app": controllerName},
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": controllerName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": controllerName},
				},
				Spec: corev1.PodSpec{
					PriorityClassName:  "system-cluster-critical",
					ServiceAccountName: svcAccountTappControllerName,
					Containers: []corev1.Container{
						{
							Name:  controllerName,
							Image: components.TappController.FullName(),
							Args:  []string{"--v", "3", "--register-admission", "true"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									// TODO: add support for configuring them
									corev1.ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(512*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}
}

func serviceTappController() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcTappControllerName,
			Labels:    map[string]string{"app": controllerName},
			Namespace: metav1.NamespaceSystem,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "https", Port: 443, TargetPort: intstr.FromInt(8443)},
			},
			Selector: map[string]string{"app": controllerName},
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }

func (c *Controller) uninstallTappController(tappController *v1.TappController) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(tappController.Spec.ClusterName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	// Deployment TappController
	deployTappControllerErr := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(deployTappControllerName, &metav1.DeleteOptions{})
	// ClusterRoleBinding TappController
	crbTappControllerErr := kubeClient.RbacV1().ClusterRoleBindings().Delete(crbTappControllerName, &metav1.DeleteOptions{})
	// ServiceAccount TappController
	svcAccountTappControllerErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(svcAccountTappControllerName, &metav1.DeleteOptions{})
	// Service TappController
	svcTappControllerErr := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(svcTappControllerName, &metav1.DeleteOptions{})

	if (deployTappControllerErr != nil && !errors.IsNotFound(deployTappControllerErr)) ||
		(crbTappControllerErr != nil && !errors.IsNotFound(crbTappControllerErr)) ||
		(svcAccountTappControllerErr != nil && !errors.IsNotFound(svcAccountTappControllerErr)) ||
		(svcTappControllerErr != nil && !errors.IsNotFound(svcTappControllerErr)) {
		return normalerrors.New("delete tapp controller error")
	}
	return nil
}

func (c *Controller) watchTappControllerHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check tapp controller in cluster health", log.String("tappControllerName", key))
		tappController, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.client.PlatformV1().Clusters().Get(tappController.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.health.Load(cluster.Name); !ok {
			log.Info("Health check over.")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		// TODO: check tapp controller service
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployTappControllerName, metav1.GetOptions{}); err != nil {
			tappController = tappController.DeepCopy()
			tappController.Status.Phase = v1.AddonPhaseFailed
			tappController.Status.Reason = "Tapp controller is not healthy."
			if err = c.persistUpdate(tappController); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
}

func (c *Controller) checkTappControllerStatus(tappController *v1.TappController, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check tapp controller health", log.String("tappControllerName", tappController.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(tappController.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.checking.Load(key); !ok {
			log.Debug("Checking over tapp controller addon status")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		tappController, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		if deploy, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployTappControllerName, metav1.GetOptions{}); err != nil ||
			(deploy.Spec.Replicas != nil && deploy.Status.AvailableReplicas < *deploy.Spec.Replicas) {
			if time.Now().After(initDelay) {
				tappController = tappController.DeepCopy()
				tappController.Status.Phase = v1.AddonPhaseFailed
				tappController.Status.Reason = "Tapp controller is not healthy."
				if err = c.persistUpdate(tappController); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		tappController = tappController.DeepCopy()
		tappController.Status.Phase = v1.AddonPhaseRunning
		tappController.Status.Reason = ""
		if err = c.persistUpdate(tappController); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) upgradeTappController(tappController *v1.TappController, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade tapp controller", log.String("tappControllerName", tappController.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(tappController.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading tapp controller", log.String("tappControllerName", tappController.Name))
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		tappController, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		newImage := images.Get(tappController.Spec.Version).TappController.FullName()

		patch := fmt.Sprintf(`[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`, newImage)
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Patch(deployTappControllerName, types.JSONPatchType, []byte(patch)); err != nil {
			if time.Now().After(initDelay) {
				tappController = tappController.DeepCopy()
				tappController.Status.Phase = v1.AddonPhaseFailed
				tappController.Status.Reason = "Failed to upgrade tapp controller."
				if err = c.persistUpdate(tappController); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		tappController = tappController.DeepCopy()
		tappController.Status.Version = tappController.Spec.Version
		tappController.Status.Phase = v1.AddonPhaseChecking
		tappController.Status.Reason = ""
		if err = c.persistUpdate(tappController); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) persistUpdate(tappController *v1.TappController) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.PlatformV1().TappControllers().UpdateStatus(tappController)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to tappController that no longer exists", log.String("clusterName", tappController.Spec.ClusterName), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to tappController '%s' that has been changed since we received it: %v", tappController.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of tappController '%s/%s'", tappController.Spec.ClusterName, tappController.Status.Phase), log.String("clusterName", tappController.Spec.ClusterName), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}
