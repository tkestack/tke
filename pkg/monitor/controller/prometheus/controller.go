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

package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"tkestack.io/tke/pkg/platform/util/addon"

	"github.com/coreos/prometheus-operator/pkg/apis/monitoring"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	influxApi "github.com/influxdata/influxdb1-client/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	clientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	monitorv1informer "tkestack.io/tke/api/client/informers/externalversions/monitor/v1"
	monitorv1lister "tkestack.io/tke/api/client/listers/monitor/v1"
	v1 "tkestack.io/tke/api/monitor/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	notifyapi "tkestack.io/tke/cmd/tke-notify-api/app"
	controllerutil "tkestack.io/tke/pkg/controller"
	"tkestack.io/tke/pkg/monitor/controller/prometheus/images"
	esutil "tkestack.io/tke/pkg/monitor/storage/es/client"
	monitorutil "tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/util/apiclient"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	utilhttp "tkestack.io/tke/pkg/util/http"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	prometheusClientRetryCount    = 5
	prometheusClientRetryInterval = 5 * time.Second

	prometheusMaxRetryCount = 5
	prometheusTimeOut       = 5 * time.Minute
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

	prometheusAdapterAuthDelegatorClusterRoleBinding  = "custom-metrics:system:auth-delegator"
	prometheusAdapterServiceAccount                   = "custom-metrics-apiserver"
	prometheusAdapterWorkLoad                         = "custom-metrics-apiserver"
	prometheusAdapterService                          = "custom-metrics-apiserver"
	prometheusAdapterClusterRole                      = "custom-metrics-server-resources"
	prometheusAdapterResourceReaderClusterRole        = "custom-metrics-resource-reader"
	prometheusAdapterResourceReaderClusterRoleBinding = "custom-metrics-resource-reader"
	prometheusAdapterHPAClusterRoleBinding            = "hpa-controller-custom-metrics"
	prometheusAdapterAuthReaderRoleBinding            = "custom-metrics-auth-reader"
	prometheusAdapterConfigMap                        = "adapter-config"
	systemAuthDelegatorClusterRole                    = "system:auth-delegator"

	nodeProblemDetectorWorkload           = "node-problem-detector"
	nodeProblemDetectorServiceAccount     = "node-problem-detector"
	nodeProblemDetectorClusterRoleBinding = "node-problem-detector"
	nodeProblemDetectorClusterRole        = "node-problem-detector"

	specialLabelName  = "k8s-submitter"
	specialLabelValue = "controller"
)

type influxdbClient struct {
	address string
	client  influxApi.Client
}

type thanosClient struct {
	address string
}

// more client will be added, now support influxdb, es and thanos
type remoteClient struct {
	influxdb *influxdbClient
	es       *esutil.Client
	thanos   *thanosClient
}

// Controller is responsible for performing actions dependent upon a prometheus phase.
type Controller struct {
	client         clientset.Interface
	platformClient platformv1client.PlatformV1Interface
	cache          *prometheusCache
	health         sync.Map
	checking       sync.Map
	upgrading      sync.Map
	queue          workqueue.RateLimitingInterface
	lister         monitorv1lister.PrometheusLister
	listerSynced   cache.InformerSynced
	stopCh         <-chan struct{}
	// RemoteAddress for prometheus
	remoteClients []remoteClient
	remoteType    string
	retentionDays int
	// NotifyApiAddress
	notifyAPIAddress string
}

// NewController creates a new Controller object.
// retentionDays is used only when storage is influx db. Controller will create database in inluxdb, and set the retention policy
func NewController(client clientset.Interface, platformClient platformv1client.PlatformV1Interface, prometheusInformer monitorv1informer.PrometheusInformer, resyncPeriod time.Duration, remoteAddress []string, remoteType string, retentionDays int) *Controller {
	// create the controller so we can inject the enqueue function
	controller := &Controller{
		client:         client,
		platformClient: platformClient,
		cache:          &prometheusCache{prometheusMap: make(map[string]*cachedPrometheus)},
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "prometheus"),
		remoteType:     remoteType,
		retentionDays:  retentionDays,
	}

	if client != nil && client.MonitorV1().RESTClient().GetRateLimiter() != nil {
		_ = metrics.RegisterMetricAndTrackRateLimiterUsage("prometheus_controller", client.MonitorV1().RESTClient().GetRateLimiter())
	}

	// configure the prometheus informer event handlers
	prometheusInformer.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueuePrometheus,
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPrometheus, ok1 := oldObj.(*v1.Prometheus)
				curPrometheus, ok2 := newObj.(*v1.Prometheus)
				if ok1 && ok2 && controller.needsUpdate(oldPrometheus, curPrometheus) {
					controller.enqueuePrometheus(newObj)
				}
			},
			DeleteFunc: controller.enqueuePrometheus,
		},
		resyncPeriod,
	)
	controller.lister = prometheusInformer.Lister()
	controller.listerSynced = prometheusInformer.Informer().HasSynced

	// construct remote client
	switch remoteType {
	case "influxdb":
		for _, remoteAddr := range remoteAddress {
			remoteClient := remoteClient{
				influxdb: &influxdbClient{},
			}
			client, err := monitorutil.ParseInfluxdb(remoteAddr)
			if err != nil {
				log.Errorf("Parse remote address %s to client failed: %v", remoteAddr, err)
				remoteClient.influxdb.client = nil
				remoteClient.influxdb.address = remoteAddr
			} else {
				attrs := strings.Split(remoteAddr, "&")
				remoteClient.influxdb.address = attrs[0]
				remoteClient.influxdb.client = client
			}
			controller.remoteClients = append(controller.remoteClients, remoteClient)
		}
	case "elasticsearch":
		for _, remoteAddr := range remoteAddress {
			remoteClient := remoteClient{
				es: &esutil.Client{},
			}
			client, err := monitorutil.ParseES(remoteAddr)
			if err != nil {
				log.Errorf("Parse remote address %s to client failed: %v", remoteAddr, err)
				remoteClient.es.URL = remoteAddr
			} else {
				remoteClient.es = &client
			}
			controller.remoteClients = append(controller.remoteClients, remoteClient)
		}
	case "thanos":
		for _, remoteAddr := range remoteAddress {
			remoteClient := remoteClient{
				thanos: &thanosClient{address: remoteAddr},
			}
			controller.remoteClients = append(controller.remoteClients, remoteClient)
		}
	default:
		log.Errorf("invalid remote type: %s", remoteType)
	}

	return controller
}

// obj could be an *v1.Prometheus, or a DeletionFinalStateUnknown marker item.
func (c *Controller) enqueuePrometheus(obj interface{}) {
	key, err := controllerutil.KeyFunc(obj)
	if err != nil {
		log.Error("Couldn't get key for object", log.Any("object", obj), log.Err(err))
		return
	}
	c.queue.Add(key)
}

