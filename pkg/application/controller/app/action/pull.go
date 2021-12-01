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
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	"tkestack.io/tke/pkg/application/util"
	chartpath "tkestack.io/tke/pkg/application/util/chartpath/v1"
)

// Pull is the action for pulling a chart.
func Pull(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc applicationprovider.UpdateStatusFunc) (string, error) {
	client, err := util.NewHelmClient(ctx, platformClient, app.Spec.TargetCluster, app.Spec.TargetNamespace)
	if err != nil {
		return "", err
	}
	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(repo, app.Spec.Chart)
	if err != nil {
		return "", err
	}

	destfile, err := client.Pull(&helmaction.PullOptions{
		ChartPathOptions: chartPathBasicOptions,
	})
	if updateStatusFunc != nil {
		newStatus := app.Status.DeepCopy()
		if err != nil {
			newStatus.Phase = applicationv1.AppPhaseChartFetchFailed
			newStatus.Message = "fetch chart failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
			updateStatusFunc(ctx, app, &app.Status, newStatus)
			return destfile, err
		}
		newStatus.Phase = applicationv1.AppPhaseChartFetched
		newStatus.Message = ""
		newStatus.Reason = ""
		newStatus.LastTransitionTime = metav1.Now()
		_, err := updateStatusFunc(ctx, app, &app.Status, newStatus)
		if err != nil {
			return destfile, err
		}
	}
	return destfile, err
}
