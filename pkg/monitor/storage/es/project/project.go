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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
	businessv1 "tkestack.io/tke/api/business/v1"
	resourceutil "tkestack.io/tke/pkg/monitor/storage/util"
	"tkestack.io/tke/pkg/util/log"
)

var availPromBeat = "http://prom-beat:8080/prometheus"

const (
	maxErrMsgLen = 256
)

func (s *ES) Collect() {
	// Maybe should use lister, get projects from cache
	projectList, err := s.businessClient.Projects().List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List project failed: %v", err)
		return
	}

	now := time.Now().Unix()
	var samples []TimeSeries

	for _, pro := range projectList.Items {
		tags := map[string]string{
			"project_name": pro.Name,
		}
		// get capacity resource metrics
		for _, clusterHard := range pro.Spec.Clusters {
			samples = addSamples(samples, tags, now, "project_capacity", clusterHard.Hard)
		}

		// get capacity resource metrics allocated to all clusters
		for _, clusterUsed := range pro.Status.Clusters {
			samples = addSamples(samples, tags, now, "project_capacity_cluster", clusterUsed.Used)
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

		samples = addSamples(samples, tags, now, "project_allocated", projectAllocated)

		for c, rList := range projectClusterCapacity {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			samples = addSamples(samples, clusterTags, now, "project_cluster_capacity", rList)
		}
		for c, rList := range projectClusterAllocated {
			clusterTags := tags
			clusterTags["cluster_name"] = c
			samples = addSamples(samples, clusterTags, now, "project_cluster_allocated", rList)
		}
	}
	log.Debugf("Writing project metrics")

	if len(samples) > 0 {
		err = sendSamples(samples)
		if err != nil {
			log.Errorf("Send samples failed: %v", err)
		}
	}
}

func sendSamples(samples []TimeSeries) error {
	req, err := buildWriteRequest(samples)
	if err != nil {
		return err
	}
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", availPromBeat, bytes.NewReader(req))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", "TKE")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	to, _ := time.ParseDuration("30s")
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	httpResp, err := client.Do(httpReq.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		return err
	}
	return err
}

func buildWriteRequest(samples []TimeSeries) ([]byte, error) {
	req := &WriteRequest{
		TimeSeries: samples,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	compressed := snappy.Encode(nil, data)
	return compressed, nil
}

func addSamples(samples []TimeSeries, tags map[string]string, now int64, resourcePrefix string, resources businessv1.ResourceList) []TimeSeries {
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
		var labels []Label
		for k, t := range tags {
			labels = append(labels, Label{Name: k, Value: t})
		}
		labels = append(labels, Label{Name: "__name__", Value: r})
		values := []Sample{
			{
				Value:     float64(v.MilliValue()) / 1000,
				Timestamp: now * 1000,
			},
		}
		s := TimeSeries{
			Labels:  labels,
			Samples: values,
		}
		samples = append(samples, s)
	}
	return samples
}
