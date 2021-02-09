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

package node_test

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/apiserver/cluster"
	tke2 "tkestack.io/tke/test/tke"
	testclient "tkestack.io/tke/test/util/client"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
	"tkestack.io/tke/test/util/env"
)

var (
	provider = tencent.NewTencentProvider()
	testTKE  *tke2.TestTKE
	cls      *platformv1.Cluster
	err      error
)

var _ = Describe("node", func() {

	SynchronizedBeforeSuite(func() []byte {
		// Download and install tke-installer
		installer := tke2.InitInstaller(provider)
		err := installer.InstallInstaller(os.Getenv("OS"), os.Getenv("ARCH"), os.Getenv("VERSION"))
		Expect(err).Should(BeNil(), "Install tke-installer failed")

		// Install tkestack with one master node
		nodes, err := provider.CreateInstances(1)
		Expect(err).Should(BeNil(), "Create instance failed")

		para := installer.CreateClusterParaTemplate(nodes)
		err = installer.Install(para)
		Expect(err).To(BeNil(), "Install failed")
		kubeconfig, err := testclient.GenerateTKEAdminKubeConfig(nodes[0])
		Expect(err).Should(BeNil(), "Generate tke admin kubeconfig failed")

		tkeClient, err := testclient.GetTKEClient([]byte(kubeconfig))
		Expect(err).Should(BeNil(), "Get tke client with admin token failed")
		testTKE = tke2.Init(tkeClient, provider)

		cls, err = testTKE.CreateCluster()
		Expect(err).To(BeNil(), "Create cluster failed")

		temp := tkeAndCluster{
			Tke: testTKE,
			Cls: cls,
		}
		data, err := json.Marshal(temp)
		Expect(err).Should(BeNil())
		return data
	}, func(data []byte) {
		temp := new(tkeAndCluster)
		Expect(json.Unmarshal(data, temp)).Should(BeNil())
		testTKE = temp.Tke
		cls = temp.Cls
	})

	SynchronizedAfterSuite(func() {
		if !env.NeedDelete() {
			return
		}
		Expect(provider.TearDown()).Should(BeNil())
	}, func() {})

	Context("Node", func() {
		var machine *platformv1.Machine

		BeforeEach(func() {
			nodes, err := testTKE.CreateInstances(1)
			Expect(err).To(BeNil(), "Create instance failed")
			machine, err = testTKE.AddNode(cls.Name, nodes[0])
			Expect(err).To(BeNil(), "Add node to cluster failed")
		})

		AfterEach(func() {
			if machine != nil {
				Expect(testTKE.DeleteNode(machine.Name)).Should(Succeed())
			}
		})

		It("Add and delete node", func() {
			// Add operation has been covered in BeforeEach
			Expect(testTKE.DeleteNode(machine.Name)).Should(Succeed())
		})

		It("Add label to node", func() {
			labelKey := "testLabelKey"
			labelValue := machine.Spec.IP

			Eventually(func() error {
				machine.Labels = map[string]string{
					labelKey: labelValue,
				}
				_, err = testTKE.TkeClient.PlatformV1().Machines().Update(context.Background(), machine, metav1.UpdateOptions{})
				if err != nil {
					// Get the latest machine object again if update failed
					machine, _ = testTKE.TkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
				}
				return err
			}, 5*time.Second, time.Second).Should(Succeed())
			machine, _ = testTKE.TkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
			Expect(machine.Labels).Should(HaveKeyWithValue(labelKey, labelValue))
		})

		It("Unschedulable node", func() {
			// Unschedule node
			node, _ := testTKE.UnscheduleNode(cls, machine.Spec.IP)
			Expect(node.Spec.Unschedulable).Should(BeTrue())

			// Cancel unschedule node
			node, _ = testTKE.CancleUnschedulableNode(cls, machine.Spec.IP)
			Expect(node.Spec.Unschedulable).Should(BeFalse())
		})

		It("Drain node", func() {
			k8sClient := testTKE.K8sClient(cls)

			node, err := k8sClient.CoreV1().Nodes().Get(context.Background(), machine.Spec.IP, metav1.GetOptions{})
			Expect(err).Should(BeNil())

			// Drain node
			Expect(cluster.DrainNode(context.Background(), k8sClient, node)).Should(BeNil(), "Drain node failed")

			node, _ = k8sClient.CoreV1().Nodes().Get(context.Background(), machine.Spec.IP, metav1.GetOptions{})
			Expect(node.Spec.Unschedulable).Should(BeTrue(), "Node was not unschedulable after draining")
		})
	})

	Context("Addon", func() {
		It("TappController", func() {
			tapp := &platformv1.TappController{
				Spec: platformv1.TappControllerSpec{
					ClusterName: cls.Name,
				},
			}
			tapp, err := testTKE.TkeClient.PlatformV1().TappControllers().Create(context.Background(), tapp, metav1.CreateOptions{})
			Expect(err).Should(BeNil())

			Eventually(func() error {
				addon, err := testTKE.TkeClient.PlatformV1().TappControllers().Get(context.Background(), tapp.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if addon.Status.Phase != "Running" {
					return errors.New(addon.Name + " Phase: " + string(addon.Status.Phase) + ", Reason: " + addon.Status.Reason)
				}
				return nil
			}, 10*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("IPAM", func() {
			ipam := &platformv1.IPAM{
				Spec: platformv1.IPAMSpec{
					ClusterName: cls.Name,
				},
			}
			ipam, err := testTKE.TkeClient.PlatformV1().IPAMs().Create(context.Background(), ipam, metav1.CreateOptions{})
			Expect(err).Should(BeNil())

			Eventually(func() error {
				addon, err := testTKE.TkeClient.PlatformV1().IPAMs().Get(context.Background(), ipam.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if addon.Status.Phase != "Running" {
					return errors.New(addon.Name + " Phase: " + string(addon.Status.Phase) + ", Reason: " + addon.Status.Reason)
				}
				return nil
			}, 10*time.Minute, 10*time.Second).Should(Succeed())
		})

		It("CronHPA", func() {
			cronHPA := &platformv1.CronHPA{
				Spec: platformv1.CronHPASpec{
					ClusterName: cls.Name,
				},
			}
			cronHPA, err := testTKE.TkeClient.PlatformV1().CronHPAs().Create(context.Background(), cronHPA, metav1.CreateOptions{})
			Expect(err).Should(BeNil())

			Eventually(func() error {
				addon, err := testTKE.TkeClient.PlatformV1().CronHPAs().Get(context.Background(), cronHPA.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if addon.Status.Phase != "Running" {
					return errors.New(addon.Name + " Phase: " + string(addon.Status.Phase) + ", Reason: " + addon.Status.Reason)
				}
				return nil
			}, 10*time.Minute, 10*time.Second).Should(Succeed())
		})
	})
})

type tkeAndCluster struct {
	Tke *tke2.TestTKE
	Cls *platformv1.Cluster
}
