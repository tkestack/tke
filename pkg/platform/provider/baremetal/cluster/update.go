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
	"time"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	certutil "k8s.io/client-go/util/cert"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/log"
)

func (p *Provider) EnsureRenewCerts(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		data, err := s.ReadFile(constants.APIServerCertName)
		if err != nil {
			return err
		}
		certs, err := certutil.ParseCertsPEM(data)
		if err != nil {
			return err
		}
		expirationDuration := time.Until(certs[0].NotAfter)
		if expirationDuration > constants.RenewCertsTimeThreshold {
			log.FromContext(ctx).Info("Skip EnsureRenewCerts because expiration duration > threshold",
				"duration", expirationDuration.String(), "threshold", constants.RenewCertsTimeThreshold.String())
			return nil
		}

		log.FromContext(ctx).Info("EnsureRenewCerts", "nodeName", s.Host)
		err = kubeadm.RenewCerts(s)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
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
				return nil
			}
		}

		log.FromContext(ctx).Info("EnsureAPIServerCert", "nodeName", s.Host)
		for _, file := range []string{constants.APIServerCertName, constants.APIServerKeyName} {
			s.CombinedOutput(fmt.Sprintf("rm -f %s", file))
		}

		err = kubeadm.Init(s, kubeadmConfig, "certs apiserver")
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

func (p *Provider) EnsureUpgradeControlPlaneNode(ctx context.Context, c *v1.Cluster) error {
	client, err := c.Clientset()
	if err != nil {
		return err
	}
	option := kubeadm.UpgradeOption{
		NodeRole:   kubeadm.NodeRoleMaster,
		Version:    c.Spec.Version,
		MaxUnready: c.Spec.Upgrade.Strategy.MaxUnready,
	}
	for i, machine := range c.Spec.Machines {
		option.MachineName = machine.Username
		option.NodeName = machine.IP
		option.BootstrapNode = i == 0
		s, err := machine.SSH()
		if err != nil {
			return err
		}
		upgraded, err := kubeadm.UpgradeNode(s, client, p.platformClient, option)
		if err != nil {
			return err
		}

		// Label next node when upgraded all master nodes and upgrade mode is auto.
		if upgraded && c.Spec.Upgrade.Mode == platformv1.UpgradeModeAuto && i == len(c.Spec.Machines)-1 {
			err = kubeadm.MarkNextUpgradeWorkerNode(client, p.platformClient, option.Version)
			if err != nil {
				return err
			}
		}
	}
	c.Status.Phase = platformv1.ClusterRunning

	return nil
}
