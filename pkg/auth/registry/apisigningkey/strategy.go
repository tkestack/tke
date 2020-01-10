/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package apisigningkey

import (
	"context"

	"tkestack.io/tke/pkg/apiserver/authentication"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"

	"tkestack.io/tke/api/auth"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for oidc api signing key.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating api signing key objects.
func NewStrategy() *Strategy {
	return &Strategy{auth.Scheme, namesutil.Generator}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {}

// NamespaceScoped is false for api signing key.
func (Strategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

// Validate validates a new api signing key.
func (Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" {
		return field.ErrorList{field.Forbidden(field.NewPath(""), "Please contact admin to create api sign key")}
	}
	return ValidateSigningKey(obj.(*auth.APISigningKey))
}

// AllowCreateOnUpdate is false for api signing key.
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate returns true if the object can be updated
// unconditionally (irrespective of the latest resource version), when there is
// no resource version specified in the object.
func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for a api signing key.
func (Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" {
		return field.ErrorList{field.Forbidden(field.NewPath(""), "Please contact admin to update api sign key")}
	}
	return ValidateSigningKeyUpdate(obj.(*auth.APISigningKey), old.(*auth.APISigningKey))
}
