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
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/provider/config"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	"tkestack.io/tke/pkg/util/log"
)

const (
	flagProviderCluster = "cluster-providers"
	flagProviderMachine = "machine-providers"
)

const (
	configProviderCluster = "providers.cluster"
	configProviderMachine = "providers.machine"
)

// Options contains configuration items related to TKE providers.
type Options struct {
	ClusterProviders []string
	MachineProviders []string
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
	fs.StringSlice(flagProviderCluster, o.ClusterProviders,
		"The plugin library file paths for cluster creation and deletion that provide different infrastructures, comma separated.")
	_ = viper.BindPFlag(configProviderCluster, fs.Lookup(flagProviderCluster))
	fs.StringSlice(flagProviderMachine, o.MachineProviders,
		"The plugin library file paths for machine creation and deletion that provide different infrastructures, comma separated.")
	_ = viper.BindPFlag(configProviderMachine, fs.Lookup(flagProviderMachine))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	o.ClusterProviders = viper.GetStringSlice(configProviderCluster)
	o.MachineProviders = viper.GetStringSlice(configProviderMachine)

	return errs
}

// ApplyTo convert feature options to provider config params.
func (o *Options) ApplyTo(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration instance that has not been initialized")
	}

	for _, p := range o.ClusterProviders {
		if err := applyClusterProviders(p, cfg); err != nil {
			return err
		}
	}

	for _, p := range o.MachineProviders {
		if err := applyMachineProviders(p, cfg); err != nil {
			return err
		}
	}

	return nil
}

func applyClusterProviders(p string, cfg *config.Config) error {
	pluginFile, configFile := pluginFileAndConfig(p)
	client, err := clusterprovider.NewClient(pluginFile, configFile)
	if err != nil {
		return err
	}
	log.Info("Successfully loaded cluster provider pluginFile", log.String("providerName", client.Name))

	cfg.ClusterProviders.Store(client.Name, client.Provider)
	cfg.Clients = append(cfg.Clients, client.Client)
	return nil
}

func applyMachineProviders(p string, cfg *config.Config) error {
	pluginFile, configFile := pluginFileAndConfig(p)
	client, err := machineprovider.NewClient(pluginFile, configFile)
	if err != nil {
		return err
	}
	log.Info("Successfully loaded machine provider plugin", log.String("providerName", client.Name))

	cfg.MachineProviders.Store(client.Name, client.Provider)
	cfg.Clients = append(cfg.Clients, client.Client)
	return nil
}

func pluginFileAndConfig(p string) (pluginFile string, configFile string) {
	a := strings.Split(p, "=")
	if len(a) == 1 {
		pluginFile = a[0]
	} else if len(a) == 2 {
		pluginFile = a[0]
		configFile = a[1]
	}
	return
}
