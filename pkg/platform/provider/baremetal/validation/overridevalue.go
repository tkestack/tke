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

package validation

// StorageOverrideValue contains all the storage info used to validate from localization OverrideConfig
type StorageOverrideValue struct {
	Global                       StorageGlobalCfg             `json:"global"`
	NfsSubdirExternalProvisioner NfsSubdirExternalProvisioner `json:"nfs-subdir-external-provisioner"`
	CephCsiCephfs                CephCsiCephfs                `json:"ceph-csi-cephfs"`
}
type StorageGlobalCfg struct {
	EnableNFS    bool `json:"enableNFS"`
	EnableCephFS bool `json:"enableCephFS"`
}
type Nfs struct {
	Server string `json:"server"`
	Path   string `json:"path"`
}
type NfsSubdirExternalProvisioner struct {
	Nfs Nfs `json:"nfs"`
}
type CsiConfig struct {
	ClusterID string   `json:"clusterID"`
	Monitors  []string `json:"monitors"`
}
type StorageClass struct {
	ClusterID string `json:"clusterID"`
	FsName    string `json:"fsName"`
}
type CephCsiSecret struct {
	AdminID  string `json:"adminID"`
	AdminKey string `json:"adminKey"`
}
type CephCsiCephfs struct {
	CsiConfig    []CsiConfig   `json:"csiConfig"`
	StorageClass StorageClass  `json:"storageClass"`
	Secret       CephCsiSecret `json:"secret"`
}

type StorageInfo struct {
	EnableNFS    bool
	EnableCephfs bool
	Nfs          Nfs
	CsiConfig    CsiConfig
}
