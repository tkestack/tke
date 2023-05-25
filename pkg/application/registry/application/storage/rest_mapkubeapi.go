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

package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/api/application"
	applicationv1 "tkestack.io/tke/api/application/v1"
	v1 "tkestack.io/tke/api/application/v1"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	"tkestack.io/tke/pkg/application/util"
	"tkestack.io/tke/pkg/util/log"
)

// MapKubeApiREST adapts a service registry into apiserver's RESTStorage model.
type MapKubeApiREST struct {
	store             ApplicationStorage
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV1Interface
}

// NewMapKubeApiREST returns a wrapper around the underlying generic storage and performs mapkubeapi of helm releases.
func NewMapKubeApiREST(
	store ApplicationStorage,
	applicationClient *applicationinternalclient.ApplicationClient,
	platformClient platformversionedclient.PlatformV1Interface,
) *MapKubeApiREST {
	rest := &MapKubeApiREST{
		store:             store,
		applicationClient: applicationClient,
		platformClient:    platformClient,
	}
	return rest
}

// New creates a new chart proxy options object
func (m *MapKubeApiREST) New() runtime.Object {
	return &applicationv1.RollbackProxyOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (m *MapKubeApiREST) ConnectMethods() []string {
	return []string{"GET"}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (m *MapKubeApiREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := m.store.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, k8serrors.NewInternalError(err)
	}
	app := obj.(*application.App)
	logger := log.WithName(fmt.Sprintf("%s/%s", app.Spec.TargetCluster, app.Spec.Name))

	appv1 := &v1.App{}
	if err = v1.Convert_application_App_To_v1_App(app, appv1, nil); err != nil {
		return nil, err
	}

	cfg, err := util.NewActionConfigWithProvider(ctx, m.platformClient, appv1)
	if err != nil {
		logger.Errorf("failed to new action config, err:%s", err.Error())
		return nil, errors.Wrap(err, "failed to get Helm action configuration")
	}

	// get the Kubernetes server version
	kubeVersionStr, err := getKubernetesServerVersion(ctx, m.platformClient, appv1)
	if err != nil {
		logger.Errorf("failed to get k8s version, err:%s", err.Error())
		return nil, err
	}

	var releaseName = appv1.Spec.Name
	releaseToMap, err := cfg.Releases.Last(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			logger.Info("release npt found, don't need to update release")
			return appv1, nil
		}
		logger.Errorf("failed to get release latest version, err:%s", err.Error())
		return nil, errors.Wrapf(err, "failed to get release '%s' latest version", releaseName)
	}
	modifiedManifest, err := ReplaceManifestUnSupportedAPIs(releaseToMap.Manifest, kubeVersionStr)
	if err != nil {
		logger.Errorf("replace manifest Unsupported APIs failed, err:%s", err.Error())
		return nil, err
	}
	if modifiedManifest == releaseToMap.Manifest {
		logger.Info("old release: %s", releaseToMap.Manifest)
		logger.Info("new release: %s", modifiedManifest)
		logger.Info("don't need to update release")
		return appv1, nil
	}

	if err = updateRelease(releaseToMap, modifiedManifest, cfg); err != nil {
		logger.Errorf("failed to update release, err:%s", err.Error())
		return nil, errors.Wrapf(err, "failed to update release '%s'", releaseName)
	}
	return appv1, nil
}

func getReleaseVersionName(rel *release.Release) string {
	return fmt.Sprintf("%s.v%d", rel.Name, rel.Version)
}

