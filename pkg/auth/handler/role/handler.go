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

package role

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"

	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
	pageutil "tkestack.io/tke/pkg/util/page"
)

// Handler handle role CRUD http request.
type Handler struct {
	tenantAdmin string
	roleService *Service
}

// NewHandler creates new role handler object.
func NewHandler(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer, tenantAdmin string) *Handler {
	return &Handler{roleService: NewRoleService(registry, policyEnforcer), tenantAdmin: tenantAdmin}
}

// Create to create a new role with policies.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	roleCreate := &types.Role{}
	if err := request.ReadEntity(roleCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && roleCreate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		roleCreate.TenantID = tenantID
	}

	if roleCreate.UserName == "" {
		roleCreate.UserName = userName
	}

	if err := h.roleService.validateRoleCreate(roleCreate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	roleCreated, err := h.roleService.CreateRole(roleCreate)
	if err != nil && err == etcd.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, roleCreated, response.ResponseWriter)
}

// Get to return a role by given ID.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
	}

	role, err := h.roleService.GetRole(tenantID, id)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("role"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, role, response.ResponseWriter)
}

// List to return roles for given options.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	roleList := &types.RoleList{Items: []*types.Role{}}
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	opt := &types.RoleOption{
		ID:       request.QueryParameter(types.IDTag),
		Name:     request.QueryParameter(types.NameTag),
		UserName: request.QueryParameter(types.UserTag),
		Keyword:  request.QueryParameter(types.KeywordTag),
		Scope:    request.QueryParameter(types.ScopeTag),
		TenantID: tenantID,
	}

	roleList, err := h.roleService.ListRoles(opt)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("role"), opt.ID).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	sort.SliceStable(roleList.Items, func(i, j int) bool {
		if roleList.Items[i].Type != roleList.Items[j].Type {
			return roleList.Items[i].Type > roleList.Items[j].Type
		}

		return roleList.Items[i].CreateAt.After(roleList.Items[j].CreateAt)
	})
	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(roleList.Items))
	pagin.Items = roleList.Items[start:end]
	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// Update to update a existing role metadata, description and name.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	roleUpdate := &types.Role{}
	if err := request.ReadEntity(roleUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}
	if tenantID == "" && roleUpdate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		roleUpdate.TenantID = tenantID
	}
	roleUpdate.ID = id

	roleGet, err := h.roleService.GetRole(tenantID, id)
	if err == nil && roleGet.Type == types.PreDefine && userName != h.tenantAdmin {
		responsewriters.WriteRawJSON(http.StatusForbidden, errors.NewForbidden(util.GroupResource("role"), roleGet.Name, fmt.Errorf("predefine role can be updated")), response.ResponseWriter)
		return
	}

	err = h.roleService.validateRoleCreate(roleUpdate)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	identityUpdate, err := h.roleService.UpdateRole(roleUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("role"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, identityUpdate, response.ResponseWriter)
}

// Delete to delete a existing role.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	roleGet, err := h.roleService.GetRole(tenantID, id)
	if err == nil && roleGet.Type == types.PreDefine && userName != h.tenantAdmin {
		responsewriters.WriteRawJSON(http.StatusForbidden, errors.NewForbidden(util.GroupResource("role"), roleGet.Name, fmt.Errorf("predefine role can be deleted")), response.ResponseWriter)
		return
	}

	err = h.roleService.DeleteRole(tenantID, id)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("role"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// ListAttachedPolicies to list all policies the role contains.
func (h *Handler) ListAttachedPolicies(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}
	attachedPolicies, err := h.roleService.ListRolePolicies(tenantID, id)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(attachedPolicies.Items))
	pagin.Items = attachedPolicies.Items[start:end]

	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// AttachPolicies to bind polices for role.
func (h *Handler) AttachPolicies(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	attachInfo := &types.AttachInfo{}
	if err := request.ReadEntity(attachInfo); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && attachInfo.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		attachInfo.TenantID = tenantID
	}
	attachInfo.ID = id

	err := h.roleService.AttachRolePolicies(attachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// DetachPolicies to unbind polices for role.
func (h *Handler) DetachPolicies(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	detachInfo := &types.AttachInfo{}
	if err := request.ReadEntity(detachInfo); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}
	if tenantID == "" && detachInfo.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		detachInfo.TenantID = tenantID
	}
	detachInfo.ID = id

	err := h.roleService.DetachRolePolicies(detachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// AttachUsers to bind role for users.
func (h *Handler) AttachUsers(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	attachInfo := &types.AttachInfo{}
	if err := request.ReadEntity(attachInfo); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && attachInfo.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		attachInfo.TenantID = tenantID
	}
	attachInfo.ID = id

	_, err := h.roleService.GetRole(attachInfo.TenantID, attachInfo.ID)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("role"), id).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	err = h.roleService.AttachUsersRole(attachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// DetachUsers to unbind role for users
func (h *Handler) DetachUsers(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	detachInfo := &types.AttachInfo{}
	if err := request.ReadEntity(detachInfo); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}
	if tenantID == "" && detachInfo.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		detachInfo.TenantID = tenantID
	}
	detachInfo.ID = id

	err := h.roleService.DetachUsersRole(detachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// ListAttachedUsers to list all users attached to the role.
func (h *Handler) ListAttachedUsers(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}
	attachedUsers, err := h.roleService.ListRoleUsers(tenantID, id)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, attachedUsers, response.ResponseWriter)
}
