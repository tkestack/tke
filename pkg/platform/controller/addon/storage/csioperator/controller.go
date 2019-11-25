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

package csioperator

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/storage/csioperator/images"

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
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	controllerutil "tkestack.io/tke/pkg/controller"
	storageutil "tkestack.io/tke/pkg/platform/controller/addon/storage/util"
	"tkestack.io/tke/pkg/platform/util"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	controllerName = "csi-operator-controller"

	crbName        = "csi-operator-role-binding"
	svcAccountName = "csi-operator"
	deploymentName = "csi-operator"

	clientRetryCount    = 5
	clientRetryInterval = 5 * time.Second

	timeOut       = 5 * time.Minute
	maxRetryCount = 5
)

// Controller is responsible for performing actions dependent upon a CSIOperator phase.
type Controller struct {
	client       clientset.Interface
	cache        *csiOperatorCache
	health       sync.Map
	checking     sync.Map
	upgrading    sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.CSIOperatorLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new CSIOperator Controller object.
func NewController(client clientset.Interface, informer platformv1informer.CSIOperatorInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{

		client: client,
		cache:  &csiOperatorCache{coMap: make(map[string]*cachedCSIOperator)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage(controllerName, client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the informer event handlers
	informer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueCSIOperator,
			UpdateFunc: func(oldObj, newObj interface{}) {
				// Use the resync mechanism to find storage vendor info update.
				controller.enqueueCSIOperator(newObj)
			},
			DeleteFunc: controller.enqueueCSIOperator,
		},
		resyncPeriod,
	)
	controller.lister = informer.Lister()
	controller.listerSynced = informer.Informer().HasSynced

	return controller
}

// obj could be an *v1.CSIOperator, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueCSIOperator(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for CSIOperator object",
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
	log.Info("Starting CSIOperator controller")
	defer log.Info("Shutting down CSIOperator controller")

	if !cache.WaitForCacheSync(stopCh, c.listerSynced) {
		return fmt.Errorf("failed to wait for CSIOperator cache to sync")
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

	err := c.syncCSIOperator(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing CSIOperator %s (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncCSIOperator will sync the CSIOperator with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncCSIOperator(key string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finished syncing CSIOperator", log.String("name", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// csiOperator holds the latest CSIOperator info from apiserver.
	csiOperator, err := c.lister.Get(name)
	switch {
	case k8serrors.IsNotFound(err):
		log.Info("CSIOperator has been deleted. Attempting to cleanup resources", log.String("name", key))
		err = c.processCSIOperatorDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve CSIOperator from store", log.String("name", key), log.Err(err))
	default:
		cachedCSIOperator := c.cache.getOrCreate(key)
		err = c.processCSIOperatorUpdate(cachedCSIOperator, csiOperator, key)
	}

	return err
}

func (c *Controller) processCSIOperatorDeletion(key string) error {
	cachedCSIOperator, ok := c.cache.get(key)
	if !ok {
		log.Error("CSIOperator not in cache even though the watcher thought it was. Ignoring the deletion", log.String("name", key))
		return nil
	}
	return c.processCSIOperatorDelete(cachedCSIOperator, key)
}

func (c *Controller) processCSIOperatorDelete(cachedCSIOperator *cachedCSIOperator, key string) error {
	log.Info("CSIOperator will be dropped", log.String("name", key))

	if c.cache.Exist(key) {
		log.Info("Delete the CSIOperator cache", log.String("name", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the CSIOperator health cache", log.String("name", key))
		c.health.Delete(key)
	}

	csiOperator := cachedCSIOperator.state
	return c.uninstallCSIOperator(csiOperator)
}

func (c *Controller) processCSIOperatorUpdate(cachedCSIOperator *cachedCSIOperator, csiOperator *v1.CSIOperator, key string) error {
	if cachedCSIOperator.state != nil {
		// exist and the cluster name changed
		if cachedCSIOperator.state.UID != csiOperator.UID {
			if err := c.processCSIOperatorDelete(cachedCSIOperator, key); err != nil {
				return err
			}
		}
	}
	err := c.createCSIOperatorIfNeeded(key, cachedCSIOperator, csiOperator)
	if err != nil {
		return err
	}

	cachedCSIOperator.state = csiOperator
	// Always update the cache upon success.
	c.cache.set(key, cachedCSIOperator)
	return nil
}

func (c *Controller) csiOperatorReinitialize(
	key string,
	cachedCSIOperator *cachedCSIOperator,
	csiOperator *v1.CSIOperator) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		_, err := c.installCSIOperator(csiOperator)
		if err == nil {
			csiOperator = csiOperator.DeepCopy()
			csiOperator.Status.Phase = v1.AddonPhaseChecking
			csiOperator.Status.Reason = ""
			csiOperator.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(csiOperator)
			if err != nil {
				return true, err
			}
			return true, nil
		}

		// First, rollback the CSIOperator.
		log.Info("Rollback CSIOperator",
			log.String("name", csiOperator.Name),
			log.String("clusterName", csiOperator.Spec.ClusterName))
		if err := c.uninstallCSIOperator(csiOperator); err != nil {
			log.Error("Uninstall CSIOperator failed", log.Err(err))
			return true, err
		}

		if csiOperator.Status.RetryCount == maxRetryCount {
			csiOperator = csiOperator.DeepCopy()
			csiOperator.Status.Phase = v1.AddonPhaseFailed
			csiOperator.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", maxRetryCount)
			err := c.persistUpdate(csiOperator)
			if err != nil {
				log.Error("Update CSIOperator failed", log.Err(err))
				return true, err
			}
			return true, nil
		}

		// Add the retry count will trigger reinitialize function from the persistent controller again.
		csiOperator = csiOperator.DeepCopy()
		csiOperator.Status.Phase = v1.AddonPhaseReinitializing
		csiOperator.Status.Reason = err.Error()
		csiOperator.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		csiOperator.Status.RetryCount++
		return true, c.persistUpdate(csiOperator)
	}
}

func (c *Controller) createCSIOperatorIfNeeded(
	key string,
	cachedCSIOperator *cachedCSIOperator,
	csiOperator *v1.CSIOperator) error {
	switch csiOperator.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info("CSIOperator will be created", log.String("name", key))
		svVersion, err := c.installCSIOperator(csiOperator)
		if err == nil {
			csiOperator = csiOperator.DeepCopy()
			fillOperatorStatus(csiOperator, svVersion)
			csiOperator.Status.Phase = v1.AddonPhaseChecking
			csiOperator.Status.Reason = ""
			csiOperator.Status.RetryCount = 0
			return c.persistUpdate(csiOperator)
		}
		// Install CSIOperator failed.
		csiOperator = csiOperator.DeepCopy()
		fillOperatorStatus(csiOperator, svVersion)
		csiOperator.Status.Phase = v1.AddonPhaseReinitializing
		csiOperator.Status.Reason = err.Error()
		csiOperator.Status.RetryCount = 1
		csiOperator.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(csiOperator)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(csiOperator.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= timeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = timeOut - interval
		}
		go func() {
			reInitialErr := wait.Poll(waitTime, timeOut,
				c.csiOperatorReinitialize(key, cachedCSIOperator, csiOperator))
			if reInitialErr != nil {
				log.Error("Reinitialize CSIOperator failed",
					log.String("name", csiOperator.Name),
					log.String("clusterName", csiOperator.Spec.ClusterName),
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
					c.checkCSIOperatorStatus(csiOperator, key, initDelay))
				if checkStatusErr != nil {
					log.Error("Check status of CSIOperator failed",
						log.String("name", csiOperator.Name),
						log.String("clusterName", csiOperator.Spec.ClusterName),
						log.Err(checkStatusErr))
				}
			}()
		}
	case v1.AddonPhaseRunning:
		if c.needUpgrade(csiOperator) {
			c.health.Delete(key)
			csiOperator = csiOperator.DeepCopy()
			csiOperator.Status.Phase = v1.AddonPhaseUpgrading
			csiOperator.Status.Reason = ""
			csiOperator.Status.RetryCount = 0
			return c.persistUpdate(csiOperator)
		}
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go func() {
				healthErr := wait.PollImmediateUntil(5*time.Minute,
					c.watchCSIOperatorHealth(key), c.stopCh)
				if healthErr != nil {
					log.Error("Watch health of CSIOperator failed",
						log.String("name", csiOperator.Name),
						log.String("clusterName", csiOperator.Spec.ClusterName),
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
					c.upgradeCSIOperator(csiOperator, key, upgradeDelay))
				if upgradeErr != nil {
					log.Error("Upgrade CSIOperator failed",
						log.String("name", csiOperator.Name),
						log.String("clusterName", csiOperator.Spec.ClusterName),
						log.Err(upgradeErr))
				}
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("CSIOperator failed", log.String("name", key))
		c.health.Delete(key)
		c.checking.Delete(key)
		c.upgrading.Delete(key)
	}
	return nil
}

func (c *Controller) needUpgrade(csiOperator *v1.CSIOperator) bool {
	if csiOperator.Spec.Version != csiOperator.Status.Version {
		return true
	}

	svVersion, err := storageutil.GetSVInfoVersion(c.client)
	return err != nil || svVersion != csiOperator.Status.StorageVendorVersion
}

func (c *Controller) installCSIOperator(csiOperator *v1.CSIOperator) (string, error) {
	cluster, err := c.client.PlatformV1().Clusters().Get(csiOperator.Spec.ClusterName, metav1.GetOptions{})
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
	if err := c.installSVC(csiOperator, kubeClient); err != nil {
		return version, err
	}

	// Create ClusterRoleBinding.
	if err := c.installCRB(csiOperator, kubeClient); err != nil {
		return version, err
	}

	// Create Deployment.
	return version, c.installDeployment(csiOperator, kubeClient, svInfo)
}

func (c *Controller) installSVC(csiOperator *v1.CSIOperator, kubeClient kubernetes.Interface) error {
	svc := genServiceAccount()
	svcClient := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem)

	_, err := svcClient.Get(svc.Name, metav1.GetOptions{})
	if err == nil {
		log.Info("ServiceAccount of CSIOperator is already created",
			log.String("name", csiOperator.Name))
		return nil
	}

	if k8serrors.IsNotFound(err) {
		_, err = svcClient.Create(svc)
		return err
	}

	return fmt.Errorf("get svc failed: %v", err)
}

func (c *Controller) installCRB(csiOperator *v1.CSIOperator, kubeClient kubernetes.Interface) error {
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
		log.Info("ClusterRoleBinding of CSIOperator is already created",
			log.String("name", csiOperator.Name))
		return nil
	}

	newCRB := oldCRB.DeepCopy()
	newCRB.RoleRef = crb.RoleRef
	newCRB.Subjects = crb.Subjects
	_, err = crbClient.Update(newCRB)

	return err
}

func (c *Controller) installDeployment(
	csiOperator *v1.CSIOperator,
	kubeClient kubernetes.Interface,
	svInfo *storageutil.SVInfo) error {
	deploy := c.genDeployment(images.Get(csiOperator.Spec.Version), svInfo)
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

func (c *Controller) genDeployment(components images.Components, svInfo *storageutil.SVInfo) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
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
					ServiceAccountName: svcAccountName,
					Containers: []corev1.Container{
						{
							Name:  deploymentName,
							Image: components.CSIOperator.FullName(),
							Args:  genContainerArgs(svInfo),
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									// TODO: add support for configuring them
									corev1.ResourceCPU:    *resource.NewMilliQuantity(200, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewScaledQuantity(256, resource.Mega),
								},
							},
						},
					},
				},
			},
		},
	}

	return deploy
}

