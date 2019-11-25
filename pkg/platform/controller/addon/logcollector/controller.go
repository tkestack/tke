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

package logcollector

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/logcollector/images"

	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/metrics"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	controllerName = "log-collector-controller"

	crbName        = "log-collector-role-binding"
	svcAccountName = "log-collector"
	daemonSetName  = "log-collector"

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	timeOut       = 5 * time.Minute
	maxRetryCount = 5

	upgradePatchTemplate = `[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`
)

// Controller is responsible for performing actions dependent upon a LogCollector phase.
type Controller struct {
	client       clientset.Interface
	cache        *logcollectorCache
	health       sync.Map
	checking     sync.Map
	upgrading    sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.LogCollectorLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new LogCollector Controller object.
func NewController(client clientset.Interface, informer platformv1informer.LogCollectorInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client: client,
		cache:  &logcollectorCache{lcMap: make(map[string]*cachedLogCollector)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueLogCollector,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldLogCollector, ok1 := oldObj.(*v1.LogCollector)
				curLogCollector, ok2 := newObj.(*v1.LogCollector)
				if ok1 && ok2 && controller.needsUpdate(oldLogCollector, curLogCollector) {
					controller.enqueueLogCollector(newObj)
				}
			},
			DeleteFunc: controller.enqueueLogCollector,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.LogCollector, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueLogCollector(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for LogCollector object",
			log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(old *v1.LogCollector, new *v1.LogCollector) bool {
	return !reflect.DeepEqual(old, new)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting LogCollector controller")
	defer log.Info("Shutting down LogCollector controller")

	if !cache.WaitForCacheSync(stopCh, c.listerSynced) {
		return fmt.Errorf("failed to wait for LogCollector cache to sync")
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

	err := c.syncLogCollector(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing LogCollector %s (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncLogCollector will sync the LogCollector with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncLogCollector(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing LogCollector", log.String("name", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// LogCollector holds the latest LogCollector info from apiserver.
	LogCollector, err := c.lister.Get(name)
	switch {
	case k8serrors.IsNotFound(err):
		log.Info("LogCollector has been deleted. Attempting to cleanup resources", log.String("name", key))
		err = c.processLogCollectorDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve LogCollector from store", log.String("name", key), log.Err(err))
	default:
		cachedLogCollector := c.cache.getOrCreate(key)
		err = c.processLogCollectorUpdate(cachedLogCollector, LogCollector, key)
	}

	return err
}

func (c *Controller) processLogCollectorDeletion(key string) error {
	cachedLogCollector, ok := c.cache.get(key)
	if !ok {
		log.Error("LogCollector not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processLogCollectorDelete(cachedLogCollector, key)
}

func (c *Controller) processLogCollectorDelete(cachedLogCollector *cachedLogCollector, key string) error {
	log.Info("LogCollector will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the LogCollector cache", log.String("name", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the LogCollector health cache", log.String("name", key))
		c.health.Delete(key)
	}

	LogCollector := cachedLogCollector.state
	return c.uninstallLogCollector(LogCollector)
}

func (c *Controller) processLogCollectorUpdate(cachedLogCollector *cachedLogCollector, LogCollector *v1.LogCollector, key string) error {
	if cachedLogCollector.state != nil {
		// exist and the cluster name changed
		if cachedLogCollector.state.UID != LogCollector.UID {
			if err := c.processLogCollectorDelete(cachedLogCollector, key); err != nil {
				return err
			}
		}
	}
	err := c.createLogCollectorIfNeeded(key, cachedLogCollector, LogCollector)
	if err != nil {
		return err
	}

	cachedLogCollector.state = LogCollector
	// Always update the cache upon success.
	c.cache.set(key, cachedLogCollector)
	return nil
}

func (c *Controller) logCollectorReinitialize(
	key string,
	cachedLogCollector *cachedLogCollector,
	LogCollector *v1.LogCollector) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installLogCollector(LogCollector)
		if err == nil {
			LogCollector = LogCollector.DeepCopy()
			LogCollector.Status.Phase = v1.AddonPhaseChecking
			LogCollector.Status.Reason = ""
			LogCollector.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(LogCollector)
			if err != nil {
				return true, err
			}
			return true, nil
		}

		// First, rollback the LogCollector.
		log.Info("Rollback LogCollector",
			log.String("name", LogCollector.Name),
			log.String("clusterName", LogCollector.Spec.ClusterName))
		if err := c.uninstallLogCollector(LogCollector); err != nil {
			log.Error("Uninstall LogCollector failed", log.Err(err))
			return true, err
		}

		if LogCollector.Status.RetryCount == maxRetryCount {
			LogCollector = LogCollector.DeepCopy()
			LogCollector.Status.Phase = v1.AddonPhaseFailed
			LogCollector.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(LogCollector)
			if err != nil {
				log.Error("Update LogCollector failed", log.Err(err))
				return true, err
			}
			return true, nil
		}

		// Add the retry count will trigger reinitialize function from the persistent controller again.
		LogCollector = LogCollector.DeepCopy()
		LogCollector.Status.Phase = v1.AddonPhaseReinitializing
		LogCollector.Status.Reason = err.Error()
		LogCollector.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		LogCollector.Status.RetryCount++
		return true, c.persistUpdate(LogCollector)
	}
}

func (c *Controller) createLogCollectorIfNeeded(
	key string,
	cachedLogCollector *cachedLogCollector,
	LogCollector *v1.LogCollector) error {
	switch LogCollector.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("LogCollector will be created", log.String("name", key))
		err := c.installLogCollector(LogCollector)
		if err == nil {
			LogCollector = LogCollector.DeepCopy()
			fillOperatorStatus(LogCollector)
			LogCollector.Status.Phase = v1.AddonPhaseChecking
			LogCollector.Status.Reason = ""
			LogCollector.Status.RetryCount = 0
			return c.persistUpdate(LogCollector)
		}
		// Install LogCollector failed.
		LogCollector = LogCollector.DeepCopy()
		fillOperatorStatus(LogCollector)
		LogCollector.Status.Phase = v1.AddonPhaseReinitializing
		LogCollector.Status.Reason = err.Error()
		LogCollector.Status.RetryCount = 1
		LogCollector.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(LogCollector)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(LogCollector.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go func() {
			reInitialErr := wait.Poll(waitTime, timeOut,
				c.logCollectorReinitialize(key, cachedLogCollector, LogCollector))
			if reInitialErr != nil {
				log.Error("Reinitialize LogCollector failed",
					log.String("name", LogCollector.Name),
					log.String("clusterName", LogCollector.Spec.ClusterName),
					log.Err(reInitialErr))
			}
		}()
	case v1.AddonPhaseChecking:
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, true)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				checkStatusErr := wait.PollImmediate(5*time.Second, 5*time.Minute+10*time.Second,
					c.checkLogCollectorStatus(LogCollector, key, initDelay))
				if checkStatusErr != nil {
					log.Error("Check status of LogCollector failed",
						log.String("name", LogCollector.Name),
						log.String("clusterName", LogCollector.Spec.ClusterName),
						log.Err(checkStatusErr))
				}
			}()
		}
	case v1.AddonPhaseRunning:
		if c.needUpgrade(LogCollector) {
			c.health.Delete(key)
			LogCollector = LogCollector.DeepCopy()
			LogCollector.Status.Phase = v1.AddonPhaseUpgrading
			LogCollector.Status.Reason = ""
			LogCollector.Status.RetryCount = 0
			return c.persistUpdate(LogCollector)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go func() {
				healthErr := wait.PollImmediateUntil(5*time.Minute,
					c.watchLogCollectorHealth(key), c.stopCh)
				if healthErr != nil {
					log.Error("Watch health of LogCollector failed",
						log.String("name", LogCollector.Name),
						log.String("clusterName", LogCollector.Spec.ClusterName),
						log.Err(healthErr))
				}
			}()
		}
	case v1.AddonPhaseUpgrading:
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, true)
			upgradeDelay := time.Now().Add(timeOut)
			go func() {
				defer c.upgrading.Delete(key)
				upgradeErr := wait.PollImmediate(5*time.Second, timeOut,
					c.upgradeLogCollector(LogCollector, key, upgradeDelay))
				if upgradeErr != nil {
					log.Error("Upgrade LogCollector failed",
						log.String("name", LogCollector.Name),
						log.String("clusterName", LogCollector.Spec.ClusterName),
						log.Err(upgradeErr))
				}
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("LogCollector failed", log.String("name", key))
		c.health.Delete(key)
		c.checking.Delete(key)
		c.upgrading.Delete(key)
	}
	return nil
}

func (c *Controller) needUpgrade(LogCollector *v1.LogCollector) bool {
	return LogCollector.Spec.Version != LogCollector.Status.Version
}

func (c *Controller) installLogCollector(LogCollector *v1.LogCollector) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(LogCollector.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	// Create ServiceAccount.
	if err := c.installSVC(LogCollector, kubeClient); err != nil {
		return err
	}

	// Create ClusterRoleBinding.
	if err := c.installCRB(LogCollector, kubeClient); err != nil {
		return err
	}

	// Create Deployment.
	return c.installDaemonSet(LogCollector, kubeClient)
}

func (c *Controller) installSVC(LogCollector *v1.LogCollector, kubeClient kubernetes.Interface) error {
	svc := genServiceAccount()
	svcClient := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem)

	_, err := svcClient.Get(svc.Name, metav1.GetOptions{})
	if err == nil {
		log.Info("ServiceAccount of LogCollector is already created",
			log.String("name", LogCollector.Name))
		return nil
	}

	if k8serrors.IsNotFound(err) {
		_, err = svcClient.Create(svc)
		return err
	}

	return fmt.Errorf("get svc failed: %v", err)
}

func (c *Controller) installCRB(LogCollector *v1.LogCollector, kubeClient kubernetes.Interface) error {
	crb := genCRB()
	crbClient := kubeClient.RbacV1().ClusterRoleBindings()

	oldCRB, err := crbClient.Get(crb.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = crbClient.Create(crb)
			return err
		}
		return fmt.Errorf("get crb failed: %v", err)
	}

	if equality.Semantic.DeepEqual(oldCRB.RoleRef, crb.RoleRef) &&
		equality.Semantic.DeepEqual(oldCRB.Subjects, crb.Subjects) {
		log.Info("ClusterRoleBinding of LogCollector is already created",
			log.String("name", LogCollector.Name))
		return nil
	}

	newCRB := oldCRB.DeepCopy()
	newCRB.RoleRef = crb.RoleRef
	newCRB.Subjects = crb.Subjects
	_, err = crbClient.Update(newCRB)

	return err
}

func (c *Controller) installDaemonSet(
	LogCollector *v1.LogCollector,
	kubeClient kubernetes.Interface) error {
	daemon := c.genDaemonSet(LogCollector.Spec.Version)
	daemonClient := kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem)

	oldDaemon, err := daemonClient.Get(daemon.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = daemonClient.Create(daemon)
			return err
		}
		return fmt.Errorf("get daemonSet failed: %v", err)
	}

	newDaemon := oldDaemon.DeepCopy()
	newDaemon.Labels = daemon.Labels
	newDaemon.Spec = daemon.Spec
	_, err = daemonClient.Update(newDaemon)

	return err
}

func genServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func genCRB() *rbacv1.ClusterRoleBinding {
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
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func (c *Controller) genDaemonSet(version string) *appsv1.DaemonSet {
	daemon := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonSetName,
			Labels:    map[string]string{"app": controllerName},
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": controllerName},
			},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": controllerName},
				},
				Spec: corev1.PodSpec{
					PriorityClassName:  "system-cluster-critical",
					ServiceAccountName: svcAccountName,
					HostNetwork:        true,
					Tolerations: []corev1.Toleration{
						{Key: "node-role.kubernetes.io/master", Effect: corev1.TaintEffectNoSchedule},
					},
					Containers: []corev1.Container{
						{
							Name:  daemonSetName,
							Image: images.Get(version).LogCollector.FullName(),
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{"SYS_ADMIN"},
								},
							},
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
							Env: []corev1.EnvVar{
								{Name: "EXCEPTION_DSN", Value: "http://fb8e7951085148148d727465de5b5d34:c2b917404bfb47ce9e1b3f7d18dfc70f@exception.log.ccs.cloud.tencent.com/3"},
								{Name: "HOST_ROOTFS", Value: "/rootfs"},
								{Name: "K8S_NODE_NAME", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}}},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "sock", MountPath: "/var/run/docker.sock"},
								{Name: "rootfs", MountPath: "/rootfs"},
								{Name: "varlogpods", MountPath: "/var/log/pods"},
								{Name: "varlogcontainers", MountPath: "/var/log/containers"},
								{Name: "varlibdockercontainers", MountPath: "/var/lib/docker/containers"},
								{Name: "datadocker", MountPath: "/data/docker"},
								{Name: "optdocker", MountPath: "/opt/docker"},
								{Name: "tdagent", MountPath: "/var/log/td-agent"},
								{Name: "localtime", MountPath: "/etc/localtime"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{Name: "sock", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/run/docker.sock"}}},
						{Name: "rootfs", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/"}}},
						{Name: "varlogpods", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/log/pods"}}},
						{Name: "varlogcontainers", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/log/containers"}}},
						{Name: "varlibdockercontainers", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/docker/containers"}}},
						{Name: "datadocker", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/data/docker"}}},
						{Name: "optdocker", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/opt/docker"}}},
						{Name: "tdagent", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/tmp/ccs-log-collector"}}},
						{Name: "localtime", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/localtime"}}},
					},
				},
			},
		},
	}

	return daemon
}

func boolPtr(value bool) *bool {
	return &value
}

func (c *Controller) uninstallLogCollector(LogCollector *v1.LogCollector) error {
	log.Info("Start to uninstall LogCollector",
		log.String("name", LogCollector.Name),
		log.String("clusterName", LogCollector.Spec.ClusterName))

	cluster, err := c.client.PlatformV1().Clusters().Get(LogCollector.Spec.ClusterName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	// Delete the operator daemonSet.
	clearDaemonSetErr := kubeClient.AppsV1().
		DaemonSets(metav1.NamespaceSystem).Delete(daemonSetName, &metav1.DeleteOptions{})
	// Delete the ClusterRoleBinding.
	clearCRBErr := kubeClient.RbacV1().
		ClusterRoleBindings().Delete(crbName, &metav1.DeleteOptions{})
	// Delete the ServiceAccount.
	clearSVCErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).
		Delete(svcAccountName, &metav1.DeleteOptions{})

	failed := false

	if clearDaemonSetErr != nil && !k8serrors.IsNotFound(clearDaemonSetErr) {
		failed = true
		log.Error("delete daemonSet for LogCollector failed",
			log.String("name", LogCollector.Name),
			log.String("clusterName", LogCollector.Spec.ClusterName),
			log.Err(clearDaemonSetErr))
	}

	if clearCRBErr != nil && !k8serrors.IsNotFound(clearCRBErr) {
		failed = true
		log.Error("delete crb for LogCollector failed",
			log.String("name", LogCollector.Name),
			log.String("clusterName", LogCollector.Spec.ClusterName),
			log.Err(clearCRBErr))
	}

	if clearSVCErr != nil && !k8serrors.IsNotFound(clearSVCErr) {
		failed = true
		log.Error("delete service account for LogCollector failed",
			log.String("name", LogCollector.Name),
			log.String("clusterName", LogCollector.Spec.ClusterName),
			log.Err(clearSVCErr))
	}

	if failed {
		return errors.New("delete LogCollector failed")
	}

	return nil
}

func (c *Controller) watchLogCollectorHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		LogCollector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		log.Info("Start check health of LogCollector", log.String("name", LogCollector.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(LogCollector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
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

		_, err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).
			Get(daemonSetName, metav1.GetOptions{})
		if err != nil {
			LogCollector = LogCollector.DeepCopy()
			LogCollector.Status.Phase = v1.AddonPhaseFailed
			LogCollector.Status.Reason = "LogCollector is not healthy."
			if err = c.persistUpdate(LogCollector); err != nil {
				return false, err
			}
			return true, nil
		}

		return false, nil
	}
}

func (c *Controller) checkLogCollectorStatus(
	LogCollector *v1.LogCollector,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check LogCollector health", log.String("name", LogCollector.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(LogCollector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}

		if _, ok := c.checking.Load(key); !ok {
			log.Debug("Checking over LogCollector addon status")
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		LogCollector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		daemon, err := kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).
			Get(daemonSetName, metav1.GetOptions{})
		if err != nil || daemon.Status.DesiredNumberScheduled == 0 ||
			daemon.Status.NumberAvailable < daemon.Status.DesiredNumberScheduled {
			if time.Now().After(initDelay) {
				LogCollector = LogCollector.DeepCopy()
				LogCollector.Status.Phase = v1.AddonPhaseFailed
				LogCollector.Status.Reason = fmt.Sprintf("Log Collector is not healthy")
				if err = c.persistUpdate(LogCollector); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		LogCollector = LogCollector.DeepCopy()
		LogCollector.Status.Phase = v1.AddonPhaseRunning
		LogCollector.Status.Reason = ""
		if err = c.persistUpdate(LogCollector); err != nil {
			return false, err
		}

		return true, nil
	}
}

func (c *Controller) upgradeLogCollector(
	LogCollector *v1.LogCollector,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade LogCollector", log.String("name", LogCollector.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(LogCollector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}

		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading LogCollector", log.String("name", LogCollector.Name))
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		LogCollector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		patch := fmt.Sprintf(upgradePatchTemplate, images.Get(LogCollector.Spec.Version).LogCollector.FullName())

		_, err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).
			Patch(daemonSetName, types.JSONPatchType, []byte(patch))
		if err != nil {
			if time.Now().After(initDelay) {
				LogCollector = LogCollector.DeepCopy()
				LogCollector.Status.Phase = v1.AddonPhaseFailed
				LogCollector.Status.Reason = "Failed to upgrade LogCollector."
				if err = c.persistUpdate(LogCollector); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		LogCollector = LogCollector.DeepCopy()
		fillOperatorStatus(LogCollector)
		LogCollector.Status.Phase = v1.AddonPhaseChecking
		LogCollector.Status.Reason = ""
		if err = c.persistUpdate(LogCollector); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) persistUpdate(LogCollector *v1.LogCollector) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.PlatformV1().LogCollectors().UpdateStatus(LogCollector)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if k8serrors.IsNotFound(err) {
			log.Info("Not persisting update to LogCollector that no longer exists",
				log.String("clusterName", LogCollector.Spec.ClusterName), log.Err(err))
			return nil
		}
		if k8serrors.IsConflict(err) {
			return fmt.Errorf("not persisting update to LogCollector %q that has been changed since we received it: %v", LogCollector.Spec.ClusterName, err)
		}
		log.Warn("Failed to persist updated status of LogCollector",
			log.String("name", LogCollector.Name),
			log.String("clusterName", LogCollector.Spec.ClusterName),
			log.String("phase", string(LogCollector.Status.Phase)), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}

func fillOperatorStatus(LogCollector *v1.LogCollector) {
	LogCollector.Status.Version = LogCollector.Spec.Version
}
