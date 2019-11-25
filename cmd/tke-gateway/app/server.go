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
	"tkestack.io/tke/cmd/tke-gateway/app/config"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/gateway"
)

// CreateServerChain creates the gateway connected via delegation.
func CreateServerChain(cfg *config.Config, stopCh <-chan struct{}) (*genericapiserver.GenericAPIServer, error) {
	gatewayConfig := createGatewayConfig(cfg)
	gatewayServer, err := CreateGateway(gatewayConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if cfg.InsecureServingInfo != nil {
		chain := handler.BuildHandlerChain(nil)
		insecureHandlerChain := chain(gatewayServer.GenericAPIServer.UnprotectedHandler(), &gatewayConfig.GenericConfig.Config)
		if err := cfg.InsecureServingInfo.Serve(insecureHandlerChain, gatewayConfig.GenericConfig.RequestTimeout, stopCh); err != nil {
			return nil, err
		}
	}

	return gatewayServer.GenericAPIServer, nil
}

// CreateGateway creates and wires a workable tke-console.
func CreateGateway(gatewayConfig *gateway.Config, delegateAPIServer genericapiserver.DelegationTarget) (*gateway.Gateway, error) {
	return gatewayConfig.Complete().New(delegateAPIServer)
}

func createGatewayConfig(cfg *config.Config) *gateway.Config {
	return &gateway.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config: *cfg.GenericAPIServerConfig,
		},
		ExtraConfig: gateway.ExtraConfig{
			ServerName:        cfg.ServerName,
			OAuthConfig:       cfg.OAuthConfig,
			OIDCHttpClient:    cfg.OIDCHTTPClient,
			OIDCAuthenticator: cfg.OIDCAuthenticator,
			GatewayConfig:     cfg.GatewayConfig,
		},
	}
}
