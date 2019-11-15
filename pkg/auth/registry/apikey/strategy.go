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
	"context"
	"fmt"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/auth"

	namesutil "tkestack.io/tke/pkg/util/names"
	//"tkestack.io/tke/api/auth"
)

// Strategy implements verification logic for project.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	keySigner util.KeySigner
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating project objects.
func NewStrategy(signer util.KeySigner) *Strategy {
	return &Strategy{auth.Scheme, namesutil.Generator, signer}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldAPIKey := old.(*auth.APIKey)
	newAPIKey, _ := obj.(*auth.APIKey)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldAPIKey.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update project information", log.String("oldTenantID", oldAPIKey.Spec.TenantID), log.String("newTenantID", newAPIKey.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		newAPIKey.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for projects.
func (Strategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {

}

// Validate validates a new project.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateAPIkey(obj.(*auth.APIKey))
}

// AllowCreateOnUpdate is false for projects.
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

// ValidateUpdate is the default update validation for an end project.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateAPIKeyUpdate(obj.(*auth.APIKey), old.(*auth.APIKey))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	project, ok := obj.(*auth.APIKey)
	if !ok {
		return nil, nil, fmt.Errorf("not a project")
	}
	return labels.Set(project.ObjectMeta.Labels), ToSelectableFields(project), nil
}

// MatchMatchAPIKeyLocalIdentity returns a generic matcher for a given label and field selector.
func MatchAPIKey(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{"spec.tenantID"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(project *auth.APIKey) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&project.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": project.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// PasswordStrategy implements password generation for password.
type PasswordStrategy struct {
	*Strategy
}

// NewStatusStrategy create the StatusStrategy object by given strategy.
func NewPasswordStrategy(strategy *Strategy) *PasswordStrategy {
	return &PasswordStrategy{strategy}
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s PasswordStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	return
}

var _ rest.RESTCreateStrategy = &PasswordStrategy{}
