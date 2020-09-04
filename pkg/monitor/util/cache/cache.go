/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/monitor"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/monitor/util"
	platformutil "tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type updateComponent func(componentStatus *corev1.ComponentStatus, health *util.ComponentHealth)

type Component string

const (
	AllNamespaces = ""

	ClusterClientSet   = "ClusterClientSet"
	WorkloadCounter    = "WorkloadCounter"
	ResourceCounter    = "ResourceCounter"
	ClusterPhase       = "ClusterPhase"
	TenantID           = "TenantID"
	ClusterDisplayName = "ClusterDisplayName"
	ComponentHealth    = "ComponentHealth"

	TAppResourceName = "tapps"
	TAppGroupName    = "apps.tkestack.io"

	Scheduler         Component = "scheduler"
	ControllerManager Component = "controller-manager"
	Etcd              Component = "etcd-0"

	EtcdPrefix = "etcd-"

	FirstLoad    = int32(1)
	NotFirstLoad = int32(0)
)

var (
	TAppResource              = schema.GroupVersionResource{Group: TAppGroupName, Version: "v1", Resource: TAppResourceName}
	UpdateComponentStatusFunc = map[Component]updateComponent{
		Scheduler: func(componentStatus *corev1.ComponentStatus, health *util.ComponentHealth) {
			health.Scheduler = isHealthy(componentStatus)
		},
		ControllerManager: func(componentStatus *corev1.ComponentStatus, health *util.ComponentHealth) {
			health.ControllerManager = isHealthy(componentStatus)
		},
		Etcd: func(componentStatus *corev1.ComponentStatus, health *util.ComponentHealth) {
			health.Etcd = isHealthy(componentStatus)
		},
	}
)

type Cacher interface {
	Reload()
	GetClusterOverviewResult(clusters []*platformv1.Cluster) *monitor.ClusterOverviewResult
}

type cacher struct {
	sync.RWMutex
	platformClient platformversionedclient.PlatformV1Interface

	// businessClient is not required, and needed to determine if it's nil
	businessClient      businessversionedclient.BusinessV1Interface
	clusterStatisticSet util.ClusterStatisticSet
	clusterClientSets   util.ClusterClientSets
	dynamicClients      util.DynamicClientSet
	clusters            util.ClusterSet
	credentials         util.ClusterCredentialSet
	clusterAbnormal     int

	firstLoad int32
}

func (c *cacher) Reload() {
	c.Lock()
	defer c.Unlock()

	c.getClusters()
	c.getProjects()

	if c.firstLoad == FirstLoad {
		atomic.StoreInt32(&c.firstLoad, NotFirstLoad)
	}
}

