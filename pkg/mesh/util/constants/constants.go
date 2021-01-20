/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
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
 *
 */

package constants

import (
	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	// local test env kubeconfig files
	//   example:
	//     export LOCAL_KUBECONFIG=~/.kube/config1.yaml:~/.kube/config2.yaml
	EnvLocalKubeConfig = "LOCAL_KUBECONFIG"
)

// istio kinds name
const (
	KindDestinationRule     = "DestinationRule"
	KindDestinationRuleList = "DestinationRuleList"
	KindEnvoyFilter         = "EnvoyFilter"
	KindEnvoyFilterList     = "EnvoyFilterList"
	KindGateway             = "Gateway"
	KindGatewayList         = "GatewayList"
	KindServiceEntry        = "ServiceEntry"
	KindServiceEntryList    = "ServiceEntryList"
	KindSidecar             = "Sidecar"
	KindSidecarList         = "SidecarList"
	KindVirtualService      = "VirtualService"
	KindVirtualServiceList  = "VirtualServiceList"
	KindWorkloadEntry       = "WorkloadEntry"
	KindWorkloadEntryList   = "WorkloadEntryList"

	KindPodList         = "PodList"
	KindDeploymentList  = "DeploymentList"
	KindStatefulSetList = "StatefulSetList"
	KindDaemonSetList   = "DaemonSetList"
)

// scheme GroupVersionKind
var (
	istioNetworkingGV   = istionetworking.SchemeGroupVersion
	DestinationRule     = istioNetworkingGV.WithKind(KindDestinationRule)
	DestinationRuleList = istioNetworkingGV.WithKind(KindDestinationRuleList)
	EnvoyFilter         = istioNetworkingGV.WithKind(KindEnvoyFilter)
	EnvoyFilterList     = istioNetworkingGV.WithKind(KindEnvoyFilterList)
	Gateway             = istioNetworkingGV.WithKind(KindGateway)
	GatewayList         = istioNetworkingGV.WithKind(KindGatewayList)
	ServiceEntry        = istioNetworkingGV.WithKind(KindServiceEntry)
	ServiceEntryList    = istioNetworkingGV.WithKind(KindServiceEntryList)
	Sidecar             = istioNetworkingGV.WithKind(KindSidecar)
	SidecarList         = istioNetworkingGV.WithKind(KindSidecarList)
	VirtualService      = istioNetworkingGV.WithKind(KindVirtualService)
	VirtualServiceList  = istioNetworkingGV.WithKind(KindVirtualServiceList)
	WorkloadEntry       = istioNetworkingGV.WithKind(KindWorkloadEntry)
	WorkloadEntryList   = istioNetworkingGV.WithKind(KindWorkloadEntryList)

	coreGV          = corev1.SchemeGroupVersion
	PodList         = coreGV.WithKind(KindPodList)
	appsGV          = appsv1.SchemeGroupVersion
	DeploymentList  = appsGV.WithKind(KindDeploymentList)
	StatefulSetList = appsGV.WithKind(KindStatefulSetList)

	IstioNetworkingListGVK = []schema.GroupVersionKind{
		DestinationRuleList,
		EnvoyFilterList,
		GatewayList,
		ServiceEntryList,
		SidecarList,
		VirtualServiceList,
		WorkloadEntryList,
	}

	IstioNetworkingGVK = []schema.GroupVersionKind{
		DestinationRule,
		EnvoyFilter,
		Gateway,
		ServiceEntry,
		Sidecar,
		VirtualService,
		WorkloadEntry,
	}

	WorkloadListGVK = []schema.GroupVersionKind{
		PodList,
		DeploymentList,
		StatefulSetList,
	}
)

// istio labels, annotations, selector
var (
	IstioInjectionAnnotation = "sidecar.istio.io/inject"
	IstioAppLabelKey         = "app"
	IstioVersionLabelKey     = "version"

	existsApp, _            = labels.NewRequirement(IstioAppLabelKey, selection.Exists, nil)
	existsVersion, _        = labels.NewRequirement(IstioVersionLabelKey, selection.Exists, nil)
	IstioAppSelector        = labels.NewSelector().Add(*existsApp)
	IstioAppVersionSelector = IstioAppSelector.Add(*existsVersion)

	// TODO configure the following exclude namespaces
	/* Exclude namespaces:
	   default
	   kube-system
	   istio-system
	   kube-public
	*/
	ExcludeNamespacesSelector, _ = fields.ParseSelector(
			"metadata.namespace!=default," +
			"metadata.namespace!=kube-system," +
			"metadata.namespace!=istio-system," +
			"metadata.namespace!=kube-public",
	)
)
