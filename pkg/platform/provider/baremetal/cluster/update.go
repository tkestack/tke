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
			log.Infof("skip EnsureRenewCerts because expiration duration(%s) > threshold(%s)", expirationDuration, constants.RenewCertsTimeThreshold)
			return nil
		}

		log.Infof("EnsureRenewCerts for %s", s.Host)
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

		log.Infof("EnsureAPIServerCert for %s", s.Host)
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
