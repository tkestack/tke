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
	"bytes"
	"context"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	clientrest "k8s.io/client-go/rest"
	"net/http"
	"strings"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
)

var (
	errorInternal = &metav1.Status{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Status",
			APIVersion: "v1",
		},
		Status:  metav1.StatusFailure,
		Code:    http.StatusInternalServerError,
		Reason:  metav1.StatusReasonInternalError,
		Message: "Internal Server Error",
	}
	errorBadName = &metav1.Status{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Status",
			APIVersion: "v1",
		},
		Status:  metav1.StatusFailure,
		Code:    http.StatusBadRequest,
		Reason:  metav1.StatusReasonInvalid,
		Message: "Name or generateName must be special",
	}
	errorHasResourceVersion = &metav1.Status{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Status",
			APIVersion: "v1",
		},
		Status:  metav1.StatusFailure,
		Code:    http.StatusConflict,
		Reason:  metav1.StatusReasonConflict,
		Message: "Not need specify resource version",
	}
)

// ApplyREST implement bucket call interface for cluster.
type ApplyREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
func (r *ApplyREST) New() runtime.Object {
	return &platform.Cluster{}
}

// Connect returns an http.Handler that will handle the request/response for a
// given API invocation.
func (r *ApplyREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, _ rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	client, err := util.ClientSetByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}

	clusterOpts := opts.(*platform.ClusterApplyOptions)

	return &handler{
		clusterName:       clusterName,
		client:            client,
		exts:              make([]*runtime.RawExtension, 0),
		groupVersionKinds: make([]*schema.GroupVersionKind, 0),
		metaAccessor:      meta.NewAccessor(),
		notUpdate:         clusterOpts.NotUpdate,
	}, nil
}

// NewConnectOptions returns an empty options object that will be used to pass
// options to the Connect method.
func (r *ApplyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.ClusterApplyOptions{}, false, ""
}

// ConnectMethods returns the list of HTTP methods handled by Connect
func (r *ApplyREST) ConnectMethods() []string {
	return []string{
		http.MethodPost,
	}
}

// ProducesMIMETypes returns a list of the MIME types the specified HTTP verb (GET, POST, DELETE,
// PATCH) can respond with.
func (r *ApplyREST) ProducesMIMETypes(_ string) []string {
	return []string{"application/json"}
}

// ProducesObject returns an object the specified HTTP verb respond with. It will overwrite storage object if
// it is not nil. Only the type of the return object matters, the value will be ignored.
func (r *ApplyREST) ProducesObject(_ string) interface{} {
	return metav1.Status{}
}

type handler struct {
	clusterName       string
	client            *kubernetes.Clientset
	exts              []*runtime.RawExtension
	groupVersionKinds []*schema.GroupVersionKind
	metaAccessor      meta.MetadataAccessor
	notUpdate         bool
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()
	if err := h.decode(req.Body); err != nil {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest(err.Error()).Status(), writer)
		return
	}
	if len(h.exts) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("must special at lease one resource").Status(), writer)
		return
	}
	status := h.apply(req.Context())
	responsewriters.WriteRawJSON(int(status.Code), status, writer)
}

func (h *handler) decode(r io.Reader) error {
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
	for {
		ext := runtime.RawExtension{}
		if err := decoder.Decode(&ext); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error parsing: %v", err)
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		obj, gkv, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			return err
		}
		ext.Object = obj
		h.groupVersionKinds = append(h.groupVersionKinds, gkv)
		h.exts = append(h.exts, &ext)
	}
}

