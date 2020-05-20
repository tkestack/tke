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
	v1 "tkestack.io/tke/api/registry/v1"
)

// ChartGroupsGetter has a method to return a ChartGroupInterface.
// A group's client should implement this interface.
type ChartGroupsGetter interface {
	ChartGroups() ChartGroupInterface
}

// ChartGroupInterface has methods to work with ChartGroup resources.
type ChartGroupInterface interface {
	Create(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.CreateOptions) (*v1.ChartGroup, error)
	Update(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.UpdateOptions) (*v1.ChartGroup, error)
	UpdateStatus(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.UpdateOptions) (*v1.ChartGroup, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ChartGroup, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ChartGroupList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ChartGroup, err error)
	ChartGroupExpansion
}

// chartGroups implements ChartGroupInterface
type chartGroups struct {
	client rest.Interface
}

// newChartGroups returns a ChartGroups
func newChartGroups(c *RegistryV1Client) *chartGroups {
	return &chartGroups{
		client: c.RESTClient(),
	}
}

// Get takes name of the chartGroup, and returns the corresponding chartGroup object, and an error if there is any.
func (c *chartGroups) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ChartGroup, err error) {
	result = &v1.ChartGroup{}
	err = c.client.Get().
		Resource("chartgroups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ChartGroups that match those selectors.
func (c *chartGroups) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ChartGroupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ChartGroupList{}
	err = c.client.Get().
		Resource("chartgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested chartGroups.
func (c *chartGroups) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("chartgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a chartGroup and creates it.  Returns the server's representation of the chartGroup, and an error, if there is any.
func (c *chartGroups) Create(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.CreateOptions) (result *v1.ChartGroup, err error) {
	result = &v1.ChartGroup{}
	err = c.client.Post().
		Resource("chartgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(chartGroup).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a chartGroup and updates it. Returns the server's representation of the chartGroup, and an error, if there is any.
func (c *chartGroups) Update(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.UpdateOptions) (result *v1.ChartGroup, err error) {
	result = &v1.ChartGroup{}
	err = c.client.Put().
		Resource("chartgroups").
		Name(chartGroup.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(chartGroup).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *chartGroups) UpdateStatus(ctx context.Context, chartGroup *v1.ChartGroup, opts metav1.UpdateOptions) (result *v1.ChartGroup, err error) {
	result = &v1.ChartGroup{}
	err = c.client.Put().
		Resource("chartgroups").
		Name(chartGroup.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(chartGroup).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the chartGroup and deletes it. Returns an error if one occurs.
func (c *chartGroups) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("chartgroups").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched chartGroup.
func (c *chartGroups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ChartGroup, err error) {
	result = &v1.ChartGroup{}
	err = c.client.Patch(pt).
		Resource("chartgroups").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
