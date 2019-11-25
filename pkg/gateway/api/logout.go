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

package api

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"tkestack.io/tke/pkg/gateway/token"
)

// Empty defines a data structure containing nothing.
type Empty struct {
}

func registerLogoutRoute(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/logout", GroupName, Version))
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON, restful.MIME_OCTET)
	ws.Route(ws.
		GET("/").
		Doc("logout current user").
		Operation("doLogout").
		Returns(http.StatusOK, "Ok", Empty{}).
		To(handleLogoutFunc()))
	container.Add(ws)
}

func handleLogoutFunc() func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		token.DeleteCookie(response.ResponseWriter)
		responsewriters.WriteRawJSON(http.StatusOK, Empty{}, response.ResponseWriter)
	}
}
