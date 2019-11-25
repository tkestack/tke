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
	"tkestack.io/tke/pkg/auth/handler/policy"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterPolicyRoute to install route for policy http handler.
func RegisterPolicyRoute(container *restful.Container, handler *policy.Handler) {
	ws := new(restful.WebService)
	ws.Path("/api/authv1/policies")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.
		GET("/").
		Doc("list all policies").
		Operation("listPolices").
		Param(ws.QueryParameter("page", "page number").DataType("string").Required(false)).
		Param(ws.QueryParameter("page_size", "page size").DataType("string").Required(false)).
		Param(ws.QueryParameter(types.IDTag, "filter by id").DataType("string").Required(false)).
		Param(ws.QueryParameter(types.NameTag, "filter by name").DataType("string").Required(false)).
		Param(ws.QueryParameter(types.UserTag, "filter by creator").DataType("string").Required(false)).
		Param(ws.QueryParameter(types.KeywordTag, "filter by keyword").DataType("string").Required(false)).
		Param(ws.QueryParameter(types.ScopeTag, "filter by scope, eg. all,local,system").DataType("string").Required(false)).
		Returns(http.StatusOK, "Ok", types.PolicyList{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.List))

	ws.Route(ws.
		POST("/").
		Doc("create a policy with action and resource").
		Operation("createPolicy").
		Reads(types.PolicyCreate{}).
		Returns(http.StatusCreated, "Created", types.Policy{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Create))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("get a policy by a given id").
		Operation("getPolicy").
		Returns(http.StatusOK, "Ok", types.Policy{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Get))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("update a policy").
		Operation("patchPolicy").
		Reads(types.Policy{}).
		Returns(http.StatusOK, "Ok", types.Policy{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Update))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}", types.IDTag)).
		Doc("delete a policy").
		Operation("deletePolicy").
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.Delete))

	ws.Route(ws.
		GET(fmt.Sprintf("/{%s}/users", types.IDTag)).
		Doc("list all users bind to the policy").
		Operation("listPolicyUsers").
		Returns(http.StatusOK, "Ok", types.AttachInfo{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.ListAttachedUsers))

	ws.Route(ws.
		PUT(fmt.Sprintf("/{%s}/users", types.IDTag)).
		Doc("bind users with the policy").
		Operation("patchPolicyUsers").
		Reads(types.AttachInfo{}).
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.AttachUsers))

	ws.Route(ws.
		DELETE(fmt.Sprintf("/{%s}/users", types.IDTag)).
		Doc("remove the policy from users").
		Operation("deletePolicyUsers").
		Reads(types.AttachInfo{}).
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(handler.DetachUsers))

	container.Add(ws)
}
