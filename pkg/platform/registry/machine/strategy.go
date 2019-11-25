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

package machine

import (
	"context"
	"fmt"
	"sync"

	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"tkestack.io/tke/pkg/util/log"

	"tkestack.io/tke/api/platform"

	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for machine.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	machineProviders *sync.Map
	platformClient   platforminternalclient.PlatformInterface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating machine objects.
func NewStrategy(machineProviders *sync.Map, platformClient platforminternalclient.PlatformInterface) *Strategy {
	return &Strategy{platform.Scheme, namesutil.Generator, machineProviders, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldMachine := old.(*platform.Machine)
	machine, _ := obj.(*platform.Machine)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldMachine.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update machine information", log.String("oldTenantID", oldMachine.Spec.TenantID), log.String("newTenantID", machine.Spec.TenantID), log.String("userTenantID", tenantID))
		}
		machine.Spec.TenantID = tenantID
	}
}

// NamespaceScoped is false for machines
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
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	machine, _ := obj.(*platform.Machine)
	if len(tenantID) != 0 {
		machine.Spec.TenantID = tenantID
	}

	if machine.Name == "" && machine.GenerateName == "" {
		machine.GenerateName = "mc-"
	}

	machine.Spec.Finalizers = []platform.FinalizerName{
		platform.MachineFinalize,
	}
}

// Validate validates a new machine
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return Validate(s.machineProviders, obj.(*platform.Machine), s.platformClient)
}

// AllowCreateOnUpdate is false for machines
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
	return ValidateUpdate(s.machineProviders, obj.(*platform.Machine), old.(*platform.Machine))
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	machine, ok := obj.(*platform.Machine)
	if !ok {
		return nil, nil, fmt.Errorf("not a machine")
	}
	return labels.Set(machine.ObjectMeta.Labels), ToSelectableFields(machine), nil
}

// SelectionPredicate returns a generic matcher for a given label and field selector.
func SelectionPredicate(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID", "spec.type", "spec.version", "status.locked", "status.version", "status.phase"},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(machine *platform.Machine) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&machine.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":    machine.Spec.TenantID,
		"spec.type":        string(machine.Spec.Type),
		"spec.clusterName": machine.Spec.ClusterName,
		"status.locked":    util.BoolPointerToSelectField(machine.Status.Locked),
		"status.phase":     string(machine.Status.Phase),
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}

// StatusStrategy implements verification logic for status of Machine.
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
	newMachine := obj.(*platform.Machine)
	oldMachine := old.(*platform.Machine)
	newMachine.Spec = oldMachine.Spec
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
	newMachine := obj.(*platform.Machine)
	oldMachine := old.(*platform.Machine)
	newMachine.Status = oldMachine.Status
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return nil
}
