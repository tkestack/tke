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

package clustercredential

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for clusterCredential.
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
	clusterCredential, _ := obj.(*platform.ClusterCredential)

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		clusterCredential.TenantID = tenantID
	}

	if clusterCredential.Name == "" && clusterCredential.GenerateName == "" {
		clusterCredential.GenerateName = "cc-"
	}
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		oldClusterCredential := old.(*platform.ClusterCredential)
		clusterCredential, _ := obj.(*platform.ClusterCredential)
		if oldClusterCredential.TenantID != tenantID {
			log.Panic("Unauthorized update clusterCredential information", log.String("oldTenantID", oldClusterCredential.TenantID), log.String("newTenantID", clusterCredential.TenantID), log.String("userTenantID", tenantID))
		}
		clusterCredential.TenantID = tenantID
	}
}

// Validate validates a new clusterCredential.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return Validate(obj.(*platform.ClusterCredential), s.platformClient)
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
	return ValidateUpdate(obj.(*platform.ClusterCredential), old.(*platform.ClusterCredential), s.platformClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	clusterCredential, _ := obj.(*platform.ClusterCredential)
	return labels.Set(clusterCredential.ObjectMeta.Labels), ToSelectableFields(clusterCredential), nil
}

// MatchClusterCredential returns a generic matcher for a given label and field selector.
func MatchClusterCredential(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"tenantID",
			"clusterName",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(clusterCredential *platform.ClusterCredential) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&clusterCredential.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"tenantID":    clusterCredential.TenantID,
		"clusterName": clusterCredential.ClusterName,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
