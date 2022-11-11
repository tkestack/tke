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

package message

import (
	"context"
	"fmt"
	"math"
	"time"

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

// Strategy implements verification logic for message.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating message objects.
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
	oldMessage := old.(*notify.Message)
	message, _ := obj.(*notify.Message)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldMessage.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update message information", log.String("oldTenantID", oldMessage.Spec.TenantID), log.String("newTenantID", message.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		message.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for messages.
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
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	message, _ := obj.(*notify.Message)
	if len(tenantID) != 0 {
		message.Spec.TenantID = tenantID
	}

	// use (math.MaxUint32 - time.Now().Unix()) as prefix of name will support
	// to list message in descending order of creationtimestamp
	if message.Name == "" && message.GenerateName == "" {
		message.GenerateName = fmt.Sprintf("msg-%d-", (math.MaxUint32 - time.Now().Unix()))
	}
}

// Validate validates a new message.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateMessage(obj.(*notify.Message))
}

// AllowCreateOnUpdate is false for messages.
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

// ValidateUpdate is the default update validation for an end message.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMessageUpdate(obj.(*notify.Message), old.(*notify.Message))
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	message, ok := obj.(*notify.Message)
	if !ok {
		return nil, nil, fmt.Errorf("not a message")
	}
	return labels.Set(message.ObjectMeta.Labels), ToSelectableFields(message), nil
}

// MatchMessage returns a generic matcher for a given label and field selector.
func MatchMessage(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.receiverName",
			"spec.username",
			"spec.channelMessageID",
			"status.phase",
			"metadata.name",
			"spec.alarmPolicyName",
			"spec.alarmPolicyType",
			"spec.receiverChannelName",
			"spec.clusterID",
			"status.alertStatus",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(message *notify.Message) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&message.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":            message.Spec.TenantID,
		"spec.receiverName":        message.Spec.ReceiverName,
		"spec.username":            message.Spec.Username,
		"spec.channelMessageID":    message.Spec.ChannelMessageID,
		"status.phase":             string(message.Status.Phase),
		"spec.alarmPolicyName":     message.Spec.AlarmPolicyName,
		"spec.alarmPolicyType":     message.Spec.AlarmPolicyType,
		"spec.receiverChannelName": message.Spec.ReceiverChannelName,
		"spec.clusterID":           message.Spec.ClusterID,
		"status.alertStatus":       message.Status.AlertStatus,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of message request.
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
	newMessage := obj.(*notify.Message)
	oldMessage := old.(*notify.Message)
	newMessage.Spec = oldMessage.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMessageUpdate(obj.(*notify.Message), old.(*notify.Message))
}
