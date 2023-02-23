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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	applicationv1 "tkestack.io/tke/api/application/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/apiclient"
)

func (t *TKE) completeExpansionApps() error {

	if len(t.Para.Config.ExpansionApps) == 0 {
		return nil
	}

	for _, expansionApp := range t.Para.Config.ExpansionApps {
		if !expansionApp.Enable {
			continue
		}
		err := t.completeChart(&expansionApp.Chart)
		if err != nil {
			return fmt.Errorf("bad platform app config. %v, %v", expansionApp.Name, err)
		}
	}

	t.log.Infof("tke platform apps config completed.")
	t.backup()
	return nil
}

func (t *TKE) completeChart(chart *types.Chart) error {

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

	if len(t.Para.Config.ExpansionApps) == 0 {
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

	for _, expansionApp := range t.Para.Config.ExpansionApps {
		if !expansionApp.Enable {
			continue
		}
		if t.applicationAlreadyInstalled(expansionApp, apps.Items) {
			t.log.Infof("application already exists. we don't override applications while installing. %v/%v", expansionApp.Chart.TargetNamespace, expansionApp.Chart.Name)
			continue
		}
		err := t.installApplication(ctx, expansionApp)
		if err != nil {
			return fmt.Errorf("install application failed. %v, %v", expansionApp.Name, err)
		}
		t.log.Infof("finish application installation %v", expansionApp.Name)
	}

	return nil
}

func (t *TKE) applicationAlreadyInstalled(expansionApp *types.ExpansionApp, installedApps []applicationv1.App) bool {

	for _, installedApp := range installedApps {
		// if there's an existed application with the same namespace+name, we consider it as already exists
		if expansionApp.Name == installedApp.Spec.Name && expansionApp.Chart.TargetNamespace == installedApp.Namespace {
			return true
		}
	}
	return false
}

func (t *TKE) installApplication(ctx context.Context, expansionApp *types.ExpansionApp) error {

	chart := expansionApp.Chart

	rawValues, err := chart.Values.YAML()
	if err != nil {
		return err
	}

	app := &applicationv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:                 chart.TargetNamespace,
			ZZZ_DeprecatedClusterName: chart.TargetCluster,
		},
		Spec: applicationv1.AppSpec{
			Type:          constants.DefaultApplicationInstallDriverType,
			TenantID:      chart.TenantID,
			Name:          expansionApp.Name,
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
func (t *TKE) initPlatformApps(ctx context.Context) error {
	defaultPlatformApps := []*types.PlatformApp{}
	if t.Para.Config.Auth.TKEAuth != nil {
		authAPIOptions, err := t.getTKEAuthAPIOptions(ctx)
		if err != nil {
			return fmt.Errorf("get tke-auth-api options failed: %v", err)
		}
		tkeAuth := &types.PlatformApp{
			HelmInstallOptions: &helmaction.InstallOptions{
				Namespace:   t.namespace,
				ReleaseName: "tke-auth",
				Values: map[string]interface{}{
					"api":        authAPIOptions,
					"controller": t.getTKEAuthControllerOptions(ctx),
				},
				DependencyUpdate: false,
				ChartPathOptions: helmaction.ChartPathOptions{},
			},
			LocalChartPath: constants.ChartDirName + "tke-auth/",
			Enable:         true,
			ConditionFunc: func() (bool, error) {
				apiOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-auth-api")
				if err != nil {
					return false, nil
				}
				controllerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-auth-controller")
				if err != nil {
					return false, nil
				}
				return apiOk && controllerOk, nil
			},
		}
		defaultPlatformApps = append(defaultPlatformApps, tkeAuth)
	}

	platformAPIOptions, err := t.getTKEPlatformAPIOptions(ctx)
	if err != nil {
		return fmt.Errorf("get tke-platform-api options failed: %v", err)
	}
	versionsByte, err := json.Marshal(spec.K8sVersions)
	if err != nil {
		return fmt.Errorf("get tke-platform k8sValidVersions options failed: %v", err)
	}
	tkePlatform := &types.PlatformApp{
		HelmInstallOptions: &helmaction.InstallOptions{
			Namespace:   t.namespace,
			ReleaseName: "tke-platform",
			Values: map[string]interface{}{
				"api":              platformAPIOptions,
				"controller":       t.getTKEPlatformControllerOptions(ctx),
				"k8sValidVersions": string(versionsByte),
				"tkeVersion":       spec.TKEVersion,
			},
			DependencyUpdate: false,
			ChartPathOptions: helmaction.ChartPathOptions{},
		},
		LocalChartPath: constants.ChartDirName + "tke-platform/",
		Enable:         true,
		ConditionFunc: func() (bool, error) {
			apiOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-platform-api")
			if err != nil {
				return false, nil
			}
			controllerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-platform-controller")
			if err != nil {
				return false, nil
			}
			return apiOk && controllerOk, nil
		},
	}
	defaultPlatformApps = append(defaultPlatformApps, tkePlatform)
	t.Para.Config.PlatformApps = append(defaultPlatformApps, t.Para.Config.PlatformApps...)
	return nil
}

func (t *TKE) preprocessPlatformApps(ctx context.Context) error {

	for _, platformApp := range t.Para.Config.PlatformApps {
		if !platformApp.Enable || platformApp.Installed {
			continue
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephRBDChartReleaseName) {
			platformApp.ConditionFunc = func() (bool, error) {
				provisionerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, platformApp.HelmInstallOptions.Namespace, "ceph-csi-rbd-provisioner")
				if err != nil {
					return false, nil
				}
				nodepluginOk, err := apiclient.CheckDaemonset(ctx, t.globalClient, platformApp.HelmInstallOptions.Namespace, "ceph-csi-rbd-nodeplugin")
				if err != nil {
					return false, nil
				}
				return provisionerOk && nodepluginOk, nil
			}
			platformApp.LocalChartPath = constants.ChartDirName + "ceph-csi-rbd/"

			if err := t.mergePlatformAppValues(platformApp); err != nil {
				return err
			}
			values := platformApp.HelmInstallOptions.Values
			if values["csiConfig"] == nil {
				return fmt.Errorf("ceph-csi-rbd platformAPP csiConfig nil")
			}
			if values["secret"] == nil || len(values["secret"].(map[string]interface{})["userID"].(string)) == 0 || len(values["secret"].(map[string]interface{})["userKey"].(string)) == 0 {
				return fmt.Errorf("ceph-csi-rbd platformAPP secret userID | userKey nil")
			}
			if values["storageClass"] == nil || len(values["storageClass"].(map[string]interface{})["clusterID"].(string)) == 0 {
				return fmt.Errorf("ceph-csi-rbd platformAPP storageClass clusterID nil")
			}

			values["images"] = map[string]interface{}{
				"enable": true,
				"nodeplugin": map[string]interface{}{
					"registrar": images.Get().CsiNodeDriverRegistrar.FullName(),
					"plugin":    images.Get().CephCsi.FullName(),
				},
				"provisioner": map[string]interface{}{
					"provisioner": images.Get().CsiProvisioner.FullName(),
					"attacher":    images.Get().CsiAttacher.FullName(),
					"resizer":     images.Get().CsiResizer.FullName(),
					"snapshotter": images.Get().CsiSnapshotter.FullName(),
				},
			}
			if values["storageClass"] == nil {
				values["storageClass"] = map[string]interface{}{
					"replicaCount": t.Config.Replicas,
					"name":         constants.CephRBDStorageClassName,
				}
			} else {
				values["storageClass"].(map[string]interface{})["replicaCount"] = t.Config.Replicas
				values["storageClass"].(map[string]interface{})["name"] = constants.CephRBDStorageClassName
			}
			values["namespace"] = platformApp.HelmInstallOptions.Namespace
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephFSChartReleaseName) {
			platformApp.ConditionFunc = func() (bool, error) {
				provisionerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, platformApp.HelmInstallOptions.Namespace, "ceph-csi-cephfs-provisioner")
				if err != nil {
					return false, nil
				}
				nodepluginOk, err := apiclient.CheckDaemonset(ctx, t.globalClient, platformApp.HelmInstallOptions.Namespace, "ceph-csi-cephfs-nodeplugin")
				if err != nil {
					return false, nil
				}
				return provisionerOk && nodepluginOk, nil
			}
			platformApp.LocalChartPath = constants.ChartDirName + "ceph-csi-cephfs/"

			if err := t.mergePlatformAppValues(platformApp); err != nil {
				return err
			}
			values := platformApp.HelmInstallOptions.Values
			if values["csiConfig"] == nil {
				return fmt.Errorf("ceph-csi-cephfs platformAPP csiConfig nil")
			}
			// TODO: little confused in preprocess step
			if values["secret"] == nil || len(values["secret"].(map[string]interface{})["adminID"].(string)) == 0 || len(values["secret"].(map[string]interface{})["adminKey"].(string)) == 0 {
				return fmt.Errorf("ceph-csi-cephfs platformAPP secret adminID | adminKey nil")
			}
			// TODO: little confused in preprocess step
			if values["storageClass"] == nil || len(values["storageClass"].(map[string]interface{})["clusterID"].(string)) == 0 || len(values["storageClass"].(map[string]interface{})["fsName"].(string)) == 0 {
				return fmt.Errorf("ceph-csi-cephfs platformAPP storageClass clusterID| fsName nil")
			}

			values["images"] = map[string]interface{}{
				"enable": true,
				"nodeplugin": map[string]interface{}{
					"registrar": images.Get().CsiNodeDriverRegistrar.FullName(),
					"plugin":    images.Get().CephCsi.FullName(),
				},
				"provisioner": map[string]interface{}{
					"provisioner": images.Get().CsiProvisioner.FullName(),
					"attacher":    images.Get().CsiAttacher.FullName(),
					"resizer":     images.Get().CsiResizer.FullName(),
					"snapshotter": images.Get().CsiSnapshotter.FullName(),
				},
			}
			if values["storageClass"] == nil {
				values["storageClass"] = map[string]interface{}{
					"replicaCount": t.Config.Replicas,
					"name":         constants.CephFSStorageClassName,
				}
			} else {
				values["storageClass"].(map[string]interface{})["replicaCount"] = t.Config.Replicas
				values["storageClass"].(map[string]interface{})["name"] = constants.CephFSStorageClassName
			}
			values["namespace"] = platformApp.HelmInstallOptions.Namespace
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.NFSChartReleaseName) {
			platformApp.ConditionFunc = func() (bool, error) {
				provisionerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, platformApp.HelmInstallOptions.Namespace, "nfs-subdir-external-provisioner")
				if err != nil {
					return false, nil
				}
				return provisionerOk, nil
			}
			platformApp.LocalChartPath = constants.ChartDirName + "nfs-subdir-external-provisioner/"

			if err := t.mergePlatformAppValues(platformApp); err != nil {
				return err
			}
			values := platformApp.HelmInstallOptions.Values
			if values == nil || values["nfs"] == nil || len(values["nfs"].(map[string]interface{})["server"].(string)) == 0 || len(values["nfs"].(map[string]interface{})["path"].(string)) == 0 {
				return fmt.Errorf("nfs-subdir-external-provisioner platformApp lack nfs.server or nfs.path")
			}
			values["replicaCount"] = t.Config.Replicas
			values["namespace"] = platformApp.HelmInstallOptions.Namespace
			if values["images"] == nil {
				values["images"] = map[string]interface{}{
					"repository": images.Get().NFSProvisioner.FullName(),
				}
			} else {
				values["images"].(map[string]interface{})["repository"] = images.Get().NFSProvisioner.FullName()
			}
		}
	}
	return nil
}

