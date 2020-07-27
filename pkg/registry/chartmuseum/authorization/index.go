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
	username, userTenantID := authentication.UsernameAndTenantID(req.Context())
	if chartGroup.Spec.Visibility == registry.VisibilityPublic {
		// visibility is public, everyone can pull
		if username == "" && userTenantID == "" {
			log.Debug("Anonymous added public chart repo",
				log.String("tenantID", tenantID),
				log.String("repo", chartGroupName))
		} else {
			log.Debug("User added public chart repo",
				log.String("tenantID", tenantID),
				log.String("repo", chartGroupName),
				log.String("userTenantID", userTenantID),
				log.String("username", username))
		}
		a.nextHandler.ServeHTTP(w, req)
		return
	}
	if username != "" && userTenantID != "" && userTenantID == tenantID {
		// authorized
		log.Debug("User added private chart repo",
			log.String("tenantID", tenantID),
			log.String("repo", chartGroupName),
			log.String("userTenantID", userTenantID),
			log.String("username", username))
		a.nextHandler.ServeHTTP(w, req)
		return
	}
	if username == "" && tenantID == "" {
		// anonymous user and chart group is private
		log.Warn("Anonymous user try add a private chart repo",
			log.String("tenantID", tenantID),
			log.String("repo", chartGroupName))
		a.notAuthenticated(w, req)
		return
	}
	log.Warn("Not authorized user try add a private chart repo",
		log.String("tenantID", tenantID),
		log.String("repo", chartGroupName),
		log.String("userTenantID", userTenantID),
		log.String("username", username))
	a.forbidden(w)
}
