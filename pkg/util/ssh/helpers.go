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

package ssh

import (
	"strconv"
	"strings"
)

// GetNetworkInterface return network interface name by ip
func GetNetworkInterface(s Interface, ip string) string {
	stdout, _, _, _ := s.Execf("ip a | grep '%s' |awk '{print $NF}'", ip)

	return stdout
}

// Timestamp returns target node timestamp.
func Timestamp(s Interface) (int, error) {
	stdout, err := s.CombinedOutput("date +%s")
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(stdout)))
}

// MemoryCapacity returns the machine's total memory from /proc/meminfo.
// Returns the total memory capacity as an uint64 (number of bytes).
func MemoryCapacity(s Interface) (uint64, error) {
	stdout, err := s.CombinedOutput(`grep 'MemTotal:' /proc/meminfo | grep -oP '\d+'`)
	if err != nil {
		return 0, err
	}

	memInKB, err := strconv.ParseUint(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}

	return memInKB * 1024, err
}

// NumCPU returns the number of logical CPUs.
func NumCPU(s Interface) (int, error) {
	stdout, err := s.CombinedOutput(`nproc --all`)
	if err != nil {
		return 0, err
	}

	cpu, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, err
	}

	return cpu, nil
}
