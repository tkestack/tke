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

package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/ldap"
	local2 "tkestack.io/tke/pkg/auth/authorization/local"

	dexstorage "github.com/dexidp/dex/storage"
	"github.com/emicklei/go-restful"
	"k8s.io/apiserver/pkg/server/mux"

	"github.com/casbin/casbin/v2"
	dexserver "github.com/dexidp/dex/server"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"

	authv1 "tkestack.io/tke/api/auth/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	"tkestack.io/tke/pkg/apiserver/storage"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"
	authnhandler "tkestack.io/tke/pkg/auth/handler/authn"
	authzhandler "tkestack.io/tke/pkg/auth/handler/authz"
	authrest "tkestack.io/tke/pkg/auth/registry/rest"
	"tkestack.io/tke/pkg/auth/route"
	"tkestack.io/tke/pkg/util/log"
)

func IgnoreAuthPathPrefixes() []string {
	return []string{
		"/oidc/",
		"/auth/",
		"/apis/auth.tkestack.io/v1/apikeys/default/password",
	}
}

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory

	OIDCExternalAddress  string
	DexConfig            *dexserver.Config
	DexStorage           dexstorage.Storage
	CasbinEnforcer       *casbin.SyncedEnforcer
	TokenAuthn           *authenticator.TokenAuthenticator
	APIKeyAuthn          *authenticator.APIKeyAuthenticator
	Authorizer           authorizer.Authorizer
	CasbinReloadInterval time.Duration
	TenantID             string
	TenantAdmin          string
	TenantAdminSecret    string
	PrivilegedUsername   string
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

// APIServer contains state for a tke api server.
type APIServer struct {
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
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*APIServer, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	dexHandler := identityprovider.DexHander{}

	hooks := c.registerHooks(&dexHandler, s)
	installHooks(s, hooks)
	installCasbinPreStopHook(s, c.ExtraConfig.CasbinEnforcer)

	c.registerRoute(&dexHandler, s.Handler.GoRestfulContainer, s.Handler.NonGoRestfulMux)

	m := &APIServer{
		GenericAPIServer: s,
	}

	// The order here is preserved in discovery.
	restStorageProviders := []storage.RESTStorageProvider{
		&authrest.StorageProvider{
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			Enforcer:             c.ExtraConfig.CasbinEnforcer,
			DexStorage:           c.ExtraConfig.DexStorage,
			PrivilegedUsername:   c.ExtraConfig.PrivilegedUsername,
		},
	}
	m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)

	return m, nil
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *APIServer) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...storage.RESTStorageProvider) {
	var apiGroupsInfo []genericapiserver.APIGroupInfo

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyVersionForGroupEnabled(groupName) {
			log.Infof("Skipping disabled API group %q.", groupName)
			continue
		}
		apiGroupInfo, enabled := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if !enabled {
			log.Warnf("Problem initializing API group %q, skipping.", groupName)
			continue
		}
		log.Infof("Enabling API group %q.", groupName)

		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				log.Fatalf("Error building PostStartHook: %v", err)
			}
			m.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}

		apiGroupsInfo = append(apiGroupsInfo, apiGroupInfo)
	}

	for i := range apiGroupsInfo {
		if err := m.GenericAPIServer.InstallAPIGroup(&apiGroupsInfo[i]); err != nil {
			log.Fatalf("Error in registering group versions: %v", err)
		}
	}
}

// DefaultAPIResourceConfigSource returns which groupVersion enabled and its
// resources enabled/disabled.
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	ret.EnableVersions(
		authv1.SchemeGroupVersion,
	)
	return ret
}

// registerRoute is used to register routes with the api server of project.
func (c completedConfig) registerRoute(dexHandler http.Handler, container *restful.Container, mux *mux.PathRecorderMux) {
	mux.HandlePrefix("/"+auth.IssuerName+"/", dexHandler)

	token := authnhandler.NewHandler(c.ExtraConfig.TokenAuthn, c.ExtraConfig.APIKeyAuthn)
	authz := authzhandler.NewHandler(c.ExtraConfig.Authorizer)
	route.RegisterAuthRoute(container, token, authz)
}

// registerHooks is used to register postStart hook to create authn provider with local oidc server.
func (c completedConfig) registerHooks(dexHandler *identityprovider.DexHander, s *genericapiserver.GenericAPIServer) []genericapiserver.PostStartHookProvider {

	authClient := authinternalclient.NewForConfigOrDie(s.LoopbackClientConfig)

	dexHook := identityprovider.NewDexHookHandler(context.Background(), c.ExtraConfig.DexConfig, c.ExtraConfig.DexStorage, dexHandler,
		c.ExtraConfig.OIDCExternalAddress, fmt.Sprintf("%s/%s", s.LoopbackClientConfig.Host, auth.IssuerName), c.ExtraConfig.TokenAuthn)

	apiSigningKeyHook := authenticator.NewAPISigningKeyHookHandler(authClient)

	identityHook := authenticator.NewAdminIdentityHookHandler(authClient, c.ExtraConfig.TenantID, c.ExtraConfig.TenantAdmin, c.ExtraConfig.TenantAdminSecret)

	localIdpHook := local.NewLocalHookHandler(authClient, c.ExtraConfig.VersionedInformers)
	ldapIdpHook := ldap.NewLdapHookHandler(authClient)

	authVersionedClient := versionedclientset.NewForConfigOrDie(s.LoopbackClientConfig)
	adapterHook := local2.NewAdapterHookHandler(authVersionedClient, c.ExtraConfig.CasbinEnforcer, c.ExtraConfig.VersionedInformers, c.ExtraConfig.CasbinReloadInterval)

	return []genericapiserver.PostStartHookProvider{dexHook, apiSigningKeyHook, identityHook, localIdpHook, ldapIdpHook, adapterHook}
}

// installCasbinPreStopHook is used to register preStop hook to stop casbin enforcer sync.
func installCasbinPreStopHook(s *genericapiserver.GenericAPIServer, enforcer *casbin.SyncedEnforcer) {
	s.AddPreShutdownHookOrDie("stop-casbin-enforcer-sync", func() error {
		enforcer.StopAutoLoadPolicy()
		return nil
	})
}

func installHooks(s *genericapiserver.GenericAPIServer, hooks []genericapiserver.PostStartHookProvider) {
	for _, hookProvider := range hooks {
		name, hook, err := hookProvider.PostStartHook()
		if err != nil {
			log.Fatal("Failed to install the post start hook", log.Err(err))
		}
		s.AddPostStartHookOrDie(name, hook)
	}
}
