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
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	"path"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	registryv1 "tkestack.io/tke/api/registry/v1"
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

	// Expansion
	var parseChartPackage = func(filepath string) (name string, version string) {
		filename := path.Base(filepath)
		filename = strings.TrimSuffix(filename, ".tgz")
		filename = strings.TrimSuffix(filename, ".tar.gz")
		arr := strings.Split(filename, "-")
		if len(arr) == 1 {
			return arr[0], ""
		}
		name = strings.Join(arr[:len(arr)-1], "-")
		version = arr[len(arr)-1]
		return
	}

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

		// Expansion
		// TODO: now we only support expanding public charts
		err = t.expansionDriver.CopyChartsToDst(chartGroup, dest)
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

		cg := cgs.Items[0]
		chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(conf, cgs.Items[0])
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, f := range files {
			// push chart with force flag
			_, err := client.Push(&helmaction.PushOptions{
				ChartPathOptions: chartPathBasicOptions,
				ChartFile:        f,
				ForceUpload:      true,
			})
			if err != nil {
				errs = append(errs, err)
				continue
			}
			chartName, _ := parseChartPackage(f)
			chart := &registryv1.Chart{
				Spec: registryv1.ChartSpec{
					Name:           chartName,
					TenantID:       cg.Spec.TenantID,
					ChartGroupName: cg.Spec.Name,
					DisplayName:    chartName,
					Visibility:     cg.Spec.Visibility,
				},
			}
			// TODO: this not works by now. cause chart CR name is not chart name
			c, err := t.registryClient.Charts(cg.Name).Get(ctx, chartName, metav1.GetOptions{})
			if err == nil {
				t.log.Infof("chart already exists %v", chart)
				continue
			}
			t.log.Infof("do not get chart by name %v, %v", chartName, c)
			if errors.IsNotFound(err) {
				_, err = t.registryClient.Charts(cg.Name).Create(ctx, chart, metav1.CreateOptions{})
				if err != nil {
					// TODO: workaround.
					if strings.Contains(err.Error(), "spec.name: Duplicate value") {
						t.log.Infof("create chart duplicated %v, %+v, %v", f, chart, err)
						continue
					}
					t.log.Errorf("create chart failed %v, %+v, %v", f, chart, err)
					errs = append(errs, err)
				}
				continue
			}
			t.log.Errorf("get chart error %v, %v", chart, err)
			errs = append(errs, err)
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
