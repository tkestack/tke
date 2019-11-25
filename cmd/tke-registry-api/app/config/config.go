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

package config

import (
	"fmt"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/klog"
	"path/filepath"
	"time"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/cmd/tke-registry-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authorization"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	"tkestack.io/tke/pkg/apiserver/storage"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryconfigvalidation "tkestack.io/tke/pkg/registry/apis/config/validation"
	"tkestack.io/tke/pkg/registry/apiserver"
	"tkestack.io/tke/pkg/registry/config/configfiles"
	"tkestack.io/tke/pkg/registry/distribution"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	"tkestack.io/tke/pkg/util/log"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Registry API"
)

// Config is the running configuration structure of the TKE registry apiserver.
type Config struct {
	ServerName                     string
	GenericAPIServerConfig         *genericapiserver.Config
	VersionedSharedInformerFactory versionedinformers.SharedInformerFactory
	StorageFactory                 *serverstorage.DefaultStorageFactory
	ExternalScheme                 string
	OIDCCAFile                     string
	OIDCTokenReviewPath            string
	OIDCIssuerURL                  string
	RegistryConfig                 *registryconfig.RegistryConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE registry apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	registryConfig, err := options.NewRegistryConfiguration()
	if err != nil {
		log.Error("Failed create default registry configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.RegistryConfig; len(configFile) > 0 {
		registryConfig, err = loadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load registry configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}

	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := registryconfigvalidation.ValidateRegistryConfiguration(registryConfig); err != nil {
		klog.Fatal(err)
	}

	genericAPIServerConfig := genericapiserver.NewConfig(registry.Codecs)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(distribution.IgnoreAuthPathPrefixes())
	genericAPIServerConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()
	genericAPIServerConfig.EnableIndex = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, generatedopenapi.GetOpenAPIDefinitions, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(registry.Codecs, registry.Scheme)
	storageFactoryConfig.APIResourceConfig = genericAPIServerConfig.MergedResourceConfig
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(opts.ETCD)
	if err != nil {
		return nil, err
	}
	storageFactory, err := completedStorageFactoryConfig.New()
	if err != nil {
		return nil, err
	}
	if err := opts.ETCD.ApplyWithStorageFactoryTo(storageFactory, genericAPIServerConfig); err != nil {
		return nil, err
	}

	// client config
	genericAPIServerConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"

	kubeClientConfig := genericAPIServerConfig.LoopbackClientConfig
	clientgoExternalClient, err := versionedclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create real external clientset: %v", err)
	}
	versionedInformers := versionedinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	debug.SetupDebug(genericAPIServerConfig, opts.Debug)

	if err := authentication.SetupAuthentication(genericAPIServerConfig, opts.Authentication); err != nil {
		return nil, err
	}

	if err := authorization.SetupAuthorization(genericAPIServerConfig, opts.Authorization); err != nil {
		return nil, err
	}

	return &Config{
		ServerName:                     serverName,
		GenericAPIServerConfig:         genericAPIServerConfig,
		VersionedSharedInformerFactory: versionedInformers,
		StorageFactory:                 storageFactory,
		ExternalScheme:                 opts.Generic.ExternalScheme,
		OIDCIssuerURL:                  opts.Authentication.OIDC.IssuerURL,
		OIDCCAFile:                     opts.Authentication.OIDC.CAFile,
		OIDCTokenReviewPath:            opts.Authentication.OIDC.TokenReviewPath,
		RegistryConfig:                 registryConfig,
	}, nil
}

func loadConfigFile(name string) (*registryconfig.RegistryConfiguration, error) {
	const errFmt = "failed to load Registry config file %s, error %v"
	// compute absolute path based on current working dir
	registryConfigFile, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	loader, err := configfiles.NewFsLoader(utilfs.DefaultFs{}, registryConfigFile)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	kc, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	return kc, err
}
