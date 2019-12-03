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
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	restclient "k8s.io/client-go/rest"
	"net/http"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Options struct {
	AdminUsername  string
	ExternalScheme string
	LoopbackConfig *restclient.Config
}

// WithAuthorization creates an http handler that tries to authorized requests
// on to handler, and returns a forbidden error otherwise.
func WithAuthorization(handler http.Handler, opts *Options) (http.Handler, error) {
	registryClient, err := registryinternalclient.NewForConfig(opts.LoopbackConfig)
	if err != nil {
		return nil, err
	}

	authorizationHandler := &authorization{
		registryClient: registryClient,
		nextHandler:    handler,
		adminUsername:  opts.AdminUsername,
		externalScheme: opts.ExternalScheme,
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
}

func (a *authorization) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	username, tenantID := authentication.GetUsernameAndTenantID(req.Context())
	if tenantID == "" && username != "" && username == a.adminUsername {
		a.nextHandler.ServeHTTP(w, req)
		return
	}
	a.router.ServeHTTP(w, req)
}

// getChart serve http get request on /chart/{tenantID}/{chartGroup}/charts/{file}
func (a *authorization) getChart(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}

// apiListChart serve http get request on /chart/api/{tenantID}/{chartGroup}/charts
func (a *authorization) apiListChart(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}

// apiGetChart serve http get request on /chart/api/{tenantID}/{chartGroup}/charts/{name}
func (a *authorization) apiGetChart(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}

// apiGetChartVersion serve http get request on /chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}
func (a *authorization) apiGetChartVersion(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}

// apiCreateProvenance serve http post request on /chart/api/{tenantID}/{chartGroup}/prov
func (a *authorization) apiCreateProvenance(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}

// apiDeleteChartVersion serve http delete request on /chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}
func (a *authorization) apiDeleteChartVersion(w http.ResponseWriter, req *http.Request) {
	// todo: authorization
	a.nextHandler.ServeHTTP(w, req)
}
