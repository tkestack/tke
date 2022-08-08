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

package v1

import (
	"fmt"
	"math/rand"
	"net"
	"path"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	applicationversiondclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/application/helm/action"
	helmconfig "tkestack.io/tke/pkg/application/helm/config"
)

const (
	defaultTimeout = 30 * time.Second
	defaultQPS     = -1
	defaultBurst   = 200
)

type Cluster struct {
	*platformv1.Cluster
	ClusterCredential   *platformv1.ClusterCredential
	restConfig          *rest.Config
	IsCredentialChanged bool
}

func (c *Cluster) Clientset() (kubernetes.Interface, error) {
	config, err := c.RESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
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

func (c *Cluster) ClientsetForBootstrap() (kubernetes.Interface, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (c *Cluster) HelmClientsetForBootstrap(namespace string) (*action.Client, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	restClientGetter := &helmconfig.RESTClientGetter{RestConfig: config}
	restClientGetter.Namespace = &namespace
	client := action.NewClient("", restClientGetter)
	return client, nil
}

func (c *Cluster) PlatformClientsetForBootstrap() (platformversionedclient.PlatformV1Interface, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	return platformversionedclient.NewForConfig(config)
}

func (c *Cluster) RegistryClientsetForBootstrap() (registryversionedclient.RegistryV1Interface, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	return registryversionedclient.NewForConfig(config)
}

func (c *Cluster) RegistryApplicationForBootstrap() (applicationversiondclient.ApplicationV1Interface, error) {
	config, err := c.RESTConfigForBootstrap()
	if err != nil {
		return nil, err
	}
	return applicationversiondclient.NewForConfig(config)
}

func (c *Cluster) RESTConfigForBootstrap() (*rest.Config, error) {
	err := c.setRESTConfigDefaults()
	if err != nil {
		return nil, err
	}
	configCopy := *c.restConfig

	return &configCopy, nil
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

func (c *Cluster) Host() (string, error) {
	addrs := make(map[platformv1.AddressType][]platformv1.ClusterAddress)
	for _, one := range c.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *platformv1.ClusterAddress
	if len(addrs[platformv1.AddressInternal]) != 0 {
		address = &addrs[platformv1.AddressInternal][rand.Intn(len(addrs[platformv1.AddressInternal]))]
	} else if len(addrs[platformv1.AddressAdvertise]) != 0 {
		address = &addrs[platformv1.AddressAdvertise][rand.Intn(len(addrs[platformv1.AddressAdvertise]))]
	} else {
		if len(addrs[platformv1.AddressReal]) != 0 {
			address = &addrs[platformv1.AddressReal][rand.Intn(len(addrs[platformv1.AddressReal]))]
		}
	}

	if address == nil {
		return "", errors.New("can't find valid address")
	}
	result := net.JoinHostPort(address.Host, fmt.Sprintf("%d", address.Port))
	if address.Path != "" {
		result = path.Join(result, path.Clean(address.Path))
		result = fmt.Sprintf("https://%s", result)
	}

	return result, nil
}

func (c *Cluster) HostForBootstrap() (string, error) {
	for _, one := range c.Status.Addresses {
		if one.Type == platformv1.AddressReal {
			return net.JoinHostPort(one.Host, fmt.Sprintf("%d", one.Port)), nil
		}
	}

	return "", errors.New("can't find bootstrap address")
}

func (c *Cluster) RegisterRestConfig(config *rest.Config) {
	c.restConfig = config
}
