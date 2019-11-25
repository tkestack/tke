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
	componentconfig "k8s.io/component-base/config"
)

const (
	flagDebugProfiling           = "profiling"
	flagDebugContentionProfiling = "contention-profiling"
)

const (
	configDebugProfiling           = "debug.profiling"
	configDebugContentionProfiling = "debug.contention_profiling"
)

// DebugOptions holds the Debugging options.
type DebugOptions struct {
	EnableProfiling           bool
	EnableContentionProfiling bool
}

// NewDebugOptions creates the default DebugOptions object.
func NewDebugOptions() *DebugOptions {
	return &DebugOptions{}
}

// AddFlags adds flags related to debugging for controller manager to the specified FlagSet.
func (o *DebugOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Bool(flagDebugProfiling, o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")
	_ = viper.BindPFlag(configDebugProfiling, fs.Lookup(flagDebugProfiling))
	fs.Bool(flagDebugContentionProfiling, o.EnableContentionProfiling,
		"Enable lock contention profiling, if profiling is enabled")
	_ = viper.BindPFlag(configDebugContentionProfiling, fs.Lookup(flagDebugContentionProfiling))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *DebugOptions) ApplyFlags() []error {
	var errs []error

	o.EnableProfiling = viper.GetBool(configDebugProfiling)
	o.EnableContentionProfiling = viper.GetBool(configDebugContentionProfiling)

	return errs
}

// ApplyTo fills up Debugging config with options.
func (o *DebugOptions) ApplyTo(cfg *componentconfig.DebuggingConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.EnableProfiling = o.EnableProfiling
	cfg.EnableContentionProfiling = o.EnableContentionProfiling

	return nil
}
