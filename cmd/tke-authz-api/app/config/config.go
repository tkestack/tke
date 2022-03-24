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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"tkestack.io/tke/api/authz"
	authzv1 "tkestack.io/tke/api/authz/v1"

	"tkestack.io/tke/pkg/apiserver/util"

	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/client-go/rest"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/cmd/tke-authz-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authorization"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	"tkestack.io/tke/pkg/apiserver/storage"
	"tkestack.io/tke/pkg/authz/apiserver"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Authz API"
)

// Config is the running configuration structure of the TKE autz apiserver.
type Config struct {
	ServerName                     string
	GenericAPIServerConfig         *genericapiserver.Config
	VersionedSharedInformerFactory versionedinformers.SharedInformerFactory
	StorageFactory                 *serverstorage.DefaultStorageFactory
	PlatformClient                 platformversionedclient.PlatformV1Interface
	DefaultPolicies                []*authzv1.Policy
	DefaultRoles                   []*authzv1.Role
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE authz apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	genericAPIServerConfig := genericapiserver.NewConfig(authz.Codecs)
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

	openapi.SetupOpenAPI(genericAPIServerConfig, generatedopenapi.GetOpenAPIDefinitions, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(authz.Codecs, authz.Scheme)
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
	platformAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.PlatformAPIClient)
	if err != nil {
		return nil, err
	}
	if !ok || platformAPIServerClientConfig == nil {
		return nil, fmt.Errorf("failed to initialize client config of platform API server")
	}
	platformClient, err := versionedclientset.NewForConfig(rest.AddUserAgent(platformAPIServerClientConfig, "tke-authz-api"))
	if err != nil {
		return nil, err
	}

	var policies []*authzv1.Policy
	if opts.Authz.DefaultPoliciesConfig != "" {
		content, err := ioutil.ReadFile(opts.Authz.DefaultPoliciesConfig)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &policies)
		if err != nil {
			return nil, err
		}
	}
	var roles []*authzv1.Role
	if opts.Authz.DefaultRolesConfig != "" {
		content, err := ioutil.ReadFile(opts.Authz.DefaultRolesConfig)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &roles)
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{
		ServerName:                     serverName,
		GenericAPIServerConfig:         genericAPIServerConfig,
		VersionedSharedInformerFactory: versionedInformers,
		StorageFactory:                 storageFactory,
		PlatformClient:                 platformClient.PlatformV1(),
		DefaultPolicies:                policies,
		DefaultRoles:                   roles,
	}
	return cfg, nil
}
