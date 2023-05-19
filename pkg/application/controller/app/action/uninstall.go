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

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/application/util"
)

// Uninstall provides the implementation of 'helm uninstall'.
func Uninstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) (*release.UninstallReleaseResponse, error) {
	hooks := getHooks(app)
	err := hooks.PreUninstall(ctx, applicationClient, platformClient, app, repo)
	if err != nil {
		return nil, err
	}
	client, err := util.NewHelmClientWithProvider(ctx, platformClient, app)
	if err != nil {
		return nil, err
	}
	resp, err := client.Uninstall(&helmaction.UninstallOptions{
		Namespace:   app.Spec.TargetNamespace,
		ReleaseName: app.Spec.Name,
		Timeout:     defaultTimeout,
	})

	if err != nil && !errors.Is(err, driver.ErrReleaseNotFound) {
		return resp, err
	}

	err = hooks.PostUninstall(ctx, applicationClient, platformClient, app, repo)
	return resp, err
}
