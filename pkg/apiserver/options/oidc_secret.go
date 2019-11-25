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

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagOIDCClientSecret = "oidc-client-secret"
)

const (
	configOIDCClientSecret = "authentication.oidc.client_secret"
)

// OIDCWithSecretOptions defines the configuration options needed to initialize
// OpenID Connect authentication.
type OIDCWithSecretOptions struct {
	*OIDCOptions
	ClientSecret string
}

// NewOIDCWithSecretOptions creates the default OIDCWithSecretOptions object.
func NewOIDCWithSecretOptions() *OIDCWithSecretOptions {
	return &OIDCWithSecretOptions{
		OIDCOptions: NewOIDCOptions(),
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *OIDCWithSecretOptions) AddFlags(fs *pflag.FlagSet) {
	o.OIDCOptions.AddFlags(fs)
	fs.String(flagOIDCClientSecret, o.ClientSecret,
		"The client secret for the OpenID connect client, must be set if oidc-issuer-url is set.")
	_ = viper.BindPFlag(configOIDCClientSecret, fs.Lookup(flagOIDCClientSecret))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *OIDCWithSecretOptions) ApplyFlags() []error {
	var errs []error

	o.ClientSecret = viper.GetString(configOIDCClientSecret)

	errs = append(errs, o.OIDCOptions.ApplyFlags()...)

	if o.ClientSecret == "" {
		errs = append(errs, fmt.Errorf("oidc-client-secret must be specified"))
	}

	if o.OIDCOptions == nil {
		panic("OIDCOptions is not initialized")
	}

	if o.ClientID == "" {
		errs = append(errs, fmt.Errorf("oidc-client-id must be specified"))
	}

	if o.IssuerURL == "" {
		errs = append(errs, fmt.Errorf("oidc-issuer-url must be specified"))
	} else {
		if o.ExternalIssuerURL == "" {
			o.ExternalIssuerURL = o.IssuerURL
		}
	}

	return errs
}
