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

package receiver

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
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for receiver.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating receiver objects.
func NewStrategy() *Strategy {
	return &Strategy{notify.Scheme, namesutil.Generator}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldReceiver := old.(*notify.Receiver)
	receiver, _ := obj.(*notify.Receiver)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldReceiver.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update receiver information", log.String("oldTenantID", oldReceiver.Spec.TenantID), log.String("newTenantID", receiver.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		receiver.Spec.TenantID = tenantID
	}

	if receiver.Name == "" && receiver.GenerateName == "" {
		receiver.GenerateName = "recv-"
	}
}

// NamespaceScoped is false for receivers.
func (Strategy) NamespaceScoped() bool {
	return false
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	receiver, _ := obj.(*notify.Receiver)
	if len(tenantID) != 0 {
		receiver.Spec.TenantID = tenantID
	}

	receiver.ObjectMeta.Name = ""
	receiver.ObjectMeta.GenerateName = "recv"
}

// Validate validates a new receiver.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateReceiver(obj.(*notify.Receiver))
}

// AllowCreateOnUpdate is false for receivers.
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

// ValidateUpdate is the default update validation for an end receiver.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateReceiverUpdate(obj.(*notify.Receiver), old.(*notify.Receiver))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	receiver, ok := obj.(*notify.Receiver)
	if !ok {
		return nil, nil, fmt.Errorf("not a receiver")
	}
	return labels.Set(receiver.ObjectMeta.Labels), ToSelectableFields(receiver), nil
}

// MatchReceiver returns a generic matcher for a given label and field selector.
func MatchReceiver(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.username",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(receiver *notify.Receiver) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&receiver.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": receiver.Spec.TenantID,
		"spec.username": receiver.Spec.Username,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
