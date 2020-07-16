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
	platformv1 "tkestack.io/tke/api/platform/v1"
)

// FakeTemplates implements TemplateInterface
type FakeTemplates struct {
	Fake *FakePlatformV1
}

var templatesResource = schema.GroupVersionResource{Group: "platform.tkestack.io", Version: "v1", Resource: "templates"}

var templatesKind = schema.GroupVersionKind{Group: "platform.tkestack.io", Version: "v1", Kind: "Template"}

// Get takes name of the template, and returns the corresponding template object, and an error if there is any.
func (c *FakeTemplates) Get(ctx context.Context, name string, options v1.GetOptions) (result *platformv1.Template, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(templatesResource, name), &platformv1.Template{})
	if obj == nil {
		return nil, err
	}
	return obj.(*platformv1.Template), err
}

// List takes label and field selectors, and returns the list of Templates that match those selectors.
func (c *FakeTemplates) List(ctx context.Context, opts v1.ListOptions) (result *platformv1.TemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(templatesResource, templatesKind, opts), &platformv1.TemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &platformv1.TemplateList{ListMeta: obj.(*platformv1.TemplateList).ListMeta}
	for _, item := range obj.(*platformv1.TemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested templates.
func (c *FakeTemplates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(templatesResource, opts))
}

// Create takes the representation of a template and creates it.  Returns the server's representation of the template, and an error, if there is any.
func (c *FakeTemplates) Create(ctx context.Context, template *platformv1.Template, opts v1.CreateOptions) (result *platformv1.Template, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(templatesResource, template), &platformv1.Template{})
	if obj == nil {
		return nil, err
	}
	return obj.(*platformv1.Template), err
}

// Update takes the representation of a template and updates it. Returns the server's representation of the template, and an error, if there is any.
func (c *FakeTemplates) Update(ctx context.Context, template *platformv1.Template, opts v1.UpdateOptions) (result *platformv1.Template, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(templatesResource, template), &platformv1.Template{})
	if obj == nil {
		return nil, err
	}
	return obj.(*platformv1.Template), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeTemplates) UpdateStatus(ctx context.Context, template *platformv1.Template, opts v1.UpdateOptions) (*platformv1.Template, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(templatesResource, "status", template), &platformv1.Template{})
	if obj == nil {
		return nil, err
	}
	return obj.(*platformv1.Template), err
}

// Delete takes name of the template and deletes it. Returns an error if one occurs.
func (c *FakeTemplates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(templatesResource, name), &platformv1.Template{})
	return err
}

// Patch applies the patch and returns the patched template.
func (c *FakeTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *platformv1.Template, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(templatesResource, name, pt, data, subresources...), &platformv1.Template{})
	if obj == nil {
		return nil, err
	}
	return obj.(*platformv1.Template), err
}
