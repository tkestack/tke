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
	"os"
	"time"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"tkestack.io/tke/pkg/util/file"
	"tkestack.io/tke/pkg/util/log"
)

// InstallOptions is the installation options to a install call.
type InstallOptions struct {
	ChartPathOptions

	DryRun           bool
	DependencyUpdate bool
	Timeout          time.Duration
	Namespace        string
	ReleaseName      string
	Description      string
	// Used by helm template to render charts with .Release.IsUpgrade. Ignored if Dry-Run is false
	IsUpgrade bool

	Values map[string]interface{}
}

// ChartPathOptions captures common options used for controlling chart paths
type ChartPathOptions struct {
	CaFile   string // --ca-file
	CertFile string // --cert-file
	KeyFile  string // --key-file
	Keyring  string // --keyring
	Password string // --password
	RepoURL  string // --repo
	Username string // --username
	Verify   bool   // --verify
	Version  string // --version

	Chart       string
	ChartRepo   string
	ExistedFile string
}

func (cp ChartPathOptions) ApplyTo(opt *action.ChartPathOptions) {
	if opt == nil {
		return
	}
	opt.CaFile = cp.CaFile
	opt.CertFile = cp.CertFile
	opt.KeyFile = cp.KeyFile
	opt.Keyring = cp.Keyring
	opt.Password = cp.Password
	opt.RepoURL = cp.RepoURL
	opt.Username = cp.Username
	opt.Verify = cp.Verify
	opt.Version = cp.Version
}

// Install installs a chart archive
func (c *Client) Install(options *InstallOptions) (*release.Release, error) {
	actionConfig, err := c.buildActionConfig(options.Namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(actionConfig)
	client.DryRun = options.DryRun
	client.DependencyUpdate = options.DependencyUpdate
	client.Timeout = options.Timeout
	client.Namespace = options.Namespace
	client.ReleaseName = options.ReleaseName
	client.Description = options.Description
	client.IsUpgrade = options.IsUpgrade

	options.ChartPathOptions.ApplyTo(&client.ChartPathOptions)

	settings, err := NewSettings(options.ChartRepo)
	if err != nil {
		return nil, err
	}

	// unpack first if need
	root := settings.RepositoryCache
	if options.ExistedFile != "" && file.IsFile(options.ExistedFile) {
		temp, err := ExpandFile(options.ExistedFile, settings.RepositoryCache)
		if err != nil {
			return nil, err
		}
		root = temp
		defer func() {
			os.RemoveAll(temp)
		}()
	}
	chartDir, err := securejoin.SecureJoin(root, options.Chart)
	if err != nil {
		return nil, err
	}

	cp, err := client.ChartPathOptions.LocateChart(chartDir, settings)
	if err != nil {
		return nil, err
	}

	p := getter.All(settings)

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		log.Warnf("This chart %s/%s is deprecated", options.ChartRepo, options.Chart)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				if err := c.DependencyUpdate(cp, p, settings, options.Verify, options.Keyring); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	return client.Run(chartRequested, options.Values)
}

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
