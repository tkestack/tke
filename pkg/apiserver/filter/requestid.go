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

package filter

import (
	"context"
	"net/http"

	"github.com/segmentio/ksuid"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/util/log"
)

const headerRequestID = "X-Remote-Extra-RequestID"
const requestIDContextKey = "requestID"

// WithRequestID adds the unique requestID to the context of the http access
// chain.
func WithRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get(headerRequestID)
		if requestID == "" {
			requestID = ksuid.New().String()
		}
		log.Debug("Received http request", log.String("requestID", requestID))
		// add request id to context
		req = req.WithContext(genericrequest.WithValue(req.Context(), requestIDContextKey, requestID))
		w.Header().Set(headerRequestID, requestID)
		handler.ServeHTTP(w, req)
	})
}

// RequestIDFrom get the request id from request context.
func RequestIDFrom(ctx context.Context) string {
	b, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}
	return b
}
