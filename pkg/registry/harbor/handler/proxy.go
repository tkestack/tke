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

package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	// "mime/multipart"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"

	"tkestack.io/tke/api/registry"

	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	harbor "tkestack.io/tke/pkg/registry/harbor/client"
	"tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/transport"

	jsoniter "github.com/json-iterator/go"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/registry/chartmuseum/model"
)

type handler struct {
	reverseProxy      *httputil.ReverseProxy
	host              string
	manifestRegexp    *regexp.Regexp
	createChartRegexp *regexp.Regexp
}

type HarborContextKey string

const manifestPattern = "/v2/.*/.*/manifests/.*"
const createChartPattern = "/api/chartrepo/.*/charts"

var registryClient *registryinternalclient.RegistryClient
var harborClient *harbor.APIClient
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NewHandler to create a reverse proxy handler and returns it.
func NewHandler(address string, cafile string, loopbackConfig *restclient.Config, registryConfig *registryconfig.RegistryConfiguration) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	tr, err := transport.NewOneWayTLSTransport(cafile, true)
	if err != nil {
		log.Error("Failed to create one-way HTTPS transport", log.String("caFile", cafile), log.Err(err))
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: u.Scheme, Host: u.Host})
	reverseProxy.Transport = tr
	reverseProxy.ModifyResponse = rewriteBody
	reverseProxy.ErrorLog = log.StdErrLogger()
	re, err := regexp.Compile(manifestPattern)
	if err != nil {
		log.Error("Failed to init harbor manifest pattern")
		return nil, err
	}
	chartRe, err := regexp.Compile(createChartPattern)
	if err != nil {
		log.Error("Failed to init harbor manifest pattern")
		return nil, err
	}
	registryLocalClient, err := registryinternalclient.NewForConfig(loopbackConfig)
	if err != nil {
		return nil, err
	}
	registryClient = registryLocalClient
	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(
		registryConfig.Security.AdminUsername+":"+registryConfig.Security.AdminPassword),
	)
	cfg := &harbor.Configuration{
		BasePath:      fmt.Sprintf("https://%s/api/v2.0", registryConfig.DomainSuffix),
		DefaultHeader: headers,
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		HTTPClient: &http.Client{
			Transport: tr,
		},
	}
	harborClient = harbor.NewAPIClient(cfg)
	return &handler{reverseProxy, u.Host, re, chartRe}, nil
}

