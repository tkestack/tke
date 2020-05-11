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

package validation

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

type ClusterGetter interface {
	Cluster(ctx context.Context, name string, options metav1.GetOptions) (*platformv1.Cluster, error)
}

type BusinessObjectGetter interface {
	Project(ctx context.Context, name string, options metav1.GetOptions) (*business.Project, error)
	Namespace(ctx context.Context, project, name string, options metav1.GetOptions) (*business.Namespace, error)
}

func NewClusterGetter(platformClient platformversionedclient.PlatformV1Interface) ClusterGetter {
	return &clusterGetter{platformClient: platformClient}
}

func NewObjectGetter(businessClient *businessinternalclient.BusinessClient) BusinessObjectGetter {
	return &businessObjectGetter{businessClient: businessClient}
}

type clusterGetter struct {
	platformClient platformversionedclient.PlatformV1Interface
}

func (getter *clusterGetter) Cluster(ctx context.Context, name string, options metav1.GetOptions) (*platformv1.Cluster, error) {
	return getter.platformClient.Clusters().Get(ctx, name, options)
}

type businessObjectGetter struct {
	businessClient *businessinternalclient.BusinessClient
}

func (getter *businessObjectGetter) Project(ctx context.Context, name string, options metav1.GetOptions) (*business.Project, error) {
	return getter.businessClient.Projects().Get(ctx, name, options)
}

func (getter *businessObjectGetter) Namespace(ctx context.Context, project, name string, options metav1.GetOptions) (*business.Namespace, error) {
	return getter.businessClient.Namespaces(project).Get(ctx, name, options)
}

// ValidateCluster validate cluster
func ValidateCluster(ctx context.Context, platformClient platforminternalclient.PlatformInterface, clusterName string) field.ErrorList {
	var allErrs field.ErrorList
	if clusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify cluster name"))
	} else {
		_, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), clusterName, fmt.Sprintf("can't get cluster:%s", err)))
		}
	}
	return allErrs
}

// ValidateClusterVersioned validate cluster
func ValidateClusterVersioned(ctx context.Context, getter ClusterGetter, clusterName string, tenantID string) field.ErrorList {
	var allErrs field.ErrorList
	if clusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify cluster name"))
	} else {
		fldPath := field.NewPath("spec", "clusterName")
		cluster, err := getter.Cluster(ctx, clusterName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			allErrs = append(allErrs, field.NotFound(fldPath, clusterName))
		} else if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		} else {
			if tenantID != "" && tenantID != cluster.Spec.TenantID {
				allErrs = append(allErrs, field.NotFound(fldPath, clusterName))
			}
		}
	}
	return allErrs
}
