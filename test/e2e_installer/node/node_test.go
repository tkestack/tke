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
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/apiserver/cluster"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	baremetalmachine "tkestack.io/tke/pkg/platform/provider/baremetal/machine"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	importedcluster "tkestack.io/tke/pkg/platform/provider/imported/cluster"
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

func copyProviderConfig() {
	wd, _ := os.Getwd()
	log.Printf("current test dir is %v", wd)
	err := os.MkdirAll(path.Dir(baremetalconstants.ConfDir), 0755)
	if err != nil {
		log.Fatalf("create dir failed: %v", err)
	}
	input, err := ioutil.ReadFile("../../../pkg/platform/" + baremetalconstants.ConfigFile)
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}
	err = ioutil.WriteFile(baremetalconstants.ConfigFile, input, 0755)
	if err != nil {
		log.Fatalf("write config failed: %v", err)
	}
}

var _ = Describe("node", func() {
	copyProviderConfig()
	baremetalProvider, _ := baremetalcluster.NewProvider()
	baremetalmachine.RegisterProvider()
	importedcluster.RegisterProvider()

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

		// Prepare a cluster
		tkeClient, err := testclient.GetTKEClient([]byte(kubeconfig))
		Expect(err).Should(BeNil(), "Get tke client with admin token failed")
		testTKE := tke2.Init(tkeClient, provider)
		cls, err := testTKE.CreateCluster()
		Expect(err).To(BeNil(), "Create cluster failed")
		baremetalProvider.PlatformClient = testTKE.TkeClient.PlatformV1()
		clusterprovider.Register(baremetalProvider.Name(), baremetalProvider)

		return []byte(cls.Name + ";" + kubeconfig)
	}, func(data []byte) {
		temp := strings.Split(string(data), ";")

		tkeClient, err := testclient.GetTKEClient([]byte(temp[1]))
		Expect(err).Should(BeNil(), "Get tke client with admin token failed")
		testTKE = tke2.Init(tkeClient, provider)
		clusters, _ := testTKE.TkeClient.PlatformV1().Clusters().List(context.Background(), metav1.ListOptions{})
		for _, c := range clusters.Items {
			if c.Name == temp[0] {
				cls = &c
			}
		}
		Expect(cls).ShouldNot(BeNil(), "Cluster %v was not found", temp[0])
	})

	SynchronizedAfterSuite(func() {
		if env.NeedDelete() {
			Expect(provider.TearDown()).Should(BeNil())
		}
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
				testTKE.DeleteNode(machine.Name)
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
			Expect(testTKE.TkeClient.PlatformV1().TappControllers().Delete(context.Background(), tapp.Name, metav1.DeleteOptions{})).Should(BeNil(), "Delete TappController failed")
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
			Expect(testTKE.TkeClient.PlatformV1().CronHPAs().Delete(context.Background(), cronHPA.Name, metav1.DeleteOptions{})).Should(BeNil(), "Delete CronHPA failed")
		})

		It("CSIOperator", func() {
			csiOperator := &platformv1.CSIOperator{
				Spec: platformv1.CSIOperatorSpec{
					ClusterName: cls.Name,
				},
			}
			csiOperator, err := testTKE.TkeClient.PlatformV1().CSIOperators().Create(context.Background(), csiOperator, metav1.CreateOptions{})
			Expect(err).Should(BeNil())

			Eventually(func() error {
				addon, err := testTKE.TkeClient.PlatformV1().CSIOperators().Get(context.Background(), csiOperator.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if addon.Status.Phase != "Running" {
					return errors.New(addon.Name + " Phase: " + string(addon.Status.Phase) + ", Reason: " + addon.Status.Reason)
				}
				return nil
			}, 10*time.Minute, 10*time.Second).Should(Succeed())
			Expect(testTKE.TkeClient.PlatformV1().CSIOperators().Delete(context.Background(), csiOperator.Name, metav1.DeleteOptions{})).Should(BeNil(), "Delete CSIOperator failed")
		})
	})
})
