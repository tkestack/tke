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

package expansion

import (
	"context"
	"fmt"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"os"
	"path"
	"strings"
	appv1 "tkestack.io/tke/api/application/v1"
	applicationv1client "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/docker"
)

// something cannot import from tke lib
const TKENamespace = "tke"
const TKEPlatformExpansionConfigmapName = "platform-expansion"
const TKEPlatformExpansionFilesConfigmapName = "platform-expansion-files"
const TKEPlatformRootDir = "/app/"
const TKEPlatformDataDir = "data/"
const TKEPlatformBase = TKEPlatformRootDir + TKEPlatformDataDir
const TKEPlatformExpansionBase = TKEPlatformBase + "expansions/"
const TKEPlatformExpansionBaseMountPath = TKEPlatformBase + "expansions"
const TKEPlatformExpansionFilesPath = TKEPlatformExpansionBase + "files_generated/"
const TKEPlatformExpansionFilesMountPath = TKEPlatformExpansionBase + "files_generated"
const InstallerProviderBasePath = "provider/"
const InstallerProviderBareMetalPath = InstallerProviderBasePath + "baremetal/"
const InstallerHooksDir = "hooks/"

const KubeconfigFileBaseName = "admin.kubeconfig"

const InstallerPreInstallHookName = "pre-install"
const InstallerPostClusterReadyHookName = "post-cluster-ready"
const InstallerPostInstallHookName = "post-install"

// MergeProvider overrides provider files in TKEStack installer by expansion specifying
func (d *ExpansionDriver) MergeProvider() error {
	if !d.enableProvider() {
		return nil
	}
	for _, f := range d.Provider {
		src := expansionProviderPath + f
		dst := InstallerProviderBasePath + f
		err := copyFile(src, dst, 0)
		if err != nil {
			return fmt.Errorf("merge provider failed, copy %v to %v, %v", src, dst, err)
		}
	}
	return nil
}

// MergeInstallerSkipSteps sets up install skip steps of TKEStack-installer by expansion specifying
func (d *ExpansionDriver) MergeInstallerSkipSteps(skipSteps []string) []string {
	// merge skip steps
	if !d.enableSkipSteps() {
		return skipSteps
	}
	if skipSteps == nil {
		skipSteps = make([]string, 0)
	}
	for _, step := range d.InstallerSkipSteps {
		if !funk.ContainsString(skipSteps, step) {
			skipSteps = append(skipSteps, step)
		}
	}
	return skipSteps
}

// MergeCluster
// 1. puts copy-files to files_generated directory and register them with TKEStack hook config.
// 2. sets up create-cluster skip steps by expansion specifying
// 3. sets up create-cluster delegation steps by expansion specifying
// 4. passes through kubernetes args from expansion specifying to TKEStack config.
func (d *ExpansionDriver) MergeCluster(cluster *v1.Cluster) {
	// merge file hook config
	if d.enableFiles() {
		if len(cluster.Spec.Features.Hooks) == 0 {
			cluster.Spec.Features.Hooks = make(map[platformv1.HookType]string)
		}
		for _, f := range d.Files {
			ff := d.toFlatPath(f)
			src, dst := expansionFilesGeneratedPath+ff, absolutePath+f
			cluster.Spec.Features.Files = append(cluster.Spec.Features.Files, platformv1.File{
				Src: src,
				Dst: dst,
			})
			hookType, ok := d.isHookScript(f)
			if ok {
				cluster.Spec.Features.Hooks[hookType] = dst
			}
		}
	}
	// merge skip conditions
	if d.enableSkipConditions() {
		if len(cluster.Spec.Features.SkipConditions) == 0 {
			cluster.Spec.Features.SkipConditions = make([]string, 0)
		}
		for _, skipC := range d.CreateClusterSkipConditions {
			if !funk.ContainsString(cluster.Spec.Features.SkipConditions, skipC) {
				cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, skipC)
			}
		}
	}

	// merge delegated conditions
	if d.enableDelegateConditions() {
		if len(cluster.Spec.Features.DelegateConditions) == 0 {
			cluster.Spec.Features.DelegateConditions = make([]string, 0)
		}
		for _, skipC := range d.CreateClusterDelegateConditions {
			if !funk.ContainsString(cluster.Spec.Features.DelegateConditions, skipC) {
				cluster.Spec.Features.DelegateConditions = append(cluster.Spec.Features.DelegateConditions, skipC)
			}
		}
	}

	// merge extra args
	mergeMap(&cluster.Spec.APIServerExtraArgs, &d.CreateClusterExtraArgs.ApiServerExtraArgs)
	d.log.Infof("APIServerExtraArgs %+v", cluster.Spec.APIServerExtraArgs)
	d.log.Infof("d.APIServerExtraArgs %+v", d.CreateClusterExtraArgs.ApiServerExtraArgs)
	mergeMap(&cluster.Spec.ControllerManagerExtraArgs, &d.CreateClusterExtraArgs.ControllerManagerExtraArgs)
	mergeMap(&cluster.Spec.SchedulerExtraArgs, &d.CreateClusterExtraArgs.SchedulerExtraArgs)
	mergeMap(&cluster.Spec.KubeletExtraArgs, &d.CreateClusterExtraArgs.KubeletExtraArgs)
	mergeMap(&cluster.Spec.DockerExtraArgs, &d.CreateClusterExtraArgs.DockerExtraArgs)
	if d.CreateClusterExtraArgs.Etcd != nil {
		cluster.Spec.Etcd = d.CreateClusterExtraArgs.Etcd
	}
}

