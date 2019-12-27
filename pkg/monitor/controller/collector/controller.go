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

package collector

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/coreos/prometheus-operator/pkg/apis/monitoring"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promopk8sutil "github.com/coreos/prometheus-operator/pkg/k8sutil"
	influxapi "github.com/influxdata/influxdb1-client/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	monitorclientset "tkestack.io/tke/api/client/clientset/versioned/typed/monitor/v1"
	platformclientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	monitorv1informer "tkestack.io/tke/api/client/informers/externalversions/monitor/v1"
	monitorv1lister "tkestack.io/tke/api/client/listers/monitor/v1"
	v1 "tkestack.io/tke/api/monitor/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	notifyapi "tkestack.io/tke/cmd/tke-notify-api/app"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/monitor/config"
	"tkestack.io/tke/pkg/monitor/storage"
	monitorutil "tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/platform/util"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	collectorClientRetryCount    = 5
	collectorClientRetryInterval = 5 * time.Second

	collectorMaxRetryCount = 5
	collectorTimeOut       = 5 * time.Minute
)

const (
	// PrometheusService is the service name for prometheus app
	PrometheusOperatorService            = "prometheus-operator"
	PrometheusOperatorServicePort        = "http"
	prometheusOperatorServiceAccount     = "prometheus-operator"
	prometheusOperatorClusterRoleBinding = "prometheus-operator"
	prometheusOperatorClusterRole        = "prometheus-operator"
	prometheusOperatorWorkLoad           = "prometheus-operator"
	PrometheusBeatService                = "prometheus-beat"
	PrometheusBeatServicePort            = "http"
	PrometheusBeatConfigmap              = "prometheus-beat-config"
	PrometheusBeatConfigFile             = "prometheusbeat.yml"
	prometheusBeatWorkLoad               = "prometheus-beat"
	// PrometheusService is the service name for prometheus app
	PrometheusService = "prometheus"
	// PrometheusServicePort is the port name for prometheus service
	PrometheusServicePort        = "http"
	PrometheusCRDName            = "k8s"
	prometheusServiceAccount     = PrometheusService + "-" + PrometheusCRDName
	prometheusClusterRoleBinding = PrometheusService + "-" + PrometheusCRDName
	prometheusClusterRole        = PrometheusService + "-" + PrometheusCRDName
	prometheusSecret             = PrometheusService + "-" + PrometheusCRDName + "-" + "additional-scrape-config"
	prometheusWorkLoad           = "prometheus"
	prometheusETCDSecret         = "prometheus-etcd"
	prometheusRuleRecord         = "prometheus-records"
	PrometheusRuleAlert          = "prometheus-alerts"
	prometheusConfigName         = "prometheus.config.yaml"
	prometheusImagePath          = "prometheus"

	// AlertManagerService defines the service for alert manager app
	AlertManagerService = "alertmanager"
	// AlertManagerWorkLoad defines the app name for alert manager
	AlertManagerWorkLoad = "alertmanager"
	// AlertManagerConfigMap defines the configmap name which stores the alertmanager config rules
	AlertManagerConfigMap = "alertmanager-config"
	// AlertManagerConfigName defines the entry name of the configmap which saves the alertmanager rules
	AlertManagerConfigName     = "alertmanager.yml"
	alertManagerImagePath      = "alertmanager"
	alertManagerServicePort    = "http"
	alertManagerCRDName        = "main"
	alertManagerServiceAccount = AlertManagerService + "-" + alertManagerCRDName
	alertManagerSecret         = AlertManagerService + "-" + alertManagerCRDName

	nodeExporterService   = "node-exporter"
	nodeExporterDaemonSet = "node-exporter"

	kubeStateService            = "kube-state-metrics"
	kubeStateServiceAccount     = "kube-state-metrics"
	kubeStateClusterRoleBinding = "kube-state-metrics"
	kubeStateClusterRole        = "kube-state-metrics"
	kubeStateRoleBinding        = "kube-state-metrics"
	kubeStateRole               = "kube-state-metrics-resizer"
	kubeStateWorkLoad           = "kube-state-metrics"

	specialLabelName  = "k8s-submitter"
	specialLabelValue = "controller"
)

var crdKinds = []monitoringv1.CrdKind{
	monitoringv1.DefaultCrdKinds.Alertmanager,
	monitoringv1.DefaultCrdKinds.Prometheus,
	monitoringv1.DefaultCrdKinds.PodMonitor,
	monitoringv1.DefaultCrdKinds.PrometheusRule,
	monitoringv1.DefaultCrdKinds.ServiceMonitor,
}

// Controller is responsible for performing actions dependent upon a collector phase.
type Controller struct {
	monitorClient  monitorclientset.MonitorV1Interface
	platformClient platformclientset.PlatformV1Interface
	cache          *collectorCache
	health         sync.Map
	checking       sync.Map
	upgrading      sync.Map
	queue          workqueue.RateLimitingInterface
	lister         monitorv1lister.CollectorLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// RemoteClient for collector
	remoteClient *storage.RemoteClient
	// NotifyApiAddress
	notifyAPIAddress string
}

// NewController creates a new Controller object.
func NewController(monitorClient monitorclientset.MonitorV1Interface, platformClient platformclientset.PlatformV1Interface, collectorInformer monitorv1informer.CollectorInformer, resyncPeriod time.Duration, remoteClient *storage.RemoteClient) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		platformClient: platformClient,
		monitorClient:  monitorClient,
		cache:          &collectorCache{collectorMap: make(map[string]*cachedCollector)},
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "collector"),
		remoteClient:   remoteClient,
	}

	if monitorClient != nil && monitorClient.RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("collector_controller", monitorClient.RESTClient().GetRateLimiter())
	}

	// configure the prometheus informer event handlers
	collectorInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueCollector,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldCollector, ok1 := oldObj.(*v1.Collector)
				curCollector, ok2 := newObj.(*v1.Collector)
				if ok1 && ok2 && controller.needsUpdate(oldCollector, curCollector) {
					controller.enqueueCollector(newObj)
				}
			},
			DeleteFunc: controller.enqueueCollector,
		},
		resyncPeriod,
	)
	controller.lister = collectorInformer.Lister()
	controller.listerSynced = collectorInformer.Informer().HasSynced

	return controller
}

