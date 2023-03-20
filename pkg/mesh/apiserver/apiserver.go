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

package apiserver

import (
	"k8s.io/kube-openapi/pkg/validation/spec"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	meshv1 "tkestack.io/tke/api/mesh/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	clusterclient "tkestack.io/tke/pkg/mesh/external/kubernetes"
	tcmeshclient "tkestack.io/tke/pkg/mesh/external/tcmesh"
	meshrest "tkestack.io/tke/pkg/mesh/registry/rest"
	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory
	PlatformClient          platformversionedclient.PlatformV1Interface
	ClusterClients          clusterclient.Client
	PrivilegedUsername      string
	MeshConfig              *meshconfig.MeshConfiguration
	TCMeshClient            *tcmeshclient.Client
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

	// restful api handlers
	r := Route{
		config: c,
	}
	s.Handler.GoRestfulContainer.Add(r.WebService())

	// The order here is preserved in discovery.
	restStorageProviders := []storage.RESTStorageProvider{
		&meshrest.StorageProvider{
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			PrivilegedUsername:   c.ExtraConfig.PrivilegedUsername,
		},
	}
	m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)

	// TCMesh clusters cache reload
	//tcmCleint := c.ExtraConfig.TCMeshClient
	//s.AddPostStartHookOrDie("mesh-cluster-cacher", func(ctx genericapiserver.PostStartHookContext) error {
	//	go wait.Until(tcmCleint.ReloadMeshes, 5*time.Minute, ctx.StopCh)
	//	return nil
	//})

	log.Info("All of http handlers registered", log.Strings("paths", m.GenericAPIServer.Handler.ListedPaths()))
	return m, nil
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Mesh API",
			Description: "Resource for Mesh",
			Version:     "1.0.0",
		},
	}
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *APIServer) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...storage.RESTStorageProvider) {
	var apiGroupsInfo []genericapiserver.APIGroupInfo

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyResourceForGroupEnabled(groupName) {
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
		meshv1.SchemeGroupVersion,
	)
	return ret
}
