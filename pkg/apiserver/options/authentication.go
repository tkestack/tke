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
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

const (
	flagAuthnClientCAFile                     = "client-ca-file"
	flagAuthnTokenFile                        = "token-auth-file"
	flagAuthnRequestHeaderUsernameHeaders     = "requestheader-username-headers"
	flagAuthnRequestHeaderGroupHeaders        = "requestheader-group-headers"
	flagAuthnRequestHeaderExtraHeaderPrefixes = "requestheader-extra-headers-prefix"
	flagAuthnRequestHeaderClientCAFile        = "requestheader-client-ca-file"
	flagAuthnRequestHeaderAllowedNames        = "requestheader-allowed-names"
)

const (
	configAuthnClientCAFile                     = "authentication.client_ca_file"
	configAuthnTokenFile                        = "authentication.token_auth_file"
	configAuthnRequestHeaderUsernameHeaders     = "authentication.requestheader.username_headers"
	configAuthnRequestHeaderGroupHeaders        = "authentication.requestheader.group_headers"
	configAuthnRequestHeaderExtraHeaderPrefixes = "authentication.requestheader.extra_headers_prefix"
	configAuthnRequestHeaderClientCAFile        = "authentication.requestheader.client_ca_file"
	configAuthnRequestHeaderAllowedNames        = "authentication.requestheader.allowed_names"
)

// AuthenticationOptions contains the options that http request authentication.
type AuthenticationOptions struct {
	ClientCert           *genericoptions.ClientCertAuthenticationOptions
	OIDC                 *OIDCOptions
	RequestHeader        *genericoptions.RequestHeaderAuthenticationOptions
	TokenFile            *TokenFileAuthenticationOptions
	TokenSuccessCacheTTL time.Duration
	TokenFailureCacheTTL time.Duration
}

// PasswordFileAuthenticationOptions defines the configuration when using static
// password file authentication
type PasswordFileAuthenticationOptions struct {
	BasicAuthFile string
}

// TokenFileAuthenticationOptions defines the configuration when using static
// token file authentication.
type TokenFileAuthenticationOptions struct {
	TokenFile string
}

// NewAuthenticationOptions creates the default AuthenticationOptions object.
func NewAuthenticationOptions() *AuthenticationOptions {
	return &AuthenticationOptions{
		ClientCert:           &genericoptions.ClientCertAuthenticationOptions{},
		OIDC:                 NewOIDCOptions(),
		RequestHeader:        &genericoptions.RequestHeaderAuthenticationOptions{},
		TokenFile:            &TokenFileAuthenticationOptions{},
		TokenSuccessCacheTTL: 10 * time.Second,
		TokenFailureCacheTTL: 0 * time.Second,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *AuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	o.ClientCert.AddFlags(fs)
	_ = viper.BindPFlag(configAuthnClientCAFile, fs.Lookup(flagAuthnClientCAFile))

	fs.String("token-auth-file", o.TokenFile.TokenFile, ""+
		"If set, the file that will be used to secure the secure port of the API server "+
		"via token authentication.")
	_ = viper.BindPFlag(configAuthnTokenFile, fs.Lookup(flagAuthnTokenFile))

	o.OIDC.AddFlags(fs)

	o.RequestHeader.AddFlags(fs)
	_ = viper.BindPFlag(configAuthnRequestHeaderUsernameHeaders, fs.Lookup(flagAuthnRequestHeaderUsernameHeaders))
	_ = viper.BindPFlag(configAuthnRequestHeaderGroupHeaders, fs.Lookup(flagAuthnRequestHeaderGroupHeaders))
	_ = viper.BindPFlag(configAuthnRequestHeaderExtraHeaderPrefixes, fs.Lookup(flagAuthnRequestHeaderExtraHeaderPrefixes))
	_ = viper.BindPFlag(configAuthnRequestHeaderClientCAFile, fs.Lookup(flagAuthnRequestHeaderClientCAFile))
	_ = viper.BindPFlag(configAuthnRequestHeaderAllowedNames, fs.Lookup(flagAuthnRequestHeaderAllowedNames))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthenticationOptions) ApplyFlags() []error {
	var errs []error

	o.ClientCert.ClientCA = viper.GetString(configAuthnClientCAFile)
	o.TokenFile.TokenFile = viper.GetString(configAuthnTokenFile)

	errs = append(errs, o.OIDC.ApplyFlags()...)

	o.RequestHeader.AllowedNames = viper.GetStringSlice(configAuthnRequestHeaderAllowedNames)
	o.RequestHeader.ClientCAFile = viper.GetString(configAuthnRequestHeaderClientCAFile)
	o.RequestHeader.ExtraHeaderPrefixes = viper.GetStringSlice(configAuthnRequestHeaderExtraHeaderPrefixes)
	o.RequestHeader.GroupHeaders = viper.GetStringSlice(configAuthnRequestHeaderGroupHeaders)
	o.RequestHeader.UsernameHeaders = viper.GetStringSlice(configAuthnRequestHeaderUsernameHeaders)

	if o.OIDC != nil && (len(o.OIDC.IssuerURL) > 0) != (len(o.OIDC.ClientID) > 0) {
		errs = append(errs, fmt.Errorf("oidc-issuer-url and oidc-client-id should be specified together"))
	}

	if o.OIDC != nil && len(o.OIDC.ExternalIssuerURL) == 0 {
		o.OIDC.ExternalIssuerURL = o.OIDC.IssuerURL
	}

	return errs
}
