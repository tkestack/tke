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
	Cluster(name string, options metav1.GetOptions) (*platformv1.Cluster, error)
}

type BusinessObjectGetter interface {
	Project(name string, options metav1.GetOptions) (*business.Project, error)
	Namespace(project, name string, options metav1.GetOptions) (*business.Namespace, error)
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

func (getter *clusterGetter) Cluster(name string, options metav1.GetOptions) (*platformv1.Cluster, error) {
	return getter.platformClient.Clusters().Get(name, options)
}

type businessObjectGetter struct {
	businessClient *businessinternalclient.BusinessClient
}

func (getter *businessObjectGetter) Project(name string, options metav1.GetOptions) (*business.Project, error) {
	return getter.businessClient.Projects().Get(name, options)
}

func (getter *businessObjectGetter) Namespace(project, name string, options metav1.GetOptions) (*business.Namespace, error) {
	return getter.businessClient.Namespaces(project).Get(name, options)
}

// ValidateCluster validate cluster
func ValidateCluster(platformClient platforminternalclient.PlatformInterface, clusterName string) field.ErrorList {
	var allErrs field.ErrorList
	if clusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify cluster name"))
	} else {
		_, err := platformClient.Clusters().Get(clusterName, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), clusterName, fmt.Sprintf("can't get cluster:%s", err)))
		}
	}
	return allErrs
}

// ValidateClusterVersioned validate cluster
func ValidateClusterVersioned(getter ClusterGetter, clusterName string, tenantID string) field.ErrorList {
	var allErrs field.ErrorList
	if clusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify cluster name"))
	} else {
		fldPath := field.NewPath("spec", "clusterName")
		cluster, err := getter.Cluster(clusterName, metav1.GetOptions{})
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

// ValidateUpdateCluster validate cluster
func ValidateUpdateCluster(newClusterName, oldClusterName string) field.ErrorList {
	var allErrs field.ErrorList
	if newClusterName != oldClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), newClusterName, "cluster name can't modify"))
	}
	return allErrs
}
