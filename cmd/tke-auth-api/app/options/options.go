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
	"fmt"

	"github.com/spf13/pflag"
	genericapiserveroptions "k8s.io/apiserver/pkg/server/options"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
	"tkestack.io/tke/pkg/auth/apiserver"
	"tkestack.io/tke/pkg/util/cachesize"
	"tkestack.io/tke/pkg/util/log"
)

// Options is the main context object for the TKE auth.
type Options struct {
	Log            *log.Options
	SecureServing  *apiserveroptions.SecureServingOptions
	Debug          *apiserveroptions.DebugOptions
	Generic        *apiserveroptions.GenericOptions
	Authentication *apiserveroptions.AuthenticationWithAPIOptions
	Authorization  *AuthorizationOptions
	ETCD           *storageoptions.ETCDStorageOptions
	Auth           *AuthOptions
	Audit          *genericapiserveroptions.AuditOptions
	ConsoleConfig  *apiserver.ConsoleConfig
}

// NewOptions creates a new Options with a default config.
func NewOptions(serverName string) *Options {
	return &Options{
		Log:            log.NewOptions(),
		SecureServing:  apiserveroptions.NewSecureServingOptions(serverName, 9451),
		Debug:          apiserveroptions.NewDebugOptions(),
		Generic:        apiserveroptions.NewGenericOptions(),
		Authentication: apiserveroptions.NewAuthenticationWithAPIOptions(),
		Authorization:  NewAuthorizationOptions(),
		ETCD:           storageoptions.NewETCDStorageOptions("/tke/auth-api"),
		Auth:           NewAuthOptions(),
		Audit:          genericapiserveroptions.NewAuditOptions(),
		ConsoleConfig:  apiserver.NewConsoleConfigOptions(),
	}
}

// AddFlags adds flags for a specific server to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.Debug.AddFlags(fs)
	o.ETCD.AddFlags(fs)
	o.Generic.AddFlags(fs)
	o.Authentication.AddFlags(fs)
	o.Authorization.AddFlags(fs)
	o.Auth.AddFlags(fs)
	o.Audit.AddFlags(fs)
	o.ConsoleConfig.AddFlags(fs)
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.Log.ApplyFlags()...)
	errs = append(errs, o.SecureServing.ApplyFlags()...)
	errs = append(errs, o.Debug.ApplyFlags()...)
	errs = append(errs, o.ETCD.ApplyFlags()...)
	errs = append(errs, o.Generic.ApplyFlags()...)
	errs = append(errs, o.Authentication.ApplyFlags()...)
	errs = append(errs, o.Authorization.ApplyFlags()...)
	errs = append(errs, o.Auth.ApplyFlags()...)
	errs = append(errs, o.ConsoleConfig.ApplyFlags()...)

	return errs
}

// Complete set default Options.
// Should be called after tke-auth flags parsed.
func (o *Options) Complete() error {
	if err := apiserveroptions.CompleteGenericAndSecureOptions(o.Generic, o.SecureServing); err != nil {
		return err
	}

	if o.ETCD.EnableWatchCache {
		log.Infof("Initializing cache sizes based on %dMB limit", o.ETCD.DefaultWatchCacheSize)
		sizes := cachesize.NewHeuristicWatchCacheSizes(o.ETCD.DefaultWatchCacheSize)
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
	if (o.Audit.WebhookOptions.ConfigFile != "" || o.Audit.LogOptions.Path != "") && o.Audit.PolicyFile == "" {
		return fmt.Errorf("audit log/webhook config specified, but audit policy file is empty")
	}

	return nil
}
