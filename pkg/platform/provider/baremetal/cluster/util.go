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
	"math"
	"net"
	"reflect"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	utilsnet "k8s.io/utils/net"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/util/ipallocator"
	"tkestack.io/tke/pkg/util/log"
)

func GetNodeCIDRMaskSize(clusterCIDR string, maxNodePodNum int32) (int32, error) {
	if maxNodePodNum <= 0 {
		return 0, errors.New("maxNodePodNum must more than 0")
	}
	_, svcSubnetCIDR, err := net.ParseCIDR(clusterCIDR)
	if err != nil {
		return 0, errors.Wrap(err, "ParseCIDR error")
	}

	nodeCidrOccupy := math.Ceil(math.Log2(float64(maxNodePodNum)))
	nodeCIDRMaskSize := 32 - int(nodeCidrOccupy)
	ones, _ := svcSubnetCIDR.Mask.Size()
	if ones > nodeCIDRMaskSize {
		return 0, errors.New("clusterCIDR IP size is less than maxNodePodNum")
	}

	return int32(nodeCIDRMaskSize), nil
}

func GetServiceCIDRAndNodeCIDRMaskSize(clusterCIDR string, maxClusterServiceNum int32, maxNodePodNum int32) (string, int32, error) {
	if maxClusterServiceNum <= 0 || maxNodePodNum <= 0 {
		return "", 0, errors.New("maxClusterServiceNum or maxNodePodNum must more than 0")
	}
	_, svcSubnetCIDR, err := net.ParseCIDR(clusterCIDR)
	if err != nil {
		return "", 0, errors.Wrap(err, "ParseCIDR error")
	}

	size := ipallocator.RangeSize(svcSubnetCIDR)
	if int32(size) < maxClusterServiceNum {
		return "", 0, errors.New("clusterCIDR IP size is less than maxClusterServiceNum")
	}
	lastIP, err := ipallocator.GetIndexedIP(svcSubnetCIDR, int(size-1))
	if err != nil {
		return "", 0, errors.Wrap(err, "get last IP error")
	}

	maskSize := int(math.Ceil(math.Log2(float64(maxClusterServiceNum))))
	_, serviceCidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", lastIP.String(), 32-maskSize))

	nodeCidrOccupy := math.Ceil(math.Log2(float64(maxNodePodNum)))
	nodeCIDRMaskSize := 32 - int(nodeCidrOccupy)
	ones, _ := svcSubnetCIDR.Mask.Size()
	if ones > nodeCIDRMaskSize {
		return "", 0, errors.New("clusterCIDR IP size is less than maxNodePodNum")
	}

	return serviceCidr.String(), int32(nodeCIDRMaskSize), nil
}

func GetIndexedIP(subnet string, index int) (net.IP, error) {
	_, svcSubnetCIDR, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't parse service subnet CIDR %q", subnet)
	}

	dnsIP, err := ipallocator.GetIndexedIP(svcSubnetCIDR, index)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get %dth IP address from service subnet CIDR %s", index, svcSubnetCIDR.String())
	}

	return dnsIP, nil
}

// GetAPIServerCertSANs returns extra APIServer's certSANs need to pass kubeadm
func GetAPIServerCertSANs(c *platformv1.Cluster) []string {
	certSANs := sets.NewString("127.0.0.1", "localhost", "::1")
	certSANs = certSANs.Insert(c.Spec.PublicAlternativeNames...)
	if c.Spec.Features.HA != nil {
		if c.Spec.Features.HA.TKEHA != nil {
			certSANs.Insert(c.Spec.Features.HA.TKEHA.VIP)
		}
		if c.Spec.Features.HA.ThirdPartyHA != nil {
			certSANs.Insert(c.Spec.Features.HA.ThirdPartyHA.VIP)
		}
	}
	for _, address := range c.Status.Addresses {
		certSANs.Insert(address.Host)
	}

	return certSANs.List()
}

func CalcNodeCidrSize(podSubnet string) (int32, bool) {
	maskSize := 24
	isIPv6 := false
	if ip, podCidr, err := net.ParseCIDR(podSubnet); err == nil {
		if utilsnet.IsIPv6(ip) {
			var nodeCidrSize int
			isIPv6 = true
			podNetSize, totalBits := podCidr.Mask.Size()
			switch {
			case podNetSize == 112:
				// Special case, allows 256 nodes, 256 pods/node
				nodeCidrSize = 120
			case podNetSize < 112:
				// Use multiple of 8 for node CIDR, with 512 to 64K nodes
				nodeCidrSize = totalBits - ((totalBits-podNetSize-1)/8-1)*8
			default:
				// Not enough bits, will fail later, when validate
				nodeCidrSize = podNetSize
			}
			maskSize = nodeCidrSize
		}
	}
	return int32(maskSize), isIPv6
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
		return scalingMachines, errors.Errorf("scale up and down master in parallel is not allowed: %v, %v", diff1.ToSlice(), diff2.ToSlice())
	}
	if diff1.Cardinality() > 0 {
		if diff1.Contains(oldCluster.Spec.Machines[0].IP) {
			return scalingMachines, errors.Errorf("master[0] can't scale down: %v", oldCluster.Spec.Machines[0].IP)
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
