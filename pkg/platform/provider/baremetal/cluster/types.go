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

package cluster

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"tkestack.io/tke/pkg/platform/provider/baremetal/config"

	"tkestack.io/tke/pkg/util/ssh"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

type Cluster struct {
	clusterprovider.Cluster
	*config.Config
	SSH map[string]ssh.Interface
}

type Address platformv1.ClusterAddress

func NewCluster(c clusterprovider.Cluster, cfg *config.Config) (*Cluster, error) {
	cluster := &Cluster{}
	cluster.Cluster = c
	cluster.Config = cfg
	cluster.SSH = make(map[string]ssh.Interface)

	for _, m := range c.Spec.Machines {
		sshConfig := &ssh.Config{
			User:       m.Username,
			Host:       m.IP,
			Port:       int(m.Port),
			Password:   string(m.Password),
			PrivateKey: m.PrivateKey,
			PassPhrase: m.PassPhrase,
		}
		s, err := ssh.New(sshConfig)
		if err != nil {
			return nil, errors.Wrap(err, "Create ssh error")
		}
		cluster.SSH[m.IP] = s
	}

	return cluster, nil
}

func (c *Cluster) Clientset() (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        fmt.Sprintf("https://%s:6443", c.Spec.Machines[0].IP), // use the first host because the rest probably not join
		BearerToken: *c.ClusterCredential.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: c.ClusterCredential.CACert,
		},
		Timeout: 5 * time.Second,
	}

	return kubernetes.NewForConfig(restConfig)
}

func (c *Cluster) Address(addrType platformv1.AddressType) *Address {
	for _, one := range c.Status.Addresses {
		if one.Type == addrType {
			a := Address(one)
			return &a
		}
	}

	return nil
}

func (c *Cluster) AddAddress(addrType platformv1.AddressType, host string, port int32) {
	addr := platformv1.ClusterAddress{
		Type: addrType,
		Host: host,
		Port: port,
	}
	// skip same address
	for _, one := range c.Cluster.Status.Addresses {
		if one == addr {
			return
		}
	}
	c.Cluster.Status.Addresses = append(c.Cluster.Status.Addresses, addr)
}

func (c *Cluster) RemoveAddress(addrType platformv1.AddressType) {
	var addrs []platformv1.ClusterAddress
	for _, one := range c.Status.Addresses {
		if one.Type == addrType {
			continue
		}
		addrs = append(addrs, one)
	}
	c.Status.Addresses = addrs
}

func (c *Cluster) SetCondition(newCondition platformv1.ClusterCondition) {
	var conditions []platformv1.ClusterCondition
	exist := false
	for _, condition := range c.Status.Conditions {
		if condition.Type == newCondition.Type {
			exist = true
			if newCondition.Status != condition.Status {
				condition.Status = newCondition.Status
			}
			if newCondition.Message != condition.Message {
				condition.Message = newCondition.Message
			}
			if newCondition.Reason != condition.Reason {
				condition.Reason = newCondition.Reason
			}
			if !newCondition.LastProbeTime.IsZero() && newCondition.LastProbeTime != condition.LastProbeTime {
				condition.LastProbeTime = newCondition.LastProbeTime
			}
			if !newCondition.LastTransitionTime.IsZero() && newCondition.LastTransitionTime != condition.LastTransitionTime {
				condition.LastTransitionTime = newCondition.LastTransitionTime
			}
		}
		conditions = append(conditions, condition)
	}
	if !exist {
		if newCondition.LastProbeTime.IsZero() {
			newCondition.LastProbeTime = metav1.Now()
		}
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = metav1.Now()
		}
		conditions = append(conditions, newCondition)
	}
	c.Status.Conditions = conditions
}

func (ca *Address) String() string {
	return fmt.Sprintf("https://%s:%d", ca.Host, ca.Port)
}
