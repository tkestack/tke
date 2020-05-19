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
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	registryv1 "tkestack.io/tke/api/registry/v1"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	authorizationutil "tkestack.io/tke/pkg/registry/util/authorization"
	"tkestack.io/tke/pkg/util/log"
)

// VersionREST adapts a service registry into apiserver's RESTStorage model.
type VersionREST struct {
	store          ChartStorage
	registryClient *registryinternalclient.RegistryClient
	registryConfig *registryconfig.RegistryConfiguration
	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string
	authorizer     authorizer.Authorizer
}

// NewVersionREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various chart.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewVersionREST(
	store ChartStorage,
	registryClient *registryinternalclient.RegistryClient,
	registryConfig *registryconfig.RegistryConfiguration,
	externalScheme string,
	externalHost string,
	externalPort int,
	externalCAFile string,
	authorizer authorizer.Authorizer,
) *VersionREST {
	rest := &VersionREST{
		store:          store,
		registryClient: registryClient,
		registryConfig: registryConfig,
		externalScheme: externalScheme,
		externalHost:   externalHost,
		externalPort:   externalPort,
		externalCAFile: externalCAFile,
		authorizer:     authorizer,
	}
	return rest
}

// New creates a new chart proxy options object
func (r *VersionREST) New() runtime.Object {
	return &registry.ChartProxyOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *VersionREST) ConnectMethods() []string {
	return []string{"DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *VersionREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &registry.ChartProxyOptions{}, false, ""
}

// Connect returns a handler for the chart proxy
func (r *VersionREST) Connect(ctx context.Context, chartName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := r.store.Get(ctx, chartName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	chart := obj.(*registry.Chart)

	cg, err := r.registryClient.ChartGroups().Get(ctx, chart.Namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	proxyOpts := opts.(*registry.ChartProxyOptions)

	if proxyOpts.Version == "" {
		return nil, errors.NewBadRequest("version is required")
	}

	return &versionProxyHandler{
		chart:          chart,
		chartGroup:     cg,
		chartVersion:   proxyOpts.Version,
		externalScheme: r.externalScheme,
		externalHost:   r.externalHost,
		externalPort:   r.externalPort,
		externalCAFile: r.externalCAFile,
		registryConfig: r.registryConfig,
		authorizer:     r.authorizer,
	}, nil
}

type versionProxyHandler struct {
	chart          *registry.Chart
	chartGroup     *registry.ChartGroup
	chartVersion   string
	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string
	registryConfig *registryconfig.RegistryConfiguration
	authorizer     authorizer.Authorizer
}

func (h *versionProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := h.check(w, req)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusUnauthorized, err.Error(), w)
	}

	host := h.externalHost
	if h.externalPort > 0 {
		host = host + ":" + strconv.Itoa(h.externalPort)
	}
	loc := &url.URL{
		Scheme: h.externalScheme,
		Host:   host,
		Path:   fmt.Sprintf("/chart/api/%s/charts/%s/%s", h.chart.Spec.ChartGroupName, h.chart.Spec.Name, h.chartVersion),
	}

	// WithContext creates a shallow clone of the request with the new context.
	newReq := req.WithContext(context.Background())
	newReq.Header = netutil.CloneHeader(req.Header)
	newReq.URL = loc
	newReq.SetBasicAuth(h.registryConfig.Security.AdminUsername, h.registryConfig.Security.AdminPassword)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: loc.Scheme, Host: loc.Host})
	reverseProxy.Transport = transport
	reverseProxy.FlushInterval = 100 * time.Millisecond
	reverseProxy.ErrorLog = log.StdErrLogger()
	reverseProxy.ServeHTTP(w, newReq)
}

func (h *versionProxyHandler) check(w http.ResponseWriter, req *http.Request) error {
	var cg = &registryv1.ChartGroup{}
	err := registryv1.Convert_registry_ChartGroup_To_v1_ChartGroup(h.chartGroup, cg, nil)
	if err != nil {
		return err
	}
	u, exist := genericapirequest.UserFrom(req.Context())
	if !exist || u == nil {
		return fmt.Errorf("empty user info, not authenticated")
	}

	authorized, err := authorizationutil.AuthorizeForChart(req.Context(), u, h.authorizer, "delete", *cg, h.chart.Name)
	if err != nil {
		return err
	}
	if !authorized {
		return fmt.Errorf("not authenticated")
	}
	return nil
}
