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

	// may remove 存量 1.3.6, DO NOT USE ANYMORE
	MasterCluster  *ClusterConfig   `json:"masterCluster"`
	RemoteClusters []*ClusterConfig `json:"remoteClusters,omitempty"`

	// MTLS enables or disables global mTLS
	MTLS bool `json:"mtls,omitempty"`

	// ControlPlaneSecurityEnabled control plane services are communicating through mTLS
	ControlPlaneSecurityEnabled bool `json:"controlPlaneSecurityEnabled,omitempty"`

	// if policy is enabled, global.disablePolicyChecks has affect.
	DisablePolicyChecks bool `json:"disablePolicyChecks"`

	TraceSampling float32 `json:"traceSampling"`

	// If SDS is configured, mTLS certificates for the sidecars will be distributed
	// through the SecretDiscoveryService instead of using K8S secrets to mount the certificates
	SDS SDSConfiguration `json:"sds,omitempty"`

	// Pilot configuration options
	Pilot PilotConfiguration `json:"pilot,omitempty"`

	// Galley configuration options
	Galley GalleyConfiguration `json:"galley,omitempty"`

	// Telemetry configuration options
	Telemetry MixerConfiguration `json:"telemetry,omitempty"`

	// Policy configuration options
	Policy MixerConfiguration `json:"policy,omitempty"`

	// Proxy configuration options
	Proxy ProxyConfiguration `json:"proxy,omitempty"`

	// Proxy Init configuration options
	ProxyInit ProxyInitConfiguration `json:"proxyInit,omitempty"`

	// Enable pod disruption budget for the control plane, which is used to
	// ensure MeshConfig control plane components are gradually upgraded or recovered
	DefaultPodDisruptionBudget PDBConfiguration `json:"defaultPodDisruptionBudget,omitempty"`

	// Set the default behavior of the sidecar for handling outbound traffic from the application (ALLOW_ANY or REGISTRY_ONLY)
	OutboundTrafficPolicy OutboundTrafficPolicyConfiguration `json:"outboundTrafficPolicy,omitempty"`

	// MeshConfig monitor(adapter) configuration options
	MeshMonitorAdapter MeshMonitorAdapterConfiguration `json:"meshMonitorAdapter,omitempty"`
}

// SDSConfiguration defines Secret Discovery Service config options
type SDSConfiguration struct {
	// If set to true, mTLS certificates for the sidecars will be
	// distributed through the SecretDiscoveryService instead of using K8S secrets to mount the certificates.
	Disabled bool `json:"disabled,omitempty"`
	// Unix Domain Socket through which envoy communicates with NodeAgent SDS to get
	// key/cert for mTLS. Use secret-mount files instead of SDS if set to empty.
	UdsPath string `json:"udsPath,omitempty"`
	// If set to true, MeshConfig will inject volumes mount for k8s service account JWT,
	// so that K8s API server mounts k8s service account JWT to envoy container, which
	// will be used to generate key/cert eventually.
	// (prerequisite: https://kubernetes.io/docs/concepts/storage/volumes/#projected)
	UseTrustworthyJwt bool `json:"useTrustworthyJwt,omitempty"`
	// If set to true, envoy will fetch normal k8s service account JWT from '/var/run/secrets/kubernetes.io/serviceaccount/token'
	// (https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#accessing-the-api-from-a-pod)
	// and pass to sds server, which will be used to request key/cert eventually
	// this flag is ignored if UseTrustworthyJwt is set
	UseNormalJwt bool `json:"useNormalJwt,omitempty"`
}

// PilotConfiguration defines config options for Pilot
type PilotConfiguration struct {
	TraceSampling float32 `json:"traceSampling"`

	CommonConfiguration
}

// CitadelConfiguration defines config options for Citadel
type CitadelConfiguration struct {
	CommonConfiguration
	SelfSigned bool `json:"selfSigned,omitempty"`
}

