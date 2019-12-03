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

package chart

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
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for chart.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	registryClient *registryinternalclient.RegistryClient
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating chart objects.
func NewStrategy(registryClient *registryinternalclient.RegistryClient) *Strategy {
	return &Strategy{registry.Scheme, namesutil.Generator, registryClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldChart := old.(*registry.Chart)
	chart, _ := obj.(*registry.Chart)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldChart.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update chart information", log.String("oldTenantID", oldChart.Spec.TenantID), log.String("newTenantID", chart.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		chart.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for repositories.
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(context.Context, runtime.Object, bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	chart, _ := obj.(*registry.Chart)
	if len(tenantID) != 0 {
		chart.Spec.TenantID = tenantID
	}
	chart.ObjectMeta.GenerateName = "repo-"
	chart.ObjectMeta.Name = ""
}

// Validate validates a new chart.
func (s *Strategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
	return ValidateChart(obj.(*registry.Chart), s.registryClient)
}

// AllowCreateOnUpdate is false for repositories.
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
func (Strategy) Canonicalize(runtime.Object) {
}

// ValidateUpdate is the default update validation for an end chart.
func (s *Strategy) ValidateUpdate(_ context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateChartUpdate(obj.(*registry.Chart), old.(*registry.Chart))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	chart, ok := obj.(*registry.Chart)
	if !ok {
		return nil, nil, fmt.Errorf("not a chart")
	}
	return chart.ObjectMeta.Labels, ToSelectableFields(chart), nil
}

// MatchChart returns a generic matcher for a given label and field selector.
func MatchChart(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.name",
			"spec.chartGroupName",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(chart *registry.Chart) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&chart.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":       chart.Spec.TenantID,
		"spec.chartGroupName": chart.Spec.ChartGroupName,
		"spec.name":           chart.Spec.Name,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of chart request.
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
	newChart := obj.(*registry.Chart)
	oldChart := old.(*registry.Chart)
	newChart.Spec = oldChart.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(_ context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateChartUpdate(obj.(*registry.Chart), old.(*registry.Chart))
}
