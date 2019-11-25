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

package config

import (
	"context"
	"fmt"
	gooidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"k8s.io/apiserver/pkg/authentication/request/anonymous"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog"
	"net/http"
	"path/filepath"
	"strings"
	"tkestack.io/tke/cmd/tke-gateway/app/options"
	"tkestack.io/tke/pkg/apiserver"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/apiserver/handler"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	"tkestack.io/tke/pkg/auth"
	"tkestack.io/tke/pkg/gateway"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	gatewayconfigvalidation "tkestack.io/tke/pkg/gateway/apis/config/validation"
	"tkestack.io/tke/pkg/gateway/config/configfiles"
	"tkestack.io/tke/pkg/registry/distribution"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName             string
	GenericAPIServerConfig *genericapiserver.Config
	InsecureServingInfo    *genericapiserver.DeprecatedInsecureServingInfo
	OAuthConfig            *oauth2.Config
	OIDCHTTPClient         *http.Client
	OIDCAuthenticator      *oidc.Authenticator
	GatewayConfig          *gatewayconfig.GatewayConfiguration
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	gatewayConfig, err := options.NewGatewayConfiguration()
	if err != nil {
		log.Error("Failed create default gateway configuration", log.Err(err))
		return nil, err
	}

	// load config file, if provided
	if configFile := opts.GatewayConfig; len(configFile) > 0 {
		gatewayConfig, err = loadConfigFile(configFile)
		if err != nil {
			log.Error("Failed to load gateway configuration file", log.String("configFile", configFile), log.Err(err))
			return nil, err
		}
	}

	// We always validate the local configuration (command line + config file).
	// This is the default "last-known-good" config for dynamic config, and must always remain valid.
	if err := gatewayconfigvalidation.ValidateGatewayConfiguration(gatewayConfig); err != nil {
		klog.Fatal(err)
	}

	genericAPIServerConfig := genericapiserver.NewConfig(apiserver.Codecs)
	var ignoreAuthPathPrefixes []string
	ignoreAuthPathPrefixes = append(ignoreAuthPathPrefixes, distribution.IgnoreAuthPathPrefixes()...)
	ignoreAuthPathPrefixes = append(ignoreAuthPathPrefixes, auth.IgnoreAuthPathPrefixes()...)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(ignoreAuthPathPrefixes)
	genericAPIServerConfig.EnableIndex = false
	genericAPIServerConfig.EnableDiscovery = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	var insecureServingInfo *genericapiserver.DeprecatedInsecureServingInfo
	if err := opts.InsecureServing.ApplyTo(&insecureServingInfo); err != nil {
		return nil, err
	}

	setupAuthentication(genericAPIServerConfig)
	setupAuthorization(genericAPIServerConfig)

	var (
		oauthConfig       *oauth2.Config
		oidcHTTPClient    *http.Client
		oidcAuthenticator *oidc.Authenticator
	)

	externalAddress := fmt.Sprintf("%s:%d", opts.Generic.ExternalHost, opts.Generic.ExternalPort)
	oauthConfig, oidcHTTPClient, err = setupOIDC(opts.OIDC, externalAddress)
	if err != nil {
		return nil, err
	}

	oidcAuthenticator, err = setupOIDCClient(opts.OIDC)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerName:             serverName,
		GenericAPIServerConfig: genericAPIServerConfig,
		InsecureServingInfo:    insecureServingInfo,
		OAuthConfig:            oauthConfig,
		OIDCHTTPClient:         oidcHTTPClient,
		OIDCAuthenticator:      oidcAuthenticator,
		GatewayConfig:          gatewayConfig,
	}, nil
}

func setupAuthentication(genericAPIServerConfig *genericapiserver.Config) {
	genericAPIServerConfig.Authentication.SupportsBasicAuth = false
	genericAPIServerConfig.Authentication.Authenticator = anonymous.NewAuthenticator()
}

func setupAuthorization(genericAPIServerConfig *genericapiserver.Config) {
	genericAPIServerConfig.Authorization.Authorizer = authorizerfactory.NewAlwaysAllowAuthorizer()
}

func setupOIDCClient(oidcOpts *apiserveroptions.OIDCWithSecretOptions) (*oidc.Authenticator, error) {
	o := &oidc.Options{
		IssuerURL:            oidcOpts.IssuerURL,
		ExternalIssuerURL:    oidcOpts.ExternalIssuerURL,
		ClientID:             oidcOpts.ClientID,
		APIAudiences:         nil,
		CAFile:               oidcOpts.CAFile,
		UsernameClaim:        oidcOpts.UsernameClaim,
		UsernamePrefix:       oidcOpts.UsernamePrefix,
		GroupsClaim:          oidcOpts.GroupsClaim,
		GroupsPrefix:         oidcOpts.GroupsPrefix,
		TenantIDClaim:        oidcOpts.TenantIDClaim,
		TenantIDPrefix:       oidcOpts.TenantIDPrefix,
		SupportedSigningAlgs: oidcOpts.SigningAlgs,
		RequiredClaims:       oidcOpts.RequiredClaims,
	}
	return oidc.New(o)
}

func setupOIDC(oidcOpts *apiserveroptions.OIDCWithSecretOptions, externalAddress string) (*oauth2.Config, *http.Client, error) {
	// construct the cert pool
	tr, err := transport.NewOneWayTLSTransport(oidcOpts.CAFile, true)
	if err != nil {
		return nil, nil, err
	}
	hc := &http.Client{
		Transport: tr,
	}
	// construct the provider
	ctx := gooidc.ClientContext(context.Background(), hc)
	providerConfig, err := oidc.GetProviderConfig(ctx, oidcOpts.IssuerURL)
	if err != nil {
		log.Error("Failed to get the OIDC provider config", log.String("issuerURL", oidcOpts.IssuerURL), log.Err(err))
		return nil, nil, err
	}

	oauthConfig := &oauth2.Config{
		ClientID:     oidcOpts.ClientID,
		ClientSecret: oidcOpts.ClientSecret,
		Endpoint:     oauth2.Endpoint{AuthURL: providerConfig.AuthURL, TokenURL: strings.Replace(providerConfig.TokenURL, oidcOpts.ExternalIssuerURL, oidcOpts.IssuerURL, -1)},
		RedirectURL:  fmt.Sprintf("https://%s%s", externalAddress, gateway.CallbackPath),
		Scopes:       []string{gooidc.ScopeOpenID, gooidc.ScopeOfflineAccess, "profile", "email", "federated:id", "groups"},
	}
	return oauthConfig, hc, nil
}

func loadConfigFile(name string) (*gatewayconfig.GatewayConfiguration, error) {
	const errFmt = "failed to load Gateway config file %s, error %v"
	// compute absolute path based on current working dir
	gatewayConfigFile, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	loader, err := configfiles.NewFsLoader(utilfs.DefaultFs{}, gatewayConfigFile)
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	kc, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf(errFmt, name, err)
	}
	return kc, err
}
