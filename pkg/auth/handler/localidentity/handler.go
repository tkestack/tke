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

package localidentity

import (
	"net/http"
	"sort"
	"strings"

	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
	pageutil "tkestack.io/tke/pkg/util/page"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

// Handler handle local identity CRUD http request.
type Handler struct {
	identityService *Service
}

// NewHandler creates new local identity handler object.
func NewHandler(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer) *Handler {
	return &Handler{identityService: NewLocalIdentityService(registry, policyEnforcer)}
}

// Create create a new TKE User and returns it.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	identityCreate := &types.LocalIdentity{}
	if err := request.ReadEntity(identityCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && identityCreate.Spec.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		identityCreate.Spec.TenantID = tenantID
	}

	if err := validateLocalIdentityCreate(identityCreate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	identityCreated, err := h.identityService.CreateLocalIdentity(identityCreate)
	if err != nil && err == etcd.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, identityCreated, response.ResponseWriter)
}

// Get to get a existing user by given name.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
	}

	identity, err := h.identityService.GetLocalIdentity(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, identity, response.ResponseWriter)
}

// List to get some existing user.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	keyword := request.QueryParameter(types.KeywordTag)
	name := request.QueryParameter(types.NameTag)
	identityList := &types.LocalIdentityList{}

	// if specify name, will get the identity named
	if len(name) != 0 {
		identity, err := h.identityService.GetLocalIdentity(tenantID, name)
		if err != nil && err == etcd.ErrNotFound {
			responsewriters.WriteRawJSON(http.StatusOK, identityList, response.ResponseWriter)
			return
		}

		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
			return
		}

		identityList.Items = append(identityList.Items, identity)
		responsewriters.WriteRawJSON(http.StatusOK, identityList, response.ResponseWriter)
		return
	}

	identityList, err := h.identityService.ListLocalIdentity(tenantID)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	if len(keyword) != 0 {
		result := &types.LocalIdentityList{}
		for _, ident := range identityList.Items {
			if util.CaseInsensitiveContains(ident.Name, keyword) {
				result.Items = append(result.Items, ident)
			} else if displayName, ok := ident.Spec.Extra["displayName"]; ok && util.CaseInsensitiveContains(displayName, keyword) {
				result.Items = append(result.Items, ident)
			}
		}
		identityList = result
	}

	sort.SliceStable(identityList.Items, func(i, j int) bool {
		return identityList.Items[i].CreateAt.After(identityList.Items[j].CreateAt)
	})

	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(identityList.Items))
	pagin.Items = identityList.Items[start:end]

	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// Update to update a existing user.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	identityUpdate := &types.LocalIdentity{}
	if err := request.ReadEntity(identityUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && identityUpdate.Spec.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		identityUpdate.Spec.TenantID = tenantID
	}
	identityUpdate.Name = name

	if err := validateLocalIdentityUpdate(identityUpdate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	identityUpdate, err := h.identityService.UpdateLocalIdentity(identityUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, identityUpdate, response.ResponseWriter)
}

// UpdateStatus to update a existing user status, locked or unlocked.
func (h *Handler) UpdateStatus(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	identityUpdate := &types.LocalIdentity{}
	if err := request.ReadEntity(identityUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && identityUpdate.Spec.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		identityUpdate.Spec.TenantID = tenantID
	}

	if name != identityUpdate.Name {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is invalid").Status(), response.ResponseWriter)
		return
	}

	identityUpdate, err := h.identityService.UpdateLocalIdentityStatus(identityUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, identityUpdate, response.ResponseWriter)
}

// UpdatePassword to update a existing user password.
func (h *Handler) UpdatePassword(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	identityUpdate := &types.LocalIdentity{}
	if err := request.ReadEntity(identityUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && identityUpdate.Spec.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		identityUpdate.Spec.TenantID = tenantID
	}

	if name != identityUpdate.Name {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is invalid").Status(), response.ResponseWriter)
		return
	}

	if err := validateLocalIdentityUpdate(identityUpdate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	identityUpdate, err := h.identityService.UpdateLocalIdentityPassword(identityUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil && strings.Contains(err.Error(), "password") {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, identityUpdate, response.ResponseWriter)
}

// Delete to delete a user by given name.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	err := h.identityService.DeleteLocalIdentity(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
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

// ListPolicies returns all policies of the user related to.
func (h *Handler) ListPolicies(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}

	// Todo check user exists with 3rd idp
	_, err := h.identityService.GetLocalIdentity(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	policyList, err := h.identityService.ListUserPolicies(tenantID, name)
	log.Info("List policies for user:", log.String("name", name), log.Any("policies", policyList))
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	sort.SliceStable(policyList.Items, func(i, j int) bool {
		if policyList.Items[i].Type != policyList.Items[j].Type {
			return policyList.Items[i].Type > policyList.Items[j].Type
		}

		return policyList.Items[i].CreateAt.After(policyList.Items[j].CreateAt)
	})

	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(policyList.Items))
	pagin.Items = policyList.Items[start:end]

	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// ListRoles returns all roles of the user related to.
func (h *Handler) ListRoles(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}

	_, err := h.identityService.GetLocalIdentity(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	roleList, err := h.identityService.ListUserRoles(tenantID, name)
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

// ListPermissions to get all permissions allowed in roles and rules of user related to.
func (h *Handler) ListPermissions(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}

	_, err := h.identityService.GetLocalIdentity(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("localIdentity"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	permissions, err := h.identityService.ListUserPerms(tenantID, name)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, permissions, response.ResponseWriter)
}
