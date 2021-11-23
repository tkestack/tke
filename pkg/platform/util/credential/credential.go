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

package credential

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

// GetClusterCredential returns the cluster's credential
func GetClusterCredential(ctx context.Context, client platforminternalclient.PlatformInterface, cluster *platform.Cluster, username string) (*platform.ClusterCredential, error) {
	var (
		credential *platform.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
	} else if client != nil {
		clusterName := cluster.Name
		fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
		clusterCredentials, err := client.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, err
		}
		if clusterCredentials == nil || clusterCredentials.Items == nil || len(clusterCredentials.Items) == 0 {
			return nil, apierrors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
		}
		credential = &clusterCredentials.Items[0]
	}

	return credential, nil
}

// GetClusterCredentialV1 returns the versioned cluster's credential
func GetClusterCredentialV1(ctx context.Context, client platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster, username string) (*platformv1.ClusterCredential, error) {
	var (
		credential *platformv1.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, err
		}
	} else if client != nil {
		clusterName := cluster.Name
		fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
		clusterCredentials, err := client.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return nil, err
		}
		if clusterCredentials == nil || clusterCredentials.Items == nil || len(clusterCredentials.Items) == 0 {
			return nil, apierrors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
		}
		credential = &clusterCredentials.Items[0]
	}

	return credential, nil
}
