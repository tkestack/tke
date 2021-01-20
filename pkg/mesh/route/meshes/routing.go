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

package meshes

import (
	"bytes"
	"compress/gzip"
	gojson "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/json"
	"tkestack.io/tke/pkg/mesh/util/proxy"
	"tkestack.io/tke/pkg/mesh/util/web"
	"tkestack.io/tke/pkg/util/log"
)

type meshClusterHandler struct {
	meshService    services.MeshClusterService
	clusterService services.ClusterService
	istioService   services.IstioService

	config meshconfig.MeshConfiguration
}

func New(config meshconfig.MeshConfiguration, clusterService services.ClusterService,
	istioService services.IstioService, meshService services.MeshClusterService) meshClusterHandler {
	return meshClusterHandler{
		meshService:    meshService,
		clusterService: clusterService,
		istioService:   istioService,
		config:         config,
	}
}

func (h meshClusterHandler) AddToWebService(ws *restful.WebService) {
	ws.Route(
		ws.GET(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/%s",
				web.MeshPrefix, web.MeshKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.MicroServicePrefix,
			)).
			To(h.ListMeshServices).
			Param(ws.PathParameter(web.MeshKey, "mesh name").
				DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").
				DataType("string").Required(true)).
			Operation("listMeshServices").
			Doc("List istio micro services").
			Returns(http.StatusOK, "List", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)
	ws.Route(
		ws.GET(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/%s/{%s}",
				web.MeshPrefix, web.MeshKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.MicroServicePrefix, web.MicroServiceKey,
			)).
			To(h.GetMeshService).
			Param(ws.PathParameter(web.MeshKey, "mesh name").
				DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").
				DataType("string").Required(true)).
			Param(ws.PathParameter(web.MicroServiceKey, "micro service name").
				DataType("string").Required(true)).
			Operation("getMeshService").
			Doc("get istio micro services").
			Returns(http.StatusOK, "Get", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	// create Istio resource
	ws.Route(
		ws.POST(
			fmt.Sprintf("/%s/{%s}/%s/{%s}/{%s}",
				web.MeshPrefix, web.MeshKey,
				web.NamespacePrefix, web.NamespaceKey,
				web.IstioResourceKind,
			)).
			To(h.CreateMeshResource).
			Param(ws.PathParameter(web.MeshKey, "mesh name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.NamespaceKey, "namespace name").DataType("string").Required(true)).
			Param(ws.PathParameter(web.IstioResourceKind, "istio resource kind").DataType("string").Required(true).
				AllowableValues(map[string]string{
					web.GatewayPrefix:         web.GatewayPrefix,
					web.VirtualservicePrefix:  web.VirtualservicePrefix,
					web.DestinationrulePrefix: web.DestinationrulePrefix,
					web.ServiceentriePrefix:   web.ServiceentriePrefix,
					web.SidecarPrefix:         web.SidecarPrefix,
					web.EnvoyfilterPrefix:     web.EnvoyfilterPrefix,
					web.WorkloadentryPrefix:   web.WorkloadentryPrefix,
				})).
			Operation("createMeshResource").
			Doc("create mesh resource").
			Returns(http.StatusOK, "Created", rest.Response{}).
			Returns(http.StatusBadRequest, "Error", rest.Response{}).
			Returns(http.StatusNotFound, "Not Found", rest.Response{}).
			Produces(restful.MIME_JSON),
	)

	// create mesh
	ws.Route(ws.POST("/meshes").
		To(h.reverseProxy("/meshes")).
		Operation("createMesh").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON),
	)

	// list mesh
	ws.Route(ws.GET("/meshes").
		To(h.reverseProxy("/meshes")).
		Operation("listMesh").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.QueryParameter("type", "type").DataType("string").
			AllowableValues(map[string]string{"simple": ""})),
	)

	// list mesh from cache
	ws.Route(ws.GET("/meshes/simple").
		To(h.reverseProxy("/meshes/simple")).
		Operation("getMeshSimple").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON),
	)

	// get mesh
	ws.Route(ws.GET("/meshes/{mesh}").
		To(h.reverseProxy("/meshes/{mesh}")).
		Operation("getMesh").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)

	// delete mesh
	ws.Route(ws.DELETE("/meshes/{mesh}").
		To(h.reverseProxy("/meshes/{mesh}")).
		Operation("deleteMesh").
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)

	// upgrade mesh
	ws.Route(ws.POST("/meshes/{mesh}").
		To(h.reverseProxy("/meshes/{mesh}")).
		Operation("patchMeshUpgrade").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)

	// auto injection namespaces
	ws.Route(ws.POST("/meshes/{mesh}/autoinjectionnamespaces").
		To(h.reverseProxy("/meshes/{mesh}/autoinjectionnamespaces")).
		Operation("replaceAutoInjectionNamespaces").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)
	// list mesh clusters namespaces
	ws.Route(ws.GET("/meshes/{mesh}/autoinjectionnamespaces").
		To(h.reverseProxy("/meshes/{mesh}/autoinjectionnamespaces")).
		Operation("getMeshAutoInjectionNamespaces").
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)

	// outbound traffic
	ws.Route(ws.POST("/meshes/{mesh}/outbound-traffic").
		To(h.reverseProxy("/meshes/{mesh}/outbound-traffic")).
		Operation("replaceMeshOutboundTraffic").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)
	// tracing sampling
	ws.Route(ws.POST("/meshes/{mesh}/sampling").
		To(h.reverseProxy("/meshes/{mesh}/sampling")).
		Operation("replaceMeshSampling").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")),
	)

	// canary upgrade confirm or rollback
	ws.Route(ws.POST("/meshes/{mesh}/canary").
		To(h.reverseProxy("/meshes/{mesh}/canary")).
		Operation("patchMeshCanaryUpgrade").
		Returns(http.StatusOK, "Created", rest.Response{}).
		Returns(http.StatusBadRequest, "Error", rest.Response{}).
		Returns(http.StatusNotFound, "Not Found", rest.Response{}).
		Produces(restful.MIME_JSON).
		Param(ws.PathParameter("mesh", "mesh name").DataType("string")).
		Param(ws.PathParameter("operate", "mesh canary operates").DataType("string").
			DefaultValue("confirm").AllowableValues(map[string]string{"confirm": "confirm", "rollback": "rollback"})),
	)
}

