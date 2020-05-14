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
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	notify "tkestack.io/tke/api/notify"
)

// FakeChannels implements ChannelInterface
type FakeChannels struct {
	Fake *FakeNotify
}

var channelsResource = schema.GroupVersionResource{Group: "notify.tkestack.io", Version: "", Resource: "channels"}

var channelsKind = schema.GroupVersionKind{Group: "notify.tkestack.io", Version: "", Kind: "Channel"}

// Get takes name of the channel, and returns the corresponding channel object, and an error if there is any.
func (c *FakeChannels) Get(ctx context.Context, name string, options v1.GetOptions) (result *notify.Channel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(channelsResource, name), &notify.Channel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*notify.Channel), err
}

// List takes label and field selectors, and returns the list of Channels that match those selectors.
func (c *FakeChannels) List(ctx context.Context, opts v1.ListOptions) (result *notify.ChannelList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(channelsResource, channelsKind, opts), &notify.ChannelList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &notify.ChannelList{ListMeta: obj.(*notify.ChannelList).ListMeta}
	for _, item := range obj.(*notify.ChannelList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested channels.
func (c *FakeChannels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(channelsResource, opts))
}

// Create takes the representation of a channel and creates it.  Returns the server's representation of the channel, and an error, if there is any.
func (c *FakeChannels) Create(ctx context.Context, channel *notify.Channel, opts v1.CreateOptions) (result *notify.Channel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(channelsResource, channel), &notify.Channel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*notify.Channel), err
}

// Update takes the representation of a channel and updates it. Returns the server's representation of the channel, and an error, if there is any.
func (c *FakeChannels) Update(ctx context.Context, channel *notify.Channel, opts v1.UpdateOptions) (result *notify.Channel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(channelsResource, channel), &notify.Channel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*notify.Channel), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeChannels) UpdateStatus(ctx context.Context, channel *notify.Channel, opts v1.UpdateOptions) (*notify.Channel, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(channelsResource, "status", channel), &notify.Channel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*notify.Channel), err
}

// Delete takes name of the channel and deletes it. Returns an error if one occurs.
func (c *FakeChannels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(channelsResource, name), &notify.Channel{})
	return err
}

// Patch applies the patch and returns the patched channel.
func (c *FakeChannels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *notify.Channel, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(channelsResource, name, pt, data, subresources...), &notify.Channel{})
	if obj == nil {
		return nil, err
	}
	return obj.(*notify.Channel), err
}
