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

package ipam

import (
	normalerrors "errors"
	"fmt"
	"reflect"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/ipam/images"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ipamMaxRetryCount = 5
	ipamTimeOut       = 5 * time.Minute
	ipamRetryInterval = 5 * time.Second
)

const (
	deployIPAMName     = "galaxy-ipam"
	svcIPAMAPIName     = "galaxy-ipam"
	svcAccountIPAMName = "galaxy-ipam"
	crbIPAMName        = "galaxy-ipam"
	crIPAMName         = "galaxy-ipam"
	cmIPAMName         = "galaxy-ipam-etc"
)

// Controller is responsible for performing actions dependent upon a ipam phase.
type Controller struct {
	client       clientset.Interface
	cache        *ipamCache
	health       *ipamHealth
	checking     *ipamChecking
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.IPAMLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object.
func NewController(client clientset.Interface, ipamInformer platformv1informer.IPAMInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{

		client:   client,
		cache:    &ipamCache{ipamMap: make(map[string]*cachedIPAM)},
		health:   &ipamHealth{ipamMap: make(map[string]*v1.IPAM)},
		checking: &ipamChecking{ipamMap: make(map[string]*v1.IPAM)},
		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ipam"),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("ipam_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the ipam informer event handlers
	ipamInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueIPAM,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldIPAM, ok1 := oldObj.(*v1.IPAM)
				curIPAM, ok2 := newObj.(*v1.IPAM)
				if ok1 && ok2 && controller.needsUpdate(oldIPAM, curIPAM) {
					controller.enqueueIPAM(newObj)
				}
			},
			DeleteFunc: controller.enqueueIPAM,
		},
		resyncPeriod,
	)
	controller.lister = ipamInformer.Lister()
	controller.listerSynced = ipamInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.IPAM, or a DeletionFinalStateUnknown marker ipam.