// ListMeshServices list mesh services, which contains the Istio "app" label
func (h meshClusterHandler) ListMeshServices(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		if result.Err != "" {
			log.Debugf("%v", result.Err)
		}
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	mesh := req.PathParameter(web.MeshKey)
	namespace := req.PathParameter(web.NamespaceKey)
	ls := req.QueryParameter(web.LabelSelectorQuery)

	if namespace == "*" {
		namespace = ""
	}

	var selector labels.Selector
	if ls != "" {
		var err error
		selector, err = labels.Parse(ls)
		if err != nil {
			result.Err = err.Error()
			return
		}
	}

	ctx := req.Request.Context()

	// 2020-11-04 list mesh services from mesh main cluster
	ret, errs := h.meshService.ListMicroServices(ctx, mesh, namespace, "", selector)
	if errs != nil && !errs.Empty() {
		result.Err = errs.Error()
	}

	status = http.StatusOK
	result.Result = true
	result.Data = ret
	return
}

// GetMeshServices get mesh service, which contains the Istio "app" label
func (h meshClusterHandler) GetMeshService(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		if result.Err != "" {
			log.Debugf("%v", result.Err)
		}
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	mesh := req.PathParameter(web.MeshKey)
	namespace := req.PathParameter(web.NamespaceKey)
	microService := req.PathParameter(web.MicroServiceKey)
	ls := req.QueryParameter(web.LabelSelectorQuery)

	if namespace == "*" {
		namespace = ""
	}

	var selector labels.Selector
	if ls != "" {
		var err error
		selector, err = labels.Parse(ls)
		if err != nil {
			result.Err = err.Error()
			return
		}
	}

	ctx := req.Request.Context()

	// 2020-11-04 get mesh microService from mesh main cluster
	ret, errs := h.meshService.ListMicroServices(ctx, mesh, namespace, microService, selector)
	if errs != nil && !errs.Empty() {
		result.Err = errs.Error()
	}

	status = http.StatusOK
	result.Result = true
	result.Data = ret
	return
}

