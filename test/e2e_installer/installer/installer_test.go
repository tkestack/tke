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
	. "github.com/onsi/gomega"
	"os"
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
		if !env.NeedDelete() {
			return
		}
		Expect(provider.TearDown()).Should(BeNil(), "")
	})

	It("最小化安装", func() {
		installInstaller()

		nodes, err := provider.CreateInstances(1)
		Expect(err).Should(BeNil(), "Create instance failed")

		para := installer.CreateClusterParaTemplate(nodes)
		err = installer.Install(para)
		Expect(err).To(BeNil(), "Install failed")
	})

	It("默认安装", func() {
		installInstaller()

		nodes, err := provider.CreateInstances(1)
		Expect(err).Should(BeNil(), "Create instance failed")

		para := installer.CreateClusterParaTemplate(nodes)
		// TODO: add more para
		err = installer.Install(para)
		Expect(err).To(BeNil(), "Install failed")
	})
})

func installInstaller() {
	// Download and install tke-installer
	installer = tke2.InitInstaller(provider)
	err := installer.InstallInstaller(os.Getenv("OS"), os.Getenv("ARCH"), os.Getenv("VERSION"))
	Expect(err).Should(BeNil(), "Install tke-installer failed")
}
