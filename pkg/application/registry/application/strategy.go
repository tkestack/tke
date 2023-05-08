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

package application

import (
	"context"
	"fmt"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/application"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for application.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	applicationClient applicationinternalclient.ApplicationInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating application objects.
func NewStrategy(applicationClient applicationinternalclient.ApplicationInterface) *Strategy {
	return &Strategy{application.Scheme, namesutil.Generator, applicationClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldApp := old.(*application.App)
	app, _ := obj.(*application.App)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldApp.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update app information", log.String("oldTenantID", oldApp.Spec.TenantID), log.String("newTenantID", app.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		app.Spec.TenantID = tenantID
	}

	if app.Spec.Values.RawValues != "" {
		app.Spec.Values.RawValues = helmutil.SafeEncodeValue(app.Spec.Values.RawValues)
	}

	// Any changes to the spec increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object.
	// See metav1.ObjectMeta description for more information on Generation.
	if !apiequality.Semantic.DeepEqual(oldApp.Spec, app.Spec) {
		app.Generation = oldApp.Generation + 1
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
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	app, _ := obj.(*application.App)
	if len(tenantID) != 0 {
		app.Spec.TenantID = tenantID
	}
	app.ObjectMeta.GenerateName = "app-"
	app.Generation = 1

	app.Spec.Finalizers = []application.FinalizerName{
		application.AppFinalize,
	}

	if app.Spec.Chart.TenantID == "" {
		app.Spec.Chart.TenantID = app.Spec.TenantID
	}

	if app.Spec.Values.RawValues != "" {
		app.Spec.Values.RawValues = helmutil.SafeEncodeValue(app.Spec.Values.RawValues)
	}
}

// Validate validates a new app.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateApplication(ctx, obj.(*application.App), s.applicationClient)
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

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(runtime.Object) {
}

// ValidateUpdate is the default update validation for an end application.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateApplicationUpdate(ctx, obj.(*application.App), old.(*application.App))
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	app, ok := obj.(*application.App)
	if !ok {
		return nil, nil, fmt.Errorf("not a application")
	}
	return app.ObjectMeta.Labels, ToSelectableFields(app), nil
}

// MatchApplication returns a generic matcher for a given label and field selector.
func MatchApplication(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.name",
			"spec.type",
			"spec.targetCluster",
			"spec.targetNamespace",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(app *application.App) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&app.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":        app.Spec.TenantID,
		"spec.name":            app.Spec.Name,
		"spec.type":            string(app.Spec.Type),
		"spec.targetCluster":   app.Spec.TargetCluster,
		"spec.targetNamespace": app.Spec.TargetNamespace,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	app, ok := obj.(*application.App)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(app.Spec.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
}

// StatusStrategy implements verification logic for status of app request.
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
	newApplication := obj.(*application.App)
	oldApplication := old.(*application.App)
	newApplication.Spec = oldApplication.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateApplicationUpdate(ctx, obj.(*application.App), old.(*application.App))
}

// FinalizeStrategy implements finalizer logic for App.
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
	newApplication := obj.(*application.App)
	oldApplication := old.(*application.App)
	newApplication.Status = oldApplication.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateApplicationUpdate(ctx, obj.(*application.App), old.(*application.App))
}