func (c *Controller) enqueueIPAM(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldIPAM *v1.IPAM, newIPAM *v1.IPAM) bool {
	return !reflect.DeepEqual(oldIPAM, newIPAM)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting ipam controller")
	defer log.Info("Shutting down ipam controller")

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

	err := c.syncIPAM(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing ipam %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncIPAM will sync the IPAM with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncIPAM(key string) error {
	startTime := time.Now()
	var cachedIPAM *cachedIPAM
	defer func() {
		log.Info("Finished syncing ipam", log.String("ipamName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// ipam holds the latest ipam info from apiserver
	ipam, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("IPAM has been deleted. Attempting to cleanup resources", log.String("ipam", key))
		err = c.processIPAMDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve ipam from store", log.String("ipam", key), log.Err(err))
	default:
		cachedIPAM = c.cache.getOrCreate(key)
		err = c.processIPAMUpdate(cachedIPAM, ipam, key)
	}
	return err
}

func (c *Controller) processIPAMDeletion(key string) error {
	cachedIPAM, ok := c.cache.get(key)
	if !ok {
		log.Error("ipam not in cache even though the watcher thought it was. Ignoring the deletion", log.String("ipamName", key))
		return nil
	}
	return c.processIPAMDelete(cachedIPAM, key)
}

func (c *Controller) processIPAMDelete(cachedIPAM *cachedIPAM, key string) error {
	log.Info("ipam will be dropped", log.String("ipamName", key))

	if c.cache.Exist(key) {
		log.Info("delete the ipam cache", log.String("ipamName", key))
		c.cache.delete(key)
	}

	if c.health.Exist(key) {
		log.Info("delete the ipam health cache", log.String("ipamName", key))
		c.health.Del(key)
	}

	ipam := cachedIPAM.state
	return c.uninstallIPAM(ipam)
}

func (c *Controller) processIPAMUpdate(cachedIPAM *cachedIPAM, ipam *v1.IPAM, key string) error {
	if cachedIPAM.state != nil {
		// exist and the cluster name changed
		if cachedIPAM.state.UID != ipam.UID {
			if err := c.processIPAMDelete(cachedIPAM, key); err != nil {
				return err
			}
		}
	}
	err := c.createIPAMIfNeeded(key, cachedIPAM, ipam)
	if err != nil {
		return err
	}

	cachedIPAM.state = ipam
	// Always update the cache upon success.
	c.cache.set(key, cachedIPAM)
	return nil
}

func (c *Controller) ipamReinitialize(key string, cachedIPAM *cachedIPAM, ipam *v1.IPAM) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installIPAM(ipam)
		if err == nil {
			ipam = ipam.DeepCopy()
			ipam.Status.Phase = v1.AddonPhaseChecking
			ipam.Status.Reason = ""
			ipam.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(ipam)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		glog.Errorf("fail to re-install galaxy-ipam %v", err)
		// First, rollback the ipam
		if err := c.uninstallIPAM(ipam); err != nil {
			log.Error("Uninstall ipam error.")
			return true, err
		}
		if ipam.Status.RetryCount == ipamMaxRetryCount {
			ipam = ipam.DeepCopy()
			ipam.Status.Phase = v1.AddonPhaseFailed
			ipam.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", ipamMaxRetryCount)
			err := c.persistUpdate(ipam)
			if err != nil {
				log.Error("Update ipam error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		ipam = ipam.DeepCopy()
		ipam.Status.Phase = v1.AddonPhaseReinitializing
		ipam.Status.Reason = err.Error()
		ipam.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		ipam.Status.RetryCount++
		err = c.persistUpdate(ipam)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createIPAMIfNeeded(key string, cachedIPAM *cachedIPAM, ipam *v1.IPAM) error {
	switch ipam.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("IPAM will be created", log.String("ipam", key))
		err := c.installIPAM(ipam)
		if err == nil {
			ipam = ipam.DeepCopy()
			ipam.Status.Phase = v1.AddonPhaseChecking
			ipam.Status.Reason = ""
			ipam.Status.RetryCount = 0
			return c.persistUpdate(ipam)
		}
		glog.Errorf("fail to create galaxy-ipam %v", err)
		ipam = ipam.DeepCopy()
		ipam.Status.Phase = v1.AddonPhaseReinitializing
		ipam.Status.Reason = err.Error()
		ipam.Status.RetryCount = 1
		ipam.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(ipam)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(ipam.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= ipamTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = ipamTimeOut - interval
		}
		go wait.Poll(waitTime, ipamTimeOut, c.ipamReinitialize(key, cachedIPAM, ipam))
	case v1.AddonPhaseChecking:
		if !c.checking.Exist(key) {
			c.checking.Set(ipam)
			initDelay := time.Now().Add(5 * time.Minute)
			go wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkIPAMStatus(ipam, key, initDelay))
		}
	case v1.AddonPhaseRunning:
		if !c.health.Exist(key) {
			c.health.Set(ipam)
			go wait.PollImmediateUntil(5*time.Minute, c.watchIPAMHealth(key), c.stopCh)
		}
	case v1.AddonPhaseFailed:
		log.Info("IPAM is error", log.String("ipam", key))
		if c.health.Exist(key) {
			c.health.Del(key)
		}
	}
	return nil
}

func (c *Controller) installIPAM(ipam *v1.IPAM) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(ipam.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	// ServiceAccount IPAM
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountIPAM()); err != nil {
		if !errors.IsAlreadyExists(err) {
			// flannel service account will create automatically
			return err
		}
	}
	// ClusterRole IPAM
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(crIPAM()); err != nil {
		if !errors.IsAlreadyExists(err) {
			// flannel service account will create automatically
			return err
		}
	}
	// ClusterRoleBinding IPAM
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(crbIPAM()); err != nil {
		if !errors.IsAlreadyExists(err) {
			// flannel service account will create automatically
			return err
		}
	}
	// ConfigMap IPAM
	if _, err := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(cmIPAM()); err != nil {
		if !errors.IsAlreadyExists(err) {
			// flannel service account will create automatically
			return err
		}
	}
	// Deployment IPAM
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deploymentIPAM(ipam.Spec.Version)); err != nil {
		return err
	}
	// Service IPAM
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(serviceIPAM()); err != nil {
		return err
	}
	log.Info("ipam installed")
	return nil
}

//installation of dep, svc, rbc and so on
func serviceAccountIPAM() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcAccountIPAMName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func crbIPAM() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crbIPAMName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     svcAccountIPAMName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcAccountIPAMName,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func crIPAM() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: crIPAMName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "namespaces", "nodes", "pods/binding"},
				Verbs:     []string{"list", "watch", "get", "patch", "create"},
			},
			{
				APIGroups: []string{"apps", "extensions"},
				Resources: []string{"statefulsets", "deployments"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "endpoints", "events"},
				Verbs:     []string{"get", "list", "watch", "update", "create", "patch"},
			},
			{
				APIGroups: []string{"galaxy.k8s.io"},
				Resources: []string{"pools", "floatingips"},
				Verbs:     []string{"get", "list", "watch", "update", "create", "patch", "delete"},
			},
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"*"},
			},
		},
	}
}

func deploymentIPAM(version string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployIPAMName,
			Labels:    map[string]string{"app": "galaxy-ipam"},
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "galaxy-ipam"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "galaxy-ipam"},
				},
				Spec: corev1.PodSpec{
					Affinity: &corev1.Affinity{
						PodAntiAffinity: &corev1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
								{
									LabelSelector: &metav1.LabelSelector{
										MatchExpressions: []metav1.LabelSelectorRequirement{
											{
												Key:      "app",
												Operator: "In",
												Values:   []string{"galaxy-ipam"},
											},
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
					ServiceAccountName: svcAccountIPAMName,
					HostNetwork:        true,
					DNSPolicy:          corev1.DNSClusterFirstWithHostNet,
					Containers: []corev1.Container{
						{
							Name:  "galaxy-ipam",
							Image: images.Get(version).IPAM.FullName(),
							Args: []string{
								"--logtostderr=true",
								"--profiling",
								"--v=3",
								"--config=/etc/galaxy/galaxy-ipam.json",
								"--port=9040",
								"--api-port=9041",
								"--leader-elect",
							},
							Command: []string{
								"/usr/bin/galaxy-ipam",
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 9040, Name: "scheduler"},
								{ContainerPort: 9041, Name: "galaxy-api"},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(150, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(80*1024*1024, resource.BinarySI),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kube-config",
									MountPath: "/etc/kubernetes/",
								},
								{
									Name:      "galaxy-ipam-etc",
									MountPath: "/etc/galaxy",
								},
							},
						},
					},
					TerminationGracePeriodSeconds: int64Ptr(30),
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Effect:   corev1.TaintEffectNoSchedule,
							Operator: corev1.TolerationOpExists,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kube-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/kubernetes/",
								},
							},
						},
						{
							Name: "galaxy-ipam-etc",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: int32Ptr(420),
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "galaxy-ipam-etc",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func cmIPAM() *corev1.ConfigMap {
	const cfg = `{
      "schedule_plugin": {
        "storageDriver": "k8s-crd"
      }
    }
`
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmIPAMName,
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"galaxy-ipam.json": cfg,
		},
	}
}

