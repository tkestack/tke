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

package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubjectAccessReview checks whether or not a user or group can perform an action.
//+k8s:openapi-gen=true
type SubjectAccessReview struct {
	metav1.TypeMeta `json:",inline"`

	// Spec holds information about the request being evaluated.
	Spec SubjectAccessReviewSpec `json:"spec"`

	// Status is filled in by the server and indicates whether the request is allowed or not.
	Status SubjectAccessReviewStatus `json:"status,omitempty"`
}

// SubjectAccessReviewSpec is a description of the access request.  Exactly one of ResourceAuthorizationAttributes
// and NonResourceAuthorizationAttributes must be set.
//+k8s:openapi-gen=true
type SubjectAccessReviewSpec struct {
	// ResourceAuthorizationAttributes describes information for a resource access request.
	ResourceAttributes *ResourceAttributes `json:"resourceAttributes,omitempty"`

	// ResourceAttributesList describes information for multi resource access request.
	ResourceAttributesList []*ResourceAttributes `json:"resourceAttributesList,omitempty"`

	// NonResourceAttributes describes information for a non-resource access request.
	NonResourceAttributes *NonResourceAttributes `json:"nonResourceAttributes,omitempty"`

	// User is the user you're testing for.
	User string `json:"user,omitempty"`
	// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer
	// it needs a reflection here.
	Extra map[string][]string `json:"extra,omitempty"`
	// UID information about the requesting user.
	UID string `json:"uid,omitempty"`
}

// SubjectAccessReviewStatus indicates whether the request is allowed or not
//+k8s:openapi-gen=true
type SubjectAccessReviewStatus struct {
	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed"`
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool `json:"denied,omitempty"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty"`
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string `json:"evaluationError,omitempty"`

	// AllowedList is the allowed response for batch authorization request.
	AllowedList []*AllowedResponse `json:"allowedList,omitempty"`
}

// ResourceAttributes includes the authorization attributes available for resource requests to the Authorizer interface.
// Only verb and resource fields could be considered, tke-auth will ignore other fields right now.
//+k8s:openapi-gen=true
type ResourceAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	// "" (empty) is defaulted for LocalSubjectAccessReviews
	// "" (empty) is empty for cluster-scoped resources
	// "" (empty) means "all" for namespace scoped resources from a SubjectAccessReview or SelfSubjectAccessReview
	Namespace string `json:"namespace,omitempty"`
	// Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	Verb string `json:"verb,omitempty"`
	// Group is the API Group of the Resource.  "*" means all.
	Group string `json:"group,omitempty"`
	// Version is the API Version of the Resource.  "*" means all.
	Version string `json:"version,omitempty"`
	// Resource is one of the existing resource types.  "*" means all.
	Resource string `json:"resource,omitempty"`
	// Subresource is one of the existing resource types.  "" means none.
	Subresource string `json:"subresource,omitempty"`
	// Name is the name of the resource being requested for a "get" or deleted for a "delete". "" (empty) means all.
	Name string `json:"name,omitempty"`
}

// NonResourceAttributes includes the authorization attributes available for non-resource requests to the Authorizer interface.
//+k8s:openapi-gen=true
type NonResourceAttributes struct {
	// Path is the URL path of the request.
	Path string `json:"path,omitempty"`
	// Verb is the standard HTTP verb.
	Verb string `json:"verb,omitempty"`
}

// AllowedResponse includes the resource access request and response.
//+k8s:openapi-gen=true
type AllowedResponse struct {
	// Path is the URL path of the request
	Resource string `json:"resource,omitempty"`
	// Verb is the standard HTTP verb
	Verb string `json:"verb,omitempty"`

	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed"`
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool `json:"denied,omitempty"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty"`
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string `json:"evaluationError,omitempty"`
}