func (c *cacher) getClusters() {
	c.getDynamicClients()

	if clusters, err := c.platformClient.Clusters().List(context.Background(), metav1.ListOptions{}); err == nil {
		wg := sync.WaitGroup{}
		wg.Add(len(clusters.Items))
		syncMap := sync.Map{}
		finished := int32(0)
		allTask := len(clusters.Items)
		started := time.Now()
		for i := range clusters.Items {
			if clusters.Items[i].Status.Phase == platformv1.ClusterFailed {
				c.clusterAbnormal++
			}
			go func(cls platformv1.Cluster) {
				defer func() {
					defer wg.Done()
					atomic.AddInt32(&finished, 1)
					log.Debugf("cacher has finished reloading (%d/%d) clusters, cluster: %s, cost: %v seconds",
						finished, allTask, cls.GetName(), time.Since(started).Seconds())
				}()
				clusterID := cls.GetName()
				clusterDisplayName := cls.Spec.DisplayName
				tenantID := cls.Spec.TenantID
				if cls.Status.Phase != platformv1.ClusterRunning {
					syncMap.Store(clusterID, map[string]interface{}{
						ClusterDisplayName: clusterDisplayName,
						ClusterPhase:       string(cls.Status.Phase),
						TenantID:           tenantID,
					})
					return
				}
				clientSet, err := platformutil.BuildExternalClientSet(context.Background(), &cls, c.platformClient)
				if err != nil {
					log.Error("create clientSet of cluster failed", log.Any("cluster", clusterID), log.Err(err))
					return
				}
				workloadCounter := c.getWorkloadCounter(clusterID, clientSet)
				resourceCounter := &util.ResourceCounter{}
				c.getNodes(clusterID, clientSet, resourceCounter)
				c.getPods(clusterID, clientSet, resourceCounter)
				if metricServerClientSet, err := c.getMetricServerClientSet(context.Background(), &cls); err == nil && metricServerClientSet != nil {
					c.getNodeMetrics(clusterID, metricServerClientSet, resourceCounter)
				}
				calResourceRate(resourceCounter)
				health := &util.ComponentHealth{}
				c.getComponentStatuses(clusterID, clientSet, health)
				syncMap.Store(clusterID, map[string]interface{}{
					ClusterClientSet:   clientSet,
					WorkloadCounter:    workloadCounter,
					ResourceCounter:    resourceCounter,
					ClusterPhase:       string(cls.Status.Phase),
					ClusterDisplayName: clusterDisplayName,
					TenantID:           tenantID,
					ComponentHealth:    health,
				})
			}(clusters.Items[i])
		}

		wg.Wait()

		log.Debugf("finish reloading all clusters, cost: %v seconds", time.Since(started).Seconds())

		c.clusterStatisticSet = make(util.ClusterStatisticSet)
		c.clusterAbnormal = 0
		syncMap.Range(func(key, value interface{}) bool {
			clusterID := key.(string)
			val := value.(map[string]interface{})
			clusterDisplayName := val[ClusterDisplayName].(string)
			tenantID := val[TenantID].(string)
			clusterPhase := val[ClusterPhase].(string)
			if clusterPhase == string(platformv1.ClusterRunning) {
				clusterClientSet := val[ClusterClientSet].(*kubernetes.Clientset)
				workloadCounter := val[WorkloadCounter].(*util.WorkloadCounter)
				resourceCounter := val[ResourceCounter].(*util.ResourceCounter)
				health := val[ComponentHealth].(*util.ComponentHealth)
				c.clusterClientSets[clusterID] = clusterClientSet
				c.clusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
					ClusterID:                clusterID,
					ClusterDisplayName:       clusterDisplayName,
					TenantID:                 tenantID,
					ClusterPhase:             clusterPhase,
					NodeCount:                int32(resourceCounter.NodeTotal),
					NodeAbnormal:             int32(resourceCounter.NodeAbnormal),
					WorkloadCount:            int32(workloadCounter.Total()),
					WorkloadAbnormal:         0,
					HasMetricServer:          resourceCounter.HasMetricServer,
					CPUUsed:                  resourceCounter.CPUUsed,
					CPURequest:               resourceCounter.CPURequest,
					CPULimit:                 resourceCounter.CPULimit,
					CPUAllocatable:           resourceCounter.CPUAllocatable,
					CPUCapacity:              resourceCounter.CPUCapacity,
					CPURequestRate:           transPercent(resourceCounter.CPURequestRate),
					CPUAllocatableRate:       transPercent(resourceCounter.CPUAllocatableRate),
					CPUUsage:                 transPercent(resourceCounter.CPUUsage),
					MemUsed:                  resourceCounter.MemUsed,
					MemRequest:               resourceCounter.MemRequest,
					MemLimit:                 resourceCounter.MemLimit,
					MemAllocatable:           resourceCounter.MemAllocatable,
					MemCapacity:              resourceCounter.MemCapacity,
					MemRequestRate:           transPercent(resourceCounter.MemRequestRate),
					MemAllocatableRate:       transPercent(resourceCounter.MemAllocatableRate),
					MemUsage:                 transPercent(resourceCounter.MemUsage),
					PodCount:                 int32(resourceCounter.PodCount),
					SchedulerHealthy:         health.Scheduler,
					ControllerManagerHealthy: health.ControllerManager,
					EtcdHealthy:              health.Etcd,
				}
			} else {
				c.clusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
					ClusterID:          clusterID,
					ClusterDisplayName: clusterDisplayName,
					ClusterPhase:       clusterPhase,
					TenantID:           tenantID,
				}
			}
			return true
		})

		log.Debugf("finish reloading all results, cost %+v seconds", time.Since(started).Seconds())
	}
}

func (c *cacher) getMetricServerClientSet(ctx context.Context, cls *platformv1.Cluster) (*metricsv.Clientset, error) {
	cc, err := platformutil.GetClusterCredentialV1(ctx, c.platformClient, cls)
	if err != nil {
		log.Error("query cluster credential failed", log.Any("cluster", cls.GetName()), log.Err(err))
		return nil, err
	}

	restConfig, err := platformutil.GetExternalRestConfig(cls, cc)
	if err != nil {
		log.Error("get rest config failed", log.Any("cluster", cls.GetName()), log.Err(err))
		return nil, err
	}

	return metricsv.NewForConfig(restConfig)
}

// TODO
func (c *cacher) getProjects() {
}

