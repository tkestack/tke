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
	apimachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"net"
	"sync"
	"time"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider"
	"tkestack.io/tke/pkg/util"
)

// ValidateClusterName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateClusterName = apimachineryValidation.NameIsDNSLabel

var nodePodNumAvails = []int32{16, 32, 64, 128, 256}
var clusterServiceNumAvails = []int32{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}

var types = sets.NewString(
	string(platform.ClusterImported),
	string(platform.ClusterBaremetal),
	string(platform.ClusterEKSHosting),
)

// ValidateCluster tests if required fields in the cluster are set.
func ValidateCluster(clusterProviders *sync.Map, obj *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apimachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateClusterName, field.NewPath("metadata"))

	if obj.Spec.Properties.MaxNodePodNum != nil {
		if !util.InInt32Slice(nodePodNumAvails, *obj.Spec.Properties.MaxNodePodNum) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "properties", "maxNodePodNum"), *obj.Spec.Properties.MaxNodePodNum, fmt.Sprintf("must in values of `%#v`", nodePodNumAvails)))
		}
	}

	if obj.Spec.Properties.MaxClusterServiceNum != nil {
		if !util.InInt32Slice(clusterServiceNumAvails, *obj.Spec.Properties.MaxClusterServiceNum) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "properties", "maxClusterServiceNum"), *obj.Spec.Properties.MaxClusterServiceNum, fmt.Sprintf("must in values of `%#v`", clusterServiceNumAvails)))
		}
	}

	if obj.Spec.Type == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "type"), fmt.Sprintf("availble type are %v", types)))
	} else {
		if !types.Has(string(obj.Spec.Type)) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("spec", "type"), obj.Spec.Type, types.List()))
		} else {
			if obj.Spec.Type == platform.ClusterImported {
				if len(obj.Status.Addresses) == 0 {
					allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses"), "must specify at least one obj access address"))
				} else {
					for _, address := range obj.Status.Addresses {
						if address.Host == "" {
							allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses", string(address.Type), "host"), "must specify the ip of address"))
						}
						if address.Port == 0 {
							allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses", string(address.Type), "port"), "must specify the port of address"))
						}
						_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address.Host, address.Port), 5*time.Second)
						if err != nil {
							allErrs = append(allErrs, field.Invalid(field.NewPath("status", "addresses"), address, err.Error()))
						}
					}
				}
			} else {
				clusterProvider, err := provider.LoadClusterProvider(clusterProviders, string(obj.Spec.Type))
				if err != nil {
					allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
				}

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
func ValidateClusterUpdate(clusterProviders *sync.Map, cluster *platform.Cluster, old *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apimachineryValidation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateCluster(clusterProviders, cluster, platformClient)...)

	if cluster.Spec.Type != "" {
		if cluster.Spec.Type != platform.ClusterImported {
			clusterProvider, err := provider.LoadClusterProvider(clusterProviders, string(cluster.Spec.Type))
			if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
			}

			resp, err := clusterProvider.ValidateUpdate(*cluster, *old)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "annotations"), err))
			}
			allErrs = append(allErrs, resp...)
		}
	}

	return allErrs
}