func serviceIPAM() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcIPAMAPIName,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"app": "galaxy-ipam"},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "galaxy-ipam"},
			Ports: []corev1.ServicePort{
				{Name: "scheduler-port", Port: 9040, TargetPort: intstr.FromInt(9040), NodePort: 32760},
				{Name: "api-port", Port: 9041, TargetPort: intstr.FromInt(9041), NodePort: 32761},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }

func int64Ptr(i int64) *int64 { return &i }

func (c *Controller) uninstallIPAM(ipam *v1.IPAM) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(ipam.Spec.ClusterName, metav1.GetOptions{})
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
	// Service ipam
	svcIPAMErr := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(svcIPAMAPIName, &metav1.DeleteOptions{})
	// Deployment ipam
	deployIPAMErr := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(deployIPAMName, &metav1.DeleteOptions{})
	// configMap ipam
	cmIPAMErr := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete(cmIPAMName, &metav1.DeleteOptions{})
	// ClusterRoleBinding ipam
	crbIPAMErr := kubeClient.RbacV1().ClusterRoleBindings().Delete(crbIPAMName, &metav1.DeleteOptions{})
	// ClusterRole ipam
	crIPAMErr := kubeClient.RbacV1().ClusterRoles().Delete(crIPAMName, &metav1.DeleteOptions{})
	// ServiceAccount ipam
	svcAccountIPAMErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(svcAccountIPAMName, &metav1.DeleteOptions{})

	if (svcIPAMErr != nil && !errors.IsNotFound(svcIPAMErr)) ||
		(deployIPAMErr != nil && !errors.IsNotFound(deployIPAMErr)) ||
		(cmIPAMErr != nil && !errors.IsNotFound(cmIPAMErr)) ||
		(crbIPAMErr != nil && !errors.IsNotFound(crbIPAMErr)) ||
		(crIPAMErr != nil && !errors.IsNotFound(crIPAMErr)) ||
		(svcAccountIPAMErr != nil && !errors.IsNotFound(svcAccountIPAMErr)) {
		return normalerrors.New("delete ipam error")
	}
	return nil
}

func (c *Controller) watchIPAMHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check ipam in cluster health", log.String("cluster", key))
		ipam, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.client.PlatformV1().Clusters().Get(ipam.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if !c.health.Exist(cluster.Name) {
			log.Info("health check over.")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployIPAMName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("fail to get ipam deployment %v", err)
			return false, err
		}

		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", svcIPAMAPIName, string(intstr.FromInt(9040).IntVal), `/healthy`, nil).DoRaw(); err != nil {
			ipam = ipam.DeepCopy()
			ipam.Status.Phase = v1.AddonPhaseFailed
			ipam.Status.Reason = "IPAM is not healthy."
			if err = c.persistUpdate(ipam); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
}

func (c *Controller) checkIPAMStatus(ipam *v1.IPAM, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check ipam health", log.String("clusterName", ipam.Spec.ClusterName))
		cluster, err := c.client.PlatformV1().Clusters().Get(ipam.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if !c.checking.Exist(key) {
			log.Debug("checking over ipam addon status")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		ipam, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Get(deployIPAMName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("fail to create ipam %v", err)
			return false, err
		}

		ipam = ipam.DeepCopy()
		ipam.Status.Phase = v1.AddonPhaseRunning
		ipam.Status.Reason = ""
		if err = c.persistUpdate(ipam); err != nil {
			return false, err
		}
		c.checking.Del(key)
		return true, nil
	}
}

func (c *Controller) persistUpdate(ipam *v1.IPAM) error {
	var err error
	for i := 0; i < ipamMaxRetryCount; i++ {
		_, err = c.client.PlatformV1().IPAMs().UpdateStatus(ipam)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to ipam that no longer exists", log.String("clusterName", ipam.Spec.ClusterName), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to ipam '%s' that has been changed since we received it: %v", ipam.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of ipam '%s/%s'", ipam.Spec.ClusterName, ipam.Status.Phase), log.String("clusterName", ipam.Spec.ClusterName), log.Err(err))
		time.Sleep(ipamRetryInterval)
	}

	return err
}
