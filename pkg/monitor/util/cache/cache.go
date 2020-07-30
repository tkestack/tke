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
	"sync"

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

	ClusterClientSet = "ClusterClientSet"
	WorkloadCounter  = "WorkloadCounter"
	ResourceCounter  = "ResourceCounter"
	ClusterPhase     = "ClusterPhase"
	ComponentHealth  = "ComponentHealth"

	TAppResourceName = "tapps"
	TAppGroupName    = "apps.tkestack.io"

	Scheduler         Component = "scheduler"
	ControllerManager Component = "controller-manager"
	Etcd              Component = "etcd-0"
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
	GetClusterOverviewResult(clusterIDs []string) *monitor.ClusterOverviewResult
}

type cacher struct {
	sync.RWMutex
	platformClient      platformversionedclient.PlatformV1Interface
	businessClient      businessversionedclient.BusinessV1Interface
	clusterStatisticSet util.ClusterStatisticSet
	clusterClientSets   util.ClusterClientSets
	dynamicClients      util.DynamicClientSet
	clusters            util.ClusterSet
	credentials         util.ClusterCredentialSet
	clusterAbnormal     int
}

func (c *cacher) Reload() {
	c.Lock()
	defer c.Unlock()

	c.getClusters()
	c.getProjects()
}

func (c *cacher) getClusters() {
	c.getDynamicClients()

	c.clusterStatisticSet = make(util.ClusterStatisticSet)
	c.clusterAbnormal = 0
	if clusters, err := c.platformClient.Clusters().List(context.Background(), metav1.ListOptions{}); err == nil {
		wg := sync.WaitGroup{}
		wg.Add(len(clusters.Items))
		syncMap := sync.Map{}
		for i := range clusters.Items {
			if clusters.Items[i].Status.Phase == platformv1.ClusterFailed {
				c.clusterAbnormal++
			}
			go func(cls platformv1.Cluster) {
				defer wg.Done()
				clusterID := cls.GetName()
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
					log.Infof("cls: %+v, counter: %+v", clusterID, resourceCounter)
				}
				calResourceRate(resourceCounter)
				health := &util.ComponentHealth{}
				c.getComponentStatuses(clusterID, clientSet, health)
				syncMap.Store(clusterID, map[string]interface{}{
					ClusterClientSet: clientSet,
					WorkloadCounter:  workloadCounter,
					ResourceCounter:  resourceCounter,
					ClusterPhase:     string(cls.Status.Phase),
					ComponentHealth:  health,
				})
			}(clusters.Items[i])
		}

		wg.Wait()

		syncMap.Range(func(key, value interface{}) bool {
			clusterID := key.(string)
			val := value.(map[string]interface{})
			clusterClientSet := val[ClusterClientSet].(*kubernetes.Clientset)
			workloadCounter := val[WorkloadCounter].(*util.WorkloadCounter)
			resourceCounter := val[ResourceCounter].(*util.ResourceCounter)
			clusterPhase := val[ClusterPhase].(string)
			health := val[ComponentHealth].(*util.ComponentHealth)
			c.clusterClientSets[clusterID] = clusterClientSet
			c.clusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
				ClusterID:                clusterID,
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
				SchedulerHealthy:         health.Scheduler,
				ControllerManagerHealthy: health.ControllerManager,
				EtcdHealthy:              health.Etcd,
			}
			return true
		})
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

func (c *cacher) GetClusterOverviewResult(clusterIDs []string) *monitor.ClusterOverviewResult {
	c.RLock()
	defer c.RUnlock()

	clusterStatistics := make([]*monitor.ClusterStatistic, 0)
	result := &monitor.ClusterOverviewResult{}
	result.ClusterCount = int32(len(clusterIDs))
	result.ClusterAbnormal = int32(c.clusterAbnormal)
	result.NodeAbnormal = 0
	result.WorkloadAbnormal = 0
	for _, clusterID := range clusterIDs {
		if clusterStatistic, ok := c.clusterStatisticSet[clusterID]; ok {
			result.NodeCount += clusterStatistic.NodeCount
			result.NodeAbnormal += clusterStatistic.NodeAbnormal
			result.WorkloadCount += clusterStatistic.WorkloadCount
			result.WorkloadAbnormal += clusterStatistic.WorkloadAbnormal
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
	if deployments, err := clientSet.AppsV1().Deployments(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil && !errors.IsNotFound(err) {
		count += len(deployments.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query deployments of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if deployments, err := clientSet.AppsV1beta1().Deployments(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil && !errors.IsNotFound(err) {
		count += len(deployments.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query deployments of v1beta11 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if deployments, err := clientSet.AppsV1beta2().Deployments(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil && !errors.IsNotFound(err) {
		count += len(deployments.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query deployments of v1beta2 failed", log.Any("clusterID", clusterID), log.Err(err))
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

	if statefulSets, err := clientSet.AppsV1beta1().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulSets of v1beta1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if statefulSets, err := clientSet.AppsV1beta2().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulSets of v1beta2 failed", log.Any("clusterID", clusterID), log.Err(err))
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

	if daemonSets, err := clientSet.AppsV1beta2().DaemonSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(daemonSets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query daemonSets of v1beta2 failed", log.Any("clusterID", clusterID), log.Err(err))
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
			UpdateComponentStatusFunc[Component(cs.GetName())](&cs, health)
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query componentStatuses failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getNodeMetrics(clusterID string, clientSet *metricsv.Clientset, counter *util.ResourceCounter) {
	if nodeMetrics, err := clientSet.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{}); err == nil {
		for _, nm := range nodeMetrics.Items {
			if resourceCPU, ok := nm.Usage[corev1.ResourceCPU]; ok {
				log.Infof("metric server cpu: %+v", float64(resourceCPU.MilliValue())/float64(1000))
				counter.CPUUsed += float64(resourceCPU.MilliValue()) / float64(1000)
			}
			if resourceMem, ok := nm.Usage[corev1.ResourceMemory]; ok {
				log.Infof("metric server node used: %+v", resourceMem.Value())
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
	counter.DaemonSet = c.getDeployments(clusterID, clientSet)
	counter.DaemonSet = c.getDaemonSets(clusterID, clientSet)
	counter.StatefulSet = c.getStatefulSets(clusterID, clientSet)
	counter.TApp = c.getTApps(clusterID)
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
