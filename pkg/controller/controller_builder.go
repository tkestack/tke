/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package controller

import (
	"k8s.io/client-go/rest"
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/pkg/util/log"
)

// ClientBuilder allows you to get clients and configs for controllers.
type ClientBuilder interface {
	Config(name string) (*rest.Config, error)
	ConfigOrDie(name string) *rest.Config
	Client(name string) (versionedclientset.Interface, error)
	ClientOrDie(name string) versionedclientset.Interface
	ClientGoClient(name string) (versionedclientset.Interface, error)
	ClientGoClientOrDie(name string) versionedclientset.Interface
}

// SimpleControllerClientBuilder returns a fixed client with different user agents.
type SimpleControllerClientBuilder struct {
	// ClientConfig is a skeleton config to clone and use as the basis for each controller client.
	ClientConfig *rest.Config
}

// Config returns a complete clientConfig for constructing clients.
func (b SimpleControllerClientBuilder) Config(name string) (*rest.Config, error) {
	clientConfig := *b.ClientConfig
	return rest.AddUserAgent(&clientConfig, name), nil
}

// ConfigOrDie returns a complete clientConfig for constructing clients.
func (b SimpleControllerClientBuilder) ConfigOrDie(name string) *rest.Config {
	clientConfig, err := b.Config(name)
	if err != nil {
		log.Fatal("Unexpected fatal error", log.Err(err))
	}
	return clientConfig
}

// Client returns the complete client set for constructing clients.
func (b SimpleControllerClientBuilder) Client(name string) (versionedclientset.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return versionedclientset.NewForConfig(clientConfig)
}

// ClientOrDie returns the complete client set for constructing clients.
func (b SimpleControllerClientBuilder) ClientOrDie(name string) versionedclientset.Interface {
	client, err := b.Client(name)
	if err != nil {
		log.Fatal("Unexpected fatal error", log.Err(err))
	}
	return client
}

// ClientGoClient returns the complete client set for constructing clients.
func (b SimpleControllerClientBuilder) ClientGoClient(name string) (versionedclientset.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return versionedclientset.NewForConfig(clientConfig)
}

// ClientGoClientOrDie returns the complete client set for constructing clients.
func (b SimpleControllerClientBuilder) ClientGoClientOrDie(name string) versionedclientset.Interface {
	client, err := b.ClientGoClient(name)
	if err != nil {
		log.Fatal("Unexpected fatal error", log.Err(err))
	}
	return client
}
