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

package storage

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	messagestrategy "tkestack.io/tke/pkg/notify/registry/message"
	"tkestack.io/tke/pkg/notify/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for messages and all sub resources.
type Storage struct {
	Message *REST
	Status  *StatusREST
}

// NewStorage returns a Storage object that will work against messages.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, privilegedUsername string, messageTTL time.Duration) *Storage {
	strategy := messagestrategy.NewStrategy()
	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &notify.Message{} },
		NewListFunc: func() runtime.Object { return &notify.MessageList{} },
		TTLFunc: func(obj runtime.Object, existing uint64, update bool) (uint64, error) {
			return uint64(messageTTL.Seconds()), nil
		},
		DefaultQualifiedResource: notify.Resource("messages"),
		PredicateFunc:            messagestrategy.MatchMessage,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    messagestrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create message etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = messagestrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = messagestrategy.NewStatusStrategy(strategy)

	return &Storage{
		Message: &REST{store, privilegedUsername},
		Status:  &StatusREST{&statusStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return Message
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	message := obj.(*notify.Message)
	if err := util.FilterMessage(ctx, message); err != nil {
		return nil, err
	}
	return message, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return Message
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	message := obj.(*notify.Message)
	if err := util.FilterMessage(ctx, message); err != nil {
		return nil, err
	}
	return message, nil
}

// REST implements a RESTStorage for messages against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"msg"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, errors.NewMethodNotSupported(notify.Resource("messages"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, messageName string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, messageName, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *REST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, name, options)
}

// Update finds a resource in the storage and updates it.
func (r *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// Delete enforces life-cycle rules for cluster termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a message request.
type StatusREST struct {
	store *registry.Store
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *StatusREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.store, name, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *StatusREST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.store, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
