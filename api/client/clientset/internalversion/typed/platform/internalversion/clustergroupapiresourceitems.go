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
	rest "k8s.io/client-go/rest"
	scheme "tkestack.io/tke/api/client/clientset/internalversion/scheme"
	platform "tkestack.io/tke/api/platform"
)

// ClusterGroupAPIResourceItemsesGetter has a method to return a ClusterGroupAPIResourceItemsInterface.
// A group's client should implement this interface.
type ClusterGroupAPIResourceItemsesGetter interface {
	ClusterGroupAPIResourceItemses() ClusterGroupAPIResourceItemsInterface
}

// ClusterGroupAPIResourceItemsInterface has methods to work with ClusterGroupAPIResourceItems resources.
type ClusterGroupAPIResourceItemsInterface interface {
	Get(ctx context.Context, name string, opts v1.GetOptions) (*platform.ClusterGroupAPIResourceItems, error)
	List(ctx context.Context, opts v1.ListOptions) (*platform.ClusterGroupAPIResourceItemsList, error)
	ClusterGroupAPIResourceItemsExpansion
}

// clusterGroupAPIResourceItemses implements ClusterGroupAPIResourceItemsInterface
type clusterGroupAPIResourceItemses struct {
	client rest.Interface
}

// newClusterGroupAPIResourceItemses returns a ClusterGroupAPIResourceItemses
func newClusterGroupAPIResourceItemses(c *PlatformClient) *clusterGroupAPIResourceItemses {
	return &clusterGroupAPIResourceItemses{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterGroupAPIResourceItems, and returns the corresponding clusterGroupAPIResourceItems object, and an error if there is any.
func (c *clusterGroupAPIResourceItemses) Get(ctx context.Context, name string, options v1.GetOptions) (result *platform.ClusterGroupAPIResourceItems, err error) {
	result = &platform.ClusterGroupAPIResourceItems{}
	err = c.client.Get().
		Resource("clustergroupapiresourceitemses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterGroupAPIResourceItemses that match those selectors.
func (c *clusterGroupAPIResourceItemses) List(ctx context.Context, opts v1.ListOptions) (result *platform.ClusterGroupAPIResourceItemsList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &platform.ClusterGroupAPIResourceItemsList{}
	err = c.client.Get().
		Resource("clustergroupapiresourceitemses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
