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
	"net/http"
	"regexp"
	"strings"
	"time"
	"tkestack.io/tke/pkg/auth"

	"github.com/casbin/casbin"
	casbinlog "github.com/casbin/casbin/log"
	casbinutil "github.com/casbin/casbin/util"
	"github.com/coreos/etcd/clientv3"
	dexserver "github.com/dexidp/dex/server"
	"github.com/dexidp/dex/storage"
	"github.com/dexidp/dex/storage/etcd"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/group"
	"k8s.io/apiserver/pkg/authentication/request/bearertoken"
	"k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/request/websocket"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/healthz"
	"tkestack.io/tke/cmd/tke-auth/app/options"
	"tkestack.io/tke/pkg/apiserver"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"
	"tkestack.io/tke/pkg/auth/authorization/aggregation"
	casbinlogger "tkestack.io/tke/pkg/auth/logger"
	authopenapi "tkestack.io/tke/pkg/auth/openapi"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/adapter"
	"tkestack.io/tke/pkg/util/log"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Auth"

	policyRulesPrefix = "policy-rules"
)

var defaultModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub)  && keyMatchCustom(r.obj, p.obj) && keyMatchCustom(r.act, p.act)
`

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName             string
	OIDCExternalAddress    string
	GenericAPIServerConfig *genericapiserver.Config
	DexServer              *dexserver.Server
	CasbinEnforcer         *casbin.SyncedEnforcer
	Registry               *registry.Registry
	TokenAuthn             *authenticator.TokenAuthenticator
	APIKeyAuthn            *authenticator.APIKeyAuthenticator
	Authorizer             authorizer.Authorizer
	PolicyFile             string
	CategoryFile           string
	TenantAdmin            string
	TenantAdminSecret      string
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE auth command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	genericAPIServerConfig := genericapiserver.NewConfig(apiserver.Codecs)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(auth.IgnoreAuthPathPrefixes())
	genericAPIServerConfig.EnableIndex = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, authopenapi.GetOpenAPIDefinitions, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)
	debug.SetupDebug(genericAPIServerConfig, opts.Debug)

	etcdClient, err := setupETCDClient(genericAPIServerConfig, opts.ETCD)
	if err != nil {
		return nil, err
	}

	enforcer, err := setupCasbinEnforcer(etcdClient, opts.Authorization)
	if err != nil {
		return nil, err
	}

	tokenAuth := authenticator.NewTokenAuthenticator()
	if err := setupAuthentication(genericAPIServerConfig, opts.Authentication, tokenAuth); err != nil {
		return nil, err
	}

	aggregateAuthz, err := aggregation.NewAuthorizer(opts.Authorization, opts.Auth, enforcer)
	if err != nil {
		return nil, err
	}
	setupAuthorization(genericAPIServerConfig, aggregateAuthz)

	dexConfig, err := setupDexConfig(opts.ETCD, opts.Auth.AssetsPath, opts.Auth.IDTokenTimeout, opts.Generic.ExternalHost, opts.Generic.ExternalPort)
	if err != nil {
		return nil, err
	}

	r, err := registry.NewRegistry(etcdClient, dexConfig.Storage)
	if err != nil {
		return nil, err
	}

	apiKeyAuth, err := authenticator.NewAPIKeyAuthenticator(opts.Auth.APIKeySignMethod, r.APIKeyStorage(), r.LocalIdentityStorage())
	if err != nil {
		return nil, err
	}

	err = setupDefaultConnectorConfig(opts.ETCD, r, opts.Auth)
	if err != nil {
		return nil, err
	}

	err = setupDefaultClient(r, opts.Auth)
	if err != nil {
		return nil, err
	}

	dexServer, err := dexserver.NewServer(context.Background(), *dexConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerName:             serverName,
		OIDCExternalAddress:    dexConfig.Issuer,
		GenericAPIServerConfig: genericAPIServerConfig,
		DexServer:              dexServer,
		CasbinEnforcer:         enforcer,
		Registry:               r,
		TokenAuthn:             tokenAuth,
		APIKeyAuthn:            apiKeyAuth,
		Authorizer:             aggregateAuthz,
		CategoryFile:           opts.Auth.CategoryFile,
		PolicyFile:             opts.Auth.PolicyFile,
		TenantAdmin:            opts.Auth.TenantAdmin,
		TenantAdminSecret:      opts.Auth.TenantAdminSecret,
	}, nil
}

func setupAuthentication(genericAPIServerConfig *genericapiserver.Config, opts *apiserveroptions.AuthenticationWithAPIOptions, tokenAuth *authenticator.TokenAuthenticator) error {
	if err := authentication.SetupAuthentication(genericAPIServerConfig, opts); err != nil {
		return nil
	}

	configAuthenticators, configDefs := genericAPIServerConfig.Authentication.Authenticator, genericAPIServerConfig.OpenAPIConfig.SecurityDefinitions

	authenticators := []genericauthenticator.Request{
		configAuthenticators,
	}
	defs := *configDefs
	authenticators = append(authenticators, bearertoken.New(tokenAuth), websocket.NewProtocolAuthenticator(tokenAuth))
	defs["BearerToken"] = &spec.SecurityScheme{
		SecuritySchemeProps: spec.SecuritySchemeProps{
			Type:        "apiKey",
			Name:        "authorization",
			In:          "header",
			Description: "Bearer Token authentication",
		},
	}

	finalAuthenticator := union.New(authenticators...)
	finalAuthenticator = group.NewAuthenticatedGroupAdder(finalAuthenticator)

	genericAPIServerConfig.Authentication.Authenticator, genericAPIServerConfig.OpenAPIConfig.SecurityDefinitions = finalAuthenticator, &defs
	return nil
}

func setupAuthorization(genericAPIServerConfig *genericapiserver.Config, authorizer authorizer.Authorizer) {
	genericAPIServerConfig.Authorization.Authorizer = authorizer
}

func setupETCDClient(genericAPIServerConfig *genericapiserver.Config, etcdOpts *storageoptions.ETCDStorageOptions) (*clientv3.Client, error) {
	client, err := etcdOpts.NewClient()
	if err != nil {
		return nil, err
	}
	healthCheck, err := etcdOpts.NewHealthCheck()
	if err != nil {
		return nil, err
	}
	genericAPIServerConfig.HealthzChecks = append(genericAPIServerConfig.HealthzChecks, healthz.NamedCheck("etcd", func(r *http.Request) error {
		return healthCheck()
	}))
	return client, nil
}

func setupDexConfig(etcdOpts *storageoptions.ETCDStorageOptions, templatePath string, tokenTimeout time.Duration, host string, port int) (*dexserver.Config, error) {
	logger := logrus.NewLogger(log.ZapLogger())
	issuer := issuer(host, port)
	namespace := etcdOpts.Prefix
	if !strings.HasSuffix(namespace, "/") {
		namespace = namespace + "/"
	}
	opts := etcd.Etcd{
		Endpoints: etcdOpts.ServerList,
		Namespace: namespace,
		SSL: etcd.SSL{
			CAFile:   etcdOpts.CAFile,
			KeyFile:  etcdOpts.KeyFile,
			CertFile: etcdOpts.CertFile,
		},
	}

	store, err := opts.Open(logger)
	if err != nil {
		return nil, err
	}

	dexConfig := dexserver.Config{
		Issuer:  issuer,
		Logger:  logger,
		Storage: store,

		Web: dexserver.WebConfig{
			Dir:   templatePath,
			Theme: "tkestack",
		},
		IDTokensValidFor:   tokenTimeout,
		SkipApprovalScreen: true,

		PrometheusRegistry: prometheus.NewRegistry(),
	}

	return &dexConfig, nil
}

// issuer returns issuer location of the OIDC server.
func issuer(advertiseAddress string, advertisePort int) string {
	var scheme, port string
	scheme = "https"
	if advertisePort != 443 {
		port = fmt.Sprintf(":%d", advertisePort)
	}
	return fmt.Sprintf("%s://%s%s/%s", scheme, advertiseAddress, port, types.IssuerName)
}

func setupCasbinEnforcer(etcdClient *clientv3.Client, authorizationOptions *options.AuthorizationOptions) (*casbin.SyncedEnforcer, error) {
	var enforcer *casbin.SyncedEnforcer
	var err error
	adpt := adapter.NewAdapter(etcdClient, policyRulesPrefix)
	if len(authorizationOptions.CasbinModelFile) == 0 {
		m := casbin.NewModel(defaultModel)
		enforcer, err = casbin.NewSyncedEnforcerSafe(m, adpt)
	} else {
		enforcer, err = casbin.NewSyncedEnforcerSafe(authorizationOptions.CasbinModelFile, adpt)
	}

	if err != nil {
		return nil, err
	}

	if authorizationOptions.Debug {
		casbinlog.SetLogger(&casbinlogger.WrapLogger{})
		enforcer.EnableLog(true)
	}

	enforcer.AddFunction("keyMatchCustom", CustomFunctionWrapper)

	enforcer.StartAutoLoadPolicy(authorizationOptions.CasbinReloadInterval)

	return enforcer, nil
}

func setupDefaultConnectorConfig(etcdOpts *storageoptions.ETCDStorageOptions, r *registry.Registry, auth *options.AuthOptions) error {
	// create dex local identity provider for tke connector.
	dexserver.ConnectorsConfig[local.TkeConnectorType] = func() dexserver.ConnectorConfig {
		return new(local.Config)
	}

	conns, err := r.DexStorage().ListConnectors()
	if err != nil {
		return err
	}

	exists := false
	for _, conn := range conns {
		if conn.ID == auth.InitTenantID {
			exists = true
			continue
		}
	}
	if !exists {
		tkeConn, err := local.NewLocalConnector(etcdOpts.ETCDClientOptions, auth.InitTenantID)
		if err != nil {
			return err
		}
		// if no connectors, create a default connector
		if err = r.DexStorage().CreateConnector(*tkeConn); err != nil {
			return err
		}
	}

	return nil
}

func setupDefaultClient(r *registry.Registry, auth *options.AuthOptions) error {
	clis, err := r.DexStorage().ListClients()
	if err != nil {
		return err
	}

	exists := false
	for _, cli := range clis {
		if cli.ID == auth.InitClientID {
			exists = true
			continue
		}
	}
	if !exists {
		cli := storage.Client{
			ID:           auth.InitClientID,
			Secret:       auth.InitClientSecret,
			Name:         auth.InitClientID,
			RedirectURIs: auth.InitClientRedirectUris,
		}

		// if no connectors, create a default connector
		if err = r.DexStorage().CreateClient(cli); err != nil {
			return err
		}
	}

	return nil
}

// CustomFunctionWrapper wraps keyMatchCustomFunction
func CustomFunctionWrapper(args ...interface{}) (interface{}, error) {
	key1 := args[0].(string)
	key2 := args[1].(string)

	return keyMatchCustomFunction(key1, key2), nil
}

// keyMatchCustomFunction determines whether key1 matches the pattern of key2 , key2 can contain a * and :*.
// For example, "/project:123/cluster:456" matches "/project:*/cluster:456", "registry:123/*" matches "registry:123/456"
func keyMatchCustomFunction(key1 string, key2 string) bool {
	// case insensitive
	key1 = strings.ToLower(key1)
	key2 = strings.ToLower(key2)

	key2 = strings.Replace(key2, "*", ".*", -1)

	re := regexp.MustCompile(`(.*):[^/]+(.*)`)
	i := 2
	for {
		if !strings.Contains(key2, "/:") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
		fmt.Printf("%d %s\n", i, key2)
		i = i + 1
	}

	return casbinutil.RegexMatch(key1, "^"+key2+"$")
}
