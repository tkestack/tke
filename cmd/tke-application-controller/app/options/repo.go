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
	appconfig "tkestack.io/tke/pkg/application/controller/app/config"
)

const (
	flagRepoHost          = "repo-host"
	flagRepoCaFile        = "repo-cafile"
	flagRepoAdmin         = "repo-admin"
	flagRepoAdminPassword = "repo-admin-password"
)

const (
	configRepoHost          = "repo.host"
	configRepoCaFile        = "repo.cafile"
	configRepoAdmin         = "repo.admin"
	configRepoAdminPassword = "repo.admin_password"
)

// RepoOptions contains configuration items related to auth attributes.
type RepoOptions struct {
	Host          string
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
	fs.String(flagRepoHost, o.Host,
		"Chart repo host.")
	_ = viper.BindPFlag(configRepoHost, fs.Lookup(flagRepoHost))

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

	o.Host = viper.GetString(configRepoHost)
	if o.Host == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagRepoHost))
	}

	o.CaFile = viper.GetString(configRepoCaFile)
	o.Admin = viper.GetString(configRepoAdmin)
	o.AdminPassword = viper.GetString(configRepoAdminPassword)

	return errs
}

// ApplyTo fills up Debugging config with options.
func (o *RepoOptions) ApplyTo(cfg *appconfig.RepoConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.Host = o.Host
	cfg.CaFile = o.CaFile
	cfg.Admin = o.Admin
	cfg.AdminPassword = o.AdminPassword

	return nil
}
