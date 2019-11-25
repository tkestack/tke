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

package cronhpa

import (
	normalerrors "errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/cronhpa/images"

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
	controllerName        = "cron-hpa-controller"
	deployCronHPAName     = "cron-hpa-controller"
	svcCronHPAName        = "cron-hpa-controller"
	svcAccountCronHPAName = "cron-hpa-controller"
	crbCronHPAName        = "cron-hpa-controller"
)

// Controller is responsible for performing actions dependent upon a CronHPA phase.
type Controller struct {
	client       clientset.Interface
	cache        *cronHPACache
	health       sync.Map
	checking     sync.Map
	upgrading    sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.CronHPALister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, informer platformv1informer.CronHPAInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client: client,
		cache:  &cronHPACache{cronHPAMap: make(map[string]*cachedCronHPA)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the CronHPA informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueCronHPA,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCronHPA, ok1 := oldObj.(*v1.CronHPA)
				curCronHPA, ok2 := newObj.(*v1.CronHPA)
				if ok1 && ok2 && controller.needsUpdate(oldCronHPA, curCronHPA) {
					controller.enqueueCronHPA(newObj)
				}
			},
			DeleteFunc: controller.enqueueCronHPA,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.CronHPA, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueCronHPA(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *v1.CronHPA, new *v1.CronHPA) bool {
	return !reflect.DeepEqual(old, new)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting CronHPA controller")
	defer log.Info("Shutting down CronHPA controller")

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

	err := c.syncCronHPA(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing CronHPA %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncCronHPA will sync the CronHPA with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncCronHPA(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing CronHPA", log.String("CronHPA", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// cronHPA holds the latest cronHPA info from apiserver
	cronHPA, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("CronHPA has been deleted. Attempting to cleanup resources", log.String("CronHPA", key))
		err = c.processCronHPADeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve CronHPA from store", log.String("CronHPA", key), log.Err(err))
	default:
		cachedCronHPA := c.cache.getOrCreate(key)
		err = c.processCronHPAUpdate(cachedCronHPA, cronHPA, key)
	}
	return err
}

func (c *Controller) processCronHPADeletion(key string) error {
	cachedCronHPA, ok := c.cache.get(key)
	if !ok {
		log.Error("CronHPA not in cache even though the watcher thought it was. Ignoring the deletion", log.String("CronHPA", key))
		return nil
	}
	return c.processCronHPADelete(cachedCronHPA, key)
}

func (c *Controller) processCronHPADelete(cachedCronHPA *cachedCronHPA, key string) error {
	log.Info("CronHPA will be dropped", log.String("CronHPA", key))

	if c.cache.Exist(key) {
		log.Info("Delete the CronHPA cache", log.String("CronHPA", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the CronHPA health cache", log.String("CronHPA", key))
		c.health.Delete(key)
	}

	cronHPA := cachedCronHPA.state
	return c.uninstallCronHPA(cronHPA)
}

func (c *Controller) processCronHPAUpdate(cachedCronHPA *cachedCronHPA, cronHPA *v1.CronHPA, key string) error {
	if cachedCronHPA.state != nil {
		// exist and the cluster name changed
		if cachedCronHPA.state.UID != cronHPA.UID {
			if err := c.processCronHPADelete(cachedCronHPA, key); err != nil {
				return err
			}
		}
	}
	err := c.createCronHPAIfNeeded(key, cachedCronHPA, cronHPA)
	if err != nil {
		return err
	}

	cachedCronHPA.state = cronHPA
	// Always update the cache upon success.
	c.cache.set(key, cachedCronHPA)
	return nil
}

func (c *Controller) cronHPAReinitialize(key string, cachedCronHPA *cachedCronHPA, cronHPA *v1.CronHPA) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installCronHPA(cronHPA)
		if err == nil {
			cronHPA = cronHPA.DeepCopy()
			cronHPA.Status.Phase = v1.AddonPhaseChecking
			cronHPA.Status.Reason = ""
			cronHPA.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(cronHPA)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		// First, rollback the cronHPA
		if err := c.uninstallCronHPA(cronHPA); err != nil {
			log.Error("Uninstall CronHPA error.")
			return true, err
		}
		if cronHPA.Status.RetryCount == maxRetryCount {
			cronHPA = cronHPA.DeepCopy()
			cronHPA.Status.Phase = v1.AddonPhaseFailed
			cronHPA.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(cronHPA)
			if err != nil {
				log.Error("Update CronHPA error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		cronHPA = cronHPA.DeepCopy()
		cronHPA.Status.Phase = v1.AddonPhaseReinitializing
		cronHPA.Status.Reason = err.Error()
		cronHPA.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		cronHPA.Status.RetryCount++
		err = c.persistUpdate(cronHPA)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createCronHPAIfNeeded(key string, cachedCronHPA *cachedCronHPA, cronHPA *v1.CronHPA) error {
	switch cronHPA.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("CronHPA will be created", log.String("CronHPA", key))
		err := c.installCronHPA(cronHPA)
		if err == nil {
			cronHPA = cronHPA.DeepCopy()
			cronHPA.Status.Version = cronHPA.Spec.Version
			cronHPA.Status.Phase = v1.AddonPhaseChecking
			cronHPA.Status.Reason = ""
			cronHPA.Status.RetryCount = 0
			return c.persistUpdate(cronHPA)
		}
		cronHPA = cronHPA.DeepCopy()
		cronHPA.Status.Version = cronHPA.Spec.Version
		cronHPA.Status.Phase = v1.AddonPhaseReinitializing
		cronHPA.Status.Reason = err.Error()
		cronHPA.Status.RetryCount = 1
		cronHPA.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(cronHPA)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(cronHPA.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go wait.Poll(waitTime, timeOut, c.cronHPAReinitialize(key, cachedCronHPA, cronHPA))
	case v1.AddonPhaseChecking:
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, true)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkCronHPAStatus(cronHPA, key, initDelay))
			}()
		}
	case v1.AddonPhaseRunning:
		if needUpgrade(cronHPA) {
			c.health.Delete(key)
			cronHPA = cronHPA.DeepCopy()
			cronHPA.Status.Phase = v1.AddonPhaseUpgrading
			cronHPA.Status.Reason = ""
			cronHPA.Status.RetryCount = 0
			return c.persistUpdate(cronHPA)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go wait.PollImmediateUntil(5*time.Minute, c.watchCronHPAHealth(key), c.stopCh)
		}
	case v1.AddonPhaseUpgrading:
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, true)
			upgradeDelay := time.Now().Add(timeOut)
			go func() {
				defer c.upgrading.Delete(key)
				wait.PollImmediate(5*time.Second, timeOut, c.upgradeCronHPA(cronHPA, key, upgradeDelay))
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("CronHPA is error", log.String("CronHPA", key))
		c.health.Delete(key)
		c.checking.Delete(key)
		c.upgrading.Delete(key)
	}
	return nil
}

func needUpgrade(cronHPA *v1.CronHPA) bool {
	return cronHPA.Spec.Version != cronHPA.Status.Version
}

func (c *Controller) installCronHPA(cronHPA *v1.CronHPA) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(cronHPA.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	// ServiceAccount CronHPA
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountCronHPA()); err != nil {
		return err
	}
	// ClusterRoleBinding CronHPA
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(crbCronHPA()); err != nil {
		return err
	}
	// Deployment CronHPA
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deploymentCronHPA(cronHPA.Spec.Version)); err != nil {
		return err
	}
	// Service CronHPA
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(serviceCronHPA()); err != nil {
		return err
	}

	return nil
}

func serviceAccountCronHPA() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountCronHPAName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func crbCronHPA() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbCronHPAName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountCronHPAName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func deploymentCronHPA(cronHPAVersion string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployCronHPAName,
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
					ServiceAccountName: svcAccountCronHPAName,
					Containers: []corev1.Container{
						{
							Name:  controllerName,
							Image: images.Get(cronHPAVersion).CronHPA.FullName(),
							Args:  []string{"--v", "3", "--stderrthreshold", "0", "--register-admission", "true"},
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

func serviceCronHPA() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcCronHPAName,
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

func (c *Controller) uninstallCronHPA(cronHPA *v1.CronHPA) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(cronHPA.Spec.ClusterName, metav1.GetOptions{})
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
	// Deployment CronHPA
	deployCronHPAErr := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(deployCronHPAName, &metav1.DeleteOptions{})
	// ClusterRoleBinding CronHPA
	crbCronHPAErr := kubeClient.RbacV1().ClusterRoleBindings().Delete(crbCronHPAName, &metav1.DeleteOptions{})
	// ServiceAccount CronHPA
	svcAccountCronHPAErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(svcAccountCronHPAName, &metav1.DeleteOptions{})
	// Service CronHPA
	svcCronHPAErr := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(svcCronHPAName, &metav1.DeleteOptions{})

	if (deployCronHPAErr != nil && !errors.IsNotFound(deployCronHPAErr)) ||
		(crbCronHPAErr != nil && !errors.IsNotFound(crbCronHPAErr)) ||
		(svcAccountCronHPAErr != nil && !errors.IsNotFound(svcAccountCronHPAErr)) ||
		(svcCronHPAErr != nil && !errors.IsNotFound(svcCronHPAErr)) {
		return normalerrors.New("delete CronHPA error")
	}
	return nil
}

func (c *Controller) watchCronHPAHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check CronHPA in cluster health", log.String("CronHPA", key))
		cronHPA, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.client.PlatformV1().Clusters().Get(cronHPA.Spec.ClusterName, metav1.GetOptions{})
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
		// TODO: check CronHPA controller service
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployCronHPAName, metav1.GetOptions{}); err != nil {
			cronHPA = cronHPA.DeepCopy()
			cronHPA.Status.Phase = v1.AddonPhaseFailed
			cronHPA.Status.Reason = "CronHPA is not healthy."
			if err = c.persistUpdate(cronHPA); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
}

func (c *Controller) checkCronHPAStatus(cronHPA *v1.CronHPA, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check CronHPA health", log.String("CronHPA", cronHPA.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(cronHPA.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.checking.Load(key); !ok {
			log.Debug("Checking over CronHPA addon status")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		cronHPA, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		if deploy, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployCronHPAName, metav1.GetOptions{}); err != nil ||
			(deploy.Spec.Replicas != nil && deploy.Status.AvailableReplicas < *deploy.Spec.Replicas) {
			if time.Now().After(initDelay) {
				cronHPA = cronHPA.DeepCopy()
				cronHPA.Status.Phase = v1.AddonPhaseFailed
				cronHPA.Status.Reason = "CronHPA is not healthy."
				if err = c.persistUpdate(cronHPA); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		cronHPA = cronHPA.DeepCopy()
		cronHPA.Status.Phase = v1.AddonPhaseRunning
		cronHPA.Status.Reason = ""
		if err = c.persistUpdate(cronHPA); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) upgradeCronHPA(cronHPA *v1.CronHPA, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade CronHPA", log.String("CronHPA", cronHPA.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(cronHPA.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading CronHPA", log.String("CronHPA", cronHPA.Name))
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		cronHPA, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		newImage := images.Get(cronHPA.Spec.Version).CronHPA.FullName()

		patch := fmt.Sprintf(`[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`, newImage)
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Patch(deployCronHPAName, types.JSONPatchType, []byte(patch)); err != nil {
			if time.Now().After(initDelay) {
				cronHPA = cronHPA.DeepCopy()
				cronHPA.Status.Phase = v1.AddonPhaseFailed
				cronHPA.Status.Reason = "Failed to upgrade CronHPA."
				if err = c.persistUpdate(cronHPA); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		cronHPA = cronHPA.DeepCopy()
		cronHPA.Status.Version = cronHPA.Spec.Version
		cronHPA.Status.Phase = v1.AddonPhaseChecking
		cronHPA.Status.Reason = ""
		if err = c.persistUpdate(cronHPA); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) persistUpdate(cronHPA *v1.CronHPA) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.PlatformV1().CronHPAs().UpdateStatus(cronHPA)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to cronHPA that no longer exists", log.String("clusterName", cronHPA.Spec.ClusterName), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to CronHPA '%s' that has been changed since we received it: %v", cronHPA.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of CronHPA '%s/%s'", cronHPA.Spec.ClusterName, cronHPA.Status.Phase), log.String("clusterName", cronHPA.Spec.ClusterName), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}
