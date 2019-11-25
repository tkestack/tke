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

package volumedecorator

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/storage/volumedecorator/images"

	"k8s.io/apimachinery/pkg/types"

	controllerutil "tkestack.io/tke/pkg/controller"
	storageutil "tkestack.io/tke/pkg/platform/controller/addon/storage/util"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/metrics"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	controllerName = "volume-decorator-controller"

	crbName        = "volume-decorator-role-binding"
	webhookName    = "volume-decorator"
	serviceName    = "volume-decorator"
	serviceLabel   = "service"
	configMapName  = "volume-decorator-config"
	svcAccountName = "volume-decorator"
	deploymentName = "volume-decorator"

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	timeOut       = 5 * time.Minute
	maxRetryCount = 5

	tlsCert  = "tls.cert"
	tlsKey   = "tls.key"
	caCert   = "ca.cert"
	cephConf = "ceph.conf"

	configVolumeName     = "config"
	configMountPath      = "/etc/volume-manager-config"
	cephConfigSubPath    = "ceph"
	webhookConfigSubPath = "webhook"
)

// Controller is responsible for performing actions dependent upon a LogCollector phase.
type Controller struct {
	client       clientset.Interface
	cache        *volumeDecoratorCache
	health       sync.Map
	checking     sync.Map
	upgrading    sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.VolumeDecoratorLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new LogCollector Controller object.
func NewController(client clientset.Interface, informer platformv1informer.VolumeDecoratorInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client: client,
		cache:  &volumeDecoratorCache{vdMap: make(map[string]*cachedVolumeDecorator)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueVolumeDecorator,
			UpdateFunc: func(oldObj, newObj interface{}) {
				// Use the resync mechanism to find storage vendor info update.
				controller.enqueueVolumeDecorator(newObj)
			},
			DeleteFunc: controller.enqueueVolumeDecorator,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.LogCollector, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueVolumeDecorator(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for LogCollector object",
			log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
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

	err := c.syncVolumeDecorator(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing LogCollector %s (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncVolumeDecorator will sync the LogCollector with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncVolumeDecorator(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing LogCollector", log.String("name", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// decorator holds the latest LogCollector info from apiserver.
	decorator, err := c.lister.Get(name)
	switch {
	case k8serrors.IsNotFound(err):
		log.Info("LogCollector has been deleted. Attempting to cleanup resources", log.String("name", key))
		err = c.processVolumeDecoratorDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve LogCollector from store", log.String("name", key), log.Err(err))
	default:
		cachedDecorator := c.cache.getOrCreate(key)
		err = c.processDecoratorUpdate(cachedDecorator, decorator, key)
	}

	return err
}

func (c *Controller) processVolumeDecoratorDeletion(key string) error {
	cachedVolumeDecorator, ok := c.cache.get(key)
	if !ok {
		log.Error("LogCollector not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processVolumeDecoratorDelete(cachedVolumeDecorator, key)
}

func (c *Controller) processVolumeDecoratorDelete(cachedDecorator *cachedVolumeDecorator, key string) error {
	log.Info("LogCollector will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the LogCollector cache", log.String("name", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the LogCollector health cache", log.String("name", key))
		c.health.Delete(key)
	}

	decorator := cachedDecorator.state
	return c.uninstallDecorator(decorator)
}

func (c *Controller) processDecoratorUpdate(cachedDecorator *cachedVolumeDecorator, decorator *v1.VolumeDecorator, key string) error {
	if cachedDecorator.state != nil {
		// exist and the cluster name changed
		if cachedDecorator.state.UID != decorator.UID {
			if err := c.processVolumeDecoratorDelete(cachedDecorator, key); err != nil {
				return err
			}
		}
	}
	err := c.createDecoratorIfNeeded(key, cachedDecorator, decorator)
	if err != nil {
		return err
	}

	cachedDecorator.state = decorator
	// Always update the cache upon success.
	c.cache.set(key, cachedDecorator)
	return nil
}

func (c *Controller) decoratorReinitialize(
	key string,
	decorator *v1.VolumeDecorator) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		_, err := c.installDecorator(decorator)
		if err == nil {
			log.Error("Install LogCollector success",
				log.String("name", decorator.Name),
				log.String("clusterName", decorator.Spec.ClusterName))
			decorator = decorator.DeepCopy()
			decorator.Status.Phase = v1.AddonPhaseChecking
			decorator.Status.Reason = ""
			decorator.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(decorator)
			if err != nil {
				return true, err
			}
			return true, nil
		}

		// First, rollback the LogCollector.
		log.Info("Rollback LogCollector",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(err))
		if err := c.uninstallDecorator(decorator); err != nil {
			log.Error("Uninstall LogCollector failed", log.Err(err))
			return true, err
		}

		if decorator.Status.RetryCount == maxRetryCount {
			decorator = decorator.DeepCopy()
			decorator.Status.Phase = v1.AddonPhaseFailed
			decorator.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(decorator)
			if err != nil {
				log.Error("Update LogCollector failed", log.Err(err))
				return true, err
			}
			return true, nil
		}

		// Add the retry count will trigger reinitialize function from the persistent controller again.
		decorator = decorator.DeepCopy()
		decorator.Status.Phase = v1.AddonPhaseReinitializing
		decorator.Status.Reason = err.Error()
		decorator.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		decorator.Status.RetryCount++
		return true, c.persistUpdate(decorator)
	}
}

func (c *Controller) createDecoratorIfNeeded(
	key string,
	cachedDecorator *cachedVolumeDecorator,
	decorator *v1.VolumeDecorator) error {
	switch decorator.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("LogCollector will be created", log.String("name", key))
		svVersion, err := c.installDecorator(decorator)
		if err == nil {
			log.Error("Install LogCollector success",
				log.String("name", decorator.Name),
				log.String("clusterName", decorator.Spec.ClusterName))
			decorator = decorator.DeepCopy()
			fillDecoratorStatus(decorator, svVersion)
			decorator.Status.Phase = v1.AddonPhaseChecking
			decorator.Status.Reason = ""
			decorator.Status.RetryCount = 0
			return c.persistUpdate(decorator)
		}
		log.Error("Install LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(err))
		// Install LogCollector failed.
		decorator = decorator.DeepCopy()
		fillDecoratorStatus(decorator, svVersion)
		decorator.Status.Phase = v1.AddonPhaseReinitializing
		decorator.Status.Reason = err.Error()
		decorator.Status.RetryCount = 1
		decorator.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(decorator)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(decorator.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go func() {
			reInitialErr := wait.Poll(waitTime, timeOut,
				c.decoratorReinitialize(key, decorator))
			if reInitialErr != nil {
				log.Error("Reinitialize LogCollector failed",
					log.String("name", decorator.Name),
					log.String("clusterName", decorator.Spec.ClusterName),
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
					c.checkDecoratorStatus(decorator, key, initDelay))
				if checkStatusErr != nil {
					log.Error("Check status of LogCollector failed",
						log.String("name", decorator.Name),
						log.String("clusterName", decorator.Spec.ClusterName),
						log.Err(checkStatusErr))
				}
			}()
		}
	case v1.AddonPhaseRunning:
		if c.needUpgrade(decorator) {
			c.health.Delete(key)
			decorator = decorator.DeepCopy()
			decorator.Status.Phase = v1.AddonPhaseUpgrading
			decorator.Status.Reason = ""
			decorator.Status.RetryCount = 0
			return c.persistUpdate(decorator)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go func() {
				healthErr := wait.PollImmediateUntil(5*time.Minute,
					c.watchDecoratorHealth(key), c.stopCh)
				if healthErr != nil {
					log.Error("Watch health of LogCollector failed",
						log.String("name", decorator.Name),
						log.String("clusterName", decorator.Spec.ClusterName),
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
					c.upgradeVolumeDecorator(decorator, key, upgradeDelay))
				if upgradeErr != nil {
					log.Error("Upgrade LogCollector failed",
						log.String("name", decorator.Name),
						log.String("clusterName", decorator.Spec.ClusterName),
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

func (c *Controller) needUpgrade(decorator *v1.VolumeDecorator) bool {
	sort.Strings(decorator.Spec.VolumeTypes)
	sort.Strings(decorator.Status.VolumeTypes)
	if decorator.Spec.Version != decorator.Status.Version ||
		decorator.Spec.WorkloadAdmission != decorator.Status.WorkloadAdmission ||
		!reflect.DeepEqual(decorator.Spec.VolumeTypes, decorator.Status.VolumeTypes) {
		return true
	}

	version, err := storageutil.GetSVInfoVersion(c.client)
	if err != nil {
		log.Errorf("Get ceph info failed: %v", err)
		return true
	}

	// Update the configMap and deployment if the storage vendor info changed.
	return version != decorator.Status.StorageVendorVersion
}

func (c *Controller) installDecorator(decorator *v1.VolumeDecorator) (string, error) {
	cluster, err := c.client.PlatformV1().Clusters().Get(decorator.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return "", err
	}

	svInfo, err := storageutil.GetSVInfo(c.client)
	if err != nil {
		return "", err
	}

	version := ""
	if svInfo != nil {
		version = svInfo.Version
	}

	// Create ServiceAccount.
	if err := c.installSVCAccount(decorator, kubeClient); err != nil {
		return version, err
	}

	// Create ClusterRoleBinding.
	if err := c.installCRB(decorator, kubeClient); err != nil {
		return version, err
	}

	// Create ConfigMap.
	if err := c.installConfigMap(decorator, kubeClient, svInfo); err != nil {
		return version, err
	}

	// Create Service.
	if err := c.installSVC(decorator, kubeClient); err != nil {
		return version, err
	}

	// Create Deployment. The decorator will create the webhook after started.
	return version, c.installDeployment(decorator, kubeClient, svInfo)
}

func (c *Controller) installSVCAccount(decorator *v1.VolumeDecorator, kubeClient kubernetes.Interface) error {
	account := genServiceAccount()
	accountClient := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem)

	_, err := accountClient.Get(account.Name, metav1.GetOptions{})
	if err == nil {
		log.Info("ServiceAccount of LogCollector is already created",
			log.String("name", decorator.Name))
		return nil
	}

	if k8serrors.IsNotFound(err) {
		_, err = accountClient.Create(account)
		return err
	}

	return fmt.Errorf("get account failed: %v", err)
}

func (c *Controller) installCRB(decorator *v1.VolumeDecorator, kubeClient kubernetes.Interface) error {
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
			log.String("name", decorator.Name))
		return nil
	}

	newCRB := oldCRB.DeepCopy()
	newCRB.RoleRef = crb.RoleRef
	newCRB.Subjects = crb.Subjects
	_, err = crbClient.Update(newCRB)

	return err
}

func (c *Controller) installConfigMap(
	decorator *v1.VolumeDecorator,
	kubeClient kubernetes.Interface,
	svInfo *storageutil.SVInfo) error {
	cm, err := genConfigMap(decorator, svInfo)
	if err != nil {
		return fmt.Errorf("generate config map failed: %v", err)
	}
	cmClient := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem)

	oldCM, err := cmClient.Get(cm.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err := cmClient.Create(cm)
			return err
		}
		return fmt.Errorf("get configMap failed: %v", err)
	}

	if equality.Semantic.DeepEqual(cm.Data, oldCM.Data) {
		log.Info("ConfigMap of LogCollector is already created",
			log.String("name", decorator.Name))
		return nil
	}

	newSVC := oldCM.DeepCopy()
	newSVC.Data = cm.Data
	_, err = cmClient.Update(newSVC)

	return err
}

func (c *Controller) installSVC(decorator *v1.VolumeDecorator, kubeClient kubernetes.Interface) error {
	if !decorator.Spec.WorkloadAdmission {
		log.Info("Workload admission disabled, skip installing service",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.ClusterName))
		return nil
	}

	svc := genService()
	svcClient := kubeClient.CoreV1().Services(metav1.NamespaceSystem)

	oldSVC, err := svcClient.Get(svc.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err := svcClient.Create(svc)
			return err
		}
		return fmt.Errorf("get account failed: %v", err)
	}

	if equality.Semantic.DeepEqual(svc.Spec, oldSVC.Spec) {
		log.Info("Service of LogCollector is already created",
			log.String("name", decorator.Name))
		return nil
	}

	newSVC := oldSVC.DeepCopy()
	newSVC.Spec = svc.Spec
	_, err = svcClient.Update(newSVC)

	return err
}

func (c *Controller) installDeployment(
	decorator *v1.VolumeDecorator,
	kubeClient kubernetes.Interface,
	svInfo *storageutil.SVInfo) error {
	deploy, err := c.genDeployment(decorator, svInfo)
	if err != nil {
		return err
	}
	deployClient := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem)

	oldDeploy, err := deployClient.Get(deploy.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = deployClient.Create(deploy)
			return err
		}
		return fmt.Errorf("get deployment failed: %v", err)
	}

	newDeploy := oldDeploy.DeepCopy()
	newDeploy.Labels = deploy.Labels
	newDeploy.Spec = deploy.Spec
	_, err = deployClient.Update(newDeploy)

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

func genConfigMap(decorator *v1.VolumeDecorator, svInfo *storageutil.SVInfo) (*corev1.ConfigMap, error) {
	certContext, err := storageutil.SetupServerCert(fmt.Sprintf("%s.%s.svc", serviceName, metav1.NamespaceSystem), "admin")
	if err != nil {
		return nil, fmt.Errorf("generate cert failed: %v", err)
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: metav1.NamespaceSystem,
		},
		Data: make(map[string]string),
	}

	if decorator.Spec.WorkloadAdmission {
		cm.Data[tlsCert] = string(certContext.Cert)
		cm.Data[tlsKey] = string(certContext.Key)
		cm.Data[caCert] = string(certContext.SigningCert)
	}

	if svInfo != nil {
		cm.Data[cephConf] = fmt.Sprintf(storageutil.CephConfTemplate, svInfo.Monitors)
		adminKeyring := fmt.Sprintf(storageutil.CephKeyringFileNameTemplate, svInfo.AdminID)
		cm.Data[adminKeyring] = fmt.Sprintf(storageutil.CephAdminKeyringTemplate, svInfo.AdminKey)
	}

	return cm, nil
}

func genService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{serviceLabel: deploymentName},
			Ports: []corev1.ServicePort{
				{Protocol: corev1.ProtocolTCP, Port: 443, TargetPort: intstr.FromInt(443)},
			},
		},
	}
}

func (c *Controller) genDeployment(decorator *v1.VolumeDecorator, svInfo *storageutil.SVInfo) (*appsv1.Deployment, error) {
	labels := map[string]string{"app": controllerName, serviceLabel: deploymentName}
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Labels:    labels,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": controllerName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					PriorityClassName:  "system-cluster-critical",
					ServiceAccountName: svcAccountName,
					HostNetwork:        true,
					HostPID:            true,
				},
			},
		},
	}

	container, volumes, err := c.genContainerAndVolumes(decorator, svInfo)
	if err != nil {
		return nil, err
	}

	deploy.Spec.Template.Spec.Volumes = volumes
	deploy.Spec.Template.Spec.Containers = []corev1.Container{*container}

	return deploy, nil
}

