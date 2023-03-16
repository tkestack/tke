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
	newApp, err := applicationClient.Apps(app.Namespace).Get(ctx, app.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	hooks := getHooks(app)

	if newApp.Status.Message != "hook pre upgrade app failed" && newApp.Status.Message != "upgrade app failed" && newApp.Status.Message != "hook post upgrade app failed" {
		newApp.Status.Message = ""
	}
	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre upgrade app failed" {
		err = hooks.PreUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "hook pre upgrade app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				newApp, updateStatusErr = updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
				if updateStatusErr != nil {
					return newApp, updateStatusErr
				}
			}
			return nil, err
		}
	}

	destfile, err := Pull(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		newStatus := newApp.Status.DeepCopy()
		if updateStatusFunc != nil {
			newStatus.Phase = applicationv1.AppPhaseUpgradFailed
			newStatus.Message = "fetch chart failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			_, updateStatusErr := updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(newApp.Spec.TargetCluster, newApp.Name).Set(1)
			if updateStatusErr != nil {
				return newApp, updateStatusErr
			}
		}
		return nil, err
	}

	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre upgrade app failed" || newApp.Status.Message == "upgrade app failed" {
		client, err := util.NewHelmClientWithProvider(ctx, platformClient, app)
		if err != nil {
			return nil, err
		}
		values, err := helmutil.MergeValues(app.Spec.Values.Values, app.Spec.Values.RawValues, string(app.Spec.Values.RawValuesType))
		if err != nil {
			return nil, err
		}
		chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(repo, newApp.Spec.Chart)
		if err != nil {
			return nil, err
		}
		chartPathBasicOptions.ExistedFile = destfile

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
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "upgrade app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				newApp, updateStatusErr = updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
				if updateStatusErr != nil {
					return newApp, updateStatusErr
				}
			}
			return nil, err
		}
	}

	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre upgrade app failed" || newApp.Status.Message == "upgrade app failed" || newApp.Status.Message == "hook post upgrade app failed" {
		err = hooks.PostUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		// 先走完hook，在更新app状态为succeed
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseUpgradFailed
				newStatus.Message = "hook post upgrade app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				newApp, updateStatusErr = updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
				if updateStatusErr != nil {
					return newApp, updateStatusErr
				}
			}
			return newApp, err
		}
	}

	if updateStatusFunc != nil {
		newStatus := newApp.Status.DeepCopy()
		var updateStatusErr error
		if err != nil {
			newStatus.Phase = applicationv1.AppPhaseUpgradFailed
			newStatus.Message = "upgrade app failed"
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
		newApp, updateStatusErr = updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
		if updateStatusErr != nil {
			return newApp, updateStatusErr
		}
	}
	return newApp, nil
}
