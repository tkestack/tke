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
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ClusterStatus(obj *ClusterStatus) {
	if obj.Phase == "" {
		obj.Phase = ClusterInitializing
	}
	if obj.Resource.Allocatable == nil {
		obj.Resource.Allocatable = make(ResourceList)
	}
	if obj.Resource.Allocated == nil {
		obj.Resource.Allocated = make(ResourceList)
	}
	if obj.Resource.Capacity == nil {
		obj.Resource.Capacity = make(ResourceList)
	}
}

func SetDefaults_MachineStatus(obj *MachineStatus) {
	if obj.Phase == "" {
		obj.Phase = MachineInitializing
	}
}

func SetDefaults_ConfigMap(obj *ConfigMap) {
	if obj.Data == nil {
		obj.Data = make(map[string]string)
	}
}

// Addon

// SetDefaults_NamespaceSetStatus sets additional defaults namespace status.
func SetDefaults_PersistentEventStatus(obj *PersistentEventStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_HelmStatus(obj *HelmStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_TappControllerStatus(obj *TappControllerStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_CSIOperatorStatus(obj *CSIOperatorStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_VolumeDecoratorStatus(obj *VolumeDecoratorStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_LogCollectorStatus(obj *LogCollectorStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_CronHPAStatus(obj *CronHPAStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_PrometheusStatus(obj *PrometheusStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_IPAMStatus(obj *IPAMStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}

func SetDefaults_LBCFStatus(obj *LBCFStatus) {
	if obj.Phase == "" {
		obj.Phase = AddonPhaseInitializing
	}
}
