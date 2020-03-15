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

package audit

import (
	"golang.org/x/oauth2"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/pkg/audit/api"
	auditconfig "tkestack.io/tke/pkg/audit/apis/config"
	"tkestack.io/tke/pkg/util/log"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName    string
	OAuthConfig   *oauth2.Config
	AuditConfig   *auditconfig.AuditConfiguration
	HeaderRequest bool
}

// Config contains the core configuration instance of server and additional
// configuration.
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embed a private pointer of Config.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Audit contains state for TKE audit server.
type Audit struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	return CompletedConfig{&c}
}

// New returns a new instance of APIServer from the given config.
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*Audit, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	if err := api.RegisterRoute(s.Handler.GoRestfulContainer, c.ExtraConfig.AuditConfig); err != nil {
		return nil, err
	}

	m := &Audit{
		GenericAPIServer: s,
	}

	log.Info("All of http handlers registered", log.Strings("paths", m.GenericAPIServer.Handler.ListedPaths()))
	return m, nil
}
