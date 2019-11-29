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

package rule


import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"

	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"

	namesutil "tkestack.io/tke/pkg/util/names"

)

// Strategy implements verification logic for rule.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating rule objects.
func NewStrategy() *Strategy {
	return &Strategy{auth.Scheme, namesutil.Generator}
}

// DefaultGarbageCollectionRule returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionRule(ctx context.Context) rest.GarbageCollectionRule {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		oldRule := old.(*auth.Rule)
		rule, _ := obj.(*auth.Rule)
		if oldRule.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update rule information", log.String("oldTenantID", oldRule.Spec.TenantID), log.String("newTenantID", rule.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		rule.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for policies.
func (Strategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	rule, _ := obj.(*auth.Rule)
	username, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		rule.Spec.TenantID = tenantID
	}

	if rule.Spec.Username == "" {
		rule.Spec.Username = username
	}

	if rule.Name == "" && rule.GenerateName == "" {
		rule.GenerateName = "rul-"
	}
}

// Validate validates a new rule.
func (Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateRule(obj.(*auth.Rule))
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

// ValidateUpdate is the default update validation for an end rule.
func (Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateRuleUpdate(obj.(*auth.Rule), old.(*auth.Rule))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	rule, ok := obj.(*auth.Rule)
	if !ok {
		return nil, nil, fmt.Errorf("not a rule")
	}
	return labels.Set(rule.ObjectMeta.Labels), ToSelectableFields(rule), nil
}

// MatchRule returns a generic matcher for a given label and field selector.
func MatchRule(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.username",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(rule *auth.Rule) fields.Set {
	objectMetaFieldsSet := generic.ObjectMetaFieldsSet(&rule.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": rule.Spec.TenantID,
		"spec.username": rule.Spec.Username,
	}
	return generic.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

