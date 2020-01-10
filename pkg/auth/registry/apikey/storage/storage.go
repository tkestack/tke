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

	"k8s.io/apimachinery/pkg/fields"

	"tkestack.io/tke/pkg/apiserver/authentication"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/auth/registry/apikey"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// Storage includes storage for identities and all sub resources.
type Storage struct {
	APIKey   *REST
	Password *PasswordREST
	Token    *TokenREST
	Status   *StatusREST
}

// NewStorage returns a Storage object that will work against identify.
func NewStorage(optsGetter generic.RESTOptionsGetter, authClient authinternalclient.AuthInterface, keySigner util.KeySigner, privilegedUsername string) *Storage {
	strategy := apikey.NewStrategy(keySigner, privilegedUsername)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.APIKey{} },
		NewListFunc:              func() runtime.Object { return &auth.APIKeyList{} },
		DefaultQualifiedResource: auth.Resource("apikeys"),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
		ExportStrategy:           strategy,
		Decorator:                apikey.Decorator,

		PredicateFunc: apikey.MatchAPIKey,
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    apikey.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create local identity etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = apikey.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = apikey.NewStatusStrategy(strategy)

	return &Storage{
		APIKey: &REST{store, keySigner, privilegedUsername},
		Password: &PasswordREST{
			apiKeyStore: store,
			keySigner:   keySigner,
			authClient:  authClient,
		},
		Token: &TokenREST{
			apiKeyStore: store,
			keySigner:   keySigner,
		},
		Status: &StatusREST{&statusStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return apiKey
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*auth.APIKey)
	if err := util.FilterAPIKey(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return apiKey
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*auth.APIKey)
	if err := util.FilterAPIKey(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

// ValidateListObject validate if list by admin, if false, filter deleted apikey.
func ValidateListObjectAndTenantID(ctx context.Context, store *registry.Store, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)

	username, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" {
		wrappedOptions.FieldSelector = fields.AndSelectors(wrappedOptions.FieldSelector, fields.OneTermEqualSelector("spec.username", username))
	}

	obj, err := store.List(ctx, wrappedOptions)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// REST implements a RESTStorage for identities against etcd.
type REST struct {
	*registry.Store

	keySigner          util.KeySigner
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"apk"}
}

func (r *REST) New() runtime.Object {
	return &auth.APIKey{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	return ValidateListObjectAndTenantID(ctx, r.Store, options)
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, apierrors.NewMethodNotSupported(auth.Resource("apiKeys"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
}

// Export an object. Fields that are not user specified are stripped out
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

// Delete enforces life-cycle rules for api key termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	userName, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return r.Store.Delete(ctx, name, deleteValidation, options)
	}

	apiKey := obj.(*auth.APIKey)
	tokenInfo, err := r.keySigner.Verify(apiKey.Spec.APIkey)
	if tokenInfo != nil {
		if tokenInfo.UserName != userName {
			return nil, false, apierrors.NewForbidden(auth.Resource("apiKeys"), name, fmt.Errorf("forbid to delete"))
		}
		return r.Store.Delete(ctx, name, deleteValidation, options)
	}

	return nil, false, apierrors.NewInternalError(err)
}
