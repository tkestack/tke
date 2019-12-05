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
	"strings"
	"time"
	"tkestack.io/tke/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"

	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util"
)

const (
	cronHPAGroupName = "extensions.tkestack.io"
)

// CronHPAREST implements proxy CronHPA request to cluster of user.
type CronHPAREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *CronHPAREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *CronHPAREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.CronHPAProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *CronHPAREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.CronHPAProxyOptions)

	location, transport, token, err := util.APIServerLocationByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}
	return &cronHPAProxyHandler{
		location:  location,
		transport: transport,
		token:     token,
		namespace: proxyOpts.Namespace,
		name:      proxyOpts.Name,
	}, nil
}

// New creates a new CronHPA proxy options object
func (r *CronHPAREST) New() runtime.Object {
	return &platform.CronHPAProxyOptions{}
}

type cronHPAProxyHandler struct {
	transport http.RoundTripper
	location  *url.URL
	token     string
	namespace string
	name      string
}

func (h *cronHPAProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	prefix := fmt.Sprintf("/apis/%s/v1", cronHPAGroupName)

	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/cronhpas", prefix)
	} else if len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/cronhpas", prefix, h.namespace)
	} else {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/cronhpas/%s", prefix, h.namespace, h.name)
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
