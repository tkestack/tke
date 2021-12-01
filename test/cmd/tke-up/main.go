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

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/test/util/cloudprovider/tencent"

	"k8s.io/apimachinery/pkg/util/wait"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	"tkestack.io/tke/test/util/cloudprovider"

	// import platform schema
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalmachine "tkestack.io/tke/pkg/platform/provider/baremetal/machine"
	importedcluster "tkestack.io/tke/pkg/platform/provider/imported/cluster"
)

func quickConfig(nodes []cloudprovider.Instance) []byte {
	para := new(types.CreateClusterPara)
	for _, one := range nodes[1:] {
		para.Cluster.Spec.Machines = append(para.Cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}
	para.Config.Auth.TKEAuth = &types.TKEAuth{}
	para.Config.Registry.ThirdPartyRegistry = &types.ThirdPartyRegistry{
		Domain:    "docker.io",
		Namespace: "tkestack",
		Username:  os.Getenv("REGISTRY_USERNAME"),
		Password:  []byte(os.Getenv("REGISTRY_PASSWORD")),
	}
	para.Config.Business = &types.Business{}
	para.Config.Monitor = &types.Monitor{
		InfluxDBMonitor: &types.InfluxDBMonitor{
			LocalInfluxDBMonitor: &types.LocalInfluxDBMonitor{},
		},
	}
	para.Config.Gateway = &types.Gateway{}

	data, _ := json.Marshal(para)

	return data
}

func main() {
	baremetalcluster.RegisterProvider()
	baremetalmachine.RegisterProvider()
	importedcluster.RegisterProvider()
	provider := tencent.NewTencentProvider()
	nodes, err := provider.CreateInstances(3)
	if err != nil {
		log.Fatal(err)
	}

	nodesSSH := make([]ssh.Interface, len(nodes))
	for i, one := range nodes {
		fmt.Printf("ensure ssh %d %s is ready\n", i, one.InternalIP)
		s, err := ssh.New(&ssh.Config{
			User:     one.Username,
			Password: one.Password,
			Host:     one.InternalIP,
			Port:     int(one.Port),
		})
		if err != nil {
			log.Fatal(err)
		}
		for j := 1; j <= 10; j++ {
			err = s.Ping()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		nodesSSH[i] = s
	}

	version := os.Getenv("VERSION")
	cmd := fmt.Sprintf("wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-x86_64-%s.run{,.sha256} && sha256sum --check --status tke-installer-x86_64-%s.run.sha256 && chmod +x tke-installer-x86_64-%s.run && ./tke-installer-x86_64-%s.run",
		version, version, version, version)
	_, err = nodesSSH[0].CombinedOutput(cmd)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s:8080/api/cluster", nodes[0].InternalIP)
	body := quickConfig(nodes)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("response code %d\n%s", resp.StatusCode, string(respBody))
	}

	err = wait.PollImmediate(5*time.Second, 30*time.Minute, func() (bool, error) {
		url := fmt.Sprintf("http://%s:8080/api/cluster/global/progress", nodes[0].InternalIP)
		resp, err := http.Get(url)
		if err != nil {
			return false, nil
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		progress := new(types.ClusterProgress)
		json.Unmarshal(data, progress)
		switch progress.Status {
		case types.StatusSuccess:
			fmt.Println(base64.StdEncoding.EncodeToString(progress.Kubeconfig))
			return true, nil

		case types.StatusUnknown, types.StatusDoing:
			return false, nil
		case types.StatusFailed:
			return false, fmt.Errorf("install failed:\n%s", progress.Data)
		default:
			return false, fmt.Errorf("unknown install progress status: %s", progress.Status)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
