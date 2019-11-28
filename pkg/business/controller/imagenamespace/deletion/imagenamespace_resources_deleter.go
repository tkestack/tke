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
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	businessv1 "tkestack.io/tke/api/business/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/pkg/util/log"
)

// ImageNamespaceResourcesDeleterInterface to delete an imageNamespace with all resources in
// it.
type ImageNamespaceResourcesDeleterInterface interface {
	Delete(projectName, imageNamespaceName string) error
}

// NewImageNamespaceResourcesDeleter to create the imageNamespaceResourcesDeleter that
// implement the ImageNamespaceResourcesDeleterInterface by given businessClient,
// registryClient and configure.
func NewImageNamespaceResourcesDeleter(registryClient registryversionedclient.RegistryV1Interface,
	businessClient v1clientset.BusinessV1Interface, finalizerToken businessv1.FinalizerName,
	deleteImageNamespaceWhenDone bool) ImageNamespaceResourcesDeleterInterface {
	d := &imageNamespaceResourcesDeleter{
		businessClient:               businessClient,
		registryClient:               registryClient,
		finalizerToken:               finalizerToken,
		deleteImageNamespaceWhenDone: deleteImageNamespaceWhenDone,
	}
	return d
}

var _ ImageNamespaceResourcesDeleterInterface = &imageNamespaceResourcesDeleter{}

// imageNamespaceResourcesDeleter is used to delete all resources in a given imageNamespace.
type imageNamespaceResourcesDeleter struct {
	// Client to manipulate the business.
	businessClient v1clientset.BusinessV1Interface
	registryClient registryversionedclient.RegistryV1Interface
	// The finalizer token that should be removed from the imageNamespace
	// when all resources in that imageNamespace have been deleted.
	finalizerToken businessv1.FinalizerName
	// Also delete the imageNamespace when all resources in the imageNamespace have been deleted.
	deleteImageNamespaceWhenDone bool
}

