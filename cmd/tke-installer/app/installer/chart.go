/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package installer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	applicationutil "tkestack.io/tke/pkg/application/util"
	"tkestack.io/tke/pkg/registry/config"
	chartpath "tkestack.io/tke/pkg/registry/util/chartpath/v1"
	"tkestack.io/tke/pkg/util/compress"
	"tkestack.io/tke/pkg/util/files"
	"tkestack.io/tke/pkg/util/log"
)

const chartFilesSuffix = ".charts.tar.gz"

var (
	needImportedChartGroups = []string{"public"}
	// defaultChartGroups must include needImportedChartGroups
	defaultChartGroups             = append(needImportedChartGroups, []string{}...)
	defaultChartGroupsStringConfig = ""
)

func init() {
	DefaultChartGroupsBytesConfig, err := json.Marshal(defaultChartGroups)
	if err != nil {
		log.Fatalf("init DefaultChartGroupsStringConfig failed", err)
	}
	defaultChartGroupsStringConfig = string(DefaultChartGroupsBytesConfig)
}

func (t *TKE) importCharts(ctx context.Context) error {
	var errs []error
	client := applicationutil.NewHelmClientWithoutRESTClient()
	for _, chartGroup := range needImportedChartGroups {
		dest, err := ioutil.TempDir("", "chartpath-")
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer os.RemoveAll(dest)

		err = compress.ExtractTarGz(chartGroup+chartFilesSuffix, dest)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		files, err := files.GetAllFiles(dest)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		cgs, err := t.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", constants.DefaultTeantID, chartGroup),
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(cgs.Items) == 0 {
			errs = append(errs, fmt.Errorf("cannot find %s chartgroup", chartGroup))
			continue
		}
		conf := config.RepoConfiguration{
			Scheme:        "http",
			DomainSuffix:  t.Para.Config.Registry.Domain(),
			Admin:         t.Para.Config.Registry.Username(),
			AdminPassword: string(t.Para.Config.Registry.Password()),
		}

		chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(conf, cgs.Items[0])
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, f := range files {
			existed, err := client.Push(&helmaction.PushOptions{
				ChartPathOptions: chartPathBasicOptions,
				ChartFile:        f,
				ForceUpload:      false,
			})
			if err != nil {
				if existed {
					log.Warn(err.Error())
				} else {
					errs = append(errs, err)
					continue
				}
			}
		}
	}
	return utilerrors.NewAggregate(errs)

}

func (t *TKE) checkNeedImportedChartgroups(ctx context.Context) error {
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		for _, chartGroup := range needImportedChartGroups {
			cgs, err := t.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", constants.DefaultTeantID, chartGroup),
			})
			if err != nil {
				return false, err
			}
			if len(cgs.Items) == 0 {
				return false, fmt.Errorf("cannot find %s chartgroup", chartGroup)
			}
		}
		return true, nil
	})
}
