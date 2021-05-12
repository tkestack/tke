/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

// expansionPath is the dir where users put there expansion files. default is data/expansions/
var expansionPath string

// expansionEnabled
var expansionEnabled bool

const envExpansionPath = "EXPANSION_PATH"
const expansionDirName = "expansions/"
const defaultExpansionPath = constants.DataDir + expansionDirName

// expansion layout
const expansionLayoutChartsDir = "charts/"
const expansionLayoutFilesDir = "files/"
const expansionLayoutFilesGeneratedDir = "files_generated/"
const expansionLayoutHooksDir = "hooks/"
const expansionLayoutProviderDir = "provider/"
const expansionLayoutImagesName = "images.tar.gz" //nolint
const expansionLayoutConfDir = "conf/"
const expansionLayoutApplicationsDir = "applications/"

var expansionLayoutMap map[string]string

// expansion resource
const expansionResourceChart = "chart"
const expansionResourceFile = "file"
const expansionResourceFileGenerated = "file_generated"
const expansionResourceHook = "hook"
const expansionResourceProvider = "provider"
const expansionResourceConf = "conf"
const expansionResourceApplication = "application"

// expansion value
const expansionValuesFile = "values.yaml" //nolint

func init() {

	expansionPath = os.Getenv(envExpansionPath)
	if expansionPath == "" {
		expansionPath = defaultExpansionPath
	}
	if fi, err := os.Stat(expansionPath); err == nil {
		if fi.IsDir() {
			expansionEnabled = true
		}
	}

	expansionLayoutMap = make(map[string]string)
	if expansionEnabled {
		expansionLayoutMap[expansionResourceChart] = path.Join(expansionPath, expansionLayoutChartsDir)
		expansionLayoutMap[expansionResourceFile] = path.Join(expansionPath, expansionLayoutFilesDir)
		expansionLayoutMap[expansionResourceFileGenerated] = path.Join(expansionPath, expansionLayoutFilesGeneratedDir)
		expansionLayoutMap[expansionResourceHook] = path.Join(expansionPath, expansionLayoutHooksDir)
		expansionLayoutMap[expansionResourceProvider] = path.Join(expansionPath, expansionLayoutProviderDir)
		expansionLayoutMap[expansionResourceConf] = path.Join(expansionPath, expansionLayoutConfDir)
		expansionLayoutMap[expansionResourceApplication] = path.Join(expansionPath, expansionLayoutApplicationsDir)
	}
}

func (t *TKE) prepareExpansionFiles(ctx context.Context) error {

	if !t.isExpansionEnabled() {
		return nil
	}

	err := t.verifyExpansionFiles()
	if err != nil {
		return err
	}

	err = t.overrideWithExpansionFiles()
	if err != nil {
		return err
	}

	return nil
}

func (t *TKE) overrideWithExpansionFiles() error {

	if t.isExpansionResourceEnabled(expansionResourceHook) {
		// rend and copy 'installer hook files'
		for _, hookFile := range []string{
			constants.PreInstallHook,
			constants.PostClusterReadyHook,
			constants.PostInstallHook,
		} {
			fileName := strings.TrimPrefix(hookFile, constants.HooksDir)
			expansionHookFile := path.Join(t.getExpansionResourcePath(expansionResourceHook), fileName)
			fi, err := os.Stat(expansionHookFile)
			if err == nil && !fi.IsDir() {
				err = t.expansionRendAndCopy(expansionHookFile, hookFile, 0755, false)
				if err != nil {
					return fmt.Errorf("rend hook file error. %v => %v, %v", expansionHookFile, hookFile, err)
				}
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

func (t *TKE) verifyExpansionFiles() error {

	if !t.isExpansionEnabled() {
		return nil
	}

	// TODO: verify tke.json with expansion specified files and paths

	return nil
}

func (t *TKE) isExpansionEnabled() bool {
	return expansionEnabled
}

func (t *TKE) getExpansionResourcePath(resource string) string {
	return expansionLayoutMap[resource]
}

func (t *TKE) isExpansionResourceEnabled(resource string) bool {

	resourcePath := t.getExpansionResourcePath(resource)
	if resourcePath == "" {
		return false
	}

	if fi, err := os.Stat(resourcePath); err == nil {
		if fi.IsDir() {
			return true
		}
	}

	return false
}

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
