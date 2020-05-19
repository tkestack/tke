/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	rest "k8s.io/client-go/rest"
	"tkestack.io/tke/api/client/clientset/versioned/scheme"
	v1 "tkestack.io/tke/api/registry/v1"
)

type RegistryV1Interface interface {
	RESTClient() rest.Interface
	ChartsGetter
	ChartGroupsGetter
	ConfigMapsGetter
	NamespacesGetter
	RepositoriesGetter
}

// RegistryV1Client is used to interact with features provided by the registry.tkestack.io group.
type RegistryV1Client struct {
	restClient rest.Interface
}

func (c *RegistryV1Client) Charts(namespace string) ChartInterface {
	return newCharts(c, namespace)
}

func (c *RegistryV1Client) ChartGroups() ChartGroupInterface {
	return newChartGroups(c)
}

func (c *RegistryV1Client) ConfigMaps() ConfigMapInterface {
	return newConfigMaps(c)
}

func (c *RegistryV1Client) Namespaces() NamespaceInterface {
	return newNamespaces(c)
}

func (c *RegistryV1Client) Repositories(namespace string) RepositoryInterface {
	return newRepositories(c, namespace)
}

// NewForConfig creates a new RegistryV1Client for the given config.
func NewForConfig(c *rest.Config) (*RegistryV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &RegistryV1Client{client}, nil
}

// NewForConfigOrDie creates a new RegistryV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *RegistryV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new RegistryV1Client for the given RESTClient.
func New(c rest.Interface) *RegistryV1Client {
	return &RegistryV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *RegistryV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
