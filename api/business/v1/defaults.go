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

func SetDefaults_ProjectStatus(obj *ProjectStatus) {
	if obj.Phase == "" {
		obj.Phase = ProjectPending
	}
	if obj.CalculatedChildProjects == nil {
		obj.CalculatedChildProjects = []string{}
	}
	if obj.CalculatedNamespaces == nil {
		obj.CalculatedNamespaces = []string{}
	}
	if obj.Clusters == nil {
		obj.Clusters = make(ClusterUsed)
	}
}

func SetDefaults_ProjectSpec(obj *ProjectSpec) {
	if obj.Members == nil {
		obj.Members = []string{}
	}
	if obj.Clusters == nil {
		obj.Clusters = make(ClusterHard)
	}
}

func SetDefaults_NamespaceSpec(obj *NamespaceSpec) {
	if obj.Hard == nil {
		obj.Hard = make(ResourceList)
	}
}

func SetDefaults_NamespaceStatus(obj *NamespaceStatus) {
	if obj.Phase == "" {
		obj.Phase = NamespacePending
	}
	if obj.Used == nil {
		obj.Used = make(ResourceList)
	}
}

func SetDefaults_ConfigMap(obj *ConfigMap) {
	if obj.Data == nil {
		obj.Data = make(map[string]string)
	}
}

func SetDefaults_ImageNamespaceStatus(obj *ImageNamespaceStatus) {
	if obj.Phase == "" {
		obj.Phase = ImageNamespacePending
	}
}

func SetDefaults_ChartGroupStatus(obj *ChartGroupStatus) {
	if obj.Phase == "" {
		obj.Phase = ChartGroupPending
	}
}

func SetDefaults_NsEmigrationStatus(obj *NsEmigrationStatus) {
	if obj.Phase == "" {
		obj.Phase = NsEmigrationPending
	}
}
