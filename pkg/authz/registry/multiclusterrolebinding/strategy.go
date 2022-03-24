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

package multiclusterrolebinding

import (
	"context"
	"encoding/json"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/client-go/tools/cache"
	"tkestack.io/tke/api/authz"
	"tkestack.io/tke/pkg/authz/constant"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for configmap.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateStrategy = &Strategy{}
var _ rest.RESTUpdateStrategy = &Strategy{}
var _ rest.RESTDeleteStrategy = &Strategy{}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating namespace set objects.
func NewStrategy() *Strategy {
	return &Strategy{authz.Scheme, namesutil.Generator}
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
	mcrb, _ := obj.(*authz.MultiClusterRoleBinding)
	roleNs, roleName, err := cache.SplitMetaNamespaceKey(mcrb.Spec.RoleName)
	if err != nil {
		return
	}
	labels := mcrb.Labels
	if labels == nil {
		labels = map[string]string{}
	}
	labels[constant.RoleNamespace] = roleNs
	labels[constant.RoleName] = roleName
	labels[constant.Username] = mcrb.Spec.Username
	if dispatchAllClusters(mcrb.Spec.Clusters) {
		labels[constant.DispatchAllClusters] = "true"
	}
	mcrb.Labels = labels

	annotation := mcrb.Annotations
	if annotation == nil {
		annotation = map[string]string{}
	}
	bytes, _ := json.Marshal(mcrb.Spec.Clusters)
	annotation[constant.LastDispatchedClusters] = string(bytes)
	mcrb.Annotations = annotation

	mcrb.Status.Phase = authz.BindingActive
	mcrb.ObjectMeta.Finalizers = []string{string(authz.MultiClusterRoleBindingFinalize)}
}

func dispatchAllClusters(clusterIDs []string) bool {
	if len(clusterIDs) == 1 && clusterIDs[0] == "*" {
		return true
	}
	return false
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	mcrb := obj.(*authz.MultiClusterRoleBinding)
	if dispatchAllClusters(mcrb.Spec.Clusters) {
		mcrb.Labels[constant.DispatchAllClusters] = "true"
	} else {
		delete(mcrb.Labels, constant.DispatchAllClusters)
	}
}

// Validate validates a new configmap.
// TODO
func (Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateMultiClusterRoleBinding(obj.(*authz.MultiClusterRoleBinding))
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
func (Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMultiClusterRoleBindingUpdate(obj.(*authz.MultiClusterRoleBinding), old.(*authz.MultiClusterRoleBinding))
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	pol, ok := obj.(*authz.MultiClusterRoleBinding)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(pol.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
}

// StatusStrategy implements verification logic for status of roletemplate request.
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
func (StatusStrategy) PrepareForUpdate(_ context.Context, obj, old runtime.Object) {
	newMultiClusterRoleBinding := obj.(*authz.MultiClusterRoleBinding)
	oldMultiClusterRoleBinding := old.(*authz.MultiClusterRoleBinding)
	status := newMultiClusterRoleBinding.Status
	newMultiClusterRoleBinding = oldMultiClusterRoleBinding
	newMultiClusterRoleBinding.Status = status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMultiClusterRoleBindingUpdate(obj.(*authz.MultiClusterRoleBinding), old.(*authz.MultiClusterRoleBinding))
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
	newBinding := obj.(*authz.MultiClusterRoleBinding)
	oldBinding := old.(*authz.MultiClusterRoleBinding)
	finalizers := newBinding.Finalizers
	newBinding = oldBinding
	newBinding.Finalizers = finalizers
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
