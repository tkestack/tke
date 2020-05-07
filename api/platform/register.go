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

package platform

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// Scheme is the default instance of runtime.Scheme to which types in the TKE API are already registered.
	Scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

// GroupName is group name used to register these schema
const GroupName = "platform.tkestack.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Kind takes an unqualified kind and returns back a IdentityProvider qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns back a IdentityProvider qualified
// GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder collects functions that add things to a scheme.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme applies all the stored functions to the scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Cluster{},
		&ClusterList{},
		&ClusterApplyOptions{},

		&ClusterCredential{},
		&ClusterCredentialList{},

		&ClusterAddon{},
		&ClusterAddonList{},
		&ClusterAddonType{},
		&ClusterAddonTypeList{},

		&Machine{},
		&MachineList{},

		&PersistentEvent{},
		&PersistentEventList{},

		&Helm{},
		&HelmList{},
		&HelmProxyOptions{},

		&IPAM{},
		&IPAMList{},
		&IPAMProxyOptions{},

		&ConfigMap{},
		&ConfigMapList{},

		&Registry{},
		&RegistryList{},

		&TappController{},
		&TappControllerList{},
		&TappControllerProxyOptions{},

		&CronHPA{},
		&CronHPAList{},
		&CronHPAProxyOptions{},

		&Prometheus{},
		&PrometheusList{},

		&CSIOperator{},
		&CSIOperatorList{},
		&CSIProxyOptions{},

		&VolumeDecorator{},
		&VolumeDecoratorList{},
		&PVCRProxyOptions{},

		&LogCollector{},
		&LogCollectorList{},
		&LogCollectorProxyOptions{},

		&LBCF{},
		&LBCFList{},
		&LBCFProxyOptions{},
	)
	return nil
}
