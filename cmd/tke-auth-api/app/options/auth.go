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

	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/ldap"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagAuthAssetsPath             = "assets-path"
	flagAuthIDTokenTimeout         = "id-token-timeout"
	flagAuthInitTenantType         = "init-tenant-type"
	flagAuthInitTenantID           = "init-tenant-id"
	flagAuthLDAPConfigFile         = "auth.ldap-config-file"
	flagAuthTenantAdmin            = "tenant-admin"
	flagAuthTenantAdminSecret      = "tenant-admin-secret"
	flagAuthInitClientID           = "init-client-id"
	flagAuthInitClientSecret       = "init-client-secret"
	flagAuthInitClientRedirectUris = "init-client-redirect-uris"
)

const (
	configAuthAssetsPath             = "auth.assets_path"
	configAuthIDTokenTimeout         = "auth.id_token_timeout"
	configAuthInitTenantType         = "auth.init_tenant_type"
	configAuthInitTenantID           = "auth.init_tenant_id"
	configAuthLDAPConfigFile         = "auth.ldap_config_file"
	configAuthTenantAdmin            = "auth.tenant_admin"
	configAuthTenantAdminSecret      = "auth.tenant_admin_secret"
	configAuthInitClientID           = "auth.init_client_id"
	configAuthInitClientSecret       = "auth.init_client_secret"
	configAuthInitClientRedirectUris = "auth.init_client_redirect_uris"
)

// AuthOptions contains configuration items related to auth attributes.
type AuthOptions struct {
	AssetsPath             string
	IDTokenTimeout         time.Duration
	InitTenantType         string
	InitTenantID           string
	LdapConfigFile         string
	TenantAdmin            string
	TenantAdminSecret      string
	InitClientID           string
	InitClientSecret       string
	InitClientRedirectUris []string
}

// NewAuthOptions creates a AuthOptions object with default parameters.
func NewAuthOptions() *AuthOptions {
	return &AuthOptions{
		IDTokenTimeout: 24 * time.Hour,
		InitTenantType: local.ConnectorType,
		InitTenantID:   "default",
		InitClientID:   "default",
	}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *AuthOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagAuthAssetsPath, o.AssetsPath,
		"Path to the OIDC front-end file assets.")
	_ = viper.BindPFlag(configAuthAssetsPath, fs.Lookup(flagAuthAssetsPath))

	fs.Duration(flagAuthIDTokenTimeout, o.IDTokenTimeout,
		"An optional field indicating the valid duration of the IDToken the OIDC generated. If blank, default value is 24h.")
	_ = viper.BindPFlag(configAuthIDTokenTimeout, fs.Lookup(flagAuthIDTokenTimeout))

	fs.String(flagAuthInitTenantType, o.InitTenantType,
		"Default tenant type for auth first started, supported tke,ldap.(default tke).")
	_ = viper.BindPFlag(configAuthInitTenantType, fs.Lookup(flagAuthInitTenantType))

	fs.String(flagAuthInitTenantID, o.InitTenantID,
		"Default tenant id will be created when started.")
	_ = viper.BindPFlag(configAuthInitTenantID, fs.Lookup(flagAuthInitTenantID))

	fs.String(flagAuthLDAPConfigFile, o.LdapConfigFile,
		"Config file path for ldap ldap, must specify if init-tenant-type is ldap.")
	_ = viper.BindPFlag(configAuthLDAPConfigFile, fs.Lookup(flagAuthLDAPConfigFile))

	fs.String(flagAuthTenantAdmin, o.TenantAdmin,
		"Default tenant admin name will be created when started.")
	_ = viper.BindPFlag(configAuthTenantAdmin, fs.Lookup(flagAuthTenantAdmin))

	fs.String(flagAuthTenantAdminSecret, o.TenantAdminSecret,
		"Secret for generate tenant admin login password.")
	_ = viper.BindPFlag(configAuthTenantAdminSecret, fs.Lookup(flagAuthTenantAdminSecret))

	fs.String(flagAuthInitClientID, o.InitClientID,
		"Default client id will be created when started.")
	_ = viper.BindPFlag(configAuthInitClientID, fs.Lookup(flagAuthInitClientID))

	fs.String(flagAuthInitClientSecret, o.InitClientSecret,
		"Default client secret will be created when started.")
	_ = viper.BindPFlag(configAuthInitClientSecret, fs.Lookup(flagAuthInitClientSecret))

	fs.StringSlice(flagAuthInitClientRedirectUris, o.InitClientRedirectUris,
		"Default client redirect uris will be created when started.")
	_ = viper.BindPFlag(configAuthInitClientRedirectUris, fs.Lookup(flagAuthInitClientRedirectUris))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthOptions) ApplyFlags() []error {
	var errs []error

	o.AssetsPath = viper.GetString(configAuthAssetsPath)

	if o.AssetsPath == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagAuthAssetsPath))
	}

	o.IDTokenTimeout = viper.GetDuration(configAuthIDTokenTimeout)

	o.InitTenantType = viper.GetString(configAuthInitTenantType)
	o.LdapConfigFile = viper.GetString(configAuthLDAPConfigFile)
	if o.InitTenantType == ldap.ConnectorType && o.LdapConfigFile == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified for ldap type tenant", flagAuthLDAPConfigFile))
	}

	o.InitTenantID = viper.GetString(configAuthInitTenantID)
	if len(o.InitTenantID) == 0 {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagAuthInitTenantID))
	}
	o.TenantAdmin = viper.GetString(configAuthTenantAdmin)
	o.TenantAdminSecret = viper.GetString(configAuthTenantAdminSecret)

	if len(o.TenantAdmin) == 0 || len(o.TenantAdminSecret) == 0 {
		errs = append(errs, fmt.Errorf("--%s and --%s must be specified", flagAuthTenantAdmin, flagAuthTenantAdminSecret))
	}

	o.InitClientID = viper.GetString(configAuthInitClientID)
	o.InitClientSecret = viper.GetString(configAuthInitClientSecret)
	if len(o.InitClientID) == 0 || len(o.InitClientSecret) == 0 {
		errs = append(errs, fmt.Errorf("--%s and --%s must be specified", flagAuthInitClientID, flagAuthInitClientSecret))
	}
	o.InitClientRedirectUris = viper.GetStringSlice(configAuthInitClientRedirectUris)

	return errs
}
