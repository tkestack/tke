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

package clusteraddontype

import (
	"bytes"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/platform"
	cronhpa "tkestack.io/tke/pkg/platform/controller/addon/cronhpa/images"
	helm "tkestack.io/tke/pkg/platform/controller/addon/helm/images"
	ipam "tkestack.io/tke/pkg/platform/controller/addon/ipam/images"
	lbcf "tkestack.io/tke/pkg/platform/controller/addon/lbcf/images"
	logcollector "tkestack.io/tke/pkg/platform/controller/addon/logcollector/images"
	persistentevent "tkestack.io/tke/pkg/platform/controller/addon/persistentevent/images"
	prometheus "tkestack.io/tke/pkg/platform/controller/addon/prometheus/images"
	volumedecorator "tkestack.io/tke/pkg/platform/controller/addon/storage/volumedecorator/images"
	tappcontroller "tkestack.io/tke/pkg/platform/controller/addon/tappcontroller/images"
	csioperator "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/registry/clusteraddontype/assets"
	"tkestack.io/tke/pkg/util/log"
)

// AddonType is a alias name of string.
type AddonType string

// These are valid type of addon.
const (
	// Helm is type for helm addon.
	Helm AddonType = "Helm"
	// PersistentEvent is type for persistent event addon.
	PersistentEvent AddonType = "PersistentEvent"
	// LogCollector is type for log collector.
	LogCollector AddonType = "LogCollector"
	// TappController is type for TappController
	TappController AddonType = "TappController"
	// CSIOperator is type for CSIOperator
	CSIOperator AddonType = "CSIOperator"
	// VolumeDecorator is type for VolumeDecorator
	VolumeDecorator AddonType = "VolumeDecorator"
	// CronHPA is type for CronHPA
	CronHPA AddonType = "CronHPA"
	// Prometheus is type for prometheus addon.
	Prometheus AddonType = "Prometheus"
	// Galaxy-IPAM
	IPAM AddonType = "IPAM"
	// LBCF is type for LBCF
	LBCF AddonType = "LBCF"
)

// Types defines the type of each plugin and the mapping table of the latest
// version number.
var Types = map[AddonType]platform.ClusterAddonType{
	Helm: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(Helm)),
		},
		Type:                  string(Helm),
		Level:                 platform.LevelEnhance,
		LatestVersion:         helm.LatestVersion,
		Description:           description("Helm.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	PersistentEvent: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(PersistentEvent)),
		},
		Type:                  string(PersistentEvent),
		Level:                 platform.LevelEnhance,
		LatestVersion:         persistentevent.LatestVersion,
		Description:           description("PersistentEvent.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	LogCollector: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(LogCollector)),
		},
		Type:                  string(LogCollector),
		Level:                 platform.LevelEnhance,
		LatestVersion:         logcollector.GetLatestVersion(),
		Description:           description("LogCollector.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	TappController: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(TappController)),
		},
		Type:                  string(TappController),
		Level:                 platform.LevelEnhance,
		LatestVersion:         tappcontroller.GetLatestVersion(),
		Description:           description("TappController.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	CSIOperator: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(CSIOperator)),
		},
		Type:                  string(CSIOperator),
		Level:                 platform.LevelBasic,
		LatestVersion:         csioperator.LatestVersion,
		Description:           description("CSIOperator.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	VolumeDecorator: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(VolumeDecorator)),
		},
		Type:                  string(VolumeDecorator),
		Level:                 platform.LevelEnhance,
		LatestVersion:         volumedecorator.LatestVersion,
		Description:           description("VolumeDecorator.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	CronHPA: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(CronHPA)),
		},
		Type:                  string(CronHPA),
		Level:                 platform.LevelEnhance,
		LatestVersion:         cronhpa.GetLatestVersion(),
		Description:           description("CronHPA.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	Prometheus: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(Prometheus)),
		},
		Type:                  string(Prometheus),
		Level:                 platform.LevelBasic,
		LatestVersion:         prometheus.LatestVersion,
		Description:           description("Prometheus.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	IPAM: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(IPAM)),
		},
		Type:                  string(IPAM),
		Level:                 platform.LevelEnhance,
		LatestVersion:         ipam.GetLatestVersion(),
		Description:           description("IPAM.md"),
		CompatibleClusterType: cluster.Providers(),
	},
	LBCF: {
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(string(LBCF)),
		},
		Type:                  string(LBCF),
		Level:                 platform.LevelBasic,
		LatestVersion:         lbcf.GetLatestVersion(),
		Description:           description("LBCF.md"),
		CompatibleClusterType: cluster.Providers(),
	},
}

func description(name string) string {
	var err error
	reader, err := assets.Open(name)
	if err != nil {
		log.Error("Failed to open asset file", log.String("name", name), log.Err(err))
		return ""
	}
	defer func() {
		_ = reader.Close()
	}()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		log.Error("Failed to read asset file", log.String("name", name), log.Err(err))
		return ""
	}
	return buf.String()
}
