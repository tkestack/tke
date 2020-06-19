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

package authorization

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/mux"
	"helm.sh/chartmuseum/pkg/repo"

	// "helm.sh/chartmuseum/pkg/repo"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
)

// getChart serve http get request on /chart/{tenantID}/{chartGroup}/charts/{file}
func (a *authorization) getChart(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tenantID, ok := vars["tenantID"]
	if !ok || tenantID == "" {
		a.notFound(w)
		return
	}
	chartGroupName, ok := vars["chartGroup"]
	if !ok || chartGroupName == "" {
		a.notFound(w)
		return
	}
	file, ok := vars["file"]
	if !ok || file == "" {
		a.notFound(w)
		return
	}
	chartName, _, ok := ChartNameVersionFromFile(file)
	if !ok {
		a.notFound(w)
		return
	}
	chartObject, err := a.validateGetChart(w, req, tenantID, chartGroupName, chartName)
	if err != nil {
		return
	}
	sw := &statusWrite{ResponseWriter: w}
	a.nextHandler.ServeHTTP(sw, req)
	if sw.status != http.StatusOK {
		return
	}
	if err := a.afterGetChart(req.Context(), chartObject); err != nil {
		log.Error("Failed to update registry chart resource", log.Err(err))
	}
}

// apiGetChart serve http get request on /chart/api/{tenantID}/{chartGroup}/charts/{name}
func (a *authorization) apiGetChart(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tenantID, ok := vars["tenantID"]
	if !ok || tenantID == "" {
		a.notFound(w)
		return
	}
	chartGroupName, ok := vars["chartGroup"]
	if !ok || chartGroupName == "" {
		a.notFound(w)
		return
	}
	chartName, ok := vars["name"]
	if !ok || chartName == "" {
		a.notFound(w)
		return
	}
	_, err := a.validateGetChart(w, req, tenantID, chartGroupName, chartName)
	if err != nil {
		return
	}
	a.nextHandler.ServeHTTP(w, req)
}

// apiGetChartVersion serve http get request on /chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}
func (a *authorization) apiGetChartVersion(w http.ResponseWriter, req *http.Request) {
	a.apiGetChart(w, req)
}

func (a *authorization) validateGetChart(w http.ResponseWriter, req *http.Request, tenantID, chartGroupName, chartName string) (*registry.Chart, error) {
	chartGroupList, err := a.registryClient.ChartGroups().List(req.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, chartGroupName),
	})
	if err != nil {
		log.Error("Failed to list chart group by tenantID and name",
			log.String("tenantID", tenantID),
			log.String("name", chartGroupName),
			log.Err(err))
		a.internalError(w)
		return nil, err
	}
	if len(chartGroupList.Items) == 0 {
		a.notFound(w)
		return nil, fmt.Errorf("not found")
	}
	chartGroup := chartGroupList.Items[0]
	if chartGroup.Status.Locked != nil && *chartGroup.Status.Locked {
		// Chart group is locked
		a.locked(w)
		return nil, fmt.Errorf("locked")
	}
	chartList, err := a.registryClient.Charts(chartGroup.ObjectMeta.Name).List(req.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, chartName),
	})
	if err != nil {
		log.Error("Failed to list chart by tenantID and name of group",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.String("name", chartName),
			log.Err(err))
		a.internalError(w)
		return nil, err
	}
	if len(chartList.Items) == 0 {
		a.notFound(w)
		return nil, fmt.Errorf("not found")
	}
	chartObject := chartList.Items[0]
	if chartObject.Spec.Visibility == registry.VisibilityPrivate {
		username, userTenantID := authentication.UsernameAndTenantID(req.Context())
		if username == "" && userTenantID == "" {
			a.notAuthenticated(w, req)
			return nil, fmt.Errorf("not authenticated")
		}
		if userTenantID != tenantID {
			a.forbidden(w)
			return nil, fmt.Errorf("forbidden")
		}
	}
	return &chartObject, nil
}

func (a *authorization) afterGetChart(ctx context.Context, chartObject *registry.Chart) error {
	chartObject.Status.PullCount = chartObject.Status.PullCount + 1
	if _, err := a.registryClient.Charts(chartObject.ObjectMeta.Namespace).UpdateStatus(ctx, chartObject, metav1.UpdateOptions{}); err != nil {
		log.Error("Failed to update repository pull count while pulled",
			log.String("tenantID", chartObject.Spec.TenantID),
			log.String("chartGroupName", chartObject.Spec.ChartGroupName),
			log.String("chartName", chartObject.Spec.Name),
			log.Err(err))
		return err
	}
	return nil
}

// ChartNameVersionFromFile returns chart name and version from chart filename.
func ChartNameVersionFromFile(file string) (name, version string, ok bool) {
	var filename string
	if strings.HasSuffix(file, fmt.Sprintf(".%s", repo.ChartPackageFileExtension)) {
		filename = strings.TrimSuffix(file, fmt.Sprintf(".%s", repo.ChartPackageFileExtension))
	} else if strings.HasSuffix(file, fmt.Sprintf(".%s", repo.ProvenanceFileExtension)) {
		filename = strings.TrimSuffix(file, fmt.Sprintf(".%s", repo.ProvenanceFileExtension))
	}
	if filename == "" {
		return
	}
	i := strings.LastIndex(filename, "-")
	if i == -1 {
		return
	}
	name = filename[:i]
	version = filename[i+1:]
	ok = true
	return
}
