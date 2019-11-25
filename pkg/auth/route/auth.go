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
	"net/http"

	"tkestack.io/tke/pkg/auth/handler/authn"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/pkg/auth/handler/authz"
	"tkestack.io/tke/pkg/auth/types"
)

// RegisterAuthRoute registers the http handlers of authz webhook for kubernetes.
func RegisterAuthRoute(container *restful.Container, authnHandler *authn.Handler, authzHandler *authz.Handler) {
	ws := new(restful.WebService)
	ws.Path("/auth")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.
		POST("/authn").
		Doc("verify token").
		Operation("getToken").
		Reads(types.TokenReviewRequest{}).
		Returns(http.StatusOK, "Ok", types.TokenReviewResponse{}).
		Returns(http.StatusUnauthorized, "Unauthorized", types.TokenReviewResponse{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(authnHandler.AuthenticateToken))

	ws.Route(ws.
		POST("/authz").
		Doc("receive a subject access review request and determine the subject access.").
		Operation("getAuthz").
		Reads(types.SubjectAccessReview{}).
		Returns(http.StatusOK, "Ok", types.SubjectAccessReview{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(authzHandler.Authorize))

	ws.Route(ws.
		POST("/restauthz").
		Doc("receive a subject access review request like k8s and determine the subject access.").
		Operation("getHttpAuthz").
		Reads(types.SubjectAccessReview{}).
		Returns(http.StatusOK, "Ok", types.SubjectAccessReview{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(authzHandler.RestAuthorize))

	ws.Route(ws.
		POST("/batchauthz").
		Doc("receive multiple subject access reviews request and return determine results.").
		Operation("getBatchAuthz").
		Reads(types.SubjectAccessReview{}).
		Returns(http.StatusOK, "Ok", types.SubjectAccessReview{}).
		Returns(http.StatusBadRequest, "BadRequest", v1.Status{}).
		To(authzHandler.BatchAuthorize))

	container.Add(ws)
}
