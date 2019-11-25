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
	"bytes"
	"context"
	"io/ioutil"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
	"net/http"
)

const clusterContextKey = "clusterName"
const requestBodyKey = "requestBody"
const fuzzyResourceContextKey = "fuzzyResourceName"

// ClusterNameHeaderKey is the header name of cluster
const ClusterNameHeaderKey = "X-TKE-ClusterName"
const fuzzyResourceNameHeaderKey = "X-TKE-FuzzyResourceName"

// RequestBody represents the body of HTTP request.
type RequestBody struct {
	Data        []byte
	ContentType string
}

// ClusterFrom get the cluster name from request context.
func ClusterFrom(ctx context.Context) string {
	clusterName, ok := ctx.Value(clusterContextKey).(string)
	if !ok {
		return ""
	}
	return clusterName
}

// FuzzyResourceFrom get the fuzzy resource name from request context.
func FuzzyResourceFrom(ctx context.Context) string {
	fuzzyResourceName, ok := ctx.Value(fuzzyResourceContextKey).(string)
	if !ok {
		return ""
	}
	return fuzzyResourceName
}

// RequestBodyFrom returns the RequestBody object.
func RequestBodyFrom(ctx context.Context) (*RequestBody, bool) {
	val := ctx.Value(requestBodyKey)
	if val == nil {
		return nil, false
	}
	obj, ok := val.(*RequestBody)
	return obj, ok
}

// WithCluster creates an http handler that tries to get the cluster name from
// the given request, and then stores any such cluster name found onto the
// provided context for the request.
func WithCluster(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		clusterName := req.Header.Get(ClusterNameHeaderKey)
		if clusterName != "" {
			req = req.WithContext(genericrequest.WithValue(req.Context(), clusterContextKey, clusterName))
		}
		handler.ServeHTTP(w, req)
	})
}

// WithFuzzyResource adds the fuzzy resource name to the context of the http
// access chain.
func WithFuzzyResource(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fuzzyResourceName := req.Header.Get(fuzzyResourceNameHeaderKey)
		if fuzzyResourceName != "" {
			req = req.WithContext(genericrequest.WithValue(req.Context(), fuzzyResourceContextKey, fuzzyResourceName))
		}
		handler.ServeHTTP(w, req)
	})
}

// WithRequestBody adds the request body to the context of the http access chain.
func WithRequestBody(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		method := req.Method
		if method == http.MethodPut || method == http.MethodPatch || method == http.MethodPost || method == http.MethodDelete {
			if req.Body != nil {
				var bodyBytes []byte
				bodyBytes, _ = ioutil.ReadAll(req.Body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				contentType := req.Header.Get("Content-Type")
				if len(contentType) == 0 {
					contentType = "application/json"
				}
				requestBody := &RequestBody{
					Data:        bodyBytes,
					ContentType: contentType,
				}
				req = req.WithContext(genericrequest.WithValue(req.Context(), requestBodyKey, requestBody))
			}
		}
		handler.ServeHTTP(w, req)
	})
}