func (c *cacher) GetClusterOverviewResult(clusters []*platformv1.Cluster) *monitor.ClusterOverviewResult {
	if atomic.LoadInt32(&c.firstLoad) == FirstLoad {
		c.RLock()
		defer c.RUnlock()
	}

	clusterStatistics := make([]*monitor.ClusterStatistic, 0)
	result := &monitor.ClusterOverviewResult{}
	result.ClusterCount = int32(len(clusters))
	result.ClusterAbnormal = int32(c.clusterAbnormal)
	result.NodeAbnormal = 0
	result.WorkloadAbnormal = 0
	for i := 0; i < len(clusters); i++ {
		cls := clusters[i]
		if clusterStatistic, ok := c.clusterStatisticSet[cls.GetName()]; ok {
			if clusterStatistic.ClusterDisplayName != cls.Spec.DisplayName && len(cls.Spec.DisplayName) > 0 {
				clusterStatistic.ClusterDisplayName = cls.Spec.DisplayName
			}
			result.NodeCount += clusterStatistic.NodeCount
			result.NodeAbnormal += clusterStatistic.NodeAbnormal
			result.WorkloadCount += clusterStatistic.WorkloadCount
			result.WorkloadAbnormal += clusterStatistic.WorkloadAbnormal
			result.CPUCapacity += clusterStatistic.CPUCapacity
			result.CPUAllocatable += clusterStatistic.CPUAllocatable
			result.MemCapacity += clusterStatistic.MemCapacity
			result.MemAllocatable += clusterStatistic.MemAllocatable
			result.PodCount += clusterStatistic.PodCount
			clusterStatistics = append(clusterStatistics, clusterStatistic)
		}
	}
	result.Clusters = clusterStatistics
	return result
}

func NewCacher(platformClient platformversionedclient.PlatformV1Interface,
	businessClient businessversionedclient.BusinessV1Interface) Cacher {
	return &cacher{
		platformClient:      platformClient,
		businessClient:      businessClient,
		clusterStatisticSet: make(util.ClusterStatisticSet),
		clusterClientSets:   make(util.ClusterClientSets),
		dynamicClients:      make(util.DynamicClientSet),
		clusters:            make(util.ClusterSet),
		credentials:         make(util.ClusterCredentialSet),
		firstLoad:           FirstLoad,
	}
}

func (c *cacher) getTApps(cluster string) int {
	count := 0
	content, err := c.dynamicClients[cluster].Resource(TAppResource).Namespace(AllNamespaces).List(context.Background(), metav1.ListOptions{})
	if content == nil || (err != nil && !errors.IsNotFound(err)) {
		log.Error("Query TApps failed", log.Any("cluster", cluster), log.Err(err))
		return 0
	}
	count += len(content.Items)
	return count
}