func genContainerArgs(svInfo *storageutil.SVInfo) []string {
	args := []string{
		"--leader-election=true",
		"--kubelet-root-dir=/var/lib/kubelet",
		"--registry-domain=" + containerregistryutil.GetPrefix(),
		"--logtostderr=true",
		"--v=5",
	}
	if svInfo != nil {
		if svInfo.Type == storageutil.Ceph {
			args = append(args,
				"--ceph-monitors="+svInfo.Monitors,
				"--ceph-admin-id="+svInfo.AdminID,
				"--ceph-admin-key="+svInfo.AdminKey)
		} else {
			args = append(args,
				"--tencent-cloud-secret-id="+svInfo.SecretID,
				"--tencent-cloud-secret-key="+svInfo.SecretKey)
		}
	}
	return args
}

func int32Ptr(i int32) *int32 { return &i }

func (c *Controller) uninstallCSIOperator(csiOperator *v1.CSIOperator) error {
	log.Info("Start to uninstall CSIOperator",
		log.String("name", csiOperator.Name),
		log.String("clusterName", csiOperator.Spec.ClusterName))

	cluster, err := c.client.PlatformV1().Clusters().Get(csiOperator.Spec.ClusterName, metav1.GetOptions{})
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
	// Delete the operator deployment.
	clearDeployErr := kubeClient.AppsV1().
		Deployments(metav1.NamespaceSystem).Delete(deploymentName, &metav1.DeleteOptions{})
	// Delete the ClusterRoleBinding.
	clearCRBErr := kubeClient.RbacV1().
		ClusterRoleBindings().Delete(crbName, &metav1.DeleteOptions{})
	// Delete the ServiceAccount.
	clearSVCErr := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).
		Delete(svcAccountName, &metav1.DeleteOptions{})

	failed := false

	if clearDeployErr != nil && !k8serrors.IsNotFound(clearDeployErr) {
		failed = true
		log.Error("delete deployment for CSIOperator failed",
			log.String("name", csiOperator.Name),
			log.String("clusterName", csiOperator.Spec.ClusterName),
			log.Err(clearDeployErr))
	}

	if clearCRBErr != nil && !k8serrors.IsNotFound(clearCRBErr) {
		failed = true
		log.Error("delete crb for CSIOperator failed",
			log.String("name", csiOperator.Name),
			log.String("clusterName", csiOperator.Spec.ClusterName),
			log.Err(clearCRBErr))
	}

	if clearSVCErr != nil && !k8serrors.IsNotFound(clearSVCErr) {
		failed = true
		log.Error("delete service account for CSIOperator failed",
			log.String("name", csiOperator.Name),
			log.String("clusterName", csiOperator.Spec.ClusterName),
			log.Err(clearSVCErr))
	}

	if failed {
		return errors.New("delete CSIOperator failed")
	}

	return nil
}

