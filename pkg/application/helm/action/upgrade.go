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
	"os"
	"time"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"tkestack.io/tke/pkg/util/file"
	"tkestack.io/tke/pkg/util/log"
)

// UpgradeOptions is the installation options to a install call.
type UpgradeOptions struct {
	ChartPathOptions

	// Install is a purely informative flag that indicates whether this upgrade was done in "install" mode.
	//
	// Applications may use this to determine whether this Upgrade operation was done as part of a
	// pure upgrade (Upgrade.Install == false) or as part of an install-or-upgrade operation
	// (Upgrade.Install == true).
	//
	// Setting this to `true` will NOT cause `Upgrade` to perform an install if the release does not exist.
	// That process must be handled by creating an Install action directly. See cmd/upgrade.go for an
	// example of how this flag is used.
	Install     bool
	DryRun      bool
	Timeout     time.Duration
	Namespace   string
	Description string
	// ResetValues will reset the values to the chart's built-ins rather than merging with existing.
	ResetValues bool
	// ReuseValues will re-use the user's last supplied values.
	ReuseValues bool
	// MaxHistory limits the maximum number of revisions saved per release
	MaxHistory int

	DependencyUpdate bool
	ReleaseName      string
	Values           map[string]interface{}
	Atomic           bool
	Wait             bool
	WaitForJobs      bool
}

// Upgrade upgrade a helm release
func (c *Client) Upgrade(ctx context.Context, options *UpgradeOptions) (*release.Release, error) {
	actionConfig, err := c.buildActionConfig(options.Namespace)
	if err != nil {
		return nil, err
	}
	if options.Install {
		// If a release does not exist, install it.
		histClient := action.NewHistory(actionConfig)
		histClient.Max = 1
		rels, err := histClient.Run(options.ReleaseName)
		if errors.Is(err, driver.ErrReleaseNotFound) {
			log.Infof("Release %d does not exist. Installing it now.", options.ReleaseName)
			return c.Install(ctx, &InstallOptions{
				DryRun:           options.DryRun,
				DependencyUpdate: options.DependencyUpdate,
				Timeout:          options.Timeout,
				Namespace:        options.Namespace,
				ReleaseName:      options.ReleaseName,
				Description:      options.Description,
				ChartPathOptions: options.ChartPathOptions,
				Values:           options.Values,
				Atomic:           options.Atomic,
				Wait:             options.Wait,
				WaitForJobs:      options.WaitForJobs,
			})
		} else if err != nil {
			return nil, err
		}
		for _, rel := range rels {
			if rel.Info.Status == release.StatusPendingInstall || rel.Info.Status == release.StatusPendingUpgrade || rel.Info.Status == release.StatusPendingRollback {
				// if release status is pending, delete it
				log.Infof("upgrade release %s is already exist, status is %s. delete it now.", options.ReleaseName, rel.Info.Status)
				actionConfig.Releases.Delete(rel.Name, rel.Version)
			}
		}
	}

	client := action.NewUpgrade(actionConfig)
	client.DryRun = options.DryRun
	client.Timeout = options.Timeout
	client.Namespace = options.Namespace
	client.Description = options.Description
	client.ResetValues = options.ResetValues
	client.ReuseValues = options.ReuseValues
	client.MaxHistory = options.MaxHistory
	client.Atomic = options.Atomic
	client.Wait = options.Wait
	client.WaitForJobs = options.WaitForJobs

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

	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if options.DependencyUpdate {
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

	if chartRequested.Metadata.Deprecated {
		log.Warnf("This chart %s/%s is deprecated", options.ChartRepo, options.Chart)
	}
	return client.RunWithContext(ctx, options.ReleaseName, chartRequested, options.Values)
}
