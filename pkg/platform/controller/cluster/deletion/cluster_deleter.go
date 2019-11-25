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
	"sync"

	"tkestack.io/tke/pkg/platform/provider"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

// ClusterDeleterInterface to delete a cluster with all resources in it.
type ClusterDeleterInterface interface {
	Delete(clusterName string) error
}

// NewClusterDeleter creates the clusterDeleter object and returns it.
func NewClusterDeleter(clusterClient v1clientset.ClusterInterface,
	platformClient v1clientset.PlatformV1Interface,
	clusterProviders *sync.Map,
	finalizerToken v1.FinalizerName,
	deleteClusterWhenDone bool) ClusterDeleterInterface {
	d := &clusterDeleter{
		clusterClient:         clusterClient,
		platformClient:        platformClient,
		clusterProviders:      clusterProviders,
		deleteClusterWhenDone: deleteClusterWhenDone,
		finalizerToken:        finalizerToken,
	}
	return d
}

var _ ClusterDeleterInterface = &clusterDeleter{}

// clusterDeleter is used to delete all resources in a given cluster.
type clusterDeleter struct {
	// Client to manipulate the cluster.
	clusterClient    v1clientset.ClusterInterface
	platformClient   v1clientset.PlatformV1Interface
	clusterProviders *sync.Map
	// The finalizer token that should be removed from the cluster
	// when all resources in that cluster have been deleted.
	finalizerToken v1.FinalizerName
	// Also delete the cluster when all resources in the cluster have been deleted.
	deleteClusterWhenDone bool
}

