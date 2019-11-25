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

package lbcf

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/controller/addon/lbcf/images"

	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/controller"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/metrics"

	"k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1informer "tkestack.io/tke/api/client/informers/externalversions/platform/v1"
	platformv1lister "tkestack.io/tke/api/client/listers/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	lbcfClientRetryCount    = 5
	lbcfClientRetryInterval = 5 * time.Second

	lbcfMaxRetryCount = 5
	lbcfTimeOut       = 5 * time.Minute
)

const (
	svcLBCFHealthCheckName = "lbcf-controller"
	svcLBCFHealthCheckPort = 11029
	svcLBCFHealthCheckPath = "/healthz"

	lbcfAPIGroup   = "lbcf.tke.cloud.tencent.com"
	lbcfAPIVersion = "v1beta1"

	crdLoadBalancerDriver = "loadbalancerdrivers"
	crdLoadBalancer       = "loadbalancers"
	crdBackendGroup       = "backendgroups"

	validatingWebhookName = "lbcf-validate"
	mutatingWebhookName   = "lbcf-mutate"
	caBundle              = `-----BEGIN CERTIFICATE-----
MIIDNDCCAhwCCQCH+2EEbqe/OTANBgkqhkiG9w0BAQsFADBcMQswCQYDVQQGEwJD
TjELMAkGA1UECAwCQkoxFjAUBgNVBAoMDXRlbmNlbnQsIEluYy4xKDAmBgNVBAMM
H2xiY2YtY29udHJvbGxlci5rdWJlLXN5c3RlbS5zdmMwHhcNMTkwNTE1MDYwMTQ5
WhcNMjIwMzA0MDYwMTQ5WjBcMQswCQYDVQQGEwJDTjELMAkGA1UECAwCQkoxFjAU
BgNVBAoMDXRlbmNlbnQsIEluYy4xKDAmBgNVBAMMH2xiY2YtY29udHJvbGxlci5r
dWJlLXN5c3RlbS5zdmMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDn
rhdUjDrFCfZTR7BLM8siq3ZH1khcbJjF2r1ikh5Kk8DZM4gulPHXrfCYm1OQB0ow
9ynI3REq0ckkQP3HfgrMaXxKTKcak4vPGtiQ9XVH/0Ga888am7PAPobIsKxSsX9Q
0N/FvRmYvR+kYECpKeU5hN7IuAFegrB8wwx0co5R7O9rIYSC/TqiK+bmmRh4ApyF
W6AioU1IZclP6XBU1nDkEUObNKMGCl8laEt4w7x/nVPxyAXeBi6cibM3uqDO0u22
0VCQsIF0iMIVYi5yTx53B1ccKLNyIZatf8xoFcKtrj7QHJPmahOruHn93bUx37nd
mboDlLjrViz8VcF8NIp9AgMBAAEwDQYJKoZIhvcNAQELBQADggEBABmrA6Cr+CW2
WqxvW45EpLvZprcyUlcFLaAtj4B+tBEBzgLoafZVTwFe+nS9hBE10QIBdXU6qdOV
+6LOUbm6hSKDoXWQ8rkyedFOBchYI3d8T9mJzOu6S9hPBbMQvBqHOG9n+RyM9E65
1xD0yV0g4oiz4AAniawTvaRVk5cmxysfXKBAQl2O8AKN+vVtAGpZbrXVCsGsLY7r
iuDxj60aNuR66FN7+Yw21eYP1awcnRFFDy/m+VPOUWJAsyPoH0Gd0apYYLpi4383
U9GSCkdsMs1M8xK3FaoD+a2ERoDwP9kdDi27sM6mumlM9KbZ7wVilLUsIIN5T61o
DSwXwCfjM58=
-----END CERTIFICATE-----`
)

var (
	validateDriverPath       = "/validate-load-balancer-driver"
	validateLBPath           = "/validate-load-balancer"
	validateBackendGroupPath = "/validate-backend-group"

	mutateLBPath           = "/mutate-load-balancer"
	mutateDriverPath       = "/mutate-load-balancer-driver"
	mutateBackendGroupPath = "/mutate-backend-broup"
)

// Controller is responsible for performing actions dependent upon a LBCF phase.
type Controller struct {
	client       clientset.Interface
	cache        *lbcfCache
	health       sync.Map
	checking     sync.Map
	queue        workqueue.RateLimitingInterface
	lister       platformv1lister.LBCFLister
	listerSynced cache.InformerSynced
	stopCh       <-chan struct{}
}

