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

package resource

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
)

// Name is the name identifying various resources in a ResourceList.
type Name string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	CPU Name = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	Memory Name = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	Storage Name = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// The resource name for EphemeralStorage is alpha and it can change across releases.
	EphemeralStorage Name = "ephemeral-storage"
)

// The following identify resource constants for Kubernetes object types
const (
	// Pods, number
	Pods Name = "pods"
	// Services, number
	Services Name = "services"
	// ReplicationControllers, number
	ReplicationControllers Name = "replicationcontrollers"
	// Quotas, number
	Quotas Name = "resourcequotas"
	// Secrets, number
	Secrets Name = "secrets"
	// ConfigMaps, number
	ConfigMaps Name = "configmaps"
	// PersistentVolumeClaims, number
	PersistentVolumeClaims Name = "persistentvolumeclaims"
	// ServicesNodePorts, number
	ServicesNodePorts Name = "services.nodeports"
	// ServicesLoadBalancers, number
	ServicesLoadBalancers Name = "services.loadbalancers"
	// CPU request, in cores. (500m = .5 cores)
	RequestsCPU Name = "requests.cpu"
	// Memory request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	RequestsMemory Name = "requests.memory"
	// Storage request, in bytes
	RequestsStorage Name = "requests.storage"
	// Local ephemeral storage request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	RequestsEphemeralStorage Name = "requests.ephemeral-storage"
	// CPU limit, in cores. (500m = .5 cores)
	LimitsCPU Name = "limits.cpu"
	// Memory limit, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	LimitsMemory Name = "limits.memory"
	// Local ephemeral storage limit, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	LimitsEphemeralStorage Name = "limits.ephemeral-storage"
)

const (
	// DefaultNamespacePrefix namespace prefix.
	DefaultNamespacePrefix = "kubernetes.io/"
	// HugePagesPrefix is huge pages prefix.
	HugePagesPrefix = "hugepages-"
)

// The following identify resource prefix for Kubernetes object types
const (
	// RequestsHugePagesPrefix request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// As burst is not supported for HugePages, we would only quota its request, and ignore the limit.
	RequestsHugePagesPrefix = "requests.hugepages-"
	// Default resource requests prefix
	DefaultResourceRequestsPrefix = "requests."
)

var standardQuotaResources = sets.NewString(
	string(CPU),
	string(Memory),
	string(EphemeralStorage),
	string(RequestsCPU),
	string(RequestsMemory),
	string(RequestsStorage),
	string(RequestsEphemeralStorage),
	string(LimitsCPU),
	string(LimitsMemory),
	string(LimitsEphemeralStorage),
	string(Pods),
	string(Quotas),
	string(Services),
	string(ReplicationControllers),
	string(Secrets),
	string(PersistentVolumeClaims),
	string(ConfigMaps),
	string(ServicesNodePorts),
	string(ServicesLoadBalancers),
)

// IsStandardQuotaResourceName returns true if the resource is known to
// the quota tracking system
func IsStandardQuotaResourceName(str string) bool {
	return standardQuotaResources.Has(str) || IsQuotaHugePageResourceName(Name(str))
}

var standardResources = sets.NewString(
	string(CPU),
	string(Memory),
	string(EphemeralStorage),
	string(RequestsCPU),
	string(RequestsMemory),
	string(RequestsEphemeralStorage),
	string(LimitsCPU),
	string(LimitsMemory),
	string(LimitsEphemeralStorage),
	string(Pods),
	string(Quotas),
	string(Services),
	string(ReplicationControllers),
	string(Secrets),
	string(ConfigMaps),
	string(PersistentVolumeClaims),
	string(Storage),
	string(RequestsStorage),
	string(ServicesNodePorts),
	string(ServicesLoadBalancers),
)

// IsStandardResourceName returns true if the resource is known to the system
func IsStandardResourceName(str string) bool {
	return standardResources.Has(str) || IsQuotaHugePageResourceName(Name(str))
}

// IsQuotaHugePageResourceName returns true if the resource name has the quota
// related huge page resource prefix.
func IsQuotaHugePageResourceName(name Name) bool {
	return strings.HasPrefix(string(name), HugePagesPrefix) || strings.HasPrefix(string(name), RequestsHugePagesPrefix)
}

var integerResources = sets.NewString(
	string(Pods),
	string(Quotas),
	string(Services),
	string(ReplicationControllers),
	string(Secrets),
	string(ConfigMaps),
	string(PersistentVolumeClaims),
	string(ServicesNodePorts),
	string(ServicesLoadBalancers),
)

// IsIntegerResourceName returns true if the resource is measured in integer values
func IsIntegerResourceName(str string) bool {
	return integerResources.Has(str) || IsExtendedResourceName(Name(str))
}

// IsExtendedResourceName returns true if:
// 1. the resource name is not in the default namespace;
// 2. resource name does not have "requests." prefix,
// to avoid confusion with the convention in quota
// 3. it satisfies the rules in IsQualifiedName() after converted into quota resource name
func IsExtendedResourceName(name Name) bool {
	if IsNativeResource(name) || strings.HasPrefix(string(name), DefaultResourceRequestsPrefix) {
		return false
	}
	// Ensure it satisfies the rules in IsQualifiedName() after converted into quota resource name
	nameForQuota := fmt.Sprintf("%s%s", DefaultResourceRequestsPrefix, string(name))
	if errs := validation.IsQualifiedName(nameForQuota); len(errs) != 0 {
		return false
	}
	return true
}

// IsNativeResource returns true if the resource name is in the
// *kubernetes.io/ namespace. Partially-qualified (unprefixed) names are
// implicitly in the kubernetes.io/ namespace.
func IsNativeResource(name Name) bool {
	return !strings.Contains(string(name), "/") ||
		strings.Contains(string(name), DefaultNamespacePrefix)
}
