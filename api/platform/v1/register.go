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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package.
const GroupName = "platform.tkestack.io"

// Version is the version name use in this package.
const Version = "v1"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: Version}

var (
	// SchemeBuilder collects functions that add things to a scheme.
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	// AddToScheme applies all the stored functions to the scheme.
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(addKnownTypes, addConversionFuncs, addDefaultingFuncs)
}

// Adds the list of known types to the given scheme.
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
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// Resource takes an unqualified resource and returns a LocalGroup qualified
// GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