func (c *Controller) needsUpdate(oldPrometheus *v1.Prometheus, newPrometheus *v1.Prometheus) bool {
	return !reflect.DeepEqual(oldPrometheus, newPrometheus)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	// Start the informer factories to begin populating the informer caches
	log.Info("Starting prometheus controller")
	defer log.Info("Shutting down prometheus controller")

	// Check if database for project existed, if not, just create
	query := influxApi.Query{
		Command:  "create database " + monitorutil.ProjectDatabaseName,
		Database: monitorutil.ProjectDatabaseName,
	}

	if c.remoteType == "influxdb" {
		// Wait unitl influxdb is OK
		_ = wait.PollImmediateInfinite(10*time.Second, func() (bool, error) {
			for _, client := range c.remoteClients {
				if client.influxdb.client == nil {
					return false, nil
				}
				log.Debugf("create database  sql: %s", query.Command)
				resp, err := client.influxdb.client.Query(query)
				if err != nil {
					log.Errorf("Create database %s for %s failed: %v", monitorutil.ProjectDatabaseName, client.influxdb.address, err)
					return false, nil
				} else if resp.Error() != nil {
					log.Errorf("Create database %s for %s failed: %v", monitorutil.ProjectDatabaseName, client.influxdb.address, resp.Error())
					return false, nil
				}
				// create of alter retention policy 'tke'
				err = createOrAlterRetention(client, monitorutil.ProjectDatabaseName, c.retentionDays)
				if err != nil {
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

	err := c.syncPrometheus(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing prometheus %v (will retry): %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

// syncPrometheus will sync the Prometheus with the given key if it has had
// its expectations fulfilled, meaning it did not expect to see any more of its
// namespaces created or deleted. This function is not meant to be invoked
// concurrently with the same key.
func (c *Controller) syncPrometheus(key string) error {
	startTime := time.Now()
	var cachedPrometheus *cachedPrometheus
	defer func() {
		log.Info("Finished syncing prometheus", log.String("prome", key), log.Duration("processTime", time.Since(startTime)))
	}()

	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// prometheus holds the latest prometheus info from apiserver
	prometheus, err := c.lister.Get(name)
	switch {
	case errors.IsNotFound(err):
		log.Info("Prometheus has been deleted. Attempting to cleanup resources", log.String("prome", key))
		err = c.processPrometheusDeletion(context.Background(), key)
	case err != nil:
		log.Warn("Unable to retrieve prometheus from store", log.String("prome", key), log.Err(err))
	default:
		cachedPrometheus = c.cache.getOrCreate(key)
		err = c.processPrometheusUpdate(context.Background(), cachedPrometheus, prometheus, key)
	}
	return err
}

func (c *Controller) processPrometheusDeletion(ctx context.Context, key string) error {
	cachedPrometheus, ok := c.cache.get(key)
	if !ok {
		log.Error("Prometheus not in cache even though the watcher thought it was. Ignoring the deletion", log.String("prome", key))
		return nil
	}
	return c.processPrometheusDelete(ctx, cachedPrometheus, key)
}

func (c *Controller) processPrometheusDelete(ctx context.Context, cachedPrometheus *cachedPrometheus, key string) error {
	log.Info("prometheus will be dropped", log.String("prome", key))

	prometheus := cachedPrometheus.state
	if prometheus != nil {
		DeleteMetricPrometheusStatusFail(prometheus.Spec.TenantID, prometheus.Spec.ClusterName, prometheus.Name)
	}
	err := c.uninstallPrometheus(ctx, prometheus, true)
	if err != nil {
		log.Errorf("Prometheus uninstall fail: %v", err)
		return err
	}

	if c.cache.Exist(key) {
		log.Info("delete the prometheus cache", log.String("prome", key))
		c.cache.delete(key)
	}

	if _, ok := c.health.Load(key); ok {
		log.Info("delete the prometheus health cache", log.String("prome", key))
		c.health.Delete(key)
	}

	return nil
}

func (c *Controller) processPrometheusUpdate(ctx context.Context, cachedPrometheus *cachedPrometheus, prometheus *v1.Prometheus, key string) error {
	if cachedPrometheus.state != nil {
		// exist and the cluster name changed
		if cachedPrometheus.state.UID != prometheus.UID {
			log.Info("Prometheus uid has chenged, just delete!", log.String("prome", key))
			if err := c.processPrometheusDelete(ctx, cachedPrometheus, key); err != nil {
				return err
			}
		}
	}
	notifyAPIConfigMap, err := c.platformClient.ConfigMaps().Get(ctx, notifyapi.NotifyApiConfigMapName, metav1.GetOptions{})
	if err == nil && notifyAPIConfigMap != nil {
		if v, ok := notifyAPIConfigMap.Annotations[notifyapi.NotifyAPIAddressKey]; ok {
			if c.notifyAPIAddress != v {
				c.notifyAPIAddress = v
			}
		}
	}
	err = c.createPrometheusIfNeeded(ctx, key, cachedPrometheus, prometheus)
	if err != nil {
		return err
	}

	cachedPrometheus.state = prometheus
	// Always update the cache upon success.
	c.cache.set(key, cachedPrometheus)
	return nil
}

func (c *Controller) prometheusReinitialize(ctx context.Context, key string, cachedPrometheus *cachedPrometheus, prometheus *v1.Prometheus) func() (bool, error) {
	// this func will always return true that keeps the poll once
	return func() (bool, error) {
		log.Info("Reinitialize, try to reinstall", log.String("prome", key))
		c.uninstallPrometheus(ctx, prometheus, false)
		err := c.installPrometheus(ctx, prometheus)
		if err == nil {
			prometheus = prometheus.DeepCopy()
			prometheus.Status.Phase = v1.AddonPhaseChecking
			prometheus.Status.Reason = ""
			prometheus.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
			err = c.persistUpdate(ctx, prometheus)
			if err != nil {
				return true, err
			}
			return true, nil
		}
		log.Info("Reinitialize, try to uninstall", log.String("prome", key))
		// First, rollback the prometheus
		if err := c.uninstallPrometheus(ctx, prometheus, false); err != nil {
			log.Error("Uninstall prometheus error.", log.String("prome", key))
			return true, err
		}
		if prometheus.Status.RetryCount == prometheusMaxRetryCount {
			log.Error("PrometheusReinitialize exceed max retry, set failed", log.String("prome", key))
			prometheus = prometheus.DeepCopy()
			prometheus.Status.Phase = v1.AddonPhaseFailed
			prometheus.Status.Reason = fmt.Sprintf("Install error and retried max(%d) times already.", prometheusMaxRetryCount)
			err := c.persistUpdate(ctx, prometheus)
			if err != nil {
				log.Error("Update prometheus error.")
				return true, err
			}
			return true, nil
		}
		// Add the retry count will trigger reinitialize function from the persistent controller again.
		prometheus = prometheus.DeepCopy()
		prometheus.Status.Phase = v1.AddonPhaseReinitializing
		prometheus.Status.Reason = err.Error()
		prometheus.Status.LastReInitializingTimestamp = metav1.NewTime(time.Now())
		prometheus.Status.RetryCount++
		err = c.persistUpdate(ctx, prometheus)
		if err != nil {
			return true, err
		}
		return true, nil
	}
}

func (c *Controller) createPrometheusIfNeeded(ctx context.Context, key string, cachedPrometheus *cachedPrometheus, prometheus *v1.Prometheus) error {
	switch prometheus.Status.Phase {
	case v1.AddonPhaseInitializing:
		log.Info("Prometheus will be created", log.String("prome", key))
		err := c.installPrometheus(ctx, prometheus)
		if err == nil {
			log.Info("Prometheus created success", log.String("prome", key))
			prometheus = prometheus.DeepCopy()
			prometheus.Status.Phase = v1.AddonPhaseChecking
			prometheus.Status.Reason = ""
			prometheus.Status.RetryCount = 0
			return c.persistUpdate(ctx, prometheus)
		}
		log.Error(fmt.Sprintf("Prometheus created failed: %v", err), log.String("prome", key))
		prometheus = prometheus.DeepCopy()
		prometheus.Status.Phase = v1.AddonPhaseReinitializing
		prometheus.Status.Reason = err.Error()
		prometheus.Status.RetryCount = 1
		return c.persistUpdate(ctx, prometheus)
	case v1.AddonPhaseReinitializing:
		log.Info("Prometheus entry Reinitializing", log.String("prome", key))
		var interval = time.Since(prometheus.Status.LastReInitializingTimestamp.Time)
		var waitTime time.Duration
		if interval >= prometheusTimeOut {
			waitTime = time.Duration(1)
		} else {
			waitTime = prometheusTimeOut - interval
		}
		go wait.Poll(waitTime, prometheusTimeOut, c.prometheusReinitialize(ctx, key, cachedPrometheus, prometheus))
	case v1.AddonPhaseChecking:
		log.Info("Prometheus entry Checking", log.String("prome", key))
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, prometheus)
			initDelay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.checking.Delete(key)
				wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkPrometheusStatus(ctx, prometheus, key, initDelay))
			}()
		}
	case v1.AddonPhaseRunning:
		log.Info("Prometheus entry Running", log.String("prome", key))
		UpdateMetricPrometheusStatusFail(prometheus.Spec.TenantID, prometheus.Spec.ClusterName, prometheus.Name, false)
		if needUpgrade(prometheus) {
			c.health.Delete(key)
			prometheus = prometheus.DeepCopy()
			prometheus.Status.Phase = v1.AddonPhaseUpgrading
			prometheus.Status.Reason = ""
			prometheus.Status.RetryCount = 0
			return c.persistUpdate(ctx, prometheus)
		}

		if _, ok := c.health.Load(key); !ok {
			c.health.Store(key, prometheus)
			go wait.PollImmediateUntil(5*time.Minute, c.watchPrometheusHealth(ctx, key), c.stopCh)
		}
	case v1.AddonPhaseUpgrading:
		log.Info("Prometheus entry upgrading", log.String("prome", key))
		if _, ok := c.upgrading.Load(key); !ok {
			c.upgrading.Store(key, prometheus)
			delay := time.Now().Add(5 * time.Minute)
			go func() {
				defer c.upgrading.Delete(key)
				wait.PollImmediate(5*time.Second, 5*time.Minute, c.checkPrometheusUpgrade(ctx, prometheus, key, delay))
			}()
		}
	case v1.AddonPhaseFailed:
		log.Info("Prometheus entry fail", log.String("prome", key))
		c.upgrading.Delete(key)
		c.health.Delete(key)

		UpdateMetricPrometheusStatusFail(prometheus.Spec.TenantID, prometheus.Spec.ClusterName, prometheus.Name, true)

		// should try check when prometheus recover again
		log.Info("Prometheus try checking after fail", log.String("prome", key))
		if _, ok := c.checking.Load(key); !ok {
			c.checking.Store(key, prometheus)
			go func() {
				defer c.checking.Delete(key)
				wait.PollImmediateUntil(60*time.Second, c.checkPrometheusStatus(ctx, prometheus, key, time.Time{}), c.stopCh)
			}()
		}
	}
	return nil
}

func (c *Controller) installPrometheus(ctx context.Context, prometheus *v1.Prometheus) error {
	if c.notifyAPIAddress == "" {
		return fmt.Errorf("empty notify api address, check if notify api exists")
	}
	components := images.Get(prometheus.Spec.Version)
	prometheus.Status.Version = prometheus.Spec.Version
	if prometheus.Status.SubVersion == nil {
		prometheus.Status.SubVersion = make(map[string]string)
	}

	cluster, err := c.platformClient.Clusters().Get(ctx, prometheus.Spec.ClusterName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get cluster failed: %v", err)
	}
	kubeClient, err := addon.BuildExternalClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return fmt.Errorf("get kubeClient failed: %v", err)
	}

	crdClient, err := addon.BuildExternalExtensionClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return fmt.Errorf("get crdClient failed: %v", err)
	}

	kaClient, err := addon.BuildKubeAggregatorClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return fmt.Errorf("get kaClient failed: %v", err)
	}

	mclient, err := addon.BuildExternalMonitoringClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return fmt.Errorf("get mclient failed: %v", err)
	}

	// Set remote write address
	var remoteWrites []string
	switch c.remoteType {
	case "influxdb":
		remoteWrites, err = c.initInfluxdb(cluster.Name)
		if err != nil {
			return fmt.Errorf("init Influxdb failed: %v", err)
		}
	case "elasticsearch":
		remoteWrites, err = c.initESAdapter(ctx, kubeClient, components)
		if err != nil {
			return fmt.Errorf("init ESAdapter failed: %v", err)
		}
		prometheus.Status.SubVersion[PrometheusBeatService] = components.PrometheusBeatWorkLoad.Tag
	case "thanos":
		remoteWrites, err = c.initThanos()
		if err != nil {
			return fmt.Errorf("init Thanos failed: %v", err)
		}
	default:
	}

	if len(prometheus.Spec.RemoteAddress.WriteAddr) > 0 {
		remoteWrites = prometheus.Spec.RemoteAddress.WriteAddr
	}
	// For remote read, just set from spec
	var remoteReads []string
	if len(prometheus.Spec.RemoteAddress.ReadAddr) > 0 {
		remoteReads = prometheus.Spec.RemoteAddress.ReadAddr
	}

	crds := getCRDs()
	for _, crd := range crds {
		if apiclient.ClusterVersionIsBefore116(kubeClient) {
			crdObj := apiextensionsv1beta1.CustomResourceDefinition{}
			err := json.Unmarshal([]byte(crd), &crdObj)
			if err != nil {
				return fmt.Errorf("Unmarshal crd failed: %v", err)
			}
			crdObj.TypeMeta.APIVersion = "apiextensions.k8s.io/v1beta1"
			_, err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(ctx, &crdObj, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("create crd failed %v", err)
				return err
			}
		} else {
			crdObj := apiextensionsv1.CustomResourceDefinition{}
			err := json.Unmarshal([]byte(crd), &crdObj)
			if err != nil {
				return fmt.Errorf("Unmarshal crd failed: %v", err)
			}

			_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, &crdObj, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("create crd failed %v", err)
				return err
			}
		}
	}

	log.Infof("Start to create prometheus-operator")
	// Service prometheus-operator
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, servicePrometheusOperator(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-operator service failed: %v", err)
	}
	// ServiceAccount for prometheus-operator
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, serviceAccountPrometheusOperator(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-operator ServiceAccount failed: %v", err)
	}
	// ClusterRole for prometheus-operator
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, clusterRolePrometheusOperator(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-operator ClusterRole failed: %v", err)
	}
	// ClusterRoleBinding prometheus-operator
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBindingPrometheusOperator(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-operator ClusterRoleBinding failed: %v", err)
	}
	// Deployment for prometheus-operator
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(ctx, deployPrometheusOperatorApps(components, prometheus), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-operator Deployment failed: %v", err)
	}

	prometheus.Status.SubVersion[PrometheusOperatorService] = components.PrometheusOperatorService.Tag

	extensionsAPIGroup := apiclient.ClusterVersionIsBefore19(kubeClient)

	// get notify webhook address
	webhookAddr, err := c.getWebhookAddr(ctx, prometheus)
	if err != nil {
		return fmt.Errorf("get webhook address failed: %v", err)
	}

	log.Infof("Start to create alertmanager")
	// secret for alertmanager
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(ctx, createSecretForAlertmanager(webhookAddr, prometheus.Spec.AlertRepeatInterval), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create alertmanager secret failed: %v", err)
	}

	// ServiceAccount for alertmanager
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, serviceAccountAlertmanager(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create alertmanager ServiceAccount failed: %v", err)
	}

	// Service for alertmanager
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, createServiceForAlerterManager(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create alertmanager Service failed: %v", err)
	}

	// Crd alertmanager instance
	if _, err := mclient.MonitoringV1().Alertmanagers(metav1.NamespaceSystem).Create(ctx, createAlertManagerCRD(components, prometheus), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create alertmanager Crd instance failed: %v", err)
	}

	prometheus.Status.SubVersion[AlertManagerService] = components.AlertManagerService.Tag

	log.Infof("Start to create prometheus")

	// Secret for prometheus-etcd
	credential, err := addon.GetClusterCredentialV1(ctx, c.platformClient, cluster)
	if err != nil {
		return fmt.Errorf("get credential failed: %v", err)
	}
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(ctx, secretETCDPrometheus(credential), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-etcd Secret failed: %v", err)
	}
	// Service Prometheus
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, servicePrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus Service failed: %v", err)
	}
	// Secret for prometheus
	if _, err := kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(ctx, createSecretForPrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus Secret failed: %v", err)
	}
	// ServiceAccount for prometheus
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, serviceAccountPrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus ServiceAccount failed: %v", err)
	}
	// ClusterRole for prometheus
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, clusterRolePrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus ClusterRole failed: %v", err)
	}
	// ClusterRoleBinding Prometheus
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBindingPrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus ClusterRoleBinding failed: %v", err)
	}
	// prometheus rule record
	if _, err := mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Create(ctx, recordsForPrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus rule record failed: %v", err)
	}
	// prometheus rule alert, empty for now, edit by tke-monitor
	if _, err := mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Create(ctx, alertsForPrometheus(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus rule alert failed: %v", err)
	}
	// Crd prometheus instance
	if _, err := mclient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).Create(ctx, createPrometheusCRD(components, prometheus, cluster, remoteWrites, remoteReads, c.remoteType), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus crd instance failed: %v", err)
	}
	prometheus.Status.SubVersion[PrometheusService] = components.PrometheusService.Tag

	log.Infof("Start to create node-exporter")
	// DaemonSet for node-exporter
	if _, err := kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(ctx, createDaemonSetForNodeExporter(components), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create node-exporter failed: %v", err)
	}
	prometheus.Status.SubVersion[nodeExporterService] = components.NodeExporterService.Tag

	log.Infof("Start to create kube-state-metrics")
	// Service for kube-state-metrics
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, createServiceForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics Service failed: %v", err)
	}
	// ServiceAccount for kube-state-metrics
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, createServiceAccountForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics ServiceAccount failed: %v", err)
	}
	// ClusterRole for kube-state-metrics
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, createClusterRoleForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics ClusterRole failed: %v", err)
	}
	// ClusterRoleBinding for kube-state-metrics
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, createClusterRoleBindingForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics ClusterRoleBinding failed: %v", err)
	}
	// Role for kube-state-metrics
	if _, err := kubeClient.RbacV1().Roles(metav1.NamespaceSystem).Create(ctx, createRoleForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics Role failed: %v", err)
	}
	// RoleBinding for kube-state-metrics
	if _, err := kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Create(ctx, createRoleBingdingForMetrics(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create kube-state-metrics RoleBinding failed: %v", err)
	}
	// Deployment for kube-state-metrics
	if extensionsAPIGroup {
		if _, err := kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Create(ctx, createExtensionDeploymentForMetrics(components, prometheus), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create kube-state-metrics Deployment failed: %v", err)
		}
	} else {
		if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(ctx, createAppsDeploymentForMetrics(components, prometheus), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create kube-state-metrics Deployment failed: %v", err)
		}
	}
	prometheus.Status.SubVersion[kubeStateService] = components.KubeStateService.Tag

	log.Infof("Start to create promtheus-adapter")
	// Service for prometheus-adapter
	if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, createServiceForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter Service failed: %v", err)
	}
	// ServiceAccount for prometheus-adapter
	if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, createServiceAccountForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter ServiceAccount failed: %v", err)
	}
	// ResourceReaderClusterRole for prometheus-adapter
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, createResourceReaderClusterRoleForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter ResourceReaderClusterRole failed: %v", err)
	}
	// ClusterRole for prometheus-adapter
	if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, createClusterRoleForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter ClusterRole failed: %v", err)
	}
	// AuthDelegatorClusterRoleBinding for prometheus-adapter
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, createAuthDelegatorClusterRoleBindingForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter AuthDelegatorClusterRoleBinding failed: %v", err)
	}
	// ResourceReaderClusterRoleBinding for prometheus-adapter
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, createResourceReaderClusterRoleBindingForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter ResourceReaderClusterRoleBinding failed: %v", err)
	}
	// HPAClusterRoleBinding for prometheus-adapter
	if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, createHPAClusterRoleBindingForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter HPAClusterRoleBinding failed: %v", err)
	}
	// AuthReaderRoleBinding for prometheus-adapter
	if _, err := kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Create(ctx, createAuthReaderRoleBingdingForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter AuthReaderRoleBinding failed: %v", err)
	}
	// APIService for prometheus-adapter
	for _, apiservice := range createAPIServiceForPrometheusAdapter() {
		if _, err := kaClient.ApiregistrationV1().APIServices().Create(ctx, apiservice, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create prometheus-adapter APIService failed: %v", err)
		}
	}
	// ConfigMap for prometheus-adapter
	if _, err := kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(ctx, createConfigMapForPrometheusAdapter(), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter ConfigMap failed: %v", err)
	}
	// Deployment for prometheus-adapter
	if _, err := kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(ctx, createAppsDeploymentForPrometheusAdapter(components, prometheus), metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("create prometheus-adapter Deployment failed: %v", err)
	}

	prometheus.Status.SubVersion[prometheusAdapterService] = components.PrometheusAdapter.Tag

	if prometheus.Spec.WithNPD {
		// ServiceAccount for node-problem-detector
		if _, err := kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Create(ctx, createServiceAccountForNPD(), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create node-problem-detector ServiceAccount failed: %v", err)
		}
		// ClusterRole for node-problem-detector
		if _, err := kubeClient.RbacV1().ClusterRoles().Create(ctx, createClusterRoleForNPD(), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create node-problem-detector ClusterRole failed: %v", err)
		}
		// ClusterRoleBinding for node-problem-detector
		if _, err := kubeClient.RbacV1().ClusterRoleBindings().Create(ctx, createClusterRoleBindingForNPD(), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create node-problem-detector ClusterRoleBinding failed: %v", err)
		}
		// DaemonSet for node-problem-detector
		if _, err := kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Create(ctx, createDaemonSetForNPD(components), metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("create node-problem-detector DaemonSet failed: %v", err)
		}
		prometheus.Status.SubVersion[nodeProblemDetectorWorkload] = components.NodeProblemDetector.Tag
	}
	return nil
}

