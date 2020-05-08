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
	"net/http"
	"strings"

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
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientrest "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util"
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

	dynamicClient, err := util.DynamicClientByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}

	clusterOpts := opts.(*platform.ClusterApplyOptions)

	return &handler{
		clusterName:       clusterName,
		client:            client,
		dynamicClient:     dynamicClient,
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
	dynamicClient     dynamic.Interface
	exts              []*runtime.RawExtension
	groupVersionKinds []*schema.GroupVersionKind
	metaAccessor      meta.MetadataAccessor
	notUpdate         bool
}

type applyObject struct {
	obj             runtime.Object
	isCreateRequest bool
	restClient      clientrest.Interface
	dynamicClient   dynamic.ResourceInterface
	namespace       string
	kind            string
	name            string
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
	applyObjects := make([]*applyObject, len(h.exts))

	for idx, ext := range h.exts {
		obj := ext.Object
		gvk := h.groupVersionKinds[idx]

		name, namespace, objectStatus := objectName(h.metaAccessor, obj)
		if objectStatus != nil {
			return objectStatus
		}

		var ao *applyObject
		var status *metav1.Status
		if dynamicResource := matchCRD(h.dynamicClient, gvk.Group, gvk.Version, gvk.Kind); dynamicResource != nil {
			ao, status = applyObjectFromDynamicClient(dynamicResource, gvk, name, namespace, h.notUpdate, h.metaAccessor, obj)
		} else {
			ao, status = applyObjectFromClientSet(ctx, h.client, gvk, name, namespace, h.notUpdate, h.metaAccessor, obj)
		}
		if status != nil {
			return status
		}
		applyObjects[idx] = ao
	}

	var messages []string
	for _, applyObj := range applyObjects {
		var message string
		var status *metav1.Status
		if applyObj.dynamicClient != nil {
			message, status = createOrUpdateFromDynamicClient(applyObj.dynamicClient, applyObj.isCreateRequest, applyObj.name, applyObj.kind, applyObj.obj)
		} else {
			message, status = createOrUpdateFromClientSet(ctx, applyObj.restClient, applyObj.isCreateRequest, applyObj.name, applyObj.namespace, applyObj.kind, applyObj.obj)
		}
		if status != nil {
			return status
		}
		messages = append(messages, message)
	}

	return &metav1.Status{
		Status:  metav1.StatusSuccess,
		Code:    http.StatusOK,
		Message: strings.Join(messages, "\n"),
	}
}

func objectName(metaAccessor meta.MetadataAccessor, obj runtime.Object) (string, string, *metav1.Status) {
	namespace, err := metaAccessor.Namespace(obj)
	if err != nil {
		return "", "", errorInternal
	}
	name, err := metaAccessor.Name(obj)
	if err != nil {
		return "", "", errorInternal
	}
	genName, err := metaAccessor.GenerateName(obj)
	if err != nil {
		return "", "", errorInternal
	}

	if len(name) == 0 && len(genName) == 0 {
		return "", "", errorBadName
	}
	return name, namespace, nil
}

func createOrUpdateFromDynamicClient(dynamicClient dynamic.ResourceInterface, isCreate bool, name, kind string, obj runtime.Object) (string, *metav1.Status) {
	unstructuredObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return "", errorInternal
	}
	if isCreate {
		_, err := dynamicClient.Create(unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			if statusError, ok := err.(*errors.StatusError); ok {
				status := statusError.Status()
				return "", &status
			}
			return "", unknownError(err)
		}
		if len(name) != 0 {
			return fmt.Sprintf("%s %s created", kind, name), nil
		}
		return fmt.Sprintf("%s generated", kind), nil
	}
	_, err := dynamicClient.Update(unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		if statusError, ok := err.(*errors.StatusError); ok {
			status := statusError.Status()
			return "", &status
		}
		return "", unknownError(err)
	}
	return fmt.Sprintf("%s %s configured", kind, name), nil
}

func createOrUpdateFromClientSet(ctx context.Context, client clientrest.Interface, isCreate bool, name, namespace, kind string, obj runtime.Object) (string, *metav1.Status) {
	if isCreate {
		// create
		result := client.Post().
			Context(ctx).
			NamespaceIfScoped(parseNamespaceIfScoped(namespace, kind)).
			Resource(util.ResourceFromKind(kind)).
			Body(obj).
			Do()
		err := result.Error()
		if err != nil {
			if statusError, ok := err.(*errors.StatusError); ok {
				status := statusError.Status()
				return "", &status
			}
			return "", unknownError(err)
		}
		if len(name) != 0 {
			return fmt.Sprintf("%s %s created", kind, name), nil
		}
		return fmt.Sprintf("%s generated", kind), nil
	}
	// update
	result := client.Put().
		Context(ctx).
		NamespaceIfScoped(parseNamespaceIfScoped(namespace, kind)).
		Resource(util.ResourceFromKind(kind)).
		Name(name).
		Body(obj).
		Do()
	err := result.Error()
	if err != nil {
		if statusError, ok := err.(*errors.StatusError); ok {
			status := statusError.Status()
			return "", &status
		}
		return "", unknownError(err)
	}
	return fmt.Sprintf("%s %s configured", kind, name), nil
}

