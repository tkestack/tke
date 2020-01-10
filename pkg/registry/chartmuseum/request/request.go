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

package request

import (
	"net/http"
	"tkestack.io/tke/pkg/apiserver/filter"
)

const headerChartRequestID = "X-Request-Id"

// WithRequestID add an interceptor to http request chain for convert requestID
// of apiserver to chartmuseum requestID.
func WithRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := filter.RequestIDFrom(req.Context())
		if requestID != "" {
			req.Header.Add(headerChartRequestID, requestID)
		}
		handler.ServeHTTP(w, req)
	})
}
