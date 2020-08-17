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
	"strconv"
	"strings"
	"time"

	securejoin "github.com/cyphar/filepath-securejoin"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/helmpath"
	"tkestack.io/tke/pkg/util/file"
)

// NewSettings return  EnvSettings that describes all of the environment settings.
func NewSettings(repoDir string) (*cli.EnvSettings, error) {
	repoDir = strings.Trim(repoDir, "/")
	registryConfig, err := securejoin.SecureJoin(repoDir, "registry.json")
	if err != nil {
		return nil, err
	}
	repositoryConfig, err := securejoin.SecureJoin(repoDir, "repositories.yaml")
	if err != nil {
		return nil, err
	}
	repositoryCache, err := securejoin.SecureJoin("repository", repoDir)
	if err != nil {
		return nil, err
	}
	env := &cli.EnvSettings{
		PluginsDirectory: helmpath.DataPath("plugins"),
		RegistryConfig:   helmpath.ConfigPath(registryConfig),
		RepositoryConfig: helmpath.ConfigPath(repositoryConfig),
		RepositoryCache:  helmpath.CachePath(repositoryCache),
	}
	return env, nil
}

// ExpandFile expand existed file to destDir
func ExpandFile(srcFile, destDir string) (string, error) {
	if srcFile != "" && file.IsFile(srcFile) {
		temp, err := securejoin.SecureJoin(destDir, strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			return "", err
		}
		err = chartutil.ExpandFile(temp, srcFile)
		if err != nil {
			return "", err
		}
		return temp, nil
	}
	return "", fmt.Errorf("file not exist")
}
