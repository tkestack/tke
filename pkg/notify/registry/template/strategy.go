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
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for template.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	notifyClient *notifyinternalclient.NotifyClient
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating template objects.
func NewStrategy(notifyClient *notifyinternalclient.NotifyClient) *Strategy {
	return &Strategy{notify.Scheme, namesutil.Generator, notifyClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldTemplate := old.(*notify.Template)
	template, _ := obj.(*notify.Template)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldTemplate.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update template information", log.String("oldTenantID", oldTemplate.Spec.TenantID), log.String("newTenantID", template.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		template.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for templates.
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	template, _ := obj.(*notify.Template)
	if len(tenantID) != 0 {
		template.Spec.TenantID = tenantID
	}

	if template.Name == "" && template.GenerateName == "" {
		template.GenerateName = "tpl-"
	}
}

// Validate validates a new template.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateTemplate(ctx, obj.(*notify.Template), s.notifyClient)
}

// AllowCreateOnUpdate is false for templates.
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

// ValidateUpdate is the default update validation for an end template.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateTemplateUpdate(ctx, obj.(*notify.Template), old.(*notify.Template), s.notifyClient)
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	template, ok := obj.(*notify.Template)
	if !ok {
		return nil, nil, fmt.Errorf("not a template")
	}
	return template.ObjectMeta.Labels, ToSelectableFields(template), nil
}

// MatchTemplate returns a generic matcher for a given label and field selector.
func MatchTemplate(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(template *notify.Template) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&template.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": template.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