func (t *TKE) mergePlatformAppValues(platformApp *types.PlatformApp) error {
	if len(platformApp.RawValues) != 0 || len(platformApp.Values) != 0 {
		values, err := helmutil.MergeValues(platformApp.Values, platformApp.RawValues, string(platformApp.RawValuesType))
		if err != nil {
			return err
		}
		platformApp.HelmInstallOptions.Values = values
		platformApp.RawValues = ""
		platformApp.Values = nil
	}
	return nil
}

func (t *TKE) installPlatformApps(ctx context.Context) error {

	if len(t.Para.Config.PlatformApps) == 0 {
		return nil
	}
	for i, platformApp := range t.Para.Config.PlatformApps {
		if !platformApp.Enable || platformApp.Installed {
			continue
		}
		err := t.installPlatformApp(ctx, platformApp)
		if err != nil {
			t.log.Errorf("Install %s failed", platformApp.HelmInstallOptions.ReleaseName)
			return err
		}
		t.Para.Config.PlatformApps[i].Installed = true
	}

	return nil
}

func (t *TKE) installPlatformApp(ctx context.Context, platformApp *types.PlatformApp) error {
	platformApp.HelmInstallOptions.Timeout = 10 * time.Minute
	if err := t.mergePlatformAppValues(platformApp); err != nil {
		return err
	}
	// TODO currently only support local chart install
	if len(platformApp.LocalChartPath) != 0 {
		t.log.Infof("Start install platform app %s in %s namespace", platformApp.HelmInstallOptions.ReleaseName, platformApp.HelmInstallOptions.Namespace)
		if _, err := t.helmClient.InstallWithLocal(ctx, platformApp.HelmInstallOptions, platformApp.LocalChartPath); err != nil {
			uninstallOptions := helmaction.UninstallOptions{
				Timeout:     10 * time.Minute,
				ReleaseName: platformApp.HelmInstallOptions.ReleaseName,
				Namespace:   platformApp.HelmInstallOptions.Namespace,
			}
			_, uninstallErr := t.helmClient.Uninstall(&uninstallOptions)
			if uninstallErr != nil {
				return fmt.Errorf("clean %s failed %v", platformApp.HelmInstallOptions.ReleaseName, uninstallErr)
			}
			return err
		}
	}
	if platformApp.ConditionFunc != nil {
		err := wait.PollImmediate(5*time.Second, 10*time.Minute, platformApp.ConditionFunc)
		if err != nil {
			return err
		}
	}
	t.log.Infof("End install platform app %s in %s namespace", platformApp.HelmInstallOptions.ReleaseName, platformApp.HelmInstallOptions.Namespace)
	return nil
}
