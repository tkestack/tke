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
	"os"
	"strings"

	"io/ioutil"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"net"
	"path"
)

func HostName() (string, error) {
	nodeName, err := os.Hostname()
	if err != nil {
		return "", err
	}
	nodeName = strings.ToLower(nodeName)
	return nodeName, nil
}

func HostIP() (net.IP, error) {
	nodeIP, err := netutil.ChooseHostInterface()
	if err != nil {
		return net.IP{}, err
	}
	return nodeIP, nil
}

func WriteFile(p string, data []byte, perm os.FileMode) error {
	if _, err := os.Stat(path.Dir(p)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
				return err
			}
		}
	}
	return ioutil.WriteFile(p, data, perm)
}
