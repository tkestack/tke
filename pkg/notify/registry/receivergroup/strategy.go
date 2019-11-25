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

package receivergroup

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

// Strategy implements verification logic for receiverGroup.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	notifyClient *notifyinternalclient.NotifyClient
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating receiverGroup objects.
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
	oldReceiverGroup := old.(*notify.ReceiverGroup)
	receiverGroup, _ := obj.(*notify.ReceiverGroup)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldReceiverGroup.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update receiverGroup information", log.String("oldTenantID", oldReceiverGroup.Spec.TenantID), log.String("newTenantID", receiverGroup.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		receiverGroup.Spec.TenantID = tenantID
	}

	if receiverGroup.Name == "" && receiverGroup.GenerateName == "" {
		receiverGroup.GenerateName = "rg-"
	}
}

// NamespaceScoped is false for receiverGroups.
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
	receiverGroup, _ := obj.(*notify.ReceiverGroup)
	if len(tenantID) != 0 {
		receiverGroup.Spec.TenantID = tenantID
	}

	receiverGroup.ObjectMeta.Name = ""
	receiverGroup.ObjectMeta.GenerateName = "recvgrp"
}

// Validate validates a new receiverGroup.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateReceiverGroup(obj.(*notify.ReceiverGroup), s.notifyClient)
}

// AllowCreateOnUpdate is false for receiverGroups.
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

// ValidateUpdate is the default update validation for an end receiverGroup.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateReceiverGroupUpdate(obj.(*notify.ReceiverGroup), old.(*notify.ReceiverGroup), s.notifyClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	receiverGroup, ok := obj.(*notify.ReceiverGroup)
	if !ok {
		return nil, nil, fmt.Errorf("not a receiverGroup")
	}
	return labels.Set(receiverGroup.ObjectMeta.Labels), ToSelectableFields(receiverGroup), nil
}

// MatchReceiverGroup returns a generic matcher for a given label and field selector.
func MatchReceiverGroup(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
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
func ToSelectableFields(receiverGroup *notify.ReceiverGroup) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&receiverGroup.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID": receiverGroup.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
