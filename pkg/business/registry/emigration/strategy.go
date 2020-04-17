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

package emigration

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
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for emigration.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	businessClient *businessinternalclient.BusinessClient
	platformClient platformversionedclient.PlatformV1Interface
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating emgration objects.
func NewStrategy(businessClient *businessinternalclient.BusinessClient, platformClient platformversionedclient.PlatformV1Interface) *Strategy {
	return &Strategy{business.Scheme, namesutil.Generator, businessClient, platformClient}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldEmigration := old.(*business.NsEmigration)
	newEmigration, _ := obj.(*business.NsEmigration)
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldEmigration.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update emigration information",
				log.String("oldTenantID", oldEmigration.Spec.TenantID),
				log.String("newTenantID", newEmigration.Spec.TenantID),
				log.String("userTenantID", tenantID))
		}
		newEmigration.Spec.TenantID = tenantID
	}
	newEmigration.Spec = oldEmigration.Spec
}

// NamespaceScoped is true for emigrations.
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (s *Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	emigration, _ := obj.(*business.NsEmigration)
	if len(tenantID) != 0 {
		emigration.Spec.TenantID = tenantID
	}

	if emigration.Spec.Namespace != "" {
		emigration.ObjectMeta.GenerateName = ""
		emigration.ObjectMeta.Name = emigration.Spec.Namespace
	} else {
		emigration.ObjectMeta.GenerateName = "nse-"
	}

	emigration.Status = business.NsEmigrationStatus{}
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object) error {
	return nil
}

// Validate validates a new emigration.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateNsEmigrationCreate(obj.(*business.NsEmigration), s.businessClient)
}

// AllowCreateOnUpdate is false for emigrations.
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

// ValidateUpdate is the default update validation for an end emigration.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateNsEmigrationUpdate(obj.(*business.NsEmigration), old.(*business.NsEmigration), s.businessClient)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	emigration, ok := obj.(*business.NsEmigration)
	if !ok {
		return nil, nil, fmt.Errorf("not a emigration")
	}
	return labels.Set(emigration.ObjectMeta.Labels), ToSelectableFields(emigration), nil
}

// MatchNsEmigration returns a generic matcher for a given label and field selector.
func MatchNsEmigration(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.namespace",
			"spec.tenantID",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(emigration *business.NsEmigration) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&emigration.ObjectMeta, true)
	specificFieldsSet := fields.Set{
		"spec.namespace": emigration.Spec.Namespace,
		"spec.tenantID":  emigration.Spec.TenantID,
	}
	return genericregistry.MergeFieldsSets(objectMetaFieldsSet, specificFieldsSet)
}
