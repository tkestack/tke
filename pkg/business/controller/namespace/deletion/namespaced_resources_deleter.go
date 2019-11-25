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
	v1 "tkestack.io/tke/api/business/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/business/util"
	platformutil "tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

// NamespacedResourcesDeleterInterface to delete a namespace with all resources in
// it.
type NamespacedResourcesDeleterInterface interface {
	Delete(projectName string, namespaceName string) error
}

// NewNamespacedResourcesDeleter to create the namespacedResourcesDeleter that
// implement the NamespacedResourcesDeleterInterface by given client and
// configure.
func NewNamespacedResourcesDeleter(platformClient platformversionedclient.PlatformV1Interface, businessClient v1clientset.BusinessV1Interface,
	finalizerToken v1.FinalizerName,
	deleteNamespaceWhenDone bool) NamespacedResourcesDeleterInterface {
	d := &namespacedResourcesDeleter{
		businessClient:          businessClient,
		platformClient:          platformClient,
		finalizerToken:          finalizerToken,
		deleteNamespaceWhenDone: deleteNamespaceWhenDone,
	}
	return d
}

var _ NamespacedResourcesDeleterInterface = &namespacedResourcesDeleter{}

// namespacedResourcesDeleter is used to delete all resources in a given namespace.
type namespacedResourcesDeleter struct {
	// Client to manipulate the business.
	businessClient v1clientset.BusinessV1Interface
	platformClient platformversionedclient.PlatformV1Interface
	// The finalizer token that should be removed from the namespace
	// when all resources in that namespace have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the namespace when all resources in the namespace have been deleted.
	deleteNamespaceWhenDone bool
}

