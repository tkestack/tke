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

package cluster

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"

	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/util/mark"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"
)

func (p *Provider) EnsureCleanClusterMark(ctx context.Context, c *typesv1.Cluster) error {
	if clientset, err := c.Clientset(); err == nil {
		mark.Delete(ctx, clientset)
	}
	return nil
}

func (p *Provider) EnsureRemoveETCDMember(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.ScalingMachines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		err = kubeadm.Reset(machineSSH, "remove-etcd-member")
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) EnsureRemoveNode(ctx context.Context, c *v1.Cluster) error {
	client, err := c.Clientset()
	if err != nil {
		return err
	}
	for _, machine := range c.Spec.ScalingMachines {
		node, err := apiclient.GetNodeByMachineIP(ctx, client, machine.IP)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			return nil
		}
		err = client.CoreV1().Nodes().Delete(context.Background(), node.Name, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
		log.FromContext(ctx).Info("deleteNode done")
	}
	return nil
}

func (p *Provider) EnsureRemoveMachine(ctx context.Context, c *v1.Cluster) error {
	log.FromContext(ctx).Info("delete machine start")
	fieldSelector := fields.OneTermEqualSelector("spec.clusterName", c.Name).String()
	machineList, err := p.PlatformClient.Machines().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return err
	}
	if len(machineList.Items) == 0 {
		return nil
	}
	for _, machine := range machineList.Items {
		if err := p.PlatformClient.Machines().Delete(ctx, machine.Name, metav1.DeleteOptions{}); err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}

		if err = wait.PollImmediate(5*time.Second, 5*time.Minute, waitForMachineDelete(ctx, p.PlatformClient, machine.Name)); err != nil {
			return err
		}
	}

	log.FromContext(ctx).Info("delete machine done")

	return nil
}

func waitForMachineDelete(ctx context.Context, c platformv1client.PlatformV1Interface, machineName string) wait.ConditionFunc {
	return func() (done bool, err error) {

		if _, err := c.Machines().Get(ctx, machineName, metav1.GetOptions{}); err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
		}

		return false, nil
	}
}
