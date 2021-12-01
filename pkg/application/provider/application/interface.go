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

package application

import (
	"context"

	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
)

type UpdateStatusFunc func(ctx context.Context, app *applicationv1.App, previousStatus, newStatus *applicationv1.AppStatus) (*applicationv1.App, error)

type ControllerProvider interface {
	OnFilter(ctx context.Context, app *applicationv1.App) bool
}

type HooksProvider interface {
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

// Provider defines a set of response interfaces for specific cluster
// types in cluster management.
type Provider interface {
	Name() string

	ControllerProvider
	HooksProvider
}

var _ Provider = &DelegateProvider{}

type DelegateProvider struct {
	ProviderName string
}

func (p *DelegateProvider) Name() string {
	if p.ProviderName == "" {
		return "unknown"
	}
	return p.ProviderName
}

func (p *DelegateProvider) OnFilter(ctx context.Context, app *applicationv1.App) (pass bool) {
	return true
}

func (DelegateProvider) PreInstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PostInstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PreUpgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PostUpgrade(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PreRollback(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PostRollback(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc UpdateStatusFunc) error {
	return nil
}

func (DelegateProvider) PreUninstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error {
	return nil
}

func (DelegateProvider) PostUninstall(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) error {
	return nil
}