// obj could be an *v1.Collector, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueueCollector(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldCollector *v1.Collector, newCollector *v1.Collector) bool {
	return !reflect.DeepEqual(oldCollector, newCollector)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	// Start the informer factories to begin populating the informer caches
	log.Info("Starting collector controller")
	defer log.Info("Shutting down collector controller")

	if len(c.remoteClient.InfluxDB) > 0 {
		// Check if database for project existed, if not, just create
		query := influxapi.Query{
			Command:  "create database " + monitorutil.ProjectDatabaseName,
			Database: monitorutil.ProjectDatabaseName,
		}
		// Wait unitl influxdb is OK
		_ = wait.PollImmediateInfinite(10*time.Second, func() (bool, error) {
			for _, client := range c.remoteClient.InfluxDB {
				log.Debugf("Query sql: %s", query.Command)
				resp, err := client.Client.Query(query)
				if err != nil {
					log.Errorf("Create database %s for %s failed: %v", monitorutil.ProjectDatabaseName, client.Address, err)
					return false, nil
				} else if resp.Error() != nil {
					log.Errorf("Create database %s for %s failed: %v", monitorutil.ProjectDatabaseName, client.Address, resp.Error())
					return false, nil
				}
			}
			log.Info("Created database projects in influxdb")
			return true, nil
		})
	}

	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

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

// worker processes the queue of objects.
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

	err := c.syncCollector(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing collector %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncCollector will sync the Prometheus with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncCollector(key string) error {
	startTime := time.Now()
	var cachedCollector *cachedCollector
	defer func() {
		log.Info("Finished syncing collector", log.String("collectorKey", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// collector holds the latest collector info from apiserver
	collector, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Collector has been deleted. Attempting to cleanup resources", log.String("collectorKey", key))
		err = c.processCollectorDeletion(key)
	case err != nil:
		log.Warn("Unable to retrieve collector from store", log.String("collectorKey", key), log.Err(err))
	default:
		cachedCollector = c.cache.getOrCreate(key)
		err = c.processCollectorUpdate(cachedCollector, collector, key)
	}
	return err
}

func (c *Controller) processCollectorDeletion(key string) error {
	cachedCollector, ok := c.cache.get(key)
	if !ok {
		log.Error("Collector not in cache even though the watcher thought it was. Ignoring the deletion", log.String("collectorKey", key))
		return nil
	}
	return c.processCollectorDelete(cachedCollector, key)
}

func (c *Controller) processCollectorDelete(cachedCollector *cachedCollector, key string) error {
	log.Info("Collector will be dropped", log.String("collectorKey", key))

	collector := cachedCollector.state
	err := c.uninstallCollector(collector, true)
	if err != nil {
		log.Errorf("Collector uninstall fail: %v", err)
		return err
	}

	if c.cache.Exist(key) {
		log.Info("Delete the collector cache", log.String("collectorKey", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("Delete the collector health cache", log.String("collectorKey", key))
		c.health.Delete(key)
	}

	return nil
}

func (c *Controller) processCollectorUpdate(cachedCollector *cachedCollector, collector *v1.Collector, key string) error {
	if cachedCollector.state != nil {
		// exist and the cluster name changed
		if cachedCollector.state.UID != collector.UID {
			log.Info("Collector uid has changed, just delete", log.String("collectorKey", key))
			if err := c.processCollectorDelete(cachedCollector, key); err != nil {
				return err
			}
		}
	}
	notifyAPIConfigMap, err := c.platformClient.ConfigMaps().Get(notifyapi.NotifyApiConfigMapName, metav1.GetOptions{})
	if err == nil && notifyAPIConfigMap != nil {
		if v, ok := notifyAPIConfigMap.Annotations[notifyapi.NotifyAPIAddressKey]; ok {
			if c.notifyAPIAddress != v {
				c.notifyAPIAddress = v
			}
		}
	}
	err = c.createCollectorIfNeeded(key, collector)
	if err != nil {
		return err
	}

	cachedCollector.state = collector
	// Always update the cache upon success.
	c.cache.set(key, cachedCollector)
	return nil
}

func (c *Controller) collectorReinitialize(key string, collector *v1.Collector) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		log.Info("Reinitialize, try to reinstall", log.String("collectorKey", key))
		if err := c.uninstallCollector(collector, false); err != nil {
			log.Error("Failed to uninstall collector", log.Err(err))
			// continue
		}
		err := c.installCollector(collector)
		if err == nil {
			collector = collector.DeepCopy()
			collector.Status.Phase = v1.CollectorPhaseChecking
			collector.Status.Reason = ""
			collector.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(collector)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		log.Info("Reinitialize, try to uninstall", log.String("collectorKey", key))
		// First, rollback the collector
		if err := c.uninstallCollector(collector, false); err != nil {
			log.Error("Uninstall collector error.", log.String("collectorKey", key))
			return true, err
		}
		if collector.Status.RetryCount == collectorMaxRetryCount {
			log.Error("Collector reinitialize exceed max retry, set failed", log.String("collectorKey", key))
			collector = collector.DeepCopy()
			collector.Status.Phase = v1.CollectorPhaseFailed
			collector.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", collectorMaxRetryCount)
			err := c.persistUpdate(collector)
			if err != nil {
				log.Error("Update collector error", log.Err(err))
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		collector = collector.DeepCopy()
		collector.Status.Phase = v1.CollectorPhaseReinitializing
		collector.Status.Reason = err.Error()
		collector.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		collector.Status.RetryCount++
		err = c.persistUpdate(collector)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createCollectorIfNeeded(key string, collector *v1.Collector) error {
	switch collector.Status.Phase {
	case v1.CollectorPhaseInitializing:
		log.Info("Collector will be created", log.String("collectorKey", key))
		err := c.installCollector(collector)
		if err == nil {
			log.Info("Collector created success", log.String("collectorKey", key))
			collector = collector.DeepCopy()
			collector.Status.Phase = v1.CollectorPhaseChecking
			collector.Status.Reason = ""
			collector.Status.RetryCount = 0
			return c.persistUpdate(collector)
		}
		log.Error(fmt.Sprintf("Collector created failed: %v", err), log.String("collectorKey", key))
		collector = collector.DeepCopy()
		collector.Status.Phase = v1.CollectorPhaseReinitializing
		collector.Status.Reason = err.Error()
		collector.Status.RetryCount = 1
		return c.persistUpdate(collector)
	case v1.CollectorPhaseReinitializing:
		log.Info("Collector entry Reinitializing", log.String("collectorKey", key))
		var interval = time.Since(collector.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= collectorTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = collectorTimeOut - interval
		}
		go func() {
			_ = wait.Poll(waitTime, collectorTimeOut, c.collectorReinitialize(key, collector))
		}()
	case v1.CollectorPhaseChecking:
		log.Info("Collector entry Checking", log.String("collectorKey", key))
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, collector)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				_ = wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkCollectorStatus(collector, key, initDelay))
			}()
		}
	case v1.CollectorPhaseRunning:
		log.Info("Collector entry Running", log.String("collectorKey", key))
		if needUpgrade(collector) {
			c.health.Delete(key)
			collector = collector.DeepCopy()
			collector.Status.Phase = v1.CollectorPhaseUpgrading
			collector.Status.Reason = ""
			collector.Status.RetryCount = 0
			return c.persistUpdate(collector)
		}

		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, collector)
			go func() {
				_ = wait.PollImmediateUntil(5*time.Minute, c.watchCollectorHealth(key), c.stopCh)
			}()
		}
	case v1.CollectorPhaseUpgrading:
		log.Info("Collector entry upgrading", log.String("collectorKey", key))
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, collector)
			delay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.upgrading.Delete(key)
				_ = wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkCollectorUpgrade(collector, key, delay))
			}()
		}
	case v1.CollectorPhaseFailed:
		log.Info("Collector entry fail", log.String("collectorKey", key))
		c.upgrading.Delete(key)
		c.health.Delete(key)

		// should try check when collector recover again
		log.Info("Collector try checking after fail", log.String("collectorKey", key))
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, collector)
			delayTime := time.Now().Add(2 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				_ = wait.PollImmediate(20*time.Second, 1*time.Minute, c.checkCollectorStatus(collector, key, delayTime))
			}()
		}
	}
	return nil
}

