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

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1 "tkestack.io/tke/api/client/clientset/versioned/typed/notify/v1"
)

type FakeNotifyV1 struct {
	*testing.Fake
}

func (c *FakeNotifyV1) Channels() v1.ChannelInterface {
	return &FakeChannels{c}
}

func (c *FakeNotifyV1) ConfigMaps() v1.ConfigMapInterface {
	return &FakeConfigMaps{c}
}

func (c *FakeNotifyV1) Messages() v1.MessageInterface {
	return &FakeMessages{c}
}

func (c *FakeNotifyV1) MessageRequests(namespace string) v1.MessageRequestInterface {
	return &FakeMessageRequests{c, namespace}
}

func (c *FakeNotifyV1) Receivers() v1.ReceiverInterface {
	return &FakeReceivers{c}
}

func (c *FakeNotifyV1) ReceiverGroups() v1.ReceiverGroupInterface {
	return &FakeReceiverGroups{c}
}

func (c *FakeNotifyV1) Templates(namespace string) v1.TemplateInterface {
	return &FakeTemplates{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeNotifyV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
