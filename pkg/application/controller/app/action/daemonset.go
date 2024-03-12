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

package action

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"helm.sh/helm/v3/pkg/kube"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	helmchart "tkestack.io/tke/pkg/helm"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/metrics"
)

const (
	damonsetsUpgradeSuccess = "daemonsets upgrade success"
	daemonsetsUpgradeFailed = "daemonsets upgrade failed"
	daemonsetsUpgradeSkiped = "daemonsets upgrade skiped"
	// daemonsetsUpgradeRunning = "daemonsets upgrade running"

	upgradeJobLabelPrefix        = "upgradejob.application.tkestack.io/"
	upgradeJobAppNameLabel       = upgradeJobLabelPrefix + "appname"
	upgradeJobAppGenerationLabel = upgradeJobLabelPrefix + "appgeneration"
)

func UpgradeDaemonset(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration,
	updateStatusFunc applicationprovider.UpdateStatusFunc) (*applicationv1.App, error) {

	if app.Status.Message == daemonsetsUpgradeFailed {
		log.Infof("upgrade ds for app %s/%s failed: %s", app.Namespace, app.Name, app.Status.Reason)
		return app, nil
	}

	ujs, err := applicationClient.UpgradeJobs(app.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			upgradeJobAppNameLabel:       app.Name,
			upgradeJobAppGenerationLabel: strconv.Itoa(int(app.Generation)),
		}).String(),
	})
	if err != nil {
		log.Errorf("get UpgradeJobs for app %s/%s failed: %v", app.Namespace, app.Name, err)
		return nil, err
	}
	log.Infof("get UpgradeJobs for app %s/%s count: %d", app.Namespace, app.Name, len(ujs.Items))

	allDone := true
	for _, uj := range ujs.Items {
		if uj.Status.Reason != nil {
			log.Errorf("upgrade ds %s for app %s/%s failed: %s", uj.Spec.Target, app.Namespace, app.Name, *uj.Status.Reason)
			// 切换到failed阶段，中断升级
			newStatus := app.Status.DeepCopy()
			newStatus.Phase = applicationv1.AppPhaseUpgradingDaemonsetFailed
			newStatus.Message = daemonsetsUpgradeFailed
			newStatus.Reason = *uj.Status.Reason
			newStatus.LastTransitionTime = metav1.Now()
			return updateStatusFunc(ctx, app, &app.Status, newStatus)
		}

		if uj.Status.BatchCompleteNum <= uj.Status.BatchOrder {
			allDone = false
		}
	}

	if !allDone {
		return app, nil
	}

	// 走后续的post流程
	log.Infof("UpgradeJobs for app %s/%s has done: begin PostUpgrade", app.Namespace, app.Name)
	hooks := getHooks(app)
	err = hooks.PostUpgrade(ctx, applicationClient, platformClient, app, repo, updateStatusFunc)
	if err != nil {
		newStatus := app.Status.DeepCopy()
		newStatus.Phase = applicationv1.AppPhaseUpgradingDaemonsetFailed
		newStatus.Message = "hook post upgrade app failed"
		newStatus.Reason = err.Error()
		newStatus.LastTransitionTime = metav1.Now()
		metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(1)
		var updateStatusErr error
		app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
		if updateStatusErr != nil {
			return app, updateStatusErr
		}
		return app, err
	}

	newStatus := app.Status.DeepCopy()
	newStatus.Phase = applicationv1.AppPhaseSucceeded
	newStatus.Message = ""
	newStatus.Reason = ""
	newStatus.LastTransitionTime = metav1.Now()
	metrics.GaugeApplicationInstallFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationUpgradeFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationRollbackFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	metrics.GaugeApplicationManifestFailed.WithLabelValues(app.Spec.TargetCluster, app.Name).Set(0)
	var updateStatusErr error
	app, updateStatusErr = updateStatusFunc(ctx, app, &app.Status, newStatus)
	if updateStatusErr != nil {
		return app, updateStatusErr
	}

	return app, nil
}

