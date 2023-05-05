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
	"os"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
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

var _ = Describe("cluster", func() {

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
		return []byte(kubeconfig)
	}, func(data []byte) {
		tkeClient, err := testclient.GetTKEClient(data)
		Expect(err).Should(BeNil(), "Get tke client with admin token failed")
		testTKE = tke2.Init(tkeClient, provider)
	})

	SynchronizedAfterSuite(func() {
		if env.NeedDelete() {
			Expect(provider.TearDown()).Should(BeNil())
		}
	}, func() {})

	AfterEach(func() {
		if cls != nil {
			testTKE.DeleteCluster(cls.Name)
		}
	})

	It("Create and Delete Baremetal cluster", func() {
		cls, err = testTKE.CreateCluster()
		Expect(err).To(BeNil(), "Create cluster failed")
		Expect(testTKE.DeleteCluster(cls.Name)).Should(BeNil(), "Delete cluster failed")
	})

	It("Import cluster", func() {
		// Prepare a cluster to be imported
		oldCls, err := testTKE.CreateCluster()
		Expect(err).To(BeNil(), "Create cluster failed")

		// Get the credential of cluster
		credential, err := testTKE.TkeClient.PlatformV1().ClusterCredentials().Get(context.Background(), oldCls.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		Expect(err).Should(BeNil(), "Get ClusterCredential failed")

		// Delete cluster from the global cluster in order to import it
		Expect(testTKE.DeleteCluster(oldCls.Name)).Should(BeNil(), "Delete cluster failed")

		// Import cluster
		cls, err = testTKE.ImportCluster(oldCls.Spec.Machines[0].IP, 6443, credential.CACert, credential.Token)
		Expect(err).Should(BeNil(), "Import cluster failed")
		Expect(cls.Name).ShouldNot(Equal(oldCls.Name), "Imported cluster name was the same with the original cluster name")
		Expect(cls.Spec.Type).Should(Equal("Imported"), "Cluster type was not 'Imported'")
	})

	/*
		DescribeTable("Upgrade cluster",
			func(oldVersion, newVersion string) {
				cls = testTKE.ClusterTemplate()
				cls.Spec.Version = oldVersion
				cls, err = testTKE.CreateClusterInternal(cls)
				Expect(err).To(BeNil(), "Create cluster failed")

				cls, err = testTKE.UpgradeCluster(cls.Name, newVersion, platformv1.UpgradeModeAuto, false)
				Expect(err).Should(BeNil(), "Upgrade cluster failed")
				Expect(cls.Spec.Version).Should(Equal(newVersion), "Cluster version is wrong")
			},
			// Entry("1.19.7->1.20.4", "1.19.7", "1.20.4"),
			Entry("1.20.6-tke.2->1.21.4-tke.3", "1.20.6-tke.2", "1.21.4-tke.3"))
	*/

	It("Cluster scaling", func() {
		// Prepare two instances
		nodes, err := provider.CreateInstances(2)
		Expect(err).Should(BeNil(), "Create instances failed")

		// Prepare a VIP by creating a CLB with the two nodes
		var ips []string
		for _, node := range nodes {
			ips = append(ips, node.InternalIP)
		}
		vip, err := provider.CreateCLB(common.StringPtrs(ips))
		Expect(err).Should(BeNil(), "Create LB failed")

		// Create cluster with the first prepared node
		cls = testTKE.ClusterTemplate(nodes[0])
		cls.Spec.Features.HA.ThirdPartyHA = &platformv1.ThirdPartyHA{
			VIP:   *vip,
			VPort: 6443,
		}
		cls, err = testTKE.CreateClusterInternal(cls)
		Expect(err).To(BeNil(), "Create cluster failed")

		// upscale
		cls, err = testTKE.ScaleUp(cls.Name, nodes[1:])
		Expect(err).Should(BeNil(), "Cluster upscale failed")
		Expect(cls.Spec.Machines).Should(HaveLen(2), "Cluster node num is wrong")

		// downscale
		cls, err = testTKE.ScaleDown(cls.Name, ips[1:])
		Expect(err).Should(BeNil(), "Cluster downscale failed")
		Expect(cls.Spec.Machines).Should(HaveLen(1), "Cluster node num is wrong")
	})

	/*
		It("Cluster bootstrap application", func() {
			nodes, err := provider.CreateInstances(1)
			Expect(err).Should(BeNil(), "Create instances failed")

			cls = testTKE.ClusterTemplate(nodes[0])
			cls.Spec.BootstrapApps = []platformv1.BootstrapApp{
				{
					App: platformv1.App{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "kube-system",
						},
						Spec: v1.AppSpec{
							Type:            "HelmV3",
							TenantID:        "default",
							Name:            "demo1",
							TargetCluster:   "",
							TargetNamespace: "",
							Chart: v1.Chart{
								ChartName:      "tke-resilience",
								ChartGroupName: "public",
								ChartVersion:   "1.0.0",
								TenantID:       "default",
							},
							Values: v1.AppValues{
								RawValues: "key2: val2-override",
							},
						},
					},
				},
				{
					App: platformv1.App{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "kube-public",
						},
						Spec: v1.AppSpec{
							Name:            "demo2",
							Type:            "HelmV3",
							TenantID:        "default",
							TargetCluster:   "",
							TargetNamespace: "kube-public",
							Chart: v1.Chart{
								ChartName:      "tke-resilience",
								ChartGroupName: "public",
								ChartVersion:   "1.0.0",
								TenantID:       "default",
							},
							Values: v1.AppValues{
								RawValues: "key2: val2-override",
							},
						},
					},
				},
			}
			cls, err = testTKE.CreateClusterInternal(cls)
			Expect(err).To(BeNil(), "Create cluster failed")

			By("验证bootstrap app已创建")
			verifyApp := func(namespace, appName string) error {
				apps, err := testTKE.TkeClient.ApplicationV1().Apps(namespace).List(context.Background(), metav1.ListOptions{})
				if err != nil {
					return fmt.Errorf("list apps in namespace %v failed", namespace)
				}
				klog.Infof("Apps in %v: %v", namespace, len(apps.Items))
				for _, app := range apps.Items {
					klog.Info(app.Name)
				}
				appFullName := fmt.Sprintf("bootstrapapp-%v-%v", namespace, appName)
				_, err = testTKE.TkeClient.ApplicationV1().Apps(namespace).Get(context.Background(), appFullName, metav1.GetOptions{})
				if err != nil {
					err = fmt.Errorf("get app %v failed", appFullName)
				}
				return err
			}
			Eventually(func() error {
				err = verifyApp("kube-system", "demo1")
				if err != nil {
					return err
				}
				return verifyApp("kube-public", "demo2")
			}, time.Minute, time.Second).Should(BeNil())
		})
	*/
})
