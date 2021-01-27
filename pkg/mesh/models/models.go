/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package models

import (
	"time"

	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

// A Namespace provide a scope for names
// This type is used to describe a set of objects.
//
// swagger:model namespace
type Namespace struct {
	// The id of the namespace.
	//
	// example:  istio-system
	// required: true
	Name string `json:"name"`

	// Creation date of the namespace.
	// There is no need to export this through the API. So, this is
	// set to be ignored by JSON package.
	//
	// required: true
	CreationTimestamp time.Time `json:"-"`
}

type BaseMetadata struct {
	Namespace Namespace `json:"namespace"`
	App       string    `json:"appruntime" validate:"required"`
}

type IstioNetworkingConfig struct {
	BaseMetadata

	VirtualService  *istionetworking.VirtualService  `json:"virtualService,omitempty"`
	Gateway         *istionetworking.Gateway         `json:"gateway,omitempty"`
	DestinationRule *istionetworking.DestinationRule `json:"destinationRule,omitempty"`
	ServiceEntry    *istionetworking.ServiceEntry    `json:"serviceEntry,omitempty"`
}
