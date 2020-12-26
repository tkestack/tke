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
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryscheme "tkestack.io/tke/pkg/registry/apis/config/scheme"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	registrycontrollerconfig "tkestack.io/tke/pkg/registry/controller/config"
)

// NewRegistryConfiguration will create a new RegistryConfiguration with default values
func NewRegistryConfiguration() (*registryconfig.RegistryConfiguration, error) {
	scheme, _, err := registryscheme.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}
	versioned := &registryconfigv1.RegistryConfiguration{}
	scheme.Default(versioned)
	config := &registryconfig.RegistryConfiguration{}
	if err := scheme.Convert(versioned, config, nil); err != nil {
		return nil, err
	}
	return config, nil
}

const (
	flagDefaultSystemChartGroups = "registry-setting-default-system-chartgroups"
	flagChartPaths               = "registry-setting-chart-paths"
)

const (
	configDefaultSystemChartGroups = "registry_setting.default_system_chartgroups"
	configChartPaths               = "registry_setting.chart_paths"
)

// ChartGroupSettingOptions contains configuration items related to registry attributes.
type ChartGroupSettingOptions struct {
	DefaultSystemChartGroups []string
	ChartPaths               []string
}

// NewChartGroupSettingOptions creates a ChartGroupSettingOptions object with default parameters.
func NewChartGroupSettingOptions() *ChartGroupSettingOptions {
	return &ChartGroupSettingOptions{}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *ChartGroupSettingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSlice(flagDefaultSystemChartGroups, o.DefaultSystemChartGroups,
		"Default chartgroups with system type and public visibility.")
	_ = viper.BindPFlag(configDefaultSystemChartGroups, fs.Lookup(flagDefaultSystemChartGroups))

	fs.StringSlice(flagChartPaths, o.ChartPaths,
		"Path to the default charts which will be load to the default platform chartgroup when started.")
	_ = viper.BindPFlag(configChartPaths, fs.Lookup(flagChartPaths))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ChartGroupSettingOptions) ApplyFlags() []error {
	var errs []error

	o.DefaultSystemChartGroups = viper.GetStringSlice(configDefaultSystemChartGroups)
	o.ChartPaths = viper.GetStringSlice(configChartPaths)

	return errs
}

// ApplyTo fills up Debugging config with options.
func (o *ChartGroupSettingOptions) ApplyTo(cfg *registrycontrollerconfig.ChartGroupSetting) error {
	if o == nil {
		return nil
	}

	cfg.DefaultSystemChartGroups = o.DefaultSystemChartGroups[:]
	cfg.ChartPaths = o.ChartPaths[:]

	return nil
}
