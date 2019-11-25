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

package category

import (
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
)

// Handler handle local category CRUD http request.
type Handler struct {
	categoryService *Service
}

// NewHandler creates new local category handler object.
func NewHandler(registry *registry.Registry) *Handler {
	return &Handler{categoryService: NewCategoryService(registry)}
}

// Create creates a new policy action category and returns it.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	categoryCreate := &types.Category{}
	if err := request.ReadEntity(categoryCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && categoryCreate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		categoryCreate.TenantID = tenantID
	}

	if err := validateCategoryCreate(categoryCreate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	categoryCreated, err := h.categoryService.CreateCategory(categoryCreate)
	if err != nil && err == etcd.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, categoryCreated, response.ResponseWriter)
}

// Get gets a category by a given name.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)

	category, err := h.categoryService.GetCategory(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("category"), name).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, category, response.ResponseWriter)
}

// List gets all existing categoryList for specify tenant.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	keyword := request.QueryParameter(types.KeywordTag)
	name := request.QueryParameter(types.NameTag)
	categoryList := &types.CategoryList{}

	// if specify name, will get the category named
	if len(name) != 0 {
		category, err := h.categoryService.GetCategory(tenantID, name)
		if err != nil && err == etcd.ErrNotFound {
			responsewriters.WriteRawJSON(http.StatusOK, categoryList, response.ResponseWriter)
			return
		}

		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
			return
		}

		categoryList.Items = append(categoryList.Items, category)
		responsewriters.WriteRawJSON(http.StatusOK, categoryList, response.ResponseWriter)
		return
	}

	categoryList, err := h.categoryService.ListCategory(tenantID)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	if len(keyword) == 0 {
		responsewriters.WriteRawJSON(http.StatusOK, categoryList, response.ResponseWriter)
	} else {
		result := &types.CategoryList{}
		for _, category := range categoryList.Items {
			if util.CaseInsensitiveContains(category.Name, keyword) || util.CaseInsensitiveContains(category.DisplayName, keyword) {
				result.Items = append(result.Items, category)
			}
		}
		responsewriters.WriteRawJSON(http.StatusOK, result, response.ResponseWriter)
	}
}

// Update updates a category metadata, include name and description, not actions.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	categoryUpdate := &types.Category{}
	if err := request.ReadEntity(categoryUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && categoryUpdate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		categoryUpdate.TenantID = tenantID
	}

	categoryUpdate.Name = name

	if err := validateCategoryCreate(categoryUpdate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	categoryUpdate, err := h.categoryService.UpdateCategory(categoryUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("category"), name).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, categoryUpdate, response.ResponseWriter)
}

// DeleteActions deletes actions in a category.
func (h *Handler) DeleteActions(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	categoryUpdate := &types.Category{}
	if err := request.ReadEntity(categoryUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && categoryUpdate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		categoryUpdate.TenantID = tenantID
	}

	categoryUpdate, err := h.categoryService.DeleteActions(categoryUpdate.TenantID, name, categoryUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("category"), name).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, categoryUpdate, response.ResponseWriter)
}

// AddActions adds actions in a category.
func (h *Handler) AddActions(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	categoryUpdate := &types.Category{}
	if err := request.ReadEntity(categoryUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if tenantID == "" && categoryUpdate.TenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("tenantID must be specified").Status(), response.ResponseWriter)
		return
	}
	if tenantID != "" {
		categoryUpdate.TenantID = tenantID
	}

	categoryUpdate, err := h.categoryService.AddActions(categoryUpdate.TenantID, name, categoryUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("category"), name).Status(), response.ResponseWriter)
		return
	}
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, categoryUpdate, response.ResponseWriter)
}

// Delete to delete a user by given name.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	name := request.PathParameter(types.NameTag)
	if len(name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("name is empty").Status(), response.ResponseWriter)
		return
	}

	err := h.categoryService.DeleteCategory(tenantID, name)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("category"), name).Status(), response.ResponseWriter)
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