// PatchPlatformWithExpansion forms configmaps up and mount them to TKEStack-platform deployments
func (d *ExpansionDriver) PatchPlatformWithExpansion(ctx context.Context, client kubernetes.Interface, app string) error {
	var supportedApps = map[string]bool{
		"tke-platform-api":        true,
		"tke-platform-controller": true,
	}
	if !supportedApps[app] {
		return fmt.Errorf("PatchPlatformWithExpansion do not support app %v", app)
	}

	expansionCmNN, err := d.createExpansionConfigmap(ctx, client)
	if err != nil {
		return err
	}
	expansionFilesCnNN, err := d.createExpansionFilesConfigmap(ctx, client)
	if err != nil {
		return err
	}
	err = apiclient.PatchDeployWithConfigmap(ctx, client, app, expansionCmNN, TKEPlatformExpansionBaseMountPath)
	if err != nil {
		return err
	}
	err = apiclient.PatchDeployWithConfigmap(ctx, client, app, expansionFilesCnNN, TKEPlatformExpansionFilesMountPath)
	if err != nil {
		return err
	}
	return nil
}

// MergeExpansionImages merges expansion image list into tkeImages, and then let TKEStack-installer tag them
func (d *ExpansionDriver) MergeExpansionImages(tkeImages *[]string) {
	if !d.enableImages() {
		return
	}
	for _, image := range d.Images {
		*tkeImages = append(*tkeImages, image)
	}
}

// CopyChartsToDst copies expansion charts to a dstDir and let TKEStack-installer upload them
func (d *ExpansionDriver) CopyChartsToDst(group string, dstDir string) error {
	if !d.enableCharts() {
		return nil
	}
	for _, chart := range d.Charts {
		if !strings.HasPrefix(chart, fmt.Sprintf("%s/", group)) {
			continue
		}
		src := fmt.Sprintf("%s%s", expansionChartsPath, chart)
		err := copyFileToDir(src, dstDir, 0644)
		if err != nil {
			return fmt.Errorf("copy file to dir error %v, %v, %v", src, dstDir, err)
		}
	}
	return nil
}

// LoadOperatorImage loads expansion operator image from expansion package
func (d *ExpansionDriver) LoadOperatorImage(ctx context.Context) error {

	// TODO:
	d.log.Errorf("mocked! loadOperatorImage not implement")
	return nil
	// return fmt.Errorf("loadOperatorImage not implement")
}

// PatchHookFiles copies installer hook files from expansion to TKEStack-installer
func (d *ExpansionDriver) PatchHookFiles(ctx context.Context) error {
	if !d.enableHooks() {
		d.log.Info("expansion hooks disabled")
		return nil
	}
	var wellKnownHookFiles = map[string]bool{
		InstallerPreInstallHookName:       true,
		InstallerPostClusterReadyHookName: true,
		InstallerPostInstallHookName:      true,
	}
	for _, hook := range d.Hooks {
		if !wellKnownHookFiles[hook] {
			d.log.Errorf("installer hook is not recognizable %v ", hook)
			continue
		}
		src := expansionHooksPath + hook
		dst := InstallerHooksDir + hook
		err := copyFile(src, dst, 0755)
		//input, err := ioutil.ReadFile(src)
		if err != nil {
			d.log.Errorf("copy hook file failed %v %v %v", src, dst, err)
			return err
		}
	}

	return nil
}

// StartOperator
func (d *ExpansionDriver) StartOperator(ctx context.Context) error {
	if !d.enableOperator() {
		return nil
	}
	return fmt.Errorf("startOperator not implement")
}

// LoadExpansionImages loads expansion images into local docker daemon
func (d *ExpansionDriver) LoadExpansionImages(ctx context.Context, docker *docker.Docker) error {

	if !d.enableImages() {
		return nil
	}
	// TODO: check image list
	if _, err := os.Stat(expansionImagesPath); err != nil {
		return err
	}
	if err := docker.LoadImages(expansionImagesPath); err != nil {
		return err
	}
	return nil
}

