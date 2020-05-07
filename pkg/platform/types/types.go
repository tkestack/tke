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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
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

func GetClusterByName(platformClient internalversion.PlatformInterface, name string) (*Cluster, error) {
	result := new(Cluster)
	cluster, err := platformClient.Clusters().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	result.Cluster = cluster
	if cluster.Spec.ClusterCredentialRef != nil {
		clusterCredential, err := platformClient.ClusterCredentials().Get(cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get cluster's credential error: %w", err)
		}
		result.ClusterCredential = clusterCredential
	}

	return result, nil
}

func GetCluster(platformClient internalversion.PlatformInterface, cluster *platform.Cluster) (*Cluster, error) {
	result := new(Cluster)
	result.Cluster = cluster
	if cluster.Spec.ClusterCredentialRef != nil {
		clusterCredential, err := platformClient.ClusterCredentials().Get(cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get cluster's credential error: %w", err)
		}
		result.ClusterCredential = clusterCredential
	}

	return result, nil
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

func (c *Cluster) RESTConfig(config *rest.Config) (*rest.Config, error) {
	if config.Host == "" {
		host, err := c.Host()
		if err != nil {
			return nil, err
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

	if c.ClusterCredential.CACert != nil {
		config.TLSClientConfig.CAData = c.ClusterCredential.CACert
	} else {
		config.TLSClientConfig.Insecure = true
	}
	if c.ClusterCredential.ClientCert != nil && c.ClusterCredential.ClientKey != nil {
		config.TLSClientConfig.CertData = c.ClusterCredential.ClientCert
		config.TLSClientConfig.KeyData = c.ClusterCredential.ClientKey
	}

	if c.ClusterCredential.Token != nil {
		config.BearerToken = *c.ClusterCredential.Token
	}

	return config, nil
}

func (c *Cluster) HostForBootstrap() (string, error) {
	for _, one := range c.Status.Addresses {
		if one.Type == platform.AddressReal {
			return fmt.Sprintf("%s:%d", one.Host, one.Port), nil
		}
	}

	return "", errors.New("can't find bootstrap address")
}
