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
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	business "tkestack.io/tke/api/business"
	businessv1 "tkestack.io/tke/api/business/v1"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	"tkestack.io/tke/cmd/tke-business-api/app/options"
	"tkestack.io/tke/pkg/apiserver/storage"
	businessrest "tkestack.io/tke/pkg/business/registry/rest"
	"tkestack.io/tke/pkg/util/log"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory
	PlatformClient          platformversionedclient.PlatformV1Interface
	RegistryClient          registryversionedclient.RegistryV1Interface
	AuthClient              authversionedclient.AuthV1Interface
	PrivilegedUsername      string
	FeatureOptions          *options.FeatureOptions
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

	// The order here is preserved in discovery.
	restStorageProviders := []storage.RESTStorageProvider{
		&businessrest.StorageProvider{
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			PlatformClient:       c.ExtraConfig.PlatformClient,
			RegistryClient:       c.ExtraConfig.RegistryClient,
			AuthClient:           c.ExtraConfig.AuthClient,
			PrivilegedUsername:   c.ExtraConfig.PrivilegedUsername,
			Features:             c.ExtraConfig.FeatureOptions,
		},
	}
	m.InstallAPIs(c.ExtraConfig, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)
	m.GenericAPIServer.AddPostStartHookOrDie("default-administrator", c.postStartHookFunc())

	return m, nil
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *APIServer) InstallAPIs(extraConfig *ExtraConfig, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...storage.RESTStorageProvider) {
	var apiGroupsInfo []genericapiserver.APIGroupInfo

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !extraConfig.APIResourceConfigSource.AnyVersionForGroupEnabled(groupName) {
			log.Infof("Skipping disabled API group %q.", groupName)
			continue
		}
		apiGroupInfo, enabled := restStorageBuilder.NewRESTStorage(extraConfig.APIResourceConfigSource, restOptionsGetter)
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
		businessv1.SchemeGroupVersion,
	)
	return ret
}

func (c completedConfig) postStartHookFunc() genericapiserver.PostStartHookFunc {
	return func(context genericapiserver.PostStartHookContext) error {
		client := businessinternalclient.NewForConfigOrDie(c.GenericConfig.LoopbackClientConfig)

		tenant := c.ExtraConfig.FeatureOptions.TenantOfInitialAdministrator
		user := c.ExtraConfig.FeatureOptions.UserOfInitialAdministrator

		_, err := client.Platforms().Get(options.DefaultPlatform, metaV1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Errorf("addAdministrator(tenant:%s, user:%s) failed, for %s", tenant, user, err)
			return err
		}

		_, err = client.Platforms().Create(&business.Platform{
			ObjectMeta: metaV1.ObjectMeta{
				Name: options.DefaultPlatform,
			},
			Spec: business.PlatformSpec{
				TenantID:       tenant,
				Administrators: []string{user},
			},
		})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				log.Infof("addAdministrator(tenant:%s, user:%s) found %s", tenant, user, options.DefaultPlatform)
				return nil
			}
			log.Errorf("addAdministrator(tenant:%s, user:%s) failed, for %s", tenant, user, err)
			return err
		}

		log.Infof("addAdministrator(tenant:%s, user:%s) created %s", tenant, user, options.DefaultPlatform)
		return err
	}
}
