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
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"k8s.io/apimachinery/pkg/util/rand"
	platformv1 "tkestack.io/tke/api/platform/v1"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

func GetMasterEndpoint(addresses []platformv1.ClusterAddress) (string, error) {
	var advertise, internal []*platformv1.ClusterAddress
	for _, one := range addresses {
		if one.Type == platformv1.AddressAdvertise {
			advertise = append(advertise, &one)
		}
		if one.Type == platformv1.AddressReal {
			internal = append(internal, &one)
		}
	}

	var address *platformv1.ClusterAddress
	if advertise != nil {
		address = advertise[rand.Intn(len(advertise))]
	} else {
		if internal != nil {
			address = internal[rand.Intn(len(internal))]
		}
	}
	if address == nil {
		return "", errors.New("no advertise or internal address for the cluster")
	}

	return fmt.Sprintf("https://%s:%d", address.Host, address.Port), nil
}

func ExcuteCustomizedHook(ctx context.Context, c *v1.Cluster, htype platformv1.HookType, machines []platformv1.ClusterMachine) error {
	hook := c.Spec.Features.Hooks[htype]
	if hook == "" {
		return nil
	}
	var buffer bytes.Buffer
	if clusterHook := strings.Contains(string(htype), "Cluster"); clusterHook {
		for k, v := range c.GetAnnotations() {
			lineStr := fmt.Sprintf("%s=%s ", k, v)
			buffer.WriteString(lineStr)
		}
		for k, v := range c.GetLabels() {
			lineStr := fmt.Sprintf("%s=%s ", k, v)
			buffer.WriteString(lineStr)
		}
	}

	cmd := strings.Split(hook, " ")[0]
	hook = fmt.Sprintf("%s %s", hook, buffer.String())

	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		machineSSH.Execf("chmod +x %s", cmd)
		_, stderr, exit, err := machineSSH.Exec(hook)
		if err != nil || exit != 0 {
			return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", hook, exit, stderr, err)
		}
	}
	return nil
}

func CleanFlannelInterfaces(cni string) error {
	var err error
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, cni) {
			cmd := exec.Command("ip", "link", "delete", iface.Name)
			if err := cmd.Run(); err != nil {
				log.Errorf("fail to delete link %s : %v", iface.Name, err)
			}
		}
	}
	return err
}
