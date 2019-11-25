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

package helm

import (
	"fmt"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/controller/addon/helm/images"
	"tkestack.io/tke/pkg/platform/util/validation"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// ValidateHelm tests if required fields in the cluster are set.
func ValidateHelm(obj *platform.Helm, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateName, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateCluster(platformClient, obj.Spec.ClusterName)...)
	allErrs = append(allErrs, validateSpec(platformClient, obj.Spec.Version, obj.Spec.ClusterName)...)
	return allErrs
}

func validateSpec(platformInterface platforminternalclient.PlatformInterface, version string, clusterName string) field.ErrorList {
	var allErrs field.ErrorList
	allErrs = append(allErrs, validateComponentVersion(version)...)
	allErrs = append(allErrs, validateUniqueAddon(platformInterface, clusterName)...)
	return allErrs
}

// validateComponentVersion validate version exists in image version map
func validateComponentVersion(version string) field.ErrorList {
	var allErrs field.ErrorList
	if err := images.Validate(version); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "version"), version, "the component version definition could not be found"))
	}
	return allErrs
}

// validateUniqueAddon validate only one addon in the same cluster
func validateUniqueAddon(platformClient platforminternalclient.PlatformInterface, clusterName string) field.ErrorList {
	var allErrs field.ErrorList
	fieldSelector := fmt.Sprintf("spec.clusterName=%s", clusterName)
	helms, err := platformClient.Helms().List(metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "clusterName"),
			fmt.Errorf("list helms of the cluster error:%s", err)))
	}
	if len(helms.Items) > 0 {
		allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "clusterName"), clusterName))
	}
	return allErrs
}

// ValidateHelmUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateHelmUpdate(new *platform.Helm, old *platform.Helm, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMeta(&new.ObjectMeta, false, ValidateName, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validation.ValidateCluster(platformClient, new.Spec.ClusterName)...)
	allErrs = append(allErrs, validateSpecUpdate(new, old)...)
	return allErrs
}

func validateSpecUpdate(new *platform.Helm, old *platform.Helm) field.ErrorList {
	var allErrs field.ErrorList
	if new.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), new.Spec.ClusterName, "disallowed change the cluster name"))
	}
	if new.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), new.Spec.TenantID, "disallowed change the tenant"))
	}
	if new.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(new.Status.Phase)))
	}
	allErrs = append(allErrs, validateComponentVersion(new.Spec.Version)...)
	return allErrs
}
