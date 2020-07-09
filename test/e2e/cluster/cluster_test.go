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

package cluster_test

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/thoas/go-funk"
	"tkestack.io/tke/test/e2e/cluster"
	"tkestack.io/tke/test/e2e/cluster/certs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	platformv1 "tkestack.io/tke/api/platform/v1"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/util/apiclient"
	testclient "tkestack.io/tke/test/util/client"
	"tkestack.io/tke/test/util/cloudprovider"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
)

var tkenamespace = "tke"

var _ = BeforeSuite(func() {

	certs.InitTmpDir()
	str, _ := os.Getwd()
	fmt.Println("Before suite:", str)

	By("start create")

	client := testclient.GetClientSet()
	tkeHostName := tkeHostName()
	fmt.Println("tkeHostName:", tkeHostName)

	By("enable master shedchule")
	err := enableMasterSchedule(context.Background(), client)
	Expect(err).To(BeNil())

	By("create namespace")
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: tkenamespace,
		},
	}

	err = apiclient.CreateOrUpdateNamespace(context.Background(), client, ns)
	Expect(err).To(BeNil())
	fmt.Println("namespace:", tkenamespace)

	By("create certmap")
	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP(tkeHostName)}
	err = certs.CreateCertMap(context.Background(), client, ips, tkenamespace)
	Expect(err).To(BeNil())

	By("create etcd")
	options := map[string]interface{}{
		"Servers":   []string{tkeHostName},
		"Namespace": tkenamespace,
	}
	err = apiclient.CreateResourceWithDir(context.Background(), client, "manifests/etcd/*.yaml", options)
	Expect(err).To(BeNil())
	err = wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), client, tkenamespace, "app=etcd")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
	Expect(err).To(BeNil())

	By("create provider configs")
	err = createProviderConfigs(context.Background(), client, tkenamespace)
	Expect(err).To(BeNil())

	By("create tke-platform-api")
	options = map[string]interface{}{
		"Image":            "tkestack/tke-platform-api-amd64:v1.2.4.155.g166fdc8.dirty",
		"ProviderResImage": "tkestack/provider-res-amd64:v1.18.3-1",
		"EnableAuth":       false,
		"EnableAudit":      false,
		"Namespace":        tkenamespace,
	}

	fmt.Println("create tke-platform-api")
	err = apiclient.CreateResourceWithDir(context.Background(), client, "manifests/tke-platform-api/*.yaml", options)
	fmt.Println("create tke-platform-api:", err)
	Expect(err).To(BeNil())

	err = wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), client, tkenamespace, "app=tke-platform-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
	Expect(err).To(BeNil())

	By("create tke-platform-controller")

	options = map[string]interface{}{
		"Image":             "tkestack/tke-platform-controller-amd64:v1.2.4.155.g166fdc8.dirty",
		"ProviderResImage":  "tkestack/provider-res-amd64:v1.18.3-1",
		"RegistryDomain":    "docker.io",
		"RegistryNamespace": "tkestack",
		"Namespace":         tkenamespace,
	}

	err = apiclient.CreateResourceWithDir(context.Background(), client, "manifests/tke-platform-controller/*.yaml", options)
	Expect(err).To(BeNil())

	err = wait.PollImmediate(5*time.Second, 2*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckPodReadyWithLabel(context.Background(), client, tkenamespace,
			"app=tke-platform-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
	Expect(err).To(BeNil())

	By("create kubeconfig")
	svc, err := client.CoreV1().Services(tkenamespace).Get(context.Background(), "tke-platform-api", metav1.GetOptions{})
	Expect(err).To(BeNil())
	nodePort := int(svc.Spec.Ports[0].NodePort)
	err = certs.WriteKubeConfig(tkeHostName, nodePort, tkenamespace)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	fmt.Println("After suite")
	certs.ClearTmpDir()
	gracePeriodSeconds := int64(0)
	deleteOptions := metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}
	client := testclient.GetClientSet()
	err := client.CoreV1().Namespaces().Delete(context.Background(), tkenamespace, deleteOptions)
	Expect(err).To(BeNil())
})

func tkeHostName() string {
	restconf := testclient.GetRESTConfig()
	host := restconf.Host
	url, _ := url.Parse(host)
	return url.Hostname()
}

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
		err := createBaremetalCluster(context.Background(), masterNodes, workerNodes)
		Expect(err).To(BeNil())
	})
})

func enableMasterSchedule(ctx context.Context, client kubernetes.Interface) error {
	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodeList.Items {
		taint := corev1.Taint{
			Key:    baremetalconstants.LabelNodeRoleMaster,
			Effect: corev1.TaintEffectNoSchedule,
		}
		if !funk.Contains(node.Spec.Taints, taint) {
			break
		}
		err = apiclient.RemoveNodeTaints(context.Background(), client, node.Name, []corev1.Taint{taint})
		if err != nil {
			return err
		}
	}

	return nil
}

func createBaremetalCluster(ctx context.Context, masterNodes []cloudprovider.Instance,
	workerNodes []cloudprovider.Instance) error {
	time.Sleep(70 * time.Second)
	cluster := &platformv1.Cluster{
		Spec: platformv1.ClusterSpec{
			TenantID:    "default",
			Version:     "1.18.3",
			ClusterCIDR: "10.244.0.0/16",
			Type:        "Baremetal",
			Features:    platformv1.ClusterFeature{EnableMasterSchedule: true},
		}}
	for _, one := range masterNodes {
		cluster.Spec.Machines = append(cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}

	restconf, err := certs.GetKubeConfig()
	Expect(err).To(BeNil())
	client := tkeclientset.NewForConfigOrDie(restconf)

	By("create cluster")
	cluster, err = client.PlatformV1().Clusters().Create(context.Background(), cluster, metav1.CreateOptions{})
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

func createProviderConfigs(ctx context.Context, client kubernetes.Interface, namespace string) error {
	configMaps := []struct {
		Name string
		File string
	}{
		{
			Name: "provider-config",
			File: cluster.ConfDir + "config.yaml",
		},
		{
			Name: "docker",
			File: cluster.ConfDir + "docker/*",
		},
		{
			Name: "kubelet",
			File: cluster.ConfDir + "kubelet/*",
		},
		{
			Name: "kubeadm",
			File: cluster.ConfDir + "kubeadm/*",
		},
	}

	for _, one := range configMaps {
		err := apiclient.CreateOrUpdateConfigMapFromFile(ctx, client,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      one.Name,
					Namespace: namespace,
				},
			},
			one.File)
		if err != nil {
			return err
		}
	}

	return nil
}
