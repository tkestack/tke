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

package client

import (
	"net/http"

	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	pageutil "tkestack.io/tke/pkg/util/page"

	"github.com/dexidp/dex/storage"
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

// Handler handle OAuth2 client http request.
type Handler struct {
	clientService *Service
}

// NewHandler creates new OAuth2 client handler object.
func NewHandler(dexStorage storage.Storage) *Handler {
	return &Handler{NewClientService(dexStorage)}
}

//Create creates a new OAuth2 client and returns it.
func (h *Handler) Create(request *restful.Request, response *restful.Response) {
	clientCreate := &types.Client{}
	if err := request.ReadEntity(clientCreate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err := validateClientCreate(clientCreate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	clientCreated, err := h.clientService.CreateClient(clientCreate)
	if err != nil && err == storage.ErrAlreadyExists {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusCreated, clientCreated, response.ResponseWriter)
}

// Get returns a OAuth2 client by given id.
func (h *Handler) Get(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)
	cli, err := h.clientService.GetClient(id)

	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("client"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, cli, response.ResponseWriter)
}

// Delete deletes a OAuth2 client by given id.
func (h *Handler) Delete(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)
	err := h.clientService.DeleteClient(id)
	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("client"), id).Status(), response.ResponseWriter)
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

// List gets a list of existing OAuth2 clients.
func (h *Handler) List(request *restful.Request, response *restful.Response) {
	keyword := request.QueryParameter(types.KeywordTag)
	id := request.QueryParameter(types.IDTag)
	clientList, err := h.clientService.ListClient(id, keyword)
	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("client"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(clientList.Items))
	pagin.Items = clientList.Items[start:end]
	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}

// Update updates a existing OAuth2 client.
func (h *Handler) Update(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(types.IDTag)

	clientUpdate := &types.Client{}
	if err := request.ReadEntity(clientUpdate); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	clientUpdate.ID = id
	if err := validateClientCreate(clientUpdate); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	cliUpdate, err := h.clientService.UpdateClient(clientUpdate)
	if err != nil && err == storage.ErrNotFound {
		responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(util.GroupResource("client"), id).Status(), response.ResponseWriter)
		return
	}

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, cliUpdate, response.ResponseWriter)
}
