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
	"net/http"
	"strings"

	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/libtrust"
	jsoniter "github.com/json-iterator/go"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	restclient "k8s.io/client-go/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/apikey"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/distribution/tenant"
	authenticationutil "tkestack.io/tke/pkg/registry/util/authentication"
	utilregistryrequest "tkestack.io/tke/pkg/registry/util/request"
	"tkestack.io/tke/pkg/util/log"
)

// Path is the authorization server url for distribution.
const Path = "/registry/auth"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Options defines the configuration of distribution auth handler.
type Options struct {
	SecurityConfig    *registryconfig.Security
	TokenReviewCAFile string
	TokenReviewPath   string
	DomainSuffix      string
	DefaultTenant     string
	LoopbackConfig    *restclient.Config
}

type handler struct {
	filterMap     map[string]accessFilter
	privateKey    libtrust.PrivateKey
	expiredHours  int64
	domainSuffix  string
	defaultTenant string
	authenticator authenticator.Request
}

// NewHandler creates a new handler object and returns it.
func NewHandler(opts *Options) (http.Handler, error) {
	if opts.SecurityConfig.AdminUsername == "" || opts.SecurityConfig.AdminPassword == "" {
		log.Warn("System image cannot be managed without setting an administrator password")
	}

	registryClient, err := registryinternalclient.NewForConfig(opts.LoopbackConfig)
	if err != nil {
		return nil, err
	}

	pk, err := libtrust.LoadKeyFile(opts.SecurityConfig.TokenPrivateKeyFile)
	if err != nil {
		return nil, err
	}

	if opts.SecurityConfig.TokenExpiredHours == nil {
		return nil, fmt.Errorf("token expired hours must be specify")
	}

	at, err := apikey.NewAPIKeyAuthenticator(&apikey.Options{
		TokenReviewCAFile: opts.TokenReviewCAFile,
		TokenReviewURL:    opts.TokenReviewPath,
		AdminUsername:     opts.SecurityConfig.AdminUsername,
		AdminPassword:     opts.SecurityConfig.AdminPassword,
	})
	if err != nil {
		return nil, err
	}

	return &handler{
		filterMap: map[string]accessFilter{
			"repository": &repositoryFilter{
				parser:         &basicParser{},
				registryClient: registryClient,
				adminUsername:  opts.SecurityConfig.AdminUsername,
			},
			"registry": &registryFilter{},
		},
		privateKey:    pk,
		expiredHours:  *opts.SecurityConfig.TokenExpiredHours,
		domainSuffix:  opts.DomainSuffix,
		defaultTenant: opts.DefaultTenant,
		authenticator: at,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	scopes := parseScopes(req.URL)
	log.Debug("Received docker registry authentication request", log.Strings("scopes", scopes))

	var username, userTenantID string
	user, authenticated := authenticationutil.RequestUser(req, h.authenticator)
	if user != nil {
		username = user.GetName()
		userTenantID = user.TenantID()
	}
	requestTenantID := utilregistryrequest.TenantID(req, h.domainSuffix, h.defaultTenant)

	if len(scopes) == 0 {
		// for docker login
		if !authenticated {
			log.Warn("User docker is not authenticated when logging in")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userTenantID != "" && userTenantID != requestTenantID {
			log.Warn("The credential tenant used by the user docker login is inconsistent with the request tenant",
				log.String("userTenantID", userTenantID),
				log.String("requestTenantID", requestTenantID))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	access := h.resourceActions(scopes, requestTenantID)
	u := &userRequest{
		username:      username,
		userTenantID:  userTenantID,
		authenticated: authenticated,
	}
	if err := filterAccess(req.Context(), access, h.filterMap, u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwtToken, err := MakeToken(username, access, h.expiredHours, h.privateKey)
	if err != nil {
		log.Error("Failed create token for docker registry authentication",
			log.String("username", username),
			log.String("tenantID", userTenantID),
			log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(jwtToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bs)
}

func (h *handler) resourceActions(scopes []string, requestTenantID string) []*token.ResourceActions {
	log.Debugf("Scopes: %+v in docker registry v2 request", scopes)
	var res []*token.ResourceActions
	for _, s := range scopes {
		if s == "" {
			continue
		}
		items := strings.Split(s, ":")
		length := len(items)

		tp := ""
		name := ""
		var actions []string

		if length == 1 {
			tp = items[0]
		} else if length == 2 {
			tp = items[0]
			name = items[1]
		} else {
			tp = items[0]
			name = strings.Join(items[1:length-1], ":")
			if len(items[length-1]) > 0 {
				actions = strings.Split(items[length-1], ",")
			}
		}

		if !strings.HasPrefix(name, fmt.Sprintf("%s/", tenant.CrossTenantNamespace)) && requestTenantID != "" {
			name = fmt.Sprintf("%s-%s", requestTenantID, name)
		}

		res = append(res, &token.ResourceActions{
			Type:    tp,
			Name:    name,
			Actions: actions,
		})
	}
	return res
}

// filterAccess iterate a list of resource actions and try to use the filter that matches the resource type to filter the actions.
func filterAccess(ctx context.Context, access []*token.ResourceActions, filters map[string]accessFilter, u *userRequest) error {
	var err error
	for _, a := range access {
		f, ok := filters[a.Type]
		if !ok {
			a.Actions = []string{}
			log.Warnf("No filter found for access type: %s, skip filter, the access of resource '%s' will be set empty.", a.Type, a.Name)
			continue
		}
		err = f.filter(ctx, a, u)
		if err != nil {
			return err
		}
	}
	return nil
}
