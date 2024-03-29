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

package messagerequest

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Strategy implements verification logic for message request.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	notifyClient *notifyinternalclient.NotifyClient
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating message request objects.
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
	oldMessageRequest := old.(*notify.MessageRequest)
	messageRequest, _ := obj.(*notify.MessageRequest)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldMessageRequest.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update message request information", log.String("oldTenantID", oldMessageRequest.Spec.TenantID), log.String("newTenantID", messageRequest.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		messageRequest.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for message requests.
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
	messageRequest, _ := obj.(*notify.MessageRequest)
	if len(tenantID) != 0 {
		messageRequest.Spec.TenantID = tenantID
	} else {
		channel, err := s.notifyClient.Channels().Get(ctx, messageRequest.ObjectMeta.Namespace, metav1.GetOptions{})
		if err == nil && channel != nil {
			messageRequest.Spec.TenantID = channel.Spec.TenantID
		} else {
			log.Panic("Cannot polyfill message request tenant", log.String("messageRequestName", messageRequest.ObjectMeta.Name), log.Err(err))
		}
	}

	if messageRequest.Name == "" && messageRequest.GenerateName == "" {
		messageRequest.GenerateName = "mr-"
	}
}

// Validate validates a new message request.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateMessageRequest(ctx, obj.(*notify.MessageRequest), s.notifyClient)
}

// AllowCreateOnUpdate is false for message requests.
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

// ValidateUpdate is the default update validation for an end message request.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMessageRequestUpdate(ctx, obj.(*notify.MessageRequest), old.(*notify.MessageRequest), s.notifyClient)
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	messageRequest, ok := obj.(*notify.MessageRequest)
	if !ok {
		return nil, nil, fmt.Errorf("not a message request")
	}
	return messageRequest.ObjectMeta.Labels, ToSelectableFields(messageRequest), nil
}

// MatchMessageRequest returns a generic matcher for a given label and field selector.
func MatchMessageRequest(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"status.phase",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(messageRequest *notify.MessageRequest) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&messageRequest.ObjectMeta, true)
	specificFieldsSet := fields.Set{
		"spec.tenantID": messageRequest.Spec.TenantID,
		"status.phase":  string(messageRequest.Status.Phase),
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
	newMessageRequest := obj.(*notify.MessageRequest)
	oldMessageRequest := old.(*notify.MessageRequest)
	newMessageRequest.Spec = oldMessageRequest.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateMessageRequestUpdate(ctx, obj.(*notify.MessageRequest), old.(*notify.MessageRequest), s.notifyClient)
}
