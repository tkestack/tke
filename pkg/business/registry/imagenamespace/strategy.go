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

package imagenamespace

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
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for ImageNamespace.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	businessClient *businessinternalclient.BusinessClient
	registryClient registryversionedclient.RegistryV1Interface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating ImageNamespace objects.
func NewStrategy(businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) *Strategy {
	return &Strategy{business.Scheme, namesutil.Generator, businessClient, registryClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldImageNamespace := old.(*business.ImageNamespace)
	imageNamespace, _ := obj.(*business.ImageNamespace)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldImageNamespace.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update imageNamespace information",
				log.String("oldTenantID", oldImageNamespace.Spec.TenantID),
				log.String("newTenantID", imageNamespace.Spec.TenantID),
				log.String("userTenantID", tenantID))
		}
		imageNamespace.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is true for ImageNamespaces.
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	imageNamespace, _ := obj.(*business.ImageNamespace)
	if len(tenantID) != 0 {
		imageNamespace.Spec.TenantID = tenantID
	}

	if imageNamespace.Spec.Name != "" {
		imageNamespace.GenerateName = ""
		imageNamespace.Name = imageNamespace.Spec.Name
	} else {
		imageNamespace.GenerateName = "imn-"
	}

	imageNamespace.Spec.Finalizers = []business.FinalizerName{
		business.ImageNamespaceFinalize,
	}
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object) error {
	return nil
}

// Validate validates a new ImageNamespace.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateImageNamespaceCreate(ctx, obj.(*business.ImageNamespace), s.businessClient, s.registryClient)
}

// AllowCreateOnUpdate is false for ImageNamespaces.
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

// ValidateUpdate is the default update validation for an end ImageNamespace.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateImageNamespaceUpdate(ctx, obj.(*business.ImageNamespace), old.(*business.ImageNamespace), s.businessClient, s.registryClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	imageNamespace, ok := obj.(*business.ImageNamespace)
	if !ok {
		return nil, nil, fmt.Errorf("not an ImageNamespace")
	}
	return imageNamespace.Labels, ToSelectableFields(imageNamespace), nil
}

// MatchNamespace returns a generic matcher for a given label and field selector.
func MatchNamespace(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.name",
			"spec.tenantID",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(imageNamespace *business.ImageNamespace) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&imageNamespace.ObjectMeta, true)
	specificFieldsSet := fields.Set{
		"spec.name":     imageNamespace.Spec.Name,
		"spec.tenantID": imageNamespace.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of ImageNamespace.
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
	newImageNamespace := obj.(*business.ImageNamespace)
	oldImageNamespace := old.(*business.ImageNamespace)
	newImageNamespace.Spec = oldImageNamespace.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateImageNamespaceUpdate(ctx, obj.(*business.ImageNamespace), old.(*business.ImageNamespace), s.businessClient, s.registryClient)
}

// FinalizeStrategy implements finalizer logic for ImageNamespace.
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
	newImageNamespace := obj.(*business.ImageNamespace)
	oldImageNamespace := old.(*business.ImageNamespace)
	newImageNamespace.Status = oldImageNamespace.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateImageNamespaceUpdate(ctx, obj.(*business.ImageNamespace), old.(*business.ImageNamespace), s.businessClient, s.registryClient)
}