func rewriteBody(resp *http.Response) (err error) {

	ctx := resp.Request.Context()
	host := ctx.Value(HarborContextKey("host"))
	externalHost := ctx.Value(HarborContextKey("exHost"))
	authHeader := resp.Header.Get("www-authenticate")
	if authHeader != "" {
		header := fmt.Sprintf("Bearer realm=\"https://%s/service/token\",service=\"harbor-registry\"", externalHost)
		log.Debug("Modify backend harbor header www-authenticate", log.String("header", header))
		resp.Header.Set("www-authenticate", header)
	}

	locationHeader := resp.Header.Get("location")
	if locationHeader != "" {
		log.Debug("Replace harbor location header", log.String("original host", host.(string)), log.String("tke host", externalHost.(string)))
		resp.Header.Set("location", strings.ReplaceAll(locationHeader, host.(string), externalHost.(string)))
	}

	manifestPattern := ctx.Value(HarborContextKey("manifestPattern"))
	createChartPattern := ctx.Value(HarborContextKey("createChartPattern"))
	if manifestPattern.(string) == "true" && resp.StatusCode < 300 && resp.StatusCode >= 200 {
		pathSpiltted := strings.Split(resp.Request.URL.Path, "/")
		harborProject := pathSpiltted[2]
		repoName := pathSpiltted[3]
		tagName := pathSpiltted[5]
		tenantID := ctx.Value(HarborContextKey("tenantID"))
		namespaceName := strings.ReplaceAll(harborProject, tenantID.(string)+"-image-", "")
		namespaceList, err := registryClient.Namespaces().List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, namespaceName),
		})
		if err != nil {
			return err
		}
		if len(namespaceList.Items) == 0 {
			return fmt.Errorf("namespace %s in tenant %s not exist", namespaceName, tenantID)
		}
		namespaceObject := namespaceList.Items[0]
		repoList, err := registryClient.Repositories(namespaceObject.ObjectMeta.Name).List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.namespaceName=%s", tenantID, repoName, namespaceName),
		})
		if err != nil {
			return err
		}
		var repoObject *registry.Repository
		if len(repoList.Items) > 0 {

			repoObject = &repoList.Items[0]

		}
		if resp.Request.Method == "PUT" {
			artifact, _, err := harborClient.ArtifactApi.GetArtifact(ctx, harborProject, repoName, tagName, nil)
			if err != nil {
				return fmt.Errorf("harbor artifact /%s/%s/%s not exist", harborProject, repoName, tagName)
			}
			util.PushRepository(ctx, registryClient, &namespaceObject, repoObject, repoName, tagName, artifact.Digest)

		} else if resp.Request.Method == "GET" {
			util.PullRepository(ctx, registryClient, &namespaceObject, repoObject, repoName, tagName)
		}

	} else if createChartPattern.(string) == "true" && resp.StatusCode < 300 && resp.StatusCode >= 200 {
		tenantID := ctx.Value(HarborContextKey("tenantID"))
		chartInfo := ctx.Value(HarborContextKey("chartInfo"))
		chartSize := ctx.Value(HarborContextKey("chartSize"))
		pathSpiltted := strings.Split(resp.Request.URL.Path, "/")
		ct := chartInfo.(*chart.Chart)
		err = afterCreateChart(ctx, resp, ct, chartSize.(int64), tenantID.(string), strings.Replace(pathSpiltted[3], tenantID.(string)+"-chart-", "", 1), "chart")
		if err != nil {
			return
		}
	}

	return nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	log.Debug("Reverse proxy to backend harbor", log.String("url", req.URL.Path))

	originalTkeHost := req.Host
	req.Host = h.host

	ctx := context.WithValue(req.Context(), HarborContextKey("host"), h.host)
	ctx = context.WithValue(ctx, HarborContextKey("exHost"), originalTkeHost)
	ctx = context.WithValue(ctx, HarborContextKey("manifestPattern"), strconv.FormatBool(h.manifestRegexp.MatchString(req.URL.Path)))
	isCreateChart := h.createChartRegexp.MatchString(req.URL.Path)
	ctx = context.WithValue(ctx, HarborContextKey("createChartPattern"), strconv.FormatBool(isCreateChart))

	if isCreateChart {

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		file, header, err := req.FormFile("chart")
		if err != nil {
			log.Error("Failed to retrieve chart file from request", log.Err(err))
			return
		}
		ct, err := loader.LoadArchive(file)
		if err != nil {
			log.Error("Failed to load chart from request body", log.Err(err))
			return
		}
		ctx = context.WithValue(ctx, HarborContextKey("chartInfo"), ct)
		ctx = context.WithValue(ctx, HarborContextKey("chartSize"), header.Size)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	req = req.WithContext(ctx)
	h.reverseProxy.ServeHTTP(w, req)

}

func getChartGroup(tenantID, chartGroupName string) (*registry.ChartGroup, error) {
	chartGroupList, err := registryClient.ChartGroups().List(context.Background(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, chartGroupName),
	})
	if err != nil {
		log.Error("Failed to list chart group by tenantID and name",
			log.String("tenantID", tenantID),
			log.String("name", chartGroupName),
			log.Err(err))
		return nil, err
	}
	if len(chartGroupList.Items) == 0 {
		// Chart group must first be created via console
		return nil, fmt.Errorf("not found")
	}
	chartGroup := chartGroupList.Items[0]
	if chartGroup.Status.Locked != nil && *chartGroup.Status.Locked {
		// Chart group is locked
		return nil, fmt.Errorf("locked")
	}

	var cg = &registryv1.ChartGroup{}
	err = registryv1.Convert_registry_ChartGroup_To_v1_ChartGroup(&chartGroup, cg, nil)
	if err != nil {
		log.Error("Failed to convert ChartGroup",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.Err(err))
		return nil, err
	}

	return &chartGroup, nil
}

