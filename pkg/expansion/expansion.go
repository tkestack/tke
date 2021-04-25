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
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

// Driver is the model to handle expansion layout.
type Driver struct {
	log               log.Logger
	ExpansionName     string `json:"expansion_name" yaml:"expansion_name"`
	ExpansionVersion  string `json:"expansion_version" yaml:"expansion_version"`
	RegistryNamespace string `json:"registry_namespace,omitempty" yaml:"registry_namespace" comment:"default is expansionName"`
	// TODO: not designed yet
	K8sVersion string `json:"k8s_version" yaml:"k8s_version"`
	Operator   string `json:"operator" yaml:"operator"`
	Values     map[string]string
	Charts     []string `json:"charts" yaml:"charts"`
	Files      []string `json:"files" yaml:"files"`
	Hooks      []string `json:"hooks" yaml:"hooks"`
	Provider   []string `json:"provider" yaml:"provider"`
	// TODO: support application expansion
	Applications                    []string `json:"applications" yaml:"applications"`
	Images                          []string `json:"images" yaml:"images"`
	globalKubeconfig                []byte   //nolint
	InstallerSkipSteps              []string `json:"installer_skip_steps" yaml:"installer_skip_steps"`
	CreateClusterSkipConditions     []string `json:"create_cluster_skip_conditions" yaml:"create_cluster_skip_conditions"`
	CreateClusterDelegateConditions []string `json:"create_cluster_delegate_conditions" yaml:"create_cluster_delegate_conditions"`
	// TODO: save image lists in case of installer restart to avoid load images again
	ImagesLoaded bool `json:"images_loaded" yaml:"images_loaded"`
	// TODO: if true, will prevent to pass the same expansion to cluster C
	DisablePassThroughToPlatform bool                    `json:"disable_pass_through_to_platform" yaml:"disable_pass_through_to_platform"`
	CreateClusterExtraArgs       *CreateClusterExtraArgs `json:"create_cluster_extra_args" yaml:"create_cluster_extra_args"`
}

// CreateClusterExtraArgs handles all kubernetes args that can be passed through from expansion to TKEStack.
type CreateClusterExtraArgs struct {
	DockerExtraArgs            map[string]string `json:"dockerExtraArgs" yaml:"dockerExtraArgs"`
	KubeletExtraArgs           map[string]string `json:"kubeletExtraArgs" yaml:"kubeletExtraArgs"`
	APIServerExtraArgs         map[string]string `json:"apiServerExtraArgs" yaml:"apiServerExtraArgs"`
	ControllerManagerExtraArgs map[string]string `json:"controllerManagerExtraArgs" yaml:"controllerManagerExtraArgs"`
	SchedulerExtraArgs         map[string]string `json:"schedulerExtraArgs" yaml:"schedulerExtraArgs"`
	Etcd                       *platformv1.Etcd  `json:"etcd" yaml:"etcd"`
}

func (d *Driver) enableOperator() bool {
	return d.Operator != ""
}
func (d *Driver) enableApplications() bool {
	return len(d.Applications) > 0
}
func (d *Driver) enableCharts() bool {
	return len(d.Charts) > 0
}
func (d *Driver) enableFiles() bool {
	return len(d.Files) > 0
}
func (d *Driver) enableHooks() bool {
	return len(d.Hooks) > 0
}
func (d *Driver) enableProvider() bool {
	return len(d.Provider) > 0
}
func (d *Driver) enableImages() bool {
	return len(d.Images) > 0
}
func (d *Driver) enableSkipSteps() bool {
	return len(d.InstallerSkipSteps) > 0
}
func (d *Driver) enableSkipConditions() bool {
	return len(d.CreateClusterSkipConditions) > 0
}
func (d *Driver) enableDelegateConditions() bool {
	return len(d.CreateClusterDelegateConditions) > 0
}

// TODO: not designed yet
//nolint
func (d *Driver) isGlobalKubeconfigReady() bool {
	return len(d.globalKubeconfig) > 0
}

// NewExpansionDriver returns an expansionDriver instance which has all expansion layout items loaded.
func NewExpansionDriver(logger log.Logger) (*Driver, error) {
	driver := &Driver{
		log:                             logger,
		Values:                          make(map[string]string),
		Charts:                          make([]string, 0),
		Files:                           make([]string, 0),
		Hooks:                           make([]string, 0),
		Provider:                        make([]string, 0),
		InstallerSkipSteps:              make([]string, 0),
		CreateClusterDelegateConditions: make([]string, 0),
		CreateClusterSkipConditions:     make([]string, 0),
		CreateClusterExtraArgs: &CreateClusterExtraArgs{
			DockerExtraArgs:            make(map[string]string),
			KubeletExtraArgs:           make(map[string]string),
			APIServerExtraArgs:         make(map[string]string),
			ControllerManagerExtraArgs: make(map[string]string),
			SchedulerExtraArgs:         make(map[string]string),
			//Etcd:                       &platformv1.Etcd{},
		},
	}
	err := driver.scan()
	if err != nil {
		return driver, err
	}
	return driver, nil
}
