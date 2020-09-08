/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package machine

import (
	"context"

	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
)

func (p *Provider) EnsureUpgrade(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	if _, ok := machine.Labels[constants.LabelNodeNeedUpgrade]; !ok {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	clientset, err := cluster.Clientset()
	if err != nil {
		return err
	}

	option := kubeadm.UpgradeOption{
		MachineName: machine.Name,
		MachineIP:   machine.Spec.IP,
		NodeRole:    kubeadm.NodeRoleWorker,
		Version:     cluster.Spec.Version,
		MaxUnready:  cluster.Spec.Upgrade.Strategy.MaxUnready,
	}
	upgraded, err := kubeadm.UpgradeNode(machineSSH, clientset, p.platformClient, option)
	if err != nil {
		return err
	}
	if !upgraded {
		return nil
	}

	// Remove upgrade label
	delete(machine.Labels, constants.LabelNodeNeedUpgrade)

	machine.Status.Phase = platformv1.MachineRunning

	// Label next node when upgraded current nodes and upgrade mode is auto.
	if cluster.Spec.Upgrade.Mode == platformv1.UpgradeModeAuto {
		err = kubeadm.MarkNextUpgradeWorkerNode(clientset, p.platformClient, option.Version)
		if err != nil {
			return err
		}
	}

	return nil
}
