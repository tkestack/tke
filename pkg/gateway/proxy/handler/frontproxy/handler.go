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

package frontproxy

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

type handler struct {
	oidcAuthenticator *oidc.Authenticator
	reverseProxy      *httputil.ReverseProxy
	usernameHeader    string
	groupsHeader      string
	extraPrefixHeader string
	protected         bool
}

// NewHandler to create a reverse proxy handler and returns it.
// This mode of reverse proxy will resolve the token in the request cookie,
// call oidc to get the user identity, and pass it to the backend service as
// HTTP header.
func NewHandler(address string, cfg *gatewayconfig.FrontProxyComponent, oidcAuthenticator *oidc.Authenticator, protected bool) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	tr, err := transport.NewTwoWayTLSTransport(cfg.CAFile, cfg.ClientCertFile, cfg.ClientKeyFile)
	if err != nil {
		log.Error("Failed to create two-way HTTPS transport",
			log.String("caFile", cfg.CAFile),
			log.String("clientCertFile", cfg.ClientCertFile),
			log.String("clientKeyFile", cfg.ClientKeyFile),
			log.Err(err))
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: u.Scheme, Host: u.Host})
	reverseProxy.Transport = tr
	reverseProxy.ErrorLog = log.StdErrLogger()
	return &handler{
		oidcAuthenticator: oidcAuthenticator,
		reverseProxy:      reverseProxy,
		usernameHeader:    cfg.UsernameHeader,
		groupsHeader:      cfg.GroupsHeader,
		extraPrefixHeader: cfg.ExtraPrefixHeader,
		protected:         protected,
	}, nil
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.protected {
		// read cookie
		t, err := token.RetrieveToken(req)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), w)
			return
		}
		r, authenticated, err := h.oidcAuthenticator.AuthenticateToken(req.Context(), t.ID)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
			return
		}
		if !authenticated {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized("invalid token"), w)
			return
		}

		header := make(http.Header, len(req.Header))
		for key, values := range req.Header {
			if strings.ToLower(key) != "cookie" {
				newValues := make([]string, len(values))
				copy(newValues, values)
				header[key] = newValues
			}
		}
		username := r.User.GetName()
		if username != "" {
			header.Set(h.usernameHeader, username)
		}
		groups := r.User.GetGroups()
		if len(groups) > 0 {
			header.Set(h.groupsHeader, strings.Join(groups, ","))
		}
		extra := r.User.GetExtra()
		if len(extra) > 0 {
			for k, v := range extra {
				header.Set(fmt.Sprintf("%s%s", h.extraPrefixHeader, k), strings.Join(v, ","))
			}
		}

		newReq := req.WithContext(context.Background())
		newReq.Header = header

		h.reverseProxy.ServeHTTP(w, newReq)
		return
	}
	h.reverseProxy.ServeHTTP(w, req)
}
