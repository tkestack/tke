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
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/pkg/auth/handler/identityprovider"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterIdentityProviderRoute to install route for identity http handler.
func RegisterIdentityProviderRoute(container *restful.Container, handler *identityprovider.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/identityproviders")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("list all existing identityproviders").
		Operation("listIdentityProviders").
		Param(ws.QueryParameter("page", "page number").DataType("string").Required(false)).
		Param(ws.QueryParameter("page_size", "page size").DataType("string").Required(false)).
		Param(ws.QueryParameter("keyword", "filter idp by keyword, query idp id or name").DataType("string").Required(false)).
		Param(ws.QueryParameter("type", "filter idp by type").DataType("string").Required(false)).
		Param(ws.QueryParameter("id", "filter idp by id").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.IdentityProviderList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.List))

	ws.Route(ws.
		POST(fmt.Sprintf("/")).
		Doc("create a new IDP and returns it").
		Operation("createIdentityProvider").
		Reads(types.IdentityProvider{}).
		Returns(http.StatusCreated, "Created", types.IdentityProvider{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Create))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("get a existing idp by given id").
		Operation("getIdentityProvider").
		Returns(http.StatusOK, "Ok", types.IdentityProvider{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Get))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("update a existing idp").
		Operation("patchIdentityProvider").
		Reads(types.IdentityProvider{}).
		Returns(http.StatusOK, "Ok", types.IdentityProvider{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Update))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("delete a idp by given id").
		Operation("deleteIdentityProvider").
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Delete))

	container.Add(ws)
}