// NewController creates a new Controller object
func NewController(client clientset.Interface, lbcfInformer platformv1informer.LBCFInformer, resyncPeriod time.Duration) *Controller {
	// create the controller so we can inject the enqueue function
	c := &Controller{
		client: client,
		cache:  &lbcfCache{lbcfMap: make(map[string]*cachedLBCF)},
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "lbcf"),
	}

	if client != nil && client.PlatformV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("lbcf_controller", client.PlatformV1().RESTClient().GetRateLimiter())
	}

	// configure the lbcf informer event handlers
	lbcfInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: c.enqueueLBCF,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldLBCF, ok1 := oldObj.(*v1.LBCF)
				curLBCF, ok2 := newObj.(*v1.LBCF)
				if ok1 && ok2 && c.needsUpdate(oldLBCF, curLBCF) {
					c.enqueueLBCF(newObj)
				}
			},
			DeleteFunc: c.enqueueLBCF,
		},
		resyncPeriod,
	)
	c.lister = lbcfInformer.Lister()
	c.listerSynced = lbcfInformer.Informer().HasSynced

	return c
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting LBCF controller")
	defer log.Info("Shutting down LBCF controller")

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

func (c *Controller) needsUpdate(old *v1.LBCF, cur *v1.LBCF) bool {
	return !reflect.DeepEqual(old, cur)
}

