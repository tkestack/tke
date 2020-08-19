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
	"net/http"

	"github.com/gorilla/mux"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/registry"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/registry/chartmuseum/model"
	"tkestack.io/tke/pkg/util/log"
)

// apiCreateChart serve http post request on /chart/api/{tenantID}/{chartGroup}/charts
func (a *authorization) apiCreateChart(w http.ResponseWriter, req *http.Request) {
	a.doAPICreateProvenance(w, req, "chart")
}

// apiCreateProvenance serve http post request on /chart/api/{tenantID}/{chartGroup}/prov
func (a *authorization) apiCreateProvenance(w http.ResponseWriter, req *http.Request) {
	a.doAPICreateProvenance(w, req, "prov")
}

func (a *authorization) doAPICreateProvenance(w http.ResponseWriter, req *http.Request, fieldName string) {
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

	file, header, err := req.FormFile(fieldName)
	if err != nil {
		log.Error("Failed to retrieve chart file from request", log.Err(err))
		return
	}
	ct, err := loader.LoadArchive(file)
	if err != nil {
		log.Error("Failed to load chart from request body", log.Err(err))
		return
	}
	if ct.Metadata == nil {
		log.Error("Chart metadata is nil after parsed", log.Any("chart", ct))
		return
	}

	chartGroup, err := a.validateAPICreateChart(w, req, tenantID, chartGroupName, ct.Metadata.Name)
	if err != nil {
		return
	}
	sw := &statusBodyWrite{ResponseWriter: w}
	a.nextHandler.ServeHTTP(sw, req)
	if sw.status != http.StatusCreated {
		return
	}
	var savedResponse model.SavedResponse
	if err := json.Unmarshal(sw.body, &savedResponse); err != nil {
		log.Error("Failed to unmarshal response of chartmuseum", log.ByteString("body", sw.body), log.Err(err))
		return
	}
	if !savedResponse.Saved {
		log.Error("Chartmuseum server that does not meet expectations", log.ByteString("body", sw.body), log.Int("status", sw.status))
		return
	}
	// helm push: error converting YAML to JSON: yaml: control characters are not allowed
	//
	// bs, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Error("Failed to read all content from chart file", log.Err(err))
	// 	return
	// }
	// ct := new(chart.Metadata)
	// if err := yaml.Unmarshal(bs, ct); err != nil {
	// 	log.Error("Failed to unmarshal chart file", log.Err(err))
	// 	return
	// }
	// if err := a.afterAPICreateChart(req.Context(), chartGroup, ct, header.Size); err != nil {
	// 	log.Error("Failed to update registry chart resource", log.Err(err))
	// }
	if err := a.afterAPICreateChart(req.Context(), chartGroup, ct.Metadata, header.Size); err != nil {
		log.Error("Failed to update registry chart resource", log.Err(err))
	}
}

func (a *authorization) validateAPICreateChart(w http.ResponseWriter, req *http.Request, tenantID, chartGroupName, chartName string) (*registry.ChartGroup, error) {
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
		// Chart group must first be created via console
		a.notFound(w)
		return nil, fmt.Errorf("not found")
	}
	chartGroup := chartGroupList.Items[0]
	if chartGroup.Status.Locked != nil && *chartGroup.Status.Locked {
		// Chart group is locked
		a.locked(w)
		return nil, fmt.Errorf("locked")
	}

	if a.isAdmin(w, req) {
		return &chartGroup, nil
	}

	var cg = &registryv1.ChartGroup{}
	err = registryv1.Convert_registry_ChartGroup_To_v1_ChartGroup(&chartGroup, cg, nil)
	if err != nil {
		log.Error("Failed to convert ChartGroup",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.String("chartName", chartName),
			log.Err(err))
		a.internalError(w)
		return nil, err
	}

	authorized, err := AuthorizeForChart(w, req, a.authorizer, "create", *cg, "")
	if err != nil {
		log.Error("Failed to authorize for chart",
			log.String("tenantID", tenantID),
			log.String("chartGroupName", chartGroupName),
			log.String("chartName", chartName),
			log.Err(err))
		a.internalError(w)
		return nil, err
	}
	if !authorized {
		a.notAuthenticated(w, req)
		return nil, fmt.Errorf("not authenticated")
	}

	return &chartGroup, nil
}

func (a *authorization) afterAPICreateChart(ctx context.Context, chartGroup *registry.ChartGroup, chartMeta *chart.Metadata, ctSize int64) error {
	chartList, err := a.registryClient.Charts(chartGroup.ObjectMeta.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.chartGroupName=%s", chartGroup.Spec.TenantID, chartMeta.Name, chartGroup.Spec.Name),
	})
	if err != nil {
		return err
	}

	newVersion := registry.ChartVersion{
		Version:     chartMeta.Version,
		ChartSize:   ctSize,
		TimeCreated: metav1.Now(),
		Description: chartMeta.Description,
		AppVersion:  chartMeta.AppVersion,
		Icon:        chartMeta.Icon,
	}
	if len(chartList.Items) == 0 {
		if _, err := a.registryClient.Charts(chartGroup.ObjectMeta.Name).Create(ctx, &registry.Chart{
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
					if _, err := a.registryClient.Charts(chartGroup.ObjectMeta.Name).UpdateStatus(ctx, &chartObject, metav1.UpdateOptions{}); err != nil {
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
			if _, err := a.registryClient.Charts(chartGroup.ObjectMeta.Name).UpdateStatus(ctx, &chartObject, metav1.UpdateOptions{}); err != nil {
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

	chartList, err = a.registryClient.Charts(chartGroup.ObjectMeta.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.chartGroupName=%s", chartGroup.Spec.TenantID, chartGroup.Spec.Name),
	})
	if err != nil {
		return err
	}
	// update chart group's chart count
	chartGroup.Status.ChartCount = int32(len(chartList.Items))
	if _, err := a.registryClient.ChartGroups().UpdateStatus(ctx, chartGroup, metav1.UpdateOptions{}); err != nil {
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