// Delete deletes all resources in the given imageNamespace.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   imageNamespace (does nothing if deletion timestamp is missing).
// * Verifies that the imageNamespace is in the "terminating" phase
//   (updates the imageNamespace phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given imageNamespace.
// * Deletes the imageNamespace if deleteImageNamespaceWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *imageNamespaceResourcesDeleter) Delete(projectName, imageNamespaceName string) error {
	// Multiple controllers may edit an imageNamespace during termination
	// first get the latest state of the imageNamespace before proceeding
	// if the imageNamespace was deleted already, don't do anything
	imageNamespace, err := d.businessClient.ImageNamespaces(projectName).Get(imageNamespaceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if imageNamespace.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("imageNamespace controller - syncImageNamespace - imageNamespace: %s, finalizerToken: %s", imageNamespace.Name, d.finalizerToken)

	// ensure that the status is up to date on the imageNamespace
	// if we get a not found error, we assume the imageNamespace is truly gone
	imageNamespace, err = d.retryOnConflictError(imageNamespace, d.updateImageNamespaceStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the imageNamespace asserts that imageNamespace is no longer deleting.
	if imageNamespace.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the imageNamespace if it is already finalized.
	if d.deleteImageNamespaceWhenDone && finalized(imageNamespace) {
		return d.deleteImageNamespace(imageNamespace)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(imageNamespace)
	if err != nil {
		return err
	}

	// we have removed all content, so mark it finalized by us.
	imageNamespace, err = d.retryOnConflictError(imageNamespace, d.finalizeImageNamespace)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do imageNamespace deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check whether we can delete it now.
	if d.deleteImageNamespaceWhenDone && finalized(imageNamespace) {
		return d.deleteImageNamespace(imageNamespace)
	}
	return nil
}

// Deletes the given imageNamespace.
func (d *imageNamespaceResourcesDeleter) deleteImageNamespace(imageNamespace *businessv1.ImageNamespace) error {
	var opts *metav1.DeleteOptions
	uid := imageNamespace.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.businessClient.ImageNamespaces(imageNamespace.Namespace).Delete(imageNamespace.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateImageNamespaceFunc is a function that makes an update to an imageNamespace
type updateImageNamespaceFunc func(imageNamespace *businessv1.ImageNamespace) (*businessv1.ImageNamespace, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in businessClient code
func (d *imageNamespaceResourcesDeleter) retryOnConflictError(imageNamespace *businessv1.ImageNamespace, fn updateImageNamespaceFunc) (result *businessv1.ImageNamespace, err error) {
	latestImageNamespace := imageNamespace
	for {
		result, err = fn(latestImageNamespace)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevImageNamespace := latestImageNamespace
		latestImageNamespace, err = d.businessClient.ImageNamespaces(latestImageNamespace.Namespace).Get(latestImageNamespace.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevImageNamespace.UID != latestImageNamespace.UID {
			return nil, fmt.Errorf("imageNamespace uid has changed across retries")
		}
	}
}

// updateImageNamespaceStatusFunc will verify that the status of the imageNamespace is correct
func (d *imageNamespaceResourcesDeleter) updateImageNamespaceStatusFunc(imageNamespace *businessv1.ImageNamespace) (*businessv1.ImageNamespace, error) {
	if imageNamespace.DeletionTimestamp.IsZero() || imageNamespace.Status.Phase == businessv1.ImageNamespaceTerminating {
		return imageNamespace, nil
	}
	newImageNamespace := businessv1.ImageNamespace{}
	newImageNamespace.ObjectMeta = imageNamespace.ObjectMeta
	newImageNamespace.Status = imageNamespace.Status
	newImageNamespace.Status.Phase = businessv1.ImageNamespaceTerminating
	return d.businessClient.ImageNamespaces(imageNamespace.Namespace).UpdateStatus(&newImageNamespace)
}

// finalized returns true if the imageNamespace.Spec.Finalizers is an empty list
func finalized(imageNamespace *businessv1.ImageNamespace) bool {
	return len(imageNamespace.Spec.Finalizers) == 0
}

// finalizeImageNamespace removes the specified finalizerToken and finalizes the imageNamespace
func (d *imageNamespaceResourcesDeleter) finalizeImageNamespace(imageNamespace *businessv1.ImageNamespace) (*businessv1.ImageNamespace, error) {
	imageNamespaceFinalize := businessv1.ImageNamespace{}
	imageNamespaceFinalize.ObjectMeta = imageNamespace.ObjectMeta
	imageNamespaceFinalize.Spec = imageNamespace.Spec
	finalizerSet := sets.NewString()
	for i := range imageNamespace.Spec.Finalizers {
		if imageNamespace.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(imageNamespace.Spec.Finalizers[i]))
		}
	}
	imageNamespaceFinalize.Spec.Finalizers = make([]businessv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		imageNamespaceFinalize.Spec.Finalizers = append(imageNamespaceFinalize.Spec.Finalizers, businessv1.FinalizerName(value))
	}

	imageNamespace = &businessv1.ImageNamespace{}
	err := d.businessClient.RESTClient().Put().
		Resource("imagenamespaces").
		Namespace(imageNamespaceFinalize.Namespace).
		Name(imageNamespaceFinalize.Name).
		SubResource("finalize").
		Body(&imageNamespaceFinalize).
		Do().
		Into(imageNamespace)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return imageNamespace, nil
		}
	}
	return imageNamespace, err
}

type deleteResourceFunc func(deleter *imageNamespaceResourcesDeleter, imageNamespace *businessv1.ImageNamespace) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteImageNamespace,
}

// deleteAllContent will use the dynamic businessClient to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *imageNamespaceResourcesDeleter) deleteAllContent(imageNamespace *businessv1.ImageNamespace) error {
	log.Debug("ImageNamespace controller - deleteAllContent", log.String("imageNamespaceName", imageNamespace.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, imageNamespace)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("ImageNamespace controller - deletedAllContent", log.String("imageNamespaceName", imageNamespace.Name))
	return nil
}

func deleteImageNamespace(deleter *imageNamespaceResourcesDeleter, imageNamespace *businessv1.ImageNamespace) error {
	namespaceList, err := deleter.registryClient.Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", imageNamespace.Spec.TenantID, imageNamespace.Name),
	})
	if err != nil {
		return err
	}
	if len(namespaceList.Items) == 0 {
		return nil
	}
	namespaceObject := namespaceList.Items[0]

	err = deleter.registryClient.Namespaces().Delete(namespaceObject.Name, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("deleteImageNamespace(%s), %s", imageNamespace.Name, err)
	}
	return nil
}
