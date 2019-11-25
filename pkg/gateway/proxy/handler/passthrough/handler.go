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

package passthrough

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

type handler struct {
	reverseProxy *httputil.ReverseProxy
	protected    bool
}

// NewHandler to create a reverse proxy handler and returns it.
// The reverse proxy will parse the requested cookie content, get the token in
// it, and append it as the http request header to the backend service component.
func NewHandler(address string, cfg *gatewayconfig.PassthroughComponent, protected bool) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	tr, err := transport.NewOneWayTLSTransport(cfg.CAFile, true)
	if err != nil {
		log.Error("Failed to create one-way HTTPS transport", log.String("caFile", cfg.CAFile), log.Err(err))
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: u.Scheme, Host: u.Host})
	reverseProxy.Transport = tr
	reverseProxy.ErrorLog = log.StdErrLogger()
	return &handler{reverseProxy, protected}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.protected {
		// read cookie
		t, err := token.RetrieveToken(req)
		if err != nil {
			log.Error("Failed to retrieve token from client", log.Err(err))
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), w)
			return
		}
		log.Debug("Reverse proxy to protected backend component", log.String("url", req.URL.Path))
		newReq := req.WithContext(context.Background())
		newReq.Header = netutil.CloneHeader(req.Header)
		newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(t.ID)))
		h.reverseProxy.ServeHTTP(w, newReq)
		return
	}
	log.Debug("Reverse proxy to unprotected backend component", log.String("url", req.URL.Path))
	h.reverseProxy.ServeHTTP(w, req)
}
