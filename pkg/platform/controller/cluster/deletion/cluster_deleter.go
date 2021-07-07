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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
)

// ClusterDeleterInterface to delete a cluster with all resources in it.
type ClusterDeleterInterface interface {
	Delete(ctx context.Context, clusterName string) error
}

// NewClusterDeleter creates the clusterDeleter object and returns it.
func NewClusterDeleter(clusterClient v1clientset.ClusterInterface,
	platformClient v1clientset.PlatformV1Interface,
	finalizerToken platformv1.FinalizerName,
	deleteClusterWhenDone bool) ClusterDeleterInterface {
	d := &clusterDeleter{
		clusterClient:         clusterClient,
		platformClient:        platformClient,
		deleteClusterWhenDone: deleteClusterWhenDone,
		finalizerToken:        finalizerToken,
	}
	return d
}

var _ ClusterDeleterInterface = &clusterDeleter{}

// clusterDeleter is used to delete all resources in a given cluster.
type clusterDeleter struct {
	// Client to manipulate the cluster.
	clusterClient  v1clientset.ClusterInterface
	platformClient v1clientset.PlatformV1Interface
	// The finalizer token that should be removed from the cluster
	// when all resources in that cluster have been deleted.
	finalizerToken platformv1.FinalizerName
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
func (d *clusterDeleter) Delete(ctx context.Context, clusterName string) error {
	ctx = log.FromContext(ctx).WithName("Delete").WithContext(ctx)

	// Multiple controllers may edit a cluster during termination
	// first get the latest state of the cluster before proceeding
	// if the cluster was deleted already, don't do anything
	cluster, err := d.clusterClient.Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if cluster.DeletionTimestamp == nil {
		return nil
	}

	// ensure that the status is up to date on the cluster
	// if we get a not found error, we assume the cluster is truly gone
	cluster, err = d.retryOnConflictError(ctx, cluster, d.updateClusterStatusFunc)
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
		return d.deleteCluster(ctx, cluster)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(ctx, cluster)
	if err != nil {
		return err
	}

	// we have removed content, so mark it finalized by us
	cluster, err = d.retryOnConflictError(ctx, cluster, d.finalizeCluster)
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
		return d.deleteCluster(ctx, cluster)
	}
	return nil
}

// Deletes the given cluster.
func (d *clusterDeleter) deleteCluster(ctx context.Context, cluster *platformv1.Cluster) error {
	var opts metav1.DeleteOptions
	uid := cluster.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.clusterClient.Delete(ctx, cluster.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateClusterFunc is a function that makes an update to a namespace
type updateClusterFunc func(ctx context.Context, cluster *platformv1.Cluster) (*platformv1.Cluster, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in client code
func (d *clusterDeleter) retryOnConflictError(ctx context.Context, cluster *platformv1.Cluster, fn updateClusterFunc) (result *platformv1.Cluster, err error) {
	latestCluster := cluster
	for {
		result, err = fn(ctx, latestCluster)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevCluster := latestCluster
		latestCluster, err = d.clusterClient.Get(ctx, latestCluster.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevCluster.UID != latestCluster.UID {
			return nil, fmt.Errorf("cluster uid has changed across retries")
		}
	}
}

// updateClusterStatusFunc will verify that the status of the cluster is correct
func (d *clusterDeleter) updateClusterStatusFunc(ctx context.Context, cluster *platformv1.Cluster) (*platformv1.Cluster, error) {
	if cluster.DeletionTimestamp.IsZero() || cluster.Status.Phase == platformv1.ClusterTerminating {
		return cluster, nil
	}
	newCluster := platformv1.Cluster{}
	newCluster.ObjectMeta = cluster.ObjectMeta
	newCluster.Status = cluster.Status
	newCluster.Status.Phase = platformv1.ClusterTerminating
	return d.clusterClient.UpdateStatus(ctx, &newCluster, metav1.UpdateOptions{})
}

// finalized returns true if the cluster.Spec.Finalizers is an empty list
func finalized(cluster *platformv1.Cluster) bool {
	return len(cluster.Spec.Finalizers) == 0
}

// finalizeCluster removes the specified finalizerToken and finalizes the cluster
func (d *clusterDeleter) finalizeCluster(ctx context.Context, cluster *platformv1.Cluster) (*platformv1.Cluster, error) {
	clusterFinalize := platformv1.Cluster{}
	clusterFinalize.ObjectMeta = cluster.ObjectMeta
	clusterFinalize.Spec = cluster.Spec

	finalizerSet := sets.NewString()
	for i := range cluster.Spec.Finalizers {
		if cluster.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(cluster.Spec.Finalizers[i]))
		}
	}
	clusterFinalize.Spec.Finalizers = make([]platformv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		clusterFinalize.Spec.Finalizers = append(clusterFinalize.Spec.Finalizers, platformv1.FinalizerName(value))
	}

	cluster = &platformv1.Cluster{}
	err := d.platformClient.RESTClient().Put().
		Resource("clusters").
		Name(clusterFinalize.Name).
		SubResource("finalize").
		Body(&clusterFinalize).
		Do(ctx).
		Into(cluster)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return cluster, nil
		}
	}
	return cluster, err
}

