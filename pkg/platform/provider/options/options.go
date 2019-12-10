/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package options

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalmachine "tkestack.io/tke/pkg/platform/provider/baremetal/machine"
	"tkestack.io/tke/pkg/platform/provider/config"
)

const (
	flagProvider = "providers"
)

const (
	configProvider = "provider.providers"
)

// Options contains configuration items related to TKE providers.
type Options struct {
	Providers []string
}

// NewOptions creates a Options object with default parameters.
func NewOptions() *Options {
	return &Options{}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *Options) Validate() []error {
	var errs []error
	return errs
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringSlice(flagProvider, o.Providers,
		"cluster and machine provider, comma separated. available values: Baremetal")
	_ = viper.BindPFlag(configProvider, fs.Lookup(flagProvider))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	o.Providers = viper.GetStringSlice(configProvider)

	return errs
}

// ApplyTo convert feature options to provider config params.
func (o *Options) ApplyTo(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration instance that has not been initialized")
	}

	if funk.ContainsString(o.Providers, "Baremetal") {
		clusterProvider, err := baremetalcluster.NewProvider()
		if err != nil {
			return err
		}
		cfg.ClusterProviders.Store(clusterProvider.Name(), clusterProvider)

		machineProvider, err := baremetalmachine.NewProvider()
		if err != nil {
			return err
		}
		cfg.MachineProviders.Store(machineProvider.Name(), machineProvider)
	}

	return nil
}
