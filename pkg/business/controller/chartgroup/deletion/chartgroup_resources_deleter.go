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

// ChartGroupResourcesDeleterInterface to delete a chartGroup with all resources in
// it.
type ChartGroupResourcesDeleterInterface interface {
	Delete(projectName, chartGroupName string) error
}

// NewChartGroupResourcesDeleter to create the chartGroupResourcesDeleter that
// implement the ChartGroupResourcesDeleterInterface by given businessClient,
// registryClient and configure.
func NewChartGroupResourcesDeleter(registryClient registryversionedclient.RegistryV1Interface,
	businessClient v1clientset.BusinessV1Interface, finalizerToken businessv1.FinalizerName,
	deleteChartGroupWhenDone bool) ChartGroupResourcesDeleterInterface {
	d := &chartGroupResourcesDeleter{
		businessClient:           businessClient,
		registryClient:           registryClient,
		finalizerToken:           finalizerToken,
		deleteChartGroupWhenDone: deleteChartGroupWhenDone,
	}
	return d
}

var _ ChartGroupResourcesDeleterInterface = &chartGroupResourcesDeleter{}

// chartGroupResourcesDeleter is used to delete all resources in a given chartGroup.
type chartGroupResourcesDeleter struct {
	// Client to manipulate the business.
	businessClient v1clientset.BusinessV1Interface
	registryClient registryversionedclient.RegistryV1Interface
	// The finalizer token that should be removed from the chartGroup
	// when all resources in that chartGroup have been deleted.
	finalizerToken businessv1.FinalizerName
	// Also delete the chartGroup when all resources in the chartGroup have been deleted.
	deleteChartGroupWhenDone bool
}

