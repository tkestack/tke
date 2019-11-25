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

package platform

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for platform.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	businessClient *businessinternalclient.BusinessClient
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating platform objects.
func NewStrategy(businessClient *businessinternalclient.BusinessClient) *Strategy {
	return &Strategy{business.Scheme, namesutil.Generator, businessClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldPlatform := old.(*business.Platform)
	platform, _ := obj.(*business.Platform)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldPlatform.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update platform information", log.String("oldTenantID", oldPlatform.Spec.TenantID), log.String("newTenantID", platform.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		platform.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for platforms.
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
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	platform, _ := obj.(*business.Platform)
	if len(tenantID) != 0 {
		platform.Spec.TenantID = tenantID
	}

	if platform.Name == "" && platform.GenerateName == "" {
		platform.GenerateName = "platform-"
	}
}

// Validate validates a new platform.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidatePlatform(obj.(*business.Platform), s.businessClient)
}

// AllowCreateOnUpdate is false for platforms.
func (Strategy) AllowCreateOnUpdate() bool {
	return true
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

// ValidateUpdate is the default update validation for an end platform.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidatePlatformUpdate(obj.(*business.Platform), old.(*business.Platform))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	platform, ok := obj.(*business.Platform)
	if !ok {
		return nil, nil, fmt.Errorf("not a platform")
	}
	return labels.Set(platform.ObjectMeta.Labels), ToSelectableFields(platform), nil
}

// MatchPlatform returns a generic matcher for a given label and field selector.
func MatchPlatform(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(platform *business.Platform) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&platform.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": platform.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
