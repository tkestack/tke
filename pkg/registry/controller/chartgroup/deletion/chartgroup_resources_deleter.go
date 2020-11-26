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
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	helm "tkestack.io/tke/pkg/registry/harbor/helmClient"
	"tkestack.io/tke/pkg/util/log"
)

// ChartGroupResourcesDeleterInterface to delete a chartGroup with all resources in
// it.
type ChartGroupResourcesDeleterInterface interface {
	Delete(ctx context.Context, chartGroupName string) error
}

// NewChartGroupResourcesDeleter to create the chartGroupResourcesDeleter that
// implement the ChartGroupResourcesDeleterInterface by given businessClient,
// registryClient and configure.
func NewChartGroupResourcesDeleter(businessClient businessversionedclient.BusinessV1Interface,
	registryClient v1clientset.RegistryV1Interface, finalizerToken registryv1.FinalizerName,
	deleteChartGroupWhenDone bool, helmClient *helm.APIClient) ChartGroupResourcesDeleterInterface {
	d := &chartGroupResourcesDeleter{
		registryClient:           registryClient,
		businessClient:           businessClient,
		finalizerToken:           finalizerToken,
		deleteChartGroupWhenDone: deleteChartGroupWhenDone,
		helmClient:               helmClient,
	}
	return d
}

var _ ChartGroupResourcesDeleterInterface = &chartGroupResourcesDeleter{}

// chartGroupResourcesDeleter is used to delete all resources in a given chartGroup.
type chartGroupResourcesDeleter struct {
	// Client to manipulate the registry.
	registryClient v1clientset.RegistryV1Interface
	businessClient businessversionedclient.BusinessV1Interface
	// The finalizer token that should be removed from the chartGroup
	// when all resources in that chartGroup have been deleted.
	finalizerToken registryv1.FinalizerName
	// Also delete the chartGroup when all resources in the chartGroup have been deleted.
	deleteChartGroupWhenDone bool
	helmClient               *helm.APIClient
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
func (d *chartGroupResourcesDeleter) Delete(ctx context.Context, chartGroupName string) error {
	// Multiple controllers may edit a chartGroup during termination
	// first get the latest state of the chartGroup before proceeding
	// if the chartGroup was deleted already, don't do anything
	chartGroup, err := d.registryClient.ChartGroups().Get(ctx, chartGroupName, metav1.GetOptions{})
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
	chartGroup, err = d.retryOnConflictError(ctx, chartGroup, d.updateChartGroupStatusFunc)
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
		return d.deleteChartGroup(ctx, chartGroup)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(ctx, chartGroup)
	if err != nil {
		return err
	}

	// we have removed all content, so mark it finalized by us.
	chartGroup, err = d.retryOnConflictError(ctx, chartGroup, d.finalizeChartGroup)
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
		return d.deleteChartGroup(ctx, chartGroup)
	}
	return nil
}

