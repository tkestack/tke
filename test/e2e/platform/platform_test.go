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
	"errors"
	"time"

	"tkestack.io/tke/pkg/platform/apiserver/cluster"
	"tkestack.io/tke/test/e2e/tke"
	tke2 "tkestack.io/tke/test/tke"
	"tkestack.io/tke/test/util"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
	"tkestack.io/tke/test/util/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1 "tkestack.io/tke/api/platform/v1"
	_ "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	_ "tkestack.io/tke/pkg/platform/provider/baremetal/machine"
	_ "tkestack.io/tke/pkg/platform/provider/imported/cluster"
	_ "tkestack.io/tke/pkg/platform/provider/registered/cluster"
	"tkestack.io/tke/test/util/cloudprovider"
)

const namespacePrefix = "platform-"

var (
	t                 = tke.TKE{Namespace: namespacePrefix + util.RandomStr(6)}
	provider          = tencent.NewTencentProvider()
	testTKE           *tke2.TestTKE
	err               error
	tkeKubeConfigFile string
)

var _ = BeforeSuite(func() {
	t.Create()

	tkeKubeConfigFile = t.GetKubeConfigFile()
	restConf, err := t.GetKubeConfig()
	Expect(err).To(BeNil())
	tkeClient := tkeclientset.NewForConfigOrDie(restConf)
	testTKE = tke2.Init(tkeClient, provider)
})

var _ = AfterSuite(func() {
	for _, cls := range testTKE.Clusters {
		testTKE.DeleteCluster(cls.Name)
	}

	t.Delete()

	if env.NeedDelete() {
		Expect(provider.TearDown()).Should(BeNil())
	}
})

var _ = Describe("Platform Test", func() {

	Context("Baremetal cluster", func() {
		var cls *platformv1.Cluster

		BeforeEach(func() {
			// Create Baremetal cluster
			if cls == nil {
				cls, err = testTKE.CreateCluster()
				Expect(err).To(BeNil(), "Create cluster failed")
			}
		})

		It("Create Baremetal cluster", func() {
			// Cluster create operation has been executed in BeforeEach, this empty 'It' indicates a independent case
		})

		Context("Node", func() {
			var workerNode cloudprovider.Instance
			var machine *platformv1.Machine

			BeforeEach(func() {
				if machine == nil {
					workerNodes, err := testTKE.CreateInstances(1)
					Expect(err).To(BeNil())
					workerNode = workerNodes[0]
					machine, err = testTKE.AddNode(cls.Name, workerNode)
					Expect(err).To(BeNil())
				}
			})

			It("Add node to cluster", func() {
				// Adding node operation has been executed in BeforeEach, this empty 'It' indicates a independent case
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
				}, 5*time.Second, time.Second).Should(BeNil())
				machine, _ = testTKE.TkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
				Expect(machine.Labels).Should(HaveKeyWithValue(labelKey, labelValue))
			})

			It("Unschedulable node", func() {
				workerNodeIP := workerNode.InternalIP

				// Unschedule node
				node, _ := testTKE.UnscheduleNode(cls, workerNodeIP)
				Expect(node.Spec.Unschedulable).Should(BeTrue())

				// Cancel unschedule node
				node, _ = testTKE.CancleUnschedulableNode(cls, workerNodeIP)
				Expect(node.Spec.Unschedulable).Should(BeFalse())
			})

			It("Drain node", func() {
				workerNodeIP := workerNode.InternalIP
				k8sClient := testTKE.K8sClient(cls)

				node, err := k8sClient.CoreV1().Nodes().Get(context.Background(), workerNodeIP, metav1.GetOptions{})
				Expect(err).Should(BeNil())

				// Drain node
				Expect(cluster.DrainNode(context.Background(), k8sClient, node)).Should(BeNil())

				node, _ = k8sClient.CoreV1().Nodes().Get(context.Background(), workerNodeIP, metav1.GetOptions{})
				Expect(node.Spec.Unschedulable).Should(BeTrue())
			})

			It("Delete node", func() {
				Expect(testTKE.DeleteNode(machine.Name)).Should(BeNil(), "Delete node failed")
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
				}, 10*time.Minute, 10*time.Second).Should(BeNil())
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
				}, 10*time.Minute, 10*time.Second).Should(BeNil())
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
				}, 10*time.Minute, 10*time.Second).Should(BeNil())
			})
		})

		It("Delete Baremetal cluster", func() {
			Expect(testTKE.DeleteCluster(cls.Name)).Should(BeNil(), "Delete cluster failed")
		})
	})

	Context("Imported cluster", func() {
		var cls *platformv1.Cluster
		var credential *platformv1.ClusterCredential

		BeforeEach(func() {
			if cls == nil {
				// Prepare a cluster to be imported
				cls, err = testTKE.CreateCluster()
				Expect(err).To(BeNil())

				// Get the credential of cluster
				credential, err = testTKE.TkeClient.PlatformV1().ClusterCredentials().Get(context.Background(), cls.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
				Expect(err).Should(BeNil())

				// Delete cluster from the global cluster to be imported
				Expect(testTKE.DeleteCluster(cls.Name)).Should(BeNil())
			}
		})

		It("Import cluster", func() {
			// Import cluster
			importedCluster, err := testTKE.ImportCluster(cls.Spec.Machines[0].IP, 6443, credential.CACert, credential.Token)
			Expect(err).Should(BeNil(), "Import cluster failed")
			Expect(importedCluster.Name).ShouldNot(Equal(cls.Name))
			Expect(importedCluster.Status.Phase).Should(Equal(platformv1.ClusterRunning))
			Expect(importedCluster.Spec.Type).Should(Equal("Imported"))
		})
	})
})
