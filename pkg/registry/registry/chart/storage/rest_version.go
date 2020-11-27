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
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/registry"
	registryv1 "tkestack.io/tke/api/registry/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	applicationutil "tkestack.io/tke/pkg/application/util"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	"tkestack.io/tke/pkg/registry/config"
	harborHandler "tkestack.io/tke/pkg/registry/harbor/handler"
	helm "tkestack.io/tke/pkg/registry/harbor/helmClient"
	registryutil "tkestack.io/tke/pkg/registry/util"
	authorizationutil "tkestack.io/tke/pkg/registry/util/authorization"
	"tkestack.io/tke/pkg/registry/util/chartpath"
	"tkestack.io/tke/pkg/registry/util/sort"
	"tkestack.io/tke/pkg/util/log"
)

// VersionREST adapts a service registry into apiserver's RESTStorage model.
type VersionREST struct {
	store          ChartStorage
	platformClient platformversionedclient.PlatformV1Interface
	registryClient *registryinternalclient.RegistryClient
	registryConfig *registryconfig.RegistryConfiguration
	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string
	authorizer     authorizer.Authorizer
	helmClient     *helm.APIClient
}

// NewVersionREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various chart.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewVersionREST(
	store ChartStorage,
	platformClient platformversionedclient.PlatformV1Interface,
	registryClient *registryinternalclient.RegistryClient,
	registryConfig *registryconfig.RegistryConfiguration,
	externalScheme string,
	externalHost string,
	externalPort int,
	externalCAFile string,
	authorizer authorizer.Authorizer,
	helmClient *helm.APIClient,
) *VersionREST {
	rest := &VersionREST{
		store:          store,
		platformClient: platformClient,
		registryClient: registryClient,
		registryConfig: registryConfig,
		externalScheme: externalScheme,
		externalHost:   externalHost,
		externalPort:   externalPort,
		externalCAFile: externalCAFile,
		authorizer:     authorizer,
		helmClient:     helmClient,
	}
	return rest
}

// New creates a new chart proxy options object
func (r *VersionREST) New() runtime.Object {
	return &registry.ChartProxyOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *VersionREST) ConnectMethods() []string {
	return []string{"DELETE", "GET"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *VersionREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &registry.ChartProxyOptions{}, true, "version"
}

// Connect returns a handler for the chart proxy
func (r *VersionREST) Connect(ctx context.Context, chartName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := r.store.Get(ctx, chartName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	chart := obj.(*registry.Chart)

	cg, err := r.registryClient.ChartGroups().Get(ctx, chart.Namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	proxyOpts := opts.(*registry.ChartProxyOptions)
	proxyOpts.Version = strings.Trim(proxyOpts.Version, "/") // should do this

	latestChartVersion := proxyOpts.Version
	if proxyOpts.Version == "" {
		if len(chart.Status.Versions) > 0 {
			log.Debug("version is empty, will use latest version")
			// return nil, errors.NewBadRequest("version is required")
			var v1chart = &registryv1.Chart{}
			err := registryv1.Convert_registry_Chart_To_v1_Chart(chart, v1chart, nil)
			if err != nil {
				return nil, errors.NewInternalError(err)
			}
			sorted := sort.ByChartVersion(v1chart.Status.Versions)
			latestChartVersion = sorted[0].Version
		}
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

	return &versionProxyHandler{
		chart:              chart,
		chartGroup:         cg,
		chartVersion:       proxyOpts.Version,
		latestChartVersion: latestChartVersion,

		externalScheme: r.externalScheme,
		externalHost:   r.externalHost,
		externalPort:   r.externalPort,
		externalCAFile: r.externalCAFile,

		registryConfig: r.registryConfig,
		authorizer:     r.authorizer,
		helmOption: helmOption{
			cluster:        proxyOpts.Cluster,
			namespace:      proxyOpts.Namespace,
			platformClient: r.platformClient,
		},
		helmClient:     r.helmClient,
		registryClient: r.registryClient,
	}, nil
}

type versionProxyHandler struct {
	chart              *registry.Chart
	chartGroup         *registry.ChartGroup
	chartVersion       string
	latestChartVersion string

	externalScheme string
	externalHost   string
	externalPort   int
	externalCAFile string

	registryConfig *registryconfig.RegistryConfiguration
	authorizer     authorizer.Authorizer

	helmOption     helmOption
	helmClient     *helm.APIClient
	registryClient *registryinternalclient.RegistryClient
}

type helmOption struct {
	cluster        string
	namespace      string
	platformClient platformversionedclient.PlatformV1Interface
}

func (h *versionProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		{
			h.ServeGetVersion(w, req)
			return
		}
	case "DELETE":
		{
			h.ServeDeleteVersion(w, req)
			return
		}
	default:
		{
			responsewriters.WriteRawJSON(http.StatusForbidden, "Method not allowed", w)
		}
	}
}

// Get chart version info
func (h *versionProxyHandler) ServeGetVersion(w http.ResponseWriter, req *http.Request) {
	client, err := applicationutil.NewHelmClient(req.Context(), h.helmOption.platformClient, h.helmOption.cluster, h.helmOption.namespace)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	host := h.externalHost
	if h.externalPort > 0 {
		host = host + ":" + strconv.Itoa(h.externalPort)
	}

	chartVersion := h.chartVersion
	if chartVersion == "" {
		chartVersion = h.latestChartVersion
	}
	if chartVersion == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, "version is required", w)
		return
	}

	var repo config.RepoConfiguration = config.RepoConfiguration{
		Scheme:        h.externalScheme,
		DomainSuffix:  host,
		CaFile:        h.externalCAFile,
		Admin:         h.registryConfig.Security.AdminUsername,
		AdminPassword: h.registryConfig.Security.AdminPassword,
	}
	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(repo, *h.chartGroup)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	chartPathBasicOptions.Chart = h.chart.Spec.Name
	chartPathBasicOptions.Version = chartVersion

	cpopt := chartPathBasicOptions
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

	var v1ChartSpec = &registryv1.ChartSpec{}
	err = registryv1.Convert_registry_ChartSpec_To_v1_ChartSpec(&h.chart.Spec, v1ChartSpec, nil)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	var v1ChartVersion = &registryv1.ChartVersion{}
	version := getTargetVersion(h.chart, chartVersion)
	err = registryv1.Convert_registry_ChartVersion_To_v1_ChartVersion(&version, v1ChartVersion, nil)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	chartInfo := &registryv1.ChartInfo{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: h.chart.Namespace,
			Name:      h.chart.Name,
		},
		Spec: registryv1.ChartInfoSpec{
			Values:       show.Values,
			Readme:       show.Readme,
			RawFiles:     files,
			ChartSpec:    *v1ChartSpec,
			ChartVersion: *v1ChartVersion,
		},
	}
	responsewriters.WriteRawJSON(http.StatusOK, chartInfo, w)
}

