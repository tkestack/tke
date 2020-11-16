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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	registryapi "tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	harbor "tkestack.io/tke/pkg/registry/harbor/client"
	harborHandler "tkestack.io/tke/pkg/registry/harbor/handler"
	namespacestrategy "tkestack.io/tke/pkg/registry/registry/namespace"
	registryutil "tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for namespaces and all sub resources.
type Storage struct {
	Namespace *REST
	Status    *StatusREST
}

// NewStorage returns a Storage object that will work against namespaces.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, registryClient *registryinternalclient.RegistryClient, privilegedUsername string, harborClient *harbor.APIClient) *Storage {
	strategy := namespacestrategy.NewStrategy(registryClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &registryapi.Namespace{} },
		NewListFunc:              func() runtime.Object { return &registryapi.NamespaceList{} },
		DefaultQualifiedResource: registryapi.Resource("namespaces"),
		PredicateFunc:            namespacestrategy.MatchNamespace,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    namespacestrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create namespace etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = namespacestrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = namespacestrategy.NewStatusStrategy(strategy)

	return &Storage{
		Namespace: &REST{store, privilegedUsername, harborClient},
		Status:    &StatusREST{&statusStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return Namespace
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*registryapi.Namespace)
	if err := registryutil.FilterNamespace(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return Namespace
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*registryapi.Namespace)
	if err := registryutil.FilterNamespace(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// REST implements a RESTStorage for namespaces against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
	harborClient       *harbor.APIClient
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"rns"}
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
		return nil, apierrors.NewMethodNotSupported(registryapi.Resource("namespaces"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {

	if r.harborClient != nil {
		o := obj.(*registryapi.Namespace)
		_, tenantID := authentication.UsernameAndTenantID(ctx)

		err := harborHandler.CreateProject(
			ctx,
			r.harborClient,
			fmt.Sprintf("%s-image-%s", tenantID, o.Spec.Name),
			o.Spec.Visibility == registryapi.VisibilityPublic,
		)
		if err != nil {
			return nil, err
		}
	}

	obj, err := r.Store.Create(ctx, obj, createValidation, options)
	if err != nil {
		if r.harborClient != nil {
			o := obj.(*registryapi.Namespace)
			_, tenantID := authentication.UsernameAndTenantID(ctx)
			// cleanup harbor project
			harborHandler.DeleteProject(ctx, r.harborClient, fmt.Sprintf("%s-image-%s", tenantID, o.Spec.Name))
		}
		return nil, err
	}

	return obj, nil
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *REST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, name, options)
}

// Update alters the object subset of an object.
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
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	if r.harborClient != nil {
		o := obj.(*registryapi.Namespace)
		// delete harbor project
		harborHandler.DeleteProject(ctx, r.harborClient, fmt.Sprintf("%s-image-%s", o.Spec.TenantID, o.Spec.Name))
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a namespace.
type StatusREST struct {
	store *registry.Store
}

// StatusREST implements Patcher.
var _ = rest.Patcher(&StatusREST{})

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
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
