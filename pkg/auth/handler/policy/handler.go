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

package policy

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	pageutil "tkestack.io/tke/pkg/util/page"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
)

// Handler to handle policy CRUD http request.
type Handler struct {
	tenantAdmin   string
	policyService *Service
}

// NewHandler creates new policy handler object.
func NewHandler(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer, tenantAdmin string) *Handler {
	return &Handler{policyService: NewPolicyService(registry, policyEnforcer), tenantAdmin: tenantAdmin}
}

// Create to create a new policy.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	policyCreate := &types.PolicyCreate{}
	if err := request.ReadEntity(policyCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && policyCreate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		policyCreate.TenantID = tenantID
	}

	var attachUsers []string
	if policyCreate.UserName != "" {
		attachUsers = strings.Split(policyCreate.UserName, ",")
	}

	pol := types.Policy{
		Name:        policyCreate.Name,
		TenantID:    policyCreate.TenantID,
		Service:     policyCreate.Service,
		Statement:   policyCreate.Statement,
		UserName:    userName,
		Description: policyCreate.Description,
	}

	if err := validatePolicyCreate(&pol); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	policyCreated, err := h.policyService.CreatePolicy(&pol, attachUsers)
	if err != nil && err == etcd.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, policyCreated, response.ResponseWriter)
}

// Get to return a policy by given id.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
	}

	policy, err := h.policyService.GetPolicy(tenantID, id)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), id).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, policy, response.ResponseWriter)
}

// List to return policies for given owner.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	opt := &types.PolicyOption{
		ID:       request.QueryParameter(types.IDTag),
		Name:     request.QueryParameter(types.NameTag),
		UserName: request.QueryParameter(types.UserTag),
		Keyword:  request.QueryParameter(types.KeywordTag),
		Scope:    request.QueryParameter(types.ScopeTag),
		TenantID: tenantID,
	}

	policyList, err := h.policyService.ListPolicies(opt)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), opt.ID).Status(), response.ResponseWriter)
		return
	}

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

// Update to update a existing policy.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	policyUpdate := &types.Policy{}
	if err := request.ReadEntity(policyUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && policyUpdate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		policyUpdate.TenantID = tenantID
	}
	policyUpdate.ID = id

	err := validatePolicyCreate(policyUpdate)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	policyGet, err := h.policyService.GetPolicy(tenantID, id)
	if err == nil && policyGet.Type == types.PreDefine && userName != h.tenantAdmin {
		responsewriters.WriteRawJSON(http.StatusForbidden, errors.NewForbidden(util.GroupResource("policy"), policyGet.Name, fmt.Errorf("predefine policy can be updated")), response.ResponseWriter)
		return
	}

	policyUpdated, err := h.policyService.UpdatePolicy(policyUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, policyUpdated, response.ResponseWriter)
}

// Delete to delete a existing policy.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	policyGet, err := h.policyService.GetPolicy(tenantID, id)
	if err == nil && policyGet.Type == types.PreDefine && userName != h.tenantAdmin {
		responsewriters.WriteRawJSON(http.StatusForbidden, errors.NewForbidden(util.GroupResource("policy"), policyGet.Name, fmt.Errorf("predefine policy can be deleted")), response.ResponseWriter)
		return
	}

	err = h.policyService.DeletePolicy(tenantID, id)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), id).Status(), response.ResponseWriter)
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

// AttachUsers to bind policy into users.
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

	_, err := h.policyService.GetPolicy(attachInfo.TenantID, attachInfo.ID)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), id).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	err = h.policyService.AttachPolicyUsers(attachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// DetachUsers to detach users from policy.
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

	_, err := h.policyService.GetPolicy(detachInfo.TenantID, detachInfo.ID)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("policy"), id).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	err = h.policyService.DetachPolicyUsers(detachInfo)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// ListAttachedUsers to list all users attached to the policy.
func (h *Handler) ListAttachedUsers(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	attachedUsers, err := h.policyService.ListPolicyUsers(tenantID, id)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, attachedUsers, response.ResponseWriter)
}

// Service returns policy handler service.
func (h *Handler) Service() *Service {
	return h.policyService
}
