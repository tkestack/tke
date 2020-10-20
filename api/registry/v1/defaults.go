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

func SetDefaults_ConfigMap(obj *ConfigMap) {
	if obj.Data == nil {
		obj.Data = make(map[string]string)
	}
}

func SetDefaults_NamespaceSpec(obj *NamespaceSpec) {
	if obj.Visibility == "" {
		obj.Visibility = VisibilityPrivate
	}
}

func SetDefaults_RepositorySpec(obj *RepositorySpec) {
	if obj.Visibility == "" {
		obj.Visibility = VisibilityPrivate
	}
}

func SetDefaults_RepositoryStatus(obj *RepositoryStatus) {
	if obj.Tags == nil {
		obj.Tags = make([]RepositoryTag, 0)
	}
}

func SetDefaults_ChartGroupSpec(obj *ChartGroupSpec) {
	if obj.Visibility == "" {
		obj.Visibility = VisibilityPrivate
	}
	if obj.Type == "" {
		obj.Type = RepoTypeProject
	}
	if obj.Projects == nil {
		obj.Projects = make([]string, 0)
	}
	if obj.Finalizers == nil || len(obj.Finalizers) == 0 {
		obj.Finalizers = []FinalizerName{
			ChartGroupFinalize,
		}
	}
}

func SetDefaults_ChartGroupStatus(obj *ChartGroupStatus) {
	if obj.Phase == "" {
		obj.Phase = ChartGroupPending
	}
}

func SetDefaults_ChartSpec(obj *ChartSpec) {
	if obj.Visibility == "" {
		obj.Visibility = VisibilityPrivate
	}
	if obj.Finalizers == nil || len(obj.Finalizers) == 0 {
		obj.Finalizers = []FinalizerName{
			ChartFinalize,
		}
	}
}

func SetDefaults_ChartStatus(obj *ChartStatus) {
	if obj.Versions == nil {
		obj.Versions = make([]ChartVersion, 0)
	}
	if obj.Phase == "" {
		obj.Phase = ChartPending
	}
}

func SetDefaults_ChartInfoSpec(obj *ChartInfoSpec) {
	if obj.Values == nil {
		obj.Values = make(map[string]string)
	}
	if obj.Readme == nil {
		obj.Readme = make(map[string]string)
	}
}
