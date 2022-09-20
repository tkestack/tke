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

// Upgrade upgrade a helm release
func Upgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc applicationprovider.UpdateStatusFunc) (*applicationv1.App, error) {
	hooks := getHooks(app)
	err := hooks.PreUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		return nil, err
	}
	client, err := util.NewHelmClientWithProvider(ctx, platformClient, app)
	if err != nil {
		return nil, err
	}

	destfile, err := Pull(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		newStatus := app.Status.DeepCopy()
		if updateStatusFunc != nil {
			if app.Status.Phase == applicationv1.AppPhaseUpgradFailed {
				log.Error(fmt.Sprintf("upgrade app failed, helm pull err: %s", err.Error()))
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				// delayed retry, queue.AddRateLimited does not meet the demand
				return app, nil
			}
			newStatus.Phase = applicationv1.AppPhaseUpgradFailed
			newStatus.Message = "fetch chart failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			updateStatusFunc(ctx, app, &app.Status, newStatus)
			metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		}
		return nil, err
	}

	newApp, err := applicationClient.Apps(app.Namespace).Get(ctx, app.Name, metav1.GetOptions{})
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
	_, err = client.Upgrade(&helmaction.UpgradeOptions{
		Namespace:        app.Spec.TargetNamespace,
		ReleaseName:      app.Spec.Name,
		DependencyUpdate: true,
		Install:          true,
		Values:           values,
		Timeout:          clientTimeOut,
		ChartPathOptions: chartPathBasicOptions,
	})

	if updateStatusFunc != nil {
		newStatus := newApp.Status.DeepCopy()
		var updateStatusErr error
		if err != nil {
			if app.Status.Phase == applicationv1.AppPhaseUpgradFailed {
				log.Error(fmt.Sprintf("upgrade app failed, helm upgrade err: %s", err.Error()))
				metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				// delayed retry, queue.AddRateLimited does not meet the demand
				return app, nil
			}
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
	if err != nil {
		return nil, err
	}
	err = hooks.PostUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	return newApp, err
}
