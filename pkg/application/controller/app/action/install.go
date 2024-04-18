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

package action

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	"tkestack.io/tke/pkg/application/util"
	chartpath "tkestack.io/tke/pkg/application/util/chartpath/v1"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

// Install installs a chart archive
func Install(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc applicationprovider.UpdateStatusFunc) (*applicationv1.App, error) {
	hooks := getHooks(app)

	err := hooks.PreInstall(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		if updateStatusFunc != nil {
			newStatus := app.Status.DeepCopy()
			var updateStatusErr error
			newStatus.Phase = applicationv1.AppPhaseInstallFailed
			newStatus.Message = "hook pre install app failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
			app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
			if updateStatusErr != nil {
				return app, updateStatusErr
			}
		}
		return nil, err
	}

	destfile, err := Pull(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		newStatus := app.Status.DeepCopy()
		if updateStatusFunc != nil {
			newStatus.Phase = applicationv1.AppPhaseInstallFailed
			newStatus.Message = "fetch chart failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			_, updateStatusErr := updateStatusFunc(ctx, app, &app.Status, newStatus)
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
			if updateStatusErr != nil {
				return nil, updateStatusErr
			}
		}
	}

	client, err := util.NewHelmClientWithProvider(ctx, platformClient, app)
	if err != nil {
		return nil, err
	}
	values, err := helmutil.MergeValues(app.Spec.Values.Values, app.Spec.Values.RawValues, string(app.Spec.Values.RawValuesType))
	if err != nil {
		return nil, err
	}
	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(repo, app.Spec.Chart)
	if err != nil {
		return nil, err
	}
	chartPathBasicOptions.ExistedFile = destfile

	/* provide compatibility with online tke addon apps */
	if app.Annotations != nil && app.Annotations[applicationprovider.AnnotationProviderNameKey] == "managecontrolplane" {
		app.Spec.Chart.InstallPara.Atomic = false
		app.Spec.Chart.InstallPara.Wait = true
		app.Spec.Chart.InstallPara.WaitForJobs = true
		if app.Annotations["ignore-install-wait"] == "true" {
			app.Spec.Chart.InstallPara.Wait = false
			app.Spec.Chart.InstallPara.WaitForJobs = false
		}
		if app.Labels != nil && app.Labels["application.tkestack.io/type"] == "internal-addon" {
			app.Spec.Chart.InstallPara.Wait = false
			app.Spec.Chart.InstallPara.WaitForJobs = false
		}
		if app.Spec.Chart.ChartName == "cranescheduler" {
			app.Spec.Chart.InstallPara.Wait = true
			app.Spec.Chart.InstallPara.WaitForJobs = true
		}
	}
	/* compatibility over, above code need to be deleted atfer the online addon apps are migrated */

	var clientTimeout = defaultTimeout
	if app.Spec.Chart.InstallPara.Timeout > 0 {
		clientTimeout = app.Spec.Chart.InstallPara.Timeout
	}

	_, err = client.Install(ctx, &helmaction.InstallOptions{
		Namespace:        app.Spec.TargetNamespace,
		ReleaseName:      app.Spec.Name,
		DependencyUpdate: true,
		Values:           values,
		Timeout:          clientTimeout,
		ChartPathOptions: chartPathBasicOptions,
		CreateNamespace:  app.Spec.Chart.InstallPara.CreateNamespace,
		Atomic:           app.Spec.Chart.InstallPara.Atomic,
		Wait:             app.Spec.Chart.InstallPara.Wait,
		WaitForJobs:      app.Spec.Chart.InstallPara.WaitForJobs,
	})

	if err != nil {
		if errors.Is(err, errors.New("chart manifest is empty")) {
			log.Errorf(fmt.Sprintf("ERROR: install cluster %s app %s manifest is empty, file %s", app.Spec.TargetCluster, app.Name, destfile))
			metrics.GaugeApplicationManifestFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		} else {
			metrics.GaugeApplicationManifestFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		}
		if updateStatusFunc != nil {
			newStatus := app.Status.DeepCopy()
			var updateStatusErr error
			newStatus.Phase = applicationv1.AppPhaseInstallFailed
			newStatus.Message = "install app failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			if hooks.NeedMetrics(ctx, applicationClient, platformClient, app, repo) {
				metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
			}
			app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
			if updateStatusErr != nil {
				return app, updateStatusErr
			}
		}
		return app, err
	}

	err = hooks.PostInstall(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	// 先走完hook，在更新app状态为succeed
	if err != nil {
		if updateStatusFunc != nil {
			newStatus := app.Status.DeepCopy()
			var updateStatusErr error
			newStatus.Phase = applicationv1.AppPhaseInstallFailed
			newStatus.Message = "hook post install app failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
			app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
			if updateStatusErr != nil {
				return app, updateStatusErr
			}
		}
		return app, err
	}

	if updateStatusFunc != nil {
		newStatus := app.Status.DeepCopy()
		var updateStatusErr error
		newStatus.Phase = applicationv1.AppPhaseSucceeded
		newStatus.Message = ""
		newStatus.Reason = ""
		newStatus.LastTransitionTime = metav1.Now()
		metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
		if updateStatusErr != nil {
			return app, updateStatusErr
		}
	}
	return app, nil
}
