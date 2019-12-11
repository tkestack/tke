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
	"sync"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/util/dryrun"
	"k8s.io/klog"
	"tkestack.io/tke/pkg/auth/util"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for identities and all sub resources.
type Storage struct {
	LocalIdentity *REST
	Password      *PasswordREST
	Status        *StatusREST
	Policy        *PolicyREST
	Role          *RoleREST
	Group         *GroupREST
	Finalize      *FinalizeREST
}

// NewStorage returns a Storage object that will work against identify.
func NewStorage(optsGetter generic.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, privilegedUsername string) *Storage {
	strategy := localidentity.NewStrategy(authClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.LocalIdentity{} },
		NewListFunc:              func() runtime.Object { return &auth.LocalIdentityList{} },
		DefaultQualifiedResource: auth.Resource("localidentities"),
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,

		PredicateFunc: localidentity.MatchLocalIdentity,
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    localidentity.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create local identity etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = localidentity.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = localidentity.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = localidentity.NewFinalizerStrategy(strategy)
	finalizeStore.ExportStrategy = localidentity.NewFinalizerStrategy(strategy)

	return &Storage{
		LocalIdentity: &REST{store, privilegedUsername},
		Password:      &PasswordREST{store, authClient},
		Status:        &StatusREST{&statusStore},
		Policy:        &PolicyREST{store, authClient, enforcer},
		Role:          &RoleREST{store, authClient, enforcer},
		Group:         &GroupREST{store, authClient, enforcer},
		Finalize:      &FinalizeREST{&finalizeStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return Project
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*auth.LocalIdentity)
	if err := util.FilterLocalIdentity(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return Project
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*auth.LocalIdentity)
	if err := util.FilterLocalIdentity(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

// ValidateListObject validate if list by admin, if true not return hashed password.
func ValidateListObjectAndTenantID(ctx context.Context, store *registry.Store, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	obj, err := store.List(ctx, wrappedOptions)
	if err != nil {
		return obj, err
	}

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return obj, err
	}

	identityList := obj.(*auth.LocalIdentityList)
	for i := range identityList.Items {
		identityList.Items[i].Spec.HashedPassword = ""
	}

	return identityList, nil
}

// REST implements a RESTStorage for identities against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"usr"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	return ValidateListObjectAndTenantID(ctx, r.Store, options)
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, apierrors.NewMethodNotSupported(auth.Resource("localIdentities"), "delete collection")
	}

	if listOptions == nil {
		listOptions = &metainternal.ListOptions{}
	} else {
		listOptions = listOptions.DeepCopy()
	}

	listObj, err := r.Store.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	items, err := meta.ExtractList(listObj)
	if err != nil {
		return nil, err
	}

	// Spawn a number of goroutines, so that we can issue requests to storage
	// in parallel to speed up deletion.
	// TODO: Make this proportional to the number of items to delete, up to
	// DeleteCollectionWorkers (it doesn't make much sense to spawn 16
	// workers to delete 10 items).
	workersNumber := r.Store.DeleteCollectionWorkers
	if workersNumber < 1 {
		workersNumber = 1
	}
	wg := sync.WaitGroup{}
	toProcess := make(chan int, 2*workersNumber)
	errs := make(chan error, workersNumber+1)

	go func() {
		defer utilruntime.HandleCrash(func(panicReason interface{}) {
			errs <- fmt.Errorf("DeleteCollection distributor panicked: %v", panicReason)
		})
		for i := 0; i < len(items); i++ {
			toProcess <- i
		}
		close(toProcess)
	}()

	wg.Add(workersNumber)
	for i := 0; i < workersNumber; i++ {
		go func() {
			// panics don't cross goroutine boundaries
			defer utilruntime.HandleCrash(func(panicReason interface{}) {
				errs <- fmt.Errorf("DeleteCollection goroutine panicked: %v", panicReason)
			})
			defer wg.Done()

			for index := range toProcess {
				accessor, err := meta.Accessor(items[index])
				if err != nil {
					errs <- err
					return
				}

				tmpOpt := options
				tmpOpt.Preconditions = nil

				if _, _, err := r.Delete(ctx, accessor.GetName(), deleteValidation, tmpOpt); err != nil && !apierrors.IsNotFound(err) {
					klog.V(4).Infof("Delete %s in DeleteCollection failed: %v", accessor.GetName(), err)
					errs <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	select {
	case err := <-errs:
		return nil, err
	default:
		return listObj, nil
	}
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
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

// Delete enforces life-cycle rules for localIdentity termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	object, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	localIdentity := object.(*auth.LocalIdentity)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &localIdentity.UID
	} else if *options.Preconditions.UID != localIdentity.UID {
		err = apierrors.NewConflict(
			auth.Resource("localidentities"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, localIdentity.UID),
		)
		return nil, false, err
	}

	// upon first request to delete, we switch the phase to start policy termination
	if localIdentity.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingLocalIdentity, ok := existing.(*auth.LocalIdentity)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *auth.LocalIdentity, got %v", existing)
				}
				if err := deleteValidation(ctx, existingLocalIdentity); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingLocalIdentity.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingLocalIdentity.DeletionTimestamp = &now
				}
				// Set the policy phase to terminating, if needed
				if existingLocalIdentity.Status.Phase != auth.LocalIdentityDeleting {
					existingLocalIdentity.Status.Phase = auth.LocalIdentityDeleting
				}

				// the current finalizers which are on namespace
				currentFinalizers := map[string]bool{}
				for _, f := range existingLocalIdentity.Finalizers {
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
					existingLocalIdentity.Finalizers = newFinalizers
				}
				return existingLocalIdentity, nil
			}),
			dryrun.IsDryRun(options.DryRun),
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, auth.Resource("localidentities"), name)
			err = storageerr.InterpretUpdateError(err, auth.Resource("localidentities"), name)
			if _, ok := err.(*apierrors.StatusError); !ok {
				err = apierrors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(localIdentity.Spec.Finalizers) != 0 {
		err = apierrors.NewConflict(auth.Resource("localidentities"), localIdentity.Name, fmt.Errorf("the system is ensuring all content is removed from this policy.  Upon completion, this policy will automatically be purged by the system"))
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a replication controller
type StatusREST struct {
	store *registry.Store
}

// StatusREST implements Patcher
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

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *FinalizeREST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.store, name, options)
}

// Update alters the status finalizers subset of an object.
func (r *FinalizeREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
