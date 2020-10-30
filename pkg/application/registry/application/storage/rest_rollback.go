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
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/rest"
	application "tkestack.io/tke/api/application"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
)

// RollbackREST adapts a service registry into apiserver's RESTStorage model.
type RollbackREST struct {
	store             ApplicationStorage
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV1Interface
}

// NewRollbackREST returns a wrapper around the underlying generic storage and performs
// rollback of helm releases.
// TODO: all transactional behavior should be supported from within generic storage
//   or the strategy.
func NewRollbackREST(
	store ApplicationStorage,
	applicationClient *applicationinternalclient.ApplicationClient,
	platformClient platformversionedclient.PlatformV1Interface,
) *RollbackREST {
	rest := &RollbackREST{
		store:             store,
		applicationClient: applicationClient,
		platformClient:    platformClient,
	}
	return rest
}

// New creates a new chart proxy options object
func (r *RollbackREST) New() runtime.Object {
	return &applicationv1.RollbackProxyOptions{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *RollbackREST) ConnectMethods() []string {
	return []string{"GET", "POST"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *RollbackREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &applicationv1.RollbackProxyOptions{}, false, ""
}

// Connect returns a handler for the chart proxy
func (r *RollbackREST) Connect(ctx context.Context, appName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	obj, err := r.store.Get(ctx, appName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	app := obj.(*application.App)
	proxyOpts := opts.(*applicationv1.RollbackProxyOptions)

	if proxyOpts.Revision == 0 {
		return nil, errors.NewBadRequest("revision is required")
	}
	if proxyOpts.Cluster == "" {
		return nil, errors.NewBadRequest("cluster is required")
	}

	return &proxyHandler{
		app:               app,
		revision:          proxyOpts.Revision,
		cluster:           proxyOpts.Cluster,
		applicationClient: r.applicationClient,
		platformClient:    r.platformClient,
	}, nil
}

type proxyHandler struct {
	app               *application.App
	revision          int64
	cluster           string
	applicationClient *applicationinternalclient.ApplicationClient
	platformClient    platformversionedclient.PlatformV1Interface
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newStatus := h.app.Status.DeepCopy()
	newStatus.Phase = application.AppPhaseRollingBack
	newStatus.RollbackRevision = h.revision
	newStatus.LastTransitionTime = metav1.Now()
	newObj := h.app.DeepCopy()
	newObj.Status = *newStatus

	updated, err := h.applicationClient.Apps(h.app.Namespace).UpdateStatus(req.Context(), newObj, metav1.UpdateOptions{})
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	responsewriters.WriteRawJSON(http.StatusOK, updated, w)
}
