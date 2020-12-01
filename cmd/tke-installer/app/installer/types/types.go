/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package types

import (
	"context"
	"path"

	"github.com/thoas/go-funk"

	v1 "tkestack.io/tke/api/platform/v1"
)

type CreateClusterPara struct {
	Cluster v1.Cluster `json:"cluster"`
	Config  Config     `json:"Config"`
}

// Config is the installer config
type Config struct {
	Basic       *Basic       `json:"basic"`
	Auth        Auth         `json:"auth"`
	Registry    Registry     `json:"registry"`
	Business    *Business    `json:"business,omitempty"`
	Monitor     *Monitor     `json:"monitor,omitempty"`
	Logagent    *Logagent    `json:"logagent,omitempty"`
	HA          *HA          `json:"ha,omitempty"`
	Gateway     *Gateway     `json:"gateway,omitempty"`
	Audit       *Audit       `json:"audit,omitempty"`
	Application *Application `json:"application,omitempty"`
	SkipSteps   []string     `json:"skipSteps,omitempty"`
}

type Basic struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Auth struct {
	TKEAuth  *TKEAuth  `json:"tke,omitempty"`
	OIDCAuth *OIDCAuth `json:"oidc,omitempty"`
}

type TKEAuth struct {
	TenantID string `json:"tenantID"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type OIDCAuth struct {
	IssuerURL string `json:"issuerURL" validate:"required"`
	ClientID  string `json:"clientID" validate:"required"`
	CACert    []byte `json:"caCert"`
}

// Registry for remote registry
type Registry struct {
	TKERegistry        *TKERegistry        `json:"tke,omitempty"`
	ThirdPartyRegistry *ThirdPartyRegistry `json:"thirdParty,omitempty"`
	UserInputRegistry  UserInputRegistry   `json:"userInputRegistry,omitempty"`
}

type Audit struct {
	ElasticSearch *ElasticSearch `json:"elasticSearch,omitempty"`
}

type ElasticSearch struct {
	Address     string `json:"address" validate:"required"`
	ReserveDays int    `json:"reserveDays" validate:"required"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func (r *Registry) Domain() string {
	if r.UserInputRegistry.Domain != "" {
		return r.UserInputRegistry.Domain
	}
	if r.ThirdPartyRegistry != nil { // first use third party when both set
		return r.ThirdPartyRegistry.Domain
	}
	return r.TKERegistry.Domain
}

func (r *Registry) Namespace() string {
	if r.UserInputRegistry.Namespace != "" {
		return r.UserInputRegistry.Namespace
	}
	if r.ThirdPartyRegistry != nil {
		return r.ThirdPartyRegistry.Namespace
	}
	return r.TKERegistry.Namespace
}

func (r *Registry) Username() string {
	if r.UserInputRegistry.Username != "" {
		return r.UserInputRegistry.Username
	}
	if r.ThirdPartyRegistry != nil {
		return r.ThirdPartyRegistry.Username
	}
	return r.TKERegistry.Username
}

func (r *Registry) Password() []byte {
	if len(r.UserInputRegistry.Password) != 0 {
		return r.UserInputRegistry.Password
	}
	if r.ThirdPartyRegistry != nil {
		return r.ThirdPartyRegistry.Password
	}
	return r.TKERegistry.Password
}

func (r *Registry) Prefix() string {
	return path.Join(r.Domain(), r.Namespace())
}

func (r *Registry) IsOfficial() bool {
	return funk.ContainsString([]string{"docker.io/tkestack", "ccr.ccs.tencentyun.com/tkestack"}, r.Prefix())
}

type TKERegistry struct {
	Domain        string `json:"domain" validate:"hostname_rfc1123"`
	HarborEnabled bool   `json:"harborEnabled"`
	HarborCAFile  string `json:"harborCAFile"`
	Namespace     string `json:"namespace"`
	Username      string `json:"username"`
	Password      []byte `json:"password"`
}

type ThirdPartyRegistry struct {
	Domain    string `json:"domain" validate:"required"`
	Namespace string `json:"namespace" validate:"required"`
	Username  string `json:"username"`
	Password  []byte `json:"password"`
}

type UserInputRegistry struct {
	Domain    string `json:"domain"`
	Namespace string `json:"namespace"`
	Username  string `json:"username"`
	Password  []byte `json:"password"`
}

type Business struct {
}

type Application struct {
	RegistryDomain   string `json:"registryDomain" validate:"hostname_rfc1123"`
	RegistryUsername string `json:"registryUsername"`
	RegistryPassword []byte `json:"registryPassword"`
}

type Monitor struct {
	ESMonitor       *ESMonitor       `json:"es,omitempty"`
	InfluxDBMonitor *InfluxDBMonitor `json:"influxDB,omitempty"`
}

type ESMonitor struct {
	URL      string `json:"url" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type InfluxDBMonitor struct {
	LocalInfluxDBMonitor    *LocalInfluxDBMonitor    `json:"local,omitempty"`
	ExternalInfluxDBMonitor *ExternalInfluxDBMonitor `json:"external,omitempty"`
}

type LocalInfluxDBMonitor struct {
}

type ExternalInfluxDBMonitor struct {
	URL      string `json:"url" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type HA struct {
	TKEHA        *TKEHA        `json:"tke,omitempty"`
	ThirdPartyHA *ThirdPartyHA `json:"thirdParty,omitempty"`
}

func (ha *HA) VIP() string {
	if ha.TKEHA != nil {
		return ha.TKEHA.VIP
	}
	return ha.ThirdPartyHA.VIP
}

type TKEHA struct {
	VIP string `json:"vip" validate:"required"`
}

type ThirdPartyHA struct {
	VIP   string `json:"vip" validate:"required"`
	VPort int32  `json:"vport"`
}

type Gateway struct {
	Domain string `json:"domain"`
	Cert   *Cert  `json:"cert"`
}

type Cert struct {
	SelfSignedCert *SelfSignedCert `json:"selfSigned,omitempty"`
	ThirdPartyCert *ThirdPartyCert `json:"thirdParty,omitempty"`
}

type SelfSignedCert struct {
}

type ThirdPartyCert struct {
	Certificate []byte `json:"certificate" validate:"required"`
	PrivateKey  []byte `json:"privateKey" validate:"required"`
}

type Keepalived struct {
	VIP string `json:"vip,omitempty"`
}

type ClusterProgress struct {
	Status     ClusterProgressStatus `json:"status"`
	Data       string                `json:"data"`
	URL        string                `json:"url,omitempty"`
	Username   string                `json:"username,omitempty"`
	Password   []byte                `json:"password,omitempty"`
	CACert     []byte                `json:"caCert,omitempty"`
	Hosts      []string              `json:"hosts,omitempty"`
	Servers    []string              `json:"servers,omitempty"`
	Kubeconfig []byte                `json:"kubeconfig,omitempty"`
}

type ClusterProgressStatus string

const (
	StatusUnknown  = "Unknown"
	StatusDoing    = "Doing"
	StatusSuccess  = "Success"
	StatusFailed   = "Failed"
	StatusRetrying = "Retrying"
)

type Handler struct {
	Name string
	Func func(ctx context.Context) error
}

type Logagent struct {
	RegistryDomain    string `json:"domain,omitempty"`
	RegistryNamespace string `json:"namespace:omitempty"`
}
