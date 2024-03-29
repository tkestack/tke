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

package storage

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/auth/registry/custompolicybinding"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for policies and all sub resources.
type Storage struct {
	CustomPolicy *REST

	Status   *StatusREST
	Finalize *FinalizeREST
}

// NewStorage returns a Storage object that will work against policies.
func NewStorage(optsGetter generic.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, privilegedUsername string) *Storage {
	strategy := custompolicybinding.NewStrategy(enforcer, authClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.CustomPolicyBinding{} },
		NewListFunc:              func() runtime.Object { return &auth.CustomPolicyBindingList{} },
		DefaultQualifiedResource: auth.Resource("CustomPolicyBindings"),
		PredicateFunc:            custompolicybinding.MatchCustomPolicy,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		ShouldDeleteDuringUpdate: custompolicybinding.ShouldDeleteDuringUpdate,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    custompolicybinding.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create policy etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = custompolicybinding.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = custompolicybinding.NewFinalizerStrategy(strategy)
	return &Storage{
		CustomPolicy: &REST{store, privilegedUsername},
		Status:       &StatusREST{&statusStore},
		Finalize:     &FinalizeREST{&finalizeStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return CustomPolicyBinding
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*auth.CustomPolicyBinding)
	if err := util.FilterCustomPolicyBinding(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// REST implements a RESTStorage for clusters against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.Creater = &REST{}
var _ rest.ShortNamesProvider = &REST{}
var _ rest.Lister = &REST{}
var _ rest.Getter = &REST{}
var _ rest.Updater = &REST{}
var _ rest.CollectionDeleter = &REST{}
var _ rest.GracefulDeleter = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cpb"}
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
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

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	obj, err := r.Store.List(ctx, wrappedOptions)
	if err != nil {
		return obj, err
	}
	return obj, nil
}

// Delete enforces life-cycle rules for policy termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	object, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	policy := object.(*auth.CustomPolicyBinding)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &policy.UID
	} else if *options.Preconditions.UID != policy.UID {
		err = apierrors.NewConflict(
			auth.Resource("CustomPolicyBindings"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, policy.UID),
		)
		return nil, false, err
	}

	// upon first request to delete, we switch the phase to start policy termination
	if policy.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingBinding, ok := existing.(*auth.CustomPolicyBinding)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *auth.CustomPolicyBinding, got %v", existing)
				}
				if err := deleteValidation(ctx, existingBinding); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingBinding.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingBinding.DeletionTimestamp = &now
				}
				// Set the policy phase to terminating, if needed
				if existingBinding.Status.Phase != auth.BindingTerminating {
					existingBinding.Status.Phase = auth.BindingTerminating
				}

				// the current finalizers which are on namespace
				currentFinalizers := map[string]bool{}
				for _, f := range existingBinding.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on rule
				shouldHaveFinalizers := map[string]bool{
					metav1.FinalizerOrphanDependents: apiserverutil.ShouldHaveOrphanFinalizer(options, currentFinalizers[metav1.FinalizerOrphanDependents]),
					metav1.FinalizerDeleteDependents: apiserverutil.ShouldHaveDeleteDependentsFinalizer(options, currentFinalizers[metav1.FinalizerDeleteDependents]),
				}
				// determine whether there are changes
				changeNeeded := false
				for finalizer, shouldHave := range shouldHaveFinalizers {
					changeNeeded = currentFinalizers[finalizer] != shouldHave || changeNeeded
					if shouldHave {
						currentFinalizers[finalizer] = true
					} else {
						delete(currentFinalizers, finalizer)
					}
				}
				// make the changes if needed
				if changeNeeded {
					var newFinalizers []string
					for f := range currentFinalizers {
						newFinalizers = append(newFinalizers, f)
					}
					existingBinding.Finalizers = newFinalizers
				}
				return existingBinding, nil
			}),
			dryrun.IsDryRun(options.DryRun),
			nil,
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, auth.Resource("CustomPolicyBindings"), name)
			err = storageerr.InterpretUpdateError(err, auth.Resource("CustomPolicyBindings"), name)
			if _, ok := err.(*apierrors.StatusError); !ok {
				err = apierrors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(policy.Spec.Finalizers) != 0 {
		err = apierrors.NewConflict(auth.Resource("CustomPolicyBindings"), policy.Name, fmt.Errorf("the system is ensuring all content is removed from this policy.  Upon completion, this policy will automatically be purged by the system"))
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a
// replication controller.
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

// FinalizeREST implements the REST endpoint for finalizing a policy.
type FinalizeREST struct {
	store *registry.Store
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *FinalizeREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *FinalizeREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.store, name, options)
}

// Update alters the status finalizers subset of an object.
func (r *FinalizeREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
