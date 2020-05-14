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
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	scheme "tkestack.io/tke/api/client/clientset/versioned/scheme"
	v1 "tkestack.io/tke/api/platform/v1"
)

// IPAMsGetter has a method to return a IPAMInterface.
// A group's client should implement this interface.
type IPAMsGetter interface {
	IPAMs() IPAMInterface
}

// IPAMInterface has methods to work with IPAM resources.
type IPAMInterface interface {
	Create(ctx context.Context, iPAM *v1.IPAM, opts metav1.CreateOptions) (*v1.IPAM, error)
	Update(ctx context.Context, iPAM *v1.IPAM, opts metav1.UpdateOptions) (*v1.IPAM, error)
	UpdateStatus(ctx context.Context, iPAM *v1.IPAM, opts metav1.UpdateOptions) (*v1.IPAM, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.IPAM, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.IPAMList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.IPAM, err error)
	IPAMExpansion
}

// iPAMs implements IPAMInterface
type iPAMs struct {
	client rest.Interface
}

// newIPAMs returns a IPAMs
func newIPAMs(c *PlatformV1Client) *iPAMs {
	return &iPAMs{
		client: c.RESTClient(),
	}
}

// Get takes name of the iPAM, and returns the corresponding iPAM object, and an error if there is any.
func (c *iPAMs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.IPAM, err error) {
	result = &v1.IPAM{}
	err = c.client.Get().
		Resource("ipams").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IPAMs that match those selectors.
func (c *iPAMs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.IPAMList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.IPAMList{}
	err = c.client.Get().
		Resource("ipams").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested iPAMs.
func (c *iPAMs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("ipams").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a iPAM and creates it.  Returns the server's representation of the iPAM, and an error, if there is any.
func (c *iPAMs) Create(ctx context.Context, iPAM *v1.IPAM, opts metav1.CreateOptions) (result *v1.IPAM, err error) {
	result = &v1.IPAM{}
	err = c.client.Post().
		Resource("ipams").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPAM).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a iPAM and updates it. Returns the server's representation of the iPAM, and an error, if there is any.
func (c *iPAMs) Update(ctx context.Context, iPAM *v1.IPAM, opts metav1.UpdateOptions) (result *v1.IPAM, err error) {
	result = &v1.IPAM{}
	err = c.client.Put().
		Resource("ipams").
		Name(iPAM.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPAM).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *iPAMs) UpdateStatus(ctx context.Context, iPAM *v1.IPAM, opts metav1.UpdateOptions) (result *v1.IPAM, err error) {
	result = &v1.IPAM{}
	err = c.client.Put().
		Resource("ipams").
		Name(iPAM.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(iPAM).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the iPAM and deletes it. Returns an error if one occurs.
func (c *iPAMs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("ipams").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched iPAM.
func (c *iPAMs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.IPAM, err error) {
	result = &v1.IPAM{}
	err = c.client.Patch(pt).
		Resource("ipams").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
