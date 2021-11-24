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

package util

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	mapset "github.com/deckarep/golang-set"
	pkgerrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/util/credential"
	"tkestack.io/tke/pkg/util/log"
)

func DynamicClientByCluster(ctx context.Context, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) (dynamic.Interface, error) {
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}

	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return nil, err
	}

	restConfig, err := provider.GetRestConfig(ctx, platformClient, cluster, username)
	if err != nil {
		return nil, err
	}
	if cluster.Status.Phase != platform.ClusterRunning {
		return nil, fmt.Errorf("cluster %s status is abnormal", cluster.ObjectMeta.Name)
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}
	return dynamic.NewForConfig(restConfig)

}

// ClientSetByCluster returns the backend kubernetes clientSet by given cluster object
func ClientSetByCluster(ctx context.Context, cluster *platform.Cluster, platformClient platforminternalclient.PlatformInterface) (*kubernetes.Clientset, error) {
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}

	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return nil, err
	}

	restConfig, err := provider.GetRestConfig(ctx, platformClient, cluster, username)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}

// ResourceFromKind returns the resource name by kind.
func ResourceFromKind(kind string) string {
	kindLower := strings.ToLower(kind)
	switch kindLower {
	case "policy":
		return "policies"
	case "storageclass":
		return "storageclasses"
	case "ingress":
		return "ingresses"
	case "networkpolicy":
		return "networkpolicies"
	case "podsecuritypolicy":
		return "podsecuritypolicies"
	case "priorityclass":
		return "priorityclasses"
	case "endpoints":
		return "endpoints"
	default:
		return kindLower + "s"
	}
}

// BuildVersionedClientSet creates the clientset of kubernetes by given
// cluster object and returns it.
func BuildVersionedClientSet(cluster *platformv1.Cluster, cc *platformv1.ClusterCredential) (*kubernetes.Clientset, error) {
	restConfig := cc.RESTConfig(cluster)
	return kubernetes.NewForConfig(restConfig)
}

// BuildExternalClientSet creates the clientset of kubernetes by given cluster
// object and returns it.
func BuildExternalClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubernetes.Clientset, error) {
	cc, err := credential.GetClusterCredentialV1(ctx, client, cluster, clusterprovider.AdminUsername)
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return BuildVersionedClientSet(cluster, cc)
}

// BuildExternalClientSetWithName creates the clientset of kubernetes by given cluster
// name and returns it.
func BuildExternalClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (*kubernetes.Clientset, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func clusterAddress(cluster *platform.Cluster) (*platform.ClusterAddress, error) {
	addrs := make(map[platform.AddressType][]platform.ClusterAddress)
	for _, one := range cluster.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *platform.ClusterAddress
	if len(addrs[platform.AddressInternal]) != 0 {
		address = &addrs[platform.AddressInternal][rand.Intn(len(addrs[platform.AddressInternal]))]
	} else if len(addrs[platform.AddressAdvertise]) != 0 {
		address = &addrs[platform.AddressAdvertise][rand.Intn(len(addrs[platform.AddressAdvertise]))]
	} else {
		if len(addrs[platform.AddressReal]) != 0 {
			address = &addrs[platform.AddressReal][rand.Intn(len(addrs[platform.AddressReal]))]
		}
	}
	if address == nil {
		return nil, pkgerrors.New("no valid address for the cluster")
	}

	return address, nil
}

func PrepareClusterScale(cluster *platform.Cluster, oldCluster *platform.Cluster) ([]platform.ClusterMachine, error) {
	allMachines, scalingMachines := []platform.ClusterMachine{}, []platform.ClusterMachine{}

	oIPs := mapset.NewSet()
	for _, machine := range oldCluster.Spec.Machines {
		oIPs.Add(machine.IP)
		allMachines = append(allMachines, machine)
	}
	IPs := mapset.NewSet()
	for _, machine := range cluster.Spec.Machines {
		IPs.Add(machine.IP)
		allMachines = append(allMachines, machine)
	}
	// nothing to do since ips not changed
	if reflect.DeepEqual(oIPs, IPs) {
		return scalingMachines, nil
	}
	// machine in oldCluster but not in cluster
	diff1 := oIPs.Difference(IPs)
	// machine in cluster but not in oldCluster
	diff2 := IPs.Difference(oIPs)
	// scaling machine ips
	diff := mapset.NewSet()
	log.Errorf("PrepareClusterScale called: diff1 -> %v, diff2 -> %v", diff1.ToSlice(), diff2.ToSlice())
	if diff1.Cardinality() > 0 && diff2.Cardinality() > 0 {
		return scalingMachines, pkgerrors.Errorf("scale up and down master in parallel is not allowed: %v, %v", diff1.ToSlice(), diff2.ToSlice())
	}
	if diff1.Cardinality() > 0 {
		if diff1.Contains(oldCluster.Spec.Machines[0].IP) {
			return scalingMachines, pkgerrors.Errorf("master[0] can't scale down: %v", oldCluster.Spec.Machines[0].IP)
		}
		diff = diff1
	}
	if diff2.Cardinality() > 0 {
		diff = diff2
	}
	for _, m := range diff.ToSlice() {
		for _, machine := range allMachines {
			if m == machine.IP {
				scalingMachines = append(scalingMachines, machine)
				break
			}
		}
	}
	return scalingMachines, nil
}