func updateRelease(origRelease *release.Release, modifiedManifest string, cfg *action.Configuration) error {
	log.Infof("Set status of release version '%s' to 'superseded'", getReleaseVersionName(origRelease))
	origRelease.Info.Status = release.StatusSuperseded
	if err := cfg.Releases.Update(origRelease); err != nil {
		return errors.Wrapf(err, "failed to update release version '%s': %s", getReleaseVersionName(origRelease))
	}
	log.Infof("Release version '%s' updated successfully", getReleaseVersionName(origRelease))

	var newRelease = origRelease
	newRelease.Manifest = modifiedManifest
	newRelease.Info.Description = "Kubernetes deprecated API upgrade - DO NOT rollback from this version"
	newRelease.Info.LastDeployed = cfg.Now()
	newRelease.Version = origRelease.Version + 1
	newRelease.Info.Status = release.StatusDeployed
	log.Infof("Add release version '%s' with updated supported APIs", getReleaseVersionName(origRelease))
	if err := cfg.Releases.Create(newRelease); err != nil {
		return errors.Wrapf(err, "failed to create new release version '%s': %s", getReleaseVersionName(origRelease))
	}
	log.Infof("Release version '%s' added successfully", getReleaseVersionName(origRelease))
	return nil
}

// TODO：mapMetadata赋值
var MetadataMappings = []Mapping{
	{
		DeprecatedAPI:    "apiVersion: extensions/v1beta1 \nkind: Deployment",
		NewAPI:           "apiVersion: apps/v1 \nkind: Deployment",
		RemovedInVersion: "v1.16",
	},
	{
		DeprecatedAPI:    "kind: Deployment \napiVersion: extensions/v1beta1",
		NewAPI:           "apiVersion: apps/v1 \nkind: Deployment",
		RemovedInVersion: "v1.16",
	},

	{
		DeprecatedAPI:    "apiVersion: apiextensions.k8s.io/v1beta1 \nkind: CustomResourceDefinition",
		NewAPI:           "apiVersion: apiextensions.k8s.io/v1 \nkind: CustomResourceDefinition",
		RemovedInVersion: "v1.22",
	},
	{
		DeprecatedAPI:    "kind: CustomResourceDefinition \napiVersion: apiextensions.k8s.io/v1beta1",
		NewAPI:           "apiVersion: apiextensions.k8s.io/v1 \nkind: CustomResourceDefinition",
		RemovedInVersion: "v1.22",
	},

	{
		DeprecatedAPI:    "apiVersion: storage.k8s.io/v1beta1 \nkind: CSIDriver",
		NewAPI:           "apiVersion: storage.k8s.io/v1 \nkind: CSIDriver",
		RemovedInVersion: "v1.19",
	},
	{
		DeprecatedAPI:    "kind: CSIDriver \napiVersion: storage.k8s.io/v1beta1",
		NewAPI:           "apiVersion: storage.k8s.io/v1 \nkind: CSIDriver",
		RemovedInVersion: "v1.19",
	},
}

type Mapping struct {
	DeprecatedAPI    string
	NewAPI           string
	RemovedInVersion string
}

func ReplaceManifestUnSupportedAPIs(modifiedManifest, kubeVersionStr string) (string, error) {
	log.Infof("kubeVersionStr %s", kubeVersionStr)
	for _, mapping := range MetadataMappings {
		if count := strings.Count(modifiedManifest, mapping.DeprecatedAPI); count > 0 {
			// 目标集群版本大于废弃k8s api的版本
			if semver.Compare(kubeVersionStr, mapping.RemovedInVersion) > 0 {
				modifiedManifest = strings.ReplaceAll(modifiedManifest, mapping.DeprecatedAPI, mapping.NewAPI)
			}
		}
	}

	return modifiedManifest, nil
}

func getKubernetesServerVersion(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, app *applicationv1.App) (string, error) {
	provider, err := applicationprovider.GetProvider(app)
	if err != nil {
		return "", err
	}
	restConfig, err := provider.GetRestConfig(ctx, platformClient, app)
	if err != nil {
		return "", err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if clientSet == nil {
		return "", errors.Errorf("kubernetes cluster unreachable")
	}
	kubeVersion, err := clientSet.ServerVersion()
	if err != nil {
		return "", errors.Wrap(err, "kubernetes cluster unreachable")
	}
	if !semver.IsValid(kubeVersion.GitVersion) {
		return "", errors.Errorf("Failed to get Kubernetes server version")
	}
	return kubeVersion.GitVersion, nil
}