func (c *Controller) installCollector(collector *v1.Collector) error {
	if c.notifyAPIAddress == "" {
		return fmt.Errorf("empty notify api address, check if notify api exists")
	}

	components := config.Get(collector.Spec.Version)
	collector.Status.Version = collector.Spec.Version
	if collector.Status.Components == nil {
		collector.Status.Components = make(map[string]string)
	}

	cluster, err := c.platformClient.Clusters().Get(collector.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	crdClient, err := util.BuildExternalExtensionClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	mclient, err := util.BuildExternalMonitoringClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	// Set remote write address
	var remoteWrites []string
	if len(c.remoteClient.InfluxDB) > 0 {
		remoteWrites, err = c.initInfluxdb(cluster.Name)
		if err != nil {
			return err
		}
	} else if len(c.remoteClient.ES) > 0 {
		remoteWrites, err = c.initESAdapter(kubeClient, &components)
		if err != nil {
			return err
		}
		collector.Status.Components[PrometheusBeatService] = components.PrometheusBeatWorkLoad.Tag
	}

	if len(collector.Spec.Storage.WriteAddr) > 0 {
		remoteWrites = collector.Spec.Storage.WriteAddr
	}
	// For remote read, just set from spec
	var remoteReads []string
	if len(collector.Spec.Storage.ReadAddr) > 0 {
		remoteReads = collector.Spec.Storage.ReadAddr
	}

	for _, crdKind := range crdKinds {
		crd := promopk8sutil.NewCustomResourceDefinition(crdKind, monitoring.GroupName, nil, true)
		_, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
		if err != nil {
			return err
		}
	}

	// Service prometheus-operator
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(servicePrometheusOperator()); err != nil {
		return err
	}
	// ServiceAccount for prometheus-operator
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountPrometheusOperator()); err != nil {
		return err
	}
	// ClusterRole for prometheus-operator
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(clusterRolePrometheusOperator()); err != nil {
		return err
	}
	// ClusterRoleBinding prometheus-operator
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(clusterRoleBindingPrometheusOperator()); err != nil {
		return err
	}
	// Deployment for prometheus-operator
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deployPrometheusOperatorApps(&components)); err != nil {
		return err
	}

	collector.Status.Components[PrometheusOperatorService] = components.PrometheusOperatorService.Tag

	extensionsAPIGroup := controllerutil.IsClusterVersionBefore1_9(kubeClient)

	// get notify webhook address
	var webhookAddr string
	if collector.Spec.NotifyWebhook != "" {
		webhookAddr = collector.Spec.NotifyWebhook
	} else {
		webhookAddr = c.notifyAPIAddress + "/webhook"
	}

	// secret for alertmanager
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(createSecretForAlertmanager(webhookAddr)); err != nil {
		return err
	}

	// ServiceAccount for alertmanager
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountAlertmanager()); err != nil {
		return err
	}

	// Service for alertmanager
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(createServiceForAlerterManager()); err != nil {
		return err
	}

	// Crd alertmanager instance
	if _, err := mclient.MonitoringV1().Alertmanagers(metav1.NamespaceSystem).Create(createAlertManagerCRD(&components)); err != nil {
		return err
	}

	collector.Status.Components[AlertManagerService] = components.AlertManagerService.Tag

	// Secret for prometheus-etcd
	credential, err := util.ClusterCredentialV1(c.platformClient, cluster.Name)
	if err != nil {
		return err
	}
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(secretETCDPrometheus(credential)); err != nil {
		return err
	}
	// Service Prometheus
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(servicePrometheus()); err != nil {
		return err
	}
	// Secret for prometheus
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(createSecretForPrometheus()); err != nil {
		return err
	}
	// ServiceAccount for prometheus
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(serviceAccountPrometheus()); err != nil {
		return err
	}
	// ClusterRole for prometheus
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(clusterRolePrometheus()); err != nil {
		return err
	}
	// ClusterRoleBinding Prometheus
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(clusterRoleBindingPrometheus()); err != nil {
		return err
	}
	// prometheus rule record
	if _, err := mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Create(recordsForPrometheus()); err != nil {
		return err
	}
	// prometheus rule alert, empty for now, edit by tke-monitor
	if _, err := mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Create(alertsForPrometheus()); err != nil {
		return err
	}
	// Crd prometheus instance
	if _, err := mclient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).Create(createPrometheusCRD(cluster.Name, remoteWrites, remoteReads, &components, len(c.remoteClient.InfluxDB) > 0)); err != nil {
		return err
	}
	collector.Status.Components[PrometheusService] = components.PrometheusService.Tag

	// DaemonSet for node-exporter
	if _, err := kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(createDaemonSetForNodeExporter(&components)); err != nil {
		return err
	}
	collector.Status.Components[nodeExporterService] = components.NodeExporterService.Tag

	// Service for kube-state-metrics
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(createServiceForMetrics()); err != nil {
		return err
	}
	// ServiceAccount for kube-state-metrics
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(createServiceAccountForMetrics()); err != nil {
		return err
	}
	// ClusterRole for kube-state-metrics
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(createClusterRoleForMetrics()); err != nil {
		return err
	}
	// ClusterRoleBinding for kube-state-metrics
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(createClusterRoleBindingForMetrics()); err != nil {
		return err
	}
	// Role for kube-state-metrics
	if _, err := kubeClient.RbacV1().Roles(metav1.NamespaceSystem).Create(createRoleForMetrics()); err != nil {
		return err
	}
	// RoleBinding for kube-state-metrics
	if _, err := kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Create(createRoleBingdingForMetrics()); err != nil {
		return err
	}
	// Deployment for kube-state-metrics
	if extensionsAPIGroup {
		if _, err := kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Create(createExtensionDeploymentForMetrics(&components)); err != nil {
			return err
		}
	} else {
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(createAppsDeploymentForMetrics(&components)); err != nil {
			return err
		}
	}
	collector.Status.Components[kubeStateService] = components.KubeStateService.Tag

	return nil
}

