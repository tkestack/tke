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
	influxclient "github.com/influxdata/influxdb1-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
	businessv1 "tkestack.io/tke/api/business/v1"
	resourceutil "tkestack.io/tke/pkg/monitor/storage/util"
	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/util/log"
)

func (s *InfluxDB) Collect() {
	// Maybe should use lister, get projects from cache
	projectList, err := s.businessClient.Projects().List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List project failed: %v", err)
		return
	}

	bp, _ := influxclient.NewBatchPoints(influxclient.BatchPointsConfig{
		Database:  util.ProjectDatabaseName,
		Precision: "us",
	})

	now := time.Now()

	for _, pro := range projectList.Items {
		tags := map[string]string{
			"project_name": pro.Name,
		}
		// get capacity resource metrics
		for _, clusterHard := range pro.Spec.Clusters {
			addBatchPoint(bp, tags, now, "project_capacity", clusterHard.Hard)
		}

		// get capacity resource metrics allocated to all clusters
		for _, clusterUsed := range pro.Status.Clusters {
			addBatchPoint(bp, tags, now, "project_capacity_cluster", clusterUsed.Used)
		}

		// get allocated resource metrics both for project and each cluster
		projectAllocated := businessv1.ResourceList{}
		projectClusterCapacity := map[string]businessv1.ResourceList{}
		projectClusterAllocated := map[string]businessv1.ResourceList{}

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
		}

		addBatchPoint(bp, tags, now, "project_allocated", projectAllocated)

		for c, rList := range projectClusterCapacity {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			addBatchPoint(bp, clusterTags, now, "project_cluster_capacity", rList)
		}
		for c, rList := range projectClusterAllocated {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			addBatchPoint(bp, clusterTags, now, "project_cluster_allocated", rList)
		}
	}
	log.Debugf("Writing project metrics: %v", bp.Points())

	for _, db := range s.clients {
		err = db.Write(bp)
		if err != nil {
			log.Errorf("Write project resource metrics failed: %v", err)
			continue
		}
		_ = db.Close()
	}
}

func addBatchPoint(bp influxclient.BatchPoints, tags map[string]string, now time.Time, resourcePrefix string, resources businessv1.ResourceList) {
	newResources := businessv1.ResourceList{}
	for r, v := range resources {
		resMetric := fmt.Sprintf("%s_%s", resourcePrefix, resourceutil.ResourceNameTranslate(r))
		if old, ok := newResources[resMetric]; ok {
			if v.Cmp(old) == 1 {
				newResources[resMetric] = v
			}
		} else {
			newResources[resMetric] = v
		}
	}

	for r, v := range newResources {
		fields := map[string]interface{}{
			"value": float64(v.MilliValue()) / 1000,
		}
		pt, err := influxclient.NewPoint(r, tags, fields, now)
		if err != nil {
			log.Errorf("New point failed: %v", err)
			continue
		}
		bp.AddPoint(pt)
	}
}
