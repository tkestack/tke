/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package cluster

import (
	"context"
	"fmt"
	"sort"
	"sync"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/types"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/platform/util/credential"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]Provider)
)

const AdminUsername = "admin"

// Register makes a provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("cluster: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("cluster: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// re register provider, if provider's name exists, new provider will replace old provider
func ReRegister(name string, provider Provider) error {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		return fmt.Errorf("cluster: Register provider is nil")
	}
	providers[name] = provider
	return nil
}

// RegisterHandler register all provider's hanlder.
func RegisterHandler(mux *mux.PathRecorderMux) {
	for _, p := range providers {
		p.RegisterHandler(mux)
	}
}

// Setup call all provider's setup method.
func Setup() error {
	for _, p := range providers {
		if err := p.Setup(); err != nil {
			return fmt.Errorf("%s.Setup error:%w", p.Name(), err)
		}
	}

	return nil
}

// Teardown call all provider's teardown method.
func Teardown() error {
	for _, p := range providers {
		if err := p.Teardown(); err != nil {
			return fmt.Errorf("%s.Teardown error:%w", p.Name(), err)
		}
	}

	return nil
}

// Providers returns a sorted list of the names of the registered providers.
func Providers() []string {
	providersMu.RLock()
	defer providersMu.RUnlock()
	var list []string
	for name := range providers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetProvider returns provider by name
func GetProvider(name string) (Provider, error) {
	providersMu.RLock()
	provider, ok := providers[name]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cluster: unknown provider %q (forgotten import?)", name)

	}

	return provider, nil
}

func GetCluster(ctx context.Context, platformClient internalversion.PlatformInterface, cluster *platform.Cluster, username string) (*types.Cluster, error) {
	result := new(types.Cluster)
	result.Cluster = cluster
	provider, err := GetProvider(cluster.Spec.Type)
	if err != nil {
		return nil, err
	}
	clusterCredential, err := credential.GetClusterCredential(ctx, platformClient, cluster, username)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, err
	}
	clusterv1 := &platformv1.Cluster{}
	err = platformv1.Convert_platform_Cluster_To_v1_Cluster(cluster, clusterv1, nil)
	if err != nil {
		return nil, err
	}
	restConfig, err := provider.GetRestConfig(ctx, clusterv1, username)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, err
	}
	result.ClusterCredential = clusterCredential
	result.RegisterRestConfig(restConfig)

	return result, nil
}

func GetClusterByName(ctx context.Context, platformClient internalversion.PlatformInterface, clsname, username string) (*types.Cluster, error) {
	cluster, err := platformClient.Clusters().Get(ctx, clsname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return GetCluster(ctx, platformClient, cluster, username)
}

func GetV1Cluster(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster, username string) (*v1.Cluster, error) {
	result := new(v1.Cluster)
	result.Cluster = cluster
	result.IsCredentialChanged = false
	provider, err := GetProvider(cluster.Spec.Type)
	if err != nil {
		return nil, err
	}
	clusterCredential, err := credential.GetClusterCredentialV1(ctx, platformClient, cluster, username)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, err
	}
	restConfig, err := provider.GetRestConfig(ctx, cluster, username)
	if err != nil && !apierrors.IsNotFound(err) {
		return result, err
	}
	result.ClusterCredential = clusterCredential
	result.RegisterRestConfig(restConfig)

	return result, nil
}

func GetV1ClusterByName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, clsname, username string) (*v1.Cluster, error) {
	cluster, err := platformClient.Clusters().Get(ctx, clsname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return GetV1Cluster(ctx, platformClient, cluster, username)
}
