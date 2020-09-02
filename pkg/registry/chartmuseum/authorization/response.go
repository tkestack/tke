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

package authorization

import (
	"fmt"
	"net/http"

	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"tkestack.io/tke/pkg/registry/chartmuseum/model"
)

func (a *authorization) notFound(w http.ResponseWriter) {
	err := &model.ErrorResponse{Error: "not found"}
	responsewriters.WriteRawJSON(http.StatusNotFound, err, w)
}

func (a *authorization) internalError(w http.ResponseWriter) {
	err := &model.ErrorResponse{Error: "internal error"}
	responsewriters.WriteRawJSON(http.StatusInternalServerError, err, w)
}

func (a *authorization) locked(w http.ResponseWriter) {
	err := &model.ErrorResponse{Error: "locked chart repository"}
	responsewriters.WriteRawJSON(http.StatusLocked, err, w)
}

// func (a *authorization) forbidden(w http.ResponseWriter) {
// 	err := &model.ErrorResponse{Error: "forbidden"}
// 	responsewriters.WriteRawJSON(http.StatusForbidden, err, w)
// }

func (a *authorization) notAuthenticated(w http.ResponseWriter, req *http.Request) {
	realm := fmt.Sprintf("%s://%s", a.externalScheme, req.Host)
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
	err := &model.ErrorResponse{Error: "unauthorized"}
	responsewriters.WriteRawJSON(http.StatusUnauthorized, err, w)
}
