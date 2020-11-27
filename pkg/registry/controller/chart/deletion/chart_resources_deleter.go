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

	pathutil "path/filepath"

	"helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	cm_repo "helm.sh/chartmuseum/pkg/repo"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	v1clientset "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	harborHandler "tkestack.io/tke/pkg/registry/harbor/handler"
	helm "tkestack.io/tke/pkg/registry/harbor/helmClient"
	"tkestack.io/tke/pkg/util/log"
)

// ChartResourcesDeleterInterface to delete a chart with all resources in
// it.
type ChartResourcesDeleterInterface interface {
	Delete(ctx context.Context, chartGroupName, chartName string) error
}

// NewChartResourcesDeleter to create the chartResourcesDeleter that
// implement the ChartResourcesDeleterInterface by given authClient,
// registryClient and configure.
func NewChartResourcesDeleter(
	registryClient v1clientset.RegistryV1Interface,
	multiTenantServer *multitenant.MultiTenantServer,
	finalizerToken registryv1.FinalizerName,
	deleteChartWhenDone bool, helmClient *helm.APIClient) ChartResourcesDeleterInterface {
	d := &chartResourcesDeleter{
		registryClient:      registryClient,
		finalizerToken:      finalizerToken,
		multiTenantServer:   multiTenantServer,
		deleteChartWhenDone: deleteChartWhenDone,
		helmClient:          helmClient,
	}
	return d
}

var _ ChartResourcesDeleterInterface = &chartResourcesDeleter{}

// chartResourcesDeleter is used to delete all resources in a given chart.
type chartResourcesDeleter struct {
	// Client to manipulate the registry.
	registryClient v1clientset.RegistryV1Interface
	// The finalizer token that should be removed from the chart
	// when all resources in that chart have been deleted.
	finalizerToken registryv1.FinalizerName
	// the registry server for chartmuseum
	multiTenantServer *multitenant.MultiTenantServer
	// Also delete the chart when all resources in the chart have been deleted.
	deleteChartWhenDone bool
	helmClient          *helm.APIClient
}

