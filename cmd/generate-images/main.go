/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"sort"
	"strings"

	"tkestack.io/tke/pkg/spec"

	"github.com/thoas/go-funk"

	installer "tkestack.io/tke/cmd/tke-installer/app/installer/images"

	baremetal "tkestack.io/tke/pkg/platform/provider/baremetal/images"

	galaxy "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"

	cronhpa "tkestack.io/tke/pkg/platform/controller/addon/cronhpa/images"
	helm "tkestack.io/tke/pkg/platform/controller/addon/helm/images"
	ipam "tkestack.io/tke/pkg/platform/controller/addon/ipam/images"
	lbcf "tkestack.io/tke/pkg/platform/controller/addon/lbcf/images"
	logcollector "tkestack.io/tke/pkg/platform/controller/addon/logcollector/images"
	persistentevent "tkestack.io/tke/pkg/platform/controller/addon/persistentevent/images"
	prometheus "tkestack.io/tke/pkg/platform/controller/addon/prometheus/images"
	csioperator "tkestack.io/tke/pkg/platform/controller/addon/storage/csioperator/images"
	volumedecorator "tkestack.io/tke/pkg/platform/controller/addon/storage/volumedecorator/images"
	tappcontroller "tkestack.io/tke/pkg/platform/controller/addon/tappcontroller/images"
)

func main() {
	funcs := []func() []string{
		//installer.List,

		//baremetal.List,

		//galaxy.List,

		cronhpa.List,
		helm.List,
		ipam.List,
		lbcf.List,
		logcollector.List,
		persistentevent.List,
		prometheus.List,
		csioperator.List,
		volumedecorator.List,
		tappcontroller.List,
	}
	var result []string
	for _, f := range funcs {
		images := f()
		result = append(result, images...)
	}
	result = funk.UniqString(result)
	for _, one := range append(baremetal.List(), append(installer.List(), galaxy.List()...)...) {
		if strings.HasPrefix(one, "nvidia-device-plugin") {
			fmt.Println(one)
			continue
		}
		for _, arch := range spec.Archs {
			fmt.Println(strings.ReplaceAll(one, ":", "-"+arch+":"))
		}
	}
	sort.Strings(result)
	for _, one := range result {
		fmt.Println(one)
	}
}