func (c *Controller) enqueueLBCF(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

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

	err := c.syncLBCF(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing lbcf %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncLBCF will sync the LBCF with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncLBCF(key string) error {
	startTime := time.Now()
	var cachedLBCF *cachedLBCF
	defer func() {
		log.Info("Finished syncing LBCF", log.String("LBCFName", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// lbcf holds the latest LBCF info from apiserver
	lbcf, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("LBCF has been deleted. Attempting to cleanup resources", log.String("LBCF", key))
		err = c.processLBCFDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve LBCF from store", log.String("LBCF", key), log.Err(err))
	default:
		cachedLBCF = c.cache.getOrCreate(key)
		err = c.processLBCFUpdate(cachedLBCF, lbcf, key)
	}
	return err
}

func (c *Controller) processLBCFDeletion(key string) error {
	cachedLBCF, ok := c.cache.get(key)
	if !ok {
		log.Error("LBCF not in cache even though the watcher thought it was. Ignoring the deletion", log.String("lbcfName", key))
		return nil
	}
	return c.processLBCFDelete(cachedLBCF, key)
}

func (c *Controller) processLBCFDelete(cachedLBCF *cachedLBCF, key string) error {
	log.Info("LBCF will be dropped", log.String("LBCFName", key))

	if c.cache.Exist(key) {
		log.Info("delete the LBCF cache", log.String("LBCFName", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("delete the LBCF health cache", log.String("LBCFName", key))
		c.health.Delete(key)
	}

	lbcf := cachedLBCF.state
	return c.uninstallLBCF(lbcf)
}

func (c *Controller) processLBCFUpdate(cachedLBCF *cachedLBCF, lbcf *v1.LBCF, key string) error {
	if cachedLBCF.state != nil {
		// exist and the cluster name changed
		if cachedLBCF.state.UID != lbcf.UID {
			if err := c.processLBCFDelete(cachedLBCF, key); err != nil {
				return err
			}
		}
	}
	err := c.createLBCFIfNeeded(key, cachedLBCF, lbcf)
	if err != nil {
		return err
	}

	cachedLBCF.state = lbcf
	// Always update the cache upon success.
	c.cache.set(key, cachedLBCF)
	return nil
}

func (c *Controller) lbcfReinitialize(key string, cachedLBCF *cachedLBCF, lbcf *v1.LBCF) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		err := c.installLBCF(lbcf)
		if err == nil {
			lbcf = lbcf.DeepCopy()
			lbcf.Status.Phase = v1.AddonPhaseChecking
			lbcf.Status.Reason = ""
			lbcf.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(lbcf)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		// First, rollback the lbcf
		if err := c.uninstallLBCF(lbcf); err != nil {
			log.Errorf("Uninstall lbcf error: %v", err)
			return true, err
		}
		if lbcf.Status.RetryCount == lbcfMaxRetryCount {
			lbcf = lbcf.DeepCopy()
			lbcf.Status.Phase = v1.AddonPhaseFailed
			lbcf.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", lbcfMaxRetryCount)
			err := c.persistUpdate(lbcf)
			if err != nil {
				log.Errorf("Update lbcf error: %v", err)
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		lbcf = lbcf.DeepCopy()
		lbcf.Status.Phase = v1.AddonPhaseReinitializing
		lbcf.Status.Reason = err.Error()
		lbcf.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		lbcf.Status.RetryCount++
		err = c.persistUpdate(lbcf)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createLBCFIfNeeded(key string, cachedLBCF *cachedLBCF, lbcf *v1.LBCF) error {
	switch lbcf.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Error("LBCF will be created", log.String("lbcf", key))
		err := c.installLBCF(lbcf)
		if err == nil {
			lbcf = lbcf.DeepCopy()
			lbcf.Status.Phase = v1.AddonPhaseChecking
			lbcf.Status.Reason = ""
			lbcf.Status.RetryCount = 0
			return c.persistUpdate(lbcf)
		}
		log.Errorf("installLBCF err: %v", err)
		lbcf = lbcf.DeepCopy()
		lbcf.Status.Phase = v1.AddonPhaseReinitializing
		lbcf.Status.Reason = err.Error()
		lbcf.Status.RetryCount = 1
		lbcf.Status.LastReInitializingTimestamp = metav1.Now()
		return c.persistUpdate(lbcf)
	case v1.AddonPhaseReinitializing:
		var interval = time.Since(lbcf.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= lbcfTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = lbcfTimeOut - interval
		}
		go wait.Poll(waitTime, lbcfTimeOut, c.lbcfReinitialize(key, cachedLBCF, lbcf))
	case v1.AddonPhaseChecking:
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, true)
			initDelay := time.Now().Add(5 * time.Minute)
			go wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkLBCFStatus(lbcf, key, initDelay))
		}
	case v1.AddonPhaseRunning:
		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, true)
			go wait.PollImmediateUntil(5*time.Minute, c.watchLBCFHealth(key), c.stopCh)
		}
	case v1.AddonPhaseFailed:
		log.Info("LBCF is error", log.String("LBCF ", key))
		if _, ok := c.health.Load(key); ok {
			c.health.Delete(key)
		}
	}
	return nil
}

func (c *Controller) installLBCF(lbcf *v1.LBCF) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(lbcf.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	crdClient, err := util.BuildExternalExtensionClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}

	if _, err := kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(validatingWebhook()); err != nil {
		return err
	}
	if _, err := kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Create(mutatingWebhook()); err != nil {
		return err
	}
	for _, crd := range crds() {
		if _, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil {
			return err
		}
	}
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccount()); err != nil {
		return err
	}
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(clusterRole()); err != nil {
		return err
	}
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(clusterRoleBinding()); err != nil {
		return err
	}
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(secret()); err != nil {
		return err
	}
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deployment(lbcf.Spec.Version)); err != nil {
		return err
	}
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(service()); err != nil {
		return err
	}
	return nil
}

var (
	driverRes       = schema.GroupVersionResource{Group: lbcfAPIGroup, Version: "v1beta1", Resource: "loadbalancerdrivers"}
	lbRes           = schema.GroupVersionResource{Group: lbcfAPIGroup, Version: "v1beta1", Resource: "loadbalancers"}
	backendRes      = schema.GroupVersionResource{Group: lbcfAPIGroup, Version: "v1beta1", Resource: "backendrecords"}
	backendGroupRes = schema.GroupVersionResource{Group: lbcfAPIGroup, Version: "v1beta1", Resource: "backendgroups"}
)

