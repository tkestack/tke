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

package authn

import (
	"context"
	"net/http"

	"k8s.io/apiserver/pkg/authentication/token/union"

	authv1 "k8s.io/api/authentication/v1"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"

	"tkestack.io/tke/pkg/util/log"
)

// Handler handle token authentication http request.
type Handler struct {
	tokenAuthenticator authenticator.Token
}

// NewHandler creates new local identity handler object.
func NewHandler(authTokenHandlers ...authenticator.Token) *Handler {
	return &Handler{union.New(authTokenHandlers...)}
}

// AuthenticateToken handles token authentication http request.
func (h *Handler) AuthenticateToken(request *restful.Request, response *restful.Response) {
	tokenReview := &authv1.TokenReview{}
	if err := request.ReadEntity(tokenReview); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	authResp, valid, err := h.tokenAuthenticator.AuthenticateToken(context.Background(), tokenReview.Spec.Token)
	if !valid || err != nil {
		log.Error("Failed to authenticate token", log.String("token", tokenReview.Spec.Token), log.Bool("valid", valid), log.Err(err))
		tokenReview.Status = authv1.TokenReviewStatus{Authenticated: false}
		responsewriters.WriteRawJSON(http.StatusUnauthorized, tokenReview, response.ResponseWriter)
		return
	}

	extra := map[string]authv1.ExtraValue{}

	for key, val := range authResp.User.GetExtra() {
		extra[key] = val
	}
	tokenResponse := &authv1.TokenReview{
		TypeMeta:   tokenReview.TypeMeta,
		ObjectMeta: tokenReview.ObjectMeta,
		Status: authv1.TokenReviewStatus{
			Authenticated: true,
			User: authv1.UserInfo{
				Username: authResp.User.GetName(),
				Groups:   authResp.User.GetGroups(),
				Extra:    extra,
			},
		},
	}

	responsewriters.WriteRawJSON(http.StatusOK, tokenResponse, response.ResponseWriter)
}
