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

package platform_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"tkestack.io/tke/test/e2e/tke"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/test/util/cloudprovider"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
)

var t = tke.TKE{Namespace: "platform"}

var _ = BeforeSuite(func() {
	t.Create()
})

var _ = AfterSuite(func() {
	t.Delete()
})

var _ = Describe("cluster lifecycle", func() {
	var provider cloudprovider.Provider

	var masterNodes []cloudprovider.Instance
	var workerNodes []cloudprovider.Instance

	BeforeEach(func() {
		var err error
		provider = tencent.NewTencentProvider()
		masterNodes, err = provider.CreateInstances(2)
		Expect(err).To(BeNil())
		workerNodes, err = provider.CreateInstances(1)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		if os.Getenv("NEED_DELETE") == "" {
			return
		}
		var instanceIDs []*string
		for i, one := range masterNodes {
			fmt.Printf("delete instance %d %s\n", i, one.InternalIP)
			instanceIDs = append(instanceIDs, &masterNodes[i].InstanceID)
		}
		for i, one := range workerNodes {
			fmt.Printf("delete instance %d %s\n", i, one.InternalIP)
			instanceIDs = append(instanceIDs, &workerNodes[i].InstanceID)
		}
		err := provider.DeleteInstances(instanceIDs)
		Expect(err).To(BeNil())
	})

	It("complete cluster lifecycle", func() {
		By("create baremetal cluster")
		restconf, err := t.GetKubeConfig()
		Expect(err).To(BeNil())
		client := tkeclientset.NewForConfigOrDie(restconf)
		err = createBaremetalCluster(context.Background(), client, masterNodes, workerNodes)
		Expect(err).To(BeNil())
	})
})

func createBaremetalCluster(ctx context.Context, client *tkeclientset.Clientset, masterNodes []cloudprovider.Instance,
	workerNodes []cloudprovider.Instance) error {
	time.Sleep(70 * time.Second)
	cluster := &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			Type:     "Baremetal",
			Features: platformv1.ClusterFeature{EnableMasterSchedule: true},
		}}
	for _, one := range masterNodes {
		cluster.Spec.Machines = append(cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}

	By("create cluster")
	cluster, err := client.PlatformV1().Clusters().Create(context.Background(), cluster, metav1.CreateOptions{})
	Expect(err).To(BeNil())
	clusterName := cluster.Name

	By(fmt.Sprintf("wait cluster(%s) status is running", clusterName))
	err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
		cluster, err = client.PlatformV1().Clusters().Get(context.Background(), clusterName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return cluster.Status.Phase == platformv1.ClusterRunning, nil
	})

	By("create wokers")
	for _, one := range workerNodes {
		machine := &platformv1.Machine{
			Spec: platformv1.MachineSpec{
				ClusterName: clusterName,
				Type:        "Baremetal",
				IP:          one.InternalIP,
				Port:        one.Port,
				Username:    one.Username,
				Password:    []byte(one.Password),
			},
		}

		By(fmt.Sprintf("add work nodes(%s) to cluster", one.InternalIP))
		machine, err = client.PlatformV1().Machines().Create(context.Background(), machine, metav1.CreateOptions{})
		Expect(err).To(BeNil())
		machineName := machine.Name

		By(fmt.Sprintf("wait worker node(%s) status is running", one.InternalIP))
		err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
			machine, err = client.PlatformV1().Machines().Get(context.Background(), machineName, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			return machine.Status.Phase == platformv1.MachineRunning, nil
		})

		By(fmt.Sprintf("delete work nodes(%s) from cluster", one.InternalIP))
		err = client.PlatformV1().Machines().Delete(context.Background(), machineName, metav1.DeleteOptions{})
		Expect(err).To(BeNil())
		err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
			_, err = client.PlatformV1().Machines().Get(context.Background(), machineName, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, nil
		})
	}

	By(fmt.Sprintf("delete cluster(%s)", clusterName))
	err = client.PlatformV1().Clusters().Delete(context.Background(), clusterName, metav1.DeleteOptions{})
	Expect(err).To(BeNil())
	err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
		_, err = client.PlatformV1().Clusters().Get(context.Background(), clusterName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, nil
	})
	Expect(err).To(BeNil())
	return nil
}