// Delete deletes all resources in the given chart.
// Before deleting resources:
// * It ensures that deletion timestamp is set on the
//   chart (does nothing if deletion timestamp is missing).
// * Verifies that the chart is in the "terminating" phase
//   (updates the chart phase if it is not yet marked terminating)
// After deleting the resources:
// * It removes finalizer token from the given chart.
// * Deletes the chart if deleteChartWhenDone is true.
//
// Returns an error if any of those steps fail.
// Returns ResourcesRemainingError if it deleted some resources but needs
// to wait for them to go away.
// Caller is expected to keep calling this until it succeeds.
func (d *chartResourcesDeleter) Delete(ctx context.Context, chartGroupName, chartName string) error {
	// Multiple controllers may edit a chart during termination
	// first get the latest state of the chart before proceeding
	// if the chart was deleted already, don't do anything
	chart, err := d.registryClient.Charts(chartGroupName).Get(ctx, chartName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if chart.DeletionTimestamp == nil {
		return nil
	}

	log.Infof("chart controller - syncChart - chart: %s, finalizerToken: %s", chart.Name, d.finalizerToken)

	// ensure that the status is up to date on the chart
	// if we get a not found error, we assume the chart is truly gone
	chart, err = d.retryOnConflictError(ctx, chart, d.updateChartStatusFunc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// the latest view of the chart asserts that chart is no longer deleting.
	if chart.DeletionTimestamp.IsZero() {
		return nil
	}

	// Delete the chart if it is already finalized.
	if d.deleteChartWhenDone && finalized(chart) {
		return d.deleteChart(ctx, chart)
	}

	// there may still be content for us to remove
	err = d.deleteAllContent(ctx, chart)
	if err != nil {
		return err
	}

	// we have removed all content, so mark it finalized by us.
	chart, err = d.retryOnConflictError(ctx, chart, d.finalizeChart)
	if err != nil {
		// in normal practice, this should not be possible, but if a deployment is running
		// two controllers to do chart deletion that share a common finalizer token it's
		// possible that a not found could occur since the other controller would have finished the delete.
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Check whether we can delete it now.
	if d.deleteChartWhenDone && finalized(chart) {
		return d.deleteChart(ctx, chart)
	}
	return nil
}

// Deletes the given chart.
func (d *chartResourcesDeleter) deleteChart(ctx context.Context, chart *registryv1.Chart) error {
	var opts metav1.DeleteOptions
	uid := chart.UID
	if len(uid) > 0 {
		opts = metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	}
	err := d.registryClient.Charts(chart.Namespace).Delete(ctx, chart.Name, opts)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// updateChartFunc is a function that makes an update to a chart
type updateChartFunc func(ctx context.Context, chart *registryv1.Chart) (*registryv1.Chart, error)

// retryOnConflictError retries the specified fn if there was a conflict error
// it will return an error if the UID for an object changes across retry operations.
// TODO RetryOnConflict should be a generic concept in authClient code
func (d *chartResourcesDeleter) retryOnConflictError(ctx context.Context, chart *registryv1.Chart, fn updateChartFunc) (result *registryv1.Chart, err error) {
	latestChart := chart
	for {
		result, err = fn(ctx, latestChart)
		if err == nil {
			return result, nil
		}
		if !errors.IsConflict(err) {
			return nil, err
		}
		prevChart := latestChart
		latestChart, err = d.registryClient.Charts(latestChart.Namespace).Get(ctx, latestChart.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if prevChart.UID != latestChart.UID {
			return nil, fmt.Errorf("chart uid has changed across retries")
		}
	}
}

// updateChartStatusFunc will verify that the status of the chart is correct
func (d *chartResourcesDeleter) updateChartStatusFunc(ctx context.Context, chart *registryv1.Chart) (*registryv1.Chart, error) {
	if chart.DeletionTimestamp.IsZero() || chart.Status.Phase == registryv1.ChartTerminating {
		return chart, nil
	}
	newChart := registryv1.Chart{}
	newChart.ObjectMeta = chart.ObjectMeta
	newChart.Status = chart.Status
	newChart.Status.Phase = registryv1.ChartTerminating
	return d.registryClient.Charts(newChart.Namespace).UpdateStatus(ctx, &newChart, metav1.UpdateOptions{})
}

// finalized returns true if the chart.Spec.Finalizers is an empty list
func finalized(chart *registryv1.Chart) bool {
	return len(chart.Spec.Finalizers) == 0
}

// finalizeChart removes the specified finalizerToken and finalizes the chart
func (d *chartResourcesDeleter) finalizeChart(ctx context.Context, chart *registryv1.Chart) (*registryv1.Chart, error) {
	chartFinalize := registryv1.Chart{}
	chartFinalize.ObjectMeta = chart.ObjectMeta
	chartFinalize.Spec = chart.Spec
	finalizerSet := sets.NewString()
	for i := range chart.Spec.Finalizers {
		if chart.Spec.Finalizers[i] != d.finalizerToken {
			finalizerSet.Insert(string(chart.Spec.Finalizers[i]))
		}
	}
	chartFinalize.Spec.Finalizers = make([]registryv1.FinalizerName, 0, len(finalizerSet))
	for _, value := range finalizerSet.List() {
		chartFinalize.Spec.Finalizers = append(chartFinalize.Spec.Finalizers, registryv1.FinalizerName(value))
	}
	chart = &registryv1.Chart{}
	var err error
	if d.registryClient.RESTClient() != nil && !reflect.ValueOf(d.registryClient.RESTClient()).IsNil() {
		err = d.registryClient.RESTClient().Put().
			Resource("charts").
			Name(chartFinalize.Name).
			Namespace(chartFinalize.Namespace).
			SubResource("finalize").
			Body(&chartFinalize).
			Do(ctx).
			Into(chart)
	} else {
		chart, err = d.registryClient.Charts(chartFinalize.Namespace).Update(ctx, &chartFinalize, metav1.UpdateOptions{})
	}

	if err != nil {
		// it was removed already, so life is good
		if errors.IsNotFound(err) {
			return chart, nil
		}
	}
	return chart, err
}

type deleteResourceFunc func(ctx context.Context, deleter *chartResourcesDeleter, chart *registryv1.Chart) error

var deleteResourceFuncs = []deleteResourceFunc{
	deleteChart,
}

// deleteAllContent will use the dynamic authClient to delete each resource identified in groupVersionResources.
// It returns an estimate of the time remaining before the remaining resources are deleted.
// If estimate > 0, not all resources are guaranteed to be gone.
func (d *chartResourcesDeleter) deleteAllContent(ctx context.Context, chart *registryv1.Chart) error {
	log.Info("Chart controller - deleteAllContent", log.String("chartGroupName", chart.Namespace), log.String("chartName", chart.Name))

	var errs []error
	for _, deleteFunc := range deleteResourceFuncs {
		err := deleteFunc(ctx, d, chart)
		if err != nil {
			// If there is an error, hold on to it but proceed with all the remaining resource.
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	log.Info("Chart controller - deletedAllContent", log.String("chartGroupName", chart.Namespace), log.String("chartName", chart.Name))
	return nil
}

func deleteChart(ctx context.Context, deleter *chartResourcesDeleter, chart *registryv1.Chart) error {
	var errs []error
	repo := chart.Spec.TenantID + "/" + chart.Spec.ChartGroupName
	name := chart.Spec.Name
	for _, v := range chart.Status.Versions {
		if deleter.helmClient != nil {
			projectName := fmt.Sprintf("%s-chart-%s", chart.Spec.TenantID, chart.Spec.ChartGroupName)
			err := harborHandler.DeleteChart(ctx, deleter.helmClient, projectName, name)
			if err != nil {
				log.Errorf("Deleting harbor chart error: %s", err.Error())
				continue
			}
		} else {
			// refer to: https://github.com/helm/chartmuseum/blob/v0.12.0/pkg/chartmuseum/server/multitenant/api.go#L81
			filename := pathutil.Join(repo, cm_repo.ChartPackageFilenameFromNameVersion(name, v.Version))
			log.Debugf("Deleting package %s from storage", filename)
			deleteObjErr := deleter.multiTenantServer.StorageBackend.DeleteObject(filename)
			if deleteObjErr != nil { //404
				errs = append(errs, fmt.Errorf("Deleting package from storage error: %s", deleteObjErr.Error()))
				log.Errorf("Deleting package from storage error: %s", deleteObjErr.Error())
				continue
			}
			provFilename := pathutil.Join(repo, cm_repo.ProvenanceFilenameFromNameVersion(name, v.Version))
			deleter.multiTenantServer.StorageBackend.DeleteObject(provFilename) // ignore error here, may be no prov file
		}
	}
	return utilerrors.NewAggregate(errs)
}
