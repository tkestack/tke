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
	"tkestack.io/tke/pkg/application/util"
)

// Rollback roll back to the previous release
func Rollback(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) (*applicationv1.App, error) {
	hooks := getHooks(app)
	err := hooks.PreRollback(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		return nil, err
	}
	client, err := util.NewHelmClient(ctx, platformClient, app.Spec.TargetCluster, app.Spec.TargetNamespace)
	if err != nil {
		return nil, err
	}

	err = client.Rollback(&helmaction.RollbackOptions{
		Namespace:   app.Spec.TargetNamespace,
		ReleaseName: app.Spec.Name,
		Revision:    app.Status.RollbackRevision,
		Timeout:     clientTimeOut,
	})
	if updateStatusFunc != nil {
		newStatus := app.Status.DeepCopy()
		var updateStatusErr error
		if err != nil {
			newStatus.Phase = applicationv1.AppPhaseRollbackFailed
			newStatus.Message = "rollback app failed"
			newStatus.Reason = err.Error()
			newStatus.LastTransitionTime = metav1.Now()
		} else {
			newStatus.Phase = applicationv1.AppPhaseRolledBack
			newStatus.Message = ""
			newStatus.Reason = ""
			newStatus.LastTransitionTime = metav1.Now()
			newStatus.RollbackRevision = 0 // clean revision
		}
		app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
		if updateStatusErr != nil {
			return app, updateStatusErr
		}
	}
	if err != nil {
		return app, err
	}
	err = hooks.PostRollback(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	return app, err
}
