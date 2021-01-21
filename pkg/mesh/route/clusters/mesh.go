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
	"net/http"

	"github.com/emicklei/go-restful"
	istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/labels"
	"tkestack.io/tke/pkg/mesh/models"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/errors"
	"tkestack.io/tke/pkg/mesh/util/web"
)

// ListAll list all istio resources
func (h clusterHandler) ListAll(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	selector := req.PathParameter(web.LabelSelectorQuery)

	sel, err := labels.Parse(selector)
	if err != nil {
		result.Err = err.Error()
		return
	}

	ctx := req.Request.Context()
	ret, errs := h.istioService.ListAllResources(ctx, cluster, namespace, "", sel)
	if errs != nil && !errs.Empty() {
		result.Err = errs.Error()
	}
	status = http.StatusOK
	result.Result = true
	result.Data = ret
}

// GetNorthTrafficGateway Get istio resource
func (h clusterHandler) GetNorthTrafficGateway(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	gateway := req.PathParameter(web.GatewayKey)

	entity := &istionetworking.Gateway{}
	entity.SetNamespace(namespace)
	entity.SetName(gateway)

	ctx := req.Request.Context()
	ret, err := h.istioService.GetNorthTrafficGateway(ctx, cluster, entity)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = ret
}

// CreateNorthTrafficGateway create istio resource
func (h clusterHandler) CreateNorthTrafficGateway(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)

	entity := &models.IstioNetworkingConfig{}
	err := req.ReadEntity(entity)
	if err != nil {
		result.Err = err.Error()
		return
	}
	entity.Namespace = models.Namespace{Name: namespace}

	ctx := req.Request.Context()
	_, err = h.istioService.CreateNorthTrafficGateway(ctx, cluster, entity)

	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
}

// UpdateNorthTrafficGateway update istio resource
func (h clusterHandler) UpdateNorthTrafficGateway(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	gateway := req.PathParameter(web.GatewayKey)

	entity := &models.IstioNetworkingConfig{}
	err := req.ReadEntity(entity)
	if err != nil {
		result.Err = err.Error()
		return
	}
	entity.Namespace = models.Namespace{Name: namespace}
	entity.Gateway.SetName(gateway)

	ctx := req.Request.Context()
	_, err = h.istioService.UpdateNorthTrafficGateway(ctx, cluster, entity)

	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
}

// DeleteNorthTrafficGateway Delete istio resource
func (h clusterHandler) DeleteNorthTrafficGateway(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)
	gateway := req.PathParameter(web.GatewayKey)

	entity := &istionetworking.Gateway{}
	entity.SetName(gateway)
	entity.SetNamespace(namespace)

	ctx := req.Request.Context()
	_, err := h.istioService.DeleteNorthTrafficGateway(ctx, cluster, entity)

	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
}
