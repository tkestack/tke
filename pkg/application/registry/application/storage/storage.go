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
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
	applicationapi "tkestack.io/tke/api/application"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	applicationtrategy "tkestack.io/tke/pkg/application/registry/application"
	applicationutil "tkestack.io/tke/pkg/application/util"
	platformfilter "tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for application and all sub resources.
type Storage struct {
	App      *GenericREST
	Status   *StatusREST
	Finalize *FinalizeREST
}

// NewStorage returns a Storage object that will work against application.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter,
	applicationClient *applicationinternalclient.ApplicationClient) *Storage {
	strategy := applicationtrategy.NewStrategy(applicationClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &applicationapi.App{} },
		NewListFunc:              func() runtime.Object { return &applicationapi.AppList{} },
		DefaultQualifiedResource: applicationapi.Resource("apps"),
		PredicateFunc:            applicationtrategy.MatchApplication,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    applicationtrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create application etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = applicationtrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = applicationtrategy.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = applicationtrategy.NewFinalizerStrategy(strategy)
	finalizeStore.ExportStrategy = applicationtrategy.NewFinalizerStrategy(strategy)

	return &Storage{
		App:      &GenericREST{store, applicationClient},
		Status:   &StatusREST{&statusStore},
		Finalize: &FinalizeREST{&finalizeStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return Message
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	repo := obj.(*applicationapi.App)
	if err := applicationutil.FilterApplication(ctx, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return App
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	repo := obj.(*applicationapi.App)
	if err := applicationutil.FilterApplication(ctx, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// GenericREST implements a RESTStorage for application against etcd.
type GenericREST struct {
	*registry.Store
	applicationClient *applicationinternalclient.ApplicationClient
}

var _ rest.ShortNamesProvider = &GenericREST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *GenericREST) ShortNames() []string {
	return []string{"app"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *GenericREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)

	clusterName := platformfilter.ClusterFrom(ctx)
	if clusterName != "" {
		wrappedOptions, _ = apiserverutil.InterceptCustomSelectorFromListOptions(wrappedOptions, "spec.targetCluster", clusterName)
		fieldSelector := fields.OneTermEqualSelector("spec.targetCluster", clusterName)
		wrappedOptions = apiserverutil.FullListOptionsFieldSelector(wrappedOptions, fieldSelector)
	}
	obj, err := r.Store.List(ctx, wrappedOptions)
	if err != nil {
		return nil, err
	}

	fuzzyResourceName := platformfilter.FuzzyResourceFrom(ctx)
	_, fuzzyResourceName = apiserverutil.InterceptFuzzyResourceNameFromListOptions(wrappedOptions, fuzzyResourceName)
	applicationList := obj.(*applicationapi.AppList)
	if fuzzyResourceName != "" {
		var newList []applicationapi.App
		for _, val := range applicationList.Items {
			if strings.Contains(strings.ToLower(val.Spec.Name), strings.ToLower(fuzzyResourceName)) {
				newList = append(newList, val)
			}
		}
		applicationList.Items = newList
	}
	return applicationList, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *GenericREST) Get(ctx context.Context, messageName string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, messageName, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *GenericREST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, name, options)
}

// Update finds a resource in the storage and updates it.
func (r *GenericREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// Delete enforces life-cycle rules for cluster termination
func (r *GenericREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	app := obj.(*applicationapi.App)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &app.UID
	} else if *options.Preconditions.UID != app.UID {
		err = errors.NewConflict(
			applicationapi.Resource("apps"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, app.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != app.ResourceVersion {
		err = errors.NewConflict(
			applicationapi.Resource("apps"),
			name,
			fmt.Errorf("precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v",
				*options.Preconditions.ResourceVersion, app.ResourceVersion),
		)
		return nil, false, err
	}

	if app.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID, ResourceVersion: options.Preconditions.ResourceVersion}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingApplication, ok := existing.(*applicationapi.App)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *application.App, got %v", existing)
				}
				if err := deleteValidation(ctx, existingApplication); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingApplication.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingApplication.DeletionTimestamp = &now
				}
				// Set the App phase to terminating, if needed
				if existingApplication.Status.Phase != applicationapi.AppPhaseTerminating {
					existingApplication.Status.Phase = applicationapi.AppPhaseTerminating
				}

				// the current finalizers which are on App
				currentFinalizers := map[string]bool{}
				for _, f := range existingApplication.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on App
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
					existingApplication.Finalizers = newFinalizers
				}
				return existingApplication, nil
			}),
			dryrun.IsDryRun(options.DryRun),
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, applicationapi.Resource("apps"), name)
			err = storageerr.InterpretUpdateError(err, applicationapi.Resource("apps"), name)
			if _, ok := err.(*errors.StatusError); !ok {
				err = errors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(app.Spec.Finalizers) != 0 {
		err = errors.NewConflict(applicationapi.Resource("apps"), app.Name,
			fmt.Errorf("the system is ensuring all content is removed from App. Upon completion, this App will automatically be purged by the system"))
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the GenericREST endpoint for changing the status of a application request.
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

// FinalizeREST implements the REST endpoint for finalizing a chartgroup.
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
func (r *FinalizeREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
