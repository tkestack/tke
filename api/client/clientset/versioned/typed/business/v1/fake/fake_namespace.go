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
	businessv1 "tkestack.io/tke/api/business/v1"
)

// FakeNamespaces implements NamespaceInterface
type FakeNamespaces struct {
	Fake *FakeBusinessV1
	ns   string
}

var namespacesResource = schema.GroupVersionResource{Group: "business.tkestack.io", Version: "v1", Resource: "namespaces"}

var namespacesKind = schema.GroupVersionKind{Group: "business.tkestack.io", Version: "v1", Kind: "Namespace"}

// Get takes name of the namespace, and returns the corresponding namespace object, and an error if there is any.
func (c *FakeNamespaces) Get(ctx context.Context, name string, options v1.GetOptions) (result *businessv1.Namespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(namespacesResource, c.ns, name), &businessv1.Namespace{})

	if obj == nil {
		return nil, err
	}
	return obj.(*businessv1.Namespace), err
}

// List takes label and field selectors, and returns the list of Namespaces that match those selectors.
func (c *FakeNamespaces) List(ctx context.Context, opts v1.ListOptions) (result *businessv1.NamespaceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(namespacesResource, namespacesKind, c.ns, opts), &businessv1.NamespaceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &businessv1.NamespaceList{ListMeta: obj.(*businessv1.NamespaceList).ListMeta}
	for _, item := range obj.(*businessv1.NamespaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested namespaces.
func (c *FakeNamespaces) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(namespacesResource, c.ns, opts))

}

// Create takes the representation of a namespace and creates it.  Returns the server's representation of the namespace, and an error, if there is any.
func (c *FakeNamespaces) Create(ctx context.Context, namespace *businessv1.Namespace, opts v1.CreateOptions) (result *businessv1.Namespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(namespacesResource, c.ns, namespace), &businessv1.Namespace{})

	if obj == nil {
		return nil, err
	}
	return obj.(*businessv1.Namespace), err
}

// Update takes the representation of a namespace and updates it. Returns the server's representation of the namespace, and an error, if there is any.
func (c *FakeNamespaces) Update(ctx context.Context, namespace *businessv1.Namespace, opts v1.UpdateOptions) (result *businessv1.Namespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(namespacesResource, c.ns, namespace), &businessv1.Namespace{})

	if obj == nil {
		return nil, err
	}
	return obj.(*businessv1.Namespace), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNamespaces) UpdateStatus(ctx context.Context, namespace *businessv1.Namespace, opts v1.UpdateOptions) (*businessv1.Namespace, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(namespacesResource, "status", c.ns, namespace), &businessv1.Namespace{})

	if obj == nil {
		return nil, err
	}
	return obj.(*businessv1.Namespace), err
}

// Delete takes name of the namespace and deletes it. Returns an error if one occurs.
func (c *FakeNamespaces) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(namespacesResource, c.ns, name), &businessv1.Namespace{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNamespaces) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(namespacesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &businessv1.NamespaceList{})
	return err
}

// Patch applies the patch and returns the patched namespace.
func (c *FakeNamespaces) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *businessv1.Namespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(namespacesResource, c.ns, name, pt, data, subresources...), &businessv1.Namespace{})

	if obj == nil {
		return nil, err
	}
	return obj.(*businessv1.Namespace), err
}
