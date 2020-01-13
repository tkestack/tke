/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	casbinlog "github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/model"
	casbinutil "github.com/casbin/casbin/v2/util"
	dexldap "github.com/dexidp/dex/connector/ldap"
	dexserver "github.com/dexidp/dex/server"
	dexstorage "github.com/dexidp/dex/storage"
	"github.com/dexidp/dex/storage/etcd"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/group"
	"k8s.io/apiserver/pkg/authentication/request/bearertoken"
	"k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/authentication/request/websocket"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"

	authapi "tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	generatedopenapi "tkestack.io/tke/api/openapi"
	"tkestack.io/tke/cmd/tke-auth-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/debug"
	"tkestack.io/tke/pkg/apiserver/handler"
	"tkestack.io/tke/pkg/apiserver/openapi"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	"tkestack.io/tke/pkg/apiserver/storage"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
	"tkestack.io/tke/pkg/auth/apiserver"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/ldap"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"
	"tkestack.io/tke/pkg/auth/authorization/aggregation"
	dexutil "tkestack.io/tke/pkg/auth/util/dex"
	casbinlogger "tkestack.io/tke/pkg/auth/util/logger"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/log/dex"
)

const (
	license = "Apache 2.0"
	title   = "Tencent Kubernetes Engine Auth API"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName                     string
	OIDCExternalAddress            string
	GenericAPIServerConfig         *genericapiserver.Config
	VersionedSharedInformerFactory versionedinformers.SharedInformerFactory
	StorageFactory                 *serverstorage.DefaultStorageFactory

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

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE auth command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	genericAPIServerConfig := genericapiserver.NewConfig(authapi.Codecs)
	genericAPIServerConfig.BuildHandlerChainFunc = handler.BuildHandlerChain(apiserver.IgnoreAuthPathPrefixes())
	genericAPIServerConfig.MergedResourceConfig = apiserver.DefaultAPIResourceConfigSource()

	genericAPIServerConfig.EnableIndex = false

	if err := opts.Generic.ApplyTo(genericAPIServerConfig); err != nil {
		return nil, err
	}
	if err := opts.SecureServing.ApplyTo(&genericAPIServerConfig.SecureServing, &genericAPIServerConfig.LoopbackClientConfig); err != nil {
		return nil, err
	}

	openapi.SetupOpenAPI(genericAPIServerConfig, generatedopenapi.GetOpenAPIDefinitions, title, license, opts.Generic.ExternalHost, opts.Generic.ExternalPort)
	debug.SetupDebug(genericAPIServerConfig, opts.Debug)

	// storageFactory
	storageFactoryConfig := storage.NewFactoryConfig(authapi.Codecs, authapi.Scheme)
	storageFactoryConfig.APIResourceConfig = genericAPIServerConfig.MergedResourceConfig
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(opts.ETCD)
	if err != nil {
		return nil, err
	}
	storageFactory, err := completedStorageFactoryConfig.New()
	if err != nil {
		return nil, err
	}
	if err := opts.ETCD.ApplyWithStorageFactoryTo(storageFactory, genericAPIServerConfig); err != nil {
		return nil, err
	}

	// client config
	genericAPIServerConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"

	kubeClientConfig := genericAPIServerConfig.LoopbackClientConfig
	clientgoExternalClient, err := versionedclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create real external clientset: %v", err)
	}
	versionedInformers := versionedinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	enforcer, err := setupCasbinEnforcer(opts.Authorization)
	if err != nil {
		return nil, err
	}

	authClient := authinternalclient.NewForConfigOrDie(genericAPIServerConfig.LoopbackClientConfig)
	tokenAuth := authenticator.NewTokenAuthenticator(authClient)
	if err := setupAuthentication(genericAPIServerConfig, opts.Authentication, tokenAuth); err != nil {
		return nil, err
	}

	aggregateAuthz, err := aggregation.NewAuthorizer(authClient, opts.Authorization, opts.Auth, enforcer, opts.Authentication.PrivilegedUsername)
	if err != nil {
		return nil, err
	}
	setupAuthorization(genericAPIServerConfig, aggregateAuthz)

	dexConfig, err := setupDexConfig(opts.ETCD, authClient, opts.Auth.AssetsPath, opts.Auth.IDTokenTimeout, opts.Generic.ExternalHost, opts.Generic.ExternalPort)
	if err != nil {
		return nil, err
	}

	apiKeyAuth, err := authenticator.NewAPIKeyAuthenticator(authClient)
	if err != nil {
		return nil, err
	}

	// create dex local identity provider for tke connector.
	dexserver.ConnectorsConfig[local.ConnectorType] = func() dexserver.ConnectorConfig {
		return new(local.DefaultIdentityProvider)
	}
	local.SetupRestClient(authClient)
	log.Info("init tenant type", log.String("type", opts.Auth.InitTenantType))
	switch opts.Auth.InitTenantType {
	case local.ConnectorType:
		err = setupDefaultConnector(versionedInformers, opts.Auth)
		if err != nil {
			return nil, err
		}
	case ldap.ConnectorType:
		err = setupLdapConnector(opts.Auth)
		if err != nil {
			return nil, err
		}
	default:
		log.Warn("Unknown init tenant type", log.String("type", opts.Auth.InitTenantType))
	}

	err = setupDefaultClient(dexConfig.Storage, opts.Auth)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerName:                     serverName,
		OIDCExternalAddress:            dexConfig.Issuer,
		GenericAPIServerConfig:         genericAPIServerConfig,
		StorageFactory:                 storageFactory,
		VersionedSharedInformerFactory: versionedInformers,
		DexConfig:                      dexConfig,
		DexStorage:                     dexConfig.Storage,
		CasbinEnforcer:                 enforcer,
		TokenAuthn:                     tokenAuth,
		APIKeyAuthn:                    apiKeyAuth,
		Authorizer:                     aggregateAuthz,
		TenantID:                       opts.Auth.InitTenantID,
		TenantAdmin:                    opts.Auth.TenantAdmin,
		TenantAdminSecret:              opts.Auth.TenantAdminSecret,
		PrivilegedUsername:             opts.Authentication.PrivilegedUsername,
		CasbinReloadInterval:           opts.Authorization.CasbinReloadInterval,
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