func afterCreateChart(ctx context.Context, resp *http.Response, chart *chart.Chart, chartSize int64, tenantID, chartGroupName, fieldName string) (err error) {

	chartMeta := chart.Metadata

	var savedResponse model.SavedResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err = json.Unmarshal(body, &savedResponse); err != nil {
		log.Error("Failed to unmarshal response of chartmuseum", log.ByteString("body", body), log.Err(err))
		return err
	}
	if !savedResponse.Saved {
		log.Error("Harbor server that does not meet expectations", log.ByteString("body", body), log.Int("status", resp.StatusCode))
		return fmt.Errorf("not saved")
	}

	chartGroup, err := getChartGroup(tenantID, chartGroupName)
	if err != nil {
		return err
	}

	chartList, err := registryClient.Charts(chartGroup.ObjectMeta.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.chartGroupName=%s", chartGroup.Spec.TenantID, chartMeta.Name, chartGroup.Spec.Name),
	})
	if err != nil {
		return err
	}

	newVersion := registry.ChartVersion{
		Version:     chartMeta.Version,
		ChartSize:   chartSize,
		TimeCreated: metav1.Now(),
		Description: chartMeta.Description,
		AppVersion:  chartMeta.AppVersion,
		Icon:        chartMeta.Icon,
	}
	if len(chartList.Items) == 0 {
		if _, err := registryClient.Charts(chartGroup.ObjectMeta.Name).Create(ctx, &registry.Chart{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: chartGroup.ObjectMeta.Name,
			},
			Spec: registry.ChartSpec{
				Name:           chartMeta.Name,
				TenantID:       chartGroup.Spec.TenantID,
				ChartGroupName: chartGroup.Spec.Name,
				Visibility:     chartGroup.Spec.Visibility,
			},
			Status: registry.ChartStatus{
				PullCount: 0,
				Versions:  []registry.ChartVersion{newVersion},
			},
		}, metav1.CreateOptions{}); err != nil {
			log.Error("Failed to create chart while pushed chart",
				log.String("tenantID", chartGroup.Spec.TenantID),
				log.String("chartGroupName", chartGroup.Spec.Name),
				log.String("chartName", chartMeta.Name),
				log.String("version", chartMeta.Version),
				log.Err(err))
			return err
		}
	} else {
		chartObject := chartList.Items[0]
		existVersion := false
		if len(chartObject.Status.Versions) > 0 {
			for k, v := range chartObject.Status.Versions {
				if v.Version == chartMeta.Version {
					existVersion = true
					chartObject.Status.Versions[k] = newVersion
					if _, err := registryClient.Charts(chartGroup.ObjectMeta.Name).UpdateStatus(ctx, &chartObject, metav1.UpdateOptions{}); err != nil {
						log.Error("Failed to update chart version while chart pushed",
							log.String("tenantID", chartGroup.Spec.TenantID),
							log.String("chartGroupName", chartGroup.Spec.Name),
							log.String("chartName", chartMeta.Name),
							log.String("version", chartMeta.Version),
							log.Err(err))
						return err
					}
					break
				}
			}
		}

		if !existVersion {
			chartObject.Status.Versions = append(chartObject.Status.Versions, newVersion)
			if _, err := registryClient.Charts(chartGroup.ObjectMeta.Name).UpdateStatus(ctx, &chartObject, metav1.UpdateOptions{}); err != nil {
				log.Error("Failed to create repository tag while received notification",
					log.String("tenantID", chartGroup.Spec.TenantID),
					log.String("chartGroupName", chartGroup.Spec.Name),
					log.String("chartName", chartMeta.Name),
					log.String("version", chartMeta.Version),
					log.Err(err))
				return err
			}
		}
	}

	chartList, err = registryClient.Charts(chartGroup.ObjectMeta.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.chartGroupName=%s", chartGroup.Spec.TenantID, chartGroup.Spec.Name),
	})
	if err != nil {
		return err
	}
	// update chart group's chart count
	chartGroup.Status.ChartCount = int32(len(chartList.Items))
	if _, err := registryClient.ChartGroups().UpdateStatus(ctx, chartGroup, metav1.UpdateOptions{}); err != nil {
		log.Error("Failed to update chart group's chart count while pushed",
			log.String("tenantID", chartGroup.Spec.TenantID),
			log.String("chartGroupName", chartGroup.Spec.Name),
			log.String("chartName", chartMeta.Name),
			log.String("version", chartMeta.Version),
			log.Err(err))
		return err
	}
	return nil

}
