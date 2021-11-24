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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/monitor"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/platform/util/addon"
	"tkestack.io/tke/pkg/util/log"
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

	SchedulerPrefix         = "kube-scheduler-"
	ControllerManagerPrefix = "kube-controller-manager-"
	EtcdPrefix              = "etcd-"

	FirstLoad    = int32(1)
	NotFirstLoad = int32(0)
)

var (
	TAppResource = schema.GroupVersionResource{Group: TAppGroupName,
		Version: "v1", Resource: TAppResourceName}
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
	c.getClusters(context.Background())
	c.getProjects()

	if c.firstLoad == FirstLoad {
		atomic.StoreInt32(&c.firstLoad, NotFirstLoad)
	}
}

func (c *cacher) updateClusters(curClusterSet util.ClusterSet, curClusterCredentialSet util.ClusterCredentialSet,
	curDynamicClientSet util.DynamicClientSet, curClusterAbnormal int, curClusterStatisticSet util.ClusterStatisticSet,
	curClusterClientSets util.ClusterClientSets) {
	c.clusters = curClusterSet
	c.credentials = curClusterCredentialSet
	c.dynamicClients = curDynamicClientSet
	c.clusterAbnormal = curClusterAbnormal
	c.clusterStatisticSet = curClusterStatisticSet
	c.clusterClientSets = curClusterClientSets
}

func (c *cacher) getClusters(ctx context.Context) {
	if atomic.LoadInt32(&c.firstLoad) == FirstLoad {
		log.Info("outer lock in getCluster")
		c.Lock()
		defer c.Unlock()
	}
	curClusterSet, curClusterCredentialSet, curDynamicClientSet := c.getDynamicClients(ctx)
	curClusterAbnormal := 0
	curClusterStatisticSet := make(util.ClusterStatisticSet)
	curClusterClientSets := make(util.ClusterClientSets)

	if clusters, err := c.platformClient.Clusters().List(ctx, metav1.ListOptions{}); err == nil {
		wg := sync.WaitGroup{}
		wg.Add(len(clusters.Items))
		syncMap := sync.Map{}
		finished := int32(0)
		allTask := len(clusters.Items)
		started := time.Now()
		for i := range clusters.Items {
			if clusters.Items[i].Status.Phase == platformv1.ClusterFailed {
				curClusterAbnormal++
			}
			go func(cls platformv1.Cluster, dynamicClientSets util.DynamicClientSet) {
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
				clientSet, err := addon.BuildExternalClientSet(ctx, &cls, c.platformClient)
				if err != nil {
					log.Error("create clientSet of cluster failed",
						log.Any("cluster", clusterID), log.Err(err))
					return
				}
				workloadCounter := c.getWorkloadCounter(ctx, dynamicClientSets, clusterID, clientSet)
				resourceCounter := &util.ResourceCounter{
					CPUCapacityMap:            map[string]map[string]float64{},
					CPUAllocatableMap:         map[string]map[string]float64{},
					CPUNotReadyCapacityMap:    map[string]map[string]float64{},
					CPUNotReadyAllocatableMap: map[string]map[string]float64{},
					MemCapacityMap:            map[string]map[string]int64{},
					MemAllocatableMap:         map[string]map[string]int64{},
					MemNotReadyCapacityMap:    map[string]map[string]int64{},
					MemNotReadyAllocatableMap: map[string]map[string]int64{},
				}
				c.getNodes(ctx, clusterID, clientSet, resourceCounter)
				c.getPods(ctx, clusterID, clientSet, resourceCounter)
				if metricServerClientSet, err := c.getMetricServerClientSet(ctx, &cls); err == nil && metricServerClientSet != nil {
					c.getNodeMetrics(ctx, clusterID, metricServerClientSet, resourceCounter)
				}
				calResourceRate(resourceCounter)
				health := &util.ComponentHealth{}
				c.getComponentStatuses(ctx, clusterID, clientSet, health)
				log.Infof("cluster: %v's components' health: %+v", clusterID, health)
				syncMap.Store(clusterID, map[string]interface{}{
					ClusterClientSet:   clientSet,
					WorkloadCounter:    workloadCounter,
					ResourceCounter:    resourceCounter,
					ClusterPhase:       string(cls.Status.Phase),
					ClusterDisplayName: clusterDisplayName,
					TenantID:           tenantID,
					ComponentHealth:    health,
				})
			}(clusters.Items[i], curDynamicClientSet)
		}

		wg.Wait()

		log.Debugf("finish reloading all clusters, cost: %v seconds", time.Since(started).Seconds())
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
				curClusterClientSets[clusterID] = clusterClientSet
				curClusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
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
					CPUNotReadyAllocatable:   resourceCounter.CPUNotReadyAllocatable,
					CPUNotReadyCapacity:      resourceCounter.CPUNotReadyCapacity,
					CPURequestRate:           transPercent(resourceCounter.CPURequestRate),
					CPUAllocatableRate:       transPercent(resourceCounter.CPUAllocatableRate),
					CPUUsage:                 transPercent(resourceCounter.CPUUsage),
					MemUsed:                  resourceCounter.MemUsed,
					MemRequest:               resourceCounter.MemRequest,
					MemLimit:                 resourceCounter.MemLimit,
					MemAllocatable:           resourceCounter.MemAllocatable,
					MemCapacity:              resourceCounter.MemCapacity,
					MemNotReadyAllocatable:   resourceCounter.MemNotReadyAllocatable,
					MemNotReadyCapacity:      resourceCounter.MemNotReadyCapacity,
					MemRequestRate:           transPercent(resourceCounter.MemRequestRate),
					MemAllocatableRate:       transPercent(resourceCounter.MemAllocatableRate),
					MemUsage:                 transPercent(resourceCounter.MemUsage),
					PodCount:                 int32(resourceCounter.PodCount),
					SchedulerHealthy:         health.Scheduler,
					ControllerManagerHealthy: health.ControllerManager,
					EtcdHealthy:              health.Etcd,
				}
			} else {
				curClusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
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
	if atomic.LoadInt32(&c.firstLoad) != FirstLoad {
		log.Info("inner lock in getCluster")
		c.Lock()
		defer c.Unlock()
	}
	c.updateClusters(curClusterSet, curClusterCredentialSet, curDynamicClientSet, curClusterAbnormal,
		curClusterStatisticSet, curClusterClientSets)
}

