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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/cloudindustry"

	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/ldap"
	local2 "tkestack.io/tke/pkg/auth/authorization/local"

	dexstorage "github.com/dexidp/dex/storage"
	"github.com/emicklei/go-restful"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	"html/template"
)

const (
	OIDCPath           = "/oidc/"
	AuthPath           = "/auth/"
	APIKeyPasswordPath = "/apis/auth.tkestack.io/v1/apikeys/default/password"

	APIKeyPath           = "/apis/auth.tkestack.io/v1/apikeys"
	DefaultTitle         = "TKEStack"
	DefaultLogoDir       = "default"
	HtmlTmplDir          = "web/auth/templates/"
	FlagConsoleTitle     = "title"
	FlagConsoleLogoDir   = "logo-dir"
	ConfigConsoleTitle   = "console_config.title"
	ConfigConsoleLogoDir = "console_config.logo_dir"
)

func IgnoreAuthPathPrefixes() []string {
	return []string{
		OIDCPath,
		AuthPath,
		APIKeyPasswordPath,
	}
}

func IgnoreAuthzPathPrefixes() []string {
	return []string{
		APIKeyPath,
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
	PrivilegedUsername   string
	ConsoleConfig        *ConsoleConfig
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

type ConsoleConfig struct {
	Title   string
	LogoDir string
}

// NewAuthOptions creates a AuthOptions object with default parameters.
func NewConsoleConfigOptions() *ConsoleConfig {
	return &ConsoleConfig{}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *ConsoleConfig) AddFlags(fs *pflag.FlagSet) {
	fs.String(FlagConsoleTitle, o.Title,
		"Custom console title.")
	_ = viper.BindPFlag(ConfigConsoleTitle, fs.Lookup(FlagConsoleTitle))

	fs.String(FlagConsoleLogoDir, o.LogoDir,
		"Custom console logo dir.")
	_ = viper.BindPFlag(ConfigConsoleLogoDir, fs.Lookup(FlagConsoleLogoDir))

}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ConsoleConfig) ApplyFlags() []error {
	var errs []error

	o.Title = viper.GetString(ConfigConsoleTitle)
	o.LogoDir = viper.GetString(ConfigConsoleLogoDir)

	return errs
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

	consoleConfig := new(ConsoleConfig)

	if c.ExtraConfig.ConsoleConfig != nil {
		consoleConfig = c.ExtraConfig.ConsoleConfig
	} else {
		consoleConfig.Title = DefaultTitle
		consoleConfig.LogoDir = DefaultLogoDir

	}

	files, err := ioutil.ReadDir(HtmlTmplDir)
	if err != nil {
		return nil, err
	}

	sourceRe := regexp.MustCompile(`\.tmpl\.html$`)
	targetRe := regexp.MustCompile(`\.tmpl`)

	for _, file := range files {
		var buf bytes.Buffer
		if !sourceRe.MatchString(file.Name()) {
			continue
		}
		t, err := template.New(file.Name()).Delims("{%", "%}").ParseFiles(HtmlTmplDir + file.Name())
		if err != nil {
			return nil, err
		}
		if err = t.Execute(&buf, c.ExtraConfig.ConsoleConfig); err != nil {
			return nil, err
		}
		// // remove .tmpl in file name
		targetFileName := targetRe.ReplaceAllString(file.Name(), "")
		if err = ioutil.WriteFile(HtmlTmplDir+targetFileName, buf.Bytes(), 0644); err != nil {
			return nil, err
		}
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
			VersionedInformers:   c.ExtraConfig.VersionedInformers,
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
		if !apiResourceConfigSource.AnyResourceForGroupEnabled(groupName) {
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

	dexHook := identityprovider.NewDexHookHandler(context.Background(), authClient, c.ExtraConfig.DexConfig, c.ExtraConfig.DexStorage, dexHandler,
		c.ExtraConfig.OIDCExternalAddress, fmt.Sprintf("%s/%s", s.LoopbackClientConfig.Host, auth.IssuerName), c.ExtraConfig.TokenAuthn)

	apiSigningKeyHook := authenticator.NewAPISigningKeyHookHandler(authClient)

	localIdpHook := local.NewLocalHookHandler(authClient)
	ldapIdpHook := ldap.NewLdapHookHandler(authClient)
	cloudIndustryIdpHook := cloudindustry.NewCloudIndustryHookHandler(authClient)

	authVersionedClient := versionedclientset.NewForConfigOrDie(s.LoopbackClientConfig)
	adapterHook := local2.NewAdapterHookHandler(authVersionedClient, c.ExtraConfig.CasbinEnforcer, c.ExtraConfig.VersionedInformers, c.ExtraConfig.CasbinReloadInterval)

	return []genericapiserver.PostStartHookProvider{dexHook, apiSigningKeyHook, localIdpHook, ldapIdpHook, cloudIndustryIdpHook, adapterHook}
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
