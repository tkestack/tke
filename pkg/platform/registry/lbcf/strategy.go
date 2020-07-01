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

package lbcf

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/controller/addon/lbcf/images"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for LBCF.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	platformClient platforminternalclient.PlatformInterface
}

var _ rest.RESTCreateStrategy = &Strategy{}
var _ rest.RESTUpdateStrategy = &Strategy{}
var _ rest.RESTDeleteStrategy = &Strategy{}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating namespace set objects.
func NewStrategy(platformClient platforminternalclient.PlatformInterface) *Strategy {
	return &Strategy{platform.Scheme, namesutil.Generator, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// NamespaceScoped is false for namespaceSets
func (Strategy) NamespaceScoped() bool {
	return false
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	lbcf, _ := obj.(*platform.LBCF)

	if len(tenantID) != 0 {
		lbcf.Spec.TenantID = tenantID
	}

	if lbcf.Name == "" && lbcf.GenerateName == "" {
		lbcf.GenerateName = "lbcf-"
	}

	if lbcf.Spec.Version == "" {
		lbcf.Spec.Version = images.LatestVersion
	}
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		oldLBCF := old.(*platform.LBCF)
		lbcf, _ := obj.(*platform.LBCF)
		if oldLBCF.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update lbcf information", log.String("oldTenantID", oldLBCF.Spec.TenantID), log.String("newTenantID", lbcf.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		lbcf.Spec.TenantID = tenantID
	}
}

// Validate validates a new LBCF.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateLBCF(ctx, obj.(*platform.LBCF), s.platformClient)
}

// AllowCreateOnUpdate is false for persistent events
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

// ValidateUpdate is the default update validation for an end namespace set.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLBCFUpdate(ctx, obj.(*platform.LBCF), old.(*platform.LBCF), s.platformClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	lbcf, ok := obj.(*platform.LBCF)
	if !ok {
		log.Fatalf("GetAttrs type = %v", reflect.TypeOf(obj))
	}
	return lbcf.ObjectMeta.Labels, ToSelectableFields(lbcf), nil
}

// MatchLBCF returns a generic matcher for a given label and field selector.
func MatchLBCF(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.clusterName",
			"spec.version",
			"status.version",
			"status.phase"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(lbcf *platform.LBCF) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&lbcf.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    lbcf.Spec.TenantID,
		"spec.clusterName": lbcf.Spec.ClusterName,
		"spec.version":     lbcf.Spec.Version,
		"status.version":   lbcf.Status.Version,
		"status.phase":     string(lbcf.Status.Phase),
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of LBCF.
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
	newLBCF := obj.(*platform.LBCF)
	oldLBCF := old.(*platform.LBCF)
	newLBCF.Spec = oldLBCF.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
