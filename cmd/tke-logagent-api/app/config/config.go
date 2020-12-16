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
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/common"

	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	"tkestack.io/tke/api/logagent"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/cmd/tke-logagent-api/app/options"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	"tkestack.io/tke/pkg/apiserver/storage"
	"tkestack.io/tke/pkg/apiserver/util"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
	"tkestack.io/tke/pkg/logagent/apiserver"
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
}

//config relies on options
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	genericAPIServerConfig := genericapiserver.NewConfig(logagent.Codecs)
	//to support file download we need to change logrun
	genericAPIServerConfig.LongRunningFunc = filter.LongRunningRequestCheck(sets.NewString("watch", "proxy", "connect"),
		sets.NewString("filedownload"), []string{})
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(nil, nil, nil)
	genericAPIServerConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()
	genericAPIServerConfig.EnableIndex = false
	genericAPIServerConfig.EnableProfiling = false

	if err := util.SetupAuditConfig(genericAPIServerConfig, opts.Audit); err != nil {
		return nil, err
	}
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
		return result
	}, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	kubeClientConfig := genericAPIServerConfig.LoopbackClientConfig
	clientgoExternalClient, err := versionedclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create real external clientset: %v", err)
	}
	versionedInformers := versionedinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	debug.SetupDebug(genericAPIServerConfig, opts.Debug)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(logagent.Codecs, logagent.Scheme)
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
	// client config for platform apiserver
	platformAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.PlatformAPIClient)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to initialize client config of platform API server")
	}
	platformClient, err := versionedclientset.NewForConfig(rest.AddUserAgent(platformAPIServerClientConfig, "tke-logagent-api"))
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
	}, nil
}
