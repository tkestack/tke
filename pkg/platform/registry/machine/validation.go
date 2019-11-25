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

package machine

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"math"
	"net"
	"sync"
	"time"

	apiMachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider"
	"tkestack.io/tke/pkg/platform/util/validation"
	"tkestack.io/tke/pkg/util/ssh"
)

var types = sets.NewString(
	string(platform.BaremetalMachine),
)

// Validate tests if required fields in the machine are set.
func Validate(machineProviders *sync.Map, obj *platform.Machine, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	var err error
	allErrs := apiMachineryvalidation.ValidateObjectMeta(&obj.ObjectMeta, false, apiMachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))

	// validate Type
	if obj.Spec.Type == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "type"), "must specify machine type"))
	} else {
		if !types.Has(string(obj.Spec.Type)) {
			allErrs = append(allErrs, field.NotSupported(field.NewPath("spec", "type"), obj.Spec.Type, types.List()))
		}
		p, err := provider.LoadMachineProvider(machineProviders, string(obj.Spec.Type))
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
		} else {
			resp, err := p.Validate(*obj)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), err))
			}
			allErrs = append(allErrs, resp...)
		}
	}

	// validate ClusterName
	var machineList *platform.MachineList
	var cluster *platform.Cluster
	clusterErrs := validation.ValidateCluster(platformClient, obj.Spec.ClusterName)
	if clusterErrs != nil {
		allErrs = append(allErrs, clusterErrs...)
	} else {
		cluster, err = platformClient.Clusters().Get(obj.Spec.ClusterName, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "clusterName"),
				fmt.Errorf("can't get cluster:%s", err)))
		} else {
			_, cidr, _ := net.ParseCIDR(cluster.Spec.ClusterCIDR)
			ones, _ := cidr.Mask.Size()
			maxNode := math.Exp2(float64(cluster.Status.NodeCIDRMaskSize - int32(ones)))

			fieldSelector := fmt.Sprintf("spec.clusterName=%s", obj.Spec.ClusterName)
			machineList, err = platformClient.Machines().List(metav1.ListOptions{FieldSelector: fieldSelector})
			if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "clusterName"),
					fmt.Errorf("list machines of the cluster error:%s", err)))
			} else {
				machineSize := len(machineList.Items)
				if machineSize >= int(maxNode) {
					allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"),
						fmt.Sprintf("the cluster's machine upper limit(%d) has been reached", int(maxNode))))
				}
			}
		}
	}

	// validate IP
	if obj.Spec.IP == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "IP"), "must specify IP"))
	} else {
		if machineList != nil {
			for _, machine := range machineList.Items {
				if machine.Spec.IP == obj.Spec.IP {
					allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "IP"), obj.Spec.IP))
				}
			}
		}
		if cluster != nil {
			for _, machine := range cluster.Spec.Machines {
				if machine.IP == obj.Spec.IP {
					allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "IP"), obj.Spec.IP))
				}
			}
		}
	}
	if obj.Spec.Port == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "IP"), "must specify Port"))
	}

	// validate ssh
	if obj.Spec.Password == nil && obj.Spec.PrivateKey == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "password"), "password or privateKey at least one"))
	}
	sshConfig := &ssh.Config{
		User:        obj.Spec.Username,
		Host:        obj.Spec.IP,
		Port:        int(obj.Spec.Port),
		Password:    string(obj.Spec.Password),
		PrivateKey:  obj.Spec.PrivateKey,
		PassPhrase:  obj.Spec.PassPhrase,
		DialTimeOut: time.Second,
		Retry:       0,
	}
	s, err := ssh.New(sshConfig)
	if err != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), err.Error()))
	} else {
		err = s.Ping()
		if err != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), err.Error()))
		}
	}

	return allErrs
}

// ValidateUpdate tests if required fields in the cluster are set during
// an update.
func ValidateUpdate(machineProviders *sync.Map, new *platform.Machine, old *platform.Machine) field.ErrorList {
	allErrs := apiMachineryvalidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateUpdateCluster(new.Spec.ClusterName, old.Spec.ClusterName)...)
	// allErrs = append(allErrs, Validate(machineProviders, new, platformClient)...)

	// if new.Spec.Type != "" {
	//	if clusterProviderClient, ok := machineProviders.Load(fmt.Sprintf("%s-new", new.Spec.Type)); ok {
	//		if clusterProvider, ok := clusterProviderClient.(clusterprovider.Provider); ok {
	//			providerErrs := machineProvider.ValidateUpdate(*new, *old)
	//			if providerErrs.RPC != "" {
	//				allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), fmt.Errorf(providerErrs.RPC)))
	//			} else {
	//				for _, e := range providerErrs.Message {
	//					allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), e))
	//				}
	//			}
	//		}
	//	}
	// }

	return allErrs
}
