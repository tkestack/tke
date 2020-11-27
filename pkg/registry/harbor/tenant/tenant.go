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
	"context"
	"fmt"
	"github.com/docker/distribution/registry/api/v2"
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	"tkestack.io/tke/pkg/registry/harbor/handler"
	utilregistryrequest "tkestack.io/tke/pkg/registry/util/request"
)

const CrossTenantNamespace = "library"

// WithTenant adds an interceptor to the original http request handle and
// converts the request for docker distribution to multi-tenant mode.
func WithTenant(handler http.Handler, registryPrefix, authPrefix, chartPrefix, chartMeseumPrefix, chartAPIPrefix, domainSuffix, defaultTenant string) http.Handler {
	router := v2.Router()
	return &tenant{handler, router, registryPrefix, authPrefix, chartPrefix, chartMeseumPrefix, chartAPIPrefix, domainSuffix, defaultTenant}
}

type tenant struct {
	handler           http.Handler
	router            *mux.Router
	registryPrefix    string
	authPrefix        string
	chartPrefix       string
	chartMeseumPrefix string
	chartAPIPrefix    string
	domainSuffix      string
	defaultTenant     string
}

func (t *tenant) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	tenant := utilregistryrequest.TenantID(r, t.domainSuffix, t.defaultTenant)
	ctx = context.WithValue(ctx, handler.HarborContextKey("tenantID"), tenant)
	r = r.WithContext(ctx)
	// convert tke chartmeseum to harbor chart prefix
	if strings.HasPrefix(r.URL.Path, fmt.Sprintf("%s%s", t.chartMeseumPrefix, "api/")) {
		r.URL.Path = strings.Replace(r.URL.Path, t.chartMeseumPrefix, "/", 1)
	} else if strings.HasPrefix(r.URL.Path, t.chartMeseumPrefix) {
		r.URL.Path = strings.Replace(r.URL.Path, t.chartMeseumPrefix, t.chartPrefix, 1)
	}
	originalPath := r.URL.Path
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
	} else if strings.HasPrefix(originalPath, t.chartAPIPrefix) {
		r.URL.Path = strings.Replace(originalPath, t.chartAPIPrefix, fmt.Sprintf("%s%s-chart-", t.chartAPIPrefix, tenant), 1)
	} else if strings.HasPrefix(originalPath, t.authPrefix) && !strings.Contains(r.URL.RawQuery, "repository%3A"+CrossTenantNamespace) {
		r.URL.RawQuery = strings.Replace(r.URL.RawQuery, "repository%3A", "repository%3A"+tenant+"-image-", -1)
	}
	t.handler.ServeHTTP(w, r)
}
