/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"tkestack.io/tke/pkg/util/log"
)

// RepoUpdateOptions is the options for chart repo update.
type RepoUpdateOptions struct {
	ChartPathOptions
}

// RepoUpdate is the action for chart repo update.
func (c *Client) RepoUpdate(options *RepoUpdateOptions) (entries map[string]repo.ChartVersions, err error) {
	settings, err := NewSettings(options.ChartRepo)
	if err != nil {
		return nil, err
	}

	repoName := path.Base(options.ChartRepo)

	cfg := repo.Entry{
		Name:                  repoName,
		URL:                   options.RepoURL,
		Username:              options.Username,
		Password:              options.Password,
		InsecureSkipTLSverify: true,
	}

	// if use repo.NewChartRepository, repo CachePath is helmpath.CachePath("repository"), but not settings.RepositoryCache
	// r, err := repo.NewChartRepository(&cfg, getter.All(settings))
	// if err != nil {
	// 	return nil, err
	// }
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	client, err := getter.All(settings).ByScheme(u.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not find protocol handler for: %s", u.Scheme)
	}

	r := &repo.ChartRepository{
		Config:    &cfg,
		IndexFile: repo.NewIndexFile(),
		Client:    client,
		CachePath: settings.RepositoryCache,
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return nil, fmt.Errorf("Unable to get an update from the %q chart repository (%s):\n\t%s", repoName, options.RepoURL, err.Error())
	}
	log.Infof("Successfully got an update from the %q chart repository", repoName)

	// Next, we need to load the index, and actually look up the chart.
	idxFile := filepath.Join(settings.RepositoryCache, helmpath.CacheIndexFile(repoName))
	i, err := repo.LoadIndexFile(idxFile)
	if err != nil {
		return nil, fmt.Errorf("no cached repo found for the %q chart repository (%s):\n\t%s", repoName, options.RepoURL, err.Error())
	}

	return i.Entries, nil
}