func (c *Controller) getWebhookAddr(ctx context.Context, prometheus *v1.Prometheus) (webhookAddr string, err error) {
	if prometheus.Spec.NotifyWebhook != "" {
		webhookAddr = prometheus.Spec.NotifyWebhook
	} else {
		// use notify api address directly in global cluster
		webhookAddr = c.notifyAPIAddress + "/webhook"
		// use tke-gateway as proxy in non-global cluster
		if prometheus.Spec.ClusterName != "global" {
			globalCluster, err := c.platformClient.Clusters().Get(ctx, "global", metav1.GetOptions{})
			if err != nil {
				return "", fmt.Errorf("get global cluster failed: %v", err)
			}
			gatewayAddr := globalCluster.Spec.Machines[0].IP
			if globalCluster.Spec.Features.HA != nil {
				if globalCluster.Spec.Features.HA.TKEHA != nil {
					gatewayAddr = globalCluster.Spec.Features.HA.TKEHA.VIP
				}
				if globalCluster.Spec.Features.HA.ThirdPartyHA != nil {
					gatewayAddr = globalCluster.Spec.Features.HA.ThirdPartyHA.VIP
				}
			}
			webhookAddr = utilhttp.MakeEndpoint("https", gatewayAddr,
				443, "/webhook")
		}
	}
	return webhookAddr, nil
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
			APIVersion: "rbac.authorization.k8s.io/v1",
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

func deployPrometheusOperatorApps(components images.Components, prometheus *v1.Prometheus) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
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
					Containers: []corev1.Container{
						{
							Name:  prometheusOperatorWorkLoad,
							Image: components.PrometheusOperatorService.FullName(),
							Args: []string{
								"--kubelet-service=kube-system/kubelet",
								"--logtostderr=true",
								"--config-reloader-image=" + components.ConfigMapReloadWorkLoad.FullName(),
								"--prometheus-config-reloader=" + components.PrometheusConfigReloaderWorkload.FullName(),
							},
							//Command:         []string{"tail", "-f"},
							//Command: []string{"/bin/sh", "-c", "./prometheus --storage.tsdb.retention=1h --storage.tsdb.path=/data --web.enable-lifecycle --config.file=config/prometheus.yml"},
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
	if prometheus.Spec.RunOnMaster {
		deploy.Spec.Template.Spec.Affinity = &corev1.Affinity{
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
		}
	}
	return deploy
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
			Name:        PrometheusService,
			Namespace:   metav1.NamespaceSystem,
			Labels:      map[string]string{"kubernetes.io/name": "Prometheus", "addonmanager.kubernetes.io/mode": "Reconcile", "kubernetes.io/cluster-service": "true"},
			Annotations: map[string]string{"prometheus.io/scrape": "true", "prometheus.io/port": "9090"},
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
			APIVersion: "rbac.authorization.k8s.io/v1",
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

func createPrometheusCRD(components images.Components, prometheus *v1.Prometheus, cluster *platformv1.Cluster, remoteWrites, remoteReads []string, remoteType string) *monitoringv1.Prometheus {
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
		switch remoteType {
		case "influxdb":
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
						Regex:        "istio_(.*)|envoy_(.*)|pilot_(.*)|k8s_(.*)|apiserver_(.*)|kube_pod_labels|kube_node_labels|kube_namespace_labels|etcd_(.*)|grpc_(.*)|process_(.*)|scheduler_(.*)|workqueue_(.*)|rest_client_requests_(.*)|go_goroutines|kubelet_(.*)|volume_manager_(.*)|storage_operation_(.*)|coredns_(.*)|up|vm_(.*)",
						Action:       "keep",
					},
				}
				rw.QueueConfig = &monitoringv1.QueueConfig{
					Capacity:          10000,
					MinShards:         1,
					MaxShards:         24,
					MaxSamplesPerSend: 500,
					BatchSendDeadline: "30s",
				}
			}
		case "elasticsearch":
			rw.WriteRelabelConfigs = []monitoringv1.RelabelConfig{
				{
					SourceLabels: []string{"__name__"},
					Regex:        "istio_(.*)|envoy_(.*)|pilot_(.*)|project_(.*)|apiserver_(.*)|k8s_(.*)|kube_pod_labels|kube_node_labels|kube_namespace_labels|etcd_(.*)|grpc_(.*)|process_(.*)|scheduler_(.*)|workqueue_(.*)|rest_client_requests_(.*)|go_goroutines|kubelet_(.*)|volume_manager_(.*)|storage_operation_(.*)|coredns_(.*)|up|vm_(.*)",
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
		case "thanos":
			rw.WriteRelabelConfigs = []monitoringv1.RelabelConfig{
				{
					SourceLabels: []string{"__name__"},
					Regex:        "istio_(.*)|envoy_(.*)|pilot_(.*)|project_(.*)|apiserver_(.*)|k8s_(.*)|kube_pod_labels|kube_node_labels|kube_namespace_labels|etcd_(.*)|grpc_(.*)|process_(.*)|scheduler_(.*)|workqueue_(.*)|rest_client_requests_(.*)|go_goroutines|kubelet_(.*)|volume_manager_(.*)|storage_operation_(.*)|coredns_(.*)|up|vm_(.*)",
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
		default:
			log.Warnf("unknown remoteType %s", remoteType)
		}

		remoteWriteSpecs = append(remoteWriteSpecs, rw)
	}
	monitorV1Prometheus := &monitoringv1.Prometheus{
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
			PodMetadata: &monitoringv1.EmbeddedObjectMetadata{
				Annotations: map[string]string{
					"prometheus.io/scrape": "false",
					"prometheus.io/port":   "9090",
				},
			},
			ExternalLabels:     map[string]string{"cluster_id": cluster.Name, "tenant_id": prometheus.Spec.TenantID, "cluster_display_name": cluster.Spec.DisplayName},
			ScrapeInterval:     "60s",
			RemoteRead:         remoteReadSpecs,
			RemoteWrite:        remoteWriteSpecs,
			Retention:          "30m",
			RetentionSize:      "5GB",
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
					corev1.ResourceCPU:    *resource.NewMilliQuantity(4000, resource.DecimalSI),
					corev1.ResourceMemory: *resource.NewQuantity(8*1024*1024*1024, resource.BinarySI),
				},
			},
			Tolerations: []corev1.Toleration{
				{
					Key:      "node-role.kubernetes.io/master",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
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
			Version:                         components.PrometheusService.Tag,
		},
	}
	if prometheus.Spec.RunOnMaster {
		monitorV1Prometheus.Spec.Affinity = &corev1.Affinity{
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
		}
	}
	if len(prometheus.Spec.Resources.Requests) > 0 {
		resources := corev1.ResourceRequirements{
			Limits:   corev1.ResourceList{},
			Requests: corev1.ResourceList{},
		}
		for k, v := range prometheus.Spec.Resources.Limits {
			resources.Limits[corev1.ResourceName(k)] = v
		}
		for k, v := range prometheus.Spec.Resources.Requests {
			resources.Requests[corev1.ResourceName(k)] = v
		}
		monitorV1Prometheus.Spec.Resources = resources
	}
	return monitorV1Prometheus
}

func createSecretForPrometheus() *corev1.Secret {
	config := scrapeConfigForPrometheus()

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusSecret,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			prometheusConfigName: []byte(config),
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

func createSecretForAlertmanager(webhookAddr string, repeatInterval string) *corev1.Secret {
	if repeatInterval == "" {
		repeatInterval = "1200s"
	}
	config := configForAlertManager(webhookAddr, repeatInterval)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alertManagerSecret,
			Namespace: metav1.NamespaceSystem,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"alertmanager.yaml": []byte(config),
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

func createAlertManagerCRD(components images.Components, prometheus *v1.Prometheus) *monitoringv1.Alertmanager {
	monitorV1Alertmanager := &monitoringv1.Alertmanager{
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
			PodMetadata: &monitoringv1.EmbeddedObjectMetadata{
				Annotations: map[string]string{
					"prometheus.io/scrape": "false",
					"prometheus.io/port":   "9093",
				},
			},
			Tolerations: []corev1.Toleration{
				{
					Key:      "node-role.kubernetes.io/master",
					Operator: corev1.TolerationOpExists,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
			BaseImage: containerregistryutil.GetImagePrefix(alertManagerImagePath),
			Replicas:  controllerutil.Int32Ptr(1),
			SecurityContext: &corev1.PodSecurityContext{
				FSGroup:      controllerutil.Int64Ptr(2000),
				RunAsNonRoot: controllerutil.BoolPtr(true),
				RunAsUser:    controllerutil.Int64Ptr(1000),
			},
			ServiceAccountName: alertManagerServiceAccount,
			Version:            components.AlertManagerService.Tag,
		},
	}
	if prometheus.Spec.RunOnMaster {
		monitorV1Alertmanager.Spec.Affinity = &corev1.Affinity{
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
		}
	}
	return monitorV1Alertmanager
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
			Annotations: map[string]string{"prometheus.io/scrape": "true", "prometheus.io/port": "9093"},
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

func createDaemonSetForNodeExporter(components images.Components) *appsv1.DaemonSet {
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
							Image: components.NodeExporterService.FullName(),
							Args: []string{
								"--path.rootfs=/host",
								"--no-collector.arp",
								"--no-collector.bcache",
								"--no-collector.bonding",
								"--no-collector.buddyinfo",
								"--no-collector.conntrack",
								"--collector.cpu",
								"--no-collector.cpufreq",
								"--collector.diskstats",
								"--no-collector.drbd",
								"--no-collector.edac",
								"--no-collector.entropy",
								"--collector.filefd",
								"--collector.filesystem",
								"--no-collector.hwmon",
								"--no-collector.infiniband",
								"--no-collector.interrupts",
								"--no-collector.ipvs",
								"--no-collector.ksmd",
								"--collector.loadavg",
								"--no-collector.logind",
								"--no-collector.mdadm",
								"--collector.meminfo",
								"--no-collector.meminfo_numa",
								"--no-collector.mountstats",
								"--collector.netdev",
								"--no-collector.netstat",
								"--no-collector.netclass",
								"--no-collector.nfs",
								"--no-collector.nfsd",
								"--no-collector.pressure",
								"--no-collector.ntp",
								"--no-collector.qdisc",
								"--no-collector.runit",
								"--collector.sockstat",
								"--collector.stat",
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
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/host",
									Name:      "root",
									ReadOnly:  true,
								},
							},
						},
					},
					HostNetwork: true,
					HostPID:     true,
					Volumes: []corev1.Volume{
						{
							Name: "root",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
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
			APIVersion: "rbac.authorization.k8s.io/v1",
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
				Resources: []string{"daemonsets", "deployments", "replicasets", "ingresses"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"statefulsets", "daemonsets", "deployments", "replicasets"},
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
			{
				APIGroups: []string{"authentication.k8s.io"},
				Resources: []string{"tokenreviews"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{"authorization.k8s.io"},
				Resources: []string{"subjectaccessreviews"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"certificates.k8s.io"},
				Resources: []string{"certificatesigningrequests"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses", "volumeattachments"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Resources: []string{"mutatingwebhookconfigurations", "validatingwebhookconfigurations"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"networkpolicies"},
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
			APIVersion: "rbac.authorization.k8s.io/v1",
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

func createExtensionDeploymentForMetrics(components images.Components, prometheus *v1.Prometheus) *extensionsv1beta1.Deployment {
	deploy := &extensionsv1beta1.Deployment{
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
							Image: components.KubeStateService.FullName(),
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
	if prometheus.Spec.RunOnMaster {
		deploy.Spec.Template.Spec.Affinity = &corev1.Affinity{
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
		}
	}
	return deploy
}

func createAppsDeploymentForMetrics(components images.Components, prometheus *v1.Prometheus) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
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
							Image: components.KubeStateService.FullName(),
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
	if prometheus.Spec.RunOnMaster {
		deploy.Spec.Template.Spec.Affinity = &corev1.Affinity{
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
		}
	}
	return deploy
}

var selectorForPrometheusAdapter = metav1.LabelSelector{
	MatchLabels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": prometheusAdapterWorkLoad},
}

func createAuthDelegatorClusterRoleBindingForPrometheusAdapter() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusAdapterAuthDelegatorClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     systemAuthDelegatorClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      prometheusAdapterServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createAuthReaderRoleBingdingForPrometheusAdapter() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusAdapterAuthReaderRoleBinding,
			Namespace: metav1.NamespaceSystem,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "extension-apiserver-authentication-reader",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      prometheusAdapterServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createResourceReaderClusterRoleBindingForPrometheusAdapter() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusAdapterResourceReaderClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusAdapterResourceReaderClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      prometheusAdapterServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createResourceReaderClusterRoleForPrometheusAdapter() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusAdapterResourceReaderClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "nodes", "nodes/stats", "configmaps"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func createServiceAccountForPrometheusAdapter() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusAdapterServiceAccount,
			Namespace: metav1.NamespaceSystem,
		},
	}
}

func createServiceForPrometheusAdapter() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusAdapterService,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorForPrometheusAdapter.MatchLabels,
			Ports: []corev1.ServicePort{
				{Port: 443, TargetPort: intstr.FromInt(6443)},
			},
		},
	}
}

func createClusterRoleForPrometheusAdapter() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusAdapterClusterRole,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"custom.metrics.k8s.io", "external.metrics.k8s.io"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
}

func createHPAClusterRoleBindingForPrometheusAdapter() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: prometheusAdapterHPAClusterRoleBinding,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     prometheusAdapterClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "horizontal-pod-autoscaler",
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

func createConfigMapForPrometheusAdapter() *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusAdapterConfigMap,
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"config.yaml": configForPrometheusAdapter(),
		},
	}
	return cm
}