func (c *Controller) watchCSIOperatorHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		csiOperator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}
		log.Info("Start check health of CSIOperator", log.String("name", csiOperator.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(csiOperator.Spec.ClusterName, metav1.GetOptions{})
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
			csiOperator = csiOperator.DeepCopy()
			csiOperator.Status.Phase = v1.AddonPhaseFailed
			csiOperator.Status.Reason = "CSIOperator is not healthy."
			if err = c.persistUpdate(csiOperator); err != nil {
				return false, err
			}
			return true, nil
		}

		return false, nil
	}
}

func (c *Controller) checkCSIOperatorStatus(
	csiOperator *v1.CSIOperator,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check CSIOperator health", log.String("name", csiOperator.Name))

		cluster, err := c.client.PlatformV1().Clusters().Get(csiOperator.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}

		if _, ok := c.checking.Load(key); !ok {
			log.Debug("Checking over CSIOperator addon status")
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		csiOperator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		deploy, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).
			Get(deploymentName, metav1.GetOptions{})
		if err != nil ||
			(deploy.Spec.Replicas != nil && deploy.Status.AvailableReplicas < *deploy.Spec.Replicas) {
			if time.Now().After(initDelay) {
				csiOperator = csiOperator.DeepCopy()
				csiOperator.Status.Phase = v1.AddonPhaseFailed
				csiOperator.Status.Reason = "CSI Operator is not healthy."
				if err = c.persistUpdate(csiOperator); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}

		csiOperator = csiOperator.DeepCopy()
		csiOperator.Status.Phase = v1.AddonPhaseRunning
		csiOperator.Status.Reason = ""
		if err = c.persistUpdate(csiOperator); err != nil {
			return false, err
		}

		return true, nil
	}
}

