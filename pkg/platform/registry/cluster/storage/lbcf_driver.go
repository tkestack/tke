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

package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"

	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"

	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util"
)

const (
	// LBCFDriverProxyAction sends request to driver via kube-apiserver
	LBCFDriverProxyAction LBCFAction = "driverProxy"
)

// LBCFDriverREST implements proxy LBCF LoadBalancerDriver request to cluster of user.
type LBCFDriverREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *LBCFDriverREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *LBCFDriverREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.LBCFProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *LBCFDriverREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.LBCFProxyOptions)

	location, transport, token, err := util.APIServerLocationByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}
	return &lbcfDriverProxyHandler{
		location:  location,
		transport: transport,
		token:     token,
		namespace: proxyOpts.Namespace,
		name:      proxyOpts.Name,
		action:    proxyOpts.Action,
	}, nil
}

// New creates a new LBCF driver proxy options object
func (r *LBCFDriverREST) New() runtime.Object {
	return &platform.LBCFProxyOptions{}
}

type lbcfDriverProxyHandler struct {
	transport http.RoundTripper
	location  *url.URL
	token     string
	namespace string
	name      string
	action    string
}

func (h *lbcfDriverProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	// todo: Change the apigroup here once the integration pipeline configuration is complete using the tapp in the tkestack group
	prefix := "/apis/lbcf.tkestack.io/v1beta1"

	if len(h.action) > 0 {
		h.serveAction(w, req)
		return
	}
	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/loadbalancerdrivers", prefix))
	} else if len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/loadbalancerdrivers", prefix, h.namespace))
	} else {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/loadbalancerdrivers/%s", prefix, h.namespace, h.name))
	}

	// WithContext creates a shallow clone of the request with the new context.
	newReq := req.WithContext(context.Background())
	newReq.Header = netutil.CloneHeader(req.Header)
	newReq.URL = &loc
	if h.token != "" {
		newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(h.token)))
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: h.location.Scheme, Host: h.location.Host})
	reverseProxy.Transport = h.transport
	reverseProxy.FlushInterval = 100 * time.Millisecond
	reverseProxy.ErrorLog = log.StdErrLogger()
	reverseProxy.ServeHTTP(w, newReq)
}

func (h *lbcfDriverProxyHandler) serveAction(w http.ResponseWriter, req *http.Request) {
	if len(h.namespace) == 0 || len(h.name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("namespace and name must be specified"), w)
		return
	}

	switch h.action {
	case string(LBCFDriverProxyAction):
		loc := *h.location
		loc.RawQuery = req.URL.RawQuery
		queries := req.URL.Query()
		apiPath, apiPort := getDriverAPI(queries)
		loc.Path = path.Join(loc.Path, fmt.Sprintf("/api/v1/namespaces/%s/services/%s:%s/proxy/%s",
			h.namespace, h.name, apiPort, apiPath))

		// delete lbcf query paramaters before send to kube-apiserver
		queries.Del("namespace")
		queries.Del("name")
		queries.Del("action")
		queries.Del("api")
		queries.Del("apiPort")
		loc.RawQuery = queries.Encode()

		newReq := req.WithContext(context.Background())
		newReq.Header = netutil.CloneHeader(req.Header)
		newReq.URL = &loc
		if h.token != "" {
			newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(h.token)))
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: h.location.Scheme, Host: h.location.Host})
		reverseProxy.Transport = h.transport
		reverseProxy.FlushInterval = 100 * time.Millisecond
		reverseProxy.ErrorLog = log.StdErrLogger()
		reverseProxy.ServeHTTP(w, newReq)
	}
}

func getDriverAPI(queries url.Values) (apiPath, apiPort string) {
	apiPath = "/"
	driverAPI := queries["api"]
	if len(driverAPI) > 0 {
		apiPath = driverAPI[0]
	}

	apiPort = "443"
	p := queries["apiPort"]
	if len(p) > 0 {
		apiPort = p[0]
	}
	return
}
