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
	"tkestack.io/tke/pkg/util/metrics"

	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	"tkestack.io/tke/pkg/application/util"
)

// Rollback roll back to the previous release
func Rollback(ctx context.Context,
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

	if newApp.Status.Message != "hook pre rollback app failed" && newApp.Status.Message != "rollback app failed" && newApp.Status.Message != "hook post rollback app failed" {
		newApp.Status.Message = ""
	}

	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre rollback app failed" {
		err = hooks.PreRollback(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseRollbackFailed
				newStatus.Message = "hook pre rollback app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
				if updateStatusErr != nil {
					return app, updateStatusErr
				}
			}
			return app, err
		}
	}

	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre rollback app failed" || newApp.Status.Message == "rollback app failed" {
		client, err := util.NewHelmClientWithProvider(ctx, platformClient, app)
		if err != nil {
			return nil, err
		}
		// rollback成功之后，会把RollbackRevision设置为0
		if app.Status.RollbackRevision != 0 {
			err = client.Rollback(&helmaction.RollbackOptions{
				Namespace:   app.Spec.TargetNamespace,
				ReleaseName: app.Spec.Name,
				Revision:    app.Status.RollbackRevision,
				Timeout:     defaultTimeout,
			})
		}

		if err != nil {
			if updateStatusFunc != nil {
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseRollbackFailed
				newStatus.Message = "rollback app failed"
				newStatus.Reason = err.Error()
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
				if updateStatusErr != nil {
					return app, updateStatusErr
				}
			}
			return app, err
		}
	}

	if newApp.Status.Message == "" || newApp.Status.Message == "hook pre rollback app failed" || newApp.Status.Message == "rollback app failed" || newApp.Status.Message == "hook post rollback app failed" {
		err = hooks.PostRollback(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
		// 先走完hook，在更新app状态为succeed
		if err != nil {
			if updateStatusFunc != nil {
				newStatus := newApp.Status.DeepCopy()
				var updateStatusErr error
				newStatus.Phase = applicationv1.AppPhaseRollbackFailed
				newStatus.Message = "hook post rollback app failed"
				newStatus.Reason = err.Error()
				newStatus.RollbackRevision = 0 // clean revision，next not do rollback again
				newStatus.LastTransitionTime = metav1.Now()
				metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
				app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
				if updateStatusErr != nil {
					return app, updateStatusErr
				}
			}
			return newApp, err
		}
	}

	if updateStatusFunc != nil {
		newStatus := newApp.Status.DeepCopy()
		var updateStatusErr error
		newStatus.Phase = applicationv1.AppPhaseSucceeded
		newStatus.Message = ""
		newStatus.Reason = ""
		newStatus.LastTransitionTime = metav1.Now()
		metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
		newApp, updateStatusErr = updateStatusFunc(ctx, newApp, &newApp.Status, newStatus)
		if updateStatusErr != nil {
			return newApp, updateStatusErr
		}
	}
	return app, nil
}