func getOndeleteDaemonsets(ctx context.Context,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	repo appconfig.RepoConfiguration) ([]string, error) {

	helper, err := helmchart.NewHelper(ctx, platformClient, app, repo)
	if err != nil {
		log.Errorf(fmt.Sprintf("new cluster %s app %s helper failed, err:%s", app.Spec.TargetCluster, app.Name, err.Error()))
		return nil, err
	}
	clusterManifest, err := helper.GetClusterManifest()
	if err != nil {
		log.Errorf(fmt.Sprintf("get cluster %s app %s cluster manifest failed, err:%s", app.Spec.TargetCluster, app.Name, err.Error()))
		return nil, err
	}
	manifest, err := helper.GetManifest()
	if err != nil {
		log.Errorf(fmt.Sprintf("get cluster %s app %s manifest failed, err:%s", app.Spec.TargetCluster, app.Name, err.Error()))
		return nil, err
	}
	if clusterManifest != manifest {
		// 两次前后渲染不一致，重新触发检查
		log.Errorf(fmt.Sprintf("cluster %s app %s, cluster manifest is not equal manifest", app.Spec.TargetCluster, app.Name))
		return nil, err
	}
	clusterResource, err := helper.GetClusterResource()
	if err != nil {
		log.Errorf(fmt.Sprintf("get cluster %s app %s resource failed, err:%s", app.Spec.TargetCluster, app.Name, err.Error()))
		return nil, err
	}

	var ds []string
	for _, v := range clusterResource {
		switch kube.AsVersioned(v).(type) {
		case *appsv1.DaemonSet:
			daemonset := appsv1.DaemonSet{}
			daemonsetObj := v.Object.(*unstructured.Unstructured)
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(daemonsetObj.UnstructuredContent(), &daemonset); err != nil {
				log.Errorf("getOndeleteDaemonsets for app %s/%s failed: %v", app.Namespace, app.Name, err)
				return nil, err
			}
			if daemonset.Spec.UpdateStrategy.Type == appsv1.OnDeleteDaemonSetStrategyType {
				ds = append(ds, daemonset.Namespace+"/"+daemonset.Name)
			}
			/* TODO: addon-charts中没有beta版本的，job-controller目前不支持beta版本的ds灰度升级
			case *extensionsv1beta1.DaemonSet:
				daemonset := extensionsv1beta1.DaemonSet{}
				daemonsetObj := v.Object.(*unstructured.Unstructured)
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(daemonsetObj.UnstructuredContent(), &daemonset); err != nil {
					log.Errorf("getOndeleteDaemonsets for app %s/%s failed: %v", app.Namespace, app.Name, err)
					return nil, err
				}
				if daemonset.Spec.UpdateStrategy.Type == extensionsv1beta1.OnDeleteDaemonSetStrategyType {
					ds = append(ds, daemonset.Namespace+"/"+daemonset.Name)
				}
			case *appsv1beta2.DaemonSet:
				daemonset := appsv1beta2.DaemonSet{}
				daemonsetObj := v.Object.(*unstructured.Unstructured)
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(daemonsetObj.UnstructuredContent(), &daemonset); err != nil {
					log.Errorf("getOndeleteDaemonsets for app %s/%s failed: %v", app.Namespace, app.Name, err)
					return nil, err
				}
				if daemonset.Spec.UpdateStrategy.Type == appsv1beta2.OnDeleteDaemonSetStrategyType {
					ds = append(ds, daemonset.Namespace+"/"+daemonset.Name)
				}
			*/
		case *appsv1.StatefulSet, *appsv1beta1.StatefulSet, *appsv1beta2.StatefulSet:
		case *appsv1.Deployment, *appsv1beta1.Deployment, *appsv1beta2.Deployment, *extensionsv1beta1.Deployment:
		default:
		}
	}
	return ds, nil
}

func createUpgradeJobsIfNotExist(ctx context.Context,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	app *applicationv1.App,
	daemonsets []string) error {

	ujs, err := applicationClient.UpgradeJobs(app.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{
			upgradeJobAppNameLabel: app.Name,
		}).String(),
	})
	if err != nil {
		log.Errorf("get UpgradeJobs for app %s/%s failed: %v", app.Namespace, app.Name, err)
		return err
	}

	existedJob := map[string]bool{}
	for _, uj := range ujs.Items {
		// 删除之前版本对应的job
		if uj.Labels[upgradeJobAppGenerationLabel] != strconv.Itoa(int(app.Generation)) {
			if err := applicationClient.UpgradeJobs(uj.Namespace).Delete(context.TODO(), uj.Name, metav1.DeleteOptions{}); err != nil {
				log.Errorf("delete UpgradeJob for app %s/%s failed: %v", app.Namespace, app.Name, err)
				return err
			}
		} else {
			existedJob[uj.Name] = true
		}
	}

	// 创建新的job
	for _, ds := range daemonsets {
		jobName := getUpgradeJobNameFromApp(ds, app)
		if existedJob[jobName] {
			// job已存在则跳过
			continue
		}
		// create job from policy
		_, err := createUpgradeJobFromPolicy(applicationClient, ds, app)
		if err != nil {
			log.Errorf("create UpgradeJob for app %s/%s failed: %v", app.Namespace, app.Name, err)
			return err
		}
	}

	return nil
}

func createUpgradeJobFromPolicy(applicationClient applicationversionedclient.ApplicationV1Interface, ds string, app *applicationv1.App) (*applicationv1.UpgradeJob, error) {
	//create new job from policy
	up, err := applicationClient.UpgradePolicies().Get(context.TODO(), app.Spec.UpgradePolicy, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	batchNum := *up.Spec.BatchNum
	batchIntervalSeconds := *up.Spec.BatchIntervalSeconds
	maxFailed := *up.Spec.MaxFailed
	maxSurge := *up.Spec.MaxSurge
	uj := &applicationv1.UpgradeJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "UpgradeJob",
			APIVersion: "application.tkestack.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getUpgradeJobNameFromApp(ds, app),
			Namespace: app.Namespace,
			Labels: map[string]string{
				upgradeJobAppNameLabel:       app.Name,
				upgradeJobAppGenerationLabel: strconv.Itoa(int(app.Generation)),
			},
		},
		Spec: applicationv1.UpgradeJobSpec{
			TenantID: app.Spec.TenantID,
			Target:   ds,
			AppRefer: app.Name,

			BatchNum:             &batchNum,
			BatchIntervalSeconds: &batchIntervalSeconds,
			MaxFailed:            &maxFailed,
			MaxSurge:             &maxSurge,
		},
	}

	return applicationClient.UpgradeJobs(app.Namespace).Create(context.TODO(), uj, metav1.CreateOptions{})
}

func getUpgradeJobNameFromApp(ds string, app *applicationv1.App) string {
	// format of ds is like kube-system/kube-proxy, use ds's name
	s := strings.Split(ds, "/")
	return fmt.Sprintf("%s-%s-%d", app.Name, s[1], app.Generation)
}
