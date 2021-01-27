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

package config

type ClusterConfig struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Role   string `json:"role"`

	KubeMesh                KubeMeshConfiguration         `json:"kubeMesh"`
	AutoInjectionNamespaces []string                      `json:"autoInjectionNamespaces"`
	Gateways                GatewaysConfiguration         `json:"gateways"`
	SystemGateways          GatewaysConfiguration         `json:"systemGateways"`
	SidecarInjector         SidecarInjectorConfiguration  `json:"sidecarInjector,omitempty"`
	DeployMode              ComponentDeployMode           `json:"deployMode,omitempty"`
	MeshKubeOperator        MeshKubeOperatorConfiguration `json:"meshKubeOperator,omitempty"`
	Istiod                  IstiodConfiguration           `json:"istiod,omitempty"`

	Proxy     ProxyConfiguration     `json:"proxy,omitempty"`
	ProxyInit ProxyInitConfiguration `json:"proxyInit,omitempty"`

	ImageHub string `json:"imageHub,omitempty"`
}

type KubeMeshConfiguration struct {
	CommonServiceConfiguration
	Port int `json:"port"`
}

type IstiodConfiguration struct {
	TraceSampling float32 `json:"traceSampling"`
	CommonConfiguration
	CommonServiceConfiguration
}

func (c *ClusterConfig) IsMaster() bool {
	return c.Role == ClusterRoleMaster
}
