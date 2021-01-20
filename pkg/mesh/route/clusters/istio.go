/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package clusters

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/constants"
	"tkestack.io/tke/pkg/mesh/util/errors"
	"tkestack.io/tke/pkg/mesh/util/web"
	"tkestack.io/tke/pkg/util/log"
)

// ListIstioResources list istio resources
func (h clusterHandler) ListIstioResources(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	kind := req.PathParameter(web.IstioResourceKind)
	labelSelector := req.QueryParameter(web.LabelSelectorQuery)

	ret := &unstructured.UnstructuredList{}
	switch kind {
	case web.GatewayPrefix:
		ret.SetGroupVersionKind(constants.GatewayList)
	case web.VirtualservicePrefix:
		ret.SetGroupVersionKind(constants.VirtualServiceList)
	case web.DestinationrulePrefix:
		ret.SetGroupVersionKind(constants.DestinationRuleList)
	case web.ServiceentriePrefix:
		ret.SetGroupVersionKind(constants.ServiceEntryList)
	case web.SidecarPrefix:
		ret.SetGroupVersionKind(constants.SidecarList)
	case web.WorkloadentryPrefix:
		ret.SetGroupVersionKind(constants.WorkloadEntryList)
	case web.EnvoyfilterPrefix:
		ret.SetGroupVersionKind(constants.EnvoyFilterList)
	default:
		result.Err = fmt.Sprintf("unsupported istio resource[%s]", kind)
		return
	}

	var selector labels.Selector
	var err error
	if labelSelector != "" {
		selector, err = labels.Parse(labelSelector)
		if err != nil {
			log.Errorf("labelSelector invalid: [%s], error: %v", labelSelector, err)
			result.Err = fmt.Sprintf("labelSelector invalid: [%s]", labelSelector)
			return
		}
	}

	ctx := req.Request.Context()
	err = h.istioService.ListResources(ctx, cluster, ret,
		&ctrlclient.ListOptions{Namespace: namespace, LabelSelector: selector})
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleApiError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = ret.Items
	return
}

// GetIstioResource get istio resource
func (h clusterHandler) GetIstioResource(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	kind := req.PathParameter(web.IstioResourceKind)
	name := req.PathParameter(web.IstioResourceKey)

	entity := &unstructured.Unstructured{}
	switch kind {
	case web.GatewayPrefix:
		entity.SetGroupVersionKind(constants.Gateway)
	case web.VirtualservicePrefix:
		entity.SetGroupVersionKind(constants.VirtualService)
	case web.DestinationrulePrefix:
		entity.SetGroupVersionKind(constants.DestinationRule)
	case web.ServiceentriePrefix:
		entity.SetGroupVersionKind(constants.ServiceEntry)
	case web.SidecarPrefix:
		entity.SetGroupVersionKind(constants.Sidecar)
	case web.WorkloadentryPrefix:
		entity.SetGroupVersionKind(constants.WorkloadEntry)
	case web.EnvoyfilterPrefix:
		entity.SetGroupVersionKind(constants.EnvoyFilter)
	default:
		result.Err = fmt.Sprintf("unsupported istio resource[%s]", kind)
		return
	}
	entity.SetNamespace(namespace)
	entity.SetName(name)

	ctx := req.Request.Context()
	err := h.istioService.GetResource(ctx, cluster, entity)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleApiError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
	return
}

// CreateIstioResource create istio resource
func (h clusterHandler) CreateIstioResource(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	kind := req.PathParameter(web.IstioResourceKind)

	switch kind {
	case web.GatewayPrefix:
	case web.VirtualservicePrefix:
	case web.DestinationrulePrefix:
	default:
		result.Err = fmt.Sprintf("unsupported istio resource[%s]", kind)
		return
	}

	entity := &unstructured.Unstructured{}
	err := req.ReadEntity(entity)
	if err != nil {
		result.Err = err.Error()
		return
	}
	entity.SetNamespace(namespace)

	ctx := req.Request.Context()
	_, err = h.istioService.CreateResource(ctx, cluster, entity)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleApiError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
	return
}

// DeleteIstioResource delete istio resource
func (h clusterHandler) DeleteIstioResource(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	kind := req.PathParameter(web.IstioResourceKind)
	name := req.PathParameter(web.IstioResourceKey)

	switch kind {
	case web.GatewayPrefix:
	case web.VirtualservicePrefix:
	case web.DestinationrulePrefix:
	case web.ServiceentriePrefix:
	case web.SidecarPrefix:
	case web.EnvoyfilterPrefix:
	case web.WorkloadentryPrefix:
	default:
		result.Err = fmt.Sprintf("unsupported istio resource[%s]", kind)
		return
	}

	entity := &unstructured.Unstructured{}
	err := req.ReadEntity(entity)
	if err != nil {
		result.Err = err.Error()
		return
	}
	entity.SetNamespace(namespace)
	entity.SetName(name)

	ctx := req.Request.Context()
	_, err = h.istioService.DeleteResource(ctx, cluster, entity)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleApiError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
	return
}