func setupDexConfig(etcdOpts *storageoptions.ETCDStorageOptions, authClient authinternalclient.AuthInterface, templatePath string, tokenTimeout time.Duration, host string, port int) (*dexserver.Config, error) {
	logger := dex.NewLogger(log.ZapLogger())
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

	wrapOpts := dexutil.NewWrapEtcdStorage(opts, authClient)
	store, err := wrapOpts.Open(logger)
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
	return fmt.Sprintf("%s://%s%s/%s", scheme, advertiseAddress, port, authapi.IssuerName)
}

func setupCasbinEnforcer(authorizationOptions *options.AuthorizationOptions) (*casbin.SyncedEnforcer, error) {
	var enforcer *casbin.SyncedEnforcer
	var err error
	if len(authorizationOptions.CasbinModelFile) == 0 {
		m, err := model.NewModelFromString(authapi.DefaultRuleModel)
		if err != nil {
			return nil, err
		}
		enforcer, err = casbin.NewSyncedEnforcer(m)
		if err != nil {
			return nil, err
		}
	} else {
		enforcer, err = casbin.NewSyncedEnforcer(authorizationOptions.CasbinModelFile)
		if err != nil {
			return nil, err
		}
	}

	if authorizationOptions.Debug {
		casbinlog.SetLogger(&casbinlogger.WrapLogger{})
		enforcer.EnableLog(true)
	}

	enforcer.AddFunction("keyMatchCustom", CustomFunctionWrapper)

	return enforcer, nil
}

func setupDefaultConnector(versionInformers versionedinformers.SharedInformerFactory, auth *options.AuthOptions) error {
	if _, ok := identityprovider.IdentityProvidersStore[auth.InitTenantID]; !ok {
		defaultIDP := local.NewDefaultIdentityProvider(auth.InitTenantID, versionInformers)
		identityprovider.IdentityProvidersStore[auth.InitTenantID] = defaultIDP
	}

	return nil
}

func setupLdapConnector(auth *options.AuthOptions) error {
	log.Info("setup ldap connector", log.Any("auth", auth))
	const errFmt = "failed to load Ldap config file %s, error %v"
	// compute absolute path based on current working dir
	ldapConfigFile, err := filepath.Abs(auth.LdapConfigFile)
	if err != nil {
		return fmt.Errorf(errFmt, ldapConfigFile, err)
	}

	bytes, err := ioutil.ReadFile(ldapConfigFile)
	var ldapConfig dexldap.Config
	if err := json.Unmarshal(bytes, &ldapConfig); err != nil {
		return fmt.Errorf(errFmt, ldapConfigFile, err)
	}

	idp, err := ldap.NewLDAPIdentityProvider(ldapConfig, auth.InitTenantID)
	if err != nil {
		return err
	}

	if _, ok := identityprovider.IdentityProvidersStore[auth.InitTenantID]; !ok {
		identityprovider.IdentityProvidersStore[auth.InitTenantID] = idp
	}

	return nil
}

func setupDefaultClient(store dexstorage.Storage, auth *options.AuthOptions) error {
	clis, err := store.ListClients()
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
		cli := dexstorage.Client{
			ID:           auth.InitClientID,
			Secret:       auth.InitClientSecret,
			Name:         auth.InitClientID,
			RedirectURIs: auth.InitClientRedirectUris,
		}

		// Create a default connector
		if err = store.CreateClient(cli); err != nil {
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
