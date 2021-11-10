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

package role

import (
	"context"
	"fmt"

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

// Strategy implements verification logic for role.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	enforcer   *casbin.SyncedEnforcer
	authClient authinternalclient.AuthInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating role objects.
func NewStrategy(enforcer *casbin.SyncedEnforcer, authClient authinternalclient.AuthInterface) *Strategy {
	return &Strategy{
		ObjectTyper:   auth.Scheme,
		NameGenerator: namesutil.Generator,
		enforcer:      enforcer,
		authClient:    authClient,
	}
}

// DefaultGarbageCollectionRole returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionRole(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	oldRole := old.(*auth.Role)
	role, _ := obj.(*auth.Role)
	if len(tenantID) != 0 {
		if oldRole.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update role information", log.String("oldTenantID", oldRole.Spec.TenantID), log.String("newTenantID", role.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		role.Spec.TenantID = tenantID
	}

	if len(role.Status.Groups) != 0 {
		role.Status.Groups = util.RemoveDuplicateSubjects(role.Status.Groups)
	} else {
		role.Status.Groups = oldRole.Status.Groups
	}

	if len(role.Status.Users) != 0 {
		role.Status.Users = util.RemoveDuplicateSubjects(role.Status.Users)
	} else {
		role.Status.Users = oldRole.Status.Users
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
	role, _ := obj.(*auth.Role)
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" {
		role.Spec.TenantID = tenantID
	}

	role.Spec.Username = username

	if role.Name == "" && role.GenerateName == "" {
		role.GenerateName = "rol-"
	}

	role.Spec.Finalizers = []auth.FinalizerName{
		auth.RoleFinalize,
	}

	role.Status.Phase = auth.RoleActive

	for i := range role.Status.Groups {
		role.Status.Groups[i].Name = ""
	}

	for i := range role.Status.Users {
		role.Status.Users[i].Name = ""
	}
}

// Validate validates a new role.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateRole(ctx, obj.(*auth.Role), s.authClient)
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

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end role.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateRoleUpdate(ctx, obj.(*auth.Role), old.(*auth.Role), s.authClient)
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	role, ok := obj.(*auth.Role)
	if !ok {
		return nil, nil, fmt.Errorf("not a role")
	}
	return role.ObjectMeta.Labels, ToSelectableFields(role), nil
}

// MatchRole returns a generic matcher for a given label and field selector.
func MatchRole(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
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
func ToSelectableFields(role *auth.Role) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&role.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    role.Spec.TenantID,
		"spec.username":    role.Spec.Username,
		"spec.displayName": role.Spec.DisplayName,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	rol, ok := obj.(*auth.Role)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(rol.Spec.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
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
	newRole := obj.(*auth.Role)
	oldRole := old.(*auth.Role)
	newRole.Spec = oldRole.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateRoleUpdate(ctx, obj.(*auth.Role), old.(*auth.Role), s.authClient)
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
	newRole := obj.(*auth.Role)
	oldRole := old.(*auth.Role)
	newRole.Status = oldRole.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateRoleUpdate(ctx, obj.(*auth.Role), old.(*auth.Role), s.authClient)
}
