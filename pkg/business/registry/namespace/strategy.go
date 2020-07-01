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

package namespace

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
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/util/validation"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for namespace.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	businessClient *businessinternalclient.BusinessClient
	platformClient platformversionedclient.PlatformV1Interface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating namespace objects.
func NewStrategy(businessClient *businessinternalclient.BusinessClient, platformClient platformversionedclient.PlatformV1Interface) *Strategy {
	return &Strategy{business.Scheme, namesutil.Generator, businessClient, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldNamespace := old.(*business.Namespace)
	namespace, _ := obj.(*business.Namespace)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldNamespace.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update namespace information", log.String("oldTenantID", oldNamespace.Spec.TenantID), log.String("newTenantID", namespace.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		namespace.Spec.TenantID = tenantID
	}
	if oldNamespace.Status.CachedSpecHard != nil {
		namespace.Status.CachedSpecHard = oldNamespace.Status.CachedSpecHard
	} else { // For historic data that has no CachedSpecHard
		namespace.Status.CachedSpecHard = oldNamespace.Spec.Hard
	}
}

// NamespaceScoped is false for namespaces.
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
	namespace, _ := obj.(*business.Namespace)
	if len(tenantID) != 0 {
		namespace.Spec.TenantID = tenantID
	}

	if namespace.Spec.ClusterName != "" && namespace.Spec.Namespace != "" {
		namespace.ObjectMeta.GenerateName = ""
		namespace.ObjectMeta.Name = fmt.Sprintf("%s-%s", namespace.Spec.ClusterName, namespace.Spec.Namespace)
	} else {
		namespace.ObjectMeta.GenerateName = "ns"
		namespace.ObjectMeta.Name = ""
	}

	namespace.Spec.Finalizers = []business.FinalizerName{
		business.NamespaceFinalize,
	}
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object) error {
	return nil
}

// Validate validates a new namespace.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateNamespace(ctx, obj.(*business.Namespace), nil,
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}

// AllowCreateOnUpdate is false for namespaces.
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

// ValidateUpdate is the default update validation for an end namespace.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateNamespaceUpdate(ctx, obj.(*business.Namespace), old.(*business.Namespace),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	namespace, ok := obj.(*business.Namespace)
	if !ok {
		return nil, nil, fmt.Errorf("not a namespace")
	}
	return namespace.ObjectMeta.Labels, ToSelectableFields(namespace), nil
}

// MatchNamespace returns a generic matcher for a given label and field selector.
func MatchNamespace(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.clusterName",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(namespace *business.Namespace) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&namespace.ObjectMeta, true)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    namespace.Spec.TenantID,
		"spec.clusterName": namespace.Spec.ClusterName,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Namespace.
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
	newNamespace := obj.(*business.Namespace)
	oldNamespace := old.(*business.Namespace)
	newNamespace.Spec = oldNamespace.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateNamespaceUpdate(ctx, obj.(*business.Namespace), old.(*business.Namespace),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}

// FinalizeStrategy implements verification logic for status of Namespace.
type FinalizeStrategy struct {
	*Strategy
}

var _ rest.RESTUpdateStrategy = &FinalizeStrategy{}

// NewFinalizeStrategy create the FinalizeStrategy object by given strategy.
func NewFinalizeStrategy(strategy *Strategy) *FinalizeStrategy {
	return &FinalizeStrategy{strategy}
}

// PrepareForUpdate is invoked on update before validation to normalize
// the object.  For example: remove fields that are not to be persisted,
// sort order-insensitive list fields, etc.  This should not remove fields
// whose presence would be considered a validation error.
func (FinalizeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newNamespace := obj.(*business.Namespace)
	oldNamespace := old.(*business.Namespace)
	newNamespace.Status = oldNamespace.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateNamespaceUpdate(ctx, obj.(*business.Namespace), old.(*business.Namespace),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}