func (c *cacher) getMetricServerClientSet(ctx context.Context, cls *platformv1.Cluster) (*metricsv.Clientset, error) {
	cc, err := addon.GetClusterCredentialV1(ctx, c.platformClient, cls)
	if err != nil {
		log.Error("query cluster credential failed", log.Any("cluster", cls.GetName()), log.Err(err))
		return nil, err
	}

	restConfig := cc.RESTConfig(cls)

	return metricsv.NewForConfig(restConfig)
}

// TODO
func (c *cacher) getProjects() {
}

func (c *cacher) GetClusterOverviewResult(clusters []*platformv1.Cluster) *monitor.ClusterOverviewResult {
	c.RLock()
	defer c.RUnlock()

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
			result.CPUNotReadyCapacity += clusterStatistic.CPUNotReadyCapacity
			result.CPUNotReadyAllocatable += clusterStatistic.CPUNotReadyAllocatable
			result.MemCapacity += clusterStatistic.MemCapacity
			result.MemAllocatable += clusterStatistic.MemAllocatable
			result.MemNotReadyCapacity += clusterStatistic.MemNotReadyCapacity
			result.MemNotReadyAllocatable += clusterStatistic.MemNotReadyAllocatable
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

func (c *cacher) getTApps(ctx context.Context, curDynamicClientSet util.DynamicClientSet, cluster string) int {
	count := 0
	content, err := curDynamicClientSet[cluster].Resource(TAppResource).
		Namespace(AllNamespaces).List(ctx, metav1.ListOptions{})
	if content == nil || (err != nil && !errors.IsNotFound(err)) {
		log.Error("Query TApps failed", log.Any("cluster", cluster), log.Err(err))
		return 0
	}
	count += len(content.Items)
	return count
}

func (c *cacher) getDeployments(ctx context.Context, clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if deployments, err := clientSet.AppsV1().Deployments(AllNamespaces).List(ctx, metav1.ListOptions{}); err == nil {
		count += len(deployments.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query deployments of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func (c *cacher) getStatefulSets(ctx context.Context, clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if statefulSets, err := clientSet.AppsV1().StatefulSets(AllNamespaces).
		List(ctx, metav1.ListOptions{}); err == nil {
		count += len(statefulSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulSets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func (c *cacher) getDaemonSets(ctx context.Context, clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if daemonSets, err := clientSet.AppsV1().DaemonSets(AllNamespaces).List(ctx, metav1.ListOptions{}); err == nil {
		count += len(daemonSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query daemonSets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func nodeIsReady(node *corev1.Node) bool {
	for _, one := range node.Status.Conditions {
		if one.Type == corev1.NodeReady && one.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func podIsReady(pod *corev1.Pod) bool {
	for _, one := range pod.Status.Conditions {
		if one.Type == corev1.PodReady && one.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsScheduled(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodScheduled && condition.Status == corev1.ConditionTrue {
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

func allPodsAreReady(pods []corev1.Pod) bool {
	if len(pods) == 0 {
		return false
	}
	for _, pod := range pods {
		if !podIsReady(&pod) {
			log.Infof("pod: %+v in %v is not ready", pod.GetName(), pod.GetNamespace())
			return false
		}
	}
	return true
}

func checkComponentsHealthFromPods(ctx context.Context, clientSet *kubernetes.Clientset) (
	*util.ComponentHealth, error) {
	podList, err := clientSet.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	schedulerPods := make([]corev1.Pod, 0)
	controllerManagerPods := make([]corev1.Pod, 0)
	etcdPods := make([]corev1.Pod, 0)
	for _, pod := range podList.Items {
		switch {
		case strings.HasPrefix(pod.GetName(), SchedulerPrefix):
			schedulerPods = append(schedulerPods, pod)
		case strings.HasPrefix(pod.GetName(), ControllerManagerPrefix):
			controllerManagerPods = append(controllerManagerPods, pod)
		case strings.HasPrefix(pod.GetName(), EtcdPrefix):
			etcdPods = append(etcdPods, pod)
		}
	}
	health := &util.ComponentHealth{
		Scheduler:         allPodsAreReady(schedulerPods),
		ControllerManager: allPodsAreReady(controllerManagerPods),
		Etcd:              allPodsAreReady(etcdPods),
	}
	return health, nil
}

func (c *cacher) getComponentStatuses(ctx context.Context, clusterID string, clientSet *kubernetes.Clientset,
	health *util.ComponentHealth) {
	if componentStatuses, err := clientSet.CoreV1().ComponentStatuses().List(ctx, metav1.ListOptions{}); err == nil {
		for _, cs := range componentStatuses.Items {
			csName := cs.GetName()
			if _, ok := UpdateComponentStatusFunc[Component(csName)]; ok {
				UpdateComponentStatusFunc[Component(csName)](&cs, health)
			} else if strings.HasPrefix(csName, EtcdPrefix) {
				health.Etcd = health.Etcd && isHealthy(&cs)
			}
		}
	}
	if !health.AllHealthy() {
		if val, err := checkComponentsHealthFromPods(ctx, clientSet); err == nil {
			*health = *val
		}
	}
}

func (c *cacher) getNodeMetrics(ctx context.Context, clusterID string,
	clientSet *metricsv.Clientset, counter *util.ResourceCounter) {
	if nodeMetrics, err := clientSet.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{}); err == nil {
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

func (c *cacher) getNodes(ctx context.Context, clusterID string,
	clientSet *kubernetes.Clientset, counter *util.ResourceCounter) {
	cpuCapacityMap := map[string]float64{}
	cpuAllocatableMap := map[string]float64{}
	cpuNotReadyCapacityMap := map[string]float64{}
	cpuNotReadyAllocatableMap := map[string]float64{}
	memCapacityMap := map[string]int64{}
	memAllocatableMap := map[string]int64{}
	memNotReadyCapacityMap := map[string]int64{}
	memNotReadyAllocatableMap := map[string]int64{}
	if val, ok := counter.CPUCapacityMap[clusterID]; ok && val != nil {
		cpuCapacityMap = val
	}
	if val, ok := counter.CPUAllocatableMap[clusterID]; ok && val != nil {
		cpuAllocatableMap = val
	}
	if val, ok := counter.CPUNotReadyCapacityMap[clusterID]; ok && val != nil {
		cpuNotReadyCapacityMap = val
	}
	if val, ok := counter.CPUNotReadyAllocatableMap[clusterID]; ok && val != nil {
		cpuNotReadyAllocatableMap = val
	}
	if val, ok := counter.MemCapacityMap[clusterID]; ok && val != nil {
		memCapacityMap = val
	}
	if val, ok := counter.MemAllocatableMap[clusterID]; ok && val != nil {
		memAllocatableMap = val
	}
	if val, ok := counter.MemNotReadyCapacityMap[clusterID]; ok && val != nil {
		memNotReadyCapacityMap = val
	}
	if val, ok := counter.MemNotReadyAllocatableMap[clusterID]; ok && val != nil {
		memNotReadyAllocatableMap = val
	}
	var cpuAllocatableInc, cpuCapacityInc, cpuNotReadyAllocatableInc, cpuNotReadyCapacityInc float64
	var memAllocatableInc, memCapacityInc, memNotReadyAllocatableInc, memNotReadyCapacityInc int64
	if nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{}); err == nil {
		counter.NodeTotal = len(nodes.Items)
		for i, node := range nodes.Items {
			if !nodeIsReady(&nodes.Items[i]) {
				counter.NodeAbnormal++
				if node.Status.Allocatable != nil && node.Status.Allocatable.Cpu() != nil {
					cpuNotReadyAllocatableInc = float64(node.Status.Allocatable.Cpu().MilliValue()) / float64(1000)
					counter.CPUNotReadyAllocatable += cpuNotReadyAllocatableInc
					cpuNotReadyAllocatableMap[node.GetName()] = cpuNotReadyAllocatableInc
				}
				if node.Status.Capacity != nil && node.Status.Capacity.Cpu() != nil {
					cpuNotReadyCapacityInc = float64(node.Status.Capacity.Cpu().MilliValue()) / float64(1000)
					counter.CPUNotReadyCapacity += cpuNotReadyCapacityInc
					cpuNotReadyCapacityMap[node.GetName()] = cpuNotReadyCapacityInc
				}
				if node.Status.Allocatable != nil && node.Status.Allocatable.Memory() != nil {
					memNotReadyAllocatableInc = node.Status.Allocatable.Memory().Value()
					counter.MemNotReadyAllocatable += memNotReadyAllocatableInc
					memNotReadyAllocatableMap[node.GetName()] = memNotReadyAllocatableInc
				}
				if node.Status.Capacity != nil && node.Status.Capacity.Memory() != nil {
					memNotReadyCapacityInc = node.Status.Allocatable.Memory().Value()
					counter.MemNotReadyCapacity += memNotReadyCapacityInc
					memNotReadyCapacityMap[node.GetName()] = memNotReadyCapacityInc
				}
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Cpu() != nil {
				cpuAllocatableInc = float64(node.Status.Allocatable.Cpu().MilliValue()) / float64(1000)
				counter.CPUAllocatable += cpuAllocatableInc
				cpuAllocatableMap[node.GetName()] = cpuAllocatableInc
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Cpu() != nil {
				cpuCapacityInc = float64(node.Status.Capacity.Cpu().MilliValue()) / float64(1000)
				counter.CPUCapacity += cpuCapacityInc
				cpuCapacityMap[node.GetName()] = cpuCapacityInc
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Memory() != nil {
				memAllocatableInc = node.Status.Allocatable.Memory().Value()
				counter.MemAllocatable += memAllocatableInc
				memAllocatableMap[node.GetName()] = memAllocatableInc
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Memory() != nil {
				memCapacityInc = node.Status.Allocatable.Memory().Value()
				counter.MemCapacity += memCapacityInc
				memCapacityMap[node.GetName()] = memCapacityInc
			}
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query nodes  failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	counter.CPUCapacityMap[clusterID] = cpuCapacityMap
	counter.CPUAllocatableMap[clusterID] = cpuAllocatableMap
	counter.CPUNotReadyCapacityMap[clusterID] = cpuNotReadyCapacityMap
	counter.CPUNotReadyAllocatableMap[clusterID] = cpuNotReadyAllocatableMap
	counter.MemCapacityMap[clusterID] = memCapacityMap
	counter.MemAllocatableMap[clusterID] = memAllocatableMap
	counter.MemNotReadyCapacityMap[clusterID] = memNotReadyCapacityMap
	counter.MemNotReadyAllocatableMap[clusterID] = memNotReadyAllocatableMap
}

func (c *cacher) getPods(ctx context.Context, clusterID string,
	clientSet *kubernetes.Clientset, counter *util.ResourceCounter) {
	if pods, err := clientSet.CoreV1().Pods(AllNamespaces).List(ctx, metav1.ListOptions{}); err == nil {
		counter.PodCount = len(pods.Items)
		nodePodMap := make(map[string][]corev1.Pod)
		for _, pod := range pods.Items {
			if !IsScheduled(&pod) {
				continue
			}
			podNode := pod.Spec.NodeName
			if _, ok := nodePodMap[podNode]; !ok {
				nodePodMap[podNode] = make([]corev1.Pod, 0)
			}
			nodePodMap[podNode] = append(nodePodMap[podNode], pod)
		}
		for nodeName := range nodePodMap {
			pods := nodePodMap[nodeName]
			nodeCPURequest := float64(0)
			nodeCPULimit := float64(0)
			nodeMemRequest := int64(0)
			nodeMemLimit := int64(0)
			for _, pod := range pods {
				for _, ctn := range pod.Spec.Containers {
					cpuRequestInc := float64(ctn.Resources.Requests.Cpu().MilliValue()) / float64(1000)
					memRequestInc := ctn.Resources.Requests.Memory().Value()
					cpuLimitInc := float64(ctn.Resources.Limits.Cpu().Value()) / float64(1000)
					memLimitInc := ctn.Resources.Limits.Memory().Value()

					if cpuRequestInc == float64(0) && cpuLimitInc > float64(0) {
						cpuRequestInc = cpuLimitInc
					} else if cpuRequestInc > float64(0) && cpuLimitInc == float64(0) {
						if outer, ok := counter.CPUAllocatableMap[clusterID]; ok && outer != nil {
							if inner, ok := outer[nodeName]; ok && inner > float64(0) {
								cpuLimitInc = inner
							}
						}
					}

					if memRequestInc == int64(0) && memLimitInc > int64(0) {
						memRequestInc = memLimitInc
					} else if memRequestInc > int64(0) && memLimitInc == int64(0) {
						if outer, ok := counter.MemAllocatableMap[clusterID]; ok && outer != nil {
							if inner, ok := outer[nodeName]; ok && inner > int64(0) {
								memLimitInc = inner
							}
						}
					}

					nodeCPURequest += cpuRequestInc
					nodeCPULimit += cpuLimitInc
					nodeMemRequest += memRequestInc
					nodeMemLimit += memLimitInc
				}
			}
			if outer, ok := counter.CPUAllocatableMap[clusterID]; ok && outer != nil {
				if inner, ok := outer[nodeName]; ok && inner > float64(0) && inner < nodeCPULimit {
					nodeCPULimit = inner
				}
			}
			if outer, ok := counter.MemAllocatableMap[clusterID]; ok && outer != nil {
				if inner, ok := outer[nodeName]; ok && inner > int64(0) && inner < nodeMemLimit {
					nodeMemLimit = inner
				}
			}
			counter.CPURequest += nodeCPURequest
			counter.CPULimit += nodeCPULimit
			counter.MemRequest += nodeMemRequest
			counter.MemLimit += nodeMemLimit
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query nodes  failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getWorkloadCounter(ctx context.Context, curDynamicClientSet util.DynamicClientSet, clusterID string,
	clientSet *kubernetes.Clientset) *util.WorkloadCounter {
	counter := &util.WorkloadCounter{}
	counter.Deployment = c.getDeployments(ctx, clusterID, clientSet)
	counter.DaemonSet = c.getDaemonSets(ctx, clusterID, clientSet)
	counter.StatefulSet = c.getStatefulSets(ctx, clusterID, clientSet)
	counter.TApp = c.getTApps(ctx, curDynamicClientSet, clusterID)
	log.Debugf("finish reloading cluster: %s's workload: %+v", clusterID, counter)
	return counter
}

func (c *cacher) getDynamicClients(ctx context.Context) (util.ClusterSet,
	util.ClusterCredentialSet, util.DynamicClientSet) {
	var err error
	var clusters *platformv1.ClusterList
	var clusterCredentials *platformv1.ClusterCredentialList
	resClusterSet := make(util.ClusterSet)
	resClusterCredentialSet := make(util.ClusterCredentialSet)
	resDynamicClientSet := make(util.DynamicClientSet)
	clusters, err = c.platformClient.Clusters().List(ctx, metav1.ListOptions{})
	if err != nil || clusters == nil {
		return resClusterSet, resClusterCredentialSet, resDynamicClientSet
	}
	clusterCredentials, err = c.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{})
	if err != nil || clusterCredentials == nil {
		return resClusterSet, resClusterCredentialSet, resDynamicClientSet
	}
	for i, cls := range clusters.Items {
		resClusterSet[cls.GetName()] = &clusters.Items[i]
	}
	for i, cc := range clusterCredentials.Items {
		clusterID := cc.ClusterName
		resClusterCredentialSet[clusterID] = &clusterCredentials.Items[i]
		if _, ok := resClusterSet[clusterID]; ok {
			if resClusterSet[clusterID] != nil && *resClusterSet[clusterID].Status.Locked {
				return resClusterSet, resClusterCredentialSet, resDynamicClientSet
			}
			restConfig := resClusterCredentialSet[clusterID].RESTConfig(resClusterSet[clusterID])
			dynamicClient, err := dynamic.NewForConfig(restConfig)
			if err == nil {
				resDynamicClientSet[clusterID] = dynamicClient
			}
		}
	}
	return resClusterSet, resClusterCredentialSet, resDynamicClientSet
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
