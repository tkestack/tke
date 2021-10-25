/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
)

// AnnotationHooksTypeKey specifies the annotation in application object.
const AnnotationHooksTypeKey = "application.tkestack.io/hooks-type"

var hooksMap = map[string]Hooks{}

type Hooks interface {
	PreInstall(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PostInstall(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PreUpgrade(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PostUpgrade(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PreRollback(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PostRollback(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration,
		updateStatusFunc UpdateStatusFunc) error
	PreUninstall(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration) error
	PostUninstall(ctx context.Context,
		applicationClient applicationversionedclient.ApplicationV1Interface,
		platformClient platformversionedclient.PlatformV1Interface,
		app *applicationv1.App,
		repo appconfig.RepoConfiguration) error
}

type EmptyHooks struct{}

func (EmptyHooks) PreInstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PostInstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PreUpgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PostUpgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PreRollback(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PostRollback(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (EmptyHooks) PreUninstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error {
	return nil
}

func (EmptyHooks) PostUninstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error {
	return nil
}

// RegisterHooks will register your hooks with the hooks type, if your hooks work for your application,
// set an annotation with key, application.tkestack.io/hooks-type, and value, the hooksType you registered.
func RegisterHooks(hooksType string, hook Hooks) {
	hooksMap[hooksType] = hook
}

func getHooks(app *applicationv1.App) Hooks {
	if app == nil {
		return EmptyHooks{}
	}
	if hook, ok := hooksMap[app.Annotations[AnnotationHooksTypeKey]]; ok {
		return hook
	}
	return EmptyHooks{}
}