// Delete chart version
func (h *versionProxyHandler) ServeDeleteVersion(w http.ResponseWriter, req *http.Request) {
	if h.chartVersion == "" {
		responsewriters.WriteRawJSON(http.StatusBadRequest, "version is required", w)
		return
	}
	err := h.check(w, req)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusUnauthorized, err.Error(), w)
	}
	if h.helmClient == nil {
		host := h.externalHost
		if h.externalPort > 0 {
			host = host + ":" + strconv.Itoa(h.externalPort)
		}
		loc := &url.URL{
			Scheme: h.externalScheme,
			Host:   registryutil.BuildTenantRegistryDomain(host, h.chart.Spec.TenantID),
			Path:   fmt.Sprintf("/chart/api/%s/charts/%s/%s", h.chart.Spec.ChartGroupName, h.chart.Spec.Name, h.chartVersion),
		}

		// WithContext creates a shallow clone of the request with the new context.
		newReq := req.WithContext(context.Background())
		newReq.Header = netutil.CloneHeader(req.Header)
		newReq.URL = loc
		newReq.SetBasicAuth(h.registryConfig.Security.AdminUsername, h.registryConfig.Security.AdminPassword)
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: loc.Scheme, Host: loc.Host})
		reverseProxy.Transport = transport
		reverseProxy.FlushInterval = 100 * time.Millisecond
		reverseProxy.ErrorLog = log.StdErrLogger()
		reverseProxy.ServeHTTP(w, newReq)
	} else {
		err := harborHandler.DeleteChart(req.Context(), h.helmClient, fmt.Sprintf("%s-chart-%s", h.chart.Spec.TenantID, h.chart.Spec.ChartGroupName), h.chart.Spec.Name)
		if err != nil {
			return
		}
		i := -1
		if len(h.chart.Status.Versions) > 0 {
			for k, v := range h.chart.Status.Versions {
				if v.Version == h.chartVersion {
					i = k
				}
			}
		}
		if i == -1 {
			return
		}
		h.chart.Status.Versions = append(h.chart.Status.Versions[:i], h.chart.Status.Versions[i+1:]...)
		if _, err := h.registryClient.Charts(h.chart.ObjectMeta.Namespace).UpdateStatus(req.Context(), h.chart, metav1.UpdateOptions{}); err != nil {
			log.Error("Failed to update repository versions while deleted",
				log.String("tenantID", h.chart.Spec.TenantID),
				log.String("chartGroupName", h.chart.Spec.ChartGroupName),
				log.String("chartName", h.chart.Spec.Name),
				log.Err(err))
			return
		}
		return
	}
}

func (h *versionProxyHandler) check(w http.ResponseWriter, req *http.Request) error {
	var cg = &registryv1.ChartGroup{}
	err := registryv1.Convert_registry_ChartGroup_To_v1_ChartGroup(h.chartGroup, cg, nil)
	if err != nil {
		return err
	}
	u, exist := genericapirequest.UserFrom(req.Context())
	if !exist || u == nil {
		return fmt.Errorf("empty user info, not authenticated")
	}

	authorized, err := authorizationutil.AuthorizeForChart(req.Context(), u, h.authorizer, "delete", *cg, h.chart.Name)
	if err != nil {
		return err
	}
	if !authorized {
		return fmt.Errorf("not authenticated")
	}
	return nil
}

func getTargetVersion(chart *registry.Chart, version string) registry.ChartVersion {
	for _, v := range chart.Status.Versions {
		if v.Version == version {
			return v
		}
	}
	return registry.ChartVersion{}
}
