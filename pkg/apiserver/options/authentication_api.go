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
)

const (
	flagAuthnAPIAudiences       = "api-audiences"
	flagAuthnPrivilegedUsername = "privileged-username"
)

const (
	configAuthnAPIAudiences       = "authentication.api_audiences"
	configAuthnPrivilegedUsername = "authentication.privileged_username"
)

// AuthenticationWithAPIOptions contains the options that http request authentication.
type AuthenticationWithAPIOptions struct {
	APIAudiences       []string
	PrivilegedUsername string
	*AuthenticationOptions
}

// NewAuthenticationWithAPIOptions creates the default AuthenticationOptions object.
func NewAuthenticationWithAPIOptions() *AuthenticationWithAPIOptions {
	return &AuthenticationWithAPIOptions{
		AuthenticationOptions: NewAuthenticationOptions(),
		PrivilegedUsername:    "admin",
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *AuthenticationWithAPIOptions) AddFlags(fs *pflag.FlagSet) {
	o.AuthenticationOptions.AddFlags(fs)

	fs.StringSlice(flagAuthnAPIAudiences, o.APIAudiences, ""+
		"Identifiers of the API. The service account token authenticator will validate that "+
		"tokens used against the API are bound to at least one of these audiences. If the "+
		"--service-account-issuer flag is configured and this flag is not, this field "+
		"defaults to a single element list containing the issuer URL .")
	_ = viper.BindPFlag(configAuthnAPIAudiences, fs.Lookup(flagAuthnAPIAudiences))

	fs.String(flagAuthnPrivilegedUsername, o.PrivilegedUsername, "A privileged username with operations to delete the entire collection of resources")
	_ = viper.BindPFlag(configAuthnPrivilegedUsername, fs.Lookup(flagAuthnPrivilegedUsername))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthenticationWithAPIOptions) ApplyFlags() []error {
	var errs []error

	o.APIAudiences = viper.GetStringSlice(configAuthnAPIAudiences)
	o.PrivilegedUsername = viper.GetString(configAuthnPrivilegedUsername)

	errs = append(errs, o.AuthenticationOptions.ApplyFlags()...)

	return errs
}
