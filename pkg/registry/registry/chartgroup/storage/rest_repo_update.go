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
	"fmt"
	"net/http"

	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	applicationutil "tkestack.io/tke/pkg/application/util"
	"tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/log"
)

// RepoUpdateREST adapts a service registry into apiserver's RESTStorage model.
type RepoUpdateREST struct {
	store          ChartGroupStorage
	registryClient *registryinternalclient.RegistryClient
	authorizer     authorizer.Authorizer
}

// NewRepoUpdateREST returns a wrapper around the underlying generic storage and performs
// allocations and deallocations of various chart.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewRepoUpdateREST(
	store ChartGroupStorage,
	registryClient *registryinternalclient.RegistryClient,
	authorizer authorizer.Authorizer,
) *RepoUpdateREST {
	rest := &RepoUpdateREST{
		store:          store,
		registryClient: registryClient,
		authorizer:     authorizer,
	}
	return rest
}

// New creates a new chart proxy options object
func (r *RepoUpdateREST) New() runtime.Object {
	return new(registry.ChartGroup)
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *RepoUpdateREST) ConnectMethods() []string {
	return []string{"POST"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *RepoUpdateREST) NewConnectOptions() (runtime.Object, bool, string) {
	return nil, false, ""
}

// Connect returns a handler for the chartgroup proxy
func (r *RepoUpdateREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := r.store.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cg := obj.(*registry.ChartGroup)

	if cg.Spec.Type != registry.RepoTypeImported {
		return nil, errors.NewInternalError(fmt.Errorf("chartGroup %s type is not %s", cg.Spec.Name, registry.RepoTypeImported))
	}
	return &repoUpdateProxyHandler{
		registryClient: r.registryClient,
		chartGroup:     cg,
		authorizer:     r.authorizer,
	}, nil
}

type repoUpdateProxyHandler struct {
	registryClient *registryinternalclient.RegistryClient
	chartGroup     *registry.ChartGroup
	authorizer     authorizer.Authorizer
}

func (h *repoUpdateProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	client := applicationutil.NewHelmClientWithoutRESTClient()

	password, err := util.VerifyDecodedPassword(h.chartGroup.Spec.ImportedInfo.Password)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	entry, err := client.RepoUpdate(&helmaction.RepoUpdateOptions{
		ChartPathOptions: helmaction.ChartPathOptions{
			ChartRepo: h.chartGroup.Spec.TenantID + "/" + h.chartGroup.Spec.Name,
			RepoURL:   h.chartGroup.Spec.ImportedInfo.Addr,
			Username:  h.chartGroup.Spec.ImportedInfo.Username,
			Password:  password,
		},
	})
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	err = h.syncChart(req.Context(), entry)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	err = h.syncChartGroup(req.Context(), entry)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, "Success", w)
}

func (h *repoUpdateProxyHandler) syncChartGroup(ctx context.Context, entries map[string]repo.ChartVersions) error {
	cg := h.chartGroup.DeepCopy()
	cg.Status.ChartCount = int32(len(entries))
	_, err := h.registryClient.ChartGroups().UpdateStatus(ctx, cg, metav1.UpdateOptions{})
	return err
}

func (h *repoUpdateProxyHandler) syncChart(ctx context.Context, entries map[string]repo.ChartVersions) error {
	for name, versions := range entries {
		newVersions := make([]registry.ChartVersion, len(versions))
		for k, v := range versions {
			newVersions[k] = registry.ChartVersion{
				Version:     v.Version,
				TimeCreated: metav1.Time{Time: v.Created},
				Description: v.Description,
				AppVersion:  v.AppVersion,
				Icon:        v.Icon,
			}
		}

		chart, found, err := h.findChart(ctx, h.chartGroup, name)
		if err != nil {
			return err
		}
		if found {
			chart.Status.Versions = newVersions
			_, err = h.registryClient.Charts(chart.ObjectMeta.Namespace).UpdateStatus(ctx, chart, metav1.UpdateOptions{})
		} else {
			_, err = h.registryClient.Charts(h.chartGroup.ObjectMeta.Name).Create(ctx, &registry.Chart{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: h.chartGroup.ObjectMeta.Name,
				},
				Spec: registry.ChartSpec{
					Name:           name,
					TenantID:       h.chartGroup.Spec.TenantID,
					ChartGroupName: h.chartGroup.Spec.Name,
					Visibility:     h.chartGroup.Spec.Visibility,
				},
				Status: registry.ChartStatus{
					PullCount: 0,
					Versions:  newVersions,
				},
			}, metav1.CreateOptions{})
		}
		if err != nil {
			log.Error("Failed to create/update chart by tenantID and name",
				log.String("tenantID", h.chartGroup.Spec.TenantID),
				log.String("name", name),
				log.Err(err))
			return err
		}
	}
	return nil
}

func (h *repoUpdateProxyHandler) findChart(ctx context.Context, cg *registry.ChartGroup, name string) (chart *registry.Chart, found bool, err error) {
	list, err := h.registryClient.Charts(cg.Name).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", cg.Spec.TenantID, name),
	})
	if err != nil {
		log.Error("Failed to list chart by tenantID and name",
			log.String("tenantID", cg.Spec.TenantID),
			log.String("name", name),
			log.Err(err))
		return nil, false, err
	}
	if len(list.Items) == 0 {
		// Chart group must first be created via console
		return nil, false, nil
	}

	return list.Items[0].DeepCopy(), true, nil
}
