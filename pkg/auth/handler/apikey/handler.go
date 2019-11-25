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

package apikey

import (
	"net/http"
	"sort"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	pageutil "tkestack.io/tke/pkg/util/page"
)

// Handler handle token authentication http request.
type Handler struct {
	apiKeyAuth *authenticator.APIKeyAuthenticator
}

// NewHandler creates new api key handler object.
func NewHandler(apiKeyAuth *authenticator.APIKeyAuthenticator) *Handler {
	return &Handler{apiKeyAuth: apiKeyAuth}
}

// CreateAPIKey generates a new api key by given user and expiration.
func (h *Handler) CreateAPIKey(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	keyReq := types.APIKeyReq{}
	if err := request.ReadEntity(&keyReq); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if userName == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("Username is empty").Status(), response.ResponseWriter)
		return
	}

	keyData, err := h.apiKeyAuth.CreateToken(tenantID, userName, keyReq.Description, keyReq.Expire.Duration)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, keyData, response.ResponseWriter)
}

// CreateAPIKeyByPassword generates a new api key by given user and expiration.
func (h *Handler) CreateAPIKeyByPassword(request *restful.Request, response *restful.Response) {
	keyReq := types.APIKeyReqPassword{}
	if err := request.ReadEntity(&keyReq); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if keyReq.TenantID == "" || keyReq.UserName == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("TenantID or username is empty").Status(), response.ResponseWriter)
		return
	}

	keyData, err := h.apiKeyAuth.CreateTokenWithPassword(keyReq.TenantID, keyReq.UserName, keyReq.Password, keyReq.Description, keyReq.Expire.Duration)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusOK, keyData, response.ResponseWriter)
}

// UpdateAPIKey used to disable or delete api key.
func (h *Handler) UpdateAPIKey(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	keyReq := types.APIKeyData{}
	if err := request.ReadEntity(&keyReq); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if userName == "" || tenantID == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("Username is empty").Status(), response.ResponseWriter)
		return
	}

	err := h.apiKeyAuth.UpdateToken(&keyReq, tenantID, userName)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
		Status: v1.StatusSuccess,
		Code:   http.StatusNoContent,
	}, response.ResponseWriter)
}

// ListAPIKeys list all apikeys for the user.
func (h *Handler) ListAPIKeys(request *restful.Request, response *restful.Response) {
	userName, tenantID := authentication.GetUsernameAndTenantID(request.Request.Context())

	if userName == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("Username is empty").Status(), response.ResponseWriter)
		return
	}

	keyList, err := h.apiKeyAuth.ListAPIKeys(tenantID, userName)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	keyword := request.QueryParameter(types.KeywordTag)
	if len(keyword) != 0 {
		log.Info("query keyword", log.String("keyword", keyword))
		result := &types.APIKeyList{}
		for _, key := range keyList.Items {
			if util.CaseInsensitiveContains(key.Description, keyword) {
				result.Items = append(result.Items, key)
			}
		}

		keyList = result
	}

	sort.SliceStable(keyList.Items, func(i, j int) bool {
		return keyList.Items[i].IssueAt.After(keyList.Items[j].IssueAt)
	})
	page, size := pageutil.ParsePageParam(request)
	start, end, pagin := pageutil.Pagein(page, size, len(keyList.Items))
	pagin.Items = keyList.Items[start:end]
	responsewriters.WriteRawJSON(http.StatusOK, pagin, response.ResponseWriter)
}
