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

package net

import "net"

// GetSourceIP return srouce ip to target ip.
func GetSourceIP(target string) (string, error) {
	conn, err := net.Dial("udp", target)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	a := conn.LocalAddr().(*net.UDPAddr)

	return a.IP.String(), nil
}

// InterfaceHasAddr use to check whether host has the specified addr.
func InterfaceHasAddr(addr string) (bool, error) {
	addrs, err := InterfaceAddrs()
	if err != nil {
		return false, err
	}

	for _, one := range addrs {
		if addr == one {
			return true, nil
		}
	}

	return false, nil
}

// InterfaceAddrs returns a list of the system's unicast interface
// addresses in string.
func InterfaceAddrs() ([]string, error) {
	var result []string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		result = append(result, addr.String())
	}

	return result, nil
}
