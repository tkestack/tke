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

// see: https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication
const (

	// TokenReviewAPIVersion is tokenReview APIVersion value
	TokenReviewAPIVersion = "authentication.k8s.io/v1beta1"

	// TokenReviewKind is tokenReview kind value
	TokenReviewKind = "TokenReview"
)

/*
/auth/authn request body:
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "spec": {
    "token": "(BEARERTOKEN)"
  }
}
*/

// TokenReviewRequest attempts to authenticate a token to a known user.
//+k8s:openapi-gen=true
type TokenReviewRequest struct {
	// APIVersion is default as "authentication.k8s.io/v1beta1".
	APIVersion string `json:"apiVersion"`
	// Kind is default as "TokenReview".
	Kind string `json:"kind"`
	// Spec holds information about the request being evaluated.
	Spec TokenReviewSpec `json:"spec"`
}

// TokenReviewSpec is a description of the token authentication request.
//+k8s:openapi-gen=true
type TokenReviewSpec struct {
	// Token is the opaque bearer token.
	Token string `json:"token"`
}

/*
/auth/authn response body:
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": true,
    "user": {
      "username": "janedoe@example.com",
      "uid": "42",
      "groups": [
        "developers",
        "qa"
      ],
      "extra": {
        "extrafield1": [
          "extravalue1",
          "extravalue2"
        ]
      }
    }
  }
}

An unsuccessful request would return:
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": false
  }
}
*/

// TokenReviewResponse is response info for authenticating a token to a known user.
//+k8s:openapi-gen=true
type TokenReviewResponse struct {
	// APIVersion is default as "authentication.k8s.io/v1beta1".
	APIVersion string `json:"apiVersion"`
	// Kind is default as "TokenReview".
	Kind string `json:"kind"`
	// Status is filled in by the server and indicates whether the request can be authenticated.
	Status TokenReviewStatus `json:"status"`
}

// TokenReviewStatus is the result of the token authentication request.
//+k8s:openapi-gen=true
type TokenReviewStatus struct {
	// Authenticated indicates that the token was associated with a known user.
	Authenticated bool `json:"authenticated"`

	// User is the UserInfo associated with the provided token.
	User TokenReviewUser `json:"user,omitempty"`
}

// TokenReviewUser holds the information about the user needed to implement the
//+k8s:openapi-gen=true
type TokenReviewUser struct {
	// The name that uniquely identifies this user among all active users.
	Username string `json:"username"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	UID string `json:"uid"`
	// The names of groups this user is a part of.
	Groups []string `json:"groups"`
	// Any additional information provided by the authenticator.
	Extra map[string][]string `json:"extra"`
}
