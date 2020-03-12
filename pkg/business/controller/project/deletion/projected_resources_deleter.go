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

package deletion

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	businessv1 "tkestack.io/tke/api/business/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	businessUtil "tkestack.io/tke/pkg/business/util"
	"tkestack.io/tke/pkg/util/log"
)

// ProjectedResourcesDeleterInterface to delete a project with all resources in
// it.
type ProjectedResourcesDeleterInterface interface {
	Delete(projectName string) error
}

// NewProjectedResourcesDeleter to create the projectedResourcesDeleter that
// implement the ProjectedResourcesDeleterInterface by given client and
// configure.
func NewProjectedResourcesDeleter(projectClient v1clientset.ProjectInterface,
	businessClient v1clientset.BusinessV1Interface,
	finalizerToken businessv1.FinalizerName,
	deleteProjectWhenDone bool,
	registryEnabled bool) ProjectedResourcesDeleterInterface {
	d := &projectedResourcesDeleter{
		projectClient:         projectClient,
		businessClient:        businessClient,
		finalizerToken:        finalizerToken,
		deleteProjectWhenDone: deleteProjectWhenDone,
		registryEnabled:       registryEnabled,
	}
	return d
}

var _ ProjectedResourcesDeleterInterface = &projectedResourcesDeleter{}

// projectedResourcesDeleter is used to delete all resources in a given project.
type projectedResourcesDeleter struct {
	// Client to manipulate the project.
	projectClient  v1clientset.ProjectInterface
	businessClient v1clientset.BusinessV1Interface
	// The finalizer token that should be removed from the project
	// when all resources in that project have been deleted.
	finalizerToken businessv1.FinalizerName
	// Also delete the project when all resources in the project have been deleted.
	deleteProjectWhenDone bool
	registryEnabled       bool
}