func createAPIServiceForPrometheusAdapter() []*apiregistrationv1.APIService {
	apiServices := []*apiregistrationv1.APIService{}
	customMetricsV1beta1 := &apiregistrationv1.APIService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIService",
			APIVersion: "apiregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "v1beta1.custom.metrics.k8s.io",
		},
		Spec: apiregistrationv1.APIServiceSpec{
			Service: &apiregistrationv1.ServiceReference{
				Namespace: metav1.NamespaceSystem,
				Name:      prometheusAdapterService,
			},
			Group:                 "custom.metrics.k8s.io",
			Version:               "v1beta1",
			InsecureSkipTLSVerify: true,
			GroupPriorityMinimum:  100,
			VersionPriority:       100,
		},
	}
	apiServices = append(apiServices, customMetricsV1beta1)
	customMetricsV1beta2 := &apiregistrationv1.APIService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIService",
			APIVersion: "apiregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "v1beta2.custom.metrics.k8s.io",
		},
		Spec: apiregistrationv1.APIServiceSpec{
			Service: &apiregistrationv1.ServiceReference{
				Namespace: metav1.NamespaceSystem,
				Name:      prometheusAdapterService,
			},
			Group:                 "custom.metrics.k8s.io",
			Version:               "v1beta2",
			InsecureSkipTLSVerify: true,
			GroupPriorityMinimum:  100,
			VersionPriority:       200,
		},
	}
	apiServices = append(apiServices, customMetricsV1beta2)
	externalMetricsV1beta1 := &apiregistrationv1.APIService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIService",
			APIVersion: "apiregistration.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "v1beta1.external.metrics.k8s.io",
		},
		Spec: apiregistrationv1.APIServiceSpec{
			Service: &apiregistrationv1.ServiceReference{
				Namespace: metav1.NamespaceSystem,
				Name:      prometheusAdapterService,
			},
			Group:                 "external.metrics.k8s.io",
			Version:               "v1beta1",
			InsecureSkipTLSVerify: true,
			GroupPriorityMinimum:  100,
			VersionPriority:       100,
		},
	}
	apiServices = append(apiServices, externalMetricsV1beta1)
	return apiServices
}