func int32Ptr(i int32) *int32 { return &i }

func boolPtr(value bool) *bool { return &value }

func (c *Controller) uninstallDecorator(decorator *v1.VolumeDecorator) error {
	log.Info("Start to uninstall LogCollector",
		log.String("name", decorator.Name),
		log.String("clusterName", decorator.Spec.ClusterName))

	kubeClient, err := c.getKubeClient(decorator.Spec.ClusterName)
	if err != nil {
		return err
	}
	if kubeClient == nil {
		// The cluster is not found.
		return nil
	}

	// Delete the webhook.
	clearWebhookErr := kubeClient.AdmissionregistrationV1beta1().
		ValidatingWebhookConfigurations().Delete(webhookName, &metav1.DeleteOptions{})
	// Delete the service.
	clearSvcErr := kubeClient.CoreV1().
		Services(metav1.NamespaceSystem).Delete(serviceName, &metav1.DeleteOptions{})
	// Delete the decorator deployment.
	clearDeployErr := kubeClient.AppsV1().
		Deployments(metav1.NamespaceSystem).Delete(deploymentName, &metav1.DeleteOptions{})
	// Delete configMap.
	clearCMErr := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).
		Delete(configMapName, &metav1.DeleteOptions{})
	// Delete the ClusterRoleBinding.
	clearCRBErr := kubeClient.RbacV1().
		ClusterRoleBindings().Delete(crbName, &metav1.DeleteOptions{})
	// Delete the ServiceAccount.
	clearSvcAccountErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).
		Delete(svcAccountName, &metav1.DeleteOptions{})

	failed := false

	if clearWebhookErr != nil && !k8serrors.IsNotFound(clearWebhookErr) {
		failed = true
		log.Error("delete webhook for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearWebhookErr))
	}

	if clearSvcErr != nil && !k8serrors.IsNotFound(clearSvcErr) {
		failed = true
		log.Error("delete service for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearSvcErr))
	}

	if clearDeployErr != nil && !k8serrors.IsNotFound(clearDeployErr) {
		failed = true
		log.Error("delete deployment for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearDeployErr))
	}

	if clearCMErr != nil && !k8serrors.IsNotFound(clearCMErr) {
		failed = true
		log.Error("delete config map for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearCMErr))
	}

	if clearCRBErr != nil && !k8serrors.IsNotFound(clearCRBErr) {
		failed = true
		log.Error("delete crb for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearCRBErr))
	}

	if clearSvcAccountErr != nil && !k8serrors.IsNotFound(clearSvcAccountErr) {
		failed = true
		log.Error("delete service account for LogCollector failed",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.Err(clearSvcAccountErr))
	}

	if failed {
		return errors.New("delete LogCollector failed")
	}

	return nil
}

func (c *Controller) watchDecoratorHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		decorator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		log.Info("Start check health of LogCollector", log.String("name", decorator.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(decorator.Spec.ClusterName, metav1.GetOptions{})
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

		_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).
			Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			decorator = decorator.DeepCopy()
			decorator.Status.Phase = v1.AddonPhaseFailed
			decorator.Status.Reason = "LogCollector is not healthy."
			if err = c.persistUpdate(decorator); err != nil {
				return false, err
			}
			return true, nil
		}

		return false, nil
	}
}

func (c *Controller) checkDecoratorStatus(
	decorator *v1.VolumeDecorator,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check LogCollector health", log.String("name", decorator.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(decorator.Spec.ClusterName, metav1.GetOptions{})
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
		decorator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		deploy, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).
			Get(deploymentName, metav1.GetOptions{})
		if err != nil ||
			(deploy.Spec.Replicas != nil && deploy.Status.AvailableReplicas < *deploy.Spec.Replicas) {
			if time.Now().After(initDelay) {
				decorator = decorator.DeepCopy()
				decorator.Status.Phase = v1.AddonPhaseFailed
				decorator.Status.Reason = "Volume Decorator is not healthy."
				if err = c.persistUpdate(decorator); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		decorator = decorator.DeepCopy()
		decorator.Status.Phase = v1.AddonPhaseRunning
		decorator.Status.Reason = ""
		if err = c.persistUpdate(decorator); err != nil {
			return false, err
		}

		return true, nil
	}
}

func (c *Controller) upgradeVolumeDecorator(
	decorator *v1.VolumeDecorator,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade LogCollector", log.String("name", decorator.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(decorator.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}

		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading LogCollector", log.String("name", decorator.Name))
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		decorator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		timeoutCheck := func(operation string, err error) (bool, error) {
			log.Error(operation+" of LogCollector failed",
				log.String("name", decorator.Name),
				log.String("clusterName", decorator.Spec.ClusterName),
				log.Err(err))
			if k8serrors.IsNotFound(err) {
				return false, err
			}
			if time.Now().After(initDelay) {
				decorator = decorator.DeepCopy()
				decorator.Status.Phase = v1.AddonPhaseFailed
				decorator.Status.Reason = "Upgrade LogCollector failed."
				if err = c.persistUpdate(decorator); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		svInfo, err := storageutil.GetSVInfo(c.client)
		if err != nil {
			return timeoutCheck("get storage vendor info", err)
		}

		if !decorator.Spec.WorkloadAdmission {
			// Delete the webhook.
			err := kubeClient.AdmissionregistrationV1beta1().
				ValidatingWebhookConfigurations().Delete(webhookName, &metav1.DeleteOptions{})
			if err != nil && !k8serrors.IsNotFound(err) {
				return timeoutCheck("delete webhook", err)
			}
		}

		// Upgrade configMap.
		if err := c.installConfigMap(decorator, kubeClient, svInfo); err != nil {
			return timeoutCheck("update configMap", err)
		}

		deployClient := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem)
		oldDeploy, err := deployClient.Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			return timeoutCheck("get old deployment", err)
		}
		oldContainer := &oldDeploy.Spec.Template.Spec.Containers[0]

		container, volumes, err := c.genContainerAndVolumes(decorator, svInfo)
		if err != nil {
			return false, err
		}
		// Copy the default fields.
		container.TerminationMessagePath = oldContainer.TerminationMessagePath
		container.TerminationMessagePolicy = oldContainer.TerminationMessagePolicy

		newDeploy := oldDeploy.DeepCopy()
		newDeploy.Spec.Template.Spec.Volumes = volumes
		newDeploy.Spec.Template.Spec.Containers = []corev1.Container{*container}
		patchData, err := storageutil.GetPatchData(oldDeploy, newDeploy)
		if err != nil {
			log.Error("get deployment patch data of LogCollector failed",
				log.String("name", decorator.Name),
				log.String("clusterName", decorator.Spec.ClusterName),
				log.Err(err))
			return false, err
		}

		if len(patchData) > 0 {
			log.Info("Upgrade deployment of LogCollector",
				log.String("name", decorator.Name),
				log.String("clusterName", decorator.Spec.ClusterName),
				log.String("patchData", string(patchData)))

			_, err = deployClient.Patch(deploymentName, types.StrategicMergePatchType, patchData)
			if err != nil {
				return timeoutCheck("patch deployment", err)
			}
		}

		version := ""
		if svInfo != nil {
			version = svInfo.Version
		}

		decorator = decorator.DeepCopy()
		fillDecoratorStatus(decorator, version)
		decorator.Status.Phase = v1.AddonPhaseChecking
		decorator.Status.Reason = ""
		if err = c.persistUpdate(decorator); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) persistUpdate(decorator *v1.VolumeDecorator) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.PlatformV1().VolumeDecorators().UpdateStatus(decorator)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if k8serrors.IsNotFound(err) {
			log.Info("Not persisting update to LogCollector that no longer exists",
				log.String("clusterName", decorator.Spec.ClusterName), log.Err(err))
			return nil
		}
		if k8serrors.IsConflict(err) {
			return fmt.Errorf("not persisting update to LogCollector %q that has been changed since we received it: %v", decorator.Spec.ClusterName, err)
		}
		log.Warn("Failed to persist updated status of LogCollector",
			log.String("name", decorator.Name),
			log.String("clusterName", decorator.Spec.ClusterName),
			log.String("phase", string(decorator.Status.Phase)), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}

func (c *Controller) getKubeClient(clusterName string) (kubernetes.Interface, error) {
	cluster, err := c.client.PlatformV1().Clusters().Get(clusterName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return util.BuildExternalClientSet(cluster, c.client.PlatformV1())
}

func (c *Controller) getImage(decorator *v1.VolumeDecorator) (string, error) {
	return images.Get(decorator.Spec.Version).VolumeDecorator.FullName(), nil
}

func (c *Controller) genContainerAndVolumes(
	decorator *v1.VolumeDecorator,
	svInfo *storageutil.SVInfo) (*corev1.Container, []corev1.Volume, error) {
	image, err := c.getImage(decorator)
	if err != nil {
		return nil, nil, err
	}
	return genContainer(image, decorator, svInfo), genVolumes(decorator, svInfo), nil
}

func genContainer(
	image string,
	decorator *v1.VolumeDecorator,
	svInfo *storageutil.SVInfo) *corev1.Container {
	procMount := corev1.DefaultProcMount
	container := &corev1.Container{
		Name:  deploymentName,
		Image: image,
		SecurityContext: &corev1.SecurityContext{
			Privileged: boolPtr(true),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"SYS_ADMIN"},
			},
			AllowPrivilegeEscalation: boolPtr(true),
			ProcMount:                &procMount,
		},
		Args: []string{
			"--leader-election=true",
			"--create-crd=true",
			"--logtostderr=true",
			"--v=5",
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				// TODO: add support for configuring them
				corev1.ResourceCPU:    *resource.NewMilliQuantity(200, resource.DecimalSI),
				corev1.ResourceMemory: *resource.NewScaledQuantity(256, resource.Mega),
			},
		},
		Lifecycle: &corev1.Lifecycle{
			PreStop: &corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"/bin/sh", "-c", "umount /tmp/cephfs-root"},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: configVolumeName, MountPath: configMountPath},
			{Name: "host-dev", MountPath: "/dev"},
			{Name: "host-sys", MountPath: "/sys"},
			{Name: "lib-modules", MountPath: "/lib/modules", ReadOnly: true},
		},
	}
	if svInfo != nil {
		cephAdminKeyring := fmt.Sprintf(storageutil.CephKeyringFileNameTemplate, svInfo.AdminID)
		container.Args = append(container.Args,
			"--ceph-config-file="+filepath.Join(configMountPath, cephConfigSubPath, cephConf),
			"--ceph-keyring-file="+filepath.Join(configMountPath, cephConfigSubPath, cephAdminKeyring))
	}

	if decorator.Spec.WorkloadAdmission {
		container.Args = append(container.Args,
			"--workload-admission=true",
			"--webhook-name="+webhookName,
			"--service-name="+serviceName,
			"--client-ca-file="+filepath.Join(configMountPath, webhookConfigSubPath, caCert),
			"--tls-cert-file="+filepath.Join(configMountPath, webhookConfigSubPath, tlsCert),
			"--tls-private-key-file="+filepath.Join(configMountPath, webhookConfigSubPath, tlsKey))
	}

	if len(decorator.Spec.VolumeTypes) > 0 {
		container.Args = append(container.Args,
			"--volume-types="+strings.Join(decorator.Spec.VolumeTypes, ","))
	}

	return container
}

