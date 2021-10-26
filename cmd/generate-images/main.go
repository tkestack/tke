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

	"github.com/spf13/pflag"

	"github.com/thoas/go-funk"
	installer "tkestack.io/tke/cmd/tke-installer/app/installer/images"
	logagent "tkestack.io/tke/pkg/logagent/controller/logagent/images"
	mesh "tkestack.io/tke/pkg/mesh/controller/meshmanager/images"
	cronhpa "tkestack.io/tke/pkg/platform/controller/addon/cronhpa/images"
	ipam "tkestack.io/tke/pkg/platform/controller/addon/ipam/images"
	lbcf "tkestack.io/tke/pkg/platform/controller/addon/lbcf/images"
	logcollector "tkestack.io/tke/pkg/platform/controller/addon/logcollector/images"
	persistentevent "tkestack.io/tke/pkg/platform/controller/addon/persistentevent/images"
	prometheus "tkestack.io/tke/pkg/platform/controller/addon/prometheus/images"
	volumedecorator "tkestack.io/tke/pkg/platform/controller/addon/storage/volumedecorator/images"
	tappcontroller "tkestack.io/tke/pkg/platform/controller/addon/tappcontroller/images"
	baremetal "tkestack.io/tke/pkg/platform/provider/baremetal/images"
	csioperator "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	galaxy "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	"tkestack.io/tke/pkg/spec"
)

var (
	specialUnsupportMultiArch = []string{"nvidia-device-plugin", "gpu"}
)

func main() {
	archsFlag := pflag.StringSliceP("archs", "a", spec.Archs, "Only list images for specified archs")
	pflag.Parse()
	unsupportMultiArchImages := []func() []string{
		cronhpa.List,
		lbcf.List,
		logcollector.List,
		persistentevent.List,
		prometheus.List,
		csioperator.List,
		volumedecorator.List,
		tappcontroller.List,
		logagent.List,
	}
	supportMultiArchImages := []func() []string{
		baremetal.List,
		installer.List,
		galaxy.List,
		ipam.List,
		mesh.List,
	}

	var result []string
	for _, f := range supportMultiArchImages {
		for _, one := range f() {
			if IsUnsupportMultiArch(one) {
				result = append(result, one)
			} else {
				for _, arch := range *archsFlag {
					result = append(result, strings.ReplaceAll(one, ":", "-"+arch+":"))
				}
			}
		}
	}

	for _, f := range unsupportMultiArchImages {
		images := f()
		result = append(result, images...)
	}

	result = funk.UniqString(result)
	sort.Strings(result)
	for _, one := range result {
		fmt.Println(one)
	}
}

func IsUnsupportMultiArch(name string) bool {
	for _, one := range specialUnsupportMultiArch {
		if strings.HasPrefix(name, one) {
			return true
		}
	}

	return false
}
