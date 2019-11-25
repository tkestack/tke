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
	"tkestack.io/tke/pkg/auth/handler/localidentity"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterIdentityRoute to install route for identity http handler.
func RegisterIdentityRoute(container *restful.Container, handler *localidentity.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/localidentities")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("list all existing users").
		Operation("listLocalIdentities").
		Param(ws.QueryParameter("page", "page number").DataType("string").Required(false)).
		Param(ws.QueryParameter("page_size", "page size").DataType("string").Required(false)).
		Param(ws.QueryParameter("keyword", "filter user by keyword, query user name or display name").DataType("string").Required(false)).
		Param(ws.QueryParameter("name", "filter user by name").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.LocalIdentityList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.List))

	ws.Route(ws.
		POST(fmt.Sprintf("/")).
		Doc("create a new TKE User and returns it").
		Operation("createLocalIdentity").
		Reads(types.LocalIdentity{}).
		Returns(http.StatusCreated, "Created", types.LocalIdentity{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Create))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("get a existing user by given name").
		Operation("getLocalIdentity").
		Returns(http.StatusOK, "Ok", types.LocalIdentity{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Get))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("update a existing user").
		Operation("patchLocalIdentity").
		Reads(types.LocalIdentity{}).
		Returns(http.StatusOK, "Ok", types.LocalIdentity{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Update))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}/status", types.NameTag)).
		Doc("update a existing user status").
		Operation("patchLocalIdentityStatus").
		Reads(types.LocalIdentity{}).
		Returns(http.StatusOK, "Ok", types.LocalIdentity{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.UpdateStatus))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}/password", types.NameTag)).
		Doc("update a existing user password").
		Operation("patchLocalIdentityPassword").
		Reads(types.LocalIdentity{}).
		Returns(http.StatusOK, "Ok", types.LocalIdentity{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.UpdatePassword))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}", types.NameTag)).
		Doc("delete a user by given name").
		Operation("deleteLocalIdentity").
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Delete))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}/policies", types.NameTag)).
		Doc("list all policies user bond to").
		Operation("listLocalIdentityPolicies").
		Returns(http.StatusOK, "Ok", types.PolicyList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.ListPolicies))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}/roles", types.NameTag)).
		Doc("list all roles user related to").
		Operation("listLocalIdentityRoles").
		Returns(http.StatusOK, "Ok", types.RoleList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.ListRoles))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}/permissions", types.NameTag)).
		Doc("list all rules user related to").
		Operation("listLocalIdentityPerms").
		Returns(http.StatusOK, "Ok", types.Permission{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.ListPermissions))

	container.Add(ws)
}