func genVolumes(
	decorator *v1.VolumeDecorator,
	svInfo *storageutil.SVInfo) []corev1.Volume {
	volume := corev1.Volume{
		Name: configVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: configMapName},
				DefaultMode:          int32Ptr(0644),
			},
		},
	}

	if svInfo != nil {
		cephAdminKeyring := fmt.Sprintf(storageutil.CephKeyringFileNameTemplate, svInfo.AdminID)

		volume.ConfigMap.Items = append(volume.ConfigMap.Items,
			corev1.KeyToPath{Key: cephConf, Path: filepath.Join(cephConfigSubPath, cephConf)},
			corev1.KeyToPath{Key: cephAdminKeyring, Path: filepath.Join(cephConfigSubPath, cephAdminKeyring)},
		)
	}

	if decorator.Spec.WorkloadAdmission {
		volume.ConfigMap.Items = append(volume.ConfigMap.Items,
			corev1.KeyToPath{Key: caCert, Path: filepath.Join(webhookConfigSubPath, caCert)},
			corev1.KeyToPath{Key: tlsCert, Path: filepath.Join(webhookConfigSubPath, tlsCert)},
			corev1.KeyToPath{Key: tlsKey, Path: filepath.Join(webhookConfigSubPath, tlsKey)},
		)
	}

	hostPathType := corev1.HostPathDirectory

	return []corev1.Volume{volume,
		{Name: "host-dev", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/dev", Type: &hostPathType}}},
		{Name: "host-sys", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/sys", Type: &hostPathType}}},
		{Name: "lib-modules", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/lib/modules", Type: &hostPathType}}},
	}
}

func fillDecoratorStatus(decorator *v1.VolumeDecorator, svVersion string) {
	decorator.Status.Version = decorator.Spec.Version
	decorator.Status.WorkloadAdmission = decorator.Spec.WorkloadAdmission
	decorator.Status.VolumeTypes = decorator.Spec.VolumeTypes
	decorator.Status.StorageVendorVersion = svVersion
}
