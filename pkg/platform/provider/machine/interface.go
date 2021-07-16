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

package machine

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
)

type APIProvider interface {
	Validate(ctx context.Context, cluster *platform.Machine) field.ErrorList
	ValidateUpdate(ctx context.Context, cluster *platform.Machine, oldCluster *platform.Machine) field.ErrorList
	PreCreate(ctx context.Context, cluster *platform.Machine) error
	AfterCreate(cluster *platform.Machine) error
}

type ControllerProvider interface {
	OnCreate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnUpdate(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
	OnDelete(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error
}

// Provider defines a set of response interfaces for specific machine
// types in machine management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
}
