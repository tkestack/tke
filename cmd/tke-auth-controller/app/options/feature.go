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
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagPolicyPath           = "policy-path"
	flagCategoryPath         = "category-path"
	flagTenantAdmin          = "tenant-admin"
	flagTenantAdminSecret    = "tenant-admin-secret"
	flagCasbinModelFile      = "casbin-model-file"
	flagCasbinReLoadInterval = "casbin-reload-interval"
)

const (
	configPolicyPath           = "features.policy_path"
	configCategoryPath         = "features.category_path"
	configTenantAdmin          = "features.tenant_admin"
	configTenantAdminSecret    = "features.tenant_admin_secret"
	configCasbinModelFile      = "features.casbin_model_file"
	configCasbinReloadInterval = "features.casbin_reload_interval"
)

type FeatureOptions struct {
	PolicyPath           string
	CategoryPath         string
	TenantAdmin          string
	TenantAdminSecret    string
	CasbinModelFile      string
	CasbinReloadInterval time.Duration
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{CasbinReloadInterval: 5*time.Second}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagPolicyPath, o.PolicyPath,
		"Path to the default policies which will be load to storage when started.")
	_ = viper.BindPFlag(configPolicyPath, fs.Lookup(flagPolicyPath))

	fs.String(flagCategoryPath, o.CategoryPath,
		"Path to the category which will be load to storage when started.")
	_ = viper.BindPFlag(configCategoryPath, fs.Lookup(flagCategoryPath))

	fs.String(flagTenantAdmin, o.TenantAdmin,
		"Default tenant admin name will be created when started.")
	_ = viper.BindPFlag(configTenantAdmin, fs.Lookup(flagTenantAdmin))

	fs.String(flagTenantAdminSecret, o.TenantAdminSecret,
		"Secret for generate tenant admin login password.")
	_ = viper.BindPFlag(configTenantAdminSecret, fs.Lookup(flagTenantAdminSecret))

	fs.String(flagCasbinModelFile, o.CasbinModelFile,
		"Casbin model file used to store ACL model.")
	_ = viper.BindPFlag(configCasbinModelFile, fs.Lookup(flagCasbinModelFile))

	fs.Duration(flagCasbinReLoadInterval, o.CasbinReloadInterval,
		"The interval of casbin reload policy from backend storage. Default 5s.")
	_ = viper.BindPFlag(configCasbinReloadInterval, fs.Lookup(flagCasbinReLoadInterval))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.PolicyPath = viper.GetString(configPolicyPath)
	o.CategoryPath = viper.GetString(configCategoryPath)

	o.TenantAdmin = viper.GetString(configTenantAdmin)
	o.TenantAdminSecret = viper.GetString(configTenantAdminSecret)

	o.CasbinModelFile = viper.GetString(configCasbinModelFile)
	o.CasbinReloadInterval = viper.GetDuration(configCasbinReloadInterval)

	if len(o.TenantAdmin) == 0 || len(o.TenantAdminSecret) == 0 {
		errs = append(errs, fmt.Errorf("%s and %s must be specified", configTenantAdmin, configTenantAdminSecret))
	}

	return errs
}
