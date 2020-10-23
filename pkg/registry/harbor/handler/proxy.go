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

package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"context"
	"fmt"
	"strings"

	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

type handler struct {
	reverseProxy *httputil.ReverseProxy
	host         string
	externalHost string
}

type harborContextKey string

// NewHandler to create a reverse proxy handler and returns it.
func NewHandler(address string, cafile string, externalHost string) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	tr, err := transport.NewOneWayTLSTransport(cafile, true)
	if err != nil {
		log.Error("Failed to create one-way HTTPS transport", log.String("caFile", cafile), log.Err(err))
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: u.Scheme, Host: u.Host})
	reverseProxy.Transport = tr
	reverseProxy.ModifyResponse = rewriteBody
	reverseProxy.ErrorLog = log.StdErrLogger()
	return &handler{reverseProxy, u.Host, externalHost}, nil
}

func rewriteBody(resp *http.Response) (err error) {

	ctx := resp.Request.Context()
	host := ctx.Value(harborContextKey("host"))
	externalHost := ctx.Value(harborContextKey("exHost"))
	authHeader := resp.Header.Get("www-authenticate")
	if authHeader != "" {
		header := fmt.Sprintf("Bearer realm=\"https://%s/service/token\",service=\"harbor-registry\"", externalHost)
		log.Debug("Modify backend harbor header www-authenticate", log.String("header", header))
		resp.Header.Set("www-authenticate", header)
	}

	locationHeader := resp.Header.Get("location")
	if locationHeader != "" {
		log.Debug("Replace harbor location header", log.String("original host", host.(string)), log.String("tke host", externalHost.(string)))
		resp.Header.Set("location", strings.ReplaceAll(locationHeader, host.(string), externalHost.(string)))
	}

	return nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	log.Debug("Reverse proxy to backend harbor", log.String("url", req.URL.Path))
	req.Host = h.host
	ctx := context.WithValue(req.Context(), harborContextKey("host"), h.host)
	ctx = context.WithValue(ctx, harborContextKey("exHost"), h.externalHost)
	req = req.WithContext(ctx)
	h.reverseProxy.ServeHTTP(w, req)
}
