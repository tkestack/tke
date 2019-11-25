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

package identityprovider

import (
	"net/http"

	"tkestack.io/tke/pkg/auth/authentication/tenant"

	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
	pageutil "tkestack.io/tke/pkg/util/page"

	"github.com/dexidp/dex/storage"
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

// Handler handle OIDC identity provider http request.
type Handler struct {
	idpService *Service
}

// NewHandler creates new OAuth2 identity provider handler object.
func NewHandler(dexStorage storage.Storage, helper *tenant.Helper) *Handler {
	helper.LoadResourceAllTenant()
	return &Handler{NewIdentidyProviderService(dexStorage, helper)}
}

// Create create a new OIDC identity provider and returns it.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	//_, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())
	idpCreate := &types.IdentityProvider{}
	if err := request.ReadEntity(idpCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err := validateIdentityProviderCreate(idpCreate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	idpCreated, err := h.idpService.CreateIdentityProvier(idpCreate)
	if err != nil && err == storage.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, idpCreated, response.ResponseWriter)
}

// Get gets a existing OIDC identity provider by given id.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)
	idp, err := h.idpService.GetIdentityProvier(id)

	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("identityProvider"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, idp, response.ResponseWriter)
}

// Delete returns a OIDC idp by given id.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)
	err := h.idpService.DeleteIdentityProvider(id)
	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("identityProvider"), id).Status(), response.ResponseWriter)
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

// List gets a list of existing OIDC idp.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	idpType := request.QueryParameter(types.TypeTag)
	id := request.QueryParameter(types.IDTag)
	keyword := request.QueryParameter(types.KeywordTag)
	log.Info("list idp", log.String("id", id), log.String("keyword", keyword))
	idpList, err := h.idpService.ListIdentityProvider(idpType, id, keyword)

	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("identityProvider"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(idpList.Items))
	pagin.Items = idpList.Items[start:end]

	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// Update updates a existing OIDC idp.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)
	if len(id) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("id is empty").Status(), response.ResponseWriter)
		return
	}

	idpUpdate := &types.IdentityProvider{}
	if err := request.ReadEntity(idpUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	idpUpdate.ID = id
	if err := validateIdentityProviderCreate(idpUpdate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	idpUpdate, err := h.idpService.UpdateIdentityProvider(idpUpdate)
	if err != nil && err == etcd.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("identityProvider"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, idpUpdate, response.ResponseWriter)
}