func createAppsDeploymentForPrometheusAdapter(components images.Components, prometheus *v1.Prometheus) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prometheusAdapterWorkLoad,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", specialLabelName: specialLabelValue, "k8s-app": prometheusAdapterWorkLoad},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: controllerutil.Int32Ptr(1),
			Selector: &selectorForPrometheusAdapter,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": prometheusAdapterWorkLoad},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: prometheusAdapterServiceAccount,
					Containers: []corev1.Container{
						{
							Name:  prometheusAdapterWorkLoad,
							Image: components.PrometheusAdapter.FullName(),
							Args: []string{
								"--secure-port=6443",
								"--logtostderr=true",
								"--prometheus-url=http://prometheus.kube-system.svc.cluster.local:9090/",
								"--metrics-relist-interval=1m",
								"--v=3",
								"--config=/etc/adapter/config.yaml",
								"--cert-dir=/tmp",
								"--requestheader-client-ca-file=/etc/kubernetes/pki/requestheader-client-ca-file",
								"--requestheader-allowed-names=front-proxy-client,admin",
								"--requestheader-extra-headers-prefix=X-Remote-Extra-",
								"--requestheader-group-headers=X-Remote-Group",
								"--requestheader-username-headers=X-Remote-User",
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 6443},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/etc/adapter/",
									Name:      "config",
									ReadOnly:  true,
								},
								{
									MountPath: "/tmp",
									Name:      "tmp-vol",
								},
								{
									MountPath: "/etc/kubernetes/pki/",
									Name:      "extension-apiserver-authentication",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "extension-apiserver-authentication",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "extension-apiserver-authentication",
									},
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: prometheusAdapterConfigMap,
									},
								},
							},
						},
						{
							Name: "tmp-vol",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
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
	return deploy
}

