/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package installer

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applicationv1 "tkestack.io/tke/api/application/v1"
	"tkestack.io/tke/cmd/tke-installer/app/config"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
)

func (t *TKE) completePlatformApps() error {

	if len(t.Config.PlatformApps) == 0 {
		return nil
	}

	for _, platformApp := range t.Config.PlatformApps {
		if !platformApp.Enable {
			continue
		}
		err := t.completeChart(&platformApp.Chart)
		if err != nil {
			return fmt.Errorf("bad platform app config. %v, %v", platformApp.Name, err)
		}
	}

	t.log.Infof("tke platform apps config completed.")
	t.backup()
	return nil
}

func (t *TKE) completeChart(chart *config.Chart) error {

	_, err := chart.Values.YAML()
	if err != nil {
		return err
	}

	if chart.Name == "" {
		return fmt.Errorf("chart name empty")
	}
	if chart.Version == "" {
		return fmt.Errorf("chart value empty")
	}
	if chart.TenantID == "" {
		chart.TenantID = constants.DefaultTeantID
	}
	if chart.ChartGroupName == "" {
		chart.ChartGroupName = constants.DefaultChartGroupName
	}
	if chart.TargetCluster == "" {
		chart.TargetCluster = constants.GlobalClusterName
	}
	if chart.TargetNamespace == "" {
		chart.TargetNamespace = metav1.NamespaceDefault
	}

	return nil
}

func (t *TKE) installApplications(ctx context.Context) error {

	if len(t.Config.PlatformApps) == 0 {
		return nil
	}

	// TODO: (workaround ) client init will only be called in createCluster.
	// If we SKIP createCluster step, all client calls will be panic
	if t.applicationClient == nil {
		err := t.initDataForDeployTKE()
		if err != nil {
			return err
		}
	}

	apps, err := t.applicationClient.Apps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("list all applications failed %v", err)
	}

	for _, platformApp := range t.Config.PlatformApps {
		if !platformApp.Enable {
			continue
		}
		if t.applicationAlreadyInstalled(platformApp, apps.Items) {
			t.log.Infof("application already exists. we don't override applications while installing. %v/%v", platformApp.Chart.TargetNamespace, platformApp.Chart.Name)
			continue
		}
		err := t.installApplication(ctx, platformApp)
		if err != nil {
			return fmt.Errorf("install application failed. %v, %v", platformApp.Name, err)
		}
		t.log.Infof("finish application installation %v", platformApp.Name)
	}

	return nil
}

func (t *TKE) applicationAlreadyInstalled(platformApp config.PlatformApp, installedApps []applicationv1.App) bool {

	for _, installedApp := range installedApps {
		// if there's an existed application with the same namespace+name, we consider it as already exists
		if platformApp.Name == installedApp.Spec.Name && platformApp.Chart.TargetNamespace == installedApp.Namespace {
			return true
		}
	}
	return false
}

func (t *TKE) installApplication(ctx context.Context, platformApp config.PlatformApp) error {

	chart := platformApp.Chart

	rawValues, err := chart.Values.YAML()
	if err != nil {
		return err
	}

	app := &applicationv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   chart.TargetNamespace,
			ClusterName: chart.TargetCluster,
		},
		Spec: applicationv1.AppSpec{
			Type:          constants.DefaultApplicationInstallDriverType,
			TenantID:      chart.TenantID,
			Name:          platformApp.Name,
			TargetCluster: chart.TargetCluster,
			Chart: applicationv1.Chart{
				TenantID:       chart.TenantID,
				ChartGroupName: chart.ChartGroupName,
				ChartName:      chart.Name,
				ChartVersion:   chart.Version,
			},
			Values: applicationv1.AppValues{
				RawValuesType: constants.DefaultApplicationInstallValueType,
				RawValues:     rawValues,
			},
		},
	}
	_, err = t.applicationClient.Apps(chart.TargetNamespace).Create(ctx, app, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create application failed %v,%v", chart.Name, err)
	}
	return nil
}
