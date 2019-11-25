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

package authz

import (
	"context"
	"net/http"

	"tkestack.io/tke/pkg/auth/filter"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"tkestack.io/tke/pkg/auth/authorization/util"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
)

// Handler handle permission authorization http request.
type Handler struct {
	authorizer authorizer.Authorizer
}

// NewHandler creates new authorizer handler object.
func NewHandler(authz authorizer.Authorizer) *Handler {
	return &Handler{authz}
}

// Authorize receive a subject access review request and determine the subject access.
func (h *Handler) Authorize(request *restful.Request, response *restful.Response) {
	accessReview := &types.SubjectAccessReview{}
	if err := request.ReadEntity(accessReview); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if errs := util.ValidateSubjectAccessReview(accessReview); len(errs) > 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(errs.ToAggregate().Error()).Status(), response.ResponseWriter)
		return
	}

	authorizationAttributes := util.AuthorizationAttributesFrom(accessReview.Spec)
	decision, reason, evaluationErr := h.authorizer.Authorize(authorizationAttributes)

	accessReview.Status = types.SubjectAccessReviewStatus{
		Allowed: decision == authorizer.DecisionAllow,
		Denied:  decision == authorizer.DecisionDeny,
		Reason:  reason,
	}
	if evaluationErr != nil {
		accessReview.Status.EvaluationError = evaluationErr.Error()
	}

	log.Debug("Receive authz request", log.Any("attribute", authorizationAttributes), log.Any("response", accessReview.Status))
	responsewriters.WriteRawJSON(http.StatusOK, accessReview, response.ResponseWriter)
}

// RestAuthorize receive a subject access review request and determine the subject access to be compatible with k8s restful attributes.
func (h *Handler) RestAuthorize(request *restful.Request, response *restful.Response) {
	accessReview := &types.SubjectAccessReview{}
	if err := request.ReadEntity(accessReview); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if errs := util.ValidateSubjectAccessReview(accessReview); len(errs) > 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(errs.ToAggregate().Error()).Status(), response.ResponseWriter)
		return
	}

	authorizationAttributes := util.AuthorizationAttributesFrom(accessReview.Spec)
	tkeAttributes := filter.ConvertTKEAttributes(context.Background(), authorizationAttributes)
	decision, reason, evaluationErr := h.authorizer.Authorize(tkeAttributes)
	accessReview.Status = types.SubjectAccessReviewStatus{
		Allowed: decision == authorizer.DecisionAllow,
		Denied:  decision == authorizer.DecisionDeny,
		Reason:  reason,
	}
	if evaluationErr != nil {
		accessReview.Status.EvaluationError = evaluationErr.Error()
	}
	log.Info("Receive reset authz request", log.Any("attribute", tkeAttributes), log.Any("response", accessReview.Status))
	responsewriters.WriteRawJSON(http.StatusOK, accessReview, response.ResponseWriter)
}

// BatchAuthorize receive multiple subject access reviews request and return determine results.
func (h *Handler) BatchAuthorize(request *restful.Request, response *restful.Response) {
	accessReview := &types.SubjectAccessReview{}
	if err := request.ReadEntity(accessReview); err != nil {
		log.Error("read entity failed", log.Err(err))
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), response.ResponseWriter)
		return
	}

	if errs := util.ValidateSubjectAccessReview(accessReview); len(errs) > 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(errs.ToAggregate().Error()).Status(), response.ResponseWriter)
		return
	}

	attributesList := util.AuthorizationAttributesListFrom(accessReview.Spec)

	accessReview.Status = types.SubjectAccessReviewStatus{AllowedList: []*types.AllowedResponse{}}
	for index, resAttr := range accessReview.Spec.ResourceAttributesList {
		decision, reason, _ := h.authorizer.Authorize(attributesList[index])
		accessReview.Status.AllowedList = append(accessReview.Status.AllowedList, &types.AllowedResponse{
			Resource: resAttr.Resource,
			Verb:     resAttr.Verb,
			Allowed:  decision == authorizer.DecisionAllow,
			Denied:   decision == authorizer.DecisionDeny,
			Reason:   reason,
		})
	}

	log.Info("Receive rest authz request", log.Any("attribute", attributesList), log.Any("response", accessReview.Status))
	responsewriters.WriteRawJSON(http.StatusOK, accessReview, response.ResponseWriter)
}
