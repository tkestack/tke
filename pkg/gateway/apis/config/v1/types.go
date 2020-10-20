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
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GatewayConfiguration contains the configuration for the Gateway
type GatewayConfiguration struct {
	metav1.TypeMeta

	// disableOIDCProxy, by default, the gateway server will proxy access requests to
	// the OIDC server. This switch is used to disable this feature.
	DisableOIDCProxy bool `json:"disableOIDCProxy"`
	// components is used to wrap all the backend component settings in the TKE.
	Components Components `json:"components"`
	Registry   *Registry  `json:"registry,omitempty"`
	Auth       *Auth      `json:"auth,omitempty"`
}

type Components struct {
	// platform is used to specify the access information of the `tke-platform-api`
	// backend service.
	// +optional
	Platform *Component `json:"platform,omitempty"`
	// business is used to specify the access information of the `tke-business-api`
	// backend service.
	// +optional
	Business *Component `json:"business,omitempty"`
	// notify is used to specify the access information of the `tke-notify-api`
	// backend service.
	// +optional
	Notify *Component `json:"notify,omitempty"`
	// monitor is used to specify the access information of the `tke-monitor-api`
	// backend service.
	// +optional
	Monitor *Component `json:"monitor,omitempty"`
	// auth is used to specify the access information of the `tke-auth`
	// backend service.
	// +optional
	Auth *Component `json:"auth,omitempty"`
	// registry is used to specify the access information of the `tke-registry`
	// backend service.
	// +optional
	Registry *Component `json:"registry,omitempty"`
	// logagent is used to specify the access information of the `tke-logagent-api`
	// backend service.
	// +optional
	LogAgent *Component `json:"logagent,omitempty"`
	// audit is used to specify the access information of the `tke-audit-api`
	// backend service.
	// +optional
	Audit *Component `json:"audit,omitempty"`
	// application is used to specify the access information of the `tke-application-api`
	// backend service.
	// +optional
	Application *Component `json:"application,omitempty"`
}

type Component struct {
	// address indicates the access address of the backend component. If it is deployed
	// in the cluster, it can be the address and port of the service.
	Address string `json:"address"`
	// frontProxy indicates that the access credentials are resolved
	// before the proxy to the backend service, and the user identity is passed to the
	// backend through the header.
	FrontProxy *FrontProxyComponent `json:"frontProxy"`
	// passthrough indicates that the credentials are passed directly
	// when the proxy requests to the backend service.
	Passthrough *PassthroughComponent `json:"passthrough"`
}

type FrontProxyComponent struct {
	// caFile is the path to a PEM-encoded certificate bundle. Trusted root certificates
	// for server.
	// +optional
	CAFile string `json:"caFile,omitempty"`
	// clientCertFile is the path to a PEM-encoded certificate bundle. If the authentication
	// is in `FrontProxy` mode, you must develop a trusted client access certificate for
	// the backend service.
	ClientCertFile string `json:"clientCertFile"`
	// clientKeyFile is the path to a PEM-encoded private key bundle. If the authentication
	// is in `FrontProxy` mode, you must develop a trusted client access private key for
	// the backend service.
	ClientKeyFile string `json:"clientKeyFile"`
	// usernameHeader is request header to inspect for username.
	// X-Remote-User is suggested.
	UsernameHeader string `json:"usernameHeader"`
	// groupsHeader is request header to inspect for groups.
	// X-Remote-Groups is suggested.
	GroupsHeader string `json:"groupsHeader"`
	// extraPrefixHeader is request header prefixes to inspect.
	// X-Remote-Extra- is suggested.
	ExtraPrefixHeader string `json:"extraPrefixHeader"`
}

type PassthroughComponent struct {
	// caFile is the path to a PEM-encoded certificate bundle. Trusted root certificates
	// for server.
	// +optional
	CAFile string `json:"caFile,omitempty"`
}

type Registry struct {
	DefaultTenant string `json:"defaultTenant"`
	// +optional
	DomainSuffix string `json:"domainSuffix"`
}

type Auth struct {
	DefaultTenant string `json:"defaultTenant"`
}
