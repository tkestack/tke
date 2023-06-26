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
	"tkestack.io/tke/pkg/util/metrics"
)

// Upgrade upgrade a helm release
func Upgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc applicationprovider.UpdateStatusFunc) (*applicationv1.App, error) {
	hooks := getHooks(app)

	if app.Status.Message != "hook pre upgrade app failed" && app.Status.Message != "upgrade app failed" && app.Status.Message != "hook post upgrade app failed" {
		app.Status.Message = ""
	}
	if app.Status.Message == "" || app.Status.Message == "hook pre upgrade app failed" {
		err := hooks.PreUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := app.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "hook pre upgrade app failed"
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
	}

	destfile, err := Pull(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		newStatus := app.Status.DeepCopy()
		if updateStatusFunc != nil {
			newStatus.Phase = applicationv1.AppPhaseUpgradFailed
			newStatus.Message = "fetch chart failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			_, updateStatusErr := updateStatusFunc(ctx, app, &app.Status, newStatus)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
			if updateStatusErr != nil {
				return app, updateStatusErr
			}
		}
		return nil, err
	}

	if app.Status.Message == "" || app.Status.Message == "hook pre upgrade app failed" || app.Status.Message == "upgrade app failed" {
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
		// * provide compatibility with online tke addon apps */
		if app.Annotations != nil && app.Annotations[applicationprovider.AnnotationProviderNameKey] == "managecontrolplane" {
			app.Spec.Chart.UpgradePara.Atomic = false
			app.Spec.Chart.UpgradePara.Wait = true
			app.Spec.Chart.UpgradePara.WaitForJobs = true
			if app.Annotations["ignore-upgrade-wait"] == "true" {
				app.Spec.Chart.UpgradePara.Wait = false
				app.Spec.Chart.UpgradePara.WaitForJobs = false
			}
		}
		// * compatibility over, above code need to be deleted atfer the online addon apps are migrated */

		var clientTimeout = defaultTimeout
		if app.Spec.Chart.UpgradePara.Timeout > 0 {
			clientTimeout = app.Spec.Chart.UpgradePara.Timeout
		}

		_, err = client.Upgrade(ctx, &helmaction.UpgradeOptions{
			Namespace:        app.Spec.TargetNamespace,
			ReleaseName:      app.Spec.Name,
			DependencyUpdate: true,
			Install:          true,
			Values:           values,
			Timeout:          clientTimeout,
			ChartPathOptions: chartPathBasicOptions,
			MaxHistory:       clientMaxHistory,
			Atomic:           app.Spec.Chart.UpgradePara.Atomic,
			Wait:             app.Spec.Chart.UpgradePara.Wait,
			WaitForJobs:      app.Spec.Chart.UpgradePara.WaitForJobs,
		})
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := app.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "upgrade app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
				if updateStatusErr != nil {
					return app, updateStatusErr
				}
			}
			return nil, err
		}
	}

	if app.Status.Message == "" || app.Status.Message == "hook pre upgrade app failed" || app.Status.Message == "upgrade app failed" || app.Status.Message == "hook post upgrade app failed" {
		err = hooks.PostUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		// 先走完hook，在更新app状态为succeed
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := app.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "hook post upgrade app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
				if updateStatusErr != nil {
					return app, updateStatusErr
				}
			}
			return app, err
		}
	}

	if updateStatusFunc != nil {
		newStatus := app.Status.DeepCopy()
		var updateStatusErr error
		if err != nil {
			newStatus.Phase = applicationv1.AppPhaseUpgradFailed
			newStatus.Message = ""
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		} else {
			newStatus.Phase = applicationv1.AppPhaseSucceeded
			newStatus.Message = ""
			newStatus.Reason = ""
			newStatus.LastTransitionTime = metav1.Now()
			metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
			metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		}
		app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
		if updateStatusErr != nil {
			return app, updateStatusErr
		}
	}
	return app, nil
}
