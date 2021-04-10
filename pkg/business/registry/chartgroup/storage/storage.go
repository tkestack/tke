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

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/business/registry/chartgroup"
	"tkestack.io/tke/pkg/business/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for ChartGroup and all sub resources.
type Storage struct {
	ChartGroup *REST
	Status     *StatusREST
	Finalize   *FinalizeREST
}

// NewStorage returns a Storage object that will work against ChartGroup sets.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, businessClient *businessinternalclient.BusinessClient,
	registryClient registryversionedclient.RegistryV1Interface, privilegedUsername string) *Storage {
	strategy := chartgroup.NewStrategy(businessClient, registryClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &business.ChartGroup{} },
		NewListFunc:              func() runtime.Object { return &business.ChartGroupList{} },
		DefaultQualifiedResource: business.Resource("chartgroups"),
		PredicateFunc:            chartgroup.MatchChartGroup,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    chartgroup.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create ChartGroup etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = chartgroup.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = chartgroup.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = chartgroup.NewFinalizerStrategy(strategy)
	finalizeStore.ExportStrategy = chartgroup.NewFinalizerStrategy(strategy)

	return &Storage{
		ChartGroup: &REST{store, privilegedUsername},
		Status:     &StatusREST{&statusStore},
		Finalize:   &FinalizeREST{&finalizeStore},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return ChartGroup
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*business.ChartGroup)
	if err := util.FilterChartGroup(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return ChartGroup
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	o := obj.(*business.ChartGroup)
	if err := util.FilterChartGroup(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// REST implements a RESTStorage for ChartGroup sets against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"rcg"}
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, errors.NewMethodNotSupported(business.Resource("chartgroups"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *REST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, name, options)
}

// Update alters the object subset of an object.
func (r *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// Delete enforces life-cycle rules for ChartGroup termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	nsObj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	ns := nsObj.(*business.ChartGroup)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &ns.UID
	} else if *options.Preconditions.UID != ns.UID {
		err = errors.NewConflict(
			business.Resource("chartgroups"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, ns.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != ns.ResourceVersion {
		err = errors.NewConflict(
			business.Resource("chartgroups"),
			name,
			fmt.Errorf("precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v",
				*options.Preconditions.ResourceVersion, ns.ResourceVersion),
		)
		return nil, false, err
	}

	if ns.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID, ResourceVersion: options.Preconditions.ResourceVersion}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingChartGroup, ok := existing.(*business.ChartGroup)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *business.ChartGroup, got %v", existing)
				}
				if err := deleteValidation(ctx, existingChartGroup); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingChartGroup.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingChartGroup.DeletionTimestamp = &now
				}
				// Set the ChartGroup phase to terminating, if needed
				if existingChartGroup.Status.Phase != business.ChartGroupTerminating {
					existingChartGroup.Status.Phase = business.ChartGroupTerminating
				}

				// the current finalizers which are on ChartGroup
				currentFinalizers := map[string]bool{}
				for _, f := range existingChartGroup.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on ChartGroup
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
					existingChartGroup.Finalizers = newFinalizers
				}
				return existingChartGroup, nil
			}),
			dryrun.IsDryRun(options.DryRun),
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, business.Resource("chartgroups"), name)
			err = storageerr.InterpretUpdateError(err, business.Resource("chartgroups"), name)
			if _, ok := err.(*errors.StatusError); !ok {
				err = errors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(ns.Spec.Finalizers) != 0 {
		err = errors.NewConflict(business.Resource("chartgroups"), ns.Name,
			fmt.Errorf("the system is ensuring all content is removed from ChartGroup. Upon completion, this ChartGroup will automatically be purged by the system"))
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
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo,
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

// FinalizeREST implements the REST endpoint for finalizing a project.
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
	_, err := ValidateGetObjectAndTenantID(ctx, r.store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
