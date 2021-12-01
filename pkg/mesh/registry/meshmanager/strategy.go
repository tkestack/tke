/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package meshmanager

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
	"tkestack.io/tke/api/mesh"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for meshmanager.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating meshmanager objects.
func NewStrategy() *Strategy {
	return &Strategy{mesh.Scheme, namesutil.Generator}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// NamespaceScoped is false for meshmanagers.
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
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	meshmanager, _ := obj.(*mesh.MeshManager)
	if len(tenantID) != 0 {
		meshmanager.Spec.TenantID = tenantID
	}

	if meshmanager.Name == "" && meshmanager.GenerateName == "" {
		meshmanager.GenerateName = "mm-"
	}

	if meshmanager.Status.Phase == "" {
		meshmanager.Status.Phase = mesh.AddonPhaseInitializing
	}
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldMeshManager := old.(*mesh.MeshManager)
	meshmanager, _ := obj.(*mesh.MeshManager)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldMeshManager.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update meshmanager information", log.String("oldTenantID", oldMeshManager.Spec.TenantID), log.String("newTenantID", meshmanager.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		meshmanager.Spec.TenantID = tenantID
	}
}

// Validate validates a new meshmanager.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateMeshManager(obj.(*mesh.MeshManager))
}

// AllowCreateOnUpdate is false for meshmanagers.
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate returns true if the object can be updated
// unconditionally (irrespective of the latest resource version), when there is
// no resource version specified in the object.
func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end meshmanager.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMeshManagerUpdate(obj.(*mesh.MeshManager), old.(*mesh.MeshManager))
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	meshmanager, ok := obj.(*mesh.MeshManager)
	if !ok {
		return nil, nil, fmt.Errorf("not a meshmanager")
	}
	return labels.Set(meshmanager.ObjectMeta.Labels), ToSelectableFields(meshmanager), nil
}

// MatchMeshManager returns a generic matcher for a given label and field selector.
func MatchMeshManager(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"status.phase",
			"metadata.name",
			"status.version",
			"status.phase",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(meshmanager *mesh.MeshManager) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&meshmanager.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":  meshmanager.Spec.TenantID,
		"status.phase":   string(meshmanager.Status.Phase),
		"status.version": meshmanager.Status.Version,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of MeshManager.
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
	newMeshManager := obj.(*mesh.MeshManager)
	oldMeshManager := old.(*mesh.MeshManager)
	newMeshManager.Spec = oldMeshManager.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMeshManagerUpdate(obj.(*mesh.MeshManager), old.(*mesh.MeshManager))
}

// FinalizeStrategy implements finalizer logic for Machine.
type FinalizeStrategy struct {
	*Strategy
}

var _ rest.RESTUpdateStrategy = &FinalizeStrategy{}

// NewFinalizerStrategy create the FinalizeStrategy object by given strategy.
func NewFinalizerStrategy(strategy *Strategy) *FinalizeStrategy {
	return &FinalizeStrategy{strategy}
}

// PrepareForUpdate is invoked on update before validation to normalize
// the object.  For example: remove fields that are not to be persisted,
// sort order-insensitive list fields, etc.  This should not remove fields
// whose presence would be considered a validation error.
func (FinalizeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newMeshManager := obj.(*mesh.MeshManager)
	oldMeshManager := old.(*mesh.MeshManager)
	newMeshManager.Status = oldMeshManager.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMeshManagerUpdate(obj.(*mesh.MeshManager), old.(*mesh.MeshManager))
}