func (h meshClusterHandler) CreateMeshResource(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		if result.Err != "" {
			log.Debugf("%v", result.Err)
		}
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	ctx := req.Request.Context()
	mesh := req.PathParameter(web.MeshKey)
	namespace := req.PathParameter(web.NamespaceKey)
	kind := req.PathParameter(web.IstioResourceKind)

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

	err = h.meshService.CreateMeshResource(ctx, mesh, entity)
	if err != nil {
		status = http.StatusInternalServerError
		result.Err = err.Error()
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = entity
	return
}

// reverse proxy to mesh
func (h meshClusterHandler) reverseProxy(reverseRestPath string, o ...proxy.Opt) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		result := rest.NewResult(false, "")
		status := http.StatusBadRequest
		defer func() {
			if result.Err != "" {
				log.Debugf("%v", result.Err)
			}
			_ = w.WriteHeaderAndEntity(status, result)
		}()

		rewritePathOpt, err := h.rewriteUrl(r, reverseRestPath)
		if err != nil {
			status = http.StatusInternalServerError
			result.Err = err.Error()
			return
		}

		// proxy to remote mesh-manger api
		opts := []proxy.Opt{
			rewritePathOpt,
			RewriteHeaderOnTCM(),
			RewriteResponseOnTCM(),
		}
		if len(o) > 0 {
			opts = append(opts, o...)
		}
		p := proxy.New(opts...)
		p.Proxy(w, r.Request)
	}
}

const (
	meshManagerContextPath = "/api"
)

func (h meshClusterHandler) rewriteUrl(req *restful.Request, subPath string) (proxy.Opt, error) {
	params := req.PathParameters()
	replaces := make([]string, 0)
	for key, value := range params {
		p := fmt.Sprintf("{%s}", key)
		replaces = append(replaces, p, value)
	}
	if len(replaces) != 0 {
		r := strings.NewReplacer(replaces...)
		subPath = r.Replace(subPath)
	}

	tcmConfig := h.config.Components.MeshManager.Address
	rawurl := tcmConfig + meshManagerContextPath + subPath
	parse, err := url.Parse(rawurl)
	if err != nil {
		log.Errorf("TCM API base url invalid: %s", rawurl)
		return nil, err
	}

	return func(p *proxy.Proxy) {
		target := &url.URL{
			Scheme:   parse.Scheme,
			Host:     parse.Host,
			Path:     parse.Path,
			RawQuery: req.Request.URL.RawQuery,
		}
		p.Request.URL = target
	}, nil
}

func RewriteResponseOnTCM() proxy.Opt {
	return func(p *proxy.Proxy) {
		p.ModifyResponse = handleTCMResponse
	}
}

// handleTCMResponse handle TCM API response data structure
// 0. success direct data
// [{
//   "name": "mesh-adj129ojm"
// }]
// 1. success without data
// {
// 	"code": 200,
// 	"message": "OK"
// }
// 2. success with data
// {
// 	"code": 200,
// 	"data": [{"test": "name"}]
// }
// 3. failure plain text
// HTTP Status 404
// 404 Not Found
func handleTCMResponse(res *http.Response) error {
	// finally return response data struct
	rsp := proxy.Response{
		ResultCode: res.StatusCode,
	}

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		reader.Close()
	default:
		reader = res.Body
	}

	bs, _ := ioutil.ReadAll(reader)
	success := res.StatusCode == http.StatusOK

	msg := json.Get(bs, "message")
	data := json.Get(bs, "data")

	hasMsg := msg.LastError() == nil
	hasData := data.LastError() == nil

	if hasMsg {
		rsp.Msg = msg.ToString()
	} else {
		log.Warnf("mesh-manager API response: %s", msg.LastError().Error())
	}

	if hasData {
		rsp.Data = gojson.RawMessage([]byte(data.ToString()))
	} else {
		log.Warnf("mesh-manager API response: %s", data.LastError().Error())
	}

	if success && !hasData && !hasMsg {
		rsp.Data = gojson.RawMessage(bs)
	}

	if !success && !hasMsg {
		rsp.Msg = string(bs)
	}

	returnResponse(res, rsp)
	return nil
}

func returnResponse(res *http.Response, body interface{}) {
	rb, _ := json.Marshal(body)
	length := len(rb)
	reader := bytes.NewReader(rb)
	res.Body = ioutil.NopCloser(reader)
	res.Header.Del("Content-Encoding")
	res.Header.Set("Content-Length", strconv.Itoa(length))
}

func RewriteHeaderOnTCM() proxy.Opt {
	return func(p *proxy.Proxy) {
		p.Request.Header = map[string][]string{
			"X-Remote-TenantID": []string{"default"},
			"Content-Type":      []string{"application/json"},
			"Accept":            []string{"*/*"},
		}
	}
}