var selectorForPrometheusOperator = metav1.LabelSelector{
	MatchLabels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "prometheus-operator"},
}

func servicePrometheusOperator() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusOperatorService,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/name": "Prometheus-Operator", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/cluster-service": "true"},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForPrometheusOperator.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: PrometheusOperatorServicePort, Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP},
			},
		},
	}
}

func serviceAccountPrometheusOperator() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusOperatorServiceAccount,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
	}
}

func clusterRolePrometheusOperator() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   prometheusOperatorClusterRole,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{monitoring.GroupName},
				Resources: []string{"alertmanagers", "prometheuses", "prometheuses/finalizers", "alertmanagers/finalizers", "servicemonitors", "podmonitors", "prometheusrules"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"list", "delete"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"services", "services/finalizers", "endpoints"},
				Verbs:     []string{"get", "create", "update", "delete"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"list", "watch", "get"},
			},
		},
	}
}

func clusterRoleBindingPrometheusOperator() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   prometheusOperatorClusterRoleBinding,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusOperatorClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      prometheusOperatorServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func deployPrometheusOperatorApps(component *config.Components) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusOperatorWorkLoad,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{specialLabelName: specialLabelValue, "k8s-app": "prometheus-operator", "kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForPrometheusOperator,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{specialLabelName: specialLabelValue, "k8s-app": "prometheus-operator", "kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
					Annotations: map[string]string{"prometheus.io/scrape": "false"},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: prometheusOperatorServiceAccount,
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "node-role.kubernetes.io/master",
												Operator: corev1.NodeSelectorOpExists,
											},
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  prometheusOperatorWorkLoad,
							Image: component.PrometheusOperatorService.FullName(),
							Args: []string{
								"--kubelet-service=kube-system/kubelet",
								"--logtostderr=true",
								"--config-reloader-image=" + component.ConfigMapReloadWorkLoad.FullName(),
								"--prometheus-config-reloader=" + component.PrometheusConfigReloaderWorkload.FullName(),
							},
							// Command:         []string{"tail", "-f"},
							// Command: []string{"/bin/sh", "-c", "./prometheus --storage.tsdb.retention=1h --storage.tsdb.path=/data --web.enable-lifecycle --config.file=config/prometheus.yml"},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(100*1024*1024, resource.BinarySI),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(200, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(200*1024*1024, resource.BinarySI),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: controllerutil.BoolPtr(false),
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser:    controllerutil.Int64Ptr(65534),
						RunAsNonRoot: controllerutil.BoolPtr(true),
					},
				},
			},
		},
	}
}

var selectorForPrometheus = metav1.LabelSelector{
	MatchLabels: map[string]string{PrometheusService: PrometheusCRDName, "app": "prometheus"},
}

func servicePrometheus() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusService,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/name": "Prometheus", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/cluster-service": "true"},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForPrometheus.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: PrometheusServicePort, Port: 9090, TargetPort: intstr.FromInt(9090), Protocol: corev1.ProtocolTCP},
			},
		},
	}
}

func serviceAccountPrometheus() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusServiceAccount,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
	}
}

func clusterRolePrometheus() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   prometheusClusterRole,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "nodes/proxy", "nodes/metrics", "services", "endpoints", "pods"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				NonResourceURLs: []string{"/metrics"},
				Verbs:           []string{"get"},
			},
		},
	}
}

