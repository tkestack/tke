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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (in *Cluster) SetCondition(newCondition ClusterCondition) {
	var conditions []ClusterCondition

	exist := false

	if newCondition.LastProbeTime.IsZero() {
		newCondition.LastProbeTime = metav1.Now()
	}
	for _, condition := range in.Status.Conditions {
		if condition.Type == newCondition.Type {
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

	return fmt.Sprintf("%s:%d", address.Host, address.Port), nil
}
