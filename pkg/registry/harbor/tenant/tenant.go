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

package tenant

import (
	"fmt"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	utilregistryrequest "tkestack.io/tke/pkg/registry/util/request"
)

const CrossTenantNamespace = "library"

// WithTenant adds an interceptor to the original http request handle and
// converts the request for docker distribution to multi-tenant mode.
func WithTenant(handler http.Handler, registryPrefix, authPrefix, chartPrefix, domainSuffix, defaultTenant string) http.Handler {
	router := v2.Router()
	return &tenant{handler, router, registryPrefix, authPrefix, chartPrefix, domainSuffix, defaultTenant}
}

type tenant struct {
	handler        http.Handler
	router         *mux.Router
	registryPrefix string
	authPrefix     string
	chartPrefix    string
	domainSuffix   string
	defaultTenant  string
}

func (t *tenant) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originalPath := r.URL.Path
	tenant := utilregistryrequest.TenantID(r, t.domainSuffix, t.defaultTenant)
	if strings.HasPrefix(originalPath, t.registryPrefix) && !strings.HasPrefix(originalPath, fmt.Sprintf("/v2/%s/", CrossTenantNamespace)) {
		var match mux.RouteMatch
		if matched := t.router.Match(r, &match); matched {
			routeName := match.Route.GetName()
			if routeName == v2.RouteNameManifest ||
				routeName == v2.RouteNameTags ||
				routeName == v2.RouteNameBlob ||
				routeName == v2.RouteNameBlobUpload {
				r.URL.Path = strings.Replace(originalPath, t.registryPrefix, fmt.Sprintf("%s%s-image-", t.registryPrefix, tenant), 1)
			}
		}
	} else if strings.HasPrefix(originalPath, t.chartPrefix) {
		r.URL.Path = strings.Replace(originalPath, t.chartPrefix, fmt.Sprintf("%s%s-chart-", t.chartPrefix, tenant), 1)
	} else if strings.HasPrefix(originalPath, t.authPrefix) {
		r.URL.RawQuery = strings.Replace(r.URL.RawQuery, "repository%3A", "repository%3A"+tenant+"-image-", -1)
	}
	t.handler.ServeHTTP(w, r)
}
