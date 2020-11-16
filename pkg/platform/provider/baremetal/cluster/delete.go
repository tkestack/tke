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

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/util/mark"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
)

func (p *Provider) EnsureCleanClusterMark(ctx context.Context, c *typesv1.Cluster) error {
	if clientset, err := c.Clientset(); err == nil {
		mark.Delete(ctx, clientset)
	}
	return nil
}

func (p *Provider) EnsureDownScaling(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.ScalingMachines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		_, err = machineSSH.CombinedOutput(`kubeadm reset -f`)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}
	return nil
}
