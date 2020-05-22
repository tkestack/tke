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
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	scheme "tkestack.io/tke/api/client/clientset/internalversion/scheme"
	platform "tkestack.io/tke/api/platform"
)

// PersistentEventsGetter has a method to return a PersistentEventInterface.
// A group's client should implement this interface.
type PersistentEventsGetter interface {
	PersistentEvents() PersistentEventInterface
}

// PersistentEventInterface has methods to work with PersistentEvent resources.
type PersistentEventInterface interface {
	Create(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.CreateOptions) (*platform.PersistentEvent, error)
	Update(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.UpdateOptions) (*platform.PersistentEvent, error)
	UpdateStatus(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.UpdateOptions) (*platform.PersistentEvent, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*platform.PersistentEvent, error)
	List(ctx context.Context, opts v1.ListOptions) (*platform.PersistentEventList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *platform.PersistentEvent, err error)
	PersistentEventExpansion
}

// persistentEvents implements PersistentEventInterface
type persistentEvents struct {
	client rest.Interface
}

// newPersistentEvents returns a PersistentEvents
func newPersistentEvents(c *PlatformClient) *persistentEvents {
	return &persistentEvents{
		client: c.RESTClient(),
	}
}

// Get takes name of the persistentEvent, and returns the corresponding persistentEvent object, and an error if there is any.
func (c *persistentEvents) Get(ctx context.Context, name string, options v1.GetOptions) (result *platform.PersistentEvent, err error) {
	result = &platform.PersistentEvent{}
	err = c.client.Get().
		Resource("persistentevents").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PersistentEvents that match those selectors.
func (c *persistentEvents) List(ctx context.Context, opts v1.ListOptions) (result *platform.PersistentEventList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &platform.PersistentEventList{}
	err = c.client.Get().
		Resource("persistentevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested persistentEvents.
func (c *persistentEvents) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("persistentevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a persistentEvent and creates it.  Returns the server's representation of the persistentEvent, and an error, if there is any.
func (c *persistentEvents) Create(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.CreateOptions) (result *platform.PersistentEvent, err error) {
	result = &platform.PersistentEvent{}
	err = c.client.Post().
		Resource("persistentevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(persistentEvent).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a persistentEvent and updates it. Returns the server's representation of the persistentEvent, and an error, if there is any.
func (c *persistentEvents) Update(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.UpdateOptions) (result *platform.PersistentEvent, err error) {
	result = &platform.PersistentEvent{}
	err = c.client.Put().
		Resource("persistentevents").
		Name(persistentEvent.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(persistentEvent).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *persistentEvents) UpdateStatus(ctx context.Context, persistentEvent *platform.PersistentEvent, opts v1.UpdateOptions) (result *platform.PersistentEvent, err error) {
	result = &platform.PersistentEvent{}
	err = c.client.Put().
		Resource("persistentevents").
		Name(persistentEvent.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(persistentEvent).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the persistentEvent and deletes it. Returns an error if one occurs.
func (c *persistentEvents) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("persistentevents").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched persistentEvent.
func (c *persistentEvents) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *platform.PersistentEvent, err error) {
	result = &platform.PersistentEvent{}
	err = c.client.Patch(pt).
		Resource("persistentevents").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
