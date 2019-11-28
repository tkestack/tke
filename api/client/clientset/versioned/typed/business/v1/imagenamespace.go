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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "tkestack.io/tke/api/business/v1"
	scheme "tkestack.io/tke/api/client/clientset/versioned/scheme"
)

// ImageNamespacesGetter has a method to return a ImageNamespaceInterface.
// A group's client should implement this interface.
type ImageNamespacesGetter interface {
	ImageNamespaces(namespace string) ImageNamespaceInterface
}

// ImageNamespaceInterface has methods to work with ImageNamespace resources.
type ImageNamespaceInterface interface {
	Create(*v1.ImageNamespace) (*v1.ImageNamespace, error)
	Update(*v1.ImageNamespace) (*v1.ImageNamespace, error)
	UpdateStatus(*v1.ImageNamespace) (*v1.ImageNamespace, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ImageNamespace, error)
	List(opts metav1.ListOptions) (*v1.ImageNamespaceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ImageNamespace, err error)
	ImageNamespaceExpansion
}

// imageNamespaces implements ImageNamespaceInterface
type imageNamespaces struct {
	client rest.Interface
	ns     string
}

// newImageNamespaces returns a ImageNamespaces
func newImageNamespaces(c *BusinessV1Client, namespace string) *imageNamespaces {
	return &imageNamespaces{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the imageNamespace, and returns the corresponding imageNamespace object, and an error if there is any.
func (c *imageNamespaces) Get(name string, options metav1.GetOptions) (result *v1.ImageNamespace, err error) {
	result = &v1.ImageNamespace{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagenamespaces").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ImageNamespaces that match those selectors.
func (c *imageNamespaces) List(opts metav1.ListOptions) (result *v1.ImageNamespaceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ImageNamespaceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagenamespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested imageNamespaces.
func (c *imageNamespaces) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("imagenamespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a imageNamespace and creates it.  Returns the server's representation of the imageNamespace, and an error, if there is any.
func (c *imageNamespaces) Create(imageNamespace *v1.ImageNamespace) (result *v1.ImageNamespace, err error) {
	result = &v1.ImageNamespace{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("imagenamespaces").
		Body(imageNamespace).
		Do().
		Into(result)
	return
}

// Update takes the representation of a imageNamespace and updates it. Returns the server's representation of the imageNamespace, and an error, if there is any.
func (c *imageNamespaces) Update(imageNamespace *v1.ImageNamespace) (result *v1.ImageNamespace, err error) {
	result = &v1.ImageNamespace{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("imagenamespaces").
		Name(imageNamespace.Name).
		Body(imageNamespace).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *imageNamespaces) UpdateStatus(imageNamespace *v1.ImageNamespace) (result *v1.ImageNamespace, err error) {
	result = &v1.ImageNamespace{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("imagenamespaces").
		Name(imageNamespace.Name).
		SubResource("status").
		Body(imageNamespace).
		Do().
		Into(result)
	return
}

// Delete takes name of the imageNamespace and deletes it. Returns an error if one occurs.
func (c *imageNamespaces) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("imagenamespaces").
		Name(name).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched imageNamespace.
func (c *imageNamespaces) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ImageNamespace, err error) {
	result = &v1.ImageNamespace{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("imagenamespaces").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
