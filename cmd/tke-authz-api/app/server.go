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
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/cmd/tke-authz-api/app/config"
	"tkestack.io/tke/pkg/authz/apiserver"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
)

// CreateServerChain creates the api servers connected via delegation.
func CreateServerChain(cfg *config.Config) (*genericapiserver.GenericAPIServer, error) {
	apiServerConfig := createAPIServerConfig(cfg)
	apiServer, err := CreateAPIServer(apiServerConfig, genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if err := registerHandler(apiServer); err != nil {
		return nil, err
	}

	apiServer.GenericAPIServer.AddPostStartHookOrDie("start-authz-api-server-informers", func(ctx genericapiserver.PostStartHookContext) error {
		cfg.VersionedSharedInformerFactory.Start(ctx.StopCh)
		return nil
	})
	apiServer.GenericAPIServer.AddPostStartHookOrDie("init-default-policies", func(ctx genericapiserver.PostStartHookContext) error {
		log.Infof("init default policies ...")
		client, err := versionedclientset.NewForConfig(ctx.LoopbackClientConfig)
		if err != nil {
			log.Warnf("failed to generate authz client, err '%#v'", err)
			return err
		}
		for _, pol := range cfg.DefaultPolicies {
			if _, err := client.AuthzV1().Policies(pol.Namespace).Create(context.TODO(), pol, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
				log.Warnf("failed to init policy '%s/%s', err '%#v'", pol.Namespace, pol.Name, err)
				return err
			}
		}
		return nil
	})
	apiServer.GenericAPIServer.AddPostStartHookOrDie("init-default-roles", func(ctx genericapiserver.PostStartHookContext) error {
		log.Infof("init default roles ...")
		client, err := versionedclientset.NewForConfig(ctx.LoopbackClientConfig)
		if err != nil {
			log.Warnf("failed to generate authz client, err '%#v'", err)
			return err
		}
		for _, rol := range cfg.DefaultRoles {
			if _, err := client.AuthzV1().Roles(rol.Namespace).Create(context.TODO(), rol, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
				log.Warnf("failed to init role '%s/%s', err '%#v'", rol.Namespace, rol.Name, err)
				return err
			}
		}
		return nil
	})

	return apiServer.GenericAPIServer, nil
}

// CreateAPIServer creates and wires a workable tke-business-api
func CreateAPIServer(apiServerConfig *apiserver.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apiserver.APIServer, error) {
	return apiServerConfig.Complete().New(delegateAPIServer)
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
			PlatformClient:          cfg.PlatformClient,
		},
	}
}

func createFilterChain(apiServer *genericapiserver.GenericAPIServer) {
	apiServer.Handler.FullHandlerChain = filter.WithFuzzyResource(apiServer.Handler.FullHandlerChain)
	apiServer.Handler.FullHandlerChain = filter.WithCluster(apiServer.Handler.FullHandlerChain)
}

func registerHandler(apiServer *apiserver.APIServer) error {
	createFilterChain(apiServer.GenericAPIServer)
	log.Info("All of http handlers registered", log.Strings("paths", apiServer.GenericAPIServer.Handler.ListedPaths()))
	return nil
}
