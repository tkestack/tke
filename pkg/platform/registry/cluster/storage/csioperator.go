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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformv1 "tkestack.io/tke/api/platform/v1"

	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util"
)

// CSIREST implements proxy CSI request to cluster of user.
type CSIREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *CSIREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *CSIREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.CSIProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *CSIREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.CSIProxyOptions)

	clusterv1 := &platformv1.Cluster{}
	err = platformv1.Convert_platform_Cluster_To_v1_Cluster(cluster, clusterv1, nil)
	if err != nil {
		return nil, err
	}

	location, transport, token, err := util.APIServerLocationByCluster(ctx, clusterv1)
	if err != nil {
		return nil, err
	}
	return &csiProxyHandler{
		location:  location,
		transport: transport,
		token:     token,
		namespace: proxyOpts.Namespace,
		name:      proxyOpts.Name,
	}, nil
}

// New creates a new CSI proxy options object
func (r *CSIREST) New() runtime.Object {
	return &platform.CSIProxyOptions{}
}

type csiProxyHandler struct {
	transport http.RoundTripper
	location  *url.URL
	token     string
	namespace string
	name      string
}

func (h *csiProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	prefix := "/apis/storage.tkestack.io/v1"

	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/csis", prefix))
	} else if len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/csis", prefix, h.namespace))
	} else {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/csis/%s", prefix, h.namespace, h.name))
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
