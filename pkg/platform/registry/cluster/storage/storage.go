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
	"k8s.io/apiserver/pkg/storage"
	storageerr "k8s.io/apiserver/pkg/storage/errors"
	"k8s.io/apiserver/pkg/util/dryrun"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	clusterstrategy "tkestack.io/tke/pkg/platform/registry/cluster"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/printers"
	printerstorage "tkestack.io/tke/pkg/util/printers/storage"
)

// Storage includes storage for clusters and all sub resources.
type Storage struct {
	Cluster           *REST
	Status            *StatusREST
	Finalize          *FinalizeREST
	Apply             *ApplyREST
	Helm              *HelmREST
	TappController    *TappControllerREST
	CSI               *CSIREST
	PVCR              *PVCRREST
	LogCollector      *LogCollectorREST
	CronHPA           *CronHPAREST
	Addon             *AddonREST
	AddonType         *AddonTypeREST
	LBCFDriver        *LBCFDriverREST
	LBCFLoadBalancer  *LBCFLoadBalancerREST
	LBCFBackendGroup  *LBCFBackendGroupREST
	LBCFBackendRecord *LBCFBackendRecordREST
	Drain             *DrainREST
	Proxy             *ProxyREST
}

// NewStorage returns a Storage object that will work against clusters.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface, host string, privilegedUsername string) *Storage {
	strategy := clusterstrategy.NewStrategy(platformClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &platform.Cluster{} },
		NewListFunc:              func() runtime.Object { return &platform.ClusterList{} },
		DefaultQualifiedResource: platform.Resource("clusters"),
		PredicateFunc:            clusterstrategy.MatchCluster,
		ReturnDeletedObject:      true,

		CreateStrategy: strategy,
		AfterCreate:    strategy.AfterCreate,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
		ExportStrategy: strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(AddHandlers)},
	}
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    clusterstrategy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create cluster etcd rest storage", log.Err(err))
	}

	statusStore := *store
	statusStore.UpdateStrategy = clusterstrategy.NewStatusStrategy(strategy)
	statusStore.ExportStrategy = clusterstrategy.NewStatusStrategy(strategy)

	finalizeStore := *store
	finalizeStore.UpdateStrategy = clusterstrategy.NewFinalizerStrategy(strategy)
	finalizeStore.UpdateStrategy = clusterstrategy.NewFinalizerStrategy(strategy)

	return &Storage{
		Cluster:  &REST{store, privilegedUsername},
		Status:   &StatusREST{&statusStore},
		Finalize: &FinalizeREST{&finalizeStore},
		Apply: &ApplyREST{
			store:          store,
			platformClient: platformClient,
		},
		Helm: &HelmREST{
			store:          store,
			platformClient: platformClient,
		},
		TappController: &TappControllerREST{
			store:          store,
			platformClient: platformClient,
		},
		CSI: &CSIREST{
			store:          store,
			platformClient: platformClient,
		},
		PVCR: &PVCRREST{
			store:          store,
			platformClient: platformClient,
		},
		LogCollector: &LogCollectorREST{
			store:          store,
			platformClient: platformClient,
		},
		CronHPA: &CronHPAREST{
			store:          store,
			platformClient: platformClient,
		},
		Addon: &AddonREST{
			store:          store,
			platformClient: platformClient,
		},
		AddonType: &AddonTypeREST{
			platformClient: platformClient,
			store:          store,
		},
		LBCFDriver: &LBCFDriverREST{
			store:          store,
			platformClient: platformClient,
		},
		LBCFLoadBalancer: &LBCFLoadBalancerREST{
			store:          store,
			platformClient: platformClient,
		},
		LBCFBackendGroup: &LBCFBackendGroupREST{
			store:          store,
			platformClient: platformClient,
		},
		LBCFBackendRecord: &LBCFBackendRecordREST{
			store:          store,
			platformClient: platformClient,
		},
		Drain: &DrainREST{
			store:          store,
			platformClient: platformClient,
		},
		Proxy: &ProxyREST{
			store: store,
			host:  host,
		},
	}
}

