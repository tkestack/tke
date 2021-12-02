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
 */

package options

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagNotFoundCRDProxy = "crd-proxy"
)

const (
	configNotFoundCRDProxy = "features.crd_proxy"
)

const (
	defaultNotFoundCRDProxy = false
)

// FeatureOptions contains the options that feature gate.
type FeatureOptions struct {
	NotFoundCRDProxy bool
}

// NewFeatureOptions creates an Options object with default feature parameters.
func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		NotFoundCRDProxy: defaultNotFoundCRDProxy,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Bool(flagNotFoundCRDProxy, o.NotFoundCRDProxy,
		"Enable url proxy that can transfer custom resource.")
	_ = viper.BindPFlag(configNotFoundCRDProxy, fs.Lookup(flagNotFoundCRDProxy))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.NotFoundCRDProxy = viper.GetBool(configNotFoundCRDProxy)

	return errs
}
