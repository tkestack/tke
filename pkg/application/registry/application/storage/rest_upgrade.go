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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/application"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	appconfig "tkestack.io/tke/pkg/application/config"
	applicationstrategy "tkestack.io/tke/pkg/application/registry/application"
)

// CanUpgradeREST adapts a service registry into apiserver's RESTStorage model.
type CanUpgradeREST struct {
	store             ApplicationStorage
	applicationClient applicationversionedclient.ApplicationV1Interface
	platformClient    platformversionedclient.PlatformV1Interface
	repo              appconfig.RepoConfiguration
}

// NewCanUpgradeREST returns a wrapper around the underlying generic storage and performs mapkubeapi of helm releases.
func NewCanUpgradeREST(
	store ApplicationStorage,
	applicationClient applicationversionedclient.ApplicationV1Interface,
	platformClient platformversionedclient.PlatformV1Interface,
	repo appconfig.RepoConfiguration,
) *CanUpgradeREST {
	rest := &CanUpgradeREST{
		store:             store,
		applicationClient: applicationClient,
		platformClient:    platformClient,
		repo:              repo,
	}
	return rest
}

// New creates a new chart upgrade options object
func (m *CanUpgradeREST) New() runtime.Object {
	return &applicationv1.AppUpgradeOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (m *CanUpgradeREST) ConnectMethods() []string {
	return []string{"GET", "POST"}
}

// NewConnectOptions returns versioned resource that represents upgrade parameters
func (m *CanUpgradeREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &applicationv1.AppUpgradeOptions{}, false, ""
}

// Connect returns a handler for the chart upgrade
func (m *CanUpgradeREST) Connect(ctx context.Context, appName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := m.store.Get(ctx, appName, &metav1.GetOptions{})
	if err != nil {
		return nil, k8serrors.NewInternalError(err)
	}
	app := obj.(*application.App)
	upgradeOpts := opts.(*applicationv1.AppUpgradeOptions)

	appv1 := &applicationv1.App{}
	if err = applicationv1.Convert_application_App_To_v1_App(app, appv1, nil); err != nil {
		return nil, err
	}
	return &upgradeHandler{
		ctx:               ctx,
		app:               appv1,
		ops:               *upgradeOpts,
		applicationClient: m.applicationClient,
		platformClient:    m.platformClient,
	}, nil
}

type upgradeHandler struct {
	ctx               context.Context
	app               *applicationv1.App
	ops               applicationv1.AppUpgradeOptions
	repo              appconfig.RepoConfiguration
	applicationClient applicationversionedclient.ApplicationV1Interface
	platformClient    platformversionedclient.PlatformV1Interface
}

func (h *upgradeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ops := h.ops
	if strings.ToLower(req.Method) == "post" {
		data, err := ioutil.ReadAll(req.Body)
		if err == nil {
			json.Unmarshal(data, &ops)
		}
	}
	hook := applicationstrategy.GetHooks(h.app)
	result, err := hook.CanUpgrade(h.ctx, h.applicationClient, h.platformClient, h.app, h.repo, ops)

	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, k8serrors.NewInternalError(err), w)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, result, w)
}
