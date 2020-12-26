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

package features

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"tkestack.io/tke/pkg/registry/config"
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

// NewRepoOptions creates a RepoOptions object with default parameters.
func NewRepoOptions() *RepoOptions {
	return &RepoOptions{}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *RepoOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagRepoScheme, o.Scheme,
		"Chart repo server scheme.")
	_ = viper.BindPFlag(configRepoScheme, fs.Lookup(flagRepoScheme))

	fs.String(flagRepoDomainSuffix, o.DomainSuffix,
		"Chart repo domain suffix.")
	_ = viper.BindPFlag(configRepoDomainSuffix, fs.Lookup(flagRepoDomainSuffix))

	fs.String(flagRepoCaFile, o.CaFile,
		"CA certificate to verify peer against.")
	_ = viper.BindPFlag(configRepoCaFile, fs.Lookup(flagRepoCaFile))

	fs.String(flagRepoAdmin, o.Admin,
		"Repo admin user.")
	_ = viper.BindPFlag(configRepoAdmin, fs.Lookup(flagRepoAdmin))

	fs.String(flagRepoAdminPassword, o.AdminPassword,
		"Repo admin user password.")
	_ = viper.BindPFlag(configRepoAdminPassword, fs.Lookup(flagRepoAdminPassword))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *RepoOptions) ApplyFlags() []error {
	var errs []error

	o.Scheme = viper.GetString(configRepoScheme)
	if o.Scheme == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagRepoScheme))
	}

	o.DomainSuffix = viper.GetString(configRepoDomainSuffix)
	if o.DomainSuffix == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagRepoDomainSuffix))
	}

	o.CaFile = viper.GetString(configRepoCaFile)
	o.Admin = viper.GetString(configRepoAdmin)
	o.AdminPassword = viper.GetString(configRepoAdminPassword)

	return errs
}

// ApplyTo fills up Debugging config with options.
func (o *RepoOptions) ApplyTo(cfg *config.RepoConfiguration) error {
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
