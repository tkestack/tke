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
	"net/http"
	"net/url"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/rest"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/registry"
	registryv1 "tkestack.io/tke/api/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/application/util"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/util/log"
)

// InfoREST adapts a service registry into apiserver's RESTStorage model.
type InfoREST struct {
	store          ChartStorage
	platformClient platformversionedclient.PlatformV1Interface
	registryConfig *registryconfig.RegistryConfiguration
	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string
}

// NewInfoREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various chart.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewInfoREST(
	store ChartStorage,
	platformClient platformversionedclient.PlatformV1Interface,
	registryConfig *registryconfig.RegistryConfiguration,
	externalScheme string,
	externalHost string,
	externalPort int,
	externalCAFile string,
) *InfoREST {
	rest := &InfoREST{
		store:          store,
		platformClient: platformClient,
		registryConfig: registryConfig,
		externalScheme: externalScheme,
		externalHost:   externalHost,
		externalPort:   externalPort,
		externalCAFile: externalCAFile,
	}
	return rest
}

// New creates a new chart proxy options object
func (r *InfoREST) New() runtime.Object {
	return &registry.ChartProxyOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *InfoREST) ConnectMethods() []string {
	return []string{"GET"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *InfoREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &registry.ChartProxyOptions{}, false, ""
}

// Connect returns a handler for the chart proxy
func (r *InfoREST) Connect(ctx context.Context, chartName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := r.store.Get(ctx, chartName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	chart := obj.(*registry.Chart)
	proxyOpts := opts.(*registry.ChartProxyOptions)

	if proxyOpts.Version == "" {
		return nil, errors.NewBadRequest("version is required")
	}
	if proxyOpts.Cluster == "" {
		log.Warn("cluster is empty but required, using default cluster: global")
		// return nil, errors.NewBadRequest("cluster is required")
		proxyOpts.Cluster = "global"
	}
	if proxyOpts.Namespace == "" {
		log.Warn("namespace is empty but required, using default cluster: default")
		// return nil, errors.NewBadRequest("default is required")
		proxyOpts.Namespace = "default"
	}

	return &proxyHandler{
		chart:          chart,
		chartVersion:   proxyOpts.Version,
		cluster:        proxyOpts.Cluster,
		namespace:      proxyOpts.Namespace,
		externalScheme: r.externalScheme,
		externalHost:   r.externalHost,
		externalPort:   r.externalPort,
		externalCAFile: r.externalCAFile,
		platformClient: r.platformClient,
		registryConfig: r.registryConfig,
	}, nil
}

type proxyHandler struct {
	chart          *registry.Chart
	chartVersion   string
	cluster        string
	namespace      string
	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string
	platformClient platformversionedclient.PlatformV1Interface
	registryConfig *registryconfig.RegistryConfiguration
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	client, err := util.NewHelmClient(req.Context(), h.platformClient, h.cluster, h.namespace)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	host := h.externalHost
	if h.externalPort > 0 {
		host = host + ":" + strconv.Itoa(h.externalPort)
	}
	url := &url.URL{
		Scheme: h.externalScheme,
		Host:   host,
	}
	cpopt := helmaction.ChartPathOptions{
		CaFile:    h.externalCAFile,
		Username:  h.registryConfig.Security.AdminUsername,
		Password:  h.registryConfig.Security.AdminPassword,
		RepoURL:   url.String() + "/chart/" + h.chart.Spec.ChartGroupName,
		ChartRepo: h.chart.Spec.TenantID + "/" + h.chart.Spec.ChartGroupName,
		Chart:     h.chart.Spec.Name,
		Version:   h.chartVersion,
	}
	destfile, err := client.Pull(&helmaction.PullOptions{
		ChartPathOptions: cpopt,
	})
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	cpopt.ExistedFile = destfile
	show, err := client.Show(&helmaction.ShowOptions{
		ChartPathOptions: cpopt,
	})
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	files := make(map[string]string)
	if show.Chart != nil {
		for _, v := range show.Chart.Raw {
			files[v.Name] = string(v.Data)
		}
	}
	chartInfo := &registryv1.ChartInfo{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: h.chart.Namespace,
			Name:      h.chart.Name,
		},
		Spec: registryv1.ChartInfoSpec{
			Values:   show.Values,
			Readme:   show.Readme,
			RawFiles: files,
		},
	}
	responsewriters.WriteRawJSON(http.StatusOK, chartInfo, w)
}
