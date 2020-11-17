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

package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"helm.sh/chartmuseum/pkg/chartmuseum/server/multitenant"
	"k8s.io/apimachinery/pkg/runtime/schema"
	authv1 "tkestack.io/tke/api/auth/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/registry/chartmuseum"
	serveroptionsv1 "tkestack.io/tke/pkg/registry/chartmuseum/serveroptions/v1"
	"tkestack.io/tke/pkg/registry/controller/chart"
	"tkestack.io/tke/pkg/registry/controller/chartgroup"
	"tkestack.io/tke/pkg/registry/controller/identityprovider"
	helm "tkestack.io/tke/pkg/registry/harbor/helmClient"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"
)

const (
	chartGroupSyncPeriod      = 30 * time.Second
	concurrentChartGroupSyncs = 10

	chartSyncPeriod      = 60 * time.Second
	concurrentChartSyncs = 10

	identityProviderSyncPeriod      = 60 * time.Second
	concurrentIdentityProviderSyncs = 10
)

func newHelmClient(ctx ControllerContext) *helm.APIClient {
	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(
		ctx.RegistryConfig.Security.AdminUsername+":"+ctx.RegistryConfig.Security.AdminPassword),
	)
	tr, _ := transport.NewOneWayTLSTransport(ctx.RegistryConfig.HarborCAFile, true)
	helmCfg := &helm.Configuration{
		BasePath:      fmt.Sprintf("https://%s/api", ctx.RegistryConfig.DomainSuffix),
		DefaultHeader: headers,
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		HTTPClient: &http.Client{
			Transport: tr,
		},
	}
	return helm.NewAPIClient(helmCfg)
}

func startChartGroupController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: registryv1.GroupName, Version: "v1", Resource: "chartgroups"}] {
		return nil, false, nil
	}

	var helmClient *helm.APIClient

	if ctx.RegistryConfig.HarborEnabled {
		helmClient = newHelmClient(ctx)
	}

	ctrl := chartgroup.NewController(
		ctx.BusinessClient,
		ctx.ClientBuilder.ClientOrDie("chartgroup-controller"),
		ctx.InformerFactory.Registry().V1().ChartGroups(),
		chartGroupSyncPeriod,
		registryv1.ChartGroupFinalize,
		helmClient,
	)

	go ctrl.Run(concurrentChartGroupSyncs, ctx.Stop)

	return nil, true, nil
}

func startChartController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: registryv1.GroupName, Version: "v1", Resource: "charts"}] {
		return nil, false, nil
	}

	multiTenantServerOptions, err := serveroptionsv1.BuildChartMuseumConfig(ctx.RegistryConfig, chartmuseum.PathPrefix, chartmuseum.MaxUploadSize)
	if err != nil {
		log.Error("Failed to initialize chartmuseum server configuration", log.Err(err))
		return nil, false, err
	}
	multiTenantServer, err := multitenant.NewMultiTenantServer(*multiTenantServerOptions)
	if err != nil {
		log.Error("Failed to create chartmuseum server", log.Err(err))
		return nil, false, err
	}

	var helmClient *helm.APIClient

	if ctx.RegistryConfig.HarborEnabled {
		helmClient = newHelmClient(ctx)
	}

	ctrl := chart.NewController(
		ctx.ClientBuilder.ClientOrDie("chart-controller"),
		ctx.InformerFactory.Registry().V1().Charts(),
		chartSyncPeriod,
		registryv1.ChartFinalize,
		multiTenantServer,
		helmClient,
	)

	go ctrl.Run(concurrentChartSyncs, ctx.Stop)

	return nil, true, nil
}

func startIdentityProviderController(ctx ControllerContext) (http.Handler, bool, error) {
	if ctx.AuthClient == nil {
		return nil, false, nil
	}

	if !ctx.AuthAvailableResources[schema.GroupVersionResource{Group: authv1.GroupName, Version: "v1", Resource: "identityproviders"}] {
		return nil, false, nil
	}

	ctrl := identityprovider.NewController(
		ctx.AuthClient,
		ctx.ClientBuilder.ClientOrDie("identityprovider-controller"),
		ctx.AuthInformerFactory.Auth().V1().IdentityProviders(),
		identityProviderSyncPeriod,
		ctx.RegistryDefaultConfiguration,
		ctx.RegistryConfig,
	)

	go ctrl.Run(concurrentIdentityProviderSyncs, ctx.Stop)

	return nil, true, nil
}
