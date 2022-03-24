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

package role

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/authz"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	authzprovider "tkestack.io/tke/pkg/authz/provider"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for configmap.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	policyGetter   rest.Getter
	platformClient platformversionedclient.PlatformV1Interface
}

const NamePrefix = "rol-"

var _ rest.RESTCreateStrategy = &Strategy{}
var _ rest.RESTUpdateStrategy = &Strategy{}
var _ rest.RESTDeleteStrategy = &Strategy{}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	pol, ok := obj.(*authz.Role)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(pol.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating namespace set objects.
func NewStrategy(policyGetter rest.Getter, platformClient platformversionedclient.PlatformV1Interface) *Strategy {
	return &Strategy{authz.Scheme, namesutil.Generator, policyGetter, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// NamespaceScoped is false for namespaceSets
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	username, _ := authentication.UsernameAndTenantID(ctx)
	tenantID := request.NamespaceValue(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	role, _ := obj.(*authz.Role)
	role.TenantID = tenantID
	if role.Username == "" {
		role.Username = username
	}
	if role.Name == "" && role.GenerateName == "" {
		role.GenerateName = NamePrefix
	}
	region := authentication.GetExtraValue("region", ctx)
	log.Debugf("region '%v'", region)
	if len(region) != 0 {
		annotations := role.Annotations
		if len(annotations) == 0 {
			annotations = map[string]string{}
		}
		annotations[authz.GroupName+"/region"] = region[0]
		role.Annotations = annotations
	}
	role.Finalizers = []string{string(authz.RoleFinalize)}
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldRole := old.(*authz.Role)
	role, _ := obj.(*authz.Role)
	if role.TenantID != oldRole.TenantID {
		log.Warnf("Unauthorized update role tenantID '%s'", oldRole.TenantID)
		role.TenantID = oldRole.TenantID
	}
}

// Validate validates a new configmap.
func (s Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	role := obj.(*authz.Role)
	provider, err := authzprovider.GetProvider(role.Annotations)
	if err == nil {
		if fieldErr := provider.Validate(context.TODO(), role, s.platformClient); fieldErr != nil {
			return field.ErrorList{fieldErr}
		}
	}
	return ValidateRole(role, s.policyGetter, s.platformClient)
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

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end namespace set.
func (s Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateRoleUpdate(ctx, obj.(*authz.Role), old.(*authz.Role), s.policyGetter, s.platformClient)
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
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
	newRole := obj.(*authz.Role)
	oldRole := old.(*authz.Role)
	finalizers := newRole.Finalizers
	newRole = oldRole
	newRole.Finalizers = finalizers
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
