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

package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MeshConfiguration contains the configuration for the Mesh
type MeshConfiguration struct {
	metav1.TypeMeta

	Region     Region
	Istio      Istio
	Components Components
}

type Region struct {
	Name string
}

type Istio struct {
	SupportedVersion []string
	// LabelSelector    *IstioLabelSelector
	Gateway *IstioGateway
}

type IstioLabelSelector struct {
	AppRuntime string
	AppRelease string
}

type IstioGateway struct {
	DefaultHttpPort int
}

type Components struct {
	MeshManager *MeshManagerConfig
}

type MeshManagerConfig struct {
	Address               string
}
