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

package authentication

import (
	"fmt"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"tkestack.io/tke/pkg/registry/chartmuseum/model"
)

// WithAuthentication creates an http handler that tries to authenticate the
// given chartmuseum request as a user, and then stores any such user found onto
// the provided context for the request.
func WithAuthentication(handler http.Handler, externalScheme string) http.Handler {
	return &authentication{
		handler:        handler,
		externalScheme: externalScheme,
	}
}

type authentication struct {
	handler        http.Handler
	externalScheme string
}

func (a *authentication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// todo: authenticate
	a.handler.ServeHTTP(w, req)
}

func (a *authentication) notAuthenticated(w http.ResponseWriter, req *http.Request) {
	realm := fmt.Sprintf("%s://%s", a.externalScheme, req.Host)
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
	err := &model.ErrorResponse{Error: "unauthorized"}
	responsewriters.WriteRawJSON(http.StatusUnauthorized, err, w)
}
