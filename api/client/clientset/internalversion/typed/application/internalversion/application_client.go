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

package internalversion

import (
	rest "k8s.io/client-go/rest"
	"tkestack.io/tke/api/client/clientset/internalversion/scheme"
)

type ApplicationInterface interface {
	RESTClient() rest.Interface
	AppsGetter
	AppHistoriesGetter
	AppResourcesGetter
	ConfigMapsGetter
}

// ApplicationClient is used to interact with features provided by the application.tkestack.io group.
type ApplicationClient struct {
	restClient rest.Interface
}

func (c *ApplicationClient) Apps(namespace string) AppInterface {
	return newApps(c, namespace)
}

func (c *ApplicationClient) AppHistories(namespace string) AppHistoryInterface {
	return newAppHistories(c, namespace)
}

func (c *ApplicationClient) AppResources(namespace string) AppResourceInterface {
	return newAppResources(c, namespace)
}

func (c *ApplicationClient) ConfigMaps() ConfigMapInterface {
	return newConfigMaps(c)
}

// NewForConfig creates a new ApplicationClient for the given config.
func NewForConfig(c *rest.Config) (*ApplicationClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ApplicationClient{client}, nil
}

// NewForConfigOrDie creates a new ApplicationClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ApplicationClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ApplicationClient for the given RESTClient.
func New(c rest.Interface) *ApplicationClient {
	return &ApplicationClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	config.APIPath = "/apis"
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("application.tkestack.io")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("application.tkestack.io")[0]
		config.GroupVersion = &gv
	}
	config.NegotiatedSerializer = scheme.Codecs

	if config.QPS == 0 {
		config.QPS = 5
	}
	if config.Burst == 0 {
		config.Burst = 10
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ApplicationClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
