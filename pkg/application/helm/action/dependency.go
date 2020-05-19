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
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
)

const garbageCollectOlderThan time.Duration = time.Minute * 5

// DependencyUpdate update chart's dependency
func (c *Client) DependencyUpdate(chartPath string, getter []getter.Provider, settings *cli.EnvSettings, verify bool, keyring string) error {
	// Garbage collect before the dependency update so that
	// anonymous files from previous runs are cleared, with
	// a safe guard time offset to not touch any files in
	// use.
	// The path is appointed by chartrepo, see in https://github.com/helm/helm/blob/v3.2.0/pkg/repo/chartrepo.go#L78
	garbageCollect(helmpath.CachePath("repository"), garbageCollectOlderThan)
	man := &downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        chartPath,
		Keyring:          keyring,
		SkipUpdate:       false,
		Getters:          getter,
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	if verify {
		man.Verify = downloader.VerifyAlways
	}
	return man.Update()
}

// GarbageCollectCacheChartsFile clean cache file, including -charts.txt and -index.yaml
func GarbageCollectCacheChartsFile() {
	garbageCollect(helmpath.CachePath("repository"), garbageCollectOlderThan)
}

// garbageCollect walks over the files in the given path and deletes
// any anonymous index file with a mod time older than the given
// duration.
func garbageCollect(path string, olderThan time.Duration) {
	now := time.Now()
	filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
		if err != nil || f.IsDir() {
			return nil
		}
		if strings.HasSuffix(f.Name(), "=-index.yaml") || strings.HasSuffix(f.Name(), "=-charts.txt") {
			if now.Sub(f.ModTime()) > olderThan {
				return os.Remove(p)
			}
		}
		return nil
	})
}
