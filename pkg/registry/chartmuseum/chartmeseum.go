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

package chartmuseum

import (
	"net/http"

	"helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/server/mux"
	restclient "k8s.io/client-go/rest"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/chartmuseum/authentication"
	"tkestack.io/tke/pkg/registry/chartmuseum/authorization"
	"tkestack.io/tke/pkg/registry/chartmuseum/request"
	"tkestack.io/tke/pkg/registry/chartmuseum/serveroptions"
	"tkestack.io/tke/pkg/registry/chartmuseum/tenant"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// PathPrefix defines the path prefix for accessing the chartmuseum server.
	PathPrefix = "/chart/"
	// MaxUploadSize defines max size of post body (in bytes).
	MaxUploadSize = 20 * 1024 * 1024 * 1024
)

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoredAuthPathPrefixes() []string {
	return []string{
		PathPrefix,
	}
}

type Options struct {
	RegistryConfig       *registryconfig.RegistryConfiguration
	LoopbackClientConfig *restclient.Config
	OIDCIssuerURL        string
	OIDCTokenReviewPath  string
	OIDCCAFile           string
	ExternalScheme       string
	Authorizer           authorizer.Authorizer
}

// RegisterRoute to register the chartmuseum server path prefix to apiserver.
func RegisterRoute(m *mux.PathRecorderMux, opts *Options) error {
	chartMuseumConfig, err := serveroptions.BuildChartMuseumConfig(opts.RegistryConfig, PathPrefix, MaxUploadSize)
	if err != nil {
		log.Error("Failed to initialize chartmuseum server configuration", log.Err(err))
		return err
	}
	multiTenantServer, err := multitenant.NewMultiTenantServer(*chartMuseumConfig)
	if err != nil {
		log.Error("Failed to create chartmuseum server", log.Err(err))
		return err
	}

	// add handler chain
	var handler http.Handler
	if opts.RegistryConfig.Security.EnableAnonymous == nil || !*opts.RegistryConfig.Security.EnableAnonymous {
		var chainErr error
		handler, chainErr = authorization.WithAuthorization(multiTenantServer.Router, &authorization.Options{
			AdminUsername:  opts.RegistryConfig.Security.AdminUsername,
			ExternalScheme: opts.ExternalScheme,
			LoopbackConfig: opts.LoopbackClientConfig,
			Authorizer:     opts.Authorizer,
		})
		if chainErr != nil {
			return chainErr
		}
		handler, chainErr = authentication.WithAuthentication(handler, &authentication.Options{
			SecurityConfig:  &opts.RegistryConfig.Security,
			ExternalScheme:  opts.ExternalScheme,
			OIDCIssuerURL:   opts.OIDCIssuerURL,
			OIDCCAFile:      opts.OIDCCAFile,
			TokenReviewPath: opts.OIDCTokenReviewPath,
		})
		if chainErr != nil {
			return chainErr
		}
	}
	handler = tenant.WithTenant(handler, PathPrefix, opts.RegistryConfig.DomainSuffix, opts.RegistryConfig.DefaultTenant)
	handler = request.WithRequestID(handler)
	m.HandlePrefix(PathPrefix, handler)

	return nil
}
