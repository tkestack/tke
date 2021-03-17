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

package installer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"os"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	tke2 "tkestack.io/tke/test/tke"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
	"tkestack.io/tke/test/util/env"
)

var (
	provider  = tencent.NewTencentProvider()
	installer *tke2.Installer
)

var _ = Describe("tke-installer", func() {

	AfterEach(func() {
		if env.NeedDelete() {
			Expect(provider.TearDown()).Should(BeNil(), "")
		}
	})

	DescribeTable("install",
		func(paraGenerator BuildCreateClusterPara) {
			installInstaller()
			para := paraGenerator()
			Expect(installer.Install(para)).To(BeNil(), "Install failed")
		},
		Entry("最小化安装", func() *types.CreateClusterPara {
			nodes, err := provider.CreateInstances(1)
			Expect(err).Should(BeNil(), "Create instance failed")
			return installer.CreateClusterParaTemplate(nodes)
		}),
		Entry("默认安装", func() *types.CreateClusterPara {
			nodes, err := provider.CreateInstances(1)
			Expect(err).Should(BeNil(), "Create instance failed")
			para := installer.CreateClusterParaTemplate(nodes)
			// 监控模块
			para.Config.Monitor = &types.Monitor{
				InfluxDBMonitor: &types.InfluxDBMonitor{
					LocalInfluxDBMonitor: &types.LocalInfluxDBMonitor{},
				},
			}
			// 开启mesh
			para.Config.Mesh = &types.Mesh{}
			// 集群设置: GPU类型: Virtual
			gpuType := platformv1.GPUVirtual
			para.Cluster.Spec.Features = platformv1.ClusterFeature{
				GPUType:             &gpuType,
				EnableMetricsServer: true,
			}

			return para
		}))
})

func installInstaller() {
	// Download and install tke-installer
	installer = tke2.InitInstaller(provider)
	err := installer.InstallInstaller(os.Getenv("OS"), os.Getenv("ARCH"), os.Getenv("VERSION"))
	Expect(err).Should(BeNil(), "Install tke-installer failed")
}

type BuildCreateClusterPara func() *types.CreateClusterPara
