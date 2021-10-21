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
	defaultQPS     = 100
	defaultBurst   = 200
)

type Cluster struct {
	*platform.Cluster
	ClusterCredential *platform.ClusterCredential
}

func (c *Cluster) Clientset() (kubernetes.Interface, error) {
	config, err := c.RESTConfig(&rest.Config{})
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func (c *Cluster) ClientsetForBootstrap() (kubernetes.Interface, error) {
	config, err := c.RESTConfigForBootstrap(&rest.Config{})
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (c *Cluster) RESTConfigForBootstrap(config *rest.Config) (*rest.Config, error) {
	host, err := c.HostForBootstrap()
	if err != nil {
		return nil, err
	}
	config.Host = host

	return c.RESTConfig(config)
}

func (c *Cluster) setRESTConfigDefaults(config *rest.Config) error {
	if config.Host == "" {
		host, err := c.Host()
		if err != nil {
			return err
		}
		config.Host = host
	}
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.QPS == 0 {
		config.QPS = defaultQPS
	}
	if config.Burst == 0 {
		config.Burst = defaultBurst
	}

	if c.ClusterCredential != nil && c.ClusterCredential.CACert != nil {
		config.TLSClientConfig.CAData = c.ClusterCredential.CACert
	} else {
		config.TLSClientConfig.Insecure = true
	}

	return nil
}

func (c *Cluster) RESTConfig(config *rest.Config) (*rest.Config, error) {
	err := c.setRESTConfigDefaults(config)
	if err != nil {
		return nil, err
	}

	if c.ClusterCredential != nil && c.ClusterCredential.ClientCert != nil && c.ClusterCredential.ClientKey != nil {
		config.TLSClientConfig.CertData = c.ClusterCredential.ClientCert
		config.TLSClientConfig.KeyData = c.ClusterCredential.ClientKey
	}

	if c.ClusterCredential != nil && c.ClusterCredential.Token != nil {
		config.BearerToken = *c.ClusterCredential.Token
	}

	if c.ClusterCredential != nil {
		config.Impersonate.UserName = c.ClusterCredential.Impersonate
		config.Impersonate.Groups = c.ClusterCredential.ImpersonateGroups
		config.Impersonate.Extra = c.ClusterCredential.ImpersonateUserExtra.ExtraToHeaders()
	}

	return config, nil
}

func (c *Cluster) RESTConfigForClientX509(config *rest.Config, clientCertData []byte,
	clientKeyData []byte) (*rest.Config, error) {
	err := c.setRESTConfigDefaults(config)
	if err != nil {
		return nil, err
	}

	config.TLSClientConfig.CertData = clientCertData
	config.TLSClientConfig.KeyData = clientKeyData

	return config, nil
}

func (c *Cluster) HostForBootstrap() (string, error) {
	for _, one := range c.Status.Addresses {
		if one.Type == platform.AddressReal {
			return net.JoinHostPort(one.Host, fmt.Sprintf("%d", one.Port)), nil
		}
	}

	return "", errors.New("can't find bootstrap address")
}