// Delete deletes all resources in the given chartGroup.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   chartGroup (does nothing if deletion timestamp is missing).
// * Verifies that the chartGroup is in the "terminating" phase
//   (updates the chartGroup phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given chartGroup.
// * Deletes the chartGroup if deleteChartGroupWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *chartGroupResourcesDeleter) Delete(projectName, chartGroupName string) error {
	// Multiple controllers may edit a chartGroup during termination
	// first get the latest state of the chartGroup before proceeding
	// if the chartGroup was deleted already, don't do anything
	chartGroup, err := d.businessClient.ChartGroups(projectName).Get(chartGroupName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if chartGroup.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("chartGroup controller - syncChartGroup - chartGroup: %s, finalizerToken: %s", chartGroup.Name, d.finalizerToken)

	// ensure that the status is up to date on the chartGroup
	// if we get a not found error, we assume the chartGroup is truly gone
	chartGroup, err = d.retryOnConflictError(chartGroup, d.updateChartGroupStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the chartGroup asserts that chartGroup is no longer deleting.
	if chartGroup.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the chartGroup if it is already finalized.
	if d.deleteChartGroupWhenDone && finalized(chartGroup) {
		return d.deleteChartGroup(chartGroup)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(chartGroup)
	if err != nil {
		return err
	}

	// we have removed all content, so mark it finalized by us.
	chartGroup, err = d.retryOnConflictError(chartGroup, d.finalizeChartGroup)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do chartGroup deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check whether we can delete it now.
	if d.deleteChartGroupWhenDone && finalized(chartGroup) {
		return d.deleteChartGroup(chartGroup)
	}
	return nil
}

// Deletes the given chartGroup.
func (d *chartGroupResourcesDeleter) deleteChartGroup(chartGroup *businessv1.ChartGroup) error {
	var opts *metav1.DeleteOptions
	uid := chartGroup.UID
	if len(uid) > 0 {
		opts = &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.businessClient.ChartGroups(chartGroup.Namespace).Delete(chartGroup.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateChartGroupFunc is a function that makes an update to a chartGroup
type updateChartGroupFunc func(chartGroup *businessv1.ChartGroup) (*businessv1.ChartGroup, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in businessClient code
func (d *chartGroupResourcesDeleter) retryOnConflictError(chartGroup *businessv1.ChartGroup, fn updateChartGroupFunc) (result *businessv1.ChartGroup, err error) {
	latestChartGroup := chartGroup
	for {
		result, err = fn(latestChartGroup)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevChartGroup := latestChartGroup
		latestChartGroup, err = d.businessClient.ChartGroups(latestChartGroup.Namespace).Get(latestChartGroup.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevChartGroup.UID != latestChartGroup.UID {
			return nil, fmt.Errorf("chartGroup uid has changed across retries")
		}
	}
}

// updateChartGroupStatusFunc will verify that the status of the chartGroup is correct
func (d *chartGroupResourcesDeleter) updateChartGroupStatusFunc(chartGroup *businessv1.ChartGroup) (*businessv1.ChartGroup, error) {
	if chartGroup.DeletionTimestamp.IsZero() || chartGroup.Status.Phase == businessv1.ChartGroupTerminating {
		return chartGroup, nil
	}
	newChartGroup := businessv1.ChartGroup{}
	newChartGroup.ObjectMeta = chartGroup.ObjectMeta
	newChartGroup.Status = chartGroup.Status
	newChartGroup.Status.Phase = businessv1.ChartGroupTerminating
	return d.businessClient.ChartGroups(chartGroup.Namespace).UpdateStatus(&newChartGroup)
}

// finalized returns true if the chartGroup.Spec.Finalizers is an empty list
func finalized(chartGroup *businessv1.ChartGroup) bool {
	return len(chartGroup.Spec.Finalizers) == 0
}

// finalizeChartGroup removes the specified finalizerToken and finalizes the chartGroup
func (d *chartGroupResourcesDeleter) finalizeChartGroup(chartGroup *businessv1.ChartGroup) (*businessv1.ChartGroup, error) {
	chartGroupFinalize := businessv1.ChartGroup{}
	chartGroupFinalize.ObjectMeta = chartGroup.ObjectMeta
	chartGroupFinalize.Spec = chartGroup.Spec
	finalizerSet := sets.NewString()
	for i := range chartGroup.Spec.Finalizers {
		if chartGroup.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(chartGroup.Spec.Finalizers[i]))
		}
	}
	chartGroupFinalize.Spec.Finalizers = make([]businessv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		chartGroupFinalize.Spec.Finalizers = append(chartGroupFinalize.Spec.Finalizers, businessv1.FinalizerName(value))
	}

	chartGroup = &businessv1.ChartGroup{}
	err := d.businessClient.RESTClient().Put().
		Resource("chartgroups").
		Namespace(chartGroupFinalize.Namespace).
		Name(chartGroupFinalize.Name).
		SubResource("finalize").
		Body(&chartGroupFinalize).
		Do().
		Into(chartGroup)

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return chartGroup, nil
		}
	}
	return chartGroup, err
}

type deleteResourceFunc func(deleter *chartGroupResourcesDeleter, chartGroup *businessv1.ChartGroup) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteChartGroup,
}

// deleteAllContent will use the dynamic businessClient to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *chartGroupResourcesDeleter) deleteAllContent(chartGroup *businessv1.ChartGroup) error {
	log.Debug("ChartGroup controller - deleteAllContent", log.String("chartGroupName", chartGroup.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(d, chartGroup)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Debug("ChartGroup controller - deletedAllContent", log.String("chartGroupName", chartGroup.Name))
	return nil
}

func deleteChartGroup(deleter *chartGroupResourcesDeleter, chartGroup *businessv1.ChartGroup) error {
	chartGroupList, err := deleter.registryClient.ChartGroups().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", chartGroup.Spec.TenantID, chartGroup.Name),
	})
	if err != nil {
		return err
	}
	if len(chartGroupList.Items) == 0 {
		return nil
	}
	chartGroupObject := chartGroupList.Items[0]

	err = deleter.registryClient.ChartGroups().Delete(chartGroupObject.Name, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("deleteChartGroup(%s), %s", chartGroup.Name, err)
	}
	return nil
}