func (c *Controller) uninstallLBCF(lbcf *v1.LBCF) error {
	cluster, err := c.client.PlatformV1().Clusters().Get(lbcf.Spec.ClusterName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	crdClient, err := util.BuildExternalExtensionClientSet(cluster, c.client.PlatformV1())
	if err != nil {
		return err
	}
	credential, err := util.ClusterCredentialV1(c.client.PlatformV1(), cluster.Name)
	if err != nil {
		return err
	}

	dynamicClient, err := util.BuildExternalDynamicClientSet(cluster, credential)
	if err != nil {
		return err
	}

	if err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Delete(validatingWebhookName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Delete(mutatingWebhookName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.RbacV1().ClusterRoles().Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.RbacV1().ClusterRoleBindings().Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	if err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(svcLBCFHealthCheckName, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	if err := removeFinalizers(dynamicClient, driverRes); err != nil {
		return err
	}
	if err := removeFinalizers(dynamicClient, lbRes); err != nil {
		return err
	}
	if err := removeFinalizers(dynamicClient, backendGroupRes); err != nil {
		return err
	}
	if err := removeFinalizers(dynamicClient, backendRes); err != nil {
		return err
	}
	return deleteCRDWithTimeout(crdClient)
}

func (c *Controller) watchLBCFHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check LBCF in cluster health", log.String("cluster", key))
		lbcf, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.client.PlatformV1().Clusters().Get(lbcf.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.health.Load(cluster.Name); !ok {
			log.Info("health check over.")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", svcLBCFHealthCheckName, strconv.Itoa(svcLBCFHealthCheckPort), svcLBCFHealthCheckPath, nil).DoRaw(); err != nil {
			lbcf = lbcf.DeepCopy()
			lbcf.Status.Phase = v1.AddonPhaseFailed
			lbcf.Status.Reason = "LBCF is not healthy."
			if err = c.persistUpdate(lbcf); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
}

func (c *Controller) checkLBCFStatus(lbcf *v1.LBCF, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check LBCF health", log.String("clusterName", lbcf.Spec.ClusterName))
		cluster, err := c.client.PlatformV1().Clusters().Get(lbcf.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Infof("checkLBCFStatus: lbcf not found")
			return false, err
		}
		if err != nil {
			log.Errorf("checkLBCFStatus err: %v", err)
			return false, nil
		}
		if _, ok := c.checking.Load(key); !ok {
			log.Infof("checking over LBCF addon status")
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.client.PlatformV1())
		if err != nil {
			return false, err
		}
		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", svcLBCFHealthCheckName, strconv.Itoa(svcLBCFHealthCheckPort), svcLBCFHealthCheckPath, nil).DoRaw(); err != nil {
			if time.Now().After(initDelay) {
				lbcf = lbcf.DeepCopy()
				lbcf.Status.Phase = v1.AddonPhaseFailed
				lbcf.Status.Reason = "LBCF is not healthy."
				if err = c.persistUpdate(lbcf); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
		lbcf = lbcf.DeepCopy()
		lbcf.Status.Phase = v1.AddonPhaseRunning
		lbcf.Status.Reason = ""
		if err = c.persistUpdate(lbcf); err != nil {
			return false, err
		}
		c.checking.Delete(key)
		return true, nil
	}
}

func (c *Controller) persistUpdate(lbcf *v1.LBCF) error {
	var err error
	for i := 0; i < lbcfClientRetryCount; i++ {
		_, err = c.client.PlatformV1().LBCFs().UpdateStatus(lbcf)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to LBDF that no longer exists", log.String("clusterName", lbcf.Spec.ClusterName), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to LBCF '%s' that has been changed since we received it: %v", lbcf.Spec.ClusterName, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of LBCF '%s/%s'", lbcf.Spec.ClusterName, lbcf.Status.Phase), log.String("clusterName", lbcf.Spec.ClusterName), log.Err(err))
		time.Sleep(lbcfClientRetryInterval)
	}
	return err
}

func validatingWebhook() *v1beta1.ValidatingWebhookConfiguration {
	failPolicy := v1beta1.Fail
	return &v1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: validatingWebhookName,
		},
		Webhooks: []v1beta1.ValidatingWebhook{
			{
				Name: "driver.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdLoadBalancerDriver},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
							v1beta1.Update,
							v1beta1.Delete,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &validateDriverPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
			{
				Name: "lb.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdLoadBalancer},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
							v1beta1.Update,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &validateLBPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
			{
				Name: "backendgroup.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdBackendGroup},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
							v1beta1.Update,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &validateBackendGroupPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
		},
	}
}

func mutatingWebhook() *v1beta1.MutatingWebhookConfiguration {
	failPolicy := v1beta1.Fail
	return &v1beta1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: mutatingWebhookName,
		},
		Webhooks: []v1beta1.MutatingWebhook{
			{
				Name: "lb.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdLoadBalancer},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &mutateLBPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
			{
				Name: "driver.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdLoadBalancerDriver},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &mutateDriverPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
			{
				Name: "backendgroup.lbcf.tke.cloud.tencent.com",
				Rules: []v1beta1.RuleWithOperations{
					{
						Rule: v1beta1.Rule{
							APIGroups:   []string{lbcfAPIGroup},
							APIVersions: []string{lbcfAPIVersion},
							Resources:   []string{crdBackendGroup},
						},
						Operations: []v1beta1.OperationType{
							v1beta1.Create,
							v1beta1.Update,
						},
					},
				},
				ClientConfig: v1beta1.WebhookClientConfig{
					CABundle: []byte(caBundle),
					Service: &v1beta1.ServiceReference{
						Name:      svcLBCFHealthCheckName,
						Namespace: metav1.NamespaceSystem,
						Path:      &mutateBackendGroupPath,
					},
				},
				FailurePolicy: &failPolicy,
			},
		},
	}
}

func crds() []*extensionsv1.CustomResourceDefinition {
	return []*extensionsv1.CustomResourceDefinition{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "loadbalancerdrivers.lbcf.tke.cloud.tencent.com",
			},
			Spec: extensionsv1.CustomResourceDefinitionSpec{
				Group: lbcfAPIGroup,
				Names: extensionsv1.CustomResourceDefinitionNames{
					Kind:     "LoadBalancerDriver",
					ListKind: "LoadBalancerDriverList",
					Plural:   "loadbalancerdrivers",
					Singular: "loadbalancerdriver",
				},
				Scope: extensionsv1.NamespaceScoped,
				Subresources: &extensionsv1.CustomResourceSubresources{
					Status: &extensionsv1.CustomResourceSubresourceStatus{},
				},
				Version: "v1beta1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "loadbalancers.lbcf.tke.cloud.tencent.com",
			},
			Spec: extensionsv1.CustomResourceDefinitionSpec{
				Group: lbcfAPIGroup,
				Names: extensionsv1.CustomResourceDefinitionNames{
					Kind:     "LoadBalancer",
					ListKind: "LoadBalancerList",
					Plural:   "loadbalancers",
					Singular: "loadbalancer",
				},
				Scope: extensionsv1.NamespaceScoped,
				Subresources: &extensionsv1.CustomResourceSubresources{
					Status: &extensionsv1.CustomResourceSubresourceStatus{},
				},
				Version: "v1beta1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backendgroups.lbcf.tke.cloud.tencent.com",
			},
			Spec: extensionsv1.CustomResourceDefinitionSpec{
				Group: lbcfAPIGroup,
				Names: extensionsv1.CustomResourceDefinitionNames{
					Kind:     "BackendGroup",
					ListKind: "BackendGroupList",
					Plural:   "backendgroups",
					Singular: "backendgroup",
				},
				Scope: extensionsv1.NamespaceScoped,
				Subresources: &extensionsv1.CustomResourceSubresources{
					Status: &extensionsv1.CustomResourceSubresourceStatus{},
				},
				Version: "v1beta1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backendrecords.lbcf.tke.cloud.tencent.com",
			},
			Spec: extensionsv1.CustomResourceDefinitionSpec{
				Group: lbcfAPIGroup,
				Names: extensionsv1.CustomResourceDefinitionNames{
					Kind:     "BackendRecord",
					ListKind: "BackendRecordList",
					Plural:   "backendrecords",
					Singular: "backendrecord",
				},
				Scope: extensionsv1.NamespaceScoped,
				Subresources: &extensionsv1.CustomResourceSubresources{
					Status: &extensionsv1.CustomResourceSubresourceStatus{},
				},
				Version: "v1beta1",
			},
		},
	}
}

var selectorForLBCF = metav1.LabelSelector{
	MatchLabels: map[string]string{"lbcf.tke.cloud.tencent.com/component": svcLBCFHealthCheckName},
}

func deployment(version string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcLBCFHealthCheckName,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForLBCF,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"lbcf.tke.cloud.tencent.com/component": svcLBCFHealthCheckName},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: svcLBCFHealthCheckName,
					PriorityClassName:  "system-node-critical",
					Containers: []corev1.Container{
						{
							Name:  "controller",
							Image: images.Get(version).LBCFController.FullName(),
							Ports: []corev1.ContainerPort{
								{ContainerPort: 443, Name: "admit-server"},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "server-tls",
									MountPath: "/etc/lbcf",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "server-tls",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: svcLBCFHealthCheckName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func serviceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcLBCFHealthCheckName,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func clusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcLBCFHealthCheckName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"pods",
					"services",
					"events",
					"nodes",
				},
				Verbs: []string{rbacv1.VerbAll},
			},
			{
				APIGroups: []string{lbcfAPIGroup},
				Resources: []string{rbacv1.ResourceAll},
				Verbs:     []string{rbacv1.VerbAll},
			},
		},
	}
}

func clusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcLBCFHealthCheckName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      svcLBCFHealthCheckName,
				Namespace: metav1.NamespaceSystem,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     svcLBCFHealthCheckName,
		},
	}
}

