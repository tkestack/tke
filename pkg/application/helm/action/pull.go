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
	"fmt"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/registry"

	"helm.sh/helm/v3/pkg/action"
	// "helm.sh/helm/v3/pkg/downloader"
	// "helm.sh/helm/v3/pkg/getter"
	// "helm.sh/helm/v3/pkg/repo"
)

// PullOptions is the options for pulling a chart.
type PullOptions struct {
	ChartPathOptions

	DestDir string
}

// Pull is the action for pulling a chart.
func (c *Client) Pull(options *PullOptions) (string, error) {
	actionConfig := new(action.Configuration)
	var err error
	actionConfig.RegistryClient, err = registry.NewClient()
	if err != nil {
		return "", err
	}
	client := action.NewPullWithOpts(action.WithConfig(actionConfig))
	settings, err := NewSettings(options.ChartRepo)
	if err != nil {
		return "", err
	}
	client.Settings = settings
	client.DestDir = options.DestDir
	if client.DestDir == "" {
		client.DestDir = settings.RepositoryCache
	}
	client.Untar = false

	options.ChartPathOptions.ApplyTo(&client.ChartPathOptions)

	GarbageCollectCacheChartsFile()
	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
		return "", err
	}
	_, err = client.Run(options.Chart)
	if err != nil {
		return "", err
	}

	destfile := filepath.Join(client.DestDir, fmt.Sprintf("%s-%s.tgz", options.Chart, options.Version))
	// get file name
	// var out strings.Builder
	// cd := downloader.ChartDownloader{
	// 	Out:     &out,
	// 	Keyring: client.Keyring,
	// 	Verify:  downloader.VerifyNever,
	// 	Getters: getter.All(client.Settings),
	// 	Options: []getter.Option{
	// 		getter.WithBasicAuth(client.Username, client.Password),
	// 		getter.WithTLSClientConfig(client.CertFile, client.KeyFile, client.CaFile),
	// 	},
	// 	RepositoryConfig: client.Settings.RepositoryConfig,
	// 	RepositoryCache:  client.Settings.RepositoryCache,
	// }
	// if client.Verify {
	// 	cd.Verify = downloader.VerifyAlways
	// } else if client.VerifyLater {
	// 	cd.Verify = downloader.VerifyLater
	// }

	// chartRef := options.Chart
	// if client.RepoURL != "" {
	// 	chartURL, err := repo.FindChartInAuthRepoURL(client.RepoURL, client.Username, client.Password, chartRef, client.Version, client.CertFile, client.KeyFile, client.CaFile, getter.All(client.Settings))
	// 	if err != nil {
	// 		return out.String(), err
	// 	}
	// 	chartRef = chartURL
	// }
	// u, err := cd.ResolveChartVersion(chartRef, client.Version)
	// if err != nil {
	// 	return "", err
	// }
	// name := filepath.Base(u.Path)
	// destfile := filepath.Join(client.DestDir, name)
	return destfile, nil
}