// ValidateGetObjectAndTenantID validate name and tenantID, if success return cluster
func ValidateGetObjectAndTenantID(ctx context.Context, store *registry.Store, clusterName string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := store.Get(ctx, clusterName, options)
	if err != nil {
		return nil, err
	}

	cluster := obj.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

// ValidateExportObjectAndTenantID validate name and tenantID, if success return cluster
func ValidateExportObjectAndTenantID(ctx context.Context, store *registry.Store, clusterName string, options metav1.ExportOptions) (runtime.Object, error) {
	obj, err := store.Export(ctx, clusterName, options)
	if err != nil {
		return nil, err
	}

	cluster := obj.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

// REST implements a RESTStorage for clusters against etcd.
type REST struct {
	*registry.Store
	privilegedUsername string
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cls"}
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
		return nil, apierrors.NewMethodNotSupported(platform.Resource("clusters"), "delete collection")
	}
	return r.Store.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, clusterName string, options *metav1.GetOptions) (runtime.Object, error) {
	return ValidateGetObjectAndTenantID(ctx, r.Store, clusterName, options)
}

// Export an object.  Fields that are not user specified are stripped out
// Returns the stripped object.
func (r *REST) Export(ctx context.Context, clusterName string, options metav1.ExportOptions) (runtime.Object, error) {
	return ValidateExportObjectAndTenantID(ctx, r.Store, clusterName, options)
}

// Update finds a resource in the storage and updates it.
func (r *REST) Update(ctx context.Context, clusterName string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	_, err := ValidateGetObjectAndTenantID(ctx, r.Store, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	// log.Info("cluster Update name = " + clusterName + ", validate tenantid ok")

	return r.Store.Update(ctx, clusterName, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}

// Delete enforces life-cycle rules for cluster termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	obj, err := ValidateGetObjectAndTenantID(ctx, r.Store, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	cluster := obj.(*platform.Cluster)
	// Ensure we have a UID precondition
	if options == nil {
		options = metav1.NewDeleteOptions(0)
	}
	if options.Preconditions == nil {
		options.Preconditions = &metav1.Preconditions{}
	}
	if options.Preconditions.UID == nil {
		options.Preconditions.UID = &cluster.UID
	} else if *options.Preconditions.UID != cluster.UID {
		err = apierrors.NewConflict(
			platform.Resource("clusters"),
			name,
			fmt.Errorf("precondition failed: UID in precondition: %v, UID in object meta: %v", *options.Preconditions.UID, cluster.UID),
		)
		return nil, false, err
	}

	// upon first request to delete, we switch the phase to start cluster termination
	if cluster.DeletionTimestamp.IsZero() {
		key, err := r.Store.KeyFunc(ctx, name)
		if err != nil {
			return nil, false, err
		}

		preconditions := storage.Preconditions{UID: options.Preconditions.UID}

		out := r.Store.NewFunc()
		err = r.Store.Storage.GuaranteedUpdate(
			ctx, key, out, false, &preconditions,
			storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
				existingCluster, ok := existing.(*platform.Cluster)
				if !ok {
					// wrong type
					return nil, fmt.Errorf("expected *platformAPI.Cluster, got %v", existing)
				}
				if err := deleteValidation(ctx, existingCluster); err != nil {
					return nil, err
				}
				// Set the deletion timestamp if needed
				if existingCluster.DeletionTimestamp.IsZero() {
					now := metav1.Now()
					existingCluster.DeletionTimestamp = &now
				}
				// Set the cluster phase to terminating, if needed
				if existingCluster.Status.Phase != platform.ClusterTerminating {
					existingCluster.Status.Phase = platform.ClusterTerminating
				}

				// the current finalizers which are on namespace
				currentFinalizers := map[string]bool{}
				for _, f := range existingCluster.Finalizers {
					currentFinalizers[f] = true
				}
				// the finalizers we should ensure on namespace
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
					existingCluster.Finalizers = newFinalizers
				}
				return existingCluster, nil
			}),
			dryrun.IsDryRun(options.DryRun),
		)

		if err != nil {
			err = storageerr.InterpretGetError(err, platform.Resource("clusters"), name)
			err = storageerr.InterpretUpdateError(err, platform.Resource("clusters"), name)
			if _, ok := err.(*apierrors.StatusError); !ok {
				err = apierrors.NewInternalError(err)
			}
			return nil, false, err
		}

		return out, false, nil
	}

	// prior to final deletion, we must ensure that finalizers is empty
	if len(cluster.Spec.Finalizers) != 0 {
		err = apierrors.NewConflict(platform.Resource("clusters"), cluster.Name, fmt.Errorf("the system is ensuring all content is removed from this cluster.  Upon completion, this cluster will automatically be purged by the system"))
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

// FinalizeREST implements the REST endpoint for finalizing a cluster.
type FinalizeREST struct {
	store *registry.Store
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *FinalizeREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the status finalizers subset of an object.
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
