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
	"k8s.io/apiserver/pkg/server/mux"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/util/log"
)

// registerAuthRoute to register route for proxy TKE auth server.
func registerAuthRoute(m *mux.PathRecorderMux, oidcHTTPClient *http.Client, oidcAuthenticator *oidc.Authenticator) error {
	u, err := url.Parse(oidcAuthenticator.IssuerURL)
	if err != nil {
		log.Error("Failed to parse OIDC issuer url", log.String("issuerURL", oidcAuthenticator.IssuerURL), log.Err(err))
		return err
	}
	pathPrefix := u.Path
	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix = pathPrefix + "/"
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: u.Scheme, Host: u.Host})
	reverseProxy.Transport = oidcHTTPClient.Transport
	reverseProxy.ErrorLog = log.StdErrLogger()
	handler := &proxyHandler{reverseProxy}
	m.HandlePrefix(pathPrefix, handler)
	return nil
}

type proxyHandler struct {
	reverseProxy *httputil.ReverseProxy
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.reverseProxy.ServeHTTP(w, req)
}
