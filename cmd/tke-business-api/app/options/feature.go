/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package options

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagTenantOfInitialAdministrator = "initial-administrator-tenant"
	flagUserOfInitialAdministrator   = "initial-administrator-user"
	flagSyncProjectsWithNamespaces   = "sync-projects-with-namespaces"
)

const (
	configTenantOfInitialAdministrator = "features.initial_administrator_tenant"
	configUserOfInitialAdministrator   = "features.initial_administrator_user"
	configSyncProjectsWithNamespaces   = "features.sync_projects_with_namespaces"
)

const (
	defaultTenantOfInitialAdministrator = "default"
	defaultUserOfInitialAdministrator   = "admin"
)

const DefaultPlatform = "platform-default"

type FeatureOptions struct {
	TenantOfInitialAdministrator string
	UserOfInitialAdministrator   string
	SyncProjectsWithNamespaces   bool
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		TenantOfInitialAdministrator: defaultTenantOfInitialAdministrator,
		UserOfInitialAdministrator:   defaultUserOfInitialAdministrator,
	}
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagTenantOfInitialAdministrator, o.TenantOfInitialAdministrator,
		"The tenant name of initial administrator.")
	_ = viper.BindPFlag(configTenantOfInitialAdministrator, fs.Lookup(flagTenantOfInitialAdministrator))
	fs.String(flagUserOfInitialAdministrator, o.UserOfInitialAdministrator,
		"The user name of initial administrator.")
	_ = viper.BindPFlag(configUserOfInitialAdministrator, fs.Lookup(flagUserOfInitialAdministrator))
	fs.Bool(flagSyncProjectsWithNamespaces, o.SyncProjectsWithNamespaces,
		"Enable creating/deleting the corresponding namespace when creating/deleting a project.")
	_ = viper.BindPFlag(configSyncProjectsWithNamespaces, fs.Lookup(flagSyncProjectsWithNamespaces))
}

func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.TenantOfInitialAdministrator = viper.GetString(configTenantOfInitialAdministrator)
	o.UserOfInitialAdministrator = viper.GetString(configUserOfInitialAdministrator)
	o.SyncProjectsWithNamespaces = viper.GetBool(configSyncProjectsWithNamespaces)

	return errs
}
