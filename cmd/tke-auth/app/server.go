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

package app

import (
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/cmd/tke-auth/app/config"
	"tkestack.io/tke/pkg/auth"
)

// CreateServerChain creates the auth connected via delegation.
func CreateServerChain(cfg *config.Config) (*genericapiserver.GenericAPIServer, error) {
	authConfig := createAuthConfig(cfg)
	authServer, err := CreateAuth(authConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	return authServer.GenericAPIServer, nil
}

// CreateAuth creates and wires a workable tke-auth.
func CreateAuth(authConfig *auth.Config, delegateAPIServer genericapiserver.DelegationTarget) (*auth.Auth, error) {
	return authConfig.Complete().New(delegateAPIServer)
}

func createAuthConfig(cfg *config.Config) *auth.Config {
	return &auth.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config: *cfg.GenericAPIServerConfig,
		},
		ExtraConfig: auth.ExtraConfig{
			ServerName:          cfg.ServerName,
			OIDCExternalAddress: cfg.OIDCExternalAddress,
			CasbinEnforcer:      cfg.CasbinEnforcer,
			DexServer:           cfg.DexServer,
			Registry:            cfg.Registry,
			TokenAuthn:          cfg.TokenAuthn,
			APIKeyAuthn:         cfg.APIKeyAuthn,
			Authorizer:          cfg.Authorizer,
			PolicyFile:          cfg.PolicyFile,
			CategoryFile:        cfg.CategoryFile,
			TenantAdmin:         cfg.TenantAdmin,
			TenantAdminSecret:   cfg.TenantAdminSecret,
		},
	}
}
