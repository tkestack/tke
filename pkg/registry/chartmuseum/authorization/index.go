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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/util/log"
)

// index serve http get request on /chart/{tenantID}/{chartGroup}/index.yaml
func (a *authorization) index(w http.ResponseWriter, req *http.Request) {
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
	chartGroupList, err := a.registryClient.ChartGroups().List(req.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", tenantID, chartGroupName),
	})
	if err != nil {
		log.Error("Failed to list chart group by tenantID and name",
			log.String("tenantID", tenantID),
			log.String("name", chartGroupName),
			log.Err(err))
		a.internalError(w)
		return
	}
	if len(chartGroupList.Items) == 0 {
		// Chart group must first be created via console
		a.notFound(w)
		return
	}
	chartGroup := chartGroupList.Items[0]
	if chartGroup.Status.Locked != nil && *chartGroup.Status.Locked {
		// Chart group is locked
		a.locked(w)
		return
	}

	if a.isAdmin(w, req) {
		a.nextHandler.ServeHTTP(w, req)
		return
	}

	var cg = &registryv1.ChartGroup{}
	err = registryv1.Convert_registry_ChartGroup_To_v1_ChartGroup(&chartGroup, cg, nil)
	if err != nil {
		log.Error("Failed to convert ChartGroup",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.Err(err))
		a.internalError(w)
		return
	}

	authorized, err := AuthorizeForChartGroup(w, req, a.authorizer, "get", *cg)
	if err != nil {
		log.Error("Failed to authorize for chartGroup",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.Err(err))
		a.internalError(w)
		return
	}
	if !authorized {
		a.notAuthenticated(w, req)
		return
	}
	a.nextHandler.ServeHTTP(w, req)
}
