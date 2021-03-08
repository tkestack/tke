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

package state

type ClusterState struct {
	Name  string `json:"name"`
	Phase string `json:"phase"`

	ImageHub string `json:"imageHub,omitempty"`

	Registry string `json:"registry,omitempty"`
	Network  string `json:"network,omitempty"` // network id

	ControlPlaneLB  string `json:"controlPlaneLB,omitempty"`
	MultiInnerLB    string `json:"multiInnerLB,omitempty"`
	KubeLB          string `json:"kubeLB,omitempty"`          // kube-mesh LB IP
	MetaTelemetryLB string `json:"metaTelemetryLB,omitempty"` // todo rename to VPCLB
	MetaProxyHost   string `json:"metaProxyHost,omitempty"`   // TODO: refactor name
	MetaProxyEniLB  string `json:"metaProxyEniLB,omitempty"`
	KubeIP          string `json:"kubeIP,omitempty"`   // endpoint ip of svc kubernetes, like 169.X.X.X
	KubePort        int32  `json:"kubePort,omitempty"` // endpoint port of svc kubernetes

	TracingReportAddress string `json:"tracing_report_address,omitempty"`

	MornitorReportAddress string `json:"mornitor_report_address,omitempty"`

	// sensitive information
	Token        string `json:"token"`
	CaCert       string `json:"cacert"`
	CaKey        string `json:"cakey"`
	CertChain    string `json:"certchain"`
	TkeClusterCA string `json:"-"`
}
