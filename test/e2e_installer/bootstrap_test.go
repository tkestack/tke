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

package e2e_installer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/test/util/cloudprovider"
	"tkestack.io/tke/test/util/cloudprovider/tencent"
)

var _ = Describe("bootstrap", func() {
	var nodes []cloudprovider.Instance
	var nodesSSH []ssh.Interface
	var provider cloudprovider.Provider
	needDelete := false

	BeforeEach(func() {
		var err error
		provider = tencent.NewTencentProvider()
		nodes, err = provider.CreateInstances(3)
		Expect(err).To(BeNil())

		nodesSSH = make([]ssh.Interface, len(nodes))
		for i, one := range nodes {
			fmt.Printf("ensure ssh %d %s is ready\n", i, one.InternalIP)
			s, err := ssh.New(&ssh.Config{
				User:     one.Username,
				Password: one.Password,
				Host:     one.InternalIP,
				Port:     int(one.Port),
			})
			Expect(err).To(BeNil())
			for j := 1; j <= 10; j++ {
				err = s.Ping()
				if err == nil {
					break
				}
				time.Sleep(5 * time.Second)
			}
			nodesSSH[i] = s
		}
	})

	AfterEach(func() {
		if !needDelete {
			return
		}
		var instanceIDs []*string
		for i, one := range nodes {
			fmt.Printf("delete instance %d %s\n", i, one.InternalIP)
			instanceIDs = append(instanceIDs, &nodes[i].InstanceID)
		}
		err := provider.DeleteInstances(instanceIDs)
		Expect(err).To(BeNil())
	})

	It("should bootstrap successfuly", func() {
		By("install installer")
		version := os.Getenv("VERSION")
		cmd := fmt.Sprintf("wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-x86_64-%s.run{,.sha256} && sha256sum --check --status tke-installer-x86_64-%s.run.sha256 && chmod +x tke-installer-x86_64-%s.run && ./tke-installer-x86_64-%s.run",
			version, version, version, version)
		_, err := nodesSSH[0].CombinedOutput(cmd)
		Expect(err).To(BeNil())

		By("prepare parametes")
		url := fmt.Sprintf("http://%s:8080/api/cluster", nodes[0].InternalIP)
		para := new(installer.CreateClusterPara)
		for _, one := range nodes[1:] {
			para.Cluster.Spec.Machines = append(para.Cluster.Spec.Machines, platformv1.ClusterMachine{
				IP:       one.InternalIP,
				Port:     one.Port,
				Username: one.Username,
				Password: []byte(one.Password),
			})
		}
		para.Config.Auth.TKEAuth = &installer.TKEAuth{}
		para.Config.Registry.TKERegistry = &installer.TKERegistry{Domain: "registry.tke.com"}
		para.Config.Business = &installer.Business{}
		para.Config.Monitor = &installer.Monitor{
			InfluxDBMonitor: &installer.InfluxDBMonitor{
				LocalInfluxDBMonitor: &installer.LocalInfluxDBMonitor{},
			},
		}
		para.Config.Gateway = &installer.Gateway{}
		body, err := json.Marshal(para)
		Expect(err).To(BeNil())

		By("post data to installer for install")
		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		Expect(err).To(BeNil())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		By("wait install finish")
		err = wait.PollImmediate(5*time.Second, 30*time.Minute, func() (bool, error) {
			url := fmt.Sprintf("http://%s:8080/api/cluster/global/progress", nodes[0].InternalIP)
			resp, err := http.Get(url)
			if err != nil {
				return false, nil
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			progress := new(installer.ClusterProgress)
			err = json.Unmarshal(data, progress)
			Expect(err).To(BeNil())
			switch progress.Status {
			case installer.StatusUnknown, installer.StatusDoing:
				return false, nil
			case installer.StatusFailed:
				return false, fmt.Errorf("install failed:\n%s", progress.Data)
			case installer.StatusSuccess:
				return true, nil
			default:
				return false, fmt.Errorf("unknown install progress status: %s", progress.Status)
			}
		})
		Expect(err).To(BeNil())
		needDelete = true
	})
})
