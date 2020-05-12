/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package filter

import (
	"context"
	"net/http"

	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
)

// HeaderTenantID is the header name of tenantID.
const HeaderTenantID = "X-TKE-TenantID"

const tenantIDContextKey = "tenantID"

// WithTenantID adds the unique tenantID to the context of the http access
// chain.
func WithTenantID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get(HeaderTenantID)
		// add request id to context
		req = req.WithContext(genericrequest.WithValue(req.Context(), tenantIDContextKey, requestID))
		w.Header().Set(HeaderTenantID, requestID)
		handler.ServeHTTP(w, req)
	})
}

// TenantIDFrom get the tenant id from request context.
func TenantIDFrom(ctx context.Context) string {
	b, ok := ctx.Value(tenantIDContextKey).(string)
	if !ok {
		return ""
	}
	return b
}
