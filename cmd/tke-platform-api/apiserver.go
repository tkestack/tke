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

package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"tkestack.io/tke/cmd/tke-platform-api/app"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalmachine "tkestack.io/tke/pkg/platform/provider/baremetal/machine"
	edgecluster "tkestack.io/tke/pkg/platform/provider/edge/cluster"
	importedcluster "tkestack.io/tke/pkg/platform/provider/imported/cluster"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	baremetalcluster.RegisterProvider()
	baremetalmachine.RegisterProvider()
	importedcluster.RegisterProvider()
	edgecluster.RegisterProvider()

	app.NewApp("tke-platform-api").Run()
}
