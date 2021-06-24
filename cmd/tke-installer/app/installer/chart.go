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
	for _, chartGroup := range needImportedChartGroups {
		err := t.pushCharts(ctx, chartGroup+chartFilesSuffix, constants.DefaultTeantID, chartGroup)
		errs = append(errs, err)
	}
	return utilerrors.NewAggregate(errs)

}

func (t *TKE) pushCharts(ctx context.Context, chartsPath, tenantID, chartGroup string) error {
	var errs []error
	client := applicationutil.NewHelmClientWithoutRESTClient()
	dest, err := ioutil.TempDir("", "chartpath-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dest)

	err = compress.ExtractTarGz(chartsPath, dest)
	if err != nil {
		return err
	}

	files, err := files.GetAllFiles(dest)
	if err != nil {
		return err
	}

	var cg registryv1.ChartGroup
	cgs, err := t.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, chartGroup),
	})
	if err != nil {
		return err
	}
	if len(cgs.Items) == 0 {
		t.log.Infof("%s chart group doesn't exist, try to create", chartGroup)
		cg, err = t.createChartgroup(ctx, tenantID, chartGroup)
		if err != nil {
			return fmt.Errorf("create chartgroup failed: %v", err)
		}
		t.log.Info("%s chart group is created", chartGroup)
	} else {
		cg = cgs.Items[0]
	}
	conf := config.RepoConfiguration{
		Scheme:        "http",
		DomainSuffix:  t.Para.Config.Registry.Domain(),
		Admin:         t.Para.Config.Registry.Username(),
		AdminPassword: string(t.Para.Config.Registry.Password()),
	}

	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(conf, cg)
	if err != nil {
		return err
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
		}
	}
	return utilerrors.NewAggregate(errs)
}

func (t *TKE) createChartgroup(ctx context.Context, tenantID, name string) (registryv1.ChartGroup, error) {
	cg, err := t.registryClient.ChartGroups().Create(ctx, &registryv1.ChartGroup{
		Spec: registryv1.ChartGroupSpec{
			Name:     name,
			TenantID: tenantID,
		},
	},
		metav1.CreateOptions{})
	if err != nil {
		return *cg, err
	}
	err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		cgs, err := t.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, name),
		})
		if err != nil {
			return false, err
		}
		if len(cgs.Items) == 0 {
			return false, fmt.Errorf("create chargroup failed, cannot find %s chartgroup", name)
		}
		return true, nil
	})
	return *cg, err
}