func createServiceAccountForNPD() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodeProblemDetectorServiceAccount,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
	}
}

func createClusterRoleForNPD() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRole",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   nodeProblemDetectorClusterRole,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "pods"},
				Verbs:     []string{"list", "get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes/status"},
				Verbs:     []string{"patch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch", "update"},
			},
		},
	}
}

func createClusterRoleBindingForNPD() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   nodeProblemDetectorClusterRoleBinding,
			Labels: map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile"},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     nodeProblemDetectorClusterRole,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      nodeProblemDetectorServiceAccount,
				Namespace: metav1.NamespaceSystem,
			},
		},
	}
}

var selectorForNPD = metav1.LabelSelector{
	MatchLabels: map[string]string{specialLabelName: specialLabelValue, "k8s-app": "node-problem-detector"},
}

func createDaemonSetForNPD(components images.Components) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodeProblemDetectorWorkload,
			Namespace: metav1.NamespaceSystem,
			Labels:    map[string]string{"kubernetes.io/cluster-service": "true", "addonmanager.kubernetes.io/mode": "Reconcile", specialLabelName: specialLabelValue, "k8s-app": "node-problem-detector"},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &selectorForNPD,
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{specialLabelName: specialLabelValue, "k8s-app": "node-problem-detector"},
					Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": "", "tke.prometheus.io/scrape": "true", "prometheus.io/scheme": "http", "prometheus.io/port": "20257"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  nodeProblemDetectorWorkload,
							Image: components.NodeProblemDetector.FullName(),
							Command: []string{
								"/node-problem-detector",
								"--logtostderr",
								"--prometheus-address=0.0.0.0",
								"--config.system-log-monitor=/config/kernel-monitor.json,/config/docker-monitor.json",
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPointer(true),
							},
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/var/log",
									Name:      "log",
									ReadOnly:  true,
								},
								{
									MountPath: "/dev/kmsg",
									Name:      "kmsg",
									ReadOnly:  true,
								},
								{
									MountPath: "/etc/localtime",
									Name:      "localtime",
									ReadOnly:  true,
								},
							},
						},
					},
					HostNetwork: true,
					HostPID:     true,
					Volumes: []corev1.Volume{
						{
							Name: "log",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
								},
							},
						},
						{
							Name: "kmsg",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev/kmsg",
								},
							},
						},
						{
							Name: "localtime",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/localtime",
								},
							},
						},
					},
					ServiceAccountName: nodeProblemDetectorServiceAccount,
					Tolerations: []corev1.Toleration{
						{
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Operator: corev1.TolerationOpExists,
							Key:      "CriticalAddonsOnly",
						},
					},
				},
			},
		},
	}
}