func secret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcLBCFHealthCheckName,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"server.key": []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAmIJebdDH7bkHcnA1EplNTS0BQbKbRK51wRdsulV8qfzq1Wl1
maC9XBpZ4rhQUw8hWpt8N5ZM4cel2WSTQpNNIRxcfJ+gIH2C7dYuat+STgI2ICVj
yeMg/LANRZEd3yDFNJJNkYKwNroTgYEGZPnErblQ3K8vY76SFm3NJ9uQpuMGq2/m
A974DGCXHju8DidmpbvUq4uGT2S1+1oeDIb7AorciZta4lX1S1GWR9WSHY4tqesq
3NfVUkOkftU1RjORpHW1RI5G5jB7aQ3uMmhe2I5puD3csdaez5Nqbh2ilfT6NFTw
j1Qc8Ox0nqJwkzEc2ICMCWYEs/R4//lZbfMn9wIDAQABAoIBAHKlPjsrQcAQ4epD
M4Jhv9yOQm2SyGnfBCI9a7y/WtGmkRoRBxiP3wmHvZ5Tk/58V0R3se9Pi0gG/0Pm
+VSIyuhjG5uLm6IQ+AW2hnpMyvzdaLbNpLA1j6yk47UyG9SKG/UjLjB+n9zkEJm/
1oC9yf4WWxUqlGNU9RjrPdgClED/LAaCAZT0hfKaRnLqCNKve1L241fbSnXe+zWF
JvrxVWF2KVy7kGpnFHO1hRrWQ0NPpTXn2A3JDmaCM4YnSFe45t2IZyN8Y+PLnL9E
tA9f6YfJ86Far/MqTApHHgPLc4fKT4QS4PdRqYOyUvamYF+dwgGx8e0rtZZBeJP0
aja8WeECgYEAxn2DRV9pH1s1gVJvgyOW6Ve6e+y/uPgGxcik6ZwLm5EpxhSAzUp8
HnD5aOobAGR7iEDEUjio13l2qrUoIYGiUWH3p8cZG+MW5/IVzmog06loYbJfLmE2
4MMkDfPKzGAHvGyIhHow+bRBjP1D/qy6YLPgbf3jy1LcOXLQ4ISDGr8CgYEAxLJV
90+4BYe2Ed81YtNjLYdxoFX2iOEiN40Bxk8UHvvsxRbj7H+LnCl5O1JM3WS2MPGR
xooXq3JOXOutJbs6eFXzeYrIjQCrrqhQ4/kFTqN++VIw8REvvvB5FlAZ842EHSAC
MBsmuqjy1Auy27+gSU7fGSeXXnscuR2+ojIO2MkCgYAZrJB3P7EcQjL4iE4uO0NA
6X0QnH3sEgDmQl66bNm/hJZPrcU/SJwnX9uS630UnuqvpBkAvZ1xSZ/E0uve8aKq
Pi7Hf+RKjCQhWlnhui6G0knTITxYhnCPwA4A1ADuUJmPkMZTxG5jTiKQdw39eiAd
dAbak1WMrioYMDa+Y8WFhwKBgQDASzBr1Q2sql4+3p5MfSgqbI2TGDcq3h4bfMjN
XKXpHJT+oUA2BwMvqgQREIaAsmLDOocvN/Wn8NnXUbg2ePHSjwS2QA2Me6lb2MUr
+llL5d7OU6HxKsIowuM+AxU724/bAV3iNckJFv4+eyliV9aVlHvbFa+P+H++Iewq
mRGWsQKBgFlukoJUVtjAsPCIYioztqyRkEehuSjOgAL5wjpWA4Xo1WbMSxI+ElqT
2yaFwxTHkeM51yF4mbbYrstZGh1opzNZ7tuCF4iDwV0gXE5KbV66vTB2AvNo5ncJ
ZAuaT9x1UI+tLu8WZSsy3uIggljZiDihFASdIPNbAaeb+4ZmIWnD
-----END RSA PRIVATE KEY-----
`),
			"server.crt": []byte(`-----BEGIN CERTIFICATE-----
