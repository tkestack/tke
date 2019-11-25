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
	flagOIDCIssuerURL         = "oidc-issuer-url"
	flagOIDCClientID          = "oidc-client-id"
	flagOIDCCAFile            = "oidc-ca-file"
	flagOIDCExternalIssuerURL = "oidc-external-issuer-url"
	flagOIDCUsernameClaim     = "oidc-username-claim"
	flagOIDCUsernamePrefix    = "oidc-username-prefix"
	flagOIDCGroupsPrefix      = "oidc-groups-prefix"
	flagOIDCGroupsClaim       = "oidc-groups-claim"
	flagOIDCTenantIDClaim     = "oidc-tenantid-claim"
	flagOIDCTenantIDPrefix    = "oidc-tenantid-prefix"
	flagOIDCSigningAlgs       = "oidc-signing-algs"
	flagOIDCRequiredClaims    = "oidc-required-claim"
	flagOIDCTokenReviewPath   = "oidc-token-review-path"
)

const (
	configOIDCIssuerURL         = "authentication.oidc.issuer_url"
	configOIDCClientID          = "authentication.oidc.client_id"
	configOIDCCAFile            = "authentication.oidc.ca_file"
	configOIDCExternalIssuerURL = "authentication.oidc.external_issuer_url"
	configOIDCUsernameClaim     = "authentication.oidc.username_claim"
	configOIDCUsernamePrefix    = "authentication.oidc.username_prefix"
	configOIDCGroupsPrefix      = "authentication.oidc.groups_prefix"
	configOIDCGroupsClaim       = "authentication.oidc.groups_claim"
	configOIDCTenantIDClaim     = "authentication.oidc.tenantid_claim"
	configOIDCTenantIDPrefix    = "authentication.oidc.tenantid_prefix"
	configOIDCSigningAlgs       = "authentication.oidc.signing_algs"
	configOIDCRequiredClaims    = "authentication.oidc.required_claim"
	configOIDCTokenReviewPath   = "authentication.oidc.token_review_path"
)

// OIDCOptions defines the configuration options needed to initialize OpenID
// Connect authentication.
type OIDCOptions struct {
	CAFile            string
	ClientID          string
	IssuerURL         string
	ExternalIssuerURL string
	UsernameClaim     string
	UsernamePrefix    string
	GroupsClaim       string
	GroupsPrefix      string
	TenantIDClaim     string
	TenantIDPrefix    string
	SigningAlgs       []string
	RequiredClaims    map[string]string
	TokenReviewPath   string
}

