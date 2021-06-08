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

const (
	applicationInstallNamespaceDefault      = "default"
	applicationInstallClusterDefault        = "global"
	applicationInstallTenantIDDefault       = constants.DefaultTeantID
	applicationInstallChartGroupNameDefault = constants.DefaultChartGroupName
	applicationInstallDriverTypeDefault     = "HelmV3"
	applicationInstallValueTypeDefault      = "yaml"
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
		chart.TenantID = applicationInstallTenantIDDefault
	}
	if chart.ChartGroupName == "" {
		chart.ChartGroupName = applicationInstallChartGroupNameDefault
	}
	if chart.TargetCluster == "" {
		chart.TargetCluster = applicationInstallClusterDefault
	}
	if chart.TargetNamespace == "" {
		chart.TargetCluster = applicationInstallNamespaceDefault
	}

	return nil
}

func (t *TKE) installApplications(ctx context.Context) error {

	if len(t.Config.PlatformApps) == 0 {
		return nil
	}

	// TODO: (workaround ) client init will only be called in createCluster.
	// If we SKIP createCluster step, all client calls will be panic
	err := t.initDataForDeployTKE()
	if err != nil {
		return err
	}

	apps, err := t.applicationClient.Apps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("list all applications failed %v", err)
	}

	for _, platformApp := range t.Config.PlatformApps {
		if !platformApp.Enable {
			continue
		}
		err := t.installApplication(ctx, platformApp, apps.Items)
		if err != nil {
			return fmt.Errorf("install application failed. %v, %v", platformApp.Name, err)
		}
		t.log.Infof("finish application installation %v", platformApp.Name)
	}

	return nil
}

func (t *TKE) installApplication(ctx context.Context, platformApp config.PlatformApp, installedApps []applicationv1.App) error {

	chart := platformApp.Chart

	var found bool
	var duplicated bool
	for _, installedApp := range installedApps {
		// if there's an existed application with the same namespace+name, we consider it as already exists
		if chart.Name == installedApp.Spec.Name && chart.TargetNamespace == installedApp.Namespace {
			found = true
			break
		}
		// if there's an existed application with different name/namespace but with same chart name+version, we consider it as duplicated
		if chart.Name == installedApp.Spec.Chart.ChartName && chart.Version == installedApp.Spec.Chart.ChartVersion {
			duplicated = true
			break
		}
	}
	if found {
		t.log.Infof("application already exists. we don't override applications while installing. %v/%v", chart.TargetNamespace, chart.Name)
		return nil
	}
	if duplicated {
		return fmt.Errorf("duplicate chart detecated. %v, %v", chart.Name, chart.Version)
	}

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
			Type:          applicationInstallDriverTypeDefault,
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
				RawValuesType: applicationInstallValueTypeDefault,
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
