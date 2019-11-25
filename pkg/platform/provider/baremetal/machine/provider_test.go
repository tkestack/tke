/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package machine

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	v1 "tkestack.io/tke/api/platform/v1"

	"github.com/stretchr/testify/assert"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

const (
	clusterData = "cluster.json"
	machineData = "machine.json"
)

func TestMachine_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skip testing in short mode")
	}

	p := new(Provider)
	_ = p.Init("conf/config.yaml")

	var machine v1.Machine
	data, err := ioutil.ReadFile(machineData)
	assert.Nil(t, err)
	_ = json.Unmarshal(data, &machine)
	machine.Status = v1.MachineStatus{} // reset for force create

	var cluster clusterprovider.Cluster
	data, err = ioutil.ReadFile(clusterData)
	assert.Nil(t, err)
	_ = json.Unmarshal(data, &cluster)

	for {
		machine, err = p.OnInitialize(machine, cluster.Cluster, cluster.ClusterCredential)
		data, _ := json.MarshalIndent(machine, "", " ")
		_ = ioutil.WriteFile(machineData, data, 0777)
		if err != nil {
			t.Fatal(err)
			return
		}
		if machine.Status.Phase == v1.MachineRunning {
			break
		}
		time.Sleep(time.Second)
	}
}