func (h *handler) apply(ctx context.Context) *metav1.Status {
	type applyObject struct {
		obj             runtime.Object
		isCreateRequest bool
		restClient      clientrest.Interface
		namespace       string
		kind            string
		name            string
	}

	applyObjects := make([]*applyObject, len(h.exts))

	for idx, ext := range h.exts {
		obj := ext.Object
		gvk := h.groupVersionKinds[idx]
		restClient := util.RESTClientFor(h.client, gvk.Group, gvk.Version)

		namespace, err := h.metaAccessor.Namespace(obj)
		if err != nil {
			return errorInternal
		}

		name, err := h.metaAccessor.Name(obj)
		if err != nil {
			return errorInternal
		}

		genName, err := h.metaAccessor.GenerateName(obj)
		if err != nil {
			return errorInternal
		}

		if len(name) == 0 && len(genName) == 0 {
			return errorBadName
		}

		if len(name) != 0 {
			result := restClient.Get().
				Context(ctx).
				NamespaceIfScoped(parseNamespaceIfScoped(namespace, gvk.Kind)).
				Resource(util.ResourceFromKind(gvk.Kind)).
				Name(name).
				Do()
			err := result.Error()
			if err != nil && !errors.IsNotFound(err) {
				if statusError, ok := err.(*errors.StatusError); ok {
					status := statusError.Status()
					return &status
				}
				return unknownError(err)
			}
			if err == nil {
				if h.notUpdate {
					return &metav1.Status{
						TypeMeta: metav1.TypeMeta{
							Kind:       "Status",
							APIVersion: "v1",
						},
						Status: metav1.StatusFailure,
						Code:   http.StatusConflict,
						Reason: metav1.StatusReasonAlreadyExists,
						Details: &metav1.StatusDetails{
							Name:  name,
							Group: gvk.Group,
							Kind:  gvk.Kind,
						},
						Message: fmt.Sprintf("%s \"%s\" already exists", gvk.Kind, name),
					}
				}
				returnedObj, err := result.Get()
				if err != nil {
					return errorInternal
				}
				resourceVersion, err := h.metaAccessor.ResourceVersion(obj)
				if err != nil {
					return errorInternal
				}
				if resourceVersion != "" {
					return errorHasResourceVersion
				}
				savedResourceVersion, err := h.metaAccessor.ResourceVersion(returnedObj)
				if err != nil {
					return errorInternal
				}
				if err := h.metaAccessor.SetResourceVersion(obj, savedResourceVersion); err != nil {
					return errorInternal
				}
				applyObjects[idx] = &applyObject{
					obj:             obj,
					isCreateRequest: false,
					restClient:      restClient,
					namespace:       namespace,
					kind:            gvk.Kind,
					name:            name,
				}
				continue
			}
		}
		// create
		applyObjects[idx] = &applyObject{
			obj:             obj,
			isCreateRequest: true,
			restClient:      restClient,
			namespace:       namespace,
			kind:            gvk.Kind,
			name:            name,
		}
	}

	var messages []string
	for _, applyObj := range applyObjects {
		if applyObj.isCreateRequest {
			// create
			result := applyObj.restClient.Post().
				Context(ctx).
				NamespaceIfScoped(parseNamespaceIfScoped(applyObj.namespace, applyObj.kind)).
				Resource(util.ResourceFromKind(applyObj.kind)).
				Body(applyObj.obj).
				Do()
			log.Debugf("Apply cluster bucket create call: %v", applyObj)
			err := result.Error()
			if err != nil {
				if statusError, ok := err.(*errors.StatusError); ok {
					status := statusError.Status()
					return &status
				}
				return unknownError(err)
			}
			if len(applyObj.name) != 0 {
				messages = append(messages, fmt.Sprintf("%s %s created", applyObj.kind, applyObj.name))
			} else {
				messages = append(messages, fmt.Sprintf("%s generated", applyObj.kind))
			}
		} else {
			// update
			result := applyObj.restClient.Put().
				Context(ctx).
				NamespaceIfScoped(parseNamespaceIfScoped(applyObj.namespace, applyObj.kind)).
				Resource(util.ResourceFromKind(applyObj.kind)).
				Name(applyObj.name).
				Body(applyObj.obj).
				Do()
			log.Debugf("Apply cluster bucket update call: %v", applyObj)
			err := result.Error()
			if err != nil {
				if statusError, ok := err.(*errors.StatusError); ok {
					status := statusError.Status()
					return &status
				}
				return unknownError(err)
			}
			messages = append(messages, fmt.Sprintf("%s %s configured", applyObj.kind, applyObj.name))
		}
	}

	return &metav1.Status{
		Status:  metav1.StatusSuccess,
		Code:    http.StatusOK,
		Message: strings.Join(messages, "\n"),
	}
}

func unknownError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusInternalServerError,
		Status:  metav1.StatusFailure,
		Reason:  metav1.StatusReasonInternalError,
		Message: err.Error(),
	}
}

func parseNamespaceIfScoped(namespace string, kind string) (string, bool) {
	kindLower := strings.ToLower(kind)
	namespaceScoped := true

	if kindLower == "namespace" ||
		kindLower == "node" ||
		kindLower == "componentstatus" ||
		kindLower == "persistentvolume" ||
		kindLower == "storageclass" ||
		kindLower == "volumeattachment" ||
		kindLower == "serviceaccount" ||
		kindLower == "clusterrole" ||
		kindLower == "clusterrolebinding" {
		namespaceScoped = false
	}

	ns := namespace
	if namespaceScoped && ns == "" {
		ns = "default"
	}
	return ns, namespaceScoped
}
