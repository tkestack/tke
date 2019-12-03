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
	"k8s.io/apiserver/pkg/authentication/authenticator"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"net/http"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/apikey"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	authenticationutil "tkestack.io/tke/pkg/registry/util/authentication"
)

type Options struct {
	SecurityConfig  *registryconfig.Security
	ExternalScheme  string
	OIDCIssuerURL   string
	OIDCCAFile      string
	TokenReviewPath string
}

// WithAuthentication creates an http handler that tries to authenticate the
// given chartmuseum request as a user, and then stores any such user found onto
// the provided context for the request.
func WithAuthentication(handler http.Handler, opts *Options) (http.Handler, error) {
	at, err := apikey.NewAPIKeyAuthenticator(&apikey.Options{
		OIDCIssuerURL:   opts.OIDCIssuerURL,
		OIDCCAFile:      opts.OIDCCAFile,
		TokenReviewPath: opts.TokenReviewPath,
		AdminUsername:   opts.SecurityConfig.AdminUsername,
		AdminPassword:   opts.SecurityConfig.AdminPassword,
	})
	if err != nil {
		return nil, err
	}
	return &authentication{
		handler:        handler,
		externalScheme: opts.ExternalScheme,
		authenticator:  at,
	}, nil
}

type authentication struct {
	handler        http.Handler
	externalScheme string
	authenticator  authenticator.Password
}

func (a *authentication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user, authenticated := authenticationutil.RequestUser(req, a.authenticator)
	if authenticated && user != nil {
		req = req.WithContext(genericapirequest.WithUser(req.Context(), user.Info))
	}
	a.handler.ServeHTTP(w, req)
}

// func (a *authentication) notAuthenticated(w http.ResponseWriter, req *http.Request) {
// 	realm := fmt.Sprintf("%s://%s", a.externalScheme, req.Host)
// 	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
// 	err := &model.ErrorResponse{Error: "unauthorized"}
// 	responsewriters.WriteRawJSON(http.StatusUnauthorized, err, w)
// }
