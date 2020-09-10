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

package keepalived

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"
)

type Option struct {
	IP              string
	VIP             string
	LoadBalance     bool
	IPVS            bool
	KubernetesSvcIP string
}

var svcPortChain = "KUBE-SVC-NPX46M4PTMTKRN6Y"

// vrid gets vrid from last byte of vip and plus 1 to prevent from vip ends with zero,
// bcs vrid ranges from 1 to 255. if vip ends with 255, then throw error.
func vrid(vip string) (int, error) {
	var lastbyte = net.ParseIP(vip).To16()[net.IPv6len-1]
	if lastbyte >= 255 {
		return 1, fmt.Errorf("unusual vip :%s", vip)
	}

	return int(net.ParseIP(vip).To16()[net.IPv6len-1]) + 1, nil
}

func Install(s ssh.Interface, option *Option) error {
	networkInterface := ssh.GetNetworkInterface(s, option.IP)
	if networkInterface == "" {
		return fmt.Errorf("can't get network interface by %s", option.IP)
	}

	vrid, err := vrid(option.VIP)
	if err != nil {
		return err
	}

	data, err := template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.conf", map[string]interface{}{
		"Interface": networkInterface,
		"VRID":      vrid,
		"VIP":       option.VIP,
	})
	if err != nil {
		return errors.Wrap(err, option.IP)
	}
	err = s.WriteFile(bytes.NewReader(data), constants.KeepavliedConfigFile)
	if err != nil {
		return errors.Wrap(err, option.IP)
	}

	data, err = template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.yaml", map[string]interface{}{
		"Image": images.Get().Keepalived.FullName(),
	})
	if err != nil {
		return errors.Wrap(err, option.IP)
	}
	err = s.WriteFile(bytes.NewReader(data), constants.KeepavlivedManifestFile)
	if err != nil {
		return errors.Wrap(err, option.IP)
	}

	if option.LoadBalance {
		err := installLoadBalanceIfKubeProxyRunning(s, option)
		if err != nil {
			return errors.Wrap(err, option.IP)
		}
	}

	return nil
}

// clearLoadbalanceRuleInIpvsMode delete all ipvs mode relative iptables rules
func clearLoadbalanceRuleInIpvsMode(s ssh.Interface, vip string, chain string, kubernetesSvcIP string) {
	for {
		cmd := fmt.Sprintf("iptables -t nat -D %s -d %s -p tcp --dport 6443 -j DNAT --to-destination %s:443", chain, vip, kubernetesSvcIP)
		_, err := s.CombinedOutput(cmd)
		log.Info(fmt.Sprintf("delete iptables %s err:%s", cmd, err))
		if err != nil {
			break
		}
	}

	for {
		cmd := fmt.Sprintf("iptables -t nat -D %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ", chain, vip)
		_, err := s.CombinedOutput(cmd)
		log.Info(fmt.Sprintf("delete iptables %s err:%s", cmd, err))
		if err != nil {
			break
		}
	}
}

// clearLoadbalanceRuleInIptablesMode delete all iptables mode relative iptables rules
func clearLoadbalanceRuleInIptablesMode(s ssh.Interface, vip string, chain string) {
	for {
		cmd := fmt.Sprintf("iptables -t nat -D %s -d %s -p tcp --dport 6443 -j %s", chain, vip, svcPortChain)
		_, err := s.CombinedOutput(cmd)
		log.Info(fmt.Sprintf("delete iptables %s err:%s", cmd, err))
		if err != nil {
			break
		}
	}

	for {
		cmd := fmt.Sprintf("iptables -t nat -D %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ", chain, vip)
		_, err := s.CombinedOutput(cmd)
		log.Info(fmt.Sprintf("delete iptables %s err:%s", cmd, err))
		if err != nil {
			break
		}
	}
}

func ClearLoadBalance(s ssh.Interface, vip string, kubernetesSvcIP string) {
	chains := []string{"PREROUTING", "OUTPUT"}
	for _, chain := range chains {
		clearLoadbalanceRuleInIpvsMode(s, vip, chain, kubernetesSvcIP)
		clearLoadbalanceRuleInIptablesMode(s, vip, chain)
	}
}

// installLoadBalanceInIpvsModeIfKubeProxyRunning dnat vip:vport to kubernetes svc cluster ip:port and mark the request.
// For more information, see the proposal: https://github.com/tkestack/tke/blob/master/docs/design-proposals/controlplane-ha-loadbalance.md
func installLoadBalanceInIpvsModeIfKubeProxyRunning(s ssh.Interface, option *Option, chain string) error {
	cmd := fmt.Sprintf("ip a | grep %s", option.KubernetesSvcIP)
	stdout, _, exit, err := s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:error %s:stdout %s", cmd, exit, err, stdout)
	}

	if !strings.Contains(stdout, option.KubernetesSvcIP) {
		return nil
	}

	cmd = fmt.Sprintf("iptables -t nat -C %s -d %s -p tcp --dport 6443 -j DNAT --to-destination %s:443",
		chain, option.VIP, option.KubernetesSvcIP)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		cmd := fmt.Sprintf("iptables -t nat -I %s -d %s -p tcp --dport 6443 -j DNAT --to-destination %s:443",
			chain, option.VIP, option.KubernetesSvcIP)
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
		}
	}

	cmd = fmt.Sprintf("iptables -t nat -C %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ",
		chain, option.VIP)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		cmd = fmt.Sprintf("iptables -t nat -I %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ",
			chain, option.VIP)
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
		}
	}

	log.Info("loadblance on ipvs mode created success.", log.String("node", option.IP), log.String("chain", chain))
	return nil
}

// installLoadBalanceInIptableModeIfKubeProxyRunning dnat vip:vport to kubernetes service chain and mark the request.
// kubernetes service chain generated by service name: kubernetes and port: 443 using hash alg
// For more information, see the proposal: https://github.com/tkestack/tke/blob/master/docs/design-proposals/controlplane-ha-loadbalance.md
func installLoadBalanceInIptableModeIfKubeProxyRunning(s ssh.Interface, option *Option, chain string) error {
	cmd := fmt.Sprintf("iptables -t nat -nxL %s", svcPortChain)
	stdout, _, exit, err := s.Exec(cmd)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:error %s:stdout %s", cmd, exit, err, stdout)
	}

	if !strings.Contains(stdout, svcPortChain) {
		return nil
	}

	cmd = fmt.Sprintf("iptables -t nat -C %s -d %s -p tcp --dport 6443 -j %s", chain, option.VIP, svcPortChain)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		cmd = fmt.Sprintf("iptables -t nat -I %s -d %s -p tcp --dport 6443 -j %s", chain, option.VIP, svcPortChain)
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
		}
	}

	cmd = fmt.Sprintf("iptables -t nat -C %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ", chain, option.VIP)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		cmd = fmt.Sprintf("iptables -t nat -I %s -d %s -p tcp --dport 6443 -j KUBE-MARK-MASQ", chain, option.VIP)
		_, err = s.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("run cmd(%s) error:%s", cmd, err)
		}
	}

	log.Info("loadblance on iptable mode created success.", log.String("node", option.IP), log.String("chain", chain))
	return nil
}

func installLoadBalanceIfKubeProxyRunning(s ssh.Interface, option *Option) error {
	chains := []string{"PREROUTING", "OUTPUT"}
	for _, chain := range chains {
		if option.IPVS {
			err := installLoadBalanceInIpvsModeIfKubeProxyRunning(s, option, chain)
			if err != nil {
				return err
			}
		} else {
			err := installLoadBalanceInIptableModeIfKubeProxyRunning(s, option, chain)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
