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

package project

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	businessv1 "tkestack.io/tke/api/business/v1"
	resourceutil "tkestack.io/tke/pkg/monitor/storage/util"
	"tkestack.io/tke/pkg/util/log"
)

func (s *Storage) Collect() {
	// Maybe should use lister, get projects from cache
	projectList, err := s.businessClient.Projects().List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List project failed: %v", err)
		return
	}

	for _, pro := range projectList.Items {
		tags := map[string]string{
			"project_name": pro.Name,
		}
		// get capacity resource metrics
		projectCapacity := businessv1.ResourceList{}
		for _, clusterHard := range pro.Spec.Clusters {
			projectCapacity = resourceutil.ResourceAdd(projectCapacity, clusterHard.Hard)
		}
		//TODO complete projectCapacity with cluster allocatable resource
		updateMetrics(tags, "project_capacity", projectCapacity)

		// get capacity resource metrics allocated to all clusters
		projectCapacityCluster := businessv1.ResourceList{}
		for _, clusterUsed := range pro.Status.Clusters {
			projectCapacityCluster = resourceutil.ResourceAdd(projectCapacityCluster, clusterUsed.Used)
		}
		updateMetrics(tags, "project_capacity_cluster", projectCapacityCluster)

		// get allocated resource metrics both for project and each cluster
		projectAllocated := businessv1.ResourceList{}
		projectClusterCapacity := map[string]businessv1.ResourceList{}
		projectClusterAllocated := map[string]businessv1.ResourceList{}
		projectNamespaceCapacity := map[string]map[string]businessv1.ResourceList{}
		projectNamespaceAllocated := map[string]map[string]businessv1.ResourceList{}

		namespacesLists, err := s.businessClient.Namespaces(pro.Name).List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("Get(%s) namespace list failed: %v", pro.Name, err)
			continue
		}
		for _, nmSpace := range namespacesLists.Items {
			projectAllocated = resourceutil.ResourceAdd(projectAllocated, nmSpace.Status.Used)
			clusterName := nmSpace.Spec.ClusterName

			if _, ok := projectClusterCapacity[clusterName]; ok {
				projectClusterCapacity[clusterName] = resourceutil.ResourceAdd(projectClusterCapacity[clusterName], nmSpace.Spec.Hard)
			} else {
				projectClusterCapacity[clusterName] = nmSpace.Spec.Hard
			}

			if _, ok := projectClusterAllocated[clusterName]; ok {
				projectClusterAllocated[clusterName] = resourceutil.ResourceAdd(projectClusterAllocated[clusterName], nmSpace.Status.Used)
			} else {
				projectClusterAllocated[clusterName] = nmSpace.Status.Used
			}

			if _, ok := projectNamespaceCapacity[clusterName]; !ok {
				projectNamespaceCapacity[clusterName] = map[string]businessv1.ResourceList{}
			}
			projectNamespaceCapacity[clusterName][nmSpace.Spec.Namespace] = nmSpace.Spec.Hard

			if _, ok := projectNamespaceAllocated[clusterName]; !ok {
				projectNamespaceAllocated[clusterName] = map[string]businessv1.ResourceList{}
			}
			projectNamespaceAllocated[clusterName][nmSpace.Spec.Namespace] = nmSpace.Status.Used
		}

		updateMetrics(tags, "project_allocated", projectAllocated)

		for c, rList := range projectClusterCapacity {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			updateMetrics(clusterTags, "project_cluster_capacity", rList)
		}
		for c, rList := range projectClusterAllocated {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			updateMetrics(clusterTags, "project_cluster_allocated", rList)
		}
		for c, nms := range projectNamespaceCapacity {
			for nm, rList := range nms {
				nsTags := tags
				nsTags["cluster_name"] = c
				nsTags["namespace"] = nm
				nsTags["namespace_name"] = c + "-" + nm
				updateMetrics(nsTags, "project_namespace_capacity", rList)
			}
		}
		for c, nms := range projectNamespaceAllocated {
			for nm, rList := range nms {
				nsTags := tags
				nsTags["cluster_name"] = c
				nsTags["namespace"] = nm
				nsTags["namespace_name"] = c + "-" + nm
				updateMetrics(nsTags, "project_namespace_allocated", rList)
			}
		}
	}
}

func updateMetrics(tags map[string]string, resourcePrefix string, resources businessv1.ResourceList) {
	fullNameResources := businessv1.ResourceList{}
	for r, v := range resources {
		metricName := fmt.Sprintf("%s_%s", resourcePrefix, resourceutil.ResourceNameTranslate(r))
		if old, ok := fullNameResources[metricName]; ok {
			if v.Cmp(old) == 1 {
				fullNameResources[metricName] = v
			}
		} else {
			fullNameResources[metricName] = v
		}
	}

	for metricName, v := range fullNameResources {
		log.Infof("metricName: %s, tags: %s", metricName, tags)
		if metric, ok := projectMetricsMap[metricName]; ok {
			metric.With(tags).Set(float64(v.MilliValue()) / 1000)
		}
	}
}