func (c *Controller) upgradeCSIOperator(
	csiOperator *v1.CSIOperator,
	key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade CSIOperator", log.String("name", csiOperator.Name))
		cluster, err := c.client.PlatformV1().Clusters().Get(csiOperator.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}

		if _, ok := c.upgrading.Load(key); !ok {
			log.Debug("Upgrading CSIOperator", log.String("name", csiOperator.Name))
			return true, nil
		}

		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		csiOperator, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		timeoutCheck := func(operation string, err error) (bool, error) {
			log.Error(operation+" of CSIOperator failed",
				log.String("name", csiOperator.Name),
				log.String("clusterName", csiOperator.Spec.ClusterName),
				log.Err(err))
			if time.Now().After(initDelay) {
				csiOperator = csiOperator.DeepCopy()
				csiOperator.Status.Phase = v1.AddonPhaseFailed
				csiOperator.Status.Reason = "Upgrade CSIOperator failed."
				if err = c.persistUpdate(csiOperator); err != nil {
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

		deployClient := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem)
		oldDeploy, err := deployClient.Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			return timeoutCheck("get deployment", err)
		}

		newDeploy := oldDeploy.DeepCopy()
		newDeploy.Spec.Template.Spec.Containers[0].Args = genContainerArgs(svInfo)
		patchData, err := storageutil.GetPatchData(oldDeploy, newDeploy)
		if err != nil {
			log.Error("get deployment patch data of CSIOperator failed",
				log.String("name", csiOperator.Name),
				log.String("clusterName", csiOperator.Spec.ClusterName),
				log.Err(err))
			return false, err
		}

		_, err = deployClient.Patch(deploymentName, types.StrategicMergePatchType, patchData)
		if err != nil {
			return timeoutCheck("update deployment", err)
		}

		version := ""
		if svInfo != nil {
			version = svInfo.Version
		}

		csiOperator = csiOperator.DeepCopy()
		fillOperatorStatus(csiOperator, version)
		csiOperator.Status.Phase = v1.AddonPhaseChecking
		csiOperator.Status.Reason = ""
		if err = c.persistUpdate(csiOperator); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) persistUpdate(csiOperator *v1.CSIOperator) error {
	var err error
	for i := 0; i < clientRetryCount; i++ {
		_, err = c.client.PlatformV1().CSIOperators().UpdateStatus(csiOperator)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if k8serrors.IsNotFound(err) {
			log.Info("Not persisting update to CSIOperator that no longer exists",
				log.String("clusterName", csiOperator.Spec.ClusterName), log.Err(err))
			return nil
		}
		if k8serrors.IsConflict(err) {
			return fmt.Errorf("not persisting update to CSIOperator %q that has been changed since we received it: %v", csiOperator.Spec.ClusterName, err)
		}
		log.Warn("Failed to persist updated status of CSIOperator",
			log.String("name", csiOperator.Name),
			log.String("clusterName", csiOperator.Spec.ClusterName),
			log.String("phase", string(csiOperator.Status.Phase)), log.Err(err))
		time.Sleep(clientRetryInterval)
	}

	return err
}

func fillOperatorStatus(csiOperator *v1.CSIOperator, svVersion string) {
	csiOperator.Status.Version = csiOperator.Spec.Version
	csiOperator.Status.StorageVendorVersion = svVersion
}
