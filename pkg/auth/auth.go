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

package auth

import (
	"context"
	"fmt"
	"github.com/casbin/casbin"
	dexserver "github.com/dexidp/dex/server"
	"github.com/emicklei/go-restful"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/auth/authentication/tenant"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	apikeyhandler "tkestack.io/tke/pkg/auth/handler/apikey"
	authnhandler "tkestack.io/tke/pkg/auth/handler/authn"
	authzhandler "tkestack.io/tke/pkg/auth/handler/authz"
	categoryhandler "tkestack.io/tke/pkg/auth/handler/category"
	clihandler "tkestack.io/tke/pkg/auth/handler/client"
	idphandler "tkestack.io/tke/pkg/auth/handler/identityprovider"
	identityhander "tkestack.io/tke/pkg/auth/handler/localidentity"
	policyhandler "tkestack.io/tke/pkg/auth/handler/policy"
	rolehandler "tkestack.io/tke/pkg/auth/handler/role"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/route"
	"tkestack.io/tke/pkg/auth/types"
)

func IgnoreAuthPathPrefixes() []string {
	return []string{
		"/oidc/",
		"/auth/",
		"/api/authv1/apikey/password",
	}
}

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName          string
	OIDCExternalAddress string
	DexServer           *dexserver.Server
	CasbinEnforcer      *casbin.SyncedEnforcer
	Registry            *registry.Registry
	TokenAuthn          *authenticator.TokenAuthenticator
	APIKeyAuthn         *authenticator.APIKeyAuthenticator
	Authorizer          authorizer.Authorizer
	CategoryFile        string
	PolicyFile          string
	TenantAdmin         string
	TenantAdminSecret   string
}

// Config contains the core configuration instance of apiserver and
// additional configuration.
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embed a private pointer of Config.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Auth contains state for a tke auth.
type Auth struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	return CompletedConfig{&c}
}

// New returns a new instance of APIServer from the given config.
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*Auth, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	c.registerRoute(s.Handler.GoRestfulContainer, s.Handler.NonGoRestfulMux)

	if err := c.registerAuthnHook(s); err != nil {
		return nil, err
	}

	if err := c.registerCasbinPreStopHook(s); err != nil {
		return nil, err
	}

	m := &Auth{
		GenericAPIServer: s,
	}

	return m, nil
}

// registerRoute is used to register routes with the api server of project.
func (c completedConfig) registerRoute(container *restful.Container, mux *mux.PathRecorderMux) {
	route.RegisterOIDCRoute(mux, c.ExtraConfig.DexServer)

	policyEnforcer := enforcer.NewPolicyEnforcer(c.ExtraConfig.CasbinEnforcer, c.ExtraConfig.Registry)

	identity := identityhander.NewHandler(c.ExtraConfig.Registry, policyEnforcer)
	route.RegisterIdentityRoute(container, identity)

	policy := policyhandler.NewHandler(c.ExtraConfig.Registry, policyEnforcer, c.ExtraConfig.TenantAdmin)
	route.RegisterPolicyRoute(container, policy)

	role := rolehandler.NewHandler(c.ExtraConfig.Registry, policyEnforcer, c.ExtraConfig.TenantAdmin)
	route.RegisterRoleRoute(container, role)

	category := categoryhandler.NewHandler(c.ExtraConfig.Registry)
	route.RegisterCategoryRoute(container, category)

	token := authnhandler.NewHandler(c.ExtraConfig.TokenAuthn, c.ExtraConfig.APIKeyAuthn)
	authz := authzhandler.NewHandler(c.ExtraConfig.Authorizer)
	route.RegisterAuthRoute(container, token, authz)

	apiKey := apikeyhandler.NewHandler(c.ExtraConfig.APIKeyAuthn)
	route.RegisterAPIKeyRoute(container, apiKey)

	helper := tenant.NewHelper(c.ExtraConfig.Registry, policy.Service(), c.ExtraConfig.PolicyFile, c.ExtraConfig.CategoryFile, c.ExtraConfig.TenantAdmin, c.ExtraConfig.TenantAdminSecret)
	idp := idphandler.NewHandler(c.ExtraConfig.Registry.DexStorage(), helper)
	route.RegisterIdentityProviderRoute(container, idp)

	cli := clihandler.NewHandler(c.ExtraConfig.Registry.DexStorage())
	route.RegisterClientRoute(container, cli)
}

// registerAuthnHook is used to register postStart hook to create authn provider with local oidc server.
func (c completedConfig) registerAuthnHook(s *genericapiserver.GenericAPIServer) error {
	authnProvider := authenticator.NewProviderHookHandler(context.Background(), c.ExtraConfig.OIDCExternalAddress, fmt.Sprintf("%s/%s", s.LoopbackClientConfig.Host, types.IssuerName), c.ExtraConfig.TokenAuthn)
	name, hook, err := authnProvider.PostStartHook()
	if err != nil {
		return err
	}

	return s.AddPostStartHook(name, hook)
}

// registerCasbinPreStopHook is used to register preStop hook to stop casbin enforcer sync.
func (c completedConfig) registerCasbinPreStopHook(s *genericapiserver.GenericAPIServer) error {
	return s.AddPreShutdownHook("stop-casbin-enforcer-sync", func() error {
		c.ExtraConfig.CasbinEnforcer.StopAutoLoadPolicy()
		return nil
	})
}
