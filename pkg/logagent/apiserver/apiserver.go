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
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/pkg/apiserver/storage"
	"k8s.io/apiserver/pkg/registry/generic"
	"tkestack.io/tke/pkg/util/log"
	logagentv1 "tkestack.io/tke/api/logagent/v1"
	logagentrest "tkestack.io/tke/pkg/logagent/registry/rest"
)


// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory
	PlatformClient          platformversionedclient.PlatformV1Interface
	PrivilegedUsername      string
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

//TODO: add storage group
// New returns a new instance of APIServer from the given config.
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*APIServer, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	m := &APIServer{
		GenericAPIServer: s,
	}
	//add provider here
	restStorageProviders := []storage.RESTStorageProvider{
		&logagentrest.StorageProvider{
		LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
		PrivilegedUsername:   c.ExtraConfig.PrivilegedUsername,
		PlatformClient: c.ExtraConfig.PlatformClient,
		},
	}
	m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)
	log.Info("All of http handlers registered in server", log.Strings("paths", m.GenericAPIServer.Handler.ListedPaths()))
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
		log.Infof("install apigroup info here %+v", apiGroupsInfo[i])
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
		logagentv1.SchemeGroupVersion,
	)
	return ret
}