func applyObjectFromDynamicClient(dynamicClient dynamic.NamespaceableResourceInterface, gvk *schema.GroupVersionKind, name, namespace string, notUpdate bool, metaAccessor meta.MetadataAccessor, obj runtime.Object) (*applyObject, *metav1.Status) {
	var resource dynamic.ResourceInterface
	if ns, namespaceScoped := parseNamespaceIfScoped(namespace, gvk.Kind); namespaceScoped {
		resource = dynamicClient.Namespace(ns)
	} else {
		resource = dynamicClient
	}
	if len(name) != 0 {
		result, err := resource.Get(name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			if statusError, ok := err.(*errors.StatusError); ok {
				status := statusError.Status()
				return nil, &status
			}
			return nil, unknownError(err)
		}
		if err == nil {
			if notUpdate {
				return nil, alreadyExistError(gvk, name)
			}
			resourceVersion, err := metaAccessor.ResourceVersion(obj)
			if err != nil {
				return nil, errorInternal
			}
			if resourceVersion != "" {
				return nil, errorHasResourceVersion
			}
			if err := metaAccessor.SetResourceVersion(obj, result.GetResourceVersion()); err != nil {
				return nil, errorInternal
			}
			return &applyObject{
				obj:             obj,
				isCreateRequest: false,
				dynamicClient:   resource,
				namespace:       namespace,
				kind:            gvk.Kind,
				name:            name,
			}, nil
		}
	}
	return &applyObject{
		obj:             obj,
		isCreateRequest: true,
		dynamicClient:   resource,
		namespace:       namespace,
		kind:            gvk.Kind,
		name:            name,
	}, nil
}

func applyObjectFromClientSet(ctx context.Context, client *kubernetes.Clientset, gvk *schema.GroupVersionKind, name, namespace string, notUpdate bool, metaAccessor meta.MetadataAccessor, obj runtime.Object) (*applyObject, *metav1.Status) {
	restClient := util.RESTClientFor(client, gvk.Group, gvk.Version)

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
				return nil, &status
			}
			return nil, unknownError(err)
		}
		if err == nil {
			if notUpdate {
				return nil, alreadyExistError(gvk, name)
			}
			returnedObj, err := result.Get()
			if err != nil {
				return nil, errorInternal
			}
			resourceVersion, err := metaAccessor.ResourceVersion(obj)
			if err != nil {
				return nil, errorInternal
			}
			if resourceVersion != "" {
				return nil, errorHasResourceVersion
			}
			savedResourceVersion, err := metaAccessor.ResourceVersion(returnedObj)
			if err != nil {
				return nil, errorInternal
			}
			if err := metaAccessor.SetResourceVersion(obj, savedResourceVersion); err != nil {
				return nil, errorInternal
			}
			return &applyObject{
				obj:             obj,
				isCreateRequest: false,
				restClient:      restClient,
				namespace:       namespace,
				kind:            gvk.Kind,
				name:            name,
			}, nil
		}
	}
	// create
	return &applyObject{
		obj:             obj,
		isCreateRequest: true,
		restClient:      restClient,
		namespace:       namespace,
		kind:            gvk.Kind,
		name:            name,
	}, nil
}

func unknownError(err error) *metav1.Status {
	return &metav1.Status{
		Code:    http.StatusInternalServerError,
		Status:  metav1.StatusFailure,
		Reason:  metav1.StatusReasonInternalError,
		Message: err.Error(),
	}
}

func alreadyExistError(gvk *schema.GroupVersionKind, name string) *metav1.Status {
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

func parseNamespaceIfScoped(namespace string, kind string) (string, bool) {
	kindLower := strings.ToLower(kind)
	namespaceScoped := true

	if kindLower == "namespace" ||
		kindLower == "node" ||
		kindLower == "componentstatus" ||
		kindLower == "persistentvolume" ||
		kindLower == "storageclass" ||
		kindLower == "volumeattachment" ||
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

func matchCRD(dynamicClient dynamic.Interface, apiGroup, apiVersion, kind string) dynamic.NamespaceableResourceInterface {
	apiGroup = strings.ToLower(apiGroup)
	crd := false
	if strings.HasPrefix(apiGroup, "networking.istio.io") {
		crd = true
	}
	if strings.HasSuffix(apiGroup, "tkestack.io") {
		crd = true
	}
	if strings.HasSuffix(apiGroup, "cloud.tencent.com") {
		crd = true
	}
	if crd {
		return dynamicClient.Resource(schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  apiVersion,
			Resource: util.ResourceFromKind(kind),
		})
	}
	return nil
}
