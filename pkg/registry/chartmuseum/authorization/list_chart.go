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
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"
)

// apiListChart serve http get request on /chart/api/{tenantID}/{chartGroup}/charts
func (a *authorization) apiListChart(w http.ResponseWriter, req *http.Request) {
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
	_, err := a.validateListChart(w, req, tenantID, chartGroupName)
	if err != nil {
		return
	}
	a.nextHandler.ServeHTTP(w, req)
}

func (a *authorization) validateListChart(w http.ResponseWriter, req *http.Request, tenantID, chartGroupName string) (*registry.ChartGroup, error) {
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
	if chartGroup.Spec.Visibility == registry.VisibilityPrivate {
		username, userTenantID := authentication.GetUsernameAndTenantID(req.Context())
		if username == "" && userTenantID == "" {
			a.notAuthenticated(w, req)
			return nil, fmt.Errorf("not authenticated")
		}
		if userTenantID != tenantID {
			a.forbidden(w)
			return nil, fmt.Errorf("forbidden")
		}
	}
	return &chartGroup, nil
}
