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
	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
)

var (
	projectMetricsMap = map[string]*metrics.GaugeVec{}
	projectMetrics    = []*metrics.GaugeVec{
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cpu",
				Help:           "Project capacity of cpu.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_memory",
				Help:           "Project capacity of memory.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_configmaps",
				Help:           "Project capacity of configmaps.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_ephemeral_storage",
				Help:           "Project capacity of ephemeral storage.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_persistentvolumeclaims",
				Help:           "Project capacity of persistentvolumeclaims.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_pods",
				Help:           "Project capacity of pods.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_resourcequotas",
				Help:           "Project capacity of resourcequotas.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_secrets",
				Help:           "Project capacity of secrets.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_services",
				Help:           "Project capacity of services.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_services_loadbalancers",
				Help:           "Project capacity of services.loadbalancers.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_services_nodeports",
				Help:           "Project capacity of services.nodeports.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_cpu",
				Help:           "Sum of cluster capacity of cpu for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_memory",
				Help:           "Sum of cluster capacity of memory for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_configmaps",
				Help:           "Sum of cluster capacity of configmaps for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_ephemeral_storage",
				Help:           "Sum of cluster capacity of ephemeral storage for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_persistentvolumeclaims",
				Help:           "Sum of cluster capacity of persistentvolumeclaims for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_pods",
				Help:           "Sum of cluster capacity of pods for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_resourcequotas",
				Help:           "Sum of cluster capacity of resourcequotas for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_secrets",
				Help:           "Sum of cluster capacity of secrets for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_services",
				Help:           "Sum of cluster capacity of services for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_services_loadbalancers",
				Help:           "Sum of cluster capacity of services loadbalancers for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_capacity_cluster_services_nodeports",
				Help:           "Sum of cluster capacity of services nodeports for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_cpu",
				Help:           "Cpu allocated to pods for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_memory",
				Help:           "Memory allocated to pods for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_configmaps",
				Help:           "Configmaps allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_ephemeral_storage",
				Help:           "Ephemeral-storage allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_persistentvolumeclaims",
				Help:           "Persistentvolumeclaims allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_pods",
				Help:           "Pods allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_resourcequotas",
				Help:           "Resourcequotas allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_secrets",
				Help:           "Secrets allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_services",
				Help:           "Services allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_services_loadbalancers",
				Help:           "Services loadbalancers allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_allocated_services_nodeports",
				Help:           "Services nodeports allocated for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_cpu",
				Help:           "Cpu capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_memory",
				Help:           "Memory capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_configmaps",
				Help:           "Configmaps capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_ephemeral_storage",
				Help:           "Ephemeral-storage capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_persistentvolumeclaims",
				Help:           "Persistentvolumeclaims capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_pods",
				Help:           "Pods capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_resourcequotas",
				Help:           "Resourcequotas capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_secrets",
				Help:           "Secrets capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_services",
				Help:           "Services capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_services_loadbalancers",
				Help:           "Services loadbalancers capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_capacity_services_nodeports",
				Help:           "Services nodeports capacity of namespaces of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_cpu",
				Help:           "Cpu allocated to pods of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_memory",
				Help:           "Memory allocated to pods of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_configmaps",
				Help:           "Configmaps allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_ephemeral_storage",
				Help:           "Ephemeral-storage allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_persistentvolumeclaims",
				Help:           "Persistentvolumeclaims allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_pods",
				Help:           "Pods allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_resourcequotas",
				Help:           "Resourcequotas allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_secrets",
				Help:           "Secrets allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_services",
				Help:           "Services allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_services_loadbalancers",
				Help:           "Services loadbalancers allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_cluster_allocated_services_nodeports",
				Help:           "Services nodeports allocated of each cluster for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_cpu",
				Help:           "Cpu capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_memory",
				Help:           "Memory capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_configmaps",
				Help:           "Configmaps capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_ephemeral_storage",
				Help:           "Ephemeral-storage capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_persistentvolumeclaims",
				Help:           "Persistentvolumeclaims capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_pods",
				Help:           "Pods capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_resourcequotas",
				Help:           "Resourcequotas capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_secrets",
				Help:           "Secrets capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_services",
				Help:           "Services capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_services_loadbalancers",
				Help:           "Services loadbalancers capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_capacity_services_nodeports",
				Help:           "Services nodeports capacity of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),

		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_cpu",
				Help:           "Cpu allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_memory",
				Help:           "Memory allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_configmaps",
				Help:           "Configmaps allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_ephemeral_storage",
				Help:           "Ephemeral-storage allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_persistentvolumeclaims",
				Help:           "Persistentvolumeclaims allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_pods",
				Help:           "Pods allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_resourcequotas",
				Help:           "Resourcequotas allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_secrets",
				Help:           "Secrets allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_services",
				Help:           "Services allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_services_loadbalancers",
				Help:           "Services loadbalancers allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
		metrics.NewGaugeVec(
			&metrics.GaugeOpts{
				Name:           "project_namespace_allocated_services_nodeports",
				Help:           "Services nodeports allocated of each namespace for each project.",
				StabilityLevel: metrics.ALPHA,
			},
			[]string{"project_name", "cluster_name", "namespace", "namespaceName"},
		),
	}
)

func init() {
	for _, metric := range projectMetrics {
		legacyregistry.MustRegister(metric)
		projectMetricsMap[metric.Name] = metric
	}
}
