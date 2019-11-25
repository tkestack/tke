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

package route

import (
	"github.com/emicklei/go-restful"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"tkestack.io/tke/pkg/auth/handler/apikey"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterAPIKeyRoute to install route for api key http handler.
func RegisterAPIKeyRoute(container *restful.Container, apiKeyHandler *apikey.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/apikey")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.
		POST("").
		Doc("generate a new api key by carrying oidc token header.").
		Operation("createApiKey").
		Reads(types.APIKeyReq{}).
		Returns(http.StatusOK, "Ok", types.APIKeyData{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(apiKeyHandler.CreateAPIKey))

	ws.Route(ws.
		POST("password").
		Doc("generate a new api key by username and password.").
		Operation("createApiKeyByPassword").
		Reads(types.APIKeyReqPassword{}).
		Returns(http.StatusOK, "Ok", types.APIKeyData{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(apiKeyHandler.CreateAPIKeyByPassword))

	ws.Route(ws.
		PUT("/").
		Doc("enable/disable or delete a apikey.").
		Operation("patchAPIKey").
		Reads(types.APIKeyData{}).
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(apiKeyHandler.UpdateAPIKey))

	ws.Route(ws.
		GET("/").
		Doc("list api keys for user.").
		Operation("listAPIKeys").
		Param(ws.QueryParameter("page", "page number").DataType("string").Required(false)).
		Param(ws.QueryParameter("page_size", "page size").DataType("string").Required(false)).
		Param(ws.QueryParameter("keyword", "filter apikey by keyword, match description").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.APIKeyList{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(apiKeyHandler.ListAPIKeys))

	container.Add(ws)
}
