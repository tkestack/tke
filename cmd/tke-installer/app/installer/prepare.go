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

package installer

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/template"

	// import platform schema
	_ "tkestack.io/tke/api/platform/install"
)

func (t *TKE) prepareCustomImagesSteps() {
	if !t.Para.Config.Registry.IsOfficial() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Login registry",
				Func: t.loginRegistry,
			},
			{
				Name: "Load images",
				Func: t.loadImages,
			},
			{
				Name: "Load custom K8s images",
				Func: t.loadCustomK8sImages,
			},
			{
				Name: "Build custom provider res image",
				Func: t.buildCustomProviderRes,
			},
			{
				Name: "Add custom version in cluster info",
				Func: t.addCustomK8sVersion,
			},
		}...)
	}

	t.steps = append(t.steps, []types.Handler{}...)

	t.steps = funk.Filter(t.steps, func(step types.Handler) bool {
		return !funk.ContainsString(t.Para.Config.SkipSteps, step.Name)
	}).([]types.Handler)

	t.log.Info("Steps:")
	for i, step := range t.steps {
		t.log.Infof("%d %s", i, step.Name)
	}
}

func (t *TKE) prepareForPrepareCustomImages(ctx context.Context) error {
	err := t.prepareForUpgrade(ctx)
	if err != nil {
		return err
	}
	tkeVersion, _, err := t.getPlatformVersions(ctx)
	if err != nil {
		return err
	}
	if tkeVersion != spec.TKEVersion {
		return errors.Errorf("cannot prepare custom images, platform version %s is not equal to installer version %s", tkeVersion, spec.TKEVersion)
	}
	return nil
}

func (t *TKE) doPrepareCustomImages() {
	ctx := t.log.WithContext(context.Background())
	taskType := "prepare custom images"
	t.prepareCustomImagesSteps()
	t.doSteps(ctx, taskType)
}

func (t *TKE) getCustomPatchVersion(ctx context.Context) (patchVersion string, err error) {
	files, err := ioutil.ReadDir(t.Config.CustomUpgradeResourceDir)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", errors.Errorf("please make sure %s dir has and only has one patch version dir", t.Config.CustomUpgradeResourceDir)
	}

	patchVersion = files[0].Name()
	dirRegexp := regexp.MustCompile(`^\d+.\d+.\d+$`)
	ressutl := dirRegexp.MatchString(patchVersion)
	if !ressutl {
		return "", errors.Errorf("your patch version dir name %s is not a version, please make sure your dir name is a version like 1.18.3", patchVersion)
	}
	return patchVersion, nil
}

func (t *TKE) loadCustomK8sImages(ctx context.Context) error {
	patchVersion, err := t.getCustomPatchVersion(ctx)
	if err != nil {
		return err
	}

	customK8sImagesDir := path.Join(t.Config.CustomUpgradeResourceDir, patchVersion, constants.CustomK8sImageDirName)
	files, err := getFilesFromDir(customK8sImagesDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = t.docker.LoadImages(path.Join(customK8sImagesDir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TKE) buildCustomProviderRes(ctx context.Context) error {
	patchVersion, err := t.getCustomPatchVersion(ctx)
	if err != nil {
		return err
	}

	amdDirExists := false
	armDirExists := false

	customK8sBinaryAmdDir := path.Join(t.Config.CustomUpgradeResourceDir, patchVersion, constants.CustomK8sBinaryAmdDirName)
	amdFiles, err := getFilesFromDir(customK8sBinaryAmdDir)
	if err == nil {
		amdDirExists = true
	} else if !os.IsNotExist(err) {
		return err
	}

	customK8sBinaryArmDir := path.Join(t.Config.CustomUpgradeResourceDir, patchVersion, constants.CustomK8sBinaryArmDirName)
	armFiles, err := getFilesFromDir(customK8sBinaryArmDir)
	if err == nil {
		armDirExists = true
	} else if !os.IsNotExist(err) {
		return err
	}

	if (len(amdFiles) == 0) && (len(armFiles) == 0) {
		return errors.Errorf("There is no file in %s and %s, cannot build custom provider res image", customK8sBinaryAmdDir, constants.CustomK8sBinaryArmDirName)
	}

	tag, err := t.getLatestProviderResTag(ctx)
	if err != nil {
		return err
	}

	values := map[string]interface{}{
		"ProviderResName": images.Get().ProviderRes.Name,
		"Arch":            "amd64",
		"Tag":             tag,
		"HasAmdDir":       amdDirExists,
		"AmdDir":          customK8sBinaryAmdDir,
		"HasArmDir":       armDirExists,
		"ArmDir":          customK8sBinaryArmDir,
	}

	amdDockerfile, err := genCustomProviderResDockerfile(values)
	if err != nil {
		return err
	}
	values["Arch"] = "arm64"
	armDockerfile, err := genCustomProviderResDockerfile(values)
	if err != nil {
		return err
	}

	nameWithoutArch := "tkestack/" + images.Get().ProviderRes.Name

	amdTarget := nameWithoutArch + "-amd64:" + patchVersion
	err = t.docker.BuildImage(amdDockerfile, amdTarget, "linux/amd64")
	if err != nil {
		return err
	}

	armTarget := nameWithoutArch + "-arm64:" + patchVersion
	err = t.docker.BuildImage(armDockerfile, armTarget, "linux/arm64")
	if err != nil {
		return err
	}
	return nil
}

func (t *TKE) addCustomK8sVersion(ctx context.Context) error {
	_, k8sValidVersions, err := t.getPlatformVersions(ctx)
	if err != nil {
		return err
	}
	patchVersion, err := t.getCustomPatchVersion(ctx)
	if err != nil {
		return err
	}
	k8sValidVersions = append(k8sValidVersions, patchVersion)
	versionsByte, err := json.Marshal(k8sValidVersions)
	if err != nil {
		return err
	}
	patchData := map[string]interface{}{
		"data": map[string]interface{}{
			"k8sValidVersions": string(versionsByte),
		},
	}
	return t.patchClusterInfo(ctx, patchData)
}

func (t *TKE) getLatestProviderResTag(ctx context.Context) (tag string, err error) {
	_, k8sValidVersions, err := t.getPlatformVersions(ctx)
	if err != nil {
		return "", err
	}
	if len(k8sValidVersions) > len(spec.K8sValidVersions) {
		return k8sValidVersions[len(k8sValidVersions)-1], nil
	}
	return images.Get().ProviderRes.Tag, nil
}

func genCustomProviderResDockerfile(values map[string]interface{}) ([]byte, error) {
	dockerfileTpl := `FROM tkestack/{{ .ProviderResName }}-{{ .Arch }}:{{ .Tag }}

WORKDIR /data

{{- if .HasAmdDir }}
COPY {{ .AmdDir }}* res/linux-amd64/
{{- end }}

{{- if .HasArmDir }}
COPY {{ .ArmDir }}* res/linux-arm64/
{{- end }}

ENTRYPOINT ["sh"]`
	return template.ParseString(dockerfileTpl, values)
}

func getFilesFromDir(path string) (files []os.FileInfo, err error) {
	files = []os.FileInfo{}
	imageDir, err := os.Stat(path)
	if err != nil {
		return files, err
	}
	if !imageDir.IsDir() {
		return files, errors.Errorf("%s is not a dir", path)
	}
	files, err = ioutil.ReadDir(path)
	return files, err
}