// WriteKubeconfigFile sends a copy of kubeconfig for sharing to expansion operator
func (d *ExpansionDriver) WriteKubeconfigFile(ctx context.Context, data []byte) error {
	return ioutil.WriteFile(expansionConfPath+"/"+KubeconfigFileBaseName, data, 0644)
}

// InstallApplications installs all application crs into global cluster
func (d *ExpansionDriver) InstallApplications(ctx context.Context, applicationClient applicationv1client.ApplicationV1Interface, tkeValues map[string]string) error {
	if !d.enableApplications() {
		return nil
	}
	for _, app := range d.Applications {
		file := expansionApplicationPath + app + fileSuffixYaml
		application := &appv1.App{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app,
				Namespace: d.ExpansionName,
			},
		}
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read application file failed %v, %v", file, err)
		}
		err = yaml.Unmarshal(b, application)
		if err != nil {
			return fmt.Errorf("unmarshal application file failed %v, %v", file, err)
		}
		for k, v := range d.Values {
			tkeValues[k] = v
		}
		vb, err := yaml.Marshal(tkeValues)
		if err != nil {
			return fmt.Errorf("marshal tkevalues failed %+v, %v", tkeValues, err)
		}
		application.Spec.Values.RawValues = string(vb)
		ns := application.Namespace
		_, err = applicationClient.Apps(ns).Create(ctx, application, metav1.CreateOptions{})
		if err != nil {
			// TODO: workaround
			if strings.Contains(err.Error(), "spec.name: Duplicate value") {
				d.log.Infof("duplicate app, skipping %v, %v", app, err)
				continue
			}
			return fmt.Errorf("create application failed %+v, %v", application, err)
		}
		d.log.Infof("application created %v", app)
	}

	return nil
}

func (d *ExpansionDriver) EnableImages() bool {
	return d.enableImages()
}

func (d *ExpansionDriver) HasNewK8sVersion() bool {
	// TODO: verify version
	return d.K8sVersion != ""
}

func (d *ExpansionDriver) NewK8sVersion() (string, error) {
	// TODO: verify version
	return d.K8sVersion, nil
}

func (d *ExpansionDriver) createExpansionConfigmap(ctx context.Context, client kubernetes.Interface) (*types.NamespacedName, error) {
	files := make([]string, 0)

	for _, f := range []string{
		defaultExpansionConfName,
		expansionValuesFileName,
	} {
		if _, err := os.Stat(defaultExpansionBase + f); err == nil {
			files = append(files, f)
		}
	}

	nn := &types.NamespacedName{
		Namespace: TKENamespace,
		Name:      TKEPlatformExpansionConfigmapName,
	}
	err := apiclient.GenerateConfigmapFromFiles(ctx, client, nn, defaultExpansionBase, files, false)
	if err != nil {
		return nn, err
	}
	return nn, nil
}

func (d *ExpansionDriver) createExpansionFilesConfigmap(ctx context.Context, client kubernetes.Interface) (*types.NamespacedName, error) {
	files := make([]string, 0)
	if d.enableFiles() {
		for _, f := range d.Files {
			ff := d.toFlatPath(f)
			files = append(files, ff)
		}
	}
	nn := &types.NamespacedName{
		Namespace: TKENamespace,
		Name:      TKEPlatformExpansionFilesConfigmapName,
	}
	err := apiclient.GenerateConfigmapFromFiles(ctx, client, nn, expansionFilesGeneratedPath, files, false)
	if err != nil {
		return nn, err
	}
	return nn, nil
}

func (d *ExpansionDriver) makeFlatFiles() error {
	if !d.enableFiles() {
		return nil
	}
	for _, f := range d.Files {
		src := expansionFilesPath + f
		ff := d.toFlatPath(f)
		flatFile := expansionFilesGeneratedPath + ff
		_, err := os.Stat(flatFile)
		if err == nil {
			continue
		}
		if os.IsNotExist(err) {
			err = copyFile(src, flatFile, 0)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("stat flatFile system error %v %v", flatFile, err)
		}
	}
	return nil
}

func (d *ExpansionDriver) toFlatPath(f string) string {
	return strings.Replace(f, string(os.PathSeparator), expansionFilePathSeparator, -1)
}

func (d *ExpansionDriver) isHookScript(fp string) (platformv1.HookType, bool) {
	var createClusterHookFileTypes = map[platformv1.HookType]bool{
		platformv1.HookPreInstall:         true,
		platformv1.HookPostInstall:        true,
		platformv1.HookPreClusterInstall:  true,
		platformv1.HookPostClusterInstall: true,
	}
	hookType := platformv1.HookType(path.Base(fp))
	if _, ok := createClusterHookFileTypes[hookType]; ok {
		return hookType, true
	}
	return "", false
}