// GalleyConfiguration defines config options for Galley
type GalleyConfiguration struct {
	CommonConfiguration
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
	GatewayConfiguration `protobuf:"bytes,21,opt,name=gatewayConfiguration"`
	ExistLBID            string     `json:"existLBID,omitempty" protobuf:"bytes,20,opt,name=resources"` // TCNP-TODO
	LBName               string     `json:"lbName,omitempty"`                                           // TCNP-TODO
	AccessType           AccessType `json:"accessType"`                                                 // TCNP-TODO
	Subnet               string     `json:"subnet,omitempty"`                                           // TCNP-TODO
}

type K8SIngressConfiguration struct {
	Disabled bool `json:"disabled,omitempty" protobuf:"varint,1,opt,name=disabled"`
}

// SidecarInjectorConfiguration defines config options for SidecarInjector
type SidecarInjectorConfiguration struct {
	CommonConfiguration
	// If true, sidecar injector will rewrite PodSpec for liveness
	// health check to redirect request to sidecar. This makes liveness check work
	// even when mTLS is enabled.
	RewriteAppHTTPProbe bool `json:"rewriteAppHTTPProbe,omitempty" protobuf:"varint,3,opt,name=rewriteAppHTTPProbe"`
}

// MixerConfiguration defines config options for Mixer
type MixerConfiguration struct {
	CommonConfiguration
}

// ProxyConfiguration defines config options for Proxy
type ProxyConfiguration struct {
	Image string `json:"image,omitempty" protobuf:"bytes,1,opt,name=image"`
	// If set to true, istio-proxy container will have privileged securityContext
	Privileged bool `json:"privileged,omitempty" protobuf:"varint,2,opt,name=privileged"`
	// If set, newly injected sidecars will have core dumps enabled.
	EnableCoreDump bool `json:"enableCoreDump,omitempty" protobuf:"varint,3,opt,name=enableCoreDump"`

	Resources *v1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,5,opt,name=resources"`
}

// ProxyInitConfiguration defines config options for Proxy Init containers
type ProxyInitConfiguration struct {
	Image string `json:"image,omitempty" protobuf:"bytes,1,opt,name=image"`
}

// PDBConfiguration holds Pod Disruption Budget related config options
type PDBConfiguration struct {
	Disabled bool `json:"disabled,omitempty" protobuf:"varint,1,opt,name=disabled"`
}

type OutboundTrafficPolicyConfiguration struct {
	Mode string `json:"mode,omitempty" protobuf:"bytes,1,opt,name=mode"`
}

type MeshMonitorAdapterConfiguration struct {
	CommonConfiguration
	Disabled    bool   `json:"disabled,omitempty" protobuf:"varint,10,opt,name=disabled"`
	Debug       bool   `json:"debug,omitempty" protobuf:"varint,5,opt,name=debug"`
	LogToStderr bool   `json:"logToStderr,omitempty" protobuf:"varint,6,opt,name=logToStderr"`
	LoggerDir   string `json:"loggerDir,omitempty" protobuf:"bytes,7,opt,name=loggerDir"`
	LoggerFile  string `json:"loggerFile,omitempty" protobuf:"bytes,8,opt,name=loggerFile"`
}

type CommonConfiguration struct {
	Image         string                   `json:"image,omitempty" protobuf:"bytes,1,opt,name=image"`
	ReplicaCount  int32                    `json:"replicaCount,omitempty" protobuf:"varint,2,opt,name=replicaCount"`
	MinReplicas   int32                    `json:"minReplicas,omitempty" protobuf:"varint,3,opt,name=minReplicas"`
	MaxReplicas   int32                    `json:"maxReplicas,omitempty" protobuf:"varint,4,opt,name=maxReplicas"`
	Resources     *v1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,5,opt,name=resources"`
	Metrics       []v2beta1.MetricSpec     `json:"metrics,omitempty" protobuf:"bytes,11,rep,name=metrics"`
	SelectedNodes []string                 `json:"selectedNodes,omitempty" protobuf:"bytes,12,rep,name=selectedNodes"`
}

type CommonServiceConfiguration struct {
	ServiceType        v1.ServiceType       `json:"serviceType,omitempty"`
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
