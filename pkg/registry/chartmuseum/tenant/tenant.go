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
	"net/http"
	"strings"

	utilregistryrequest "tkestack.io/tke/pkg/registry/util/request"
)

// WithTenant adds an interceptor to the original http request handle and
// converts the request for chartmuseum to multi-tenant mode.
func WithTenant(handler http.Handler, pathPrefix, domainSuffix, defaultTenant string) http.Handler {
	return &tenant{handler, pathPrefix, domainSuffix, defaultTenant}
}

type tenant struct {
	handler       http.Handler
	pathPrefix    string
	domainSuffix  string
	defaultTenant string
}

// ServeHTTP responds to an HTTP request.
func (t *tenant) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, t.pathPrefix) {
		tenant := utilregistryrequest.TenantID(r, t.domainSuffix, t.defaultTenant)
		if strings.HasPrefix(path, fmt.Sprintf("%sapi/", t.pathPrefix)) {
			r.URL.Path = strings.Replace(path, fmt.Sprintf("%sapi/", t.pathPrefix), fmt.Sprintf("%sapi/%s/", t.pathPrefix, tenant), 1)
		} else {
			r.URL.Path = strings.Replace(path, t.pathPrefix, fmt.Sprintf("%s%s/", t.pathPrefix, tenant), 1)
		}
	}
	t.handler.ServeHTTP(w, r)
}
