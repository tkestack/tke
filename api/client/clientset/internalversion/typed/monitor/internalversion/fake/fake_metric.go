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
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
	monitor "tkestack.io/tke/api/monitor"
)

// FakeMetrics implements MetricInterface
type FakeMetrics struct {
	Fake *FakeMonitor
}

var metricsResource = schema.GroupVersionResource{Group: "monitor.tkestack.io", Version: "", Resource: "metrics"}

var metricsKind = schema.GroupVersionKind{Group: "monitor.tkestack.io", Version: "", Kind: "Metric"}

// Create takes the representation of a metric and creates it.  Returns the server's representation of the metric, and an error, if there is any.
func (c *FakeMetrics) Create(ctx context.Context, metric *monitor.Metric, opts v1.CreateOptions) (result *monitor.Metric, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(metricsResource, metric), &monitor.Metric{})
	if obj == nil {
		return nil, err
	}
	return obj.(*monitor.Metric), err
}
