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
	rest "k8s.io/client-go/rest"
)

// AppResourcesGetter has a method to return a AppResourceInterface.
// A group's client should implement this interface.
type AppResourcesGetter interface {
	AppResources(namespace string) AppResourceInterface
}

// AppResourceInterface has methods to work with AppResource resources.
type AppResourceInterface interface {
	AppResourceExpansion
}

// appResources implements AppResourceInterface
type appResources struct {
	client rest.Interface
	ns     string
}

// newAppResources returns a AppResources
func newAppResources(c *ApplicationV1Client, namespace string) *appResources {
	return &appResources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}
