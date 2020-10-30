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

package authorization

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	restclient "k8s.io/client-go/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	authorizationutil "tkestack.io/tke/pkg/registry/util/authorization"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Options struct {
	AdminUsername  string
	ExternalScheme string
	LoopbackConfig *restclient.Config
	Authorizer     authorizer.Authorizer
}

// WithAuthorization creates an http handler that tries to authorized requests
// on to handler, and returns a forbidden error otherwise.
func WithAuthorization(handler http.Handler, opts *Options) (http.Handler, error) {
	registryClient, err := registryinternalclient.NewForConfig(opts.LoopbackConfig)
	if err != nil {
		return nil, err
	}
	if opts.Authorizer == nil {
		return nil, fmt.Errorf("chartmuseum authorizer is nil")
	}

	authorizationHandler := &authorization{
		registryClient: registryClient,
		nextHandler:    handler,
		adminUsername:  opts.AdminUsername,
		externalScheme: opts.ExternalScheme,
		authorizer:     opts.Authorizer,
	}
	router := mux.NewRouter()
	router.HandleFunc("/chart/{tenantID}/{chartGroup}/index.yaml", authorizationHandler.index).Methods(http.MethodGet)
	router.HandleFunc("/chart/{tenantID}/{chartGroup}/charts/{file}", authorizationHandler.getChart).Methods(http.MethodGet)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/charts", authorizationHandler.apiListChart).Methods(http.MethodGet)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/charts/{name}", authorizationHandler.apiGetChart).Methods(http.MethodGet)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}", authorizationHandler.apiGetChartVersion).Methods(http.MethodGet)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/charts", authorizationHandler.apiCreateChart).Methods(http.MethodPost)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/prov", authorizationHandler.apiCreateProvenance).Methods(http.MethodPost)
	router.HandleFunc("/chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}", authorizationHandler.apiDeleteChartVersion).Methods(http.MethodDelete)

	authorizationHandler.router = router
	return authorizationHandler, nil
}

type authorization struct {
	router         *mux.Router
	registryClient *registryinternalclient.RegistryClient
	nextHandler    http.Handler
	adminUsername  string
	externalScheme string
	authorizer     authorizer.Authorizer
}

func (a *authorization) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func (a *authorization) isAdmin(w http.ResponseWriter, req *http.Request) bool {
	username, tenantID := authentication.UsernameAndTenantID(req.Context())
	if tenantID == "" && username != "" && username == a.adminUsername {
		return true
	}
	return false
}

// AuthorizeForChart check if chart resource is authorized
func AuthorizeForChart(w http.ResponseWriter, req *http.Request, authzer authorizer.Authorizer, verb string, cg registryv1.ChartGroup, chartName string) (passed bool, err error) {
	switch verb {
	case "get":
		{
			if cg.Spec.Visibility == registryv1.VisibilityPublic {
				return true, nil
			}
			break
		}
	}
	u, exist := genericapirequest.UserFrom(req.Context())
	if !exist || u == nil {
		return false, fmt.Errorf("empty user info, not authenticated")
	}
	return authorizationutil.AuthorizeForChart(req.Context(), u, authzer, verb, cg, chartName)
}

// AuthorizeForChartGroup check if chartgroup resource is authorized
func AuthorizeForChartGroup(w http.ResponseWriter, req *http.Request, authzer authorizer.Authorizer, verb string, cg registryv1.ChartGroup) (passed bool, err error) {
	switch verb {
	case "get":
		{
			if cg.Spec.Visibility == registryv1.VisibilityPublic {
				return true, nil
			}
			break
		}
	}
	u, exist := genericapirequest.UserFrom(req.Context())
	if !exist || u == nil {
		return false, fmt.Errorf("empty user info, not authenticated")
	}
	return authorizationutil.AuthorizeForChartGroup(req.Context(), u, authzer, verb, cg)
}
