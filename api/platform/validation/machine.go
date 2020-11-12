/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package validation

import (
	"context"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsnet "k8s.io/utils/net"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	utilmath "tkestack.io/tke/pkg/util/math"
	"tkestack.io/tke/pkg/util/ssh"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

const MaxTimeOffset = 5 * 300

// ValidateMachine validates a given machine.
func ValidateMachine(ctx context.Context, machine *platform.Machine, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&machine.ObjectMeta, false, apimachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))
	// allErrs = append(allErrs, ValidateMachineSpec(ctx, &machine.Spec, field.NewPath("spec"), platformClient)...)
	// allErrs = append(allErrs, ValidateMachineByProvider(machine)...)

	return allErrs
}

// ValidateMachineUpdate tests if an update to a machine is valid.
func ValidateMachineUpdate(ctx context.Context, machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&machine.ObjectMeta, &oldMachine.ObjectMeta, field.NewPath("metadata"))
	// fldPath := field.NewPath("spec")
	// allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Type, oldMachine.Spec.Type, fldPath.Child("type"))...)
	// allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.ClusterName, oldMachine.Spec.ClusterName, fldPath.Child("clusterName"))...)
	// allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.IP, oldMachine.Spec.IP, fldPath.Child("ip"))...)
	// allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Labels, oldMachine.Spec.Labels, fldPath.Child("labels"))...)
	// allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Taints, oldMachine.Spec.Taints, fldPath.Child("taints"))...)
	// allErrs = append(allErrs, ValidateSSH(fldPath, machine.Spec.IP, int(machine.Spec.Port), machine.Spec.Username, machine.Spec.Password, machine.Spec.PrivateKey, machine.Spec.PassPhrase)...)

	return allErrs
}

// ValidateMachineSpec validates a given machine spec.
func ValidateMachineSpec(ctx context.Context, spec *platform.MachineSpec, fldPath *field.Path, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateMachineSpecType(spec.Type, fldPath.Child("type"))...)
	cluster := new(platform.Cluster)
	allErrs = append(allErrs, ValidateClusterName(ctx, spec.ClusterName, fldPath.Child("clusterName"), cluster, platformClient)...)
	if cluster.Name != "" {
		allErrs = append(allErrs, ValidateMachineWithCluster(ctx, spec.IP, fldPath.Child("ip"), cluster, platformClient)...)
	}
	sshErrors := ValidateSSH(fldPath, spec.IP, int(spec.Port), spec.Username, spec.Password, spec.PrivateKey, spec.PassPhrase)
	if sshErrors != nil {
		allErrs = append(allErrs, sshErrors...)
	} else {
		var masters []*ssh.SSH
		worker, _ := spec.SSH()
		for _, one := range cluster.Spec.Machines {
			master, _ := one.SSH()
			masters = append(masters, master)
		}
		allErrs = append(allErrs, ValidateWorkerTimeOffset(fldPath, worker, masters)...)
	}

	return allErrs
}

// ValidateMachineByProvider validates a given machine by machine provider.
func ValidateMachineByProvider(machine *platform.Machine) field.ErrorList {
	p, err := machineprovider.GetProvider(machine.Spec.Type)
	if err != nil {
		return nil
	}

	return p.Validate(machine)
}

// ValidateMachineSpecType validates a given type and call provider.Validate.
func ValidateMachineSpecType(machineType string, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(machineType, fldPath, machineprovider.Providers())
}

// ValidateWorkerTimeOffset validates a given worker time offset with masters.
func ValidateWorkerTimeOffset(fldPath *field.Path, worker *ssh.SSH, masters []*ssh.SSH) field.ErrorList {
	allErrs := field.ErrorList{}

	workerTimestamp, err := ssh.Timestamp(worker)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
		return allErrs
	}

	times := make([]float64, 0, len(masters))
	for _, one := range masters {
		t, err := ssh.Timestamp(one)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
			return allErrs
		}
		times = append(times, float64(t))
	}
	minIndex, minTime := utilmath.Min(times)
	offset := workerTimestamp - int(*minTime)
	if offset > MaxTimeOffset {
		allErrs = append(allErrs, field.Invalid(fldPath, worker.Host,
			fmt.Sprintf("the time offset(%v-%v=%v) between node(%v) with node(%v) exceeds %d seconds, please unify machine time between nodes by using ntp or manual", workerTimestamp, int(*minTime), offset, worker.Host, masters[*minIndex].Host, MaxTimeOffset)))
	}

	return allErrs
}

// ValidateSSH validates a given ssh config.
func ValidateSSH(fldPath *field.Path, ip string, port int, user string, password []byte, privateKey []byte, passPhrase []byte) field.ErrorList {
	allErrs := field.ErrorList{}

	for _, msg := range validation.IsValidIP(ip) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ip"), ip, msg))

	}
	for _, msg := range validation.IsValidPortNum(port) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), port, msg))
	}
	if password == nil && privateKey == nil {
		allErrs = append(allErrs, field.Required(fldPath, "must specify password or privateKey"))
	}

	if len(allErrs) != 0 {
		return allErrs
	}

	sshConfig := &ssh.Config{
		User:        user,
		Host:        ip,
		Port:        port,
		Password:    string(password),
		PrivateKey:  privateKey,
		PassPhrase:  passPhrase,
		DialTimeOut: time.Second,
		Retry:       0,
	}
	s, err := ssh.New(sshConfig)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, "", err.Error()))
	} else {
		output, err := s.CombinedOutput("whoami")
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, "", err.Error()))
		}
		if strings.TrimSpace(string(output)) != "root" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("user"), user, `must be root or set sudo without password`))
		}
	}

	return allErrs
}

// ValidateMachineWithCluster validates a given machine by ip with cluster.
func ValidateMachineWithCluster(ctx context.Context, ip string, fldPath *field.Path, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, machine := range cluster.Spec.Machines {
		if machine.IP == ip {
			allErrs = append(allErrs, field.Duplicate(fldPath, ip))
		}
	}
	cidrs := strings.Split(cluster.Spec.ClusterCIDR, ",")
	for _, cidr := range cidrs {
		if utilsnet.IsIPv6CIDRString(cidr) {
			return allErrs
		}
	}

	_, cidr, _ := net.ParseCIDR(cluster.Spec.ClusterCIDR)
	ones, _ := cidr.Mask.Size()
	maxNode := math.Exp2(float64(cluster.Status.NodeCIDRMaskSize - int32(ones)))

	fieldSelector := fmt.Sprintf("spec.clusterName=%s", cluster.Name)
	machineList, err := platformClient.Machines().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
	} else {
		machineSize := len(machineList.Items)
		if machineSize >= int(maxNode) {
			allErrs = append(allErrs, field.Forbidden(fldPath, fmt.Sprintf("the cluster's machine upper limit(%d) has been reached", int(maxNode))))
		}
	}
	for _, machine := range machineList.Items {
		if machine.Spec.IP == ip {
			allErrs = append(allErrs, field.Duplicate(fldPath, ip))
		}
	}

	return allErrs
}

// ValidateClusterName validates a given clusterName and return cluster if exists.
func ValidateClusterName(ctx context.Context, clusterName string, fldPath *field.Path, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := field.ErrorList{}

	if clusterName == "" {
		allErrs = append(allErrs, field.Required(fldPath, "must specify cluster name"))
	} else {
		c, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				allErrs = append(allErrs, field.NotFound(fldPath, clusterName))
			} else {
				allErrs = append(allErrs, field.InternalError(fldPath, err))
			}
		}
		*cluster = *c
	}

	return allErrs
}
