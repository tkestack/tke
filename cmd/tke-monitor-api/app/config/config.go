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
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"k8s.io/kube-openapi/pkg/common"
	"path/filepath"
	"time"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	"tkestack.io/tke/api/monitor"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/cmd/tke-monitor-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authorization"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	"tkestack.io/tke/pkg/apiserver/storage"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	monitorconfigvalidation "tkestack.io/tke/pkg/monitor/apis/config/validation"
	"tkestack.io/tke/pkg/monitor/apiserver"
	"tkestack.io/tke/pkg/monitor/config/configfiles"
	monitoropenapi "tkestack.io/tke/pkg/monitor/openapi"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	"tkestack.io/tke/pkg/util/log"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Monitor API"
)

// Config is the running configuration structure of the TKE monitor.
type Config struct {
	ServerName                     string
	GenericAPIServerConfig         *genericapiserver.Config
	VersionedSharedInformerFactory versionedinformers.SharedInformerFactory
	StorageFactory                 *serverstorage.DefaultStorageFactory
	PrivilegedUsername             string
	PlatformClient                 platformversionedclient.PlatformV1Interface
	MonitorConfig                  *monitorconfig.MonitorConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE monitor command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	monitorConfig, err := options.NewMonitorConfiguration()
	if err != nil {
		log.Error("Failed create default monitor configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.MonitorConfig; len(configFile) > 0 {
		monitorConfig, err = LoadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load monitor configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}

	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := monitorconfigvalidation.ValidateMonitorConfiguration(monitorConfig); err != nil {
		klog.Fatal(err)
	}

	genericAPIServerConfig := genericapiserver.NewConfig(monitor.Codecs)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(nil)
	genericAPIServerConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()
	genericAPIServerConfig.EnableIndex = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, func(callback common.ReferenceCallback) map[string]common.OpenAPIDefinition {
		result := make(map[string]common.OpenAPIDefinition)
		generated := generatedopenapi.GetOpenAPIDefinitions(callback)
		for k, v := range generated {
			result[k] = v
		}
		customs := monitoropenapi.GetOpenAPIDefinitions(callback)
		for k, v := range customs {
			result[k] = v
		}
		return result
	}, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(monitor.Codecs, monitor.Scheme)
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

	// client config for platform apiserver
	platformAPIServerClientConfig, err := controllerconfig.BuildClientConfig(opts.PlatformAPIClient)
	if err != nil {
		return nil, err
	}
	platformClient, err := versionedclientset.NewForConfig(rest.AddUserAgent(platformAPIServerClientConfig, "tke-monitor-api"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerName:                     serverName,
		GenericAPIServerConfig:         genericAPIServerConfig,
		VersionedSharedInformerFactory: versionedInformers,
		StorageFactory:                 storageFactory,
		PlatformClient:                 platformClient.PlatformV1(),
		PrivilegedUsername:             opts.Authentication.PrivilegedUsername,
		MonitorConfig:                  monitorConfig,
	}, nil
}

func LoadConfigFile(name string) (*monitorconfig.MonitorConfiguration, error) {
	const errFmt = "failed to load Monitor config file %s, error %v"
	// compute absolute path based on current working dir
	monitorConfigFile, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	loader, err := configfiles.NewFsLoader(utilfs.DefaultFs{}, monitorConfigFile)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	kc, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	return kc, err
}
