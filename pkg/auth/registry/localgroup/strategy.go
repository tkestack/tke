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

package localgroup

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
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

// Strategy implements verification logic for group.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	authClient authinternalclient.AuthInterface
	enforcer   *casbin.SyncedEnforcer
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating group objects.
func NewStrategy(authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer) *Strategy {
	return &Strategy{
		ObjectTyper:   auth.Scheme,
		NameGenerator: namesutil.Generator,
		authClient:    authClient,
		enforcer:      enforcer,
	}
}

// DefaultGarbageCollectionGroup returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionGroup(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (s *Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	oldGroup := old.(*auth.LocalGroup)
	group, _ := obj.(*auth.LocalGroup)

	if len(tenantID) != 0 {
		if oldGroup.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update group information", log.String("oldTenantID", oldGroup.Spec.TenantID), log.String("newTenantID", group.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		group.Spec.TenantID = tenantID
	}

	// Update bind users, use binding api
	if len(group.Status.Users) > 0 {
		group.Status.Users = util.RemoveDuplicateSubjects(group.Status.Users)
	} else {
		group.Status.Users = oldGroup.Status.Users
	}

	_ = util.HandleGroupPoliciesUpdate(ctx, s.authClient, s.enforcer, group)
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
	group, _ := obj.(*auth.LocalGroup)
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" {
		group.Spec.TenantID = tenantID
	}

	group.Spec.Username = username

	if group.Name == "" && group.GenerateName == "" {
		group.GenerateName = "grp-"
	}

	group.Spec.Finalizers = []auth.FinalizerName{
		auth.LocalGroupFinalize,
	}

	for i := range group.Status.Users {
		group.Status.Users[i].Name = ""
	}

	group.Status.Phase = auth.GroupActive
}

// Validate validates a new group.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateLocalGroup(ctx, obj.(*auth.LocalGroup), s.authClient)
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

// ValidateUpdate is the default update validation for an end group.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLocalGroupUpdate(ctx, obj.(*auth.LocalGroup), old.(*auth.LocalGroup), s.authClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	group, ok := obj.(*auth.LocalGroup)
	if !ok {
		return nil, nil, fmt.Errorf("not a group")
	}
	return group.ObjectMeta.Labels, ToSelectableFields(group), nil
}

// MatchGroup returns a generic matcher for a given label and field selector.
func MatchGroup(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.username",
			"spec.displayName",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(group *auth.LocalGroup) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&group.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    group.Spec.TenantID,
		"spec.username":    group.Spec.Username,
		"spec.displayName": group.Spec.DisplayName,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	localGroup, ok := obj.(*auth.LocalGroup)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(localGroup.Spec.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
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
	newGroup := obj.(*auth.LocalGroup)
	oldGroup := old.(*auth.LocalGroup)
	newGroup.Spec = oldGroup.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLocalGroupUpdate(ctx, obj.(*auth.LocalGroup), old.(*auth.LocalGroup), s.authClient)
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
	newGroup := obj.(*auth.LocalGroup)
	oldGroup := old.(*auth.LocalGroup)
	newGroup.Status = oldGroup.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateLocalGroupUpdate(ctx, obj.(*auth.LocalGroup), old.(*auth.LocalGroup), s.authClient)
}
