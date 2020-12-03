/*
 * Tencent is pleased to support the open source community by making TKEStack available.
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
 */

package filter

import (
	"context"
	"fmt"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"net/http"

	platformv1 "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericfilters "k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type Inspector interface {
	Inspect(handler http.Handler, c *genericapiserver.Config) http.Handler
}

type clusterInspector struct {
	platformClient     platformv1.PlatformV1Interface
	privilegedUsername string
}

func NewClusterInspector(platformClient platformv1.PlatformV1Interface, privilegedUsername string) Inspector {
	return &clusterInspector{
		platformClient:     platformClient,
		privilegedUsername: privilegedUsername,
	}
}

func (i *clusterInspector) Inspect(handler http.Handler, c *genericapiserver.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		username, tenantID := authentication.UsernameAndTenantID(ctx)
		if (username == i.privilegedUsername || username == "system:apiserver") && tenantID == "" {
			handler.ServeHTTP(w, req)
			return
		}
		ae := request.AuditEventFrom(ctx)
		attributes, err := genericfilters.GetAuthorizerAttributes(ctx)
		if err != nil {
			responsewriters.InternalError(w, req, err)
			return
		}
		tkeAttributes := ConvertTKEAttributes(ctx, attributes)
		verb := tkeAttributes.GetVerb()
		clusterNames := ExtractClusterNames(ctx, req, tkeAttributes.GetResource())
		if len(clusterNames) == 0 {
			handler.ServeHTTP(w, req)
			return
		}
		if len(clusterNames) > maxCheckClusterNameCount &&
			verb != createProjectAction && verb != updateProjectAction {
			ForbiddenResponse(ctx, tkeAttributes, w, req, ae, c.Serializer,
				"invalid request: too many clusterName in request")
			return
		}
		log.Infof("WithTKEAuthorization clusterNames: %+v, username: %+v, tenant: %+v, "+
			"action: %+v, resource: %+v, name: %+v",
			clusterNames, username, tenantID, tkeAttributes.GetVerb(),
			tkeAttributes.GetResource(), tkeAttributes.GetName())
		reason, valid := CheckClustersTenant(ctx, tenantID, clusterNames, i.platformClient, verb)
		if !valid {
			ForbiddenResponse(ctx, tkeAttributes, w, req, ae, c.Serializer, reason)
			return
		}
		handler.ServeHTTP(w, req)
	})
}

func CheckClustersTenant(ctx context.Context, tenantID string, clusterNames []string,
	platformClient platformv1.PlatformV1Interface, verb string) (string, bool) {
	for _, clusterName := range clusterNames {
		cluster, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
		if err != nil {
			if verb == updateProjectAction && k8serrors.IsNotFound(err) {
				continue
			}
			return fmt.Sprintf("check cluster: %+v err: %+v", clusterName, err), false
		}
		if tenantID != cluster.Spec.TenantID {
			return fmt.Sprintf("cluster: %+v has invalid tenantID", clusterName), false
		}
	}
	return "", true
}
