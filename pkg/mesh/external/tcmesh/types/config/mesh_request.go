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

type MeshRequest struct {
	MeshMetadata `json:",inline"`

	Components ComponentsConfig       `json:"components"`
	DeployMode ComponentDeployMode    `json:"deployMode,omitempty"`
	Proxy      ProxyConfiguration     `json:"proxy,omitempty"`
	ProxyInit  ProxyInitConfiguration `json:"proxyInit,omitempty"`
}

type MeshMetadata struct {
	MeshTitle string           `json:"meshTitle"`
	Region    string           `json:"region"`
	Version   string           `json:"version"`
	TenantID  string           `json:"tenantID"`
	Topology  string           `json:"topology"`
	Mode      string           `json:"mode"`
	Clusters  []*ClusterConfig `json:"clusters,omitempty"`

	OutboundTrafficPolicy string   `json:"outboundTrafficPolicy,omitempty"`
	TraceSampling         *float32 `json:"traceSampling,omitempty"`
	DisablePolicyChecks   *bool    `json:"disablePolicyChecks,omitempty"`
	MeshMonitorEnable     *bool    `json:"meshMonitorEnable,omitempty"`
}

type ComponentsConfig struct {
	Pilot       CommonConfiguration `json:"pilot,omitempty"`
	Telemetry   CommonConfiguration `json:"telemetry,omitempty"`
	Policy      CommonConfiguration `json:"policy,omitempty"`
	Galley      CommonConfiguration `json:"galley,omitempty"`
	MeshMonitor MeshMonitorConfig   `json:"meshMonitor,omitempty"`
}

type MeshMonitorConfig struct {
	CommonConfiguration `json:",inline"`

	Debug       bool   `json:"debug,omitempty"`
	LogToStderr bool   `json:"logToStderr,omitempty"`
	LoggerDir   string `json:"loggerDir,omitempty"`
	LoggerFile  string `json:"loggerFile,omitempty"`
}

type SidecarInjectorConfig struct {
	CommonConfiguration `json:",inline"`

	// If true, sidecar injector will rewrite PodSpec for liveness
	// health check to redirect request to sidecar. This makes liveness check work
	// even when mTLS is enabled.
	RewriteAppHTTPProbe bool `json:"rewriteAppHTTPProbe,omitempty"`
}
