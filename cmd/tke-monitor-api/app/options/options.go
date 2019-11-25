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

package options

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	genericapiserveroptions "k8s.io/apiserver/pkg/server/options"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
	controlleroptions "tkestack.io/tke/pkg/controller/options"
	"tkestack.io/tke/pkg/util/cachesize"
	"tkestack.io/tke/pkg/util/log"
)

const (
	flagMonitorConfig   = "monitor-config"
	configMonitorConfig = "monitor_config"
)

// Options is the main context object for the TKE monitor.
type Options struct {
	Log               *log.Options
	SecureServing     *apiserveroptions.SecureServingOptions
	Debug             *apiserveroptions.DebugOptions
	ETCD              *storageoptions.ETCDStorageOptions
	Generic           *apiserveroptions.GenericOptions
	Authentication    *apiserveroptions.AuthenticationWithAPIOptions
	Authorization     *apiserveroptions.AuthorizationOptions
	PlatformAPIClient *controlleroptions.APIServerClientOptions
	// The Registry will load its initial configuration from this file.
	// The path may be absolute or relative; relative paths are under the Monitor's current working directory.
	MonitorConfig string
}

// NewOptions creates a new Options with a default config.
func NewOptions(serverName string) *Options {
	return &Options{
		Log:               log.NewOptions(),
		Debug:             apiserveroptions.NewDebugOptions(),
		ETCD:              storageoptions.NewETCDStorageOptions("/tke/monitor"),
		SecureServing:     apiserveroptions.NewSecureServingOptions(serverName, 9455),
		Generic:           apiserveroptions.NewGenericOptions(),
		Authentication:    apiserveroptions.NewAuthenticationWithAPIOptions(),
		Authorization:     apiserveroptions.NewAuthorizationOptions(),
		PlatformAPIClient: controlleroptions.NewAPIServerClientOptions("platform"),
	}
}

// AddFlags adds flags for a specific server to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)
	o.Debug.AddFlags(fs)
	o.ETCD.AddFlags(fs)
	o.Generic.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.Authentication.AddFlags(fs)
	o.Authorization.AddFlags(fs)
	o.PlatformAPIClient.AddFlags(fs)

	fs.String(flagMonitorConfig, o.MonitorConfig,
		"The Monitor will load its initial configuration from this file. The path may be absolute or relative; relative paths start at the Monitor's current working directory. Omit this flag to use the built-in default configuration values.")
	_ = viper.BindPFlag(configMonitorConfig, fs.Lookup(flagMonitorConfig))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.Log.ApplyFlags()...)
	errs = append(errs, o.Debug.ApplyFlags()...)
	errs = append(errs, o.ETCD.ApplyFlags()...)
	errs = append(errs, o.Generic.ApplyFlags()...)
	errs = append(errs, o.SecureServing.ApplyFlags()...)
	errs = append(errs, o.Authentication.ApplyFlags()...)
	errs = append(errs, o.Authorization.ApplyFlags()...)
	errs = append(errs, o.PlatformAPIClient.ApplyFlags()...)

	o.MonitorConfig = viper.GetString(configMonitorConfig)

	return errs
}

// Complete set default Options.
// Should be called after tke-monitor flags parsed.
func (o *Options) Complete() error {
	if err := apiserveroptions.CompleteGenericAndSecureOptions(o.Generic, o.SecureServing); err != nil {
		return err
	}

	if o.ETCD.EnableWatchCache {
		log.Infof("Initializing cache sizes based on %dMB limit", o.Generic.TargetRAMMB)
		sizes := cachesize.NewHeuristicWatchCacheSizes(o.Generic.TargetRAMMB)
		if userSpecified, err := genericapiserveroptions.ParseWatchCacheSizes(o.ETCD.WatchCacheSizes); err == nil {
			for resource, size := range userSpecified {
				sizes[resource] = size
			}
		}

		watchCacheSizes, err := genericapiserveroptions.WriteWatchCacheSizes(sizes)
		if err != nil {
			return err
		}
		o.ETCD.WatchCacheSizes = watchCacheSizes
	}
	return nil
}
