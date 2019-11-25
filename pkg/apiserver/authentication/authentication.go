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

package authentication

import (
	"fmt"
	"github.com/go-openapi/spec"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/pkg/apiserver/options"
)

// SetupAuthentication config the generic apiserver by authentication options.
func SetupAuthentication(genericAPIServerConfig *genericapiserver.Config, authenticationOpts *options.AuthenticationWithAPIOptions) error {
	genericAPIServerConfig.Authentication.APIAudiences = authenticationOpts.APIAudiences
	return SetupAuthenticationWithoutAudiences(genericAPIServerConfig, authenticationOpts.AuthenticationOptions, authenticationOpts.APIAudiences)
}

// SetupAuthenticationWithoutAudiences config the generic apiserver by
// authentication options.
func SetupAuthenticationWithoutAudiences(genericAPIServerConfig *genericapiserver.Config, authenticationOpts *options.AuthenticationOptions, apiAudiences []string) error {
	var err error
	genericAPIServerConfig.Authentication.Authenticator, genericAPIServerConfig.OpenAPIConfig.SecurityDefinitions, err = buildAuthenticator(authenticationOpts, apiAudiences)
	if err != nil {
		return fmt.Errorf("invalid authentication config: %v", err)
	}
	if authenticationOpts.ClientCert != nil {
		if err = genericAPIServerConfig.Authentication.ApplyClientCert(authenticationOpts.ClientCert.ClientCA, genericAPIServerConfig.SecureServing); err != nil {
			return fmt.Errorf("unable to load client CA file: %v", err)
		}
	}
	if authenticationOpts.RequestHeader != nil {
		if err = genericAPIServerConfig.Authentication.ApplyClientCert(authenticationOpts.RequestHeader.ClientCAFile, genericAPIServerConfig.SecureServing); err != nil {
			return fmt.Errorf("unable to load client CA file: %v", err)
		}
	}
	return nil
}

// buildAuthenticator constructs the authenticator.
func buildAuthenticator(o *options.AuthenticationOptions, apiAudiences []string) (authenticator.Request, *spec.SecurityDefinitions, error) {
	ret := Config{
		APIAudiences:         apiAudiences,
		TokenSuccessCacheTTL: o.TokenSuccessCacheTTL,
		TokenFailureCacheTTL: o.TokenFailureCacheTTL,
	}

	if o.ClientCert != nil {
		ret.ClientCAFile = o.ClientCert.ClientCA
	}

	if o.TokenFile != nil {
		ret.TokenAuthFile = o.TokenFile.TokenFile
	}

	if o.OIDC != nil {
		ret.OIDCCAFile = o.OIDC.CAFile
		ret.OIDCClientID = o.OIDC.ClientID
		ret.OIDCGroupsClaim = o.OIDC.GroupsClaim
		ret.OIDCGroupsPrefix = o.OIDC.GroupsPrefix
		ret.OIDCIssuerURL = o.OIDC.IssuerURL
		ret.OIDCExternalIssuerURL = o.OIDC.ExternalIssuerURL
		ret.OIDCUsernameClaim = o.OIDC.UsernameClaim
		ret.OIDCUsernamePrefix = o.OIDC.UsernamePrefix
		ret.OIDCTenantIDClaim = o.OIDC.TenantIDClaim
		ret.OIDCTenantIDPrefix = o.OIDC.TenantIDPrefix
		ret.OIDCSigningAlgs = o.OIDC.SigningAlgs
		ret.OIDCRequiredClaims = o.OIDC.RequiredClaims
	}

	if o.RequestHeader != nil {
		ret.RequestHeaderConfig = o.RequestHeader.ToAuthenticationRequestHeaderConfig()
	}

	return ret.New()
}