func clusterRoleBindingPrometheus() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   prometheusClusterRoleBinding,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      prometheusServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createPrometheusCRD(clusterName string, remoteWrites, remoteReads []string, component *config.Components, isInfluxDB bool) *monitoringv1.Prometheus {
	var remoteReadSpecs []monitoringv1.RemoteReadSpec
	for _, r := range remoteReads {
		if r == "nil" {
			continue
		}
		rr := monitoringv1.RemoteReadSpec{
			URL: r,
		}
		remoteReadSpecs = append(remoteReadSpecs, rr)
	}
	var remoteWriteSpecs []monitoringv1.RemoteWriteSpec
	for _, w := range remoteWrites {
		if w == "nil" {
			continue
		}
		rw := monitoringv1.RemoteWriteSpec{
			URL: w,
		}
		if isInfluxDB {
			if strings.Contains(w, "db=projects") {
				rw.WriteRelabelConfigs = []monitoringv1.RelabelConfig{
					{
						SourceLabels: []string{"__name__"},
						Regex:        "project_(.*)",
						Action:       "keep",
					},
				}
				rw.QueueConfig = &monitoringv1.QueueConfig{
					Capacity:          100,
					MinShards:         10,
					MaxShards:         10,
					MaxSamplesPerSend: 100,
					BatchSendDeadline: "30s",
				}
			} else {
				rw.WriteRelabelConfigs = []monitoringv1.RelabelConfig{
					{
						SourceLabels: []string{"__name__"},
						Regex:        "k8s_(.*)|kube_pod_labels|etcd_server_leader_changes_seen_total|etcd_debugging_mvcc_db_total_size_in_bytes",
						Action:       "keep",
					},
				}
				rw.QueueConfig = &monitoringv1.QueueConfig{
					Capacity:          10000,
					MinShards:         1000,
					MaxShards:         1000,
					MaxSamplesPerSend: 1000,
					BatchSendDeadline: "30s",
				}
			}
		} else {
			rw.WriteRelabelConfigs = []monitoringv1.RelabelConfig{
				{
					SourceLabels: []string{"__name__"},
					Regex:        "project_(.*)|k8s_(.*)|kube_pod_labels|etcd_server_leader_changes_seen_total|etcd_debugging_mvcc_db_total_size_in_bytes",
					Action:       "keep",
				},
			}
			rw.QueueConfig = &monitoringv1.QueueConfig{
				Capacity:          10000,
				MinShards:         1000,
				MaxShards:         1000,
				MaxSamplesPerSend: 1000,
				BatchSendDeadline: "30s",
			}
		}

		remoteWriteSpecs = append(remoteWriteSpecs, rw)
	}
	return &monitoringv1.Prometheus{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.GroupName + "/v1",
			Kind:       monitoringv1.PrometheusesKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusCRDName,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{specialLabelName: specialLabelValue, "k8s-app": "prometheus", "kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Spec: monitoringv1.PrometheusSpec{
			PodMetadata: &metav1.ObjectMeta{
				CreationTimestamp: metav1.Now(), // For validation only: https://github.com/coreos/prometheus-operator/issues/2399
				Annotations: map[string]string{
					"prometheus.io/scrape": "true",
					"prometheus.io/port":   "9090",
				},
			},
			ExternalLabels:     map[string]string{"cluster_id": clusterName},
			ScrapeInterval:     "60s",
			RemoteRead:         remoteReadSpecs,
			RemoteWrite:        remoteWriteSpecs,
			EvaluationInterval: "1m",
			AdditionalScrapeConfigs: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: prometheusSecret},
				Key:                  prometheusConfigName,
				Optional:             controllerutil.BoolPtr(false),
			},
			Secrets: []string{prometheusETCDSecret},
			Alerting: &monitoringv1.AlertingSpec{
				Alertmanagers: []monitoringv1.AlertmanagerEndpoints{
					{
						Namespace: metav1.NamespaceSystem,
						Name:      AlertManagerService,
						Port:      intstr.FromString(alertManagerServicePort),
					},
				},
			},
			BaseImage: containerregistryutil.GetImagePrefix(prometheusImagePath),
			Replicas:  controllerutil.Int32Ptr(1),
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
					corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    *resource.NewMilliQuantity(1000, resource.DecimalSI),
					corev1.ResourceMemory: *resource.NewQuantity(2*1024*1024*1024, resource.BinarySI),
				},
			},
			Tolerations: []corev1.Toleration{
				{
					Key:      "node-role.kubernetes.io/master",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
			Affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "node-role.kubernetes.io/master",
										Operator: corev1.NodeSelectorOpExists,
									},
								},
							},
						},
					},
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				FSGroup:      controllerutil.Int64Ptr(2000),
				RunAsNonRoot: controllerutil.BoolPtr(true),
				RunAsUser:    controllerutil.Int64Ptr(1000),
			},
			RuleSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{PrometheusService: PrometheusCRDName, "role": "alert-rules"},
			},
			ServiceAccountName:              prometheusServiceAccount,
			ServiceMonitorNamespaceSelector: &metav1.LabelSelector{},
			ServiceMonitorSelector:          &metav1.LabelSelector{},
			PodMonitorNamespaceSelector:     &metav1.LabelSelector{},
			PodMonitorSelector:              &metav1.LabelSelector{},
			Version:                         component.PrometheusService.Tag,
		},
	}
}

func createSecretForPrometheus() *corev1.Secret {
	cfg := scrapeConfigForPrometheus()

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusSecret,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			prometheusConfigName: []byte(cfg),
		},
	}
}

func secretETCDPrometheus(cred *platformv1.ClusterCredential) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusETCDSecret,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"etcd-ca.crt":     cred.ETCDCACert,
			"etcd-client.crt": cred.ETCDAPIClientCert,
			"etcd-client.key": cred.ETCDAPIClientKey,
		},
	}
}

func recordsForPrometheus() *monitoringv1.PrometheusRule {
	records := recordRulesForPrometheus()
	reader := strings.NewReader(records)
	prometheusRuleSpec := &monitoringv1.PrometheusRuleSpec{}
	err := yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(prometheusRuleSpec)
	if err != nil {
		log.Error("decode record err", log.String("err", err.Error()))
		return nil
	}
	return &monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.GroupName + "/v1",
			Kind:       monitoringv1.PrometheusRuleKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusRuleRecord,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{PrometheusService: PrometheusCRDName, "role": "alert-rules"},
		},
		Spec: *prometheusRuleSpec,
	}
}

func alertsForPrometheus() *monitoringv1.PrometheusRule {
	return &monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.GroupName + "/v1",
			Kind:       monitoringv1.PrometheusRuleKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusRuleAlert,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{PrometheusService: PrometheusCRDName, "role": "alert-rules"},
		},
		Spec: monitoringv1.PrometheusRuleSpec{Groups: []monitoringv1.RuleGroup{}},
	}
}

var selectorForAlertManager = metav1.LabelSelector{
	MatchLabels: map[string]string{"alertmanager": alertManagerCRDName, "app": "alertmanager"},
}

func createSecretForAlertmanager(webhookAddr string) *corev1.Secret {
	cfg := configForAlertManager(webhookAddr)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alertManagerSecret,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"alertmanager.yaml": []byte(cfg),
		},
	}
}

func serviceAccountAlertmanager() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      alertManagerServiceAccount,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
	}
}

func createAlertManagerCRD(component *config.Components) *monitoringv1.Alertmanager {
	return &monitoringv1.Alertmanager{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.GroupName + "/v1",
			Kind:       monitoringv1.AlertmanagersKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      alertManagerCRDName,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{specialLabelName: specialLabelValue, "k8s-app": "alertmanager", "kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Spec: monitoringv1.AlertmanagerSpec{
			PodMetadata: &metav1.ObjectMeta{
				CreationTimestamp: metav1.Now(), // For validation only: https://github.com/coreos/prometheus-operator/issues/2399
				Annotations: map[string]string{
					"prometheus.io/scrape": "true",
					"prometheus.io/port":   "9093",
				},
			},
			BaseImage: containerregistryutil.GetImagePrefix(alertManagerImagePath),
			Replicas:  controllerutil.Int32Ptr(3),
			SecurityContext: &corev1.PodSecurityContext{
				FSGroup:      controllerutil.Int64Ptr(2000),
				RunAsNonRoot: controllerutil.BoolPtr(true),
				RunAsUser:    controllerutil.Int64Ptr(1000),
			},
			ServiceAccountName: alertManagerServiceAccount,
			Version:            component.AlertManagerService.Tag,
		},
	}
}

