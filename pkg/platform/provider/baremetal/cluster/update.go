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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	certutil "k8s.io/client-go/util/cert"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

func (p *Provider) EnsureRenewCerts(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		logger := log.FromContext(ctx).WithValues("node", machine.IP)
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		data, err := s.ReadFile(constants.APIServerCertName)
		if err != nil {
			logger.Error(err, "read cert file error")
			return nil
		}
		certs, err := certutil.ParseCertsPEM(data)
		if err != nil {
			logger.Error(err, "ParseCertsPEM error")
			return nil
		}
		expirationDuration := time.Until(certs[0].NotAfter)
		if expirationDuration > constants.RenewCertsTimeThreshold {
			logger.Info("Skip EnsureRenewCerts because expiration duration > threshold",
				"duration", expirationDuration.String(),
				"threshold", constants.RenewCertsTimeThreshold.String(),
			)
			return nil
		}

		logger.Info("RenewCerts doing")
		err = kubeadm.RenewCerts(s)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
		logger.Info("RenewCerts done")
	}

	return nil
}

func (p *Provider) EnsureAPIServerCert(ctx context.Context, c *v1.Cluster) error {
	kubeadmConfig := p.getKubeadmInitConfig(c)
	exptectCertSANs := GetAPIServerCertSANs(c.Cluster)

	needUpload := false
	for _, machine := range c.Spec.Machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		data, err := s.ReadFile(constants.APIServerCertName)
		if err == nil {
			certs, err := certutil.ParseCertsPEM(data)
			if err != nil {
				return err
			}
			actualCertSANs := certs[0].DNSNames
			for _, ip := range certs[0].IPAddresses {
				actualCertSANs = append(actualCertSANs, ip.String())
			}
			if reflect.DeepEqual(funk.IntersectString(actualCertSANs, exptectCertSANs), exptectCertSANs) {
				continue
			}
			log.FromContext(ctx).Info("EnsureAPIServerCert",
				"nodeName", s.Host,
				"exptectCertSANs", exptectCertSANs,
				"actualCertSANs", actualCertSANs,
			)
		}

		var preActions []string
		for _, file := range []string{constants.APIServerCertName, constants.APIServerKeyName} {
			preActions = append(preActions, fmt.Sprintf("rm -f %s", file))
		}

		err = kubeadm.Init(s, kubeadmConfig, "certs apiserver", preActions...)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
		err = kubeadm.RestartContainerByFilter(s, kubeadm.DockerFilterForControlPlane("kube-apiserver"))
		if err != nil {
			return err
		}

		needUpload = true
	}

	if needUpload {
		err := p.EnsureKubeadmInitPhaseUploadConfig(ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsurePreClusterUpgradeHook(ctx context.Context, c *v1.Cluster) error {
	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPreClusterUpgrade, c.Spec.Machines[:1])
}

func (p *Provider) EnsureUpgradeControlPlaneNode(ctx context.Context, c *v1.Cluster) error {
	// check all machines are upgraded before upgrade cluster
	requirement, err := labels.NewRequirement(constants.LabelNodeNeedUpgrade, selection.Exists, []string{})
	if err != nil {
		return err
	}
	machines, err := p.platformClient.Machines().List(context.TODO(), metav1.ListOptions{
		LabelSelector: requirement.String(),
		FieldSelector: fields.OneTermEqualSelector(platformv1.MachineClusterField, c.Name).String(),
	})
	if err != nil {
		return err
	}
	if len(machines.Items) != 0 {
		var itemsName []string
		for _, item := range machines.Items {
			itemsName = append(itemsName, item.Name)
		}
		return fmt.Errorf("some machines, [%s], need to be upgraded", strings.Join(itemsName, ","))
	}

	client, err := c.Clientset()
	if err != nil {
		return err
	}
	option := kubeadm.UpgradeOption{
		NodeRole:               kubeadm.NodeRoleMaster,
		Version:                c.Spec.Version,
		MaxUnready:             c.Spec.Features.Upgrade.Strategy.MaxUnready,
		DrainNodeBeforeUpgrade: c.Spec.Features.Upgrade.Strategy.DrainNodeBeforeUpgrade,
	}
	logger := log.FromContext(ctx).WithName("Cluster upgrade")
	for i, machine := range c.Spec.Machines {
		option.MachineName = machine.Username
		option.MachineIP = machine.IP
		option.BootstrapNode = i == 0
		s, err := machine.SSH()
		if err != nil {
			return err
		}
		upgraded, err := kubeadm.UpgradeNode(s, client, p.platformClient, logger, c, option)
		if err != nil {
			return err
		}

		if i == len(c.Spec.Machines)-1 && upgraded {
			var labelValue string
			if c.Spec.Features.Upgrade.Mode == platformv1.UpgradeModeAuto {
				// set willUpgrade value to all worker node when upgraded all master nodes and upgrade mode is auto.
				labelValue = kubeadm.WillUpgrade
			}
			if err := kubeadm.AddNeedUpgradeLabel(p.platformClient, c.Name, labelValue); err != nil {
				return err
			}
			err = kubeadm.MarkNextUpgradeWorkerNode(client, p.platformClient, option.Version, c.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Provider) EnsurePostClusterUpgradeHook(ctx context.Context, c *v1.Cluster) error {

	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPostClusterUpgrade, c.Spec.Machines[:1])
}
