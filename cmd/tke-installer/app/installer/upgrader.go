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
	"fmt"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	platformv1 "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	configv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/file"

	// import platform schema
	_ "tkestack.io/tke/api/platform/install"
)

const (
	registryCmName = "tke-registry-api"
	registryCmKey  = "tke-registry-config.yaml"
)

func (t *TKE) upgradeSteps() {
	if !t.Para.Config.Registry.IsOfficial() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Login registry",
				Func: t.loginRegistryForUpgrade,
			},
			{
				Name: "Load images",
				Func: t.loadImages,
			},
			{
				Name: "Push images",
				Func: t.pushImages,
			},
		}...)
	}

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Update tke-platform-api",
			Func: t.updateTKEPlatformAPI,
		},
		{
			Name: "Update tke-platform-controller",
			Func: t.updateTKEPlatformController,
		},
	}...)

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Patch platform versions in cluster info",
			Func: t.patchPlatformVersion,
		},
	}...)

	t.steps = funk.Filter(t.steps, func(step types.Handler) bool {
		return !funk.ContainsString(t.Para.Config.SkipSteps, step.Name)
	}).([]types.Handler)

	t.log.Info("Steps:")
	for i, step := range t.steps {
		t.log.Infof("%d %s", i, step.Name)
	}
}

func (t *TKE) updateTKEPlatformAPI(ctx context.Context) error {
	com := "tke-platform-api"
	depl, err := t.globalClient.AppsV1().Deployments(t.namespace).Get(ctx, com, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(depl.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("%s has no containers", com)
	}
	depl.Spec.Template.Spec.Containers[0].Image = images.Get().TKEPlatformAPI.FullName()

	_, err = t.globalClient.AppsV1().Deployments(t.namespace).Update(ctx, depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, com)
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) updateTKEPlatformController(ctx context.Context) error {
	com := "tke-platform-controller"
	depl, err := t.globalClient.AppsV1().Deployments(t.namespace).Get(ctx, com, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if len(depl.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("%s has no containers", com)
	}
	depl.Spec.Template.Spec.Containers[0].Image = images.Get().TKEPlatformController.FullName()

	if len(depl.Spec.Template.Spec.InitContainers) == 0 {
		return fmt.Errorf("%s has no initContainers", com)
	}
	depl.Spec.Template.Spec.InitContainers[0].Image = images.Get().ProviderRes.FullName()

	_, err = t.globalClient.AppsV1().Deployments(t.namespace).Update(ctx, depl, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, com)
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) prepareForUpgrade(ctx context.Context) error {
	t.namespace = namespace

	_ = t.loadTKEData()

	if !file.Exists(t.Config.Kubeconfig) || !file.IsFile(t.Config.Kubeconfig) {
		return fmt.Errorf("kubeconfig %s doesn't exist", t.Config.Kubeconfig)
	}
	config, err := clientcmd.BuildConfigFromFlags("", t.Config.Kubeconfig)
	if err != nil {
		return err
	}
	t.globalClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	platformClient, err := platformv1.NewForConfig(config)
	if err != nil {
		return err
	}
	t.Cluster, err = typesv1.GetClusterByName(ctx, platformClient, "global")
	if err != nil {
		return err
	}
	t.Para.Cluster = *t.Cluster.Cluster
	t.Para.Config.Registry.UserInputRegistry.Domain = t.Config.RegistryDomain
	t.Para.Config.Registry.UserInputRegistry.Username = t.Config.RegistryUsername
	t.Para.Config.Registry.UserInputRegistry.Password = []byte(t.Config.RegistryPassword)
	t.Para.Config.Registry.UserInputRegistry.Namespace = t.Config.RegistryNamespace
	err = t.loadRegistry(ctx)
	if err != nil {
		if apierrors.IsNotFound(err) {
			t.log.Infof("Not found %s", registryCmName)
			if t.Para.Config.Registry.ThirdPartyRegistry == nil {
				t.log.Infof("Not found third party registry")
				t.Para.Config.Registry.ThirdPartyRegistry = &types.ThirdPartyRegistry{}
			}
		} else {
			return err
		}
	}
	return nil
}

func (t *TKE) loadRegistry(ctx context.Context) error {
	registryCm, err := t.globalClient.CoreV1().ConfigMaps(t.namespace).Get(ctx, registryCmName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	registryConfig := &configv1.RegistryConfiguration{}
	err = yaml.Unmarshal([]byte(registryCm.Data[registryCmKey]), registryConfig)
	if err != nil {
		return err
	}
	t.Para.Config.Registry.TKERegistry = &types.TKERegistry{
		Domain:        registryConfig.DomainSuffix,
		HarborEnabled: registryConfig.HarborEnabled,
		HarborCAFile:  registryConfig.HarborCAFile,
		Namespace:     "library",
		Username:      registryConfig.Security.AdminUsername,
		Password:      []byte(registryConfig.Security.AdminPassword),
	}
	return nil
}

func (t *TKE) loginRegistryForUpgrade(ctx context.Context) error {
	containerregistry.Init(t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace())
	cmd := exec.Command("docker", "login",
		"--username", t.Para.Config.Registry.Username(),
		"--password", string(t.Para.Config.Registry.Password()),
		t.Para.Config.Registry.Domain(),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return errors.New(string(out))
		}
		return err
	}
	return nil
}