// Delete deletes all resources in the given cluster.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   cluster (does nothing if deletion timestamp is missing).
// * Verifies that the cluster is in the "terminating" phase
//   (updates the cluster phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given cluster.
// * Deletes the cluster if deleteClusterWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *clusterDeleter) Delete(clusterName string) error {
	// Multiple controllers may edit a cluster during termination
	// first get the latest state of the cluster before proceeding
	// if the cluster was deleted already, don't do anything
	cluster, err := d.clusterClient.Get(clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if cluster.DeletionTimestamp == nil {
		return nil
	}

	log.Info("Cluster controller - cluster deleter", log.String("clusterName", cluster.Name), log.String("finalizerToken", string(d.finalizerToken)))

	// ensure that the status is up to date on the cluster
	// if we get a not found error, we assume the cluster is truly gone
	cluster, err = d.retryOnConflictError(cluster, d.updateClusterStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the cluster asserts that cluster is no longer deleting..
	if cluster.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the cluster if it is already finalized.
	if d.deleteClusterWhenDone && finalized(cluster) {
		return d.deleteCluster(cluster)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(cluster)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	cluster, err = d.retryOnConflictError(cluster, d.finalizeCluster)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do cluster deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check if we can delete now.
	if d.deleteClusterWhenDone && finalized(cluster) {
		return d.deleteCluster(cluster)
	}
	return nil
}

// Deletes the given cluster.
func (d *clusterDeleter) deleteCluster(cluster *v1.Cluster) error {
	var opts *metav1.DeleteOptions
	uid := cluster.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.clusterClient.Delete(cluster.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateClusterFunc is a function that makes an update to a namespace
type updateClusterFunc func(cluster *v1.Cluster) (*v1.Cluster, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *clusterDeleter) retryOnConflictError(cluster *v1.Cluster, fn updateClusterFunc) (result *v1.Cluster, err error) {
	latestCluster := cluster
	for {
		result, err = fn(latestCluster)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevCluster := latestCluster
		latestCluster, err = d.clusterClient.Get(latestCluster.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevCluster.UID != latestCluster.UID {
			return nil, fmt.Errorf("cluster uid has changed across retries")
		}
	}
}

// updateClusterStatusFunc will verify that the status of the cluster is correct
func (d *clusterDeleter) updateClusterStatusFunc(cluster *v1.Cluster) (*v1.Cluster, error) {
	if cluster.DeletionTimestamp.IsZero() || cluster.Status.Phase == v1.ClusterTerminating {
		return cluster, nil
	}
	newCluster := v1.Cluster{}
	newCluster.ObjectMeta = cluster.ObjectMeta
	newCluster.Status = cluster.Status
	newCluster.Status.Phase = v1.ClusterTerminating
	return d.clusterClient.UpdateStatus(&newCluster)
}

// finalized returns true if the cluster.Spec.Finalizers is an empty list
func finalized(cluster *v1.Cluster) bool {
	return len(cluster.Spec.Finalizers) == 0
}

// finalizeCluster removes the specified finalizerToken and finalizes the cluster
func (d *clusterDeleter) finalizeCluster(cluster *v1.Cluster) (*v1.Cluster, error) {
	clusterFinalize := v1.Cluster{}
	clusterFinalize.ObjectMeta = cluster.ObjectMeta
	clusterFinalize.Spec = cluster.Spec

	finalizerSet := sets.NewString()
	for i := range cluster.Spec.Finalizers {
		if cluster.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(cluster.Spec.Finalizers[i]))
		}
	}
	clusterFinalize.Spec.Finalizers = make([]v1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		clusterFinalize.Spec.Finalizers = append(clusterFinalize.Spec.Finalizers, v1.FinalizerName(value))
	}

	cluster = &v1.Cluster{}
	err := d.platformClient.RESTClient().Put().
		Resource("clusters").
		Name(clusterFinalize.Name).
		SubResource("finalize").
		Body(&clusterFinalize).
		Do().
		Into(cluster)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return cluster, nil
		}
	}
	return cluster, err
}

type deleteResourceFunc func(deleter *clusterDeleter, cluster *v1.Cluster) error

// todo: delete more addons
var deleteResourceFuncs = []deleteResourceFunc{
	deletePersistentEvent,
	deleteHelm,
	deleteClusterProvider,
}

// deleteAllContent will use the client to delete each resource identified in cluster.
func (d *clusterDeleter) deleteAllContent(cluster *v1.Cluster) error {
	log.Debug("Cluster controller - deleteAllContent", log.String("clusterName", cluster.ObjectMeta.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, cluster)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("Cluster controller - deletedAllContent", log.String("clusterName", cluster.ObjectMeta.Name))
	return nil
}

func deletePersistentEvent(deleter *clusterDeleter, cluster *v1.Cluster) error {
	log.Debug("Cluster controller - deletePersistentEvent", log.String("clusterName", cluster.ObjectMeta.Name))

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.ObjectMeta.Name),
	}
	helmList, err := deleter.platformClient.PersistentEvents().List(listOpt)
	if err != nil {
		return err
	}
	if len(helmList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
	for _, pe := range helmList.Items {
		if err := deleter.platformClient.PersistentEvents().Delete(pe.ObjectMeta.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func deleteHelm(deleter *clusterDeleter, cluster *v1.Cluster) error {
	log.Debug("Cluster controller - deleteHelm", log.String("clusterName", cluster.ObjectMeta.Name))

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.ObjectMeta.Name),
	}
	helmList, err := deleter.platformClient.Helms().List(listOpt)
	if err != nil {
		return err
	}
	if len(helmList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := &metav1.DeleteOptions{PropagationPolicy: &background}
	for _, helm := range helmList.Items {
		if err := deleter.platformClient.Helms().Delete(helm.ObjectMeta.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func deleteClusterProvider(deleter *clusterDeleter, cluster *v1.Cluster) error {
	log.Debug("Cluster controller - deleteClusterProvider", log.String("clusterName", cluster.ObjectMeta.Name))

	if cluster.Spec.Type == v1.ClusterImported {
		return nil
	}
	clusterProvider, err := provider.LoadClusterProvider(deleter.clusterProviders, string(cluster.Spec.Type))
	if err != nil {
		panic(err)
	}

	return clusterProvider.OnDelete(*cluster)
}
