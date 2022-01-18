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
	prometheus "tkestack.io/tke/pkg/monitor/controller/prometheus/images"
	cronhpa "tkestack.io/tke/pkg/platform/controller/addon/cronhpa/images"
	persistentevent "tkestack.io/tke/pkg/platform/controller/addon/persistentevent/images"
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
	baseFlag := pflag.Bool("base", false, "Only list base components")
	extraFlag := pflag.Bool("extra", false, "Only list extra components")
	pflag.Parse()
	if *baseFlag {
		for _, one := range baseComponents(*archsFlag) {
			fmt.Println(one)
		}
		return
	}
	if *extraFlag {
		for _, one := range exComponents(*archsFlag) {
			fmt.Println(one)
		}
		return
	}
	for _, one := range append(baseComponents(*archsFlag), exComponents(*archsFlag)...) {
		fmt.Println(one)
	}
}

func baseComponents(archsFlag []string) []string {
	supportMultiArchImages := []func() []string{
		baremetal.List,
		installer.ListBaseComponents,
		galaxy.List,
	}

	var result []string
	for _, f := range supportMultiArchImages {
		for _, one := range f() {
			if IsUnsupportMultiArch(one) {
				result = append(result, one)
			} else {
				for _, arch := range archsFlag {
					result = append(result, strings.ReplaceAll(one, ":", "-"+arch+":"))
				}
			}
		}
	}

	result = funk.UniqString(result)
	sort.Strings(result)
	return result
}

func exComponents(archsFlag []string) []string {
	unsupportMultiArchImages := []func() []string{
		cronhpa.List,
		persistentevent.List,
		prometheus.List,
		csioperator.List,
		tappcontroller.List,
		logagent.List,
	}
	supportMultiArchImages := []func() []string{
		installer.ListExComponents,
		mesh.List,
	}

	var result []string
	for _, f := range supportMultiArchImages {
		for _, one := range f() {
			if IsUnsupportMultiArch(one) {
				result = append(result, one)
			} else {
				for _, arch := range archsFlag {
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
	return result
}

func IsUnsupportMultiArch(name string) bool {
	for _, one := range specialUnsupportMultiArch {
		if strings.HasPrefix(name, one) {
			return true
		}
	}

	return false
}
