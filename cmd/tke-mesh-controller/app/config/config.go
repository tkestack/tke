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

package config

import (
	"fmt"
	"net"

	"k8s.io/apiserver/pkg/authentication/request/anonymous"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	apiserver "k8s.io/apiserver/pkg/server"
	restclient "k8s.io/client-go/rest"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	meshapiconfig "tkestack.io/tke/cmd/tke-mesh-api/app/config"
	meshapioptions "tkestack.io/tke/cmd/tke-mesh-api/app/options"
	"tkestack.io/tke/cmd/tke-mesh-controller/app/options"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
	controlleroptions "tkestack.io/tke/pkg/controller/options"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	meshconfigvalidation "tkestack.io/tke/pkg/mesh/apis/config/validation"
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
	// the rest config for the mesh apiserver
	MeshAPIServerClientConfig *restclient.Config
	// the rest config for the platform apiserver
	PlatformAPIServerClientConfig *restclient.Config
	Component                     controlleroptions.ComponentConfiguration
	MeshConfig                    *meshconfig.MeshConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	if err := opts.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	meshAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.MeshAPIClient)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to initialize client config of mesh API server")
	}

	// shallow copy, do not modify the apiServerClientConfig.Timeout.
	config := *meshAPIServerClientConfig
	config.Timeout = opts.Component.LeaderElection.RenewDeadline
	leaderElectionClient := versionedclientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))

	meshConfig, err := meshapioptions.NewMeshConfiguration()
	if err != nil {
		log.Error("Failed create default mesh configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.MeshConfig; len(configFile) > 0 {
		meshConfig, err = meshapiconfig.LoadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load mesh configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}

	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := meshconfigvalidation.ValidateConfiguration(meshConfig); err != nil {
		log.Error("Failed to validate mesh configuration", log.Err(err))
		return nil, err
	}

	controllerManagerConfig := &Config{
		ServerName:                serverName,
		LeaderElectionClient:      leaderElectionClient,
		MeshAPIServerClientConfig: meshAPIServerClientConfig,
		Authorization: apiserver.AuthorizationInfo{
			Authorizer: authorizerfactory.NewAlwaysAllowAuthorizer(),
		},
		Authentication: apiserver.AuthenticationInfo{
			Authenticator: anonymous.NewAuthenticator(),
		},
		MeshConfig: meshConfig,
	}

	platformAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.PlatformAPIClient)
	if err != nil {
		return nil, err
	}
	if ok && platformAPIServerClientConfig != nil {
		controllerManagerConfig.PlatformAPIServerClientConfig = platformAPIServerClientConfig
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
	return controllerManagerConfig, nil
}
