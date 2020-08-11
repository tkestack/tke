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
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/controller/app/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/application/util"
	registryutil "tkestack.io/tke/pkg/registry/util"
)

// Pull is the action for pulling a chart.
func Pull(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc updateStatusFunc) (string, error) {
	client, err := util.NewHelmClient(ctx, platformClient, app.Spec.TargetCluster, app.Namespace)
	if err != nil {
		return "", err
	}
	loc := &url.URL{
		Scheme: repo.Scheme,
		Host:   registryutil.BuildTenantRegistryDomain(repo.DomainSuffix, app.Spec.Chart.TenantID),
		Path:   fmt.Sprintf("/chart/%s", app.Spec.Chart.ChartGroupName),
	}
	destfile, err := client.Pull(&helmaction.PullOptions{
		ChartPathOptions: helmaction.ChartPathOptions{
			CaFile:    repo.CaFile,
			Username:  repo.Admin,
			Password:  repo.AdminPassword,
			RepoURL:   loc.String(),
			ChartRepo: app.Spec.Chart.TenantID + "/" + app.Spec.Chart.ChartGroupName,
			Chart:     app.Spec.Chart.ChartName,
			Version:   app.Spec.Chart.ChartVersion,
		},
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
