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

package v1

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"tkestack.io/tke/pkg/util/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/pkg/util/ssh"
)

func (in *ClusterMachine) SSH() (*ssh.SSH, error) {
	sshConfig := &ssh.Config{
		User:        in.Username,
		Host:        in.IP,
		Port:        int(in.Port),
		Password:    string(in.Password),
		PrivateKey:  in.PrivateKey,
		PassPhrase:  in.PassPhrase,
		DialTimeOut: time.Second,
		Retry:       0,
	}
	switch in.Proxy.Type {
	case SSHJumpServer:
		proxy := ssh.JumpServer{}
		proxy.Host = in.Proxy.IP
		proxy.Port = int(in.Proxy.Port)
		proxy.User = in.Proxy.Username
		proxy.Password = string(in.Proxy.Password)
		proxy.PrivateKey = in.Proxy.PrivateKey
		proxy.PassPhrase = in.Proxy.PassPhrase
		proxy.DialTimeOut = time.Second
		proxy.Retry = 0
		sshConfig.Proxy = proxy
	case SOCKS5:
		proxy := ssh.SOCKS5{}
		proxy.Host = in.Proxy.IP
		proxy.Port = int(in.Proxy.Port)
		proxy.DialTimeOut = time.Second
		sshConfig.Proxy = proxy
	}
	return ssh.New(sshConfig)
}

func (in *Cluster) Address(addrType AddressType) *ClusterAddress {
	for _, one := range in.Status.Addresses {
		if one.Type == addrType {
			return &one
		}
	}

	return nil
}

func (in *Cluster) AddAddress(addrType AddressType, host string, port int32) {
	addr := ClusterAddress{
		Type: addrType,
		Host: host,
		Port: port,
	}
	for _, one := range in.Status.Addresses {
		if one == addr {
			return
		}
	}
	in.Status.Addresses = append(in.Status.Addresses, addr)
}

func (in *Cluster) RemoveAddress(addrType AddressType) {
	var addrs []ClusterAddress
	for _, one := range in.Status.Addresses {
		if one.Type == addrType {
			continue
		}
		addrs = append(addrs, one)
	}
	in.Status.Addresses = addrs
}

func (in *Cluster) GetCondition(conditionType string) *ClusterCondition {
	for _, condition := range in.Status.Conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}

	return nil
}

func (in *Cluster) KeepHistory(keepHistory bool, condition ClusterCondition) bool {
	if !keepHistory {
		return false
	}
	return condition.Status == ConditionTrue
}

func (in *Cluster) SetCondition(newCondition ClusterCondition, keepHistory bool) {
	var conditions []ClusterCondition

	exist := false

	if newCondition.LastProbeTime.IsZero() {
		newCondition.LastProbeTime = metav1.Now()
	}
	for _, condition := range in.Status.Conditions {
		if condition.Type == newCondition.Type && !in.KeepHistory(keepHistory, condition) {
			exist = true
			if newCondition.LastTransitionTime.IsZero() {
				newCondition.LastTransitionTime = condition.LastTransitionTime
			}
			condition = newCondition
		}
		conditions = append(conditions, condition)
	}

	if !exist {
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = metav1.Now()
		}
		conditions = append(conditions, newCondition)
	}

	in.Status.Conditions = conditions
	switch newCondition.Status {
	case ConditionFalse:
		in.Status.Reason = newCondition.Reason
		in.Status.Message = newCondition.Message
	default:
		in.Status.Reason = ""
		in.Status.Message = ""
	}
}

func (in *Cluster) Host() (string, error) {
	addrs := make(map[AddressType][]ClusterAddress)
	for _, one := range in.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *ClusterAddress
	if len(addrs[AddressInternal]) != 0 {
		address = &addrs[AddressInternal][rand.Intn(len(addrs[AddressInternal]))]
	} else if len(addrs[AddressAdvertise]) != 0 {
		address = &addrs[AddressAdvertise][rand.Intn(len(addrs[AddressAdvertise]))]
	} else {
		if len(addrs[AddressReal]) != 0 {
			address = &addrs[AddressReal][rand.Intn(len(addrs[AddressReal]))]
		}
	}

	if address == nil {
		return "", errors.New("can't find valid address")
	}

	return net.JoinHostPort(address.Host, fmt.Sprintf("%d", address.Port)), nil
}

func (in *Cluster) AuthzWebhookEnabled() bool {
	// for anyhwere case authz is always enable
	if in.Spec.Type == "Anywhere" {
		return true
	}
	return in.Spec.Features.AuthzWebhookAddr != nil &&
		(in.Spec.Features.AuthzWebhookAddr.Builtin != nil || in.Spec.Features.AuthzWebhookAddr.External != nil)
}

func (in *Cluster) AuthzWebhookExternEndpoint() (string, bool) {
	if in.Spec.Features.AuthzWebhookAddr == nil {
		return "", false
	}

	if in.Spec.Features.AuthzWebhookAddr.External == nil {
		return "", false
	}

	ip := in.Spec.Features.AuthzWebhookAddr.External.IP
	port := int(in.Spec.Features.AuthzWebhookAddr.External.Port)
	return http.MakeEndpoint("https", ip, port, "/auth/authz"), true
}

func (in *Cluster) AuthzWebhookBuiltinEndpoint() (string, bool) {
	if in.Spec.Features.AuthzWebhookAddr == nil {
		return "", false
	}

	if in.Spec.Features.AuthzWebhookAddr.Builtin == nil {
		return "", false
	}

	endPointHost := in.Spec.Machines[0].IP

	// use VIP in HA situation
	if in.Spec.Features.HA != nil {
		if in.Spec.Features.HA.TKEHA != nil {
			endPointHost = in.Spec.Features.HA.TKEHA.VIP
		}
		if in.Spec.Features.HA.ThirdPartyHA != nil {
			endPointHost = in.Spec.Features.HA.ThirdPartyHA.VIP
		}
	}

	return http.MakeEndpoint("https", endPointHost,
		constants.AuthzWebhookNodePort, "/auth/authz"), true
}
