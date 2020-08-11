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

	"github.com/prometheus/client_golang/prometheus"
)

var (
	projectMetricsMap = map[string]*prometheus.GaugeVec{}
	resources         = []string{
		"cpu",
		"memory",
		"configmaps",
		"ephemeral_storage",
		"persistentvolumeclaims",
		"pods",
		"resourcequotas",
		"secrets",
		"services",
		"services_loadbalancers",
		"services_nodeports",
	}
)

func init() {
	for _, resource := range resources {
		prefix := "project_capacity"
		name := fmt.Sprintf("%s_%s", prefix, resource)
		help := fmt.Sprintf("Project capacity of %s", resource)
		labelNames := []string{"project_name"}
		metric := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_capacity_cluster"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("Sum of cluster capacity of %s for each project.", resource)
		labelNames = []string{"project_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_allocated"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("%s allocated for each project.", resource)
		labelNames = []string{"project_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_cluster_capacity"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("%s capacity of namespaces of each cluster for each project.", resource)
		labelNames = []string{"project_name", "cluster_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_cluster_allocated"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("%s allocated of each cluster for each project.", resource)
		labelNames = []string{"project_name", "cluster_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_namespace_capacity"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("%s capacity of each namespace for each project.", resource)
		labelNames = []string{"project_name", "cluster_name", "namespace", "namespace_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric

		prefix = "project_namespace_allocated"
		name = fmt.Sprintf("%s_%s", prefix, resource)
		help = fmt.Sprintf("%s allocated of each namespace for each project.", resource)
		labelNames = []string{"project_name", "cluster_name", "namespace", "namespace_name"}
		metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name, Help: help}, labelNames)
		prometheus.MustRegister(metric)
		projectMetricsMap[name] = metric
	}
}
