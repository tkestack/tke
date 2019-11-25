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
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	controlleroptions "tkestack.io/tke/pkg/controller/options"
	"tkestack.io/tke/pkg/util/log"
)

// Options is the main context object for the TKE controller manager.
type Options struct {
	Log               *log.Options
	Debug             *apiserveroptions.DebugOptions
	SecureServing     *apiserveroptions.SecureServingOptions
	Component         *controlleroptions.ComponentOptions
	PlatformAPIClient *controlleroptions.APIServerClientOptions
	BusinessAPIClient *controlleroptions.APIServerClientOptions
}

// NewOptions creates a new Options with a default config.
func NewOptions(serverName string, allControllers []string, disabledByDefaultControllers []string) *Options {
	return &Options{
		Log:               log.NewOptions(),
		Debug:             apiserveroptions.NewDebugOptions(),
		SecureServing:     apiserveroptions.NewSecureServingOptions(serverName, 9449),
		Component:         controlleroptions.NewComponentOptions(allControllers, disabledByDefaultControllers),
		PlatformAPIClient: controlleroptions.NewAPIServerClientOptions("platform"),
		BusinessAPIClient: controlleroptions.NewAPIServerClientOptions("business"),
	}
}

// AddFlags adds flags for a specific server to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)
	o.Debug.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.Component.AddFlags(fs)
	o.PlatformAPIClient.AddFlags(fs)
	o.BusinessAPIClient.AddFlags(fs)
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.Log.ApplyFlags()...)
	errs = append(errs, o.Debug.ApplyFlags()...)
	errs = append(errs, o.SecureServing.ApplyFlags()...)
	errs = append(errs, o.Component.ApplyFlags()...)
	errs = append(errs, o.PlatformAPIClient.ApplyFlags()...)
	errs = append(errs, o.BusinessAPIClient.ApplyFlags()...)

	return errs
}
