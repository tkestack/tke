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

package cluster

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	v1 "tkestack.io/tke/api/platform/v1"

	"github.com/stretchr/testify/assert"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

const dataFile = "testdata/cluster.json"

func init() {
	_ = os.Chdir("..")
}

func TestCluster_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skip testing in short mode")
	}

	p := new(Provider)
	err := p.Init("conf/config.yaml")
	assert.Nil(t, err)

	var cluster clusterprovider.Cluster
	data, err := ioutil.ReadFile(dataFile)
	assert.Nil(t, err)
	_ = json.Unmarshal(data, &cluster)
	// cluster.Status = v1.ClusterStatus{} // reset for force create

	for {
		cluster, err = p.OnInitialize(cluster)
		data, _ := json.MarshalIndent(cluster, "", " ")
		_ = ioutil.WriteFile(dataFile, data, 0777)
		if err != nil {
			t.Fatal(err)
			return
		}
		if cluster.Cluster.Status.Phase == v1.ClusterRunning {
			break
		}
		time.Sleep(time.Second)
	}
}
