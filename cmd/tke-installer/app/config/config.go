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

package config

import (
	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/apimachinery/pkg/util/wait"
	applicationv1 "tkestack.io/tke/api/application/v1"
	"tkestack.io/tke/cmd/tke-installer/app/options"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName                 string
	ListenAddr                 string
	NoUI                       bool
	Config                     string
	Force                      bool
	SyncProjectsWithNamespaces bool
	Replicas                   int
	Upgrade                    bool
	PrepareCustomK8sImages     bool
	PrepareCustomCharts        bool
	Kubeconfig                 string
	RegistryUsername           string
	RegistryPassword           string
	RegistryDomain             string
	RegistryNamespace          string
	CustomUpgradeResourceDir   string
	// CustomChartsName
	// when upgrading, it is chart tar file name under data/ directory
	// when installing, it is chart tar file name under data/expansions/ directory
	// default custom.charts.tar.gz
	CustomChartsName string
	// EnableCustomExpansion will enable expansion. default false
	EnableCustomExpansion bool
	// CustomExpansionDir path to expansions. default `data/expansions`
	CustomExpansionDir string
	ExpansionApps      []ExpansionApp
	PlatformApps       []PlatformApp
}
type ExpansionApp struct {
	Name   string
	Enable bool
	Chart  Chart
}

type Chart struct {
	Name           string
	TenantID       string
	ChartGroupName string
	// install options
	Version string
	// install options
	TargetCluster string
	// install options
	TargetNamespace string
	// install options
	// chartutil.ReadValues/ReadValuesFile
	Values chartutil.Values
}

type PlatformApp struct {
	HelmInstallOptions helmaction.InstallOptions
	LocalChartPath     string
	ConditionFunc      wait.ConditionFunc
	Enable             bool
	Installed          bool
	// rawValues: json format or yaml format
	RawValues     string
	RawValuesType applicationv1.RawValuesType
	// values: can specify multiple or separate values: key1=val1,key2=val2
	Values []string
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	log.Infof("Available cluster providers: %v", clusterprovider.Providers())

	return &Config{
		ServerName:                 serverName,
		ListenAddr:                 *opts.ListenAddr,
		NoUI:                       *opts.NoUI,
		Config:                     *opts.Config,
		Force:                      *opts.Force,
		SyncProjectsWithNamespaces: *opts.SyncProjectsWithNamespaces,
		Replicas:                   *opts.Replicas,
		Upgrade:                    *opts.Upgrade,
		PrepareCustomK8sImages:     *opts.PrepareCustomK8sImages,
		PrepareCustomCharts:        *opts.PrepareCustomCharts,
		Kubeconfig:                 *opts.Kubeconfig,
		RegistryUsername:           *opts.RegistryUsername,
		RegistryPassword:           *opts.RegistryPassword,
		RegistryDomain:             *opts.RegistryDomain,
		RegistryNamespace:          *opts.RegistryNamespace,
		CustomUpgradeResourceDir:   *opts.CustomUpgradeResourceDir,
		CustomChartsName:           *opts.CustomChartsName,
	}, nil
}
