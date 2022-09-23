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

package types

import (
	"errors"
	"fmt"
	"net"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"tkestack.io/tke/api/platform"
)

const (
	defaultTimeout = 30 * time.Second
	defaultQPS     = -1
	defaultBurst   = 200
)

type Cluster struct {
	*platform.Cluster
	ClusterCredential *platform.ClusterCredential
	restConfig        *rest.Config
}

func (c *Cluster) setRESTConfigDefaults() error {
	if c.restConfig == nil {
		c.restConfig = &rest.Config{}
	}
	if c.restConfig.Host == "" {
		host, err := c.Host()
		if err != nil {
			return err
		}
		c.restConfig.Host = host
	}
	if c.restConfig.Timeout == 0 {
		c.restConfig.Timeout = defaultTimeout
	}
	if c.restConfig.QPS == 0 {
		c.restConfig.QPS = defaultQPS
	}
	if c.restConfig.Burst == 0 {
		c.restConfig.Burst = defaultBurst
	}

	return nil
}

func (c *Cluster) RESTConfig() (*rest.Config, error) {
	err := c.setRESTConfigDefaults()
	if err != nil {
		return nil, err
	}
	return c.restConfig, nil
}

func (c *Cluster) GetMainIP() string {
	mainIP := c.Spec.Machines[0].IP
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil {
			mainIP = c.Spec.Features.HA.TKEHA.VIP
		}
		if c.Spec.Features.HA.ThirdPartyHA != nil {
			mainIP = c.Spec.Features.HA.ThirdPartyHA.VIP
		}
	}
	return mainIP
}

func (c *Cluster) GetDirectorIP() string {
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil {
			return c.Spec.Features.HA.TKEHA.DirectorIP
		}
	}
	return ""
}

func (c *Cluster) GetInnerVIPs() []string {
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil && c.Spec.Features.HA.TKEHA.VIPPool != nil {
			return c.Spec.Features.HA.TKEHA.VIPPool.Inner
		}
	}
	return nil
}

func (c *Cluster) GetMasterIPs() []string {
	var IPs []string
	for _, machine := range c.Spec.Machines {
		IPs = append(IPs, machine.IP)
	}
	return IPs
}

func (c *Cluster) Clientset() (kubernetes.Interface, error) {
	config, err := c.RESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func (c *Cluster) ClientsetForBootstrap() (kubernetes.Interface, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (c *Cluster) RESTConfigForBootstrap() (*rest.Config, error) {
	err := c.setRESTConfigDefaults()
	if err != nil {
		return nil, err
	}
	host, err := c.HostForBootstrap()
	if err != nil {
		return nil, err
	}
	configCopy := *c.restConfig
	configCopy.Host = host

	return &configCopy, nil
}

func (c *Cluster) RESTConfigForClientX509(clientCertData []byte,
	clientKeyData []byte) (*rest.Config, error) {
	err := c.setRESTConfigDefaults()
	if err != nil {
		return nil, err
	}
	configCopy := *c.restConfig

	configCopy.TLSClientConfig.CertData = clientCertData
	configCopy.TLSClientConfig.KeyData = clientKeyData

	return &configCopy, nil
}

func (c *Cluster) HostForBootstrap() (string, error) {
	for _, one := range c.Status.Addresses {
		if one.Type == platform.AddressReal {
			return net.JoinHostPort(one.Host, fmt.Sprintf("%d", one.Port)), nil
		}
	}

	return "", errors.New("can't find bootstrap address")
}

func (c *Cluster) RegisterRestConfig(config *rest.Config) {
	c.restConfig = config
}
