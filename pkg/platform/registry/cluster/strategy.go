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

package cluster

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/api/platform/validation"
	"tkestack.io/tke/pkg/apiserver/authentication"
	clusterutil "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for cluster.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	platformClient platforminternalclient.PlatformInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating cluster objects.
func NewStrategy(platformClient platforminternalclient.PlatformInterface) *Strategy {
	return &Strategy{platform.Scheme, namesutil.Generator, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldCluster := old.(*platform.Cluster)
	cluster, _ := obj.(*platform.Cluster)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldCluster.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update cluster information", log.String("oldTenantID", oldCluster.Spec.TenantID), log.String("newTenantID", cluster.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		cluster.Spec.TenantID = tenantID
	}
	if cluster.Spec.Version != oldCluster.Spec.Version && cluster.Spec.Version != cluster.Status.Version {
		cluster.Status.Phase = platform.ClusterUpgrading
	}
	if len(cluster.Spec.Machines) > len(oldCluster.Spec.Machines) {
		cluster.Status.Phase = platform.ClusterUpscaling
		cluster.Spec.ScalingMachines, _ = clusterutil.PrepareClusterScale(cluster, oldCluster)
	}
	if len(cluster.Spec.Machines) < len(oldCluster.Spec.Machines) {
		cluster.Status.Phase = platform.ClusterDownscaling
		cluster.Spec.ScalingMachines, _ = clusterutil.PrepareClusterScale(cluster, oldCluster)
	}
}

// NamespaceScoped is false for clusters
func (Strategy) NamespaceScoped() bool {
	return false
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	cluster, _ := obj.(*platform.Cluster)
	if tenantID != "" {
		cluster.Spec.TenantID = tenantID
	}

	if cluster.Name == "" && cluster.GenerateName == "" {
		cluster.GenerateName = "cls-"
	}

	cluster.Spec.Finalizers = []platform.FinalizerName{
		platform.ClusterFinalize,
	}

	if cluster.Spec.DNSDomain == "" {
		cluster.Spec.DNSDomain = "cluster.local"
	}

	clusterProvider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return // avoid panic validate will be report error
	}
	clusterWrapper, err := types.GetCluster(ctx, s.platformClient, cluster)
	if err != nil {
		panic(err)
	}
	err = clusterProvider.PreCreate(clusterWrapper)
	if err != nil {
		panic(err)
	}
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object) error {
	cluster, _ := obj.(*platform.Cluster)
	clusterProvider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return err
	}
	clusterWrapper, err := types.GetCluster(context.Background(), s.platformClient, cluster)
	if err != nil {
		return err
	}
	err = clusterProvider.AfterCreate(clusterWrapper)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates a new cluster
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	cluster, _ := obj.(*platform.Cluster)
	clusterWrapper, err := types.GetCluster(ctx, s.platformClient, cluster)
	if err != nil {
		return field.ErrorList{field.InternalError(field.NewPath(""), err)}
	}
	return validation.ValidateCluster(clusterWrapper)
}

// AllowCreateOnUpdate is false for clusters
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate returns true if the object can be updated
// unconditionally (irrespective of the latest resource version), when there is
// no resource version specified in the object.
func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end cluster.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	cluster, _ := obj.(*platform.Cluster)
	oldCluster, _ := old.(*platform.Cluster)
	clusterWrapper, err := types.GetCluster(ctx, s.platformClient, cluster)
	if err != nil {
		return field.ErrorList{field.InternalError(field.NewPath(""), err)}
	}
	oldClusterWrapper, err := types.GetCluster(ctx, s.platformClient, oldCluster)
	if err != nil {
		return field.ErrorList{field.InternalError(field.NewPath(""), err)}
	}
	return validation.ValidateClusterUpdate(clusterWrapper, oldClusterWrapper)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	cluster, ok := obj.(*platform.Cluster)
	if !ok {
		return nil, nil, fmt.Errorf("not a cluster")
	}
	return cluster.ObjectMeta.Labels, ToSelectableFields(cluster), nil
}

// MatchCluster returns a generic matcher for a given label and field selector.
func MatchCluster(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID", "spec.type", "spec.version", "status.locked", "status.version", "status.phase"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(cluster *platform.Cluster) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&cluster.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":  cluster.Spec.TenantID,
		"spec.type":      cluster.Spec.Type,
		"spec.version":   cluster.Spec.Version,
		"status.locked":  util.BoolPointerToSelectField(cluster.Status.Locked),
		"status.version": cluster.Status.Version,
		"status.phase":   string(cluster.Status.Phase),
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Cluster.
type StatusStrategy struct {
	*Strategy
}

var _ rest.RESTUpdateStrategy = &StatusStrategy{}

// NewStatusStrategy create the StatusStrategy object by given strategy.
func NewStatusStrategy(strategy *Strategy) *StatusStrategy {
	return &StatusStrategy{strategy}
}

// PrepareForUpdate is invoked on update before validation to normalize
// the object.  For example: remove fields that are not to be persisted,
// sort order-insensitive list fields, etc.  This should not remove fields
// whose presence would be considered a validation error.
func (StatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCluster := obj.(*platform.Cluster)
	oldCluster := old.(*platform.Cluster)
	newCluster.Spec = oldCluster.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// FinalizeStrategy implements finalizer logic for Machine.
type FinalizeStrategy struct {
	*Strategy
}

var _ rest.RESTUpdateStrategy = &FinalizeStrategy{}

// NewFinalizerStrategy create the FinalizeStrategy object by given strategy.
func NewFinalizerStrategy(strategy *Strategy) *FinalizeStrategy {
	return &FinalizeStrategy{strategy}
}

// PrepareForUpdate is invoked on update before validation to normalize
// the object.  For example: remove fields that are not to be persisted,
// sort order-insensitive list fields, etc.  This should not remove fields
// whose presence would be considered a validation error.
func (FinalizeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newCluster := obj.(*platform.Cluster)
	oldCluster := old.(*platform.Cluster)
	newCluster.Status = oldCluster.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
