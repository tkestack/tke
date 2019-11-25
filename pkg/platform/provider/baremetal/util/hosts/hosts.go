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

package hosts

import (
	"fmt"
	"regexp"
	"runtime"
)

const (
	linuxHostfile   = "/etc/hosts"
	windowsHostFile = "C:/Windows/System32/drivers/etc/hosts"
)

// Hostser for hosts
type Hostser interface {
	Data() ([]byte, error)
	Set(ip string) error
}

func hostFile() string {
	var hostfile string
	if runtime.GOOS == "windows" {
		hostfile = windowsHostFile
	} else {
		hostfile = linuxHostfile
	}
	return hostfile
}

func setHosts(data []byte, host, ip string) ([]byte, error) {
	item := fmt.Sprintf("%s %s", ip, host)
	var re = regexp.MustCompile(fmt.Sprintf(".* %s", host))
	var newData string
	if re.Match(data) {
		newData = re.ReplaceAllString(string(data), item)
	} else {
		newData = fmt.Sprintf("%s\n%s\n", data, item)
	}

	return []byte(newData), nil
}
