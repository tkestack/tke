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

package apiserver

import (
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/chartmuseum"
	"tkestack.io/tke/pkg/registry/distribution"
	"tkestack.io/tke/pkg/registry/harbor"
	registryrest "tkestack.io/tke/pkg/registry/registry/rest"
	"tkestack.io/tke/pkg/util/log"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory
	ExternalScheme          string
	ExternalHost            string
	ExternalPort            int
	ExternalCAFile          string
	OIDCCAFile              string
	OIDCTokenReviewPath     string
	OIDCIssuerURL           string
	RegistryConfig          *registryconfig.RegistryConfiguration
	AuthClient              authversionedclient.AuthV1Interface
	BusinessClient          businessversionedclient.BusinessV1Interface
	PlatformClient          platformversionedclient.PlatformV1Interface
}

// Config contains the core configuration instance of apiserver and
// additional configuration.
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

// APIServer contains state for a tke api server.
type APIServer struct {
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
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*APIServer, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	m := &APIServer{
		GenericAPIServer: s,
	}
	if c.ExtraConfig.RegistryConfig.HarborEnabled {
		harborOpts := &harbor.Options{
			RegistryConfig:       c.ExtraConfig.RegistryConfig,
			ExternalHost: c.ExtraConfig.ExternalHost,
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
		}
		if err := harbor.RegisterRoute(s.Handler.NonGoRestfulMux, harborOpts); err != nil {
			return nil, err
		}
	} else {
		distributionOpts := &distribution.Options{
			RegistryConfig:       c.ExtraConfig.RegistryConfig,
			ExternalScheme:       c.ExtraConfig.ExternalScheme,
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			OIDCCAFile:           c.ExtraConfig.OIDCCAFile,
			OIDCTokenReviewPath:  c.ExtraConfig.OIDCTokenReviewPath,
			OIDCIssuerURL:        c.ExtraConfig.OIDCIssuerURL,
		}
		if err := distribution.RegisterRoute(s.Handler.NonGoRestfulMux, distributionOpts); err != nil {
			return nil, err
		}

		chartmuseumOpts := &chartmuseum.Options{
			RegistryConfig:       c.ExtraConfig.RegistryConfig,
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			OIDCCAFile:           c.ExtraConfig.OIDCCAFile,
			OIDCTokenReviewPath:  c.ExtraConfig.OIDCTokenReviewPath,
			OIDCIssuerURL:        c.ExtraConfig.OIDCIssuerURL,
			ExternalScheme:       c.ExtraConfig.ExternalScheme,
			Authorizer:           c.GenericConfig.Authorization.Authorizer,
		}
		if err := chartmuseum.RegisterRoute(s.Handler.NonGoRestfulMux, chartmuseumOpts); err != nil {
			return nil, err
		}
	}	

	// The order here is preserved in discovery.
	restStorageProviders := []storage.RESTStorageProvider{
		&registryrest.StorageProvider{
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			ExternalScheme:       c.ExtraConfig.ExternalScheme,
			ExternalHost:         c.ExtraConfig.ExternalHost,
			ExternalPort:         c.ExtraConfig.ExternalPort,
			ExternalCAFile:       c.ExtraConfig.ExternalCAFile,
			AuthClient:           c.ExtraConfig.AuthClient,
			BusinessClient:       c.ExtraConfig.BusinessClient,
			PlatformClient:       c.ExtraConfig.PlatformClient,
			RegistryConfig:       c.ExtraConfig.RegistryConfig,
			Authorizer:           c.GenericConfig.Authorization.Authorizer,
		},
	}
	m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)

	log.Info("All of http handlers registered", log.Strings("paths", m.GenericAPIServer.Handler.ListedPaths()))

	return m, nil
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *APIServer) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...storage.RESTStorageProvider) {
	var apiGroupsInfo []genericapiserver.APIGroupInfo

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyVersionForGroupEnabled(groupName) {
			log.Infof("Skipping disabled API group %q.", groupName)
			continue
		}
		apiGroupInfo, enabled := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if !enabled {
			log.Warnf("Problem initializing API group %q, skipping.", groupName)
			continue
		}
		log.Infof("Enabling API group %q.", groupName)

		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				log.Fatalf("Error building PostStartHook: %v", err)
			}
			m.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}

		apiGroupsInfo = append(apiGroupsInfo, apiGroupInfo)
	}

	for i := range apiGroupsInfo {
		if err := m.GenericAPIServer.InstallAPIGroup(&apiGroupsInfo[i]); err != nil {
			log.Fatalf("Error in registering group versions: %v", err)
		}
	}
}

// DefaultAPIResourceConfigSource returns which groupVersion enabled and its
// resources enabled/disabled.
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	ret.EnableVersions(
		registryv1.SchemeGroupVersion,
	)
	return ret
}
