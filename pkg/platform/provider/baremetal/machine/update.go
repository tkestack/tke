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
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

func (p *Provider) EnsurePreUpgradeHook(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {

	mc := []platformv1.ClusterMachine{
		{
			IP:       machine.Spec.IP,
			Port:     machine.Spec.Port,
			Username: machine.Spec.Username,
			Password: machine.Spec.Password,
		},
	}
	return util.ExcuteCustomizedHook(ctx, cluster, platformv1.HookPreUpgrade, mc)
}

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
		MachineName:            machine.Name,
		MachineIP:              machine.Spec.IP,
		NodeRole:               kubeadm.NodeRoleWorker,
		Version:                cluster.Spec.Version,
		MaxUnready:             cluster.Spec.Features.Upgrade.Strategy.MaxUnready,
		DrainNodeBeforeUpgrade: cluster.Spec.Features.Upgrade.Strategy.DrainNodeBeforeUpgrade,
	}
	logger := log.FromContext(ctx).WithName("Cluster upgrade")
	upgraded, err := kubeadm.UpgradeNode(machineSSH, clientset, p.platformClient, logger, cluster, option)
	if err != nil {
		return err
	}
	if !upgraded {
		return nil
	}

	err = kubeadm.RemoveUpgradeLabel(p.platformClient, machine)
	if err != nil {
		return err
	}

	err = kubeadm.MarkNextUpgradeWorkerNode(clientset, p.platformClient, option.Version, cluster.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsurePostUpgradeHook(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {

	mc := []platformv1.ClusterMachine{
		{
			IP:       machine.Spec.IP,
			Port:     machine.Spec.Port,
			Username: machine.Spec.Username,
			Password: machine.Spec.Password,
		},
	}
	return util.ExcuteCustomizedHook(ctx, cluster, platformv1.HookPostUpgrade, mc)
}