func createServiceForAlerterManager() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        AlertManagerService,
			Namespace:   metav1.NamespaceSystem,
			Labels:      map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/name": "Alertmanager"},
			Annotations: map[string]string{"prometheus.io/scrape": "false"},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForAlertManager.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: alertManagerServicePort, Port: 80, TargetPort: intstr.FromInt(9093), Protocol: corev1.ProtocolTCP},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

var selectorForNodeExporter = metav1.LabelSelector{
	MatchLabels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "node-exporter"},
}

func createDaemonSetForNodeExporter(component *config.Components) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodeExporterDaemonSet,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", specialLabelName: specialLabelValue, "k8s-app": "node-exporter"},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &selectorForNodeExporter,
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{specialLabelName: specialLabelValue, "k8s-app": "node-exporter"},
					Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": "", "tke.prometheus.io/scrape": "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  nodeExporterDaemonSet,
							Image: component.NodeExporterService.FullName(),
							Args: []string{
								"--path.procfs=/host/proc",
								"--path.sysfs=/host/sys",
								"--no-collector.arp",
								"--no-collector.bcache",
								"--no-collector.bonding",
								"--no-collector.buddyinfo",
								"--no-collector.conntrack",
								"--no-collector.cpu",
								"--collector.diskstats",
								"--no-collector.drbd",
								"--no-collector.edac",
								"--no-collector.entropy",
								"--no-collector.filefd",
								"--collector.filesystem",
								"--no-collector.gmond",
								"--no-collector.hwmon",
								"--no-collector.infiniband",
								"--no-collector.interrupts",
								"--no-collector.ipvs",
								"--no-collector.ksmd",
								"--no-collector.loadavg",
								"--no-collector.logind",
								"--no-collector.mdadm",
								"--no-collector.megacli",
								"--no-collector.meminfo",
								"--no-collector.meminfo_numa",
								"--no-collector.mountstats",
								"--collector.netdev",
								"--no-collector.netstat",
								"--no-collector.nfs",
								"--no-collector.ntp",
								"--no-collector.qdisc",
								"--no-collector.runit",
								"--collector.sockstat",
								"--no-collector.stat",
								"--no-collector.supervisord",
								"--no-collector.systemd",
								"--no-collector.tcpstat",
								"--no-collector.textfile",
								"--no-collector.time",
								"--no-collector.uname",
								"--no-collector.vmstat",
								"--no-collector.wifi",
								"--no-collector.xfs",
								"--no-collector.zfs",
								"--no-collector.timex",
							},
							Ports: []corev1.ContainerPort{
								{Name: "metrics", ContainerPort: 9100, HostPort: 9100},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(1*1024*1024*1024, resource.BinarySI),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/host/proc",
									Name:      "proc",
									ReadOnly:  true,
								},
								{
									MountPath: "/host/sys",
									Name:      "sys",
									ReadOnly:  true,
								},
							},
						},
					},
					HostNetwork: true,
					HostPID:     true,
					Volumes: []corev1.Volume{
						{
							Name: "proc",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/proc",
								},
							},
						},
						{
							Name: "sys",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/sys",
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
				},
			},
		},
	}
}

var selectorForMetrics = metav1.LabelSelector{
	MatchLabels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "kube-state-metrics"},
}

func createServiceForMetrics() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        kubeStateService,
			Namespace:   metav1.NamespaceSystem,
			Labels:      map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/name": "kube-state-metrics"},
			Annotations: map[string]string{"tke.prometheus.io/scrape": "true"},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForMetrics.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: "http-metrics", Port: 8080, TargetPort: intstr.FromString("http-metrics"), Protocol: corev1.ProtocolTCP},
			},
		},
	}
}

func createServiceAccountForMetrics() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeStateServiceAccount,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
	}
}

func createClusterRoleForMetrics() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   kubeStateClusterRole,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps", "secrets", "nodes", "resourcequotas", "services", "endpoints", "pods", "limitranges",
					"replicationcontrollers", "persistentvolumeclaims", "persistentvolumes", "namespaces"},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"daemonsets", "deployments", "replicasets"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets", "daemonsets", "deployments"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{"cronjobs", "jobs"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"autoscaling"},
				Resources: []string{"horizontalpodautoscalers"},
				Verbs:     []string{"list", "watch"},
			},
		},
	}
}

func createClusterRoleBindingForMetrics() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   kubeStateClusterRoleBinding,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     kubeStateClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      kubeStateServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createRoleForMetrics() *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeStateRole,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups:     []string{"extensions"},
				Resources:     []string{"deployments"},
				ResourceNames: []string{"kube-state-metrics"},
				Verbs:         []string{"get", "update"},
			},
		},
	}
}

func createRoleBingdingForMetrics() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeStateRoleBinding,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     kubeStateRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      kubeStateServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createExtensionDeploymentForMetrics(component *config.Components) *extensionsv1beta1.Deployment {
	return &extensionsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeStateWorkLoad,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", specialLabelName: specialLabelValue, "k8s-app": "kube-state-metrics"},
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForMetrics,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "kube-state-metrics"},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: kubeStateServiceAccount,
					Containers: []corev1.Container{
						{
							Name:  kubeStateWorkLoad,
							Image: component.KubeStateService.FullName(),
							Args: []string{
								"--port=8080",
								"--telemetry-port=8081",
							},
							Ports: []corev1.ContainerPort{
								{Name: "http-metrics", ContainerPort: 8080},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(1000, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(2*1024*1024*1024, resource.BinarySI),
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
				},
			},
		},
	}
}

func createAppsDeploymentForMetrics(component *config.Components) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeStateWorkLoad,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", specialLabelName: specialLabelValue, "k8s-app": "kube-state-metrics"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForMetrics,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "kube-state-metrics"},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: kubeStateServiceAccount,
					Containers: []corev1.Container{
						{
							Name:  kubeStateWorkLoad,
							Image: component.KubeStateService.FullName(),
							Args: []string{
								"--port=8080",
								"--telemetry-port=8081",
							},
							Ports: []corev1.ContainerPort{
								{Name: "http-metrics", ContainerPort: 8080},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(1000, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewQuantity(2*1024*1024*1024, resource.BinarySI),
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
				},
			},
		},
	}
}

