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

package clustercredential

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// Validate tests if required fields in the cluster are set.
func Validate(obj *platform.ClusterCredential, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateName, field.NewPath("metadata"))
	if obj.ClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("clusterName"), "must specify clusterName"))
	}
	cluster, err := platformClient.Clusters().Get(obj.ClusterName, metav1.GetOptions{})
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), obj.ClusterName, "no such cluster:%s"))
	} else {
		clusterProvider, err := clusterprovider.GetProvider(cluster.Spec.Type)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
		}

		args := clusterprovider.InternalCluster{
			Cluster:           *cluster,
			ClusterCredential: *obj,
		}
		resp, err := clusterProvider.ValidateCredential(args)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
		}
		allErrs = append(allErrs, resp...)
	}

	return allErrs
}

// ValidateUpdate tests if required fields in the core dns are set during an
// update.
func ValidateUpdate(newObj *platform.ClusterCredential, oldObj *platform.ClusterCredential, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&newObj.ObjectMeta, &oldObj.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, Validate(newObj, platformClient)...)

	if newObj.ClusterName != oldObj.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("clusterName"), newObj.ClusterName, "disallowed change the clusterName"))
	}

	return allErrs
}