MIIDNDCCAhwCCQDW7a4S+vONGjANBgkqhkiG9w0BAQsFADBcMQswCQYDVQQGEwJD
TjELMAkGA1UECAwCQkoxFjAUBgNVBAoMDXRlbmNlbnQsIEluYy4xKDAmBgNVBAMM
H2xiY2YtY29udHJvbGxlci5rdWJlLXN5c3RlbS5zdmMwHhcNMTkwNTE1MDYwMTUz
WhcNMjAwOTI2MDYwMTUzWjBcMQswCQYDVQQGEwJDTjELMAkGA1UECAwCQkoxFjAU
BgNVBAoMDXRlbmNlbnQsIEluYy4xKDAmBgNVBAMMH2xiY2YtY29udHJvbGxlci5r
dWJlLXN5c3RlbS5zdmMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCY
gl5t0MftuQdycDUSmU1NLQFBsptErnXBF2y6VXyp/OrVaXWZoL1cGlniuFBTDyFa
m3w3lkzhx6XZZJNCk00hHFx8n6AgfYLt1i5q35JOAjYgJWPJ4yD8sA1FkR3fIMU0
kk2RgrA2uhOBgQZk+cStuVDcry9jvpIWbc0n25Cm4warb+YD3vgMYJceO7wOJ2al
u9Sri4ZPZLX7Wh4MhvsCityJm1riVfVLUZZH1ZIdji2p6yrc19VSQ6R+1TVGM5Gk
dbVEjkbmMHtpDe4yaF7Yjmm4Pdyx1p7Pk2puHaKV9Po0VPCPVBzw7HSeonCTMRzY
gIwJZgSz9Hj/+Vlt8yf3AgMBAAEwDQYJKoZIhvcNAQELBQADggEBABjXjVI0RzhC
bPQw1fPMYzq+trilMuZZe0EhBS9pNE/OEf54ECgYSz0U0i/K9VqTJrsd5D9wwVV1
slMejeRPZsojPwcyM+6th2hMyijWF2TaLlvsB/I1rFzo19YRYbicqp9qby61OL++
9fAfxhGAPSTOPaL6DWDZM2DCSrA0WXQNLFiMLjxCsKxfCYPWM51VY60JCxvOIuWF
1jaecz6ymopms0glE0e0XxcPk56yBjH/8aTxtyGjmund7+t6qv7keGPAVHYdWeMw
H32tNdB5RHM/w9U8suobL2ZUGUjP4W9aspvMl7k28EYZgnqeO1xuBq9dZZ7qYCd4
bojVejw/us8=
-----END CERTIFICATE-----
`),
		},
	}
}

func service() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcLBCFHealthCheckName,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"lbcf.tke.cloud.tencent.com/component": svcLBCFHealthCheckName,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "admit-server",
					Port:       443,
					TargetPort: intstr.IntOrString{IntVal: 443},
				},
				{
					Name:       "healthz",
					Port:       11029,
					TargetPort: intstr.IntOrString{IntVal: 11029},
				},
			},
		},
	}
}

func removeFinalizers(dynamicClient dynamic.Interface, resource schema.GroupVersionResource) error {
	list, err := dynamicClient.Resource(resource).Namespace(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, obj := range list.Items {
		if len(obj.GetFinalizers()) > 0 {
			cpy := obj.DeepCopy()
			cpy.SetFinalizers(nil)
			if _, err := dynamicClient.Resource(resource).Namespace(cpy.GetNamespace()).Update(cpy, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteCRDWithTimeout(crdClient *apiextensionsclient.Clientset) error {
	ch := make(chan error)
	defer close(ch)

	go func() {
		for _, crd := range crds() {
			log.Infof("Start delete LBCF crd %s", crd.Name)
			if err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("delete LBCF crd failed: %v", err)
				if !errors.IsNotFound(err) {
					ch <- err
					return
				}
			}
			log.Infof("Finish delete LBCF crd %s", crd.Name)
		}
		ch <- nil
	}()

	select {
	case err := <-ch:
		return err
	case <-time.After(20 * time.Second):
		return fmt.Errorf("delete CRD timeout")
	}
}