func (c *Controller) uninstallCollector(collector *v1.Collector, dropData bool) error {
	var err error
	var errs []error

	cluster, err := c.platformClient.Clusters().Get(collector.Spec.ClusterName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	kubeClient, err := util.BuildExternalClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	crdClient, err := util.BuildExternalExtensionClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	mclient, err := util.BuildExternalMonitoringClientSet(cluster, c.platformClient)
	if err != nil {
		return err
	}

	extensionsAPIGroup := controllerutil.IsClusterVersionBefore1_9(kubeClient)

	// delete prometheus
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(prometheusETCDSecret, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(prometheusSecret, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(prometheusClusterRoleBinding, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(prometheusClusterRole, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(PrometheusService, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(prometheusServiceAccount, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	err = mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Delete(prometheusRuleRecord, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Delete(PrometheusRuleAlert, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	err = mclient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).Delete(PrometheusCRDName, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete alertmanager
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(AlertManagerService, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(alertManagerSecret, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(alertManagerServiceAccount, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = mclient.MonitoringV1().Alertmanagers(metav1.NamespaceSystem).Delete(alertManagerCRDName, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete node-exporter
	err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Delete(nodeExporterDaemonSet, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete kube-state-metrics
	if extensionsAPIGroup {
		err = kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Delete(kubeStateWorkLoad, &metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}

		// For extension group, should delete replicaset and pod additionally
		selector, err := metav1.LabelSelectorAsSelector(&selectorForMetrics)
		if err != nil {
			errs = append(errs, err)
		} else {
			options := metav1.ListOptions{
				LabelSelector: selector.String(),
			}
			err = controllerutil.DeleteReplicaSetApp(kubeClient, options)
			if err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(kubeStateWorkLoad, &metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(kubeStateClusterRoleBinding, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(kubeStateClusterRole, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Delete(kubeStateRoleBinding, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().Roles(metav1.NamespaceSystem).Delete(kubeStateRole, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(kubeStateServiceAccount, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(kubeStateService, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete prometheus-operator
	err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(prometheusOperatorWorkLoad, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(prometheusOperatorClusterRoleBinding, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(prometheusOperatorClusterRole, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(PrometheusOperatorService, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(prometheusOperatorServiceAccount, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	for _, crdKind := range crdKinds {
		crd := promopk8sutil.NewCustomResourceDefinition(crdKind, monitoring.GroupName, nil, true)
		err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, &metav1.DeleteOptions{})
		if err != nil {
			errs = append(errs, err)
		}
	}

	// drop influxdb data, may take long time
	if dropData {
		if len(c.remoteClient.InfluxDB) > 0 {
			go func() {
				if err := c.dropInfluxdb(cluster.Name); err != nil {
					log.Error("Failed to drop influxdb", log.Err(err))
				}
			}()
		}
	}

	if len(c.remoteClient.ES) > 0 {
		err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(prometheusBeatWorkLoad, &metav1.DeleteOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete(PrometheusBeatConfigmap, &metav1.DeleteOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(PrometheusBeatService, &metav1.DeleteOptions{})
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		errMsg := ""
		for _, e := range errs {
			errMsg += e.Error() + ";"
		}
		return fmt.Errorf("delete prometheus fail:%s", errMsg)
	}

	return nil
}

func (c *Controller) watchCollectorHealth(key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check collector in cluster health", log.String("collectorKey", key))

		collector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.platformClient.Clusters().Get(collector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.health.Load(collector.Name); !ok {
			log.Info("Collector health check over", log.String("collectorKey", key))
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", PrometheusService, PrometheusServicePort, `/-/healthy`, nil).DoRaw(); err != nil {
			collector = collector.DeepCopy()
			collector.Status.Phase = v1.CollectorPhaseFailed
			collector.Status.Reason = "Collector is not healthy."
			if err = c.persistUpdate(collector); err != nil {
				return false, err
			}
			return true, nil
		}
		log.Debug("Collector health is ok", log.String("collectorKey", key))
		return false, nil
	}
}

func (c *Controller) checkCollectorStatus(collector *v1.Collector, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check collector status", log.String("collectorKey", key))

		cluster, err := c.platformClient.Clusters().Get(collector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.checking.Load(key); !ok {
			log.Info("Collector status checking over", log.String("collectorKey", key))
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		collector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", PrometheusService, PrometheusServicePort, `/-/healthy`, nil).DoRaw(); err != nil {
			if time.Now().After(initDelay) {
				collector = collector.DeepCopy()
				collector.Status.Phase = v1.CollectorPhaseFailed
				collector.Status.Reason = "Collector is not healthy."
				if err = c.persistUpdate(collector); err != nil {
					return false, err
				}
				return true, nil
			}
			log.Error("collector status has not healthy", log.String("collectorKey", key), log.Err(err))
			return false, nil
		}
		collector = collector.DeepCopy()
		collector.Status.Phase = v1.CollectorPhaseRunning
		collector.Status.Reason = ""
		if err = c.persistUpdate(collector); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) checkCollectorUpgrade(collector *v1.Collector, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade collector", log.String("collectorKey", key))

		cluster, err := c.platformClient.Clusters().Get(collector.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.upgrading.Load(key); !ok {
			log.Info("Collector upgrade over", log.String("collectorKey", key))
			return true, nil
		}
		kubeClient, err := util.BuildExternalClientSet(cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		collector, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		newComponents := config.Get(collector.Spec.Version)
		for name, version := range collector.Status.Components {
			if image := newComponents.Get(name); image != nil {
				if image.Tag != version {
					patch := fmt.Sprintf(`[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`, image.FullName())
					err = upgradeVersion(kubeClient, name, patch)
					if err != nil {
						if time.Now().After(initDelay) {
							collector = collector.DeepCopy()
							collector.Status.Phase = v1.CollectorPhaseUpgrading
							collector.Status.Reason = fmt.Sprintf("Upgrade timeout, %v", err)
							if err = c.persistUpdate(collector); err != nil {
								return false, err
							}
							return true, nil
						}
						return false, nil
					}
					collector.Status.Components[name] = image.Tag
				}
			}
		}

		collector = collector.DeepCopy()
		collector.Status.Version = collector.Spec.Version
		collector.Status.Phase = v1.CollectorPhaseChecking
		collector.Status.Reason = ""
		if err = c.persistUpdate(collector); err != nil {
			return false, err
		}

		return true, nil
	}
}

func needUpgrade(collector *v1.Collector) bool {
	if collector.Status.Components == nil {
		log.Errorf("Nil component version when checking upgrade!")
		return false
	}

	if collector.Spec.Version != collector.Status.Version {
		return true
	}

	return false
}

func upgradeVersion(kubeClient *kubernetes.Clientset, workLoad, patch string) error {
	var err error

	switch workLoad {
	case kubeStateWorkLoad, prometheusWorkLoad, AlertManagerWorkLoad:
		extensionsAPIGroup := controllerutil.IsClusterVersionBefore1_9(kubeClient)
		if extensionsAPIGroup {
			_, err = kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Patch(workLoad, types.JSONPatchType, []byte(patch))
		} else {
			_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Patch(workLoad, types.JSONPatchType, []byte(patch))
		}
	case nodeExporterDaemonSet:
		_, err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Patch(workLoad, types.JSONPatchType, []byte(patch))
	default:
		return fmt.Errorf("wrong workload: %s", workLoad)
	}

	return err
}

func (c *Controller) persistUpdate(collector *v1.Collector) error {
	var err error
	for i := 0; i < collectorClientRetryCount; i++ {
		_, err = c.monitorClient.Collectors().UpdateStatus(collector)
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to collector that no longer exists", log.String("collectorName", collector.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to collector '%s' that has been changed since we received it: %v", collector.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of collector '%s/%s'", collector.Name, collector.Status.Phase), log.String("collectorName", collector.Name), log.Err(err))
		time.Sleep(collectorClientRetryInterval)
	}

	return err
}

func (c *Controller) initInfluxdb(dbName string) ([]string, error) {
	log.Infof("Starting create influxdb table: %s", dbName)

	// generate password
	str := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var passwd []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		passwd = append(passwd, bytes[r.Intn(len(bytes))])
	}

	// remove invalid character
	db := monitorutil.RenameInfluxDB(dbName)
	usr := db

	// drop users if user already existed
	cmdUser := fmt.Sprintf("drop user %s", usr)
	queryUser := influxapi.Query{
		Command:  cmdUser,
		Database: monitorutil.ProjectDatabaseName,
	}

	// create database/user, and grant privilege
	cmdDB := fmt.Sprintf("create database %s; create user %s with password '%s'; grant all on %s to %s; grant write on %s to %s",
		db, usr, passwd, db, usr, monitorutil.ProjectDatabaseName, usr)
	log.Debugf("Create influxdb table: %s", cmdDB)
	queryDB := influxapi.Query{
		Command:  cmdDB,
		Database: monitorutil.ProjectDatabaseName,
	}

	var queryStr []string
	for _, client := range c.remoteClient.InfluxDB {
		_, _ = client.Client.Query(queryUser)

		resp, err := client.Client.Query(queryDB)
		if err != nil {
			return nil, err
		} else if resp.Error() != nil {
			return nil, resp.Error()
		}
		queryStr = append(queryStr, []string{
			fmt.Sprintf("%s/api/v1/prom/write?db=%s&u=%s&p=%s", client.Address, db, usr, passwd),
			fmt.Sprintf("%s/api/v1/prom/write?db=%s&u=%s&p=%s", client.Address, monitorutil.ProjectDatabaseName, usr, passwd),
		}...)
	}

	return queryStr, nil
}

func (c *Controller) dropInfluxdb(dbName string) error {
	log.Infof("Starting drop influxdb table: %s", dbName)

	// remove invalid character
	db := monitorutil.RenameInfluxDB(dbName)
	usr := db

	// drop user and database
	cmd := fmt.Sprintf("drop user %s; drop database %s", usr, db)
	query := influxapi.Query{
		Command:  cmd,
		Database: monitorutil.ProjectDatabaseName,
	}

	// just continue when error
	for _, client := range c.remoteClient.InfluxDB {
		resp, err := client.Client.Query(query)
		if err != nil {
			log.Errorf("Drop database(%s) for %s err: %v", dbName, client.Address, err)
		} else if resp.Error() != nil {
			log.Errorf("Drop database(%s) for %s err: %v", dbName, client.Address, resp.Error())
		}
	}

	return nil
}

func (c *Controller) initESAdapter(kubeClient *kubernetes.Clientset, component *config.Components) ([]string, error) {
	var (
		remoteWrites []string
		hosts        []string
		user         string
		password     string
	)

	selectorForPrometheusBeat := metav1.LabelSelector{
		MatchLabels: map[string]string{
			specialLabelName: specialLabelValue,
			"k8s-app":        "prometheus-beat",
		},
	}
	for _, client := range c.remoteClient.ES {
		hosts = append(hosts, client.URL)
		user = client.Username
		password = client.Password
	}
	// create prom-beat service
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusBeatService,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/name": "Prometheus-Beat", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/cluster-service": "true"},
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForPrometheusBeat.MatchLabels,
			Ports: []corev1.ServicePort{
				{Name: PrometheusBeatServicePort, Port: 8080, TargetPort: intstr.FromInt(8080), Protocol: corev1.ProtocolTCP},
			},
		},
	}
	_, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(svc)
	if err != nil {
		return remoteWrites, err
	}

	// create prom-beat configmap
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusBeatConfigmap,
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			PrometheusBeatConfigFile: configForPrometheusBeat(hosts, user, password),
		},
	}
	_, err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(cm)
	if err != nil {
		return remoteWrites, err
	}

	// create prom-beat deployment
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusBeatWorkLoad,
			Namespace: metav1.NamespaceSystem,
			Labels:    selectorForPrometheusBeat.MatchLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(2),
			Selector: &selectorForPrometheusBeat,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      selectorForPrometheusBeat.MatchLabels,
					Annotations: map[string]string{"prometheus.io/scrape": "false"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  prometheusBeatWorkLoad,
							Image: component.PrometheusBeatWorkLoad.FullName(),
							Args: []string{
								"-c",
								"config/prometheusbeat.yml",
								"-e",
								"-d",
								"*",
							},
							Command: []string{"./prometheusbeat"},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/config",
									Name:      PrometheusBeatConfigmap,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: PrometheusBeatConfigmap,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: PrometheusBeatConfigmap,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(deployment)
	if err != nil {
		return remoteWrites, err
	}
	remoteWrites = append(remoteWrites, fmt.Sprintf("http://%s:%d/prometheus", svc.Name, 8080))
	return remoteWrites, nil
}
