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
	"tkestack.io/tke/cmd/tke-logagent-api/app/config"
	"tkestack.io/tke/pkg/logagent/apiserver"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// CreateServerChain creates the apiserver connected via delegation.
func CreateServerChain(cfg *config.Config) (*genericapiserver.GenericAPIServer, error) {
	apiServerConfig := createAPIServerConfig(cfg)
	apiServer, err := CreateAPIServer(apiServerConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if err := registerHandler(apiServer); err != nil {
		return nil, err
	}

	apiServer.GenericAPIServer.AddPostStartHookOrDie("start-logagent-api-server-informers", func(context genericapiserver.PostStartHookContext) error {
		cfg.VersionedSharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return apiServer.GenericAPIServer, nil
}

// CreateAPIServer creates and wires a workable tke-logagent.
func CreateAPIServer(logagentConfig *apiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apiserver.APIServer, error) {
	return logagentConfig.Complete().New(delegateAPIServer)
}

func createAPIServerConfig(cfg *config.Config) *apiserver.Config {
	return &apiserver.Config{
		GenericConfig: &genericapiserver.RecommendedConfig{
			Config: *cfg.GenericAPIServerConfig,
		},
		ExtraConfig: apiserver.ExtraConfig{
			ServerName:              cfg.ServerName,
			VersionedInformers:      cfg.VersionedSharedInformerFactory,
			StorageFactory:          cfg.StorageFactory,
			APIResourceConfigSource: cfg.StorageFactory.APIResourceConfigSource,
			PrivilegedUsername:      cfg.PrivilegedUsername,
			PlatformClient:          cfg.PlatformClient,
		},
	}
}

func createFilterChain(apiServer *genericapiserver.GenericAPIServer) {
	apiServer.Handler.FullHandlerChain = filter.WithCluster(apiServer.Handler.FullHandlerChain)
	apiServer.Handler.FullHandlerChain = filter.WithRequestBody(apiServer.Handler.FullHandlerChain)
	apiServer.Handler.FullHandlerChain = filter.WithFuzzyResource(apiServer.Handler.FullHandlerChain)
}

func registerHandler(apiServer *apiserver.APIServer) error {
	createFilterChain(apiServer.GenericAPIServer)
	log.Info("All of http handlers registered in cmd", log.Strings("paths", apiServer.GenericAPIServer.Handler.ListedPaths()))
	return nil
}
