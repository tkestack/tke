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
	"net"

	"k8s.io/apiserver/pkg/authentication/request/anonymous"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	apiserver "k8s.io/apiserver/pkg/server"
	restclient "k8s.io/client-go/rest"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/cmd/tke-registry-controller/app/options"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
	controlleroptions "tkestack.io/tke/pkg/controller/options"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	registryconfigvalidation "tkestack.io/tke/pkg/registry/apis/config/validation"
	registryconfig "tkestack.io/tke/pkg/registry/config"
	"tkestack.io/tke/pkg/registry/config/configfiles"
	registrycontrollerconfig "tkestack.io/tke/pkg/registry/controller/config"
	"tkestack.io/tke/pkg/util/log"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	SecureServing *apiserver.SecureServingInfo
	// LoopbackClientConfig is a config for a privileged loopback connection
	LoopbackClientConfig *restclient.Config
	Authentication       apiserver.AuthenticationInfo
	Authorization        apiserver.AuthorizationInfo
	ServerName           string
	// the client only used for leader election
	LeaderElectionClient *versionedclientset.Clientset
	// the rest config for the registry apiserver
	RegistryAPIServerClientConfig *restclient.Config
	// the rest config for the business apiserver
	BusinessAPIServerClientConfig *restclient.Config
	// the rest config for the auth apiserver
	AuthAPIServerClientConfig *restclient.Config
	// the registry config for chartmuseum/image
	RegistryConfig    *registryconfigv1.RegistryConfiguration
	RepoConfiguration registryconfig.RepoConfiguration
	ChartGroupSetting registrycontrollerconfig.ChartGroupSetting

	Component controlleroptions.ComponentConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	if err := opts.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	registryConfig, err := options.NewRegistryConfiguration()
	if err != nil {
		log.Error("Failed create default registry configuration", log.Err(err))
		return nil, err
	}
	// load config file, if provided
	if configFile := opts.RegistryConfig; len(configFile) > 0 {
		registryConfig, err = configfiles.LoadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load registry configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}
	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := registryconfigvalidation.ValidateRegistryConfiguration(registryConfig); err != nil {
		log.Error("Failed to validate registry configuration", log.Err(err))
		return nil, err
	}
	v1RegistryConfig := &registryconfigv1.RegistryConfiguration{}
	if err := registryconfigv1.Convert_config_RegistryConfiguration_To_v1_RegistryConfiguration(registryConfig, v1RegistryConfig, nil); err != nil {
		log.Error("Failed to convert registry configuration", log.Err(err))
		return nil, err
	}

	registryAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.RegistryAPIClient)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to initialize client config of registry API server")
	}

	// optional
	businessAPIServerClientConfig, _, err := controllerconfig.BuildClientConfig(opts.BusinessAPIClient)
	if err != nil {
		return nil, err
	}

	// optional
	authAPIServerClientConfig, _, err := controllerconfig.BuildClientConfig(opts.AuthAPIClient)
	if err != nil {
		return nil, err
	}

	// shallow copy, do not modify the apiServerClientConfig.Timeout.
	config := *registryAPIServerClientConfig
	config.Timeout = opts.Component.LeaderElection.RenewDeadline
	leaderElectionClient := versionedclientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))

	controllerManagerConfig := &Config{
		ServerName:           serverName,
		LeaderElectionClient: leaderElectionClient,
		Authorization: apiserver.AuthorizationInfo{
			Authorizer: authorizerfactory.NewAlwaysAllowAuthorizer(),
		},
		Authentication: apiserver.AuthenticationInfo{
			Authenticator: anonymous.NewAuthenticator(),
		},
		BusinessAPIServerClientConfig: businessAPIServerClientConfig,
		RegistryAPIServerClientConfig: registryAPIServerClientConfig,
		AuthAPIServerClientConfig:     authAPIServerClientConfig,
		RegistryConfig:                v1RegistryConfig,
	}

	if err := opts.Component.ApplyTo(&controllerManagerConfig.Component); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&controllerManagerConfig.SecureServing, &controllerManagerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}
	if err := opts.Debug.ApplyTo(&controllerManagerConfig.Component.Debugging); err != nil {
		return nil, err
	}

	if err := opts.ChartGroupSettingOptions.ApplyTo(&controllerManagerConfig.ChartGroupSetting); err != nil {
		return nil, err
	}

	if err := (&opts.FeatureOptions.Repo).ApplyTo(&controllerManagerConfig.RepoConfiguration); err != nil {
		return nil, err
	}

	return controllerManagerConfig, nil
}