// Delete deletes all resources in the given namespace.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   namespace (does nothing if deletion timestamp is missing).
// * Verifies that the namespace is in the "terminating" phase
//   (updates the namespace phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given namespace.
// * Deletes the namespace if deleteNamespaceWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *namespacedResourcesDeleter) Delete(projectName string, namespaceName string) error {
	// Multiple controllers may edit a namespace during termination
	// first get the latest state of the namespace before proceeding
	// if the namespace was deleted already, don't do anything
	namespace, err := d.businessClient.Namespaces(projectName).Get(namespaceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if namespace.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("namespace controller - syncNamespace - namespace: %s, finalizerToken: %s", namespace.Name, d.finalizerToken)

	// ensure that the status is up to date on the namespace
	// if we get a not found error, we assume the namespace is truly gone
	namespace, err = d.retryOnConflictError(namespace, d.updateNamespaceStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the namespace asserts that namespace is no longer deleting..
	if namespace.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the namespace if it is already finalized.
	if d.deleteNamespaceWhenDone && finalized(namespace) {
		return d.deleteNamespace(namespace)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(namespace)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	namespace, err = d.retryOnConflictError(namespace, d.finalizeNamespace)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do namespace deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteNamespaceWhenDone && finalized(namespace) {
		return d.deleteNamespace(namespace)
	}
	return nil
}

// Deletes the given namespace.
func (d *namespacedResourcesDeleter) deleteNamespace(namespace *v1.Namespace) error {
	var opts *metav1.DeleteOptions
	uid := namespace.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.businessClient.Namespaces(namespace.ObjectMeta.Namespace).Delete(namespace.ObjectMeta.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateNamespaceFunc is a function that makes an update to a namespace
type updateNamespaceFunc func(namespace *v1.Namespace) (*v1.Namespace, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *namespacedResourcesDeleter) retryOnConflictError(namespace *v1.Namespace, fn updateNamespaceFunc) (result *v1.Namespace, err error) {
	latestNamespace := namespace
	for {
		result, err = fn(latestNamespace)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevNamespace := latestNamespace
		latestNamespace, err = d.businessClient.Namespaces(latestNamespace.ObjectMeta.Namespace).Get(latestNamespace.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevNamespace.UID != latestNamespace.UID {
			return nil, fmt.Errorf("namespace uid has changed across retries")
		}
	}
}

// updateNamespaceStatusFunc will verify that the status of the namespace is correct
func (d *namespacedResourcesDeleter) updateNamespaceStatusFunc(namespace *v1.Namespace) (*v1.Namespace, error) {
	if namespace.DeletionTimestamp.IsZero() || namespace.Status.Phase == v1.NamespaceTerminating {
		return namespace, nil
	}
	newNamespace := v1.Namespace{}
	newNamespace.ObjectMeta = namespace.ObjectMeta
	newNamespace.Status = namespace.Status
	newNamespace.Status.Phase = v1.NamespaceTerminating
	return d.businessClient.Namespaces(newNamespace.ObjectMeta.Namespace).UpdateStatus(&newNamespace)
}

// finalized returns true if the namespace.Spec.Finalizers is an empty list
func finalized(namespace *v1.Namespace) bool {
	return len(namespace.Spec.Finalizers) == 0
}

// finalizeNamespace removes the specified finalizerToken and finalizes the namespace
func (d *namespacedResourcesDeleter) finalizeNamespace(namespace *v1.Namespace) (*v1.Namespace, error) {
	namespaceFinalize := v1.Namespace{}
	namespaceFinalize.ObjectMeta = namespace.ObjectMeta
	namespaceFinalize.Spec = namespace.Spec
	finalizerSet := sets.NewString()
	for i := range namespace.Spec.Finalizers {
		if namespace.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(namespace.Spec.Finalizers[i]))
		}
	}
	namespaceFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		namespaceFinalize.Spec.Finalizers = append(namespaceFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	namespace = &v1.Namespace{}
	err := d.businessClient.RESTClient().Put().
		Resource("namespaces").
		Name(namespaceFinalize.Name).
		SubResource("finalize").
		Body(&namespaceFinalize).
		Do().
		Into(namespace)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return namespace, nil
		}
	}
	return namespace, err
}

type deleteResourceFunc func(deleter *namespacedResourcesDeleter, namespace *v1.Namespace) error

var deleteResourceFuncs = []deleteResourceFunc{
	recalculateProjectUsed,
	deleteNamespaceOnCluster,
}

// deleteAllContent will use the dynamic client to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *namespacedResourcesDeleter) deleteAllContent(namespace *v1.Namespace) error {
	log.Debug("Namespace controller - deleteAllContent", log.String("namespaceName", namespace.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, namespace)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Namespace controller - deletedAllContent", log.String("namespaceName", namespace.ObjectMeta.Name))
	return nil
}

func recalculateProjectUsed(deleter *namespacedResourcesDeleter, namespace *v1.Namespace) error {
	log.Debug("Namespace controller - recalculateProjectUsed", log.String("namespaceName", namespace.ObjectMeta.Name))

	project, err := deleter.businessClient.Projects().Get(namespace.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get the project", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("projectName", namespace.ObjectMeta.Namespace), log.Err(err))
		return err
	}
	calculatedNamespaceNames := sets.NewString(project.Status.CalculatedNamespaces...)
	if calculatedNamespaceNames.Has(namespace.ObjectMeta.Name) {
		calculatedNamespaceNames.Delete(namespace.ObjectMeta.Name)
		project.Status.CalculatedNamespaces = calculatedNamespaceNames.List()
		if project.Status.Clusters != nil {
			clusterUsed, clusterUsedExist := project.Status.Clusters[namespace.Spec.ClusterName]
			if clusterUsedExist {
				for k, v := range namespace.Spec.Hard {
					usedValue, ok := clusterUsed.Used[k]
					if ok {
						usedValue.Sub(v)
						clusterUsed.Used[k] = usedValue
					}
				}
				project.Status.Clusters[namespace.Spec.ClusterName] = clusterUsed
			}
		}
		_, err := deleter.businessClient.Projects().Update(project)
		if err != nil {
			log.Error("Failed to update the project status", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("projectName", namespace.ObjectMeta.Namespace), log.Err(err))
			return err
		}
	}

	return nil
}

func deleteNamespaceOnCluster(deleter *namespacedResourcesDeleter, namespace *v1.Namespace) error {
	kubeClient, err := platformutil.BuildExternalClientSetWithName(deleter.platformClient, namespace.Spec.ClusterName)
	if err != nil {
		log.Error("Failed to create the kubernetes client", log.String("namespaceName", namespace.ObjectMeta.Name), log.String("clusterName", namespace.Spec.ClusterName), log.Err(err))
		return err
	}
	ns, err := kubeClient.CoreV1().Namespaces().Get(namespace.Spec.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	projectName, ok := ns.ObjectMeta.Labels[util.LabelProjectName]
	if !ok {
		return fmt.Errorf("no project label were found on the namespace within the business cluster")
	}
	if projectName != namespace.ObjectMeta.Namespace {
		return fmt.Errorf("the namespace in the business cluster currently belongs to another project")
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
	if err := kubeClient.CoreV1().Namespaces().Delete(namespace.Spec.Namespace, deleteOpt); err != nil {
		log.Error("Failed to delete the namespace in the business cluster", log.String("clusterName", namespace.Spec.ClusterName), log.String("namespace", namespace.Spec.Namespace), log.Err(err))
		return err
	}
	return nil
}
