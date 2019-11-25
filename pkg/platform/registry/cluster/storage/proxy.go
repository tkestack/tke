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
	"crypto/tls"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/localtrust"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/platform/util"
)

// ProxyREST implements proxy native api request to cluster of user.
type ProxyREST struct {
	rest.Storage
	store *registry.Store
	host  string
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *ProxyREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.HelmProxyOptions{}, false, "path"
}

// Connect returns a handler for the native api proxy
func (r *ProxyREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.HelmProxyOptions)

	if proxyOpts.Path == "" {
		return nil, errors.NewBadRequest("invalid path")
	}

	if strings.HasPrefix(proxyOpts.Path, "/apis/platform.tkestack.io/") {
		return nil, errors.NewBadRequest("cycle dispatch")
	}

	u, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.NewUnauthorized("unknown user")
	}
	token, err := localtrust.GenerateToken(u)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	uri, err := makeURL(r.host, proxyOpts.Path)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	return &httputil.ReverseProxy{
		Director: makeDirector(cluster.ObjectMeta.Name, uri, token),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}, nil
}

// New creates a new helm proxy options object
func (r *ProxyREST) New() runtime.Object {
	return &platform.HelmProxyOptions{}
}

func makeDirector(clusterName string, uri *url.URL, token string) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Set(filter.ClusterNameHeaderKey, clusterName)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.URL = uri
	}
}

func makeURL(host, path string) (*url.URL, error) {
	var port int64
	hostSegment := strings.Split(host, ":")
	if len(hostSegment) == 0 {
		port = 443
	} else {
		var err error
		port, err = strconv.ParseInt(hostSegment[len(hostSegment)-1], 10, 32)
		if err != nil {
			port = 443
		}
	}

	p := strings.TrimPrefix(path, "/")

	return url.Parse(fmt.Sprintf("https://127.0.0.1:%d/%s", port, p))
}