// Deletes the given chartGroup.
func (d *chartGroupResourcesDeleter) deleteChartGroup(ctx context.Context, chartGroup *registryv1.ChartGroup) error {
	var opts metav1.DeleteOptions
	uid := chartGroup.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.registryClient.ChartGroups().Delete(ctx, chartGroup.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateChartGroupFunc is a function that makes an update to a chartGroup
type updateChartGroupFunc func(ctx context.Context, chartGroup *registryv1.ChartGroup) (*registryv1.ChartGroup, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in businessClient code
func (d *chartGroupResourcesDeleter) retryOnConflictError(ctx context.Context, chartGroup *registryv1.ChartGroup, fn updateChartGroupFunc) (result *registryv1.ChartGroup, err error) {
	latestChartGroup := chartGroup
	for {
		result, err = fn(ctx, latestChartGroup)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevChartGroup := latestChartGroup
		latestChartGroup, err = d.registryClient.ChartGroups().Get(ctx, latestChartGroup.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevChartGroup.UID != latestChartGroup.UID {
			return nil, fmt.Errorf("chartGroup uid has changed across retries")
		}
	}
}

// updateChartGroupStatusFunc will verify that the status of the chartGroup is correct
func (d *chartGroupResourcesDeleter) updateChartGroupStatusFunc(ctx context.Context, chartGroup *registryv1.ChartGroup) (*registryv1.ChartGroup, error) {
	if chartGroup.DeletionTimestamp.IsZero() || chartGroup.Status.Phase == registryv1.ChartGroupTerminating {
		return chartGroup, nil
	}
	newChartGroup := registryv1.ChartGroup{}
	newChartGroup.ObjectMeta = chartGroup.ObjectMeta
	newChartGroup.Status = chartGroup.Status
	newChartGroup.Status.Phase = registryv1.ChartGroupTerminating
	return d.registryClient.ChartGroups().UpdateStatus(ctx, &newChartGroup, metav1.UpdateOptions{})
}

// finalized returns true if the chartGroup.Spec.Finalizers is an empty list
func finalized(chartGroup *registryv1.ChartGroup) bool {
	return len(chartGroup.Spec.Finalizers) == 0
}

// finalizeChartGroup removes the specified finalizerToken and finalizes the chartGroup
func (d *chartGroupResourcesDeleter) finalizeChartGroup(ctx context.Context, chartGroup *registryv1.ChartGroup) (*registryv1.ChartGroup, error) {
	chartGroupFinalize := registryv1.ChartGroup{}
	chartGroupFinalize.ObjectMeta = chartGroup.ObjectMeta
	chartGroupFinalize.Spec = chartGroup.Spec
	finalizerSet := sets.NewString()
	for i := range chartGroup.Spec.Finalizers {
		if chartGroup.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(chartGroup.Spec.Finalizers[i]))
		}
	}
	chartGroupFinalize.Spec.Finalizers = make([]registryv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		chartGroupFinalize.Spec.Finalizers = append(chartGroupFinalize.Spec.Finalizers, registryv1.FinalizerName(value))
	}
	chartGroup = &registryv1.ChartGroup{}
	var err error
	if d.registryClient.RESTClient() != nil && !reflect.ValueOf(d.registryClient.RESTClient()).IsNil() {
		err = d.registryClient.RESTClient().Put().
			Resource("chartgroups").
			Name(chartGroupFinalize.Name).
			SubResource("finalize").
			Body(&chartGroupFinalize).
			Do(ctx).
			Into(chartGroup)
	} else {
		chartGroup, err = d.registryClient.ChartGroups().Update(ctx, &chartGroupFinalize, metav1.UpdateOptions{})
	}

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return chartGroup, nil
		}
	}
	return chartGroup, err
}

type deleteResourceFunc func(ctx context.Context, deleter *chartGroupResourcesDeleter, chartGroup *registryv1.ChartGroup) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteChartGroup,
	deleteChart,
}

// deleteAllContent will use the dynamic businessClient to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *chartGroupResourcesDeleter) deleteAllContent(ctx context.Context, chartGroup *registryv1.ChartGroup) error {
	log.Info("ChartGroup controller - deleteAllContent", log.String("chartGroupName", chartGroup.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(ctx, d, chartGroup)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Info("ChartGroup controller - deletedAllContent", log.String("chartGroupName", chartGroup.Name))
	return nil
}

func deleteChartGroup(ctx context.Context, deleter *chartGroupResourcesDeleter, chartGroup *registryv1.ChartGroup) error {
	if deleter.businessClient == nil {
		return nil
	}

	var errs []error
	for _, projectID := range chartGroup.Spec.Projects {
		businessChartGroup, err := deleter.businessClient.ChartGroups(projectID).Get(ctx, chartGroup.Spec.Name, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("deleteBusinessChartGroup(%s/%s), %s", projectID, chartGroup.Spec.Name, err))
			}
			continue
		}
		err = deleter.businessClient.ChartGroups(projectID).Delete(ctx, businessChartGroup.Name, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("deleteBusinessChartGroup(%s/%s), %s", projectID, chartGroup.Spec.Name, err))
			}
			continue
		}
		log.Info("ChartGroup controller - deleteBusinessChartGroup",
			log.String("projectID", projectID),
			log.String("chartGroupName", businessChartGroup.Name))
	}
	return utilerrors.NewAggregate(errs)
}

func deleteChart(ctx context.Context, deleter *chartGroupResourcesDeleter, chartGroup *registryv1.ChartGroup) error {
	var errs []error
	charts, err := deleter.registryClient.Charts(chartGroup.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.chartGroupName=%s", chartGroup.Spec.TenantID, chartGroup.Spec.Name),
	})
	if err != nil {
		return fmt.Errorf("deleteChart(%s/*), %s", chartGroup.Name, err)

	}
	for _, c := range charts.Items {
		var opts metav1.DeleteOptions
		uid := c.UID
		if len(uid) > 0 {
			opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
		}
		err := deleter.registryClient.Charts(chartGroup.Name).Delete(ctx, c.Name, opts)
		if err != nil {
			if !errors.IsNotFound(err) {
				errs = append(errs, fmt.Errorf("deleteChart(%s/%s), %s", chartGroup.Name, c.Name, err))
			}
			continue
		}
		log.Info("ChartGroup controller - deleteChart", log.String("chartName", c.Name))
	}
	return utilerrors.NewAggregate(errs)
}
