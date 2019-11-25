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

package machine

import (
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/config"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util/apiclient"
	"tkestack.io/tke/pkg/util/ssh"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Machine struct {
	platformv1.Machine
	Cluster           *platformv1.Cluster
	ClusterCredential *platformv1.ClusterCredential
	*config.Config
	ssh.Interface
	ClientSet kubernetes.Interface
}

func NewMachine(m platformv1.Machine, c *platformv1.Cluster, credential *platformv1.ClusterCredential, cfg *config.Config) (*Machine, error) {
	var err error

	machine := &Machine{}
	machine.Machine = m
	machine.Cluster = c
	machine.ClusterCredential = credential

	machine.Config = cfg

	sshConfig := &ssh.Config{
		User:       m.Spec.Username,
		Host:       m.Spec.IP,
		Port:       int(m.Spec.Port),
		Password:   string(m.Spec.Password),
		PrivateKey: m.Spec.PrivateKey,
		PassPhrase: m.Spec.PassPhrase,
	}
	machine.Interface, err = ssh.New(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Create ssh error")
	}

	masterEndpoint, err := util.GetMasterEndpoint(c.Status.Addresses)
	if err != nil {
		return nil, errors.Wrap(err, "GetMasterEndpoint error")
	}
	machine.ClientSet, err = apiclient.GetClientset(masterEndpoint, *credential.Token, credential.CACert)
	if err != nil {
		return nil, errors.Wrap(err, "GetClientset error")
	}

	return machine, nil
}

func (c *Machine) SetCondition(newCondition platformv1.MachineCondition) {
	var conditions []platformv1.MachineCondition
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

func IsGPU(labels map[string]string) bool {
	isGPU := labels["nvidia-device-enable"]
	return isGPU == "enable"
}
