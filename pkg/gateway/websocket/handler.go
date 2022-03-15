/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
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

package websocket

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
	"k8s.io/klog"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
)

type handler struct {
	config rest.Config
}

// NewHandler to create a reverse proxy handler and returns it.
// The reverse proxy will parse the requested cookie content, get the token in
// it, and append it as the http request header to the backend service component.
func NewHandler(address string) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	cfg := rest.Config{
		Host: u.Host,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	return &handler{
		config: cfg,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("in websocket")
	clusterName := req.URL.Query().Get("clusterName")
	if clusterName == "" {
		log.Error("Failed to get cluster name from request")
		http.Error(w, "Invalid cluster name", http.StatusBadRequest)
		return
	}

	encodePath := req.URL.Query().Get("encodePath")
	decodePathBytes, err := base64.StdEncoding.DecodeString(encodePath)
	if err != nil {
		log.Error("Failed to get/decode encode path from request")
		http.Error(w, fmt.Sprintf("Invalid decode path: %v", err), http.StatusBadRequest)
		return
	}
	decodePath := string(decodePathBytes)

	t, err := token.RetrieveToken(req)
	if err != nil {
		log.Error("Failed to retrieve token from webtty", log.Err(err))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	cfg := h.config
	cfg.BearerToken = t.ID
	target, err := url.Parse(decodePath)
	if err != nil {
		log.Error("Failed to generate decode path url", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	target.Scheme = "https"
	target.Host = cfg.Host

	transport, err := rest.TransportFor(&cfg)
	if err != nil {
		log.Error("Failed initialize transport", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	upgradeTransport, err := makeUpgradeTransport(&cfg, 30*time.Second)
	if err != nil {
		log.Error("Failed initialize upgrade transport", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	responder := &responder{}
	proxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, responder)
	proxy.UpgradeTransport = upgradeTransport
	proxy.Location = target

	reqClone := utilnet.CloneRequest(req)
	reqClone.URL = target
	reqClone.Header.Add(filter.ClusterNameHeaderKey, clusterName)

	proxy.ServeHTTP(w, reqClone)
}

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.Errorf("Error while proxying request: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func makeUpgradeTransport(config *rest.Config, keepalive time.Duration) (proxy.UpgradeRequestRoundTripper, error) {
	transportConfig, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	tlsConfig, err := transport.TLSConfigFor(transportConfig)
	if err != nil {
		return nil, err
	}
	rt := utilnet.SetOldTransportDefaults(&http.Transport{
		TLSClientConfig: tlsConfig,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: keepalive,
		}).DialContext,
	})

	upgrader, err := transport.HTTPWrappersForConfig(transportConfig, proxy.MirrorRequest)
	if err != nil {
		return nil, err
	}
	return proxy.NewUpgradeRequestRoundTripper(rt, upgrader), nil
}
