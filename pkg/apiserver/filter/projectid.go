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

	"tkestack.io/tke/pkg/platform/apiserver/filter"

	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
)

const (
	// ProjectIDKey defines the key representing the project id in the additional
	// information mapping table of the user information.
	ProjectIDKey = "projectid"

	projectContextKey = "projectID"
)

func WithProject(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		clusterName := req.Header.Get(filter.ProjectNameHeaderKey)
		if clusterName != "" {
			req = req.WithContext(genericrequest.WithValue(req.Context(), projectContextKey, clusterName))
		}
		handler.ServeHTTP(w, req)
	})
}

// ProjectIDFrom get the project id from request context.
func ProjectIDFrom(ctx context.Context) string {
	b, ok := ctx.Value(projectContextKey).(string)
	if !ok {
		return ""
	}
	return b
}
