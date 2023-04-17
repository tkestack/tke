/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package installer

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
)

const expansionDirName = "expansions/"
const defaultExpansionDir = constants.DataDir + expansionDirName

// expansion layout
const expansionLayoutHooksDir = "hooks/"

func (t *TKE) initExpansion(ctx context.Context) error {

	if t.Config.CustomExpansionDir == "" {
		t.Config.CustomExpansionDir = defaultExpansionDir
		t.backup()
	}
	if t.Config.CustomChartsName == "" {
		t.Config.CustomChartsName = constants.DefaultCustomChartsName
		t.backup()
	}

	err := t.verifyExpansionPath()
	if err != nil {
		return err
	}
	t.log.Infof("TKEStack installer expansion enabled, with expansion path: %v", t.Config.CustomExpansionDir)

	return nil
}

func (t *TKE) verifyExpansionPath() error {

	_, err := os.Stat(t.Config.CustomExpansionDir)
	return err
}

func (t *TKE) prepareExpansionFiles(ctx context.Context) error {

	err := t.verifyExpansionFiles()
	if err != nil {
		return err
	}

	err = t.overrideWithExpansionFiles()
	if err != nil {
		return err
	}

	// TODO: verify chart files VS platformApps

	return nil
}

func (t *TKE) verifyExpansionFiles() error {

	// TODO: verify tke.json with expansion specified files and paths

	return nil
}

func (t *TKE) overrideWithExpansionFiles() error {

	// rend and copy 'installer hook files'
	hooksDir := path.Join(t.Config.CustomExpansionDir, expansionLayoutHooksDir)
	for _, hookFile := range []string{
		constants.PreInstallHook,
		constants.PostClusterReadyHook,
		constants.PostInstallHook,
	} {
		fileName := strings.TrimPrefix(hookFile, constants.HooksDir)
		expansionHookFile := path.Join(hooksDir, fileName)
		fi, err := os.Stat(expansionHookFile)
		if err == nil && !fi.IsDir() {
			err = t.expansionRendAndCopy(expansionHookFile, hookFile, 0755, false)
			if err != nil {
				return fmt.Errorf("rend hook file error. %v => %v, %v", expansionHookFile, hookFile, err)
			}
		}
	}

	// TODO: rend and copy 'provider files'

	// TODO: rend and copy 'copy files'

	return nil
}

func (t *TKE) expansionRendAndCopy(from, to string, perm os.FileMode, needRend bool) error {

	fi, err := os.Stat(from)
	if err != nil || fi.IsDir() {
		return fmt.Errorf("state from file error, %v", from)
	}

	//nolint
	if needRend {
		// TODO: rend file with tke.json and expansion.yaml
	}

	_ = os.MkdirAll(path.Dir(to), 0755)

	return copyFile(from, to, perm)
}

// importExpansionCharts uploads expansion charts
// func (t *TKE) importExpansionCharts(ctx context.Context) error {
// 	chartPath := path.Join(t.Config.CustomExpansionDir, t.Config.CustomChartsName)
// 	if info, err := os.Stat(chartPath); err != nil {
// 		if os.IsNotExist(err) {
// 			// no chart file found means no charts.
// 			return nil
// 		}
// 		return fmt.Errorf("locate chart file get system error %v,%v", chartPath, err)
// 	} else if info.IsDir() {
// 		return fmt.Errorf("chart file is a directory %v", chartPath)
// 	}

// 	// TODO: (workaround ) client init will only be called in createCluster.
// 	// If we SKIP createCluster step, all client call will be panic
// 	if t.registryClient == nil {
// 		err := t.initDataForDeployTKE()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	err := t.pushCharts(ctx, chartPath, constants.DefaultTeantID, constants.DefaultExpansionChartGroupName)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }

func copyFile(src, dst string, perm os.FileMode) error {

	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if perm == 0 {
		fi, _ := os.Stat(src)
		perm = fi.Mode()
	}

	err = ioutil.WriteFile(dst, input, perm)
	if err != nil {
		return err
	}

	return nil
}
