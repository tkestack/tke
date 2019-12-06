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

package policy

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/util/log"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for policy.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	enforcer   *enforcer.PolicyEnforcer
	authClient authinternalclient.AuthInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating policy objects.
func NewStrategy(enforcer *enforcer.PolicyEnforcer, authClient authinternalclient.AuthInterface) *Strategy {
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
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		oldPolicy := old.(*auth.Policy)
		policy, _ := obj.(*auth.Policy)
		if oldPolicy.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update policy information", log.String("oldTenantID", oldPolicy.Spec.TenantID), log.String("newTenantID", policy.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		policy.Spec.TenantID = tenantID
	}
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
	policy, _ := obj.(*auth.Policy)
	username, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" {
		policy.Spec.TenantID = tenantID
		policy.Spec.Type = auth.PolicyCustom
		policy.Spec.Username = username
	} else {
		if policy.Spec.Type == "" {
			policy.Spec.Type = auth.PolicyCustom
		}

		if policy.Spec.Username == "" {
			policy.Spec.Username = "admin"
		}
	}

	if policy.Name == "" && policy.GenerateName == "" {
		policy.GenerateName = "pol-"
	}

	policy.Status.Phase = auth.PolicyActive

	policy.Spec.Finalizers = []auth.FinalizerName{
		auth.PolicyFinalize,
	}
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object) error {
	policy, _ := obj.(*auth.Policy)
	if err := func() error {
		return s.enforcer.AddPolicy(policy)
	}(); err != nil {
		return fmt.Errorf("failed to create policy '%s', for '%s'", policy.Name, err)
	}

	return nil
}

// Validate validates a new policy.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidatePolicy(obj.(*auth.Policy), s.authClient)
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
	return ValidatePolicyUpdate(obj.(*auth.Policy), old.(*auth.Policy), s.authClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	policy, ok := obj.(*auth.Policy)
	if !ok {
		return nil, nil, fmt.Errorf("not a policy")
	}
	return labels.Set(policy.ObjectMeta.Labels), ToSelectableFields(policy), nil
}

// MatchPolicy returns a generic matcher for a given label and field selector.
func MatchPolicy(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.username",
			"spec.policyName",
			"spec.category",
			"spec.type",
			"spec.displayName",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(policy *auth.Policy) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&policy.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    policy.Spec.TenantID,
		"spec.username":    policy.Spec.Username,
		"spec.category":    policy.Spec.Category,
		"spec.type": string(policy.Spec.Type),
		"spec.policyName": policy.Spec.PolicyName,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
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
	newPolicy := obj.(*auth.Policy)
	oldPolicy := old.(*auth.Policy)
	newPolicy.Spec = oldPolicy.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidatePolicyUpdate(obj.(*auth.Policy), old.(*auth.Policy), s.authClient)
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
	newPolicy := obj.(*auth.Policy)
	oldPolicy := old.(*auth.Policy)
	newPolicy.Status = oldPolicy.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidatePolicyUpdate(obj.(*auth.Policy), old.(*auth.Policy), s.authClient)
}
