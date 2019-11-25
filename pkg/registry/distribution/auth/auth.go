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
	"bytes"
	"fmt"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/libtrust"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"net/http"
	"net/url"
	"strings"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/distribution/tenant"
	utilregistryrequest "tkestack.io/tke/pkg/registry/util/request"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

const Path = "/registry/auth"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Options struct {
	SecurityConfig  *registryconfig.Security
	OIDCIssuerURL   string
	OIDCCAFile      string
	TokenReviewPath string
	DomainSuffix    string
	DefaultTenant   string
	LoopbackConfig  *restclient.Config
}

type handler struct {
	filterMap            map[string]accessFilter
	privateKey           libtrust.PrivateKey
	expiredHours         int64
	tokenReviewURL       string
	tokenReviewTransport *http.Transport
	domainSuffix         string
	defaultTenant        string
	adminUsername        string
	adminPassword        string
}

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

	issuerURL, err := url.Parse(opts.OIDCIssuerURL)
	if err != nil {
		return nil, err
	}
	issuerURL.Path = opts.TokenReviewPath

	tr, err := transport.NewOneWayTLSTransport(opts.OIDCCAFile, true)
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
		privateKey:           pk,
		expiredHours:         *opts.SecurityConfig.TokenExpiredHours,
		tokenReviewURL:       issuerURL.String(),
		tokenReviewTransport: tr,
		domainSuffix:         opts.DomainSuffix,
		defaultTenant:        opts.DefaultTenant,
		adminPassword:        opts.SecurityConfig.AdminPassword,
		adminUsername:        opts.SecurityConfig.AdminUsername,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	scopes := parseScopes(req.URL)
	log.Debug("Received docker registry authentication request", log.Strings("scopes", scopes))

	username, userTenantID, authenticated := h.userFromRequest(req)
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
	if err := filterAccess(access, h.filterMap, u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwtToken, err := makeToken(username, access, h.expiredHours, h.privateKey)
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

func (h *handler) userFromRequest(r *http.Request) (user, tenantID string, authenticated bool) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Warn("Docker login did not carry the http basic authentication credentials when logging in")
		return
	}
	if h.adminPassword != "" &&
		h.adminUsername != "" &&
		username == h.adminUsername &&
		password == h.adminPassword {
		log.Debug("Docker login as system administrator")
		user = h.adminUsername
		tenantID = ""
		authenticated = true
		return
	}
	if h.tokenReviewURL == "" {
		log.Warn("Token review url not specify, failed to review token")
		return
	}
	log.Debug("Start review token", log.String("username", username), log.String("tokenReviewURL", h.tokenReviewURL))
	tokenReviewRequest := &v1.TokenReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TokenReview",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "tke",
		},
		Spec: v1.TokenReviewSpec{
			Token: password,
		},
	}
	bs, err := json.Marshal(tokenReviewRequest)
	if err != nil {
		log.Error("Failed to marshal token review request", log.Any("tokenReview", tokenReviewRequest), log.Err(err))
		return
	}

	req, err := http.NewRequest(http.MethodPost, h.tokenReviewURL, bytes.NewBuffer(bs))
	if err != nil {
		log.Error("Failed to create token review request", log.Err(err))
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Transport: h.tokenReviewTransport,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to request token review", log.Err(err))
	}

	defer func() {
		_ = res.Body.Close()
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read response body", log.Err(err))
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Error("Authentication server returns http status code error when token review",
			log.Int("statusCode", res.StatusCode),
			log.ByteString("responseBody", body))
		return
	}
	var tokenReviewResponse v1.TokenReview
	if err := json.Unmarshal(body, &tokenReviewResponse); err != nil {
		log.Error("Failed to unmarshal token review response", log.ByteString("body", body), log.Err(err))
		return
	}
	if tokenReviewResponse.Status.Error != "" {
		log.Error("Authentication server returns error when token review",
			log.String("error", tokenReviewResponse.Status.Error))
		return
	}
	if !tokenReviewResponse.Status.Authenticated {
		log.Error("Authentication server returns not authenticated when token review",
			log.Any("tokenReviewStatus", tokenReviewResponse.Status))
		return
	}

	user = tokenReviewResponse.Status.User.Username
	authenticated = true
	if len(tokenReviewResponse.Status.User.Extra) > 0 {
		if t, ok := tokenReviewResponse.Status.User.Extra[genericoidc.TenantIDKey]; ok {
			if len(t) > 0 {
				tenantID = t[0]
			}
		}
	}
	log.Debug("Docker login verifies that the user identity is successful",
		log.String("username", user),
		log.String("tenantID", tenantID))
	return
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
func filterAccess(access []*token.ResourceActions, filters map[string]accessFilter, u *userRequest) error {
	var err error
	for _, a := range access {
		f, ok := filters[a.Type]
		if !ok {
			a.Actions = []string{}
			log.Warnf("No filter found for access type: %s, skip filter, the access of resource '%s' will be set empty.", a.Type, a.Name)
			continue
		}
		err = f.filter(a, u)
		if err != nil {
			return err
		}
	}
	return nil
}
