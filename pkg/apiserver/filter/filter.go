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
	"github.com/segmentio/ksuid"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
	"net"
	"net/http"
	"tkestack.io/tke/pkg/util/log"
)

const localRequestContextKey = "localRequest"
const localRequestContextValue = "true"

// WithRequestID adds the unique requestID to the context of the http access
// chain.
func WithRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Remote-Extra-RequestID")
		if requestID == "" {
			requestID = ksuid.New().String()
		}
		log.Debug("Received http request", log.String("requestID", requestID))
		w.Header().Set("X-Remote-Extra-RequestID", requestID)
		handler.ServeHTTP(w, req)
	})
}

// WithLocal adds the local identify to the context of the http access chain.
func WithLocal(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			userIP := net.ParseIP(ip)
			if userIP.IsLoopback() {
				req = req.WithContext(genericrequest.WithValue(req.Context(), localRequestContextKey, localRequestContextValue))
			}
		}
		handler.ServeHTTP(w, req)
	})
}

// LocalFrom get the local identity from request context.
func LocalFrom(ctx context.Context) bool {
	b, ok := ctx.Value(localRequestContextKey).(string)
	if !ok {
		return false
	}
	return b == localRequestContextValue
}
