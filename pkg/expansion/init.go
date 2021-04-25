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

package expansion

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// expansions layout
const defaultExpansionBase = TKEPlatformDataDir + "expansions/"
const defaultExpansionConfName = "expansion.yaml"
const defaultExpansionConf = defaultExpansionBase + defaultExpansionConfName
const expansionValuesFileName = "values.yaml"
const expansionValuesPath = defaultExpansionBase + expansionValuesFileName
const expansionChartsPath = defaultExpansionBase + "charts/"
const expansionFilesPath = defaultExpansionBase + "files/"
const expansionFilesGeneratedPath = defaultExpansionBase + "files_generated/"
const expansionHooksPath = defaultExpansionBase + "hooks/"
const expansionProviderPath = defaultExpansionBase + "provider/"
const expansionImagesPath = defaultExpansionBase + "images.tar.gz"
const expansionConfPath = defaultExpansionBase + "conf/"
const expansionApplicationPath = defaultExpansionBase + "applications/"

// expansion path and separators
const absolutePath = "/"
const expansionFilePathSeparator = "__"

//
const fileSuffixYaml = ".yaml"

// scan looks up for all files/charts/images that specified in expansion.yaml, verifies them and generate flat-named files.
func (d *Driver) scan() error {
	// check dir
	_, err := os.Stat(defaultExpansionBase)
	if err != nil {
		if os.IsNotExist(err) {
			d.log.Info("defaultExpansionBase not exists, skip loading expansions")
			return nil
		}
		d.log.Errorf("scan defaultExpansionBase failed %v", err)
		return err
	}
	err = d.readConfig()
	if err != nil {
		d.log.Errorf("read expansion config failed %v", err)
		return err
	}
	d.log.Infof("extra args: %+v", d.CreateClusterExtraArgs.Etcd)
	err = d.verify()
	if err != nil {
		d.log.Errorf("read expansion failed %v", err)
		return err
	}
	_, err = os.Stat(expansionValuesPath)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Errorf("scan expansionValuesPath failed %v", err)
			return err
		}
	} else {
		d.Values, err = kvFromPath(expansionValuesPath)
		if err != nil {
			return err
		}
	}

	// prepare generated dir
	for _, dir := range []string{
		expansionConfPath,
		//expansionManifestGeneratedPath,
		expansionFilesGeneratedPath,
	} {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			d.log.Errorf("mkdir expansion generate path failed %v,%v", dir, err)
			return err
		}
	}

	// verify
	err = d.verify()
	if err != nil {
		d.log.Errorf("expansion verify failed %v", err)
		return err
	}

	// prepare copy files
	err = d.makeFlatFiles()
	if err != nil {
		d.log.Errorf("expansion makeFlatFiles failed %v", err)
		return err
	}

	return nil
}

// readConfig
func (d *Driver) readConfig() error {
	_, err := os.Stat(defaultExpansionConf)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Info("stat defaultExpansionConf failed %v, %v", defaultExpansionConf, err)
			return err
		}
		return nil
	}
	b, err := ioutil.ReadFile(defaultExpansionConf)
	if err != nil {
		d.log.Errorf("read defaultExpansionConf failed %v, %v", defaultExpansionConf, err)
		return err
	}

	err = yaml.Unmarshal(b, d)
	if err != nil {
		d.log.Errorf("yaml unmarshal defaultExpansionConf file failed %v, %v", defaultExpansionConf, err)
		return err
	}
	return nil
}

func (d *Driver) verify() error {
	// TODO: verify files,charts,hooks,provider,images

	// verify application files
	if d.enableApplications() {
		for _, app := range d.Applications {
			file := expansionApplicationPath + app + fileSuffixYaml
			_, err := os.Stat(file)
			if err != nil {
				return fmt.Errorf("application file not exists %v", file)
			}
		}
	}
	return nil
}

//nolint
func (d *Driver) filter() error {
	// TODO: filter files,charts,hooks,provider,images
	return nil
}

//nolint
func (d *Driver) backup() error {
	// TODO: backup expansion config file
	return nil
}
