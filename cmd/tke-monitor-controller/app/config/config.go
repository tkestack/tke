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
	monitorapiconfig "tkestack.io/tke/cmd/tke-monitor-api/app/config"
	monitorapioptions "tkestack.io/tke/cmd/tke-monitor-api/app/options"
	"tkestack.io/tke/cmd/tke-monitor-controller/app/options"
	controllerconfig "tkestack.io/tke/pkg/controller/config"
	controlleroptions "tkestack.io/tke/pkg/controller/options"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	monitorconfigvalidation "tkestack.io/tke/pkg/monitor/apis/config/validation"
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
	// the rest config for the monitor apiserver
	MonitorAPIServerClientConfig *restclient.Config
	// the rest config for the business apiserver
	BusinessAPIServerClientConfig *restclient.Config
	// the rest config for the platform apiserver
	PlatformAPIServerClientConfig *restclient.Config
	Component                     controlleroptions.ComponentConfiguration
	MonitorConfig                 *monitorconfig.MonitorConfiguration
	Features                      *options.FeatureOptions
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	if err := opts.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	monitorAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.MonitorAPIClient)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to initialize client config of monitor API server")
	}

	// shallow copy, do not modify the apiServerClientConfig.Timeout.
	config := *monitorAPIServerClientConfig
	config.Timeout = opts.Component.LeaderElection.RenewDeadline
	leaderElectionClient := versionedclientset.NewForConfigOrDie(restclient.AddUserAgent(&config, "leader-election"))

	monitorConfig, err := monitorapioptions.NewMonitorConfiguration()
	if err != nil {
		log.Error("Failed create default monitor configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.MonitorConfig; len(configFile) > 0 {
		var defaultRetentionDays int = 15 //set default retention value as 15 days for influxdb
		monitorConfig, err = monitorapiconfig.LoadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load monitor configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}

		if monitorConfig.Storage.InfluxDB != nil && monitorConfig.Storage.InfluxDB.RetentionDays == nil {
			log.Info("don't set retention times in config, use the default one")
			monitorConfig.Storage.InfluxDB.RetentionDays = &defaultRetentionDays
		}
	}

	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := monitorconfigvalidation.ValidateMonitorConfiguration(monitorConfig); err != nil {
		log.Error("Failed to validate monitor configuration", log.Err(err))
		return nil, err
	}

	controllerManagerConfig := &Config{
		ServerName:                   serverName,
		LeaderElectionClient:         leaderElectionClient,
		MonitorAPIServerClientConfig: monitorAPIServerClientConfig,
		Authorization: apiserver.AuthorizationInfo{
			Authorizer: authorizerfactory.NewAlwaysAllowAuthorizer(),
		},
		Authentication: apiserver.AuthenticationInfo{
			Authenticator: anonymous.NewAuthenticator(),
		},
		MonitorConfig: monitorConfig,
		Features:      opts.FeatureOptions,
	}

	businessAPIServerClientConfig, ok, err := controllerconfig.BuildClientConfig(opts.BusinessAPIClient)
	if err != nil {
		return nil, err
	}
	if ok && businessAPIServerClientConfig != nil {
		controllerManagerConfig.BusinessAPIServerClientConfig = businessAPIServerClientConfig
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
