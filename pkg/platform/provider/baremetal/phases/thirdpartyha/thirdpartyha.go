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

package thirdpartyha

import (
	"fmt"
	"strings"

	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
)

type Option struct {
	IP    string
	VIP   string
	VPort int32
}

func Clear(s ssh.Interface, option *Option) {
	for {
		cmd := fmt.Sprintf("iptables -w 30 -t nat -D OUTPUT -p tcp -d %s --dport %d -j REDIRECT --to-ports 6443",
			option.VIP, option.VPort)
		_, err := s.CombinedOutput(cmd)
		log.Info(fmt.Sprintf("delete iptables %s err:%s", cmd, err))
		if err != nil {
			break
		}
	}
}

// Install solve request roll back problem by rediect request to local host
func Install(s ssh.Interface, option *Option) error {
	cmd := fmt.Sprintf("iptables -w 30 -t nat -C OUTPUT -p tcp -d %s --dport %d -j REDIRECT --to-ports 6443",
		option.VIP, option.VPort)
	_, err := s.CombinedOutput(cmd)
	if err != nil && strings.Contains(err.Error(), "rule exist") { // redirect request to local port
		cmd := fmt.Sprintf("iptables -w 30 -t nat -I OUTPUT -p tcp -d %s --dport %d -j REDIRECT --to-ports 6443",
			option.VIP, option.VPort)
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
		}
	}

	log.Info("redirect vip:vport request to local host and local port.", log.String("node", option.IP))

	return nil
}