func (c *cacher) getDeployments(clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if deployments, err := clientSet.AppsV1().Deployments(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(deployments.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query deployments of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func (c *cacher) getStatefulSets(clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if statefulSets, err := clientSet.AppsV1().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulSets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func (c *cacher) getDaemonSets(clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if daemonSets, err := clientSet.AppsV1().DaemonSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(daemonSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query daemonSets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func isReady(node *corev1.Node) bool {
	for _, one := range node.Status.Conditions {
		if one.Type == corev1.NodeReady && one.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isHealthy(component *corev1.ComponentStatus) bool {
	for _, one := range component.Conditions {
		if one.Type == corev1.ComponentHealthy && one.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (c *cacher) getComponentStatuses(clusterID string, clientSet *kubernetes.Clientset, health *util.ComponentHealth) {
	if componentStatuses, err := clientSet.CoreV1().ComponentStatuses().List(context.Background(), metav1.ListOptions{}); err == nil {
		for _, cs := range componentStatuses.Items {
			csName := cs.GetName()
			if _, ok := UpdateComponentStatusFunc[Component(csName)]; ok {
				UpdateComponentStatusFunc[Component(csName)](&cs, health)
			} else if strings.HasPrefix(csName, EtcdPrefix) {
				health.Etcd = health.Etcd && isHealthy(&cs)
			}
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query componentStatuses failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getNodeMetrics(clusterID string, clientSet *metricsv.Clientset, counter *util.ResourceCounter) {
	if nodeMetrics, err := clientSet.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{}); err == nil {
		for _, nm := range nodeMetrics.Items {
			if resourceCPU, ok := nm.Usage[corev1.ResourceCPU]; ok {
				counter.CPUUsed += float64(resourceCPU.MilliValue()) / float64(1000)
			}
			if resourceMem, ok := nm.Usage[corev1.ResourceMemory]; ok {
				counter.MemUsed += resourceMem.Value()
			}
		}
		if len(nodeMetrics.Items) > 0 {
			counter.HasMetricServer = true
		} else {
			counter.HasMetricServer = false
		}
	} else {
		counter.HasMetricServer = false
		log.Error("query node metrics from metric server failed", log.Any("cluster", clusterID), log.Err(err))
	}
}

func (c *cacher) getNodes(clusterID string, clientSet *kubernetes.Clientset, counter *util.ResourceCounter) {
	if nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{}); err == nil {
		counter.NodeTotal = len(nodes.Items)
		for i, node := range nodes.Items {
			if !isReady(&nodes.Items[i]) {
				counter.NodeAbnormal++
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Cpu() != nil {
				counter.CPUAllocatable += float64(node.Status.Allocatable.Cpu().MilliValue()) / float64(1000)
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Cpu() != nil {
				counter.CPUCapacity += float64(node.Status.Capacity.Cpu().MilliValue()) / float64(1000)
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Memory() != nil {
				counter.MemAllocatable += node.Status.Allocatable.Memory().Value()
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Memory() != nil {
				counter.MemCapacity += node.Status.Allocatable.Memory().Value()
			}
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query nodes  failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getPods(clusterID string, clientSet *kubernetes.Clientset, counter *util.ResourceCounter) {
	if pods, err := clientSet.CoreV1().Pods(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		counter.PodCount = len(pods.Items)
		for _, pod := range pods.Items {
			for _, ctn := range pod.Spec.Containers {
				counter.CPURequest += float64(ctn.Resources.Requests.Cpu().MilliValue()) / float64(1000)
				counter.CPULimit += float64(ctn.Resources.Limits.Cpu().Value()) / float64(1000)
				counter.MemRequest += ctn.Resources.Requests.Memory().Value()
				counter.MemLimit += ctn.Resources.Limits.Memory().Value()
			}
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query nodes  failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getWorkloadCounter(clusterID string, clientSet *kubernetes.Clientset) *util.WorkloadCounter {
	counter := &util.WorkloadCounter{}
	counter.Deployment = c.getDeployments(clusterID, clientSet)
	counter.DaemonSet = c.getDaemonSets(clusterID, clientSet)
	counter.StatefulSet = c.getStatefulSets(clusterID, clientSet)
	counter.TApp = c.getTApps(clusterID)
	log.Debugf("finish reloading cluster: %s's workload: %+v", clusterID, counter)
	return counter
}

func (c *cacher) getDynamicClients() {
	var err error
	var clusters *platformv1.ClusterList
	var clusterCredentials *platformv1.ClusterCredentialList
	clusters, err = c.platformClient.Clusters().List(context.Background(), metav1.ListOptions{})
	if err != nil || clusters == nil {
		return
	}
	clusterCredentials, err = c.platformClient.ClusterCredentials().List(context.Background(), metav1.ListOptions{})
	if err != nil || clusterCredentials == nil {
		return
	}
	for i, cls := range clusters.Items {
		c.clusters[cls.GetName()] = &clusters.Items[i]
	}
	for i, cc := range clusterCredentials.Items {
		clusterID := cc.ClusterName
		c.credentials[clusterID] = &clusterCredentials.Items[i]
		if _, ok := c.clusters[clusterID]; ok {
			dynamicClient, err := platformutil.BuildExternalDynamicClientSet(c.clusters[clusterID], c.credentials[clusterID])
			if err == nil {
				c.dynamicClients[clusterID] = dynamicClient
			}
		}
	}
}

func calResourceRate(counter *util.ResourceCounter) {
	counter.CPURequestRate = float64(0)
	counter.CPUAllocatableRate = float64(0)
	counter.CPUUsage = float64(0)
	if counter.CPUCapacity > float64(0) {
		counter.CPURequestRate = counter.CPURequest / counter.CPUCapacity
		counter.CPUAllocatableRate = counter.CPUAllocatable / counter.CPUCapacity
		counter.CPUUsage = counter.CPUUsed / counter.CPUCapacity
	}
	counter.MemRequestRate = float64(0)
	counter.MemAllocatableRate = float64(0)
	counter.MemUsage = float64(0)
	if counter.MemCapacity > int64(0) {
		counter.MemRequestRate = float64(counter.MemRequest) / float64(counter.MemCapacity)
		counter.MemAllocatableRate = float64(counter.MemAllocatable) / float64(counter.MemCapacity)
		counter.MemUsage = float64(counter.MemUsed) / float64(counter.MemCapacity)
	}
}

func transPercent(value float64) string {
	if value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", value*float64(100)), 64); err == nil {
		return fmt.Sprintf("%v%%", value)
	}
	return "0%"
}