type deleteResourceFunc func(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error

// todo: delete more addons
var deleteResourceFuncs = []deleteResourceFunc{
	deletePersistentEvent,
	deleteHelm,
	deleteIPAM,
	deleteTappControllers,
	deleteClusterProvider,
	deleteMachine,
}

// deleteAllContent will use the client to delete each resource identified in cluster.
func (d *clusterDeleter) deleteAllContent(ctx context.Context, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteAllContent doing")

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(ctx, d, cluster)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.FromContext(ctx).Info("deleteAllContent done")

	return nil
}

func deletePersistentEvent(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deletePersistentEvent doing")

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.Name),
	}
	helmList, err := deleter.platformClient.PersistentEvents().List(ctx, listOpt)
	if err != nil {
		return err
	}
	if len(helmList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := metav1.DeleteOptions{PropagationPolicy: &background}
	for _, pe := range helmList.Items {
		if err := deleter.platformClient.PersistentEvents().Delete(ctx, pe.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deletePersistentEvent done")

	return nil
}

func deleteHelm(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteHelm doing")

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.Name),
	}
	helmList, err := deleter.platformClient.Helms().List(ctx, listOpt)
	if err != nil {
		return err
	}
	if len(helmList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := metav1.DeleteOptions{PropagationPolicy: &background}
	for _, helm := range helmList.Items {
		if err := deleter.platformClient.Helms().Delete(ctx, helm.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deleteHelm done")

	return nil
}

func deleteIPAM(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteIPAM doing")

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.Name),
	}
	ipamList, err := deleter.platformClient.IPAMs().List(ctx, listOpt)
	if err != nil {
		return err
	}
	if len(ipamList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := metav1.DeleteOptions{PropagationPolicy: &background}
	for _, ipam := range ipamList.Items {
		if err := deleter.platformClient.IPAMs().Delete(ctx, ipam.ObjectMeta.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deleteIPAM done")

	return nil
}

func deleteTappControllers(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteTappControllers doing")

	listOpt := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.clusterName=%s", cluster.Name),
	}
	tappControllerList, err := deleter.platformClient.TappControllers().List(ctx, listOpt)
	if err != nil {
		return err
	}
	if len(tappControllerList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationBackground
	deleteOpt := metav1.DeleteOptions{PropagationPolicy: &background}
	for _, item := range tappControllerList.Items {
		if err := deleter.platformClient.TappControllers().Delete(ctx, item.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deleteTappControllers done")

	return nil
}

func deleteClusterProvider(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteClusterProvider doing")

	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		panic(err)
	}
	clusterWrapper, err := clusterprovider.GetV1Cluster(ctx, deleter.platformClient, cluster)
	if err != nil {
		return err
	}

	err = provider.OnDelete(ctx, clusterWrapper)
	if err != nil {
		return err
	}

	log.FromContext(ctx).Info("deleteClusterProvider done")

	return nil
}

/*
func deleteClusterCredential(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteClusterCredential doing")

	if cluster.Spec.ClusterCredentialRef != nil {
		if err := deleter.platformClient.ClusterCredentials().Delete(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	fieldSelector := fields.OneTermEqualSelector("clusterName", cluster.Name).String()
	clusterCredentialList, err := deleter.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return err
	}
	for _, item := range clusterCredentialList.Items {
		if err := deleter.platformClient.ClusterCredentials().Delete(ctx, item.Name, metav1.DeleteOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deleteClusterCredential done")

	return nil
}
*/

func deleteMachine(ctx context.Context, deleter *clusterDeleter, cluster *platformv1.Cluster) error {
	log.FromContext(ctx).Info("deleteMachine doing")

	fieldSelector := fields.OneTermEqualSelector("spec.clusterName", cluster.Name).String()
	machineList, err := deleter.platformClient.Machines().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return err
	}
	if len(machineList.Items) == 0 {
		return nil
	}
	background := metav1.DeletePropagationForeground
	deleteOpt := metav1.DeleteOptions{PropagationPolicy: &background}
	for _, machine := range machineList.Items {
		if err := deleter.platformClient.Machines().Delete(ctx, machine.Name, deleteOpt); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	log.FromContext(ctx).Info("deleteMachine done")

	return nil
}
