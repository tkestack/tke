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
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"os"
	"os/exec"
	"time"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/test/e2e/tke"

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

var (
	t                 = tke.TKE{Namespace: "platform"}
	tkeClient         *tkeclientset.Clientset
	client            *kubernetes.Clientset
	provider          = tencent.NewTencentProvider()
	clusterName       string
	masterNodes       []cloudprovider.Instance
	workerNodes       []cloudprovider.Instance
	machines          []*platformv1.Machine
	tkeKubeConfigFile string
)

var _ = BeforeSuite(func() {
	t.Create()

	tkeKubeConfigFile = t.GetKubeConfigFile()
	restConf, err := t.GetKubeConfig()
	Expect(err).To(BeNil())
	tkeClient = tkeclientset.NewForConfigOrDie(restConf)
	client = kubernetes.NewForConfigOrDie(restConf)

	// Create cluster
	masterNodes, err = provider.CreateInstances(1)
	Expect(err).To(BeNil())
	time.Sleep(30 * time.Second)
	cls, err := createCluster(masterNodes)
	Expect(err).To(BeNil())
	clusterName = cls.Name
})

var _ = AfterSuite(func() {
	// Delete cluster if it exist
	if clusterName != "" {
		Expect(deleteCluster(clusterName)).Should(Succeed())
	}

	t.Delete()

	if os.Getenv("NEED_DELETE") == "" {
		return
	}
	var instanceIDs []*string
	for _, one := range masterNodes {
		instanceIDs = append(instanceIDs, &one.InstanceID)
	}
	for _, one := range workerNodes {
		instanceIDs = append(instanceIDs, &one.InstanceID)
	}
	Expect(provider.DeleteInstances(instanceIDs)).Should(Succeed())
})

var _ = Describe("Platform Test", func() {

	It("Create Baremetal cluster", func() {
		out, err := runCmd("kubectl get clusters --kubeconfig " + tkeKubeConfigFile + " | grep " + clusterName)
		Expect(err).Should(BeNil())
		Expect(out).Should(ContainSubstring("Running"))
	})

	Context("Node", func() {
		BeforeEach(func() {
			if len(machines) > 0 {
				return
			}

			var err error
			workerNodes, err = provider.CreateInstances(1)
			Expect(err).To(BeNil())
			time.Sleep(30 * time.Second)
			for _, node := range workerNodes {
				machine, err := addNode(clusterName, node)
				Expect(err).To(BeNil())
				machines = append(machines, machine)
			}
		})

		It("Add node to cluster", func() {
			for _, machine := range machines {
				out, err := runCmd("kubectl get mc --kubeconfig " + tkeKubeConfigFile + " | grep " + machine.Name)
				Expect(err).Should(BeNil())
				Expect(out).Should(ContainSubstring("Running"))
			}
		})

		It("Add tag to node", func() {
			labelKey := "testLabelKey"
			labelValue := machines[0].Spec.IP
			machines[0].Labels = map[string]string{
				labelKey: labelValue,
			}

			machine, err := tkeClient.PlatformV1().Machines().Update(context.Background(), machines[0], metav1.UpdateOptions{})
			Expect(err).Should(BeNil())

			out, err := runCmd("kubectl get mc --show-labels --kubeconfig " + tkeKubeConfigFile + " | grep " + machine.Name)
			Expect(err).Should(BeNil())
			Expect(out).Should(ContainSubstring(labelKey + "=" + labelValue))
		})

		It("Delete node", func() {
			for _, machine := range machines {
				Expect(deleteNode(machine.Name)).Should(Succeed())
			}
			for _, machine := range machines {
				out, _ := runCmd("kubectl get mc --kubeconfig " + tkeKubeConfigFile)
				Expect(out).ShouldNot(ContainSubstring(machine.Name))
			}
		})
	})

	Context("Addon", func() {
		It("TappController", func() {
			tapp := &platformv1.TappController{
				Spec: platformv1.TappControllerSpec{
					ClusterName: clusterName,
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

			out, _ := runCmd("kubectl describe tc " + tapp.Name + " -n kube-system --kubeconfig " + tkeKubeConfigFile + " | grep Phase")
			Expect(out).Should(ContainSubstring("Running"))
		})

		It("IPAM", func() {
			ipam := &platformv1.IPAM{
				Spec: platformv1.IPAMSpec{
					ClusterName: clusterName,
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

			out, _ := runCmd("kubectl describe ipam " + ipam.Name + " -n kube-system --kubeconfig " + tkeKubeConfigFile + " | grep Phase")
			Expect(out).Should(ContainSubstring("Running"))
		})

		It("CronHPA", func() {
			cronHPA := &platformv1.CronHPA{
				Spec: platformv1.CronHPASpec{
					ClusterName: clusterName,
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

			out, _ := runCmd("kubectl describe cronhpa " + cronHPA.Name + " -n kube-system --kubeconfig " + tkeKubeConfigFile + " | grep Phase")
			Expect(out).Should(ContainSubstring("Running"))
		})

		It("Prometheus", func() {
			prometheus := &platformv1.Prometheus{
				Spec: platformv1.PrometheusSpec{
					ClusterName: clusterName,
				},
			}
			prometheus, err := tkeClient.PlatformV1().Prometheuses().Create(context.Background(), prometheus, metav1.CreateOptions{})
			Expect(err).Should(BeNil())

			Eventually(func() error {
				addon, err := tkeClient.PlatformV1().Prometheuses().Get(context.Background(), prometheus.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if addon.Status.Phase != "Running" {
					return errors.New(addon.Name + " Phase: " + string(addon.Status.Phase) + ", Reason: " + addon.Status.Reason)
				}
				return nil
			}, 20*time.Minute, 10*time.Second).Should(Succeed())

			out, _ := runCmd("kubectl describe prom " + prometheus.Name + " -n kube-system --kubeconfig " + tkeKubeConfigFile + " | grep Phase")
			Expect(out).Should(ContainSubstring("Running"))
		})
	})

	It("Delete cluster", func() {
		Expect(deleteCluster(clusterName)).Should(Succeed())

		out, _ := runCmd("kubectl get clusters --kubeconfig " + tkeKubeConfigFile)
		Expect(out).ShouldNot(ContainSubstring(clusterName))
	})
})

func createCluster(masterNodes []cloudprovider.Instance) (cluster *platformv1.Cluster, err error) {
	klog.Info("Create cluster")
	cluster = &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			Type:          "Baremetal",
			Features:      platformv1.ClusterFeature{EnableMasterSchedule: true},
			Version:       "1.18.3",
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
	clusterName := cluster.Name
	klog.Info("Cluster name: ", clusterName)

	klog.Info("Wait cluster status to be running")
	err = wait.Poll(10*time.Second, 10*time.Minute, func() (bool, error) {
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
	klog.Info("Add node. InstanceId: ", workerNode.InstanceID, ", PublicIP: ", workerNode.PublicIP, ", InternalIP: ", workerNode.InternalIP)
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

func nodeSSH(ins cloudprovider.Instance) *ssh.SSH {
	s, err := ssh.New(&ssh.Config{
		User:     ins.Username,
		Password: ins.Password,
		Host:     ins.PublicIP,
		Port:     int(ins.Port),
	})
	Expect(err).To(BeNil())
	return s
}

func runCmd(cmd string) (string, error) {
	klog.Info("Run cmd: ", cmd)
	command := exec.Command("bash", "-c", cmd)
	out, err := command.CombinedOutput()
	return string(out), err
}
