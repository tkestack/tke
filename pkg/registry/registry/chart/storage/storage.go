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
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	registryapi "tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	platformfilter "tkestack.io/tke/pkg/platform/apiserver/filter"
	chartstrategy "tkestack.io/tke/pkg/registry/registry/chart"
	registryutil "tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for charts and all sub resources.
type Storage struct {
	Chart    *REST
	Status   *StatusREST
	Finalize *FinalizeREST
}

// NewStorage returns a Storage object that will work against charts.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter,
	registryClient *registryinternalclient.RegistryClient,
	authClient authversionedclient.AuthV1Interface,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) *Storage {
	strategy := chartstrategy.NewStrategy(registryClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &registryapi.Chart{} },
		NewListFunc:              func() runtime.Object { return &registryapi.ChartList{} },
		DefaultQualifiedResource: registryapi.Resource("charts"),
		PredicateFunc:            chartstrategy.MatchChart,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    chartstrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create chart etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = chartstrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = chartstrategy.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = chartstrategy.NewFinalizerStrategy(strategy)
	finalizeStore.ExportStrategy = chartstrategy.NewFinalizerStrategy(strategy)

	return &Storage{
		Chart:    &REST{store, registryClient, authClient, businessClient, privilegedUsername},
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

	repo := obj.(*registryapi.Chart)
	if err := registryutil.FilterChart(ctx, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return Chart
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, name string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, name, options)
	if err != nil {
		return nil, err
	}

	repo := obj.(*registryapi.Chart)
	if err := registryutil.FilterChart(ctx, repo); err != nil {
		return nil, err
	}
	return repo, nil
}

// REST implements a RESTStorage for charts against etcd.
type REST struct {
	*registry.Store
	registryClient     *registryinternalclient.RegistryClient
	authClient         authversionedclient.AuthV1Interface
	businessClient     businessversionedclient.BusinessV1Interface
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"chart"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	var obj runtime.Object
	var err error

	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	// repoType is custom label, which not exist in chart
	repoType := ""
	defaultType := "__internal"
	targetProjectID := ""
	wrappedOptions, repoType = apiserverutil.InterceptCustomSelectorFromListOptions(wrappedOptions, "repoType", defaultType)
	wrappedOptions, targetProjectID = apiserverutil.InterceptCustomSelectorFromListOptions(wrappedOptions, "projectID", "")

	switch repoType {
	case registryapi.ScopeTypeUser:
		obj, err = registryutil.ListUserChartsFromStore(ctx, wrappedOptions, r.businessClient, r.registryClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypeProject:
		obj, err = registryutil.ListProjectChartsFromStore(ctx, wrappedOptions, targetProjectID, r.businessClient, r.authClient, r.registryClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypePublic:
		obj, err = registryutil.ListPublicChartsFromStore(ctx, wrappedOptions, r.businessClient, r.registryClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypeAll:
		obj, err = registryutil.ListAllChartsFromStore(ctx, wrappedOptions, targetProjectID, r.businessClient, r.authClient, r.registryClient, r.privilegedUsername, r.Store)
	case defaultType:
		obj, err = r.Store.List(ctx, wrappedOptions)
	default:
		return nil, errors.NewBadRequest(fmt.Sprintf("unsupport spec.repoType: %s", repoType))
	}

	if err != nil {
		return nil, err
	}

	fuzzyResourceName := platformfilter.FuzzyResourceFrom(ctx)
	_, fuzzyResourceName = apiserverutil.InterceptFuzzyResourceNameFromListOptions(wrappedOptions, fuzzyResourceName)
	chartList := obj.(*registryapi.ChartList)
	if fuzzyResourceName != "" {
		var newList []registryapi.Chart
		for _, val := range chartList.Items {
			if strings.Contains(strings.ToLower(val.Spec.Name), strings.ToLower(fuzzyResourceName)) ||
				strings.Contains(strings.ToLower(val.Spec.DisplayName), strings.ToLower(fuzzyResourceName)) {
				newList = append(newList, val)
			}
		}
		chartList.Items = newList
	}
	return chartList, nil
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *REST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, errors.NewMethodNotSupported(registryapi.Resource("charts"), "delete collection")
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
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	chart := obj.(*registryapi.Chart)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &chart.UID
	} else if *options.Preconditions.UID != chart.UID {
		err = errors.NewConflict(
			registryapi.Resource("charts"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, chart.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != chart.ResourceVersion {
		err = errors.NewConflict(
			registryapi.Resource("charts"),
			name,
			fmt.Errorf("precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v",
				*options.Preconditions.ResourceVersion, chart.ResourceVersion),
		)
		return nil, false, err
	}

	if chart.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID, ResourceVersion: options.Preconditions.ResourceVersion}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingChart, ok := existing.(*registryapi.Chart)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *registry.Chart, got %v", existing)
				}
				if err := deleteValidation(ctx, existingChart); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingChart.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingChart.DeletionTimestamp = &now
				}
				// Set the Chart phase to terminating, if needed
				if existingChart.Status.Phase != registryapi.ChartTerminating {
					existingChart.Status.Phase = registryapi.ChartTerminating
				}

				// the current finalizers which are on Chart
				currentFinalizers := map[string]bool{}
				for _, f := range existingChart.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on Chart
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
					existingChart.Finalizers = newFinalizers
				}
				return existingChart, nil
			}),
			dryrun.IsDryRun(options.DryRun),
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, registryapi.Resource("charts"), name)
			err = storageerr.InterpretUpdateError(err, registryapi.Resource("charts"), name)
			if _, ok := err.(*errors.StatusError); !ok {
				err = errors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(chart.Spec.Finalizers) != 0 {
		err = errors.NewConflict(registryapi.Resource("charts"), chart.Name,
			fmt.Errorf("the system is ensuring all content is removed from Chart. Upon completion, this Chart will automatically be purged by the system"))
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the REST endpoint for changing the status of a chart request.
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

// FinalizeREST implements the REST endpoint for finalizing a chart.
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
