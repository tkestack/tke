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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	chartgroupstrategy "tkestack.io/tke/pkg/registry/registry/chartgroup"
	registryutil "tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for chart groups and all sub resources.
type Storage struct {
	ChartGroup *GenericREST
	Status     *StatusREST
	Finalize   *FinalizeREST
}

// NewStorage returns a Storage object that will work against chart groups.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter,
	registryClient *registryinternalclient.RegistryClient,
	authClient authversionedclient.AuthV1Interface,
	businessClient businessversionedclient.BusinessV1Interface,
	privilegedUsername string) *Storage {
	strategy := chartgroupstrategy.NewStrategy(registryClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &registryapi.ChartGroup{} },
		NewListFunc:              func() runtime.Object { return &registryapi.ChartGroupList{} },
		DefaultQualifiedResource: registryapi.Resource("chartgroups"),
		PredicateFunc:            chartgroupstrategy.MatchChartGroup,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    chartgroupstrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create chart group etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = chartgroupstrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = chartgroupstrategy.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = chartgroupstrategy.NewFinalizerStrategy(strategy)
	finalizeStore.ExportStrategy = chartgroupstrategy.NewFinalizerStrategy(strategy)

	return &Storage{
		ChartGroup: &GenericREST{store, authClient, businessClient, privilegedUsername},
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

	o := obj.(*registryapi.ChartGroup)
	if err := registryutil.FilterChartGroup(ctx, o); err != nil {
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

	o := obj.(*registryapi.ChartGroup)
	if err := registryutil.FilterChartGroup(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

// GenericREST implements a RESTStorage for chart groups against etcd.
type GenericREST struct {
	*registry.Store
	authClient         authversionedclient.AuthV1Interface
	businessClient     businessversionedclient.BusinessV1Interface
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &GenericREST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *GenericREST) ShortNames() []string {
	return []string{"rcg"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *GenericREST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	var obj runtime.Object
	var err error

	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	// repoType is custom label, which not exist in chartgroup
	repoType := ""
	defaultType := "__internal"
	targetProjectID := ""
	wrappedOptions, repoType = apiserverutil.InterceptCustomSelectorFromListOptions(wrappedOptions, "repoType", defaultType)
	wrappedOptions, targetProjectID = apiserverutil.InterceptCustomSelectorFromListOptions(wrappedOptions, "projectID", "")

	switch repoType {
	case registryapi.ScopeTypeUser:
		obj, err = registryutil.ListUserChartGroupsFromStore(ctx, wrappedOptions, r.businessClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypeProject:
		obj, err = registryutil.ListProjectChartGroupsFromStore(ctx, wrappedOptions, targetProjectID, r.businessClient, r.authClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypePublic:
		obj, err = registryutil.ListPublicChartGroupsFromStore(ctx, wrappedOptions, r.businessClient, r.privilegedUsername, r.Store)
	case registryapi.ScopeTypeAll:
		obj, err = registryutil.ListAllChartGroupsFromStore(ctx, wrappedOptions, targetProjectID, r.businessClient, r.authClient, r.privilegedUsername, r.Store)
	case defaultType:
		obj, err = r.Store.List(ctx, wrappedOptions)
	default:
		return nil, errors.NewBadRequest(fmt.Sprintf("unsupport repoType: %s", repoType))
	}
	if err != nil {
		return nil, err
	}

	fuzzyResourceName := platformfilter.FuzzyResourceFrom(ctx)
	_, fuzzyResourceName = apiserverutil.InterceptFuzzyResourceNameFromListOptions(wrappedOptions, fuzzyResourceName)
	chartGroupList := obj.(*registryapi.ChartGroupList)
	if fuzzyResourceName != "" {
		var newList []registryapi.ChartGroup
		for _, val := range chartGroupList.Items {
			if strings.Contains(strings.ToLower(val.Spec.Name), strings.ToLower(fuzzyResourceName)) ||
				strings.Contains(strings.ToLower(val.Spec.DisplayName), strings.ToLower(fuzzyResourceName)) {
				newList = append(newList, val)
			}
		}
		chartGroupList.Items = newList
	}
	return chartGroupList, nil
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them.
func (r *GenericREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternal.ListOptions) (runtime.Object, error) {
	if !authentication.IsAdministrator(ctx, r.privilegedUsername) {
		return nil, apierrors.NewMethodNotSupported(registryapi.Resource("chartgroups"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *GenericREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, name, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *GenericREST) Export(ctx context.Context, name string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, name, options)
}

// Update alters the object subset of an object.
func (r *GenericREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	return r.Store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// Delete enforces life-cycle rules for chart group termination
func (r *GenericREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	cg := obj.(*registryapi.ChartGroup)

	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &cg.UID
	} else if *options.Preconditions.UID != cg.UID {
		err = errors.NewConflict(
			registryapi.Resource("chartgroups"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, cg.UID),
		)
		return nil, false, err
	}
	if options.Preconditions.ResourceVersion != nil && *options.Preconditions.ResourceVersion != cg.ResourceVersion {
		err = errors.NewConflict(
			registryapi.Resource("chartgroups"),
			name,
			fmt.Errorf("precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v",
				*options.Preconditions.ResourceVersion, cg.ResourceVersion),
		)
		return nil, false, err
	}

	if cg.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID, ResourceVersion: options.Preconditions.ResourceVersion}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingChartGroup, ok := existing.(*registryapi.ChartGroup)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *registry.ChartGroup, got %v", existing)
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
				if existingChartGroup.Status.Phase != registryapi.ChartGroupTerminating {
					existingChartGroup.Status.Phase = registryapi.ChartGroupTerminating
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
			err = storageerr.InterpretGetError(err, registryapi.Resource("chartgroups"), name)
			err = storageerr.InterpretUpdateError(err, registryapi.Resource("chartgroups"), name)
			if _, ok := err.(*errors.StatusError); !ok {
				err = errors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(cg.Spec.Finalizers) != 0 {
		err = errors.NewConflict(registryapi.Resource("chartgroups"), cg.Name,
			fmt.Errorf("the system is ensuring all content is removed from ChartGroup. Upon completion, this ChartGroup will automatically be purged by the system"))
		return nil, false, err
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}

// StatusREST implements the GenericREST endpoint for changing the status of a chart group.
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

// FinalizeREST implements the GenericREST endpoint for finalizing a chartgroup.
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
