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

package gateway

import (
	"context"
	"net/http"
	"net/url"
	"tkestack.io/tke/pkg/gateway/auth"
	"tkestack.io/tke/pkg/gateway/token"

	gooidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/pkg/util/log"
)

// CallbackPath is the callback URL path for OAuth2 authorization
const CallbackPath = "/callback"

func registerCallbackRoute(m *mux.PathRecorderMux, oauthConfig *oauth2.Config, oidcHTTPClient *http.Client, disableOIDCProxy bool) {
	m.HandleFunc(CallbackPath, func(writer http.ResponseWriter, request *http.Request) {
		state := request.URL.Query().Get("state")
		if state == "" {
			log.Error("Failed to get state from oidc callback")
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		code := request.URL.Query().Get("code")
		if code == "" {
			log.Error("Failed to get code from oidc callback")
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		ctx := gooidc.ClientContext(context.Background(), oidcHTTPClient)
		opts := authCodeOptions(request, disableOIDCProxy)
		oauth2Token, err := oauthConfig.Exchange(ctx, code, opts...)
		if err != nil {
			log.Error("Failed to exchange oauth2 token by given code", log.String("code", code), log.Err(err))
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		if err := token.ResponseToken(oauth2Token, writer); err != nil {
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, state, http.StatusFound)
	})
}

func authCodeOptions(req *http.Request, disableOIDCProxy bool) []oauth2.AuthCodeOption {
	var opts []oauth2.AuthCodeOption
	if !disableOIDCProxy {
		u := &url.URL{
			Path: CallbackPath,
			Host: req.Host,
		}
		if req.TLS == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
		opts = append(opts, oauth2.SetAuthURLParam(auth.RedirectURIKey, u.String()))
	}
	return opts
}
