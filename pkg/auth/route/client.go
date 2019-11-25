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
	"tkestack.io/tke/pkg/auth/handler/client"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterClientRoute to install route for oath2 client http handler.
func RegisterClientRoute(container *restful.Container, handler *client.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/clients")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("list all existing clients").
		Operation("listClients").
		Param(ws.QueryParameter("page", "page number").DataType("string").Required(false)).
		Param(ws.QueryParameter("page_size", "page size").DataType("string").Required(false)).
		Param(ws.QueryParameter("keyword", "filter client by keyword, query client name or display name").DataType("string").Required(false)).
		Param(ws.QueryParameter("id", "filter client by id").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.ClientList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.List))

	ws.Route(ws.
		POST(fmt.Sprintf("/")).
		Doc("create a new client and returns it").
		Operation("createClient").
		Reads(types.Client{}).
		Returns(http.StatusCreated, "Created", types.Client{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Create))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("get a existing client by given id").
		Operation("getClient").
		Returns(http.StatusOK, "Ok", types.Client{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Get))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("update a existing client").
		Operation("patchClient").
		Reads(types.Client{}).
		Returns(http.StatusOK, "Ok", types.Client{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Update))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("delete a client by given id").
		Operation("deleteClient").
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Delete))

	container.Add(ws)
}
