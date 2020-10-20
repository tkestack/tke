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

func SetDefaults_AppSpec(obj *AppSpec) {
	if obj.Finalizers == nil || len(obj.Finalizers) == 0 {
		obj.Finalizers = []FinalizerName{
			AppFinalize,
		}
	}
	if obj.Values.Values == nil || len(obj.Values.Values) == 0 {
		obj.Values.Values = make([]string, 0)
	}
}

func SetDefaults_AppStatus(obj *AppStatus) {
	if obj.Phase == "" {
		obj.Phase = AppPhaseInstalling
	}
}

func SetDefaults_AppHistorySpec(obj *AppHistorySpec) {
	if obj.Histories == nil || len(obj.Histories) == 0 {
		obj.Histories = make([]History, 0)
	}
}

func SetDefaults_ConfigMap(obj *ConfigMap) {
	if obj.Data == nil {
		obj.Data = make(map[string]string)
	}
}
