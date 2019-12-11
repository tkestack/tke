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

package localidentity

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"

	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for oidc identity.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	authClient authinternalclient.AuthInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating identity objects.
func NewStrategy(authClient authinternalclient.AuthInterface) *Strategy {
	return &Strategy{auth.Scheme, namesutil.Generator, authClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldLocalIdentity := old.(*auth.LocalIdentity)
	localIdentity, _ := obj.(*auth.LocalIdentity)

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldLocalIdentity.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update local identity information", log.String("oldTenantID", oldLocalIdentity.Spec.TenantID), log.String("newTenantID", localIdentity.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		localIdentity.Spec.TenantID = tenantID
	}

	localIdentity.Status.LastUpdateTime = metav1.Now()
}

// NamespaceScoped is false for identities.
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
	localIdentity, _ := obj.(*auth.LocalIdentity)

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		localIdentity.Spec.TenantID = tenantID
	}
	if localIdentity.Name == "" && localIdentity.GenerateName == "" {
		localIdentity.GenerateName = "usr-"
	}

	localIdentity.Spec.Finalizers = []auth.FinalizerName{
		auth.LocalIdentityFinalize,
	}
}

// Validate validates a new identity.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateLocalIdentity(s.authClient, obj.(*auth.LocalIdentity), false)
}

// AllowCreateOnUpdate is false for identities.
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

// ValidateUpdate is the default update validation for an identity.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLocalIdentityUpdate(s.authClient, obj.(*auth.LocalIdentity), old.(*auth.LocalIdentity))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	localIdentity, ok := obj.(*auth.LocalIdentity)
	if !ok {
		return nil, nil, fmt.Errorf("not a localIdentity")
	}
	return labels.Set(localIdentity.ObjectMeta.Labels), ToSelectableFields(localIdentity), nil
}

// MatchLocalIdentity returns a generic matcher for a given label and field selector.
func MatchLocalIdentity(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.username",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(localIdentity *auth.LocalIdentity) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&localIdentity.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": localIdentity.Spec.TenantID,
		"spec.username": localIdentity.Spec.Username,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Registry.
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
	newLocalIdentity := obj.(*auth.LocalIdentity)
	oldLocalIdentity := old.(*auth.LocalIdentity)
	newLocalIdentity.Spec = oldLocalIdentity.Spec
	newLocalIdentity.Status.LastUpdateTime = metav1.Now()
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
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
	newPolicy := obj.(*auth.LocalIdentity)
	oldPolicy := old.(*auth.LocalIdentity)
	newPolicy.Status = oldPolicy.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLocalIdentityUpdate(s.authClient, obj.(*auth.LocalIdentity), old.(*auth.LocalIdentity))
}
