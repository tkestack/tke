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

// ClusterCredentialsGetter has a method to return a ClusterCredentialInterface.
// A group's client should implement this interface.
type ClusterCredentialsGetter interface {
	ClusterCredentials() ClusterCredentialInterface
}

// ClusterCredentialInterface has methods to work with ClusterCredential resources.
type ClusterCredentialInterface interface {
	Create(ctx context.Context, clusterCredential *v1.ClusterCredential, opts metav1.CreateOptions) (*v1.ClusterCredential, error)
	Update(ctx context.Context, clusterCredential *v1.ClusterCredential, opts metav1.UpdateOptions) (*v1.ClusterCredential, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ClusterCredential, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ClusterCredentialList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterCredential, err error)
	ClusterCredentialExpansion
}

// clusterCredentials implements ClusterCredentialInterface
type clusterCredentials struct {
	client rest.Interface
}

// newClusterCredentials returns a ClusterCredentials
func newClusterCredentials(c *PlatformV1Client) *clusterCredentials {
	return &clusterCredentials{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterCredential, and returns the corresponding clusterCredential object, and an error if there is any.
func (c *clusterCredentials) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.ClusterCredential, err error) {
	result = &v1.ClusterCredential{}
	err = c.client.Get().
		Resource("clustercredentials").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterCredentials that match those selectors.
func (c *clusterCredentials) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ClusterCredentialList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ClusterCredentialList{}
	err = c.client.Get().
		Resource("clustercredentials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterCredentials.
func (c *clusterCredentials) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clustercredentials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterCredential and creates it.  Returns the server's representation of the clusterCredential, and an error, if there is any.
func (c *clusterCredentials) Create(ctx context.Context, clusterCredential *v1.ClusterCredential, opts metav1.CreateOptions) (result *v1.ClusterCredential, err error) {
	result = &v1.ClusterCredential{}
	err = c.client.Post().
		Resource("clustercredentials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterCredential).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterCredential and updates it. Returns the server's representation of the clusterCredential, and an error, if there is any.
func (c *clusterCredentials) Update(ctx context.Context, clusterCredential *v1.ClusterCredential, opts metav1.UpdateOptions) (result *v1.ClusterCredential, err error) {
	result = &v1.ClusterCredential{}
	err = c.client.Put().
		Resource("clustercredentials").
		Name(clusterCredential.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterCredential).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterCredential and deletes it. Returns an error if one occurs.
func (c *clusterCredentials) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clustercredentials").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterCredential.
func (c *clusterCredentials) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterCredential, err error) {
	result = &v1.ClusterCredential{}
	err = c.client.Patch(pt).
		Resource("clustercredentials").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
