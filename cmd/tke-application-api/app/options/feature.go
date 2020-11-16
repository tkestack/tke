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
	appconfig "tkestack.io/tke/pkg/application/config"
)

const (
	flagRepoScheme        = "features-repo-scheme"
	flagRepoDomainSuffix  = "features-repo-domain-suffix"
	flagRepoCaFile        = "features-repo-cafile"
	flagRepoAdmin         = "features-repo-admin"
	flagRepoAdminPassword = "features-repo-admin-password"
)

const (
	configRepoScheme        = "features.repo.scheme"
	configRepoDomainSuffix  = "features.repo.domain_suffix"
	configRepoCaFile        = "features.repo.cafile"
	configRepoAdmin         = "features.repo.admin"
	configRepoAdminPassword = "features.repo.admin_password"
)

// RepoOptions contains configuration items related to application attributes.
type RepoOptions struct {
	Scheme        string
	DomainSuffix  string
	CaFile        string
	Admin         string
	AdminPassword string
}

// FeatureOptions contains configuration items related to application attributes.
type FeatureOptions struct {
	Repo RepoOptions
}

// NewFeatureOptions creates a FeatureOptions object with default parameters.
func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagRepoScheme, o.Repo.Scheme,
		"Chart repo server scheme.")
	_ = viper.BindPFlag(configRepoScheme, fs.Lookup(flagRepoScheme))

	fs.String(flagRepoDomainSuffix, o.Repo.DomainSuffix,
		"Chart repo domain suffix.")
	_ = viper.BindPFlag(configRepoDomainSuffix, fs.Lookup(flagRepoDomainSuffix))

	fs.String(flagRepoCaFile, o.Repo.CaFile,
		"CA certificate to verify peer against.")
	_ = viper.BindPFlag(configRepoCaFile, fs.Lookup(flagRepoCaFile))

	fs.String(flagRepoAdmin, o.Repo.Admin,
		"Repo admin user.")
	_ = viper.BindPFlag(configRepoAdmin, fs.Lookup(flagRepoAdmin))

	fs.String(flagRepoAdminPassword, o.Repo.AdminPassword,
		"Repo admin user password.")
	_ = viper.BindPFlag(configRepoAdminPassword, fs.Lookup(flagRepoAdminPassword))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.Repo.Scheme = viper.GetString(configRepoScheme)
	if o.Repo.Scheme == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagRepoScheme))
	}

	o.Repo.DomainSuffix = viper.GetString(configRepoDomainSuffix)
	if o.Repo.DomainSuffix == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagRepoDomainSuffix))
	}

	o.Repo.CaFile = viper.GetString(configRepoCaFile)
	o.Repo.Admin = viper.GetString(configRepoAdmin)
	o.Repo.AdminPassword = viper.GetString(configRepoAdminPassword)

	return errs
}

// ApplyTo fills up Debugging config with options.
func (o *RepoOptions) ApplyTo(cfg *appconfig.RepoConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.Scheme = o.Scheme
	cfg.DomainSuffix = o.DomainSuffix
	cfg.CaFile = o.CaFile
	cfg.Admin = o.Admin
	cfg.AdminPassword = o.AdminPassword

	return nil
}
