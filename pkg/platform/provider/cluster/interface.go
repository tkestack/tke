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

package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/server/mux"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/types"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
)

type APIProvider interface {
	RegisterHandler(mux *mux.PathRecorderMux)
	Validate(ctx context.Context, cluster *types.Cluster) field.ErrorList
	ValidateUpdate(ctx context.Context, cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList
	PreCreate(ctx context.Context, cluster *types.Cluster) error
	AfterCreate(cluster *types.Cluster) error
}

type ControllerProvider interface {
	// Setup called by controller to give an chance for plugin do some init work.
	Setup() error
	// Teardown called by controller for plugin do some clean job.
	Teardown() error

	OnCreate(ctx context.Context, cluster *v1.Cluster) error
	OnUpdate(ctx context.Context, cluster *v1.Cluster) error
	OnDelete(ctx context.Context, cluster *v1.Cluster) error
	// OnFilter called by cluster controller informer for plugin
	// do the filter on the cluster obj for specific case:
	// return bool:
	//  false: drop the object to the queue
	//  true: add the object to queue, AddFunc and UpdateFunc will
	//  go through later
	OnFilter(ctx context.Context, cluster *platformv1.Cluster) bool
	// OnRunning call on first running.
	OnRunning(ctx context.Context, cluster *v1.Cluster) error
}

type CredentialProvider interface {
	GetClusterCredential(ctx context.Context, client platforminternalclient.PlatformInterface, cluster *platform.Cluster, username string) (*platform.ClusterCredential, error)
	GetClusterCredentialV1(ctx context.Context, client platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster, username string) (*platformv1.ClusterCredential, error)
}

// Provider defines a set of response interfaces for specific cluster
// types in cluster management.
type Provider interface {
	Name() string

	APIProvider
	ControllerProvider
	CredentialProvider
}
