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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	clientrest "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/platform/proxy"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

// ProxyREST implements proxy native api request to cluster of user.
type ProxyREST struct {
	rest.Storage
	store *registry.Store
	host  string

	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *ProxyREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *ProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.ProxyOptions{}, false, "path"
}

// Connect returns a handler for the native api proxy
func (r *ProxyREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, _ rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.ProxyOptions)

	if proxyOpts.Path == "" {
		return nil, errors.NewBadRequest("invalid path")
	}

	if strings.HasPrefix(proxyOpts.Path, "/apis/platform.tkestack.io/") {
		return nil, errors.NewBadRequest("cycle dispatch")
	}

	config, err := proxy.GetConfig(ctx, r.platformClient)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	userName, tenantID := authentication.UsernameAndTenantID(ctx)
	uri, err := makeURL(config.Host, proxyOpts.Path)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	transport, err := clientrest.TransportFor(config)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director:  makeDirector(cluster.ObjectMeta.Name, userName, tenantID, uri, config.BearerToken),
		Transport: transport,
	}, nil
}

// New creates a new helm proxy options object
func (r *ProxyREST) New() runtime.Object {
	return &platform.ProxyOptions{}
}

func makeDirector(clusterName, userName, tenantID string, uri *url.URL, token string) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Set(filter.ClusterNameHeaderKey, clusterName)
		if token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
		req.URL = uri
	}
}

//proxyPath have been decoded somewhere before passing to makeURL
func makeURL(host, proxyPath string) (*url.URL, error) {
	u, err := url.Parse(host) //will returen error if a host not contains a schema
	if err != nil {
		log.Errorf("parse host error %s\n", err)
		return nil, err
	}

	/* a host without a path will have a emplty u.Path, and a proxyPath may not start with "/"
	In order to make the newPath begin with only one "/", add a "/" to empty u.Path
	*/
	if u.Path == "" {
		u.Path = "/"
	}

	newPath := path.Join(u.Path, proxyPath) // ensure newPath begin with "/"

	newURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, newPath)
	u, err = url.Parse(newURL)
	if err != nil {
		log.Errorf("parse new url error %s\n", err)
		return nil, err
	}
	return u, nil
}
