/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package services

import (
	"context"
	"net/http"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"tkestack.io/tke/pkg/mesh/models"

	"tkestack.io/tke/pkg/mesh/services/rest"

	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/mesh/util/errors"
)

// ClusterService cluster service interface
type ClusterService interface {
	// Get cluster
	Get(ctx context.Context, clusterName string) (*platformv1.Cluster, error)
	// List cluster
	List(ctx context.Context) ([]platformv1.Cluster, error)
	// ListNamespaces list cluster's namespaces
	ListNamespaces(ctx context.Context, clusterName string) ([]corev1.Namespace, error)
	// ListServices
	ListServices(ctx context.Context, clusterName, namespace string, selector labels.Selector) ([]corev1.Service, error)
	// Proxy proxy to Kubernetes api server
	Proxy(ctx context.Context, clusterName string) (transport http.RoundTripper, host string, err error)
}

// IstioService istio resource service interface
type IstioService interface {
	ListAllResources(
		ctx context.Context, clusterName, namespace, kind string, selector labels.Selector,
	) (map[string][]unstructured.Unstructured, *errors.MultiError)

	ListGateways(ctx context.Context, clusterName, namespace string) ([]istionetworking.Gateway, error)
	ListVirtualServices(ctx context.Context, clusterName, namespace string) ([]istionetworking.VirtualService, error)
	ListDestinationRules(ctx context.Context, clusterName, namespace string) ([]istionetworking.DestinationRule, error)
	ListServiceEntries(ctx context.Context, clusterName, namespace string) ([]istionetworking.ServiceEntry, error)

	ListResources(ctx context.Context, clusterName string, obj runtime.Object, opt ...ctrlclient.ListOption) error
	GetResource(ctx context.Context, clusterName string, obj runtime.Object) error
	CreateResource(ctx context.Context, clusterName string, obj runtime.Object) (bool, error)
	UpdateResource(ctx context.Context, clusterName string, obj runtime.Object) (bool, error)
	DeleteResource(ctx context.Context, clusterName string, obj runtime.Object) (bool, error)

	CreateNorthTrafficGateway(ctx context.Context, clusterName string, obj *models.IstioNetworkingConfig) (bool, error)
	UpdateNorthTrafficGateway(ctx context.Context, clusterName string, obj *models.IstioNetworkingConfig) (bool, error)
	DeleteNorthTrafficGateway(ctx context.Context, clusterName string, gateway *istionetworking.Gateway) (bool, error)
	GetNorthTrafficGateway(
		ctx context.Context, clusterName string, gateway *istionetworking.Gateway,
	) (*models.IstioNetworkingConfig, error)
}

// MeshClusterService
type MeshClusterService interface {
	// ListMicroServices
	//  if namespace empty, list all namespaces services
	//  if serviceName empty, list all services
	ListMicroServices(
		ctx context.Context, meshName, namespace, serviceName string, selector labels.Selector,
	) ([]rest.MicroService, *errors.MultiError)

	// CreateMeshResource create istio resource to all mesh's main clusters
	CreateMeshResource(ctx context.Context, meshName string, obj *unstructured.Unstructured) error
}

// MonitorService
type MonitorService interface {
	GetTracingMetricsData(ctx context.Context, query rest.MetricQuery) (*rest.MetricData, error)
	GetMonitorMetricsData(ctx context.Context, query *rest.MetricQuery) (*rest.MetricData, error)
	GetMonitorTopologyData(ctx context.Context, query *rest.TopoQuery) (*rest.TopoData, error)
}