// NewOIDCOptions creates the default OIDCOptions object.
func NewOIDCOptions() *OIDCOptions {
	return &OIDCOptions{
		UsernameClaim: "sub",
		SigningAlgs:   []string{"RS256"},
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *OIDCOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagOIDCIssuerURL, o.IssuerURL, ""+
		"The URL of the OpenID internal issuer, only HTTPS scheme will be accepted. "+
		"If set, it will be used to verify the OIDC JSON Web Token (JWT).")
	_ = viper.BindPFlag(configOIDCIssuerURL, fs.Lookup(flagOIDCIssuerURL))
	fs.String(flagOIDCClientID, o.ClientID,
		"The client ID for the OpenID Connect client, must be set if oidc-issuer-url is set.")
	_ = viper.BindPFlag(configOIDCClientID, fs.Lookup(flagOIDCClientID))
	fs.String(flagOIDCCAFile, o.CAFile, ""+
		"If set, the OpenID server's certificate will be verified by one of the authorities "+
		"in the oidc-ca-file, otherwise the host's root CA set will be used.")
	_ = viper.BindPFlag(configOIDCCAFile, fs.Lookup(flagOIDCCAFile))
	fs.String(flagOIDCUsernameClaim, o.UsernameClaim, ""+
		"The OpenID claim to use as the user name. Note that claims other than the default ('sub') "+
		"is not guaranteed to be unique and immutable.")
	fs.String(flagOIDCExternalIssuerURL, o.ExternalIssuerURL, ""+
		"The URL of the OpenID external issuer, only HTTPS scheme will be accepted.  "+
		"is set, it will be used to verify issuer url with the OIDC provider, otherwise the oidc-issuer-url will be used.")
	_ = viper.BindPFlag(configOIDCExternalIssuerURL, fs.Lookup(flagOIDCExternalIssuerURL))
	_ = viper.BindPFlag(configOIDCUsernameClaim, fs.Lookup(flagOIDCUsernameClaim))
	fs.String(flagOIDCUsernamePrefix, o.UsernamePrefix, ""+
		"If provided, all usernames will be prefixed with this value. If not provided, "+
		"username claims other than 'email' are prefixed by the issuer URL to avoid "+
		"clashes. To skip any prefixing, provide the value '-'.")
	_ = viper.BindPFlag(configOIDCUsernamePrefix, fs.Lookup(flagOIDCUsernamePrefix))
	fs.String(flagOIDCTenantIDClaim, o.TenantIDClaim,
		"If provided, the name of a custom OpenID Connect claim for specifying user tenant id.")
	_ = viper.BindPFlag(configOIDCTenantIDClaim, fs.Lookup(flagOIDCTenantIDClaim))
	fs.String(flagOIDCTenantIDPrefix, o.TenantIDPrefix,
		"If provided, all tenant ids will be prefixed with this value to prevent conflicts with "+
			"other authentication strategies.")
	_ = viper.BindPFlag(configOIDCTenantIDPrefix, fs.Lookup(flagOIDCTenantIDPrefix))
	fs.String(flagOIDCGroupsClaim, o.GroupsClaim, ""+
		"If provided, the name of a custom OpenID Connect claim for specifying user groups. "+
		"The claim value is expected to be a string or array of strings. ")
	_ = viper.BindPFlag(configOIDCGroupsClaim, fs.Lookup(flagOIDCGroupsClaim))
	fs.String(flagOIDCGroupsPrefix, o.GroupsPrefix, ""+
		"If provided, all groups will be prefixed with this value to prevent conflicts with "+
		"other authentication strategies.")
	_ = viper.BindPFlag(configOIDCGroupsPrefix, fs.Lookup(flagOIDCGroupsPrefix))
	fs.StringSlice(flagOIDCSigningAlgs, o.SigningAlgs, ""+
		"Comma-separated list of allowed JOSE asymmetric signing algorithms. JWTs with a "+
		"'alg' header value not in this list will be rejected. "+
		"Values are defined by RFC 7518 https://tools.ietf.org/html/rfc7518#section-3.1.")
	_ = viper.BindPFlag(configOIDCSigningAlgs, fs.Lookup(flagOIDCSigningAlgs))
	fs.String(flagOIDCRequiredClaims, "", ""+
		"A key=value pair that describes a required claim in the ID Token. "+
		"If set, the claim is verified to be present in the ID Token with a matching value. "+
		"Repeat this flag to specify multiple claims.")
	_ = viper.BindPFlag(configOIDCRequiredClaims, fs.Lookup(flagOIDCRequiredClaims))
	fs.String(flagOIDCTokenReviewPath, "", ""+
		"Set the access path used by the OIDC server to verify the validity of the ID Token "+
		"generated by it.")
	_ = viper.BindPFlag(configOIDCTokenReviewPath, fs.Lookup(flagOIDCTokenReviewPath))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *OIDCOptions) ApplyFlags() []error {
	var errs []error
	o.CAFile = viper.GetString(configOIDCCAFile)
	o.ClientID = viper.GetString(configOIDCClientID)
	o.GroupsClaim = viper.GetString(configOIDCGroupsClaim)
	o.GroupsPrefix = viper.GetString(configOIDCGroupsPrefix)
	o.TenantIDPrefix = viper.GetString(configOIDCTenantIDPrefix)
	o.TenantIDClaim = viper.GetString(configOIDCTenantIDClaim)
	o.IssuerURL = viper.GetString(configOIDCIssuerURL)
	o.ExternalIssuerURL = viper.GetString(configOIDCExternalIssuerURL)
	o.RequiredClaims = viper.GetStringMapString(configOIDCRequiredClaims)
	o.SigningAlgs = viper.GetStringSlice(configOIDCSigningAlgs)
	o.UsernameClaim = viper.GetString(configOIDCUsernameClaim)
	o.UsernamePrefix = viper.GetString(configOIDCUsernamePrefix)
	o.TokenReviewPath = viper.GetString(configOIDCTokenReviewPath)

	return errs
}
