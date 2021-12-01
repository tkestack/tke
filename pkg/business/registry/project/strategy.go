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

package project

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/cmd/tke-business-api/app/options"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/util/validation"
	"tkestack.io/tke/pkg/util/log"

	//businessUtil "tkestack.io/tke/pkg/business/util"
	//platformUtil "tkestack.io/tke/pkg/platform/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	clientrest "k8s.io/client-go/rest"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for project.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	businessClient *businessinternalclient.BusinessClient
	platformClient platformversionedclient.PlatformV1Interface

	features *options.FeatureOptions
}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating project objects.
func NewStrategy(businessClient *businessinternalclient.BusinessClient,
	platformClient platformversionedclient.PlatformV1Interface,
	features *options.FeatureOptions) *Strategy {
	return &Strategy{
		ObjectTyper:    business.Scheme,
		NameGenerator:  namesutil.Generator,
		businessClient: businessClient,
		platformClient: platformClient,
		features:       features,
	}
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldProject := old.(*business.Project)
	project, _ := obj.(*business.Project)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) != 0 {
		if oldProject.Spec.TenantID != tenantID {
			log.Panic("Unauthorized update project information",
				log.String("oldTenantID", oldProject.Spec.TenantID),
				log.String("newTenantID", project.Spec.TenantID),
				log.String("userTenantID", tenantID))
		}
		project.Spec.TenantID = tenantID
	}
	if oldProject.Status.CachedSpecClusters != nil {
		project.Status.CachedSpecClusters = oldProject.Status.CachedSpecClusters
	} else { // For historic data that has no CachedSpecClusters
		project.Status.CachedSpecClusters = oldProject.Spec.Clusters
	}

	if oldProject.Status.CachedParent != nil {
		project.Status.CachedParent = oldProject.Status.CachedParent
	} else { // For historic data that has no CachedParent
		project.Status.CachedParent = &oldProject.Spec.ParentProjectName
	}

	project.Spec.Members = oldProject.Spec.Members
}

// NamespaceScoped is false for projects.
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
	project, _ := obj.(*business.Project)
	if len(tenantID) != 0 {
		project.Spec.TenantID = tenantID
	}

	if project.Name == "" && project.GenerateName == "" {
		project.GenerateName = "prj-"
	}

	project.Spec.Finalizers = []business.FinalizerName{
		business.ProjectFinalize,
	}

	locked := true
	project.Status.Locked = &locked
}

// AfterCreate implements a further operation to run after a resource is
// created and before it is decorated, optional.
func (s *Strategy) AfterCreate(obj runtime.Object, options *metav1.CreateOptions) {
	project, _ := obj.(*business.Project)

	/* for in-cluster mode, create a corresponding namespace */
	if s.features.SyncProjectsWithNamespaces {
		config, err := clientrest.InClusterConfig()
		if err != nil {
			log.Error("after create project failed", log.Any("project", project.Name), log.Err(err))
			return
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("after create project failed", log.Any("project", project.Name), log.Err(err))
			return
		}
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: project.Name,
			},
		}
		_, err = client.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Error("after create project failed", log.Any("project", project.Name), log.Err(err))
		}
	}
}

// AfterDelete implements a further operation to run after a resource
// has been deleted.
func (s *Strategy) AfterDelete(obj runtime.Object, options *metav1.DeleteOptions) {
	project, _ := obj.(*business.Project)

	/* for in-cluster mode, delete the corresponding namespace within the cluster 'global' */
	if s.features.SyncProjectsWithNamespaces {
		config, err := clientrest.InClusterConfig()
		if err != nil {
			log.Error("after delete project failed", log.Any("project", project.Name), log.Err(err))
			return
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("after delete project failed", log.Any("project", project.Name), log.Err(err))
			return
		}
		err = client.CoreV1().Namespaces().Delete(context.Background(), project.Name, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.Error("after delete project failed", log.Any("project", project.Name), log.Err(err))
			return
		}
	}
}

// Validate validates a new project.
func (s *Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return ValidateProject(ctx, obj.(*business.Project), nil,
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}

// AllowCreateOnUpdate is false for projects.
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate returns true if the object can be updated
// unconditionally (irrespective of the latest resource version), when there is
// no resource version specified in the object.
func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end project.
func (s *Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateProjectUpdate(ctx, obj.(*business.Project), old.(*business.Project),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	project, ok := obj.(*business.Project)
	if !ok {
		return nil, nil, fmt.Errorf("not a project")
	}
	return project.ObjectMeta.Labels, ToSelectableFields(project), nil
}

// MatchProject returns a generic matcher for a given label and field selector.
func MatchProject(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
		IndexFields: []string{
			"spec.tenantID",
			"spec.parentProjectName",
			"status.phase",
			"metadata.name",
		},
	}
}

// ToSelectableFields returns a field set that represents the object
func ToSelectableFields(project *business.Project) fields.Set {
	objectMetaFieldsSet := genericregistry.ObjectMetaFieldsSet(&project.ObjectMeta, false)
	specificFieldsSet := fields.Set{
		"spec.tenantID":          project.Spec.TenantID,
		"spec.parentProjectName": project.Spec.ParentProjectName,
		"status.phase":           string(project.Status.Phase),
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
	newProject := obj.(*business.Project)
	oldProject := old.(*business.Project)
	newProject.Spec = oldProject.Spec
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *StatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateProjectUpdate(ctx, obj.(*business.Project), old.(*business.Project),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
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
	newProject := obj.(*business.Project)
	oldProject := old.(*business.Project)

	childProjects := newProject.Status.CalculatedChildProjects
	childNamespaces := newProject.Status.CalculatedNamespaces

	newProject.Status = oldProject.Status

	newProject.Status.CalculatedChildProjects = childProjects
	newProject.Status.CalculatedNamespaces = childNamespaces

	newProject.Spec.Members = oldProject.Spec.Members
}

// ValidateUpdate is invoked after default fields in the object have been
// filled in before the object is persisted.  This method should not mutate
// the object.
func (s *FinalizeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidateProjectUpdate(ctx, obj.(*business.Project), old.(*business.Project),
		validation.NewObjectGetter(s.businessClient), validation.NewClusterGetter(s.platformClient))
}