// Delete deletes all resources in the given project.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   project (does nothing if deletion timestamp is missing).
// * Verifies that the project is in the "terminating" phase
//   (updates the project phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given project.
// * Deletes the project if deleteProjectWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *projectedResourcesDeleter) Delete(projectName string) error {
	// Multiple controllers may edit a project during termination
	// first get the latest state of the project before proceeding
	// if the project was deleted already, don't do anything
	project, err := d.projectClient.Get(projectName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if project.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("project controller - syncProject - project: %s, finalizerToken: %s", project.Name, d.finalizerToken)

	// ensure that the status is up to date on the project
	// if we get a not found error, we assume the project is truly gone
	project, err = d.retryOnConflictError(project, d.updateProjectStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the project asserts that project is no longer deleting..
	if project.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the project if it is already finalized.
	if d.deleteProjectWhenDone && finalized(project) {
		return d.deleteProject(project)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(project)
	if err != nil {
		return err
	}

	if !d.hasAllChildrenDeleted(project) {
		return fmt.Errorf("%s, waiting for all children projects or namespaces to be deleted", project.Name)
	}

	// we have removed content, so mark it finalized by us
	project, err = d.retryOnConflictError(project, d.finalizeProject)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do project deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteProjectWhenDone && finalized(project) {
		return d.deleteProject(project)
	}
	return nil
}

// Deletes the given project.
func (d *projectedResourcesDeleter) deleteProject(project *businessv1.Project) error {
	var opts *metav1.DeleteOptions
	uid := project.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.projectClient.Delete(project.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateProjectFunc is a function that makes an update to a project
type updateProjectFunc func(project *businessv1.Project) (*businessv1.Project, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *projectedResourcesDeleter) retryOnConflictError(project *businessv1.Project, fn updateProjectFunc) (result *businessv1.Project, err error) {
	latestProject := project
	for {
		result, err = fn(latestProject)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevProject := latestProject
		latestProject, err = d.projectClient.Get(latestProject.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevProject.UID != latestProject.UID {
			return nil, fmt.Errorf("project uid has changed across retries")
		}
	}
}

// updateProjectStatusFunc will verify that the status of the project is correct
func (d *projectedResourcesDeleter) updateProjectStatusFunc(project *businessv1.Project) (*businessv1.Project, error) {
	if project.DeletionTimestamp.IsZero() || project.Status.Phase == businessv1.ProjectTerminating {
		return project, nil
	}
	newProject := businessv1.Project{}
	newProject.ObjectMeta = project.ObjectMeta
	newProject.Status = project.Status
	newProject.Status.Phase = businessv1.ProjectTerminating
	return d.projectClient.UpdateStatus(&newProject)
}

// finalized returns true if the project.Spec.Finalizers is an empty list
func finalized(project *businessv1.Project) bool {
	return len(project.Spec.Finalizers) == 0
}

// finalizeProject removes the specified finalizerToken and finalizes the project
func (d *projectedResourcesDeleter) finalizeProject(project *businessv1.Project) (*businessv1.Project, error) {
	projectFinalize := businessv1.Project{}
	projectFinalize.ObjectMeta = project.ObjectMeta
	projectFinalize.Spec = project.Spec
	finalizerSet := sets.NewString()
	for i := range project.Spec.Finalizers {
		if project.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(project.Spec.Finalizers[i]))
		}
	}
	projectFinalize.Spec.Finalizers = make([]businessv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		projectFinalize.Spec.Finalizers = append(projectFinalize.Spec.Finalizers, businessv1.FinalizerName(value))
	}

	project = &businessv1.Project{}
	err := d.businessClient.RESTClient().Put().
		Resource("projects").
		Name(projectFinalize.Name).
		SubResource("finalize").
		Body(&projectFinalize).
		Do().
		Into(project)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return project, nil
		}
	}
	return project, err
}

type deleteResourceFunc func(deleter *projectedResourcesDeleter, project *businessv1.Project) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteChartGroups,
	deleteImageNamespaces,
	deleteNamespaces,
	deleteChildProjects,
	recalculateParentProjectUsed,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *projectedResourcesDeleter) deleteAllContent(project *businessv1.Project) error {
	log.Debug("Project controller - deleteAllContent", log.String("projectName", project.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, project)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Project controller - deletedAllContent", log.String("projectName", project.ObjectMeta.Name))
	return nil
}

func recalculateParentProjectUsed(deleter *projectedResourcesDeleter, project *businessv1.Project) error {
	log.Debug("Project controller - recalculateParentProjectUsed", log.String("projectName", project.ObjectMeta.Name))

	if project.Spec.ParentProjectName != "" {
		parentProject, err := deleter.businessClient.Projects().Get(project.Spec.ParentProjectName, metav1.GetOptions{})
		if err != nil {
			log.Error("Failed to get the parent project", log.String("projectName", project.ObjectMeta.Name), log.String("parentProjectName", project.Spec.ParentProjectName), log.Err(err))
			return err
		}
		calculatedChildProjectNames := sets.NewString(parentProject.Status.CalculatedChildProjects...)
		if calculatedChildProjectNames.Has(project.ObjectMeta.Name) {
			calculatedChildProjectNames.Delete(project.ObjectMeta.Name)
			parentProject.Status.CalculatedChildProjects = calculatedChildProjectNames.List()
			if parentProject.Status.Clusters != nil {
				release := project.Spec.Clusters // For historic data that has no CachedSpecClusters
				if project.Status.CachedSpecClusters != nil {
					release = project.Status.CachedSpecClusters
				}
				businessUtil.SubClusterHardFromUsed(&parentProject.Status.Clusters, release)
			}
			_, err := deleter.businessClient.Projects().Update(parentProject)
			if err != nil {
				log.Error("Failed to update the parent project status", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
				return err
			}
		}
	}
	return nil
}

func deleteChildProjects(deleter *projectedResourcesDeleter, project *businessv1.Project) error {
	log.Debug("Project controller - deleteChildProjects", log.String("projectName", project.ObjectMeta.Name))

	childProjectList, err := deleter.businessClient.Projects().List(metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.parentProjectName", project.ObjectMeta.Name).String(),
	})
	if err != nil {
		log.Error("Project controller - failed to list child projects", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return err
	}
	for _, childProject := range childProjectList.Items {
		background := metav1.DeletePropagationBackground
		deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
		if err := deleter.businessClient.Projects().Delete(childProject.ObjectMeta.Name, deleteOpt); err != nil {
			log.Info("Project controller - failed to delete child project", log.String("projectName", project.ObjectMeta.Name), log.String("childProjectName", childProject.ObjectMeta.Name), log.Err(err))
		}
	}

	return nil
}

func deleteNamespaces(deleter *projectedResourcesDeleter, project *businessv1.Project) error {
	log.Debug("Project controller - deleteNamespaces", log.String("projectName", project.ObjectMeta.Name))

	namespaceList, err := deleter.businessClient.Namespaces(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list namespaces", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return err
	}
	for _, namespace := range namespaceList.Items {
		background := metav1.DeletePropagationBackground
		deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
		if err := deleter.businessClient.Namespaces(project.ObjectMeta.Name).Delete(namespace.ObjectMeta.Name, deleteOpt); err != nil {
			log.Info("Project controller - failed to delete namespace", log.String("projectName", project.ObjectMeta.Name), log.String("namespace", namespace.ObjectMeta.Name), log.Err(err))
		}
	}

	return nil
}

func deleteImageNamespaces(deleter *projectedResourcesDeleter, project *businessv1.Project) error {
	if !deleter.registryEnabled {
		return nil
	}
	log.Debug("Project controller - deleteImageNamespaces", log.String("projectName", project.ObjectMeta.Name))

	imageNamespaceList, err := deleter.businessClient.ImageNamespaces(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list imageNamespaces", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return err
	}
	for _, imageNamespace := range imageNamespaceList.Items {
		background := metav1.DeletePropagationBackground
		deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
		if err := deleter.businessClient.ImageNamespaces(project.ObjectMeta.Name).Delete(imageNamespace.ObjectMeta.Name, deleteOpt); err != nil {
			log.Info("Project controller - failed to delete imageNamespace", log.String("projectName", project.ObjectMeta.Name), log.String("imageNamespace", imageNamespace.ObjectMeta.Name), log.Err(err))
		}
	}

	return nil
}

func deleteChartGroups(deleter *projectedResourcesDeleter, project *businessv1.Project) error {
	if !deleter.registryEnabled {
		return nil
	}
	log.Debug("Project controller - deleteChartGroups", log.String("projectName", project.ObjectMeta.Name))

	chartGroupList, err := deleter.businessClient.ChartGroups(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list chartGroups", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return err
	}
	for _, chartGroup := range chartGroupList.Items {
		background := metav1.DeletePropagationBackground
		deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
		if err := deleter.businessClient.ChartGroups(project.ObjectMeta.Name).Delete(chartGroup.ObjectMeta.Name, deleteOpt); err != nil {
			log.Info("Project controller - failed to delete chartGroup", log.String("projectName", project.ObjectMeta.Name), log.String("chartGroup", chartGroup.ObjectMeta.Name), log.Err(err))
		}
	}

	return nil
}

func (d *projectedResourcesDeleter) hasAllChildrenDeleted(project *businessv1.Project) bool {
	return d.hasAllChildProjectsDeleted(project) &&
		d.hasAllNamespacesDeleted(project) &&
		d.hasAllImageNamespacesDeleted(project) &&
		d.hasAllChartGroupsDeleted(project)
}

func (d *projectedResourcesDeleter) hasAllChildProjectsDeleted(project *businessv1.Project) bool {
	childProjectList, err := d.businessClient.Projects().List(metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.parentProjectName", project.ObjectMeta.Name).String(),
	})
	if err != nil {
		log.Error("Project controller - failed to list child projects", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return false
	}
	for _, childProject := range childProjectList.Items {
		_, err = d.businessClient.Projects().Get(childProject.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil || !errors.IsNotFound(err) {
			return false
		}
	}

	return true
}

func (d *projectedResourcesDeleter) hasAllNamespacesDeleted(project *businessv1.Project) bool {
	namespaceList, err := d.businessClient.Namespaces(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list namespaces", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return false
	}
	for _, namespace := range namespaceList.Items {
		_, err = d.businessClient.Namespaces(project.ObjectMeta.Name).Get(namespace.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil || !errors.IsNotFound(err) {
			return false
		}
	}

	return true
}

func (d *projectedResourcesDeleter) hasAllImageNamespacesDeleted(project *businessv1.Project) bool {
	if !d.registryEnabled {
		return true
	}
	imageNamespaceList, err := d.businessClient.ImageNamespaces(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list imageNamespaces", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return false
	}
	for _, imageNamespace := range imageNamespaceList.Items {
		_, err = d.businessClient.ImageNamespaces(project.ObjectMeta.Name).Get(imageNamespace.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil || !errors.IsNotFound(err) {
			return false
		}
	}

	return true
}

func (d *projectedResourcesDeleter) hasAllChartGroupsDeleted(project *businessv1.Project) bool {
	if !d.registryEnabled {
		return true
	}
	chartGroupList, err := d.businessClient.ChartGroups(project.ObjectMeta.Name).List(metav1.ListOptions{})
	if err != nil {
		log.Error("Project controller - failed to list chartGroups", log.String("projectName", project.ObjectMeta.Name), log.Err(err))
		return false
	}
	for _, chartGroup := range chartGroupList.Items {
		_, err = d.businessClient.ChartGroups(project.ObjectMeta.Name).Get(chartGroup.ObjectMeta.Name, metav1.GetOptions{})
		if err == nil || !errors.IsNotFound(err) {
			return false
		}
	}

	return true
}
