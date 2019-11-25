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

package util

import (
	"fmt"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "tkestack.io/tke/api/client/clientset/versioned"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// TODO: Make this configurable?
	storageConfigName = "storage-config"

	// CephConfTemplate is the content template of /etc/ceph/ceph.conf.
	CephConfTemplate = `    [global]
    mon host                     = %s

    auth_cluster_required = cephx
    auth_service_required = cephx
    auth_client_required = cephx

    # Workaround for http://tracker.ceph.com/issues/23446
    fuse_set_user_groups = false`

	// CephAdminKeyringTemplate is the content template of ceph client keyring.
	CephAdminKeyringTemplate = `    [client.admin]
      key = %s
      auid = 0
      caps mds = "allow *"
      caps mon = "allow *"
      caps osd = "allow *"
      caps mgr = "allow *"`

	// CephKeyringFileNameTemplate is the template of ceph client keyring file name.
	CephKeyringFileNameTemplate = "ceph.client.%s.keyring"
)

const (
	// Ceph defines the storage type of ceph
	Ceph = "Ceph"
	// TencentCloud defines the storage type of tencent cloud
	TencentCloud = "TencentCloud"
)

// GetSVInfo returns the information of storage vendor.
func GetSVInfo(client clientset.Interface) (*SVInfo, error) {
	config, err := getConfig(client)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	info := &SVInfo{
		Version: config.ResourceVersion,
		Type:    config.Data["type"],
		cephInfo: cephInfo{
			Monitors: config.Data["monitors"],
			AdminID:  config.Data["adminID"],
			AdminKey: config.Data["adminKey"],
		},
		tencentCloudInfo: tencentCloudInfo{
			SecretID:  config.Data["secretID"],
			SecretKey: config.Data["secretKey"],
		},
	}

	if info.AdminID == "" {
		info.AdminID = "admin"
	}

	var infoErr error
	switch info.Type {
	case TencentCloud:
		infoErr = validateTencentCloudInfo(info)
	case Ceph:
	default:
		infoErr = validateCephInfo(info)
	}
	if infoErr != nil {
		return nil, infoErr
	}

	log.Debugf("Found ceph config: %+v", info)

	return info, nil
}

func validateCephInfo(info *SVInfo) error {
	missedConfigs := []string{}

	if info.Monitors == "" {
		missedConfigs = append(missedConfigs, "monitors")
	}
	if info.AdminKey == "" {
		missedConfigs = append(missedConfigs, "adminKey")
	}

	if len(missedConfigs) > 0 {
		return fmt.Errorf("ceph config %s missed", strings.Join(missedConfigs, ","))
	}

	return nil
}

func validateTencentCloudInfo(info *SVInfo) error {
	missedConfigs := []string{}

	if info.SecretKey == "" {
		missedConfigs = append(missedConfigs, "secretID")
	}
	if info.SecretID == "" {
		missedConfigs = append(missedConfigs, "secretKey")
	}

	if len(missedConfigs) > 0 {
		return fmt.Errorf("tecent cloud config %s missed", strings.Join(missedConfigs, ","))
	}

	return nil
}

// GetSVInfoVersion returns storage vendor info ConfigMap's ResourceVersion.
func GetSVInfoVersion(client clientset.Interface) (string, error) {
	config, err := getConfig(client)
	if err != nil {
		return "", err
	}
	if config == nil {
		return "", nil
	}
	return config.ResourceVersion, nil
}

func getConfig(client clientset.Interface) (*v1.ConfigMap, error) {
	config, err := client.PlatformV1().ConfigMaps().Get(storageConfigName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Ceph cluster not enabled.
			return nil, nil
		}
		return nil, fmt.Errorf("get storage config failed: %v", err)
	}
	return config, nil
}

// SVInfo is a bunch of information of storage vendor.
type SVInfo struct {
	Version string
	Type    string
	cephInfo
	// TODO: Add tencent cloud storage info.\
	tencentCloudInfo
}

type cephInfo struct {
	Monitors string
	AdminID  string
	AdminKey string
}

type tencentCloudInfo struct {
	SecretID  string
	SecretKey string
}
