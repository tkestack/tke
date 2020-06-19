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

package projectpolicybinding

import (
	"context"
	"fmt"

	"tkestack.io/tke/pkg/apiserver/filter"

	"github.com/casbin/casbin/v2"
	"k8s.io/apiserver/pkg/registry/generic/registry"

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
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for policy.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	enforcer   *casbin.SyncedEnforcer
	authClient authinternalclient.AuthInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating policy objects.
func NewStrategy(enforcer *casbin.SyncedEnforcer, authClient authinternalclient.AuthInterface) *Strategy {
	return &Strategy{
		ObjectTyper:   auth.Scheme,
		NameGenerator: namesutil.Generator,
		enforcer:      enforcer,
		authClient:    authClient,
	}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	oldBinding, _ := old.(*auth.ProjectPolicyBinding)
	newBinding, _ := obj.(*auth.ProjectPolicyBinding)
	if len(tenantID) != 0 {
		if oldBinding.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update policy information", log.String("oldTenantID", oldBinding.Spec.TenantID), log.String("newTenantID", newBinding.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		newBinding.Spec.TenantID = tenantID
	}

	newBinding.Spec.Groups = util.RemoveDuplicateSubjects(newBinding.Spec.Groups)
	newBinding.Spec.Users = util.RemoveDuplicateSubjectsByIDOrName(newBinding.Spec.Users)
}

// NamespaceScoped is false for policies.
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
	binding, _ := obj.(*auth.ProjectPolicyBinding)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		binding.Spec.TenantID = tenantID
	}

	projectID := filter.ProjectIDFrom(ctx)
	if projectID != "" {
		binding.Spec.ProjectID = projectID
	}

	binding.Status.Phase = auth.BindingActive

	binding.Spec.Finalizers = []auth.FinalizerName{
		auth.BindingFinalize,
	}

	if binding.Spec.ProjectID != "" && binding.Spec.PolicyID != "" {
		binding.Name = util.ProjectPolicyName(binding.Spec.ProjectID, binding.Spec.PolicyID)
	}

	for i := range binding.Spec.Groups {
		binding.Spec.Groups[i].Name = ""
	}

	for i := range binding.Spec.Users {
		binding.Spec.Users[i].Name = ""
	}

	binding.Spec.Groups = util.RemoveDuplicateSubjects(binding.Spec.Groups)
	binding.Spec.Users = util.RemoveDuplicateSubjectsByIDOrName(binding.Spec.Users)
}

// Validate validates a new policy.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateProjectPolicyBinding(ctx, obj.(*auth.ProjectPolicyBinding), s.authClient)
}

// AllowCreateOnUpdate is false for policies.
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

// ValidateUpdate is the default update validation for an end policy.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateProjectPolicyBindingUpdate(ctx, obj.(*auth.ProjectPolicyBinding), old.(*auth.ProjectPolicyBinding), s.authClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	binding, ok := obj.(*auth.ProjectPolicyBinding)
	if !ok {
		return nil, nil, fmt.Errorf("not a projectpolicybinding")
	}
	return binding.ObjectMeta.Labels, ToSelectableFields(binding), nil
}

// MatchPolicy returns a generic matcher for a given label and field selector.
func MatchProjectPolicy(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.projectID",
			"spec.policyID",
			"spec.tenantID",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(binding *auth.ProjectPolicyBinding) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&binding.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.projectID": binding.Spec.ProjectID,
		"spec.policyID":  binding.Spec.PolicyID,
		"spec.tenantID":  binding.Spec.TenantID,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	pol, ok := obj.(*auth.ProjectPolicyBinding)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(pol.Spec.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
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
	newBinding := obj.(*auth.ProjectPolicyBinding)
	oldBinding := old.(*auth.ProjectPolicyBinding)
	newBinding.Spec = oldBinding.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateProjectPolicyBindingUpdate(ctx, obj.(*auth.ProjectPolicyBinding), old.(*auth.ProjectPolicyBinding), s.authClient)
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
	newBinding := obj.(*auth.ProjectPolicyBinding)
	oldBinding := old.(*auth.ProjectPolicyBinding)
	finalizers := newBinding.Spec.Finalizers
	newBinding.Status = oldBinding.Status
	newBinding.Spec = oldBinding.Spec
	newBinding.Spec.Finalizers = finalizers
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
