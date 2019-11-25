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
	"tkestack.io/tke/pkg/auth/handler/category"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterCategoryRoute to install route for category Category http handler.
func RegisterCategoryRoute(container *restful.Container, handler *category.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/categories")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("list all categories").
		Operation("listCategories").
		Param(ws.QueryParameter("name", "filter by name").DataType("string").Required(false)).
		Param(ws.QueryParameter("keyword", "filter by keyword").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.CategoryList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.List))

	ws.Route(ws.
		POST("/").
		Doc("create a category with actions").
		Operation("createCategory").
		Reads(types.Category{}).
		Returns(http.StatusCreated, "Created", types.Category{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Create))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("get a category by a given nanme").
		Operation("getCategory").
		Returns(http.StatusOK, "Ok", types.Category{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Get))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("update a category metadata").
		Operation("patchCategory").
		Reads(types.Category{}).
		Returns(http.StatusOK, "Ok", types.Category{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Update))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("delete a category").
		Operation("deleteCategory").
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Delete))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}/actions", types.NameTag)).
		Doc("add new actions for the category").
		Operation("createActions").
		Reads(types.Category{}).
		Returns(http.StatusOK, "Ok", types.AttachInfo{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.AddActions))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}/actions", types.NameTag)).
		Doc("delete actions of the category").
		Operation("deleteActions").
		Reads(types.AttachInfo{}).
		Returns(http.StatusOK, "Ok", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.DeleteActions))

	container.Add(ws)
}
