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
	"fmt"

	"github.com/thoas/go-funk"
	apimachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

// ValidateClusterName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateClusterName = apimachineryValidation.NameIsDNSLabel

// ValidateCluster tests if required fields in the cluster are set.
func ValidateCluster(obj *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apimachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateClusterName, field.NewPath("metadata"))

	types := clusterprovider.Providers()
	if obj.Spec.Type == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "type"), fmt.Sprintf("availble type are %v", types)))
	} else {
		if !funk.ContainsString(types, obj.Spec.Type) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("spec", "type"), obj.Spec.Type, types))
		} else {
			clusterProvider, err := clusterprovider.GetProvider(obj.Spec.Type)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
			} else {
				resp, err := clusterProvider.Validate(*obj)
				if err != nil {
					allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
				}
				allErrs = append(allErrs, resp...)
			}
		}
	}

	return allErrs
}

// ValidateClusterUpdate tests if required fields in the cluster are set during
// an update.
func ValidateClusterUpdate(cluster *platform.Cluster, old *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apimachineryValidation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateCluster(cluster, platformClient)...)

	return allErrs
}
