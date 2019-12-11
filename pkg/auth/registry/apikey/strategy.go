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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/util"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for project.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	keySigner          util.KeySigner
	privilegedUsername string
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating project objects.
func NewStrategy(keySigner util.KeySigner, privilegedUsername string) *Strategy {
	return &Strategy{auth.Scheme, namesutil.Generator, keySigner, privilegedUsername}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newAPIKey, _ := obj.(*auth.APIKey)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		newAPIKey.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for projects.
func (Strategy) NamespaceScoped() bool {
	return false
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	apiKey := obj.(*auth.APIKey)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		apiKey.Spec.TenantID = tenantID
	}

	if apiKey.Name == "" && apiKey.GenerateName == "" {
		apiKey.GenerateName = "apk"
	}
	return
}

// Validate validates a new project.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateAPIkey(obj.(*auth.APIKey), s.keySigner, s.privilegedUsername)
}

// AllowCreateOnUpdate is false for projects.
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// Decorator is intended for
// removing hashed password for identity or list of identities on returned from the
// underlying storage, since they cannot be watched.
func Decorator(obj runtime.Object) error {
	now := metav1.Now()
	if apiKey, ok := obj.(*auth.APIKey); ok {
		if apiKey.Spec.ExpireAt.Before(&now) {
			apiKey.Status.Expired = true
		} else {
			apiKey.Status.Expired = false
		}
		return nil
	}

	if apiKeyList, ok := obj.(*auth.APIKeyList); ok {
		for i := range apiKeyList.Items {
			if apiKeyList.Items[i].Spec.ExpireAt.Before(&now) {
				apiKeyList.Items[i].Status.Expired = true
			} else {
				apiKeyList.Items[i].Status.Expired = false
			}
		}
		return nil
	}

	return fmt.Errorf("unknown type")
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
		return nil, nil, fmt.Errorf("not a apikey")
	}
	return labels.Set(project.ObjectMeta.Labels), ToSelectableFields(project), nil
}

// MatchMatchAPIKeyLocalIdentity returns a generic matcher for a given label and field selector.
func MatchAPIKey(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{"spec.tenantID", "spec.apiKey", "spec.username"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(apiKey *auth.APIKey) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&apiKey.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": apiKey.Spec.TenantID,
		"spec.apiKey":   apiKey.Spec.APIkey,
		"spec.username": apiKey.Spec.Username,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Machine.
type StatusStrategy struct {
	*Strategy
}

var _ rest.RESTUpdateStrategy = &StatusStrategy{}

// NewStatusStrategy create the StatusStrategy object by given strategy.
func NewStatusStrategy(strategy *Strategy) *StatusStrategy {
	return &StatusStrategy{strategy}
}

// PrepareForUpdate is invoked on update before validation to normalize
// the object.  For example: remove fields that are not to be persisted,
// sort order-insensitive list fields, etc.  This should not remove fields
// whose presence would be considered a validation error.
func (StatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newAPIKey := obj.(*auth.APIKey)
	oldAPIKey := old.(*auth.APIKey)
	newAPIKey.Spec = oldAPIKey.Spec
	newAPIKey.Status.Expired = oldAPIKey.Status.Expired
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
