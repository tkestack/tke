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
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/registry/chartmuseum/model"
	"tkestack.io/tke/pkg/util/log"
)

// apiDeleteChartVersion serve http delete request on /chart/api/{tenantID}/{chartGroup}/charts/{name}/{version}
func (a *authorization) apiDeleteChartVersion(w http.ResponseWriter, req *http.Request) {
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
	chartVersion, ok := vars["version"]
	if !ok || chartVersion == "" {
		a.notFound(w)
		return
	}
	chartObject, err := a.validateGetChart(w, req, tenantID, chartGroupName, chartName)
	if err != nil {
		return
	}
	sw := &statusBodyWrite{ResponseWriter: w}
	a.nextHandler.ServeHTTP(sw, req)
	if sw.status != http.StatusOK {
		return
	}
	var deletedResponse model.DeletedResponse
	if err := json.Unmarshal(sw.body, &deletedResponse); err != nil {
		log.Error("Failed to unmarshal response of chartmuseum", log.ByteString("body", sw.body), log.Err(err))
		return
	}
	if !deletedResponse.Deleted {
		log.Error("Chartmuseum server that does not meet expectations", log.ByteString("body", sw.body), log.Int("status", sw.status))
		return
	}
	if err := a.afterAPIDeleteChartVersion(req.Context(), chartObject, chartVersion); err != nil {
		log.Error("Failed to delete chart version from resource", log.Err(err))
	}
}

func (a *authorization) afterAPIDeleteChartVersion(ctx context.Context, chartObject *registry.Chart, version string) error {
	i := -1
	if len(chartObject.Status.Versions) > 0 {
		for k, v := range chartObject.Status.Versions {
			if v.Version == version {
				i = k
			}
		}
	}
	if i == -1 {
		return nil
	}
	chartObject.Status.Versions = append(chartObject.Status.Versions[:i], chartObject.Status.Versions[i+1:]...)
	if _, err := a.registryClient.Charts(chartObject.ObjectMeta.Namespace).UpdateStatus(ctx, chartObject, metav1.UpdateOptions{}); err != nil {
		log.Error("Failed to update repository versions while deleted",
			log.String("tenantID", chartObject.Spec.TenantID),
			log.String("chartGroupName", chartObject.Spec.ChartGroupName),
			log.String("chartName", chartObject.Spec.Name),
			log.Err(err))
		return err
	}
	return nil
}