func (c *Controller) uninstallPrometheus(ctx context.Context, prometheus *v1.Prometheus, dropData bool) error {
	var err error
	var errs []error

	cluster, err := c.platformClient.Clusters().Get(ctx, prometheus.Spec.ClusterName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	kubeClient, err := addon.BuildExternalClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return err
	}

	kaClient, err := addon.BuildKubeAggregatorClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return fmt.Errorf("get kaClient failed: %v", err)
	}

	crdClient, err := addon.BuildExternalExtensionClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return err
	}

	mclient, err := addon.BuildExternalMonitoringClientSet(ctx, cluster, c.platformClient)
	if err != nil {
		return err
	}

	extensionsAPIGroup := apiclient.ClusterVersionIsBefore19(kubeClient)

	// delete prometheus
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(ctx, prometheusETCDSecret, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(ctx, prometheusSecret, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, prometheusClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, prometheusClusterRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, PrometheusService, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, prometheusServiceAccount, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	err = mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Delete(ctx, prometheusRuleRecord, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = mclient.MonitoringV1().PrometheusRules(metav1.NamespaceSystem).Delete(ctx, PrometheusRuleAlert, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	err = mclient.MonitoringV1().Prometheuses(metav1.NamespaceSystem).Delete(ctx, PrometheusCRDName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete alertmanager
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, AlertManagerService, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Delete(ctx, alertManagerSecret, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, alertManagerServiceAccount, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = mclient.MonitoringV1().Alertmanagers(metav1.NamespaceSystem).Delete(ctx, alertManagerCRDName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete node-exporter
	err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Delete(ctx, nodeExporterDaemonSet, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	// delete kube-state-metrics
	if extensionsAPIGroup {
		err = kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Delete(ctx, kubeStateWorkLoad, metav1.DeleteOptions{})
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
			err = controllerutil.DeleteReplicaSetApp(ctx, kubeClient, options)
			if err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(ctx, kubeStateWorkLoad, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, kubeStateClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, kubeStateClusterRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Delete(ctx, kubeStateRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().Roles(metav1.NamespaceSystem).Delete(ctx, kubeStateRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, kubeStateServiceAccount, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, kubeStateService, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	if prometheus.Spec.WithNPD {
		err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Delete(ctx, nodeProblemDetectorWorkload, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
		err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, nodeProblemDetectorClusterRoleBinding, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
		err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, nodeProblemDetectorClusterRole, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
		err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, nodeProblemDetectorServiceAccount, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
	}

	// delete prometheus-adapter
	err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(ctx, prometheusAdapterWorkLoad, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, prometheusAdapterService, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, prometheusAdapterServiceAccount, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, prometheusAdapterResourceReaderClusterRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, prometheusAdapterClusterRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, prometheusAdapterAuthDelegatorClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, prometheusAdapterResourceReaderClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, prometheusAdapterHPAClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().RoleBindings(metav1.NamespaceSystem).Delete(ctx, prometheusAdapterAuthReaderRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete(ctx, prometheusAdapterConfigMap, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	for _, apiservice := range createAPIServiceForPrometheusAdapter() {
		if err := kaClient.ApiregistrationV1().APIServices().Delete(ctx, apiservice.Name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
	}

	// delete prometheus-operator
	err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(ctx, prometheusOperatorWorkLoad, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(ctx, prometheusOperatorClusterRoleBinding, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.RbacV1().ClusterRoles().Delete(ctx, prometheusOperatorClusterRole, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, PrometheusOperatorService, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}
	err = kubeClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(ctx, prometheusOperatorServiceAccount, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		errs = append(errs, err)
	}

	crds := getCRDs()
	for _, crd := range crds {
		if apiclient.ClusterVersionIsBefore116(kubeClient) {
			crdObj := apiextensionsv1beta1.CustomResourceDefinition{}
			err := json.Unmarshal([]byte(crd), &crdObj)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			crdObj.TypeMeta.APIVersion = "apiextensions.k8s.io/v1beta1"
			err = crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(ctx, crdObj.Name, metav1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				errs = append(errs, err)
			}
		} else {
			crdObj := apiextensionsv1.CustomResourceDefinition{}
			err := json.Unmarshal([]byte(crd), &crdObj)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, crdObj.Name, metav1.DeleteOptions{})
			if err != nil && !errors.IsNotFound(err) {
				errs = append(errs, err)
			}
		}
	}

	// drop influxdb data, may take long time
	if dropData {
		if c.remoteType == "influxdb" {
			go c.dropInfluxdb(cluster.Name)
		}
	}

	if c.remoteType == "elasticsearch" {
		err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete(ctx, prometheusBeatWorkLoad, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
		err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Delete(ctx, PrometheusBeatConfigmap, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			errs = append(errs, err)
		}
		err = kubeClient.CoreV1().Services(metav1.NamespaceSystem).Delete(ctx, PrometheusBeatService, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
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

func (c *Controller) watchPrometheusHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start check prometheus in cluster health", log.String("prome", key))

		prometheus, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		cluster, err := c.platformClient.Clusters().Get(ctx, prometheus.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.health.Load(key); !ok {
			log.Info("Prometheus health check over", log.String("prome", key))
			return true, nil
		}
		kubeClient, err := addon.BuildExternalClientSet(ctx, cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", PrometheusService, PrometheusServicePort, `/-/healthy`, nil).DoRaw(ctx); err != nil {
			prometheus = prometheus.DeepCopy()
			prometheus.Status.Phase = v1.AddonPhaseFailed
			prometheus.Status.Reason = "Prometheus is not healthy."
			if err = c.persistUpdate(ctx, prometheus); err != nil {
				return false, err
			}
			return true, nil
		}
		log.Debug("Prometheus health is ok", log.String("prome", key))
		return false, nil
	}
}

func (c *Controller) checkPrometheusStatus(ctx context.Context, prometheus *v1.Prometheus, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to check prometheus status", log.String("prome", key))

		cluster, err := c.platformClient.Clusters().Get(ctx, prometheus.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.checking.Load(key); !ok {
			log.Info("Prometheus status checking over", log.String("prome", key))
			return true, nil
		}
		kubeClient, err := addon.BuildExternalClientSet(ctx, cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		prometheus, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		if _, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).ProxyGet("http", PrometheusService, PrometheusServicePort, `/-/healthy`, nil).DoRaw(ctx); err != nil {
			if !initDelay.IsZero() && time.Now().After(initDelay) {
				prometheus = prometheus.DeepCopy()
				prometheus.Status.Phase = v1.AddonPhaseFailed
				prometheus.Status.Reason = "Prometheus is not healthy."
				if err = c.persistUpdate(ctx, prometheus); err != nil {
					return false, err
				}
				return true, nil
			}
			log.Error("prometheus status has not healthy", log.String("prome", key), log.Err(err))
			return false, nil
		}
		prometheus = prometheus.DeepCopy()
		prometheus.Status.Phase = v1.AddonPhaseRunning
		prometheus.Status.Reason = ""
		if err = c.persistUpdate(ctx, prometheus); err != nil {
			return false, err
		}
		return true, nil
	}
}

func (c *Controller) checkPrometheusUpgrade(ctx context.Context, prometheus *v1.Prometheus, key string, initDelay time.Time) func() (bool, error) {
	return func() (bool, error) {
		log.Info("Start to upgrade prometheus", log.String("prome", key))

		cluster, err := c.platformClient.Clusters().Get(ctx, prometheus.Spec.ClusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			return false, err
		}
		if err != nil {
			return false, nil
		}
		if _, ok := c.upgrading.Load(key); !ok {
			log.Info("Prometheus upgrade over", log.String("prome", key))
			return true, nil
		}
		kubeClient, err := addon.BuildExternalClientSet(ctx, cluster, c.platformClient)
		if err != nil {
			return false, err
		}
		prometheus, err := c.lister.Get(key)
		if err != nil {
			return false, err
		}

		newComponents := images.Get(prometheus.Spec.Version)
		for name, version := range prometheus.Status.SubVersion {
			if image := newComponents.Get(name); image != nil {
				if image.Tag != version {
					patch := fmt.Sprintf(`[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`, image.FullName())
					err = upgradeVersion(ctx, kubeClient, name, patch)
					if err != nil {
						if time.Now().After(initDelay) {
							prometheus = prometheus.DeepCopy()
							prometheus.Status.Phase = v1.AddonPhaseUpgrading
							prometheus.Status.Reason = fmt.Sprintf("Upgrade timeout, %v", err)
							if err = c.persistUpdate(ctx, prometheus); err != nil {
								return false, err
							}
							return true, nil
						}
						return false, nil
					}
					prometheus.Status.SubVersion[name] = image.Tag
				}
			}
		}

		prometheus = prometheus.DeepCopy()
		prometheus.Status.Version = prometheus.Spec.Version
		prometheus.Status.Phase = v1.AddonPhaseChecking
		prometheus.Status.Reason = ""
		if err = c.persistUpdate(ctx, prometheus); err != nil {
			return false, err
		}

		return true, nil
	}
}

func needUpgrade(prom *v1.Prometheus) bool {
	if prom.Status.SubVersion == nil {
		log.Errorf("Nil component version when checking upgrade!")
		return false
	}

	if prom.Spec.SubVersion != nil {
		for com, version := range prom.Spec.SubVersion {
			if _, ok := prom.Status.SubVersion[com]; ok {
				if version != prom.Status.SubVersion[com] {
					return true
				}
			} else {
				log.Errorf("Not find %s version when checking upgrade!", com)
				return false
			}
		}
	}

	if prom.Spec.Version != prom.Status.Version {
		return true
	}

	return false
}

func upgradeVersion(ctx context.Context, kubeClient *kubernetes.Clientset, workLoad, patch string) error {
	var err error

	switch workLoad {
	case kubeStateWorkLoad, prometheusWorkLoad, AlertManagerWorkLoad:
		extensionsAPIGroup := apiclient.ClusterVersionIsBefore19(kubeClient)
		if extensionsAPIGroup {
			_, err = kubeClient.ExtensionsV1beta1().Deployments(metav1.NamespaceSystem).Patch(ctx, workLoad, types.JSONPatchType, []byte(patch), metav1.PatchOptions{})
		} else {
			_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Patch(ctx, workLoad, types.JSONPatchType, []byte(patch), metav1.PatchOptions{})
		}
	case nodeExporterDaemonSet:
		_, err = kubeClient.AppsV1().DaemonSets(metav1.NamespaceSystem).Patch(ctx, workLoad, types.JSONPatchType, []byte(patch), metav1.PatchOptions{})
	default:
		return fmt.Errorf("wrong workload: %s", workLoad)
	}

	return err
}

func (c *Controller) persistUpdate(ctx context.Context, prometheus *v1.Prometheus) error {
	var err error
	for i := 0; i < prometheusClientRetryCount; i++ {
		_, err = c.client.MonitorV1().Prometheuses().UpdateStatus(ctx, prometheus, metav1.UpdateOptions{})
		if err == nil {
			return nil
		}
		// If the object no longer exists, we don't want to recreate it. Just bail
		// out so that we can process the delete, which we should soon be receiving
		// if we haven't already.
		if errors.IsNotFound(err) {
			log.Info("Not persisting update to prometheus that no longer exists", log.String("prome", prometheus.Name), log.Err(err))
			return nil
		}
		if errors.IsConflict(err) {
			return fmt.Errorf("not persisting update to prometheus '%s' that has been changed since we received it: %v", prometheus.Name, err)
		}
		log.Warn(fmt.Sprintf("Failed to persist updated status of prometheus '%s/%s'", prometheus.Name, prometheus.Status.Phase), log.String("prome", prometheus.Name), log.Err(err))
		time.Sleep(prometheusClientRetryInterval)
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
	queryUser := influxApi.Query{
		Command:  cmdUser,
		Database: monitorutil.ProjectDatabaseName,
	}

	// create database/user, and grant privilege
	cmdDB := fmt.Sprintf("create database %s; create user %s with password '%s'; grant all on %s to %s; grant write on %s to %s",
		db, usr, passwd, db, usr, monitorutil.ProjectDatabaseName, usr)

	log.Debugf("Create influxdb table: %s", cmdDB)
	queryDB := influxApi.Query{
		Command:  cmdDB,
		Database: monitorutil.ProjectDatabaseName,
	}

	var queryStr []string
	for _, client := range c.remoteClients {
		_, _ = client.influxdb.client.Query(queryUser)

		resp, err := client.influxdb.client.Query(queryDB)
		if err != nil {
			return nil, err
		} else if resp.Error() != nil {
			return nil, resp.Error()
		}

		err = createOrAlterRetention(client, db, c.retentionDays)
		if err != nil {
			return nil, err
		}

		queryStr = append(queryStr, []string{
			fmt.Sprintf("%s/api/v1/prom/write?db=%s&u=%s&p=%s", client.influxdb.address, db, usr, passwd),
			fmt.Sprintf("%s/api/v1/prom/write?db=%s&u=%s&p=%s", client.influxdb.address, monitorutil.ProjectDatabaseName, usr, passwd),
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
	query := influxApi.Query{
		Command:  cmd,
		Database: monitorutil.ProjectDatabaseName,
	}

	// just continue when error
	for _, client := range c.remoteClients {
		resp, err := client.influxdb.client.Query(query)
		if err != nil {
			log.Errorf("Drop database(%s) for %s err: %v", dbName, client.influxdb.address, err)
		} else if resp.Error() != nil {
			log.Errorf("Drop database(%s) for %s err: %v", dbName, client.influxdb.address, resp.Error())
		}
	}

	return nil
}

func (c *Controller) initESAdapter(ctx context.Context, kubeClient *kubernetes.Clientset, components images.Components) ([]string, error) {
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
	for _, client := range c.remoteClients {
		hosts = append(hosts, client.es.URL)
		user = client.es.Username
		password = client.es.Password
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
	_, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).Create(ctx, svc, metav1.CreateOptions{})
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
	_, err = kubeClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(ctx, cm, metav1.CreateOptions{})
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
							Image: components.PrometheusBeatWorkLoad.FullName(),
							Args: []string{
								"-c",
								"config/prometheusbeat.yml",
								"-e",
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
	_, err = kubeClient.AppsV1().Deployments(metav1.NamespaceSystem).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return remoteWrites, err
	}
	remoteWrites = append(remoteWrites, fmt.Sprintf("http://%s:%d/prometheus", svc.Name, 8080))
	return remoteWrites, nil
}

func (c *Controller) initThanos() ([]string, error) {
	var remoteWrites []string
	for _, client := range c.remoteClients {
		remoteWrites = append(remoteWrites, fmt.Sprintf("%s/api/v1/receive", client.thanos.address))
	}

	return remoteWrites, nil
}

func boolPointer(value bool) *bool {
	return &value
}

func createOrAlterRetention(client remoteClient, db string, duration int) error {
	createRPSQL := influxApi.Query{
		Command:  fmt.Sprintf("create retention policy tke on %s duration %dd replication 1 default", db, duration),
		Database: monitorutil.ProjectDatabaseName,
	}
	alterRPSQL := influxApi.Query{
		Command:  fmt.Sprintf("alter retention policy tke on %s duration %dd", db, duration),
		Database: monitorutil.ProjectDatabaseName,
	}

	log.Debugf("create retention policy sql: %s", createRPSQL.Command)
	resp, err := client.influxdb.client.Query(createRPSQL) // create retention policy first
	if err != nil {
		log.Errorf("create retention policy 'tke' failed on db %s with client %s , err: %v", db, client.influxdb.address, err)
		return err
	}
	if resp.Error() != nil {
		if !strings.EqualFold(resp.Error().Error(), "retention policy already exists") { //other errors
			log.Errorf("create retention policy 'tke' failed on db %s with client %s , err: %v", db, client.influxdb.address, resp.Error())
			return resp.Error()
		}

		log.Debugf("alter retention policy sql: %s", alterRPSQL.Command)
		resp, err = client.influxdb.client.Query(alterRPSQL) //alter exist retetntion policy
		if err != nil {
			log.Errorf("alter retention policy 'tke' failed on db %s with client %s , err: %v", db, client.influxdb.address, err)
			return err
		}
		if resp.Error() != nil {
			log.Errorf("alter retention policy 'tke' failed on db %s with client %s , err: %v", db, client.influxdb.address, resp.Error())
			return resp.Error()
		}
		log.Infof("success to alter retention policy 'tke' on db %s with client %s", db, client.influxdb.address)
		return nil
	}
	log.Infof("success to create retention policy 'tke' on db %s with client %s", db, client.influxdb.address)
	return nil
}
