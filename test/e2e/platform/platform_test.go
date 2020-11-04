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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"time"
	"tkestack.io/tke/pkg/platform/apiserver/cluster"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/test/e2e/tke"
	"tkestack.io/tke/test/util"
	"tkestack.io/tke/test/util/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/test/util/cloudprovider"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
)

const namespacePrefix = "platform-"

var (
	t                 = tke.TKE{Namespace: namespacePrefix + util.RandomStr(6)}
	tkeClient         *tkeclientset.Clientset
	provider          = tencent.NewTencentProvider()
	err               error
	tkeKubeConfigFile string
	clusterNames      []string
)

var _ = BeforeSuite(func() {
	t.Create()

	tkeKubeConfigFile = t.GetKubeConfigFile()
	restConf, err := t.GetKubeConfig()
	Expect(err).To(BeNil())
	tkeClient = tkeclientset.NewForConfigOrDie(restConf)
})

var _ = AfterSuite(func() {
	for _, name := range clusterNames {
		Expect(deleteCluster(name)).Should(Succeed())
	}

	t.Delete()

	if env.NeedDelete() == "" {
		return
	}
	Expect(provider.DeleteAllInstances()).Should(Succeed())
})

var _ = Describe("Platform Test", func() {

	Context("Baremetal cluster", func() {
		var cls *platformv1.Cluster

		BeforeEach(func() {
			// Create Baremetal cluster
			if cls == nil {
				cls, err = createCluster(1)
				Expect(err).To(BeNil())
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
					workerNodes, err := provider.CreateInstances(1)
					Expect(err).To(BeNil())
					workerNode = workerNodes[0]
					time.Sleep(10 * time.Second)
					machine, err = addNode(cls.Name, workerNode)
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
					_, err = tkeClient.PlatformV1().Machines().Update(context.Background(), machine, metav1.UpdateOptions{})
					if err != nil {
						// Get the latest machine object again if update failed
						machine, _ = tkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
					}
					return err
				}, 5*time.Second, time.Second).Should(Succeed())
				machine, _ = tkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
				Expect(machine.Labels).Should(HaveKeyWithValue(labelKey, labelValue))
			})

			It("Unschedulable node", func() {
				workerNodeIP := workerNode.InternalIP
				clusterWrapper, err := typesv1.GetCluster(context.Background(), tkeClient.PlatformV1(), cls)
				Expect(err).Should(BeNil())

				client, _ := clusterWrapper.Clientset()

				// Unschedule node
				node := updateNode(client, workerNodeIP, true)
				Expect(node.Spec.Unschedulable).Should(BeTrue())

				// Cancel unschedule node
				node = updateNode(client, workerNodeIP, false)
				Expect(node.Spec.Unschedulable).Should(BeFalse())
			})

			It("Drain node", func() {
				workerNodeIP := workerNode.InternalIP
				clusterWrapper, err := typesv1.GetCluster(context.Background(), tkeClient.PlatformV1(), cls)
				Expect(err).Should(BeNil())

				client, _ := clusterWrapper.Clientset()

				node, err := client.CoreV1().Nodes().Get(context.Background(), workerNodeIP, metav1.GetOptions{})
				Expect(err).Should(BeNil())

				// Drain node
				Expect(cluster.DrainNode(context.Background(), client, node)).Should(Succeed())

				node, _ = client.CoreV1().Nodes().Get(context.Background(), workerNodeIP, metav1.GetOptions{})
				Expect(node.Spec.Unschedulable).Should(BeTrue())
			})

			It("Delete node", func() {
				Expect(deleteNode(machine.Name)).Should(Succeed())
			})
		})

		Context("Addon", func() {
			It("TappController", func() {
				tapp := &platformv1.TappController{
					Spec: platformv1.TappControllerSpec{
						ClusterName: cls.Name,
					},
				}
				tapp, err := tkeClient.PlatformV1().TappControllers().Create(context.Background(), tapp, metav1.CreateOptions{})
				Expect(err).Should(BeNil())

				Eventually(func() error {
					addon, err := tkeClient.PlatformV1().TappControllers().Get(context.Background(), tapp.Name, metav1.GetOptions{})
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
				ipam, err := tkeClient.PlatformV1().IPAMs().Create(context.Background(), ipam, metav1.CreateOptions{})
				Expect(err).Should(BeNil())

				Eventually(func() error {
					addon, err := tkeClient.PlatformV1().IPAMs().Get(context.Background(), ipam.Name, metav1.GetOptions{})
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
				cronHPA, err := tkeClient.PlatformV1().CronHPAs().Create(context.Background(), cronHPA, metav1.CreateOptions{})
				Expect(err).Should(BeNil())

				Eventually(func() error {
					addon, err := tkeClient.PlatformV1().CronHPAs().Get(context.Background(), cronHPA.Name, metav1.GetOptions{})
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

		It("Delete Baremetal cluster", func() {
			Expect(deleteCluster(cls.Name)).Should(Succeed())
		})
	})

	Context("Imported cluster", func() {
		var cls *platformv1.Cluster
		var credential *platformv1.ClusterCredential

		BeforeEach(func() {
			if cls == nil {
				// Prepare a cluster to be imported
				cls, err = createCluster(1)
				Expect(err).To(BeNil())

				// Get the credential of cluster
				credential, err = tkeClient.PlatformV1().ClusterCredentials().Get(context.Background(), cls.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
				Expect(err).Should(BeNil())

				// Delete cluster from the global cluster to be imported
				Expect(deleteCluster(cls.Name)).Should(Succeed())
			}
		})

		It("Import cluster", func() {
			// Import cluster
			importedCluster, err := importCluster(cls.Spec.Machines[0].IP, 6443, credential.CACert, credential.Token)
			Expect(err).Should(BeNil())
			Expect(importedCluster.Name).ShouldNot(Equal(cls.Name))
			Expect(importedCluster.Status.Phase).Should(Equal(platformv1.ClusterRunning))
			Expect(importedCluster.Spec.Type).Should(Equal("Imported"))
		})
	})
})

func createCluster(masterNodeNum int64) (cluster *platformv1.Cluster, err error) {
	masterNodes, err := provider.CreateInstances(masterNodeNum)
	if err != nil {
		return nil, err
	}
	time.Sleep(10 * time.Second)
	return createClusterWithMasterNodes(masterNodes)
}

func createClusterWithMasterNodes(masterNodes []cloudprovider.Instance) (cluster *platformv1.Cluster, err error) {
	klog.Info("Create cluster")
	cluster = &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			Type:          "Baremetal",
			Features:      platformv1.ClusterFeature{EnableMasterSchedule: true},
			Version:       env.K8sVersion(),
			ClusterCIDR:   "10.244.0.0/16",
			NetworkDevice: "eth0",
		}}
	for _, one := range masterNodes {
		cluster.Spec.Machines = append(cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}
	cluster, err = tkeClient.PlatformV1().Clusters().Create(context.Background(), cluster, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info("Cluster name: ", cluster.Name)
	clusterNames = append(clusterNames, cluster.Name)
	return waitClusterToBeRunning(cluster.Name)
}

func importCluster(host string, port int32, caCert []byte, token *string) (cluster *platformv1.Cluster, err error) {
	credential := &platformv1.ClusterCredential{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "clustercredential",
		},
		CACert: caCert,
		Token:  token,
	}
	credential, err = tkeClient.PlatformV1().ClusterCredentials().Create(context.Background(), credential, metav1.CreateOptions{})
	if err != nil {
		return
	}

	cluster = &platformv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cls",
		},
		Spec: platformv1.ClusterSpec{
			//DisplayName: baremetalClusterName,
			Type: "Imported",
			ClusterCredentialRef: &corev1.LocalObjectReference{
				Name: credential.Name,
			},
		},
		Status: platformv1.ClusterStatus{
			Addresses: []platformv1.ClusterAddress{
				{
					Host: host,
					Path: "",
					Port: port,
					Type: platformv1.AddressAdvertise,
				},
			},
		},
	}
	cluster, err = tkeClient.PlatformV1().Clusters().Create(context.Background(), cluster, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info("Cluster name: ", cluster.Name)
	clusterNames = append(clusterNames, cluster.Name)
	return waitClusterToBeRunning(cluster.Name)
}

func waitClusterToBeRunning(clusterName string) (cluster *platformv1.Cluster, err error) {
	klog.Info("Wait cluster status to be running")
	err = wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		cluster, err = tkeClient.PlatformV1().Clusters().Get(context.Background(), clusterName, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			return false, nil
		}
		if len(cluster.Status.Conditions) > 0 {
			lastCondition := cluster.Status.Conditions[len(cluster.Status.Conditions)-1]
			klog.Info("Phase: ", cluster.Status.Phase, ", Type: ", lastCondition.Type, ", message: ", lastCondition.Message)
		}
		return cluster.Status.Phase == platformv1.ClusterRunning, nil
	})
	return
}

func addNode(clusterName string, workerNode cloudprovider.Instance) (machine *platformv1.Machine, err error) {
	klog.Info("Add node. InstanceId: ", workerNode.InstanceID, ", InternalIP: ", workerNode.InternalIP)
	machine = &platformv1.Machine{
		Spec: platformv1.MachineSpec{
			ClusterName: clusterName,
			Type:        "Baremetal",
			IP:          workerNode.InternalIP,
			Port:        workerNode.Port,
			Username:    workerNode.Username,
			Password:    []byte(workerNode.Password),
		},
	}
	machine, err = tkeClient.PlatformV1().Machines().Create(context.Background(), machine, metav1.CreateOptions{})
	if err != nil {
		klog.Error(err)
		return
	}

	klog.Info("Wait node status to be running")
	err = wait.Poll(5*time.Second, 5*time.Minute, func() (bool, error) {
		machine, err = tkeClient.PlatformV1().Machines().Get(context.Background(), machine.Name, metav1.GetOptions{})
		if err != nil {
			klog.Error(err)
			return false, nil
		}
		if len(machine.Status.Conditions) > 0 {
			lastCondition := machine.Status.Conditions[len(machine.Status.Conditions)-1]
			klog.Info("Phase: ", machine.Status.Phase, ", Type: ", lastCondition.Type, ", message: ", lastCondition.Message)
		}
		return machine.Status.Phase == platformv1.MachineRunning, nil
	})
	return
}

func deleteNode(machineName string) (err error) {
	klog.Info("Delete node: ", machineName)
	err = tkeClient.PlatformV1().Machines().Delete(context.Background(), machineName, metav1.DeleteOptions{})
	if err != nil {
		return
	}

	klog.Info("Wait node to be deleted")
	return wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err = tkeClient.PlatformV1().Machines().Get(context.Background(), machineName, metav1.GetOptions{})
		if k8serror.IsNotFound(err) {
			klog.Info("Node was deleted")
			return true, nil
		}
		return false, nil
	})
}

func updateNode(client kubernetes.Interface, nodeName string, unschedulable bool) *corev1.Node {
	Eventually(func() error {
		node, err := client.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		node.Spec.Unschedulable = unschedulable
		_, err = client.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
		return err
	}, 5*time.Second, time.Second).Should(Succeed())

	node, _ := client.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	return node
}

func deleteCluster(clusterName string) (err error) {
	klog.Info("Delete cluster: ", clusterName)
	err = tkeClient.PlatformV1().Clusters().Delete(context.Background(), clusterName, metav1.DeleteOptions{})
	if k8serror.IsNotFound(err) {
		klog.Info("Cluster was not found")
		return nil
	}
	if err != nil {
		klog.Error(err)
		return err
	}
	klog.Info("Wait cluster to be deleted")
	return wait.Poll(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err := tkeClient.PlatformV1().Clusters().Get(context.Background(), clusterName, metav1.GetOptions{})
		if k8serror.IsNotFound(err) {
			klog.Info("Cluster was deleted")
			return true, nil
		}
		return false, nil
	})
}
