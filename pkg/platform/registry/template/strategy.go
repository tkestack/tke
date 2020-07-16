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

package template

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for persistent event.
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
	return &Strategy{platform.Scheme, namesutil.Generator}
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
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	template, _ := obj.(*platform.Template)

	if len(tenantID) != 0 {
		template.Spec.TenantID = tenantID
	}

	if template.Name == "" && template.GenerateName == "" {
		template.GenerateName = "tl-"
	}

	if template.Spec.Content.Deployment != nil {
		template.Spec.Type = Deployment
	}
	if template.Spec.Content.StatefulSet != nil {
		template.Spec.Type = StatefulSet
	}
	if template.Spec.Content.DaemonSet != nil {
		template.Spec.Type = DaemonSet
	}
	if template.Spec.Content.Job != nil {
		template.Spec.Type = Job
	}
	if template.Spec.Content.CronJob != nil {
		template.Spec.Type = CronJob
	}
	if template.Spec.Content.Tapp != nil {
		template.Spec.Type = TApp
	}

}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		oldTemplate := old.(*platform.Template)
		template, _ := obj.(*platform.Template)
		if oldTemplate.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update persistent event information", log.String("oldTenantID", oldTemplate.Spec.TenantID), log.String("newTenantID", template.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		template.Spec.TenantID = tenantID
	}
}

// Validate validates a new persistent event.
func (Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateTemplate(obj.(*platform.Template))
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
func (Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateTemplateUpdate(obj.(*platform.Template), old.(*platform.Template))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	template, _ := obj.(*platform.Template)
	return labels.Set(template.ObjectMeta.Labels), ToSelectableFields(template), nil
}

// MatchTemplate returns a generic matcher for a given label and field selector.
func MatchTemplate(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.type",
			"spec.username",
			"status.phase"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(template *platform.Template) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&template.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": template.Spec.TenantID,
		"spec.type":     template.Spec.Type,
		"spec.username": template.Spec.Username,
		"status.phase":  string(template.Status.Phase),
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Template.
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
	newTemplate := obj.(*platform.Template)
	oldTemplate := old.(*platform.Template)
	newTemplate.Spec = oldTemplate.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
