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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sync"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/monitor"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/monitor/util"
	platformutil "tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

const (
	AllNamespaces = ""

	ClusterSet      = "ClusterSet"
	WorkloadCounter = "WorkloadCounter"
	ResourceCounter = "ResourceCounter"

	TAppResourceName = "tapps"
	TAppGroupName    = "apps.tkestack.io"
)

var (
	TAppResource = schema.GroupVersionResource{Group: TAppGroupName, Version: "v1", Resource: TAppResourceName}
)

type Cacher interface {
	Reload()
	GetClusterOverviewResult(clusterIDs []string) *monitor.ClusterOverviewResult
}

type cacher struct {
	sync.RWMutex
	platformClient      platformversionedclient.PlatformV1Interface
	clusterStatisticSet util.ClusterStatisticSet
	clusterClientSets   util.ClusterClientSets
	dynamicClients      util.DynamicClientSet
	clusters            util.ClusterSet
	credentials         util.ClusterCredentialSet
}

func (c *cacher) Reload() {
	c.Lock()
	defer c.Unlock()

	c.getDynamicClients()

	c.clusterStatisticSet = make(util.ClusterStatisticSet)
	if clusters, err := c.platformClient.Clusters().List(context.Background(), metav1.ListOptions{}); err == nil {
		wg := sync.WaitGroup{}
		wg.Add(len(clusters.Items))
		syncMap := sync.Map{}
		for i := range clusters.Items {
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
				syncMap.Store(clusterID, map[string]interface{}{
					ClusterSet:      clientSet,
					WorkloadCounter: workloadCounter,
					ResourceCounter: resourceCounter,
				})
			}(clusters.Items[i])
		}

		wg.Wait()

		syncMap.Range(func(key, value interface{}) bool {
			clusterID := key.(string)
			val := value.(map[string]interface{})
			clientSet := val[ClusterSet].(*kubernetes.Clientset)
			workloadCounter := val[WorkloadCounter].(*util.WorkloadCounter)
			resourceCounter := val[ResourceCounter].(*util.ResourceCounter)
			c.clusterClientSets[clusterID] = clientSet
			c.clusterStatisticSet[clusterID] = &monitor.ClusterStatistic{
				ClusterID:        clusterID,
				NodeCount:        resourceCounter.NodeTotal,
				NodeAbnormal:     resourceCounter.NodeAbnormal,
				WorkloadCount:    workloadCounter.Total(),
				WorkloadAbnormal: 0,
				CPURequest:       resourceCounter.CPURequest,
				CPULimit:         resourceCounter.CPULimit,
				CPUAllocatable:   resourceCounter.CPUAllocatable,
				CPUCapacity:      resourceCounter.CPUCapacity,
				MemRequest:       resourceCounter.MemRequest,
				MemLimit:         resourceCounter.MemLimit,
				MemAllocatable:   resourceCounter.MemAllocatable,
				MemCapacity:      resourceCounter.MemCapacity,
			}
			return true
		})
	}
}

func (c *cacher) GetClusterOverviewResult(clusterIDs []string) *monitor.ClusterOverviewResult {
	c.RLock()
	defer c.RUnlock()

	clusterStatistics := make([]*monitor.ClusterStatistic, 0)
	result := &monitor.ClusterOverviewResult{}
	result.ClusterCount = len(clusterIDs)
	for _, clusterID := range clusterIDs {
		if clusterStatistic, ok := c.clusterStatisticSet[clusterID]; ok {
			result.NodeCount += clusterStatistic.NodeCount
			result.WorkloadCount += clusterStatistic.WorkloadCount
			clusterStatistics = append(clusterStatistics, clusterStatistic)
		}
	}
	result.Clusters = clusterStatistics
	return result
}

func NewCacher(platformClient platformversionedclient.PlatformV1Interface) Cacher {
	return &cacher{
		platformClient:      platformClient,
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
	if statefulsets, err := clientSet.AppsV1().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulsets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulsets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if statefulsets, err := clientSet.AppsV1beta1().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulsets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulsets of v1beta1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if statefulsets, err := clientSet.AppsV1beta2().StatefulSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(statefulsets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query statefulsets of v1beta2 failed", log.Any("clusterID", clusterID), log.Err(err))
	}
	return count
}

func (c *cacher) getDaemonSets(clusterID string, clientSet *kubernetes.Clientset) int {
	count := 0
	if daemonsets, err := clientSet.AppsV1().DaemonSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(daemonsets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query daemonsets of v1 failed", log.Any("clusterID", clusterID), log.Err(err))
	}

	if daemonsets, err := clientSet.AppsV1beta2().DaemonSets(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		count += len(daemonsets.Items)
	} else if !errors.IsNotFound(err) {
		log.Error("Query daemonsets of v1beta2 failed", log.Any("clusterID", clusterID), log.Err(err))
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

func (c *cacher) getNodes(clusterID string, clientSet *kubernetes.Clientset, resourceCounter *util.ResourceCounter) {
	if nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{}); err == nil {
		resourceCounter.NodeTotal = len(nodes.Items)
		for i, node := range nodes.Items {
			if !isReady(&nodes.Items[i]) {
				resourceCounter.NodeAbnormal++
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Cpu() != nil {
				resourceCounter.CPUAllocatable += float64(node.Status.Allocatable.Cpu().MilliValue()) / float64(1000)
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Cpu() != nil {
				resourceCounter.CPUCapacity += float64(node.Status.Capacity.Cpu().MilliValue()) / float64(1000)
			}
			if node.Status.Allocatable != nil && node.Status.Allocatable.Memory() != nil {
				resourceCounter.MemAllocatable += node.Status.Allocatable.Memory().Value()
			}
			if node.Status.Capacity != nil && node.Status.Capacity.Memory() != nil {
				resourceCounter.MemCapacity += node.Status.Allocatable.Memory().Value()
			}
		}
	} else if !errors.IsNotFound(err) {
		log.Error("Query nodes  failed", log.Any("clusterID", clusterID), log.Err(err))
	}
}

func (c *cacher) getPods(clusterID string, clientSet *kubernetes.Clientset, resourceCounter *util.ResourceCounter) {
	if pods, err := clientSet.CoreV1().Pods(AllNamespaces).List(context.Background(), metav1.ListOptions{}); err == nil {
		for _, pod := range pods.Items {
			for _, ctn := range pod.Spec.Containers {
				resourceCounter.CPURequest += float64(ctn.Resources.Requests.Cpu().MilliValue()) / float64(1000)
				resourceCounter.CPULimit += float64(ctn.Resources.Limits.Cpu().Value()) / float64(1000)
				resourceCounter.MemRequest += ctn.Resources.Requests.Memory().Value()
				resourceCounter.MemLimit += ctn.Resources.Limits.Memory().Value()
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
