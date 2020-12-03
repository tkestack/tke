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

package ipam

import (
	"context"

	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util/validation"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// ValidateIPAM tests if required fields in the cluster are set.
func ValidateIPAM(ctx context.Context, platformClient platforminternalclient.PlatformInterface, obj *platform.IPAM) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateName, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateCluster(ctx, platformClient, obj.Spec.ClusterName)...)

	if len(obj.Spec.ClusterName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify a cluster name"))
	}

	return allErrs
}

// ValidateIPAMUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateIPAMUpdate(ctx context.Context, platformClient platforminternalclient.PlatformInterface, obj *platform.IPAM, old *platform.IPAM) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&obj.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateIPAM(ctx, platformClient, obj)...)

	if obj.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), obj.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if obj.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), obj.Spec.TenantID, "disallowed change the tenant"))
	}
	if obj.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(obj.Status.Phase)))
	}

	return allErrs
}
