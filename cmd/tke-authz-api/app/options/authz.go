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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagDefaultPoliciesConfig = "default-policies-config"
	flagDefaultRolesConfig    = "default-roles-config"
)

const (
	configDefaultPoliciesConfig = "authz.default_policies_config"
	configDefaultRolesConfig    = "authz.default_roles_config"
)

type AuthzOptions struct {
	DefaultPoliciesConfig string
	DefaultRolesConfig    string
}

func NewAuthzOptions() *AuthzOptions {
	return &AuthzOptions{}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *AuthzOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.DefaultPoliciesConfig, flagDefaultPoliciesConfig, o.DefaultPoliciesConfig, "Default policies config")
	_ = viper.BindPFlag(configDefaultPoliciesConfig, fs.Lookup(flagDefaultPoliciesConfig))

	fs.StringVar(&o.DefaultRolesConfig, flagDefaultRolesConfig, o.DefaultRolesConfig, "Default roles config")
	_ = viper.BindPFlag(configDefaultRolesConfig, fs.Lookup(flagDefaultRolesConfig))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthzOptions) ApplyFlags() []error {
	var errs []error
	o.DefaultPoliciesConfig = viper.GetString(configDefaultPoliciesConfig)
	o.DefaultRolesConfig = viper.GetString(configDefaultRolesConfig)
	return errs
}
