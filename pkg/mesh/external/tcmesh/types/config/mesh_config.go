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

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"k8s.io/api/autoscaling/v2beta1"
	v1 "k8s.io/api/core/v1"
)

type MeshConfig struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`

	TenantID    string `json:"tenantID"`
	Version     string `json:"version"`
	Revision    string `json:"revision"`
	Region      string `json:"region"`
	DisplayName string `json:"displayName"`

	Mode     string `json:"mode"`
	Topology string `json:"topology"`

	Clusters []*ClusterConfig `json:"clusters"`

	// TODO： may remove， from 1.3.6 DO NOT USE ANYMORE
	MasterCluster  *ClusterConfig   `json:"masterCluster"`
	RemoteClusters []*ClusterConfig `json:"remoteClusters,omitempty"`

	// MTLS enables or disables global mTLS
	MTLS bool `json:"mtls,omitempty"`

	// ControlPlaneSecurityEnabled control plane services are communicating through mTLS
	ControlPlaneSecurityEnabled bool `json:"controlPlaneSecurityEnabled,omitempty"`

	// if policy is enabled, global.disablePolicyChecks has affect.
	DisablePolicyChecks bool `json:"disablePolicyChecks"`

	TraceSampling float32 `json:"traceSampling"`

	// Istiod configuration options
	Istiod IstiodConfiguration `json:"istiod,omitempty"`

	// Proxy configuration options
	Proxy ProxyConfiguration `json:"proxy,omitempty"`

	// Proxy Init configuration options
	ProxyInit ProxyInitConfiguration `json:"proxyInit,omitempty"`

	// Enable pod disruption budget for the control plane, which is used to
	// ensure MeshConfig control plane components are gradually upgraded or recovered
	DefaultPodDisruptionBudget PDBConfiguration `json:"defaultPodDisruptionBudget,omitempty"`

	// Set the default behavior of the sidecar for handling outbound traffic from the application (ALLOW_ANY or REGISTRY_ONLY)
	OutboundTrafficPolicy OutboundTrafficPolicyConfiguration `json:"outboundTrafficPolicy,omitempty"`

	// MeshTracing  configuration options
	MeshTracing MeshTracingConfiguration `json:"meshTracing,omitempty"`
}

// GatewaysConfiguration defines config options for Gateways
type GatewaysConfiguration struct {
	IngressConfigs []IngressGatewayConfiguration `json:"ingressGateways,omitempty"`
	EgressConfigs  []GatewayConfiguration        `json:"egressGateways,omitempty"`
}

type GatewaySDSConfiguration struct {
	Enabled bool   `json:"enabled,omitempty"`
	Image   string `json:"image,omitempty"`
}

type GatewayConfiguration struct {
	Name               string                  `json:"name"`
	Namespace          string                  `json:"namespace"`
	ServiceType        v1.ServiceType          `json:"serviceType,omitempty"`
	ServiceAnnotations map[string]string       `json:"serviceAnnotations,omitempty"`
	ServiceLabels      map[string]string       `json:"serviceLabels,omitempty"`
	SDS                GatewaySDSConfiguration `json:"sds,omitempty"`
	Inner              bool                    `json:"inner,omitempty"` // add label `managed-by: tke-mesh`, not editable

	CommonConfiguration
}

type IngressGatewayConfiguration struct {
	GatewayConfiguration
	ExistLBID  string     `json:"existLBID,omitempty"` // TCNP-TODO
	LBName     string     `json:"lbName,omitempty"`    // TCNP-TODO
	AccessType AccessType `json:"accessType"`          // TCNP-TODO
	Subnet     string     `json:"subnet,omitempty"`    // TCNP-TODO
}

type K8SIngressConfiguration struct {
	Disabled bool `json:"disabled,omitempty"`
}

// SidecarInjectorConfiguration defines config options for SidecarInjector
type SidecarInjectorConfiguration struct {
	CommonConfiguration
	// If true, sidecar injector will rewrite PodSpec for liveness
	// health check to redirect request to sidecar. This makes liveness check work
	// even when mTLS is enabled.
	RewriteAppHTTPProbe bool `json:"rewriteAppHTTPProbe,omitempty"`
}

// MixerConfiguration defines config options for Mixer
type MixerConfiguration struct {
	CommonConfiguration
}

// ProxyConfiguration defines config options for Proxy
type ProxyConfiguration struct {
	Image string `json:"image,omitempty"`
	// If set to true, istio-proxy container will have privileged securityContext
	Privileged bool `json:"privileged,omitempty"`
	// If set, newly injected sidecars will have core dumps enabled.
	EnableCoreDump bool `json:"enableCoreDump,omitempty"`

	Resources *v1.ResourceRequirements `json:"resources,omitempty"`
}

// ProxyInitConfiguration defines config options for Proxy Init containers
type ProxyInitConfiguration struct {
	Image string `json:"image,omitempty"`
}

// PDBConfiguration holds Pod Disruption Budget related config options
type PDBConfiguration struct {
	Disabled bool `json:"disabled,omitempty"`
}

type OutboundTrafficPolicyConfiguration struct {
	Mode string `json:"mode,omitempty"`
}

type MeshTracingConfiguration struct {
	CommonConfiguration
	Disabled    bool   `json:"disabled,omitempty"`
	Debug       bool   `json:"debug,omitempty"`
	LogToStderr bool   `json:"logToStderr,omitempty"`
	LoggerDir   string `json:"loggerDir,omitempty"`
	LoggerFile  string `json:"loggerFile,omitempty"`
}

type CommonConfiguration struct {
	Image         string                   `json:"image,omitempty"`
	ReplicaCount  int32                    `json:"replicaCount,omitempty"`
	MinReplicas   int32                    `json:"minReplicas,omitempty"`
	MaxReplicas   int32                    `json:"maxReplicas,omitempty"`
	Resources     *v1.ResourceRequirements `json:"resources,omitempty"`
	Metrics       []v2beta1.MetricSpec     `json:"metrics,omitempty"`
	SelectedNodes []string                 `json:"selectedNodes,omitempty"`
}

type CommonServiceConfiguration struct {
	ServiceType        v1.ServiceType    `json:"serviceType,omitempty"`
	ServiceAnnotations map[string]string `json:"serviceAnnotations,omitempty"`
	ServiceLabels      map[string]string `json:"serviceLabels,omitempty"`
}

type MeshKubeOperatorConfiguration struct {
	CommonConfiguration
}

type ComponentDeployMode struct {
	ProprietaryNodesDeploy bool              `json:"proprietaryNodesDeploy"`
	Labels                 map[string]string `json:"labels,omitempty"`
	Nodes                  []string          `json:"nodes,omitempty"`
}

// Scan implement sql.Scanner
func (c *MeshConfig) Scan(value interface{}) error {
	if bytes, ok := value.([]byte); ok {
		if len(bytes) == 0 { // todo
			return nil
		}
		return json.Unmarshal(bytes, c)
	}
	return errors.New(fmt.Sprint("Failed to unmarshal MeshConfig JSON from DB", value))
}

// Value implement driver.Valuer
func (c MeshConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}
