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
 */

package addon

// These clients are for addon, decoupling with provider

import (
	"context"
	"fmt"

	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
)

// BuildExternalMonitoringClientSetNoStatus creates the monitoring clientset of prometheus operator by given
// cluster object and returns it.
func BuildExternalMonitoringClientSetNoStatus(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}
	restConfig := credential.RESTConfig(cluster)
	return monitoringclient.NewForConfig(restConfig)
}

// BuildExternalMonitoringClientSet creates the monitoring clientset of  prometheus operator by given cluster
// object and returns it.
func BuildExternalMonitoringClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (monitoringclient.Interface, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Phase != platformv1.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	return BuildExternalMonitoringClientSetNoStatus(ctx, cluster, client)
}

// BuildExternalMonitoringClientSetWithName creates the clientset of prometheus operator by given cluster
// name and returns it.
func BuildExternalMonitoringClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (monitoringclient.Interface, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalMonitoringClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildKubeAggregatorClientSet creates the kube-aggregator clientset of kubernetes by given cluster
// object and returns it.
func BuildKubeAggregatorClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubeaggregatorclientset.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Phase != platformv1.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	return BuildKubeAggregatorClientSetNoStatus(ctx, cluster, client)
}

// BuildExternalExtensionClientSetNoStatus creates the api extension clientset of kubernetes by given
// cluster object and returns it.
func BuildExternalExtensionClientSetNoStatus(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}
	restConfig := credential.RESTConfig(cluster)
	return apiextensionsclient.NewForConfig(restConfig)
}

// BuildExternalExtensionClientSet creates the api extension clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalExtensionClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*apiextensionsclient.Clientset, error) {
	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Phase != platformv1.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	return BuildExternalExtensionClientSetNoStatus(ctx, cluster, client)
}

// BuildKubeAggregatorClientSetNoStatus creates the kube-aggregator clientset of kubernetes by given
// cluster object and returns it.
func BuildKubeAggregatorClientSetNoStatus(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubeaggregatorclientset.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}
	restConfig := credential.RESTConfig(cluster)
	return kubeaggregatorclientset.NewForConfig(restConfig)
}

// BuildExternalClientSetWithName creates the clientset of kubernetes by given cluster
// name and returns it.
func BuildExternalClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (*kubernetes.Clientset, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalClientSet creates the clientset of kubernetes by given cluster object and returns it.
func BuildExternalClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubernetes.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Phase != platformv1.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	return util.BuildVersionedClientSet(cluster, credential)
}

// GetClusterCredentialV1 returns the versioned cluster's credential
func GetClusterCredentialV1(ctx context.Context, client platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster) (*platformv1.ClusterCredential, error) {
	var (
		credential *platformv1.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	} else {
		return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), cluster.Name)
	}

	return credential, nil
}
