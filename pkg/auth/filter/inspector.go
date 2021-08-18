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
	"net/http"
	"regexp"
	"time"
	"tkestack.io/tke/pkg/apiserver/util"

	platformv1 "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/log"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	genericfilters "k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	rbaclisters "k8s.io/client-go/listers/rbac/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	serviceAccountRegExp = regexp.MustCompile(`^system:serviceaccount:([^:]+):(.+)$`)
)

type Inspector interface {
	Inspect(handler http.Handler, c *genericapiserver.Config) http.Handler
}

type clusterInspector struct {
	k8sClient          kubernetes.Interface
	crbLister          rbaclisters.ClusterRoleBindingLister
	crLister           rbaclisters.ClusterRoleLister
	platformClient     platformv1.PlatformV1Interface
	privilegedUsername string
}

func NewClusterInspector(platformClient platformv1.PlatformV1Interface, privilegedUsername string) (Inspector, error) {
	k8sClient, err := apiclient.BuildKubeClient()
	if err != nil {
		return nil, err
	}
	informerFactory := informers.NewSharedInformerFactory(k8sClient, time.Minute)
	clusterRoleBindingInformer := informerFactory.Rbac().V1().ClusterRoleBindings()
	clusterRoleBindingLister := clusterRoleBindingInformer.Lister()
	clusterRoleInformer := informerFactory.Rbac().V1().ClusterRoles()
	clusterRoleLister := clusterRoleInformer.Lister()
	stopCh := util.SetupSignalHandler()
	informerFactory.Start(stopCh)
	if ok := cache.WaitForCacheSync(stopCh, clusterRoleBindingInformer.Informer().HasSynced,
		clusterRoleInformer.Informer().HasSynced); !ok {
		return nil, fmt.Errorf("failed to wait for namespaces caches to sync")
	}
	return &clusterInspector{
		k8sClient:          k8sClient,
		crbLister:          clusterRoleBindingLister,
		crLister:           clusterRoleLister,
		platformClient:     platformClient,
		privilegedUsername: privilegedUsername,
	}, nil
}

func isClusterAdmin(rules []rbacv1.PolicyRule) bool {
	if len(rules) != 2 {
		return false
	}
	isAdmin := true
	for _, rul := range rules {
		if len(rul.APIGroups) == 1 && rul.APIGroups[0] == "*" &&
			len(rul.Resources) == 1 && rul.Resources[0] == "*" &&
			len(rul.Verbs) == 1 && rul.Verbs[0] == "*" {
			continue
		}
		if len(rul.NonResourceURLs) == 1 && rul.NonResourceURLs[0] == "*" &&
			len(rul.Verbs) == 1 && rul.Verbs[0] == "*" {
			continue
		}
		isAdmin = false
		break
	}
	return isAdmin
}

func (i *clusterInspector) needInspect(ctx context.Context, privilegedUsername string) bool {
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if (username == privilegedUsername || username == "system:apiserver") && tenantID == "" {
		return false
	}

	clusterRoleBindings, err := i.crbLister.List(labels.Everything())
	if err != nil {
		log.Errorf("query clusterRoleBindings failed: %+v", err)
		return true
	}
	matches := serviceAccountRegExp.FindStringSubmatch(username)
	if len(matches) != 3 {
		return true
	}
	namespace := matches[1]
	username = matches[2]
	for _, crb := range clusterRoleBindings {
		for _, sub := range crb.Subjects {
			if sub.Name == username && sub.Namespace == namespace {
				cr, err := i.crLister.Get(crb.RoleRef.Name)
				if err != nil {
					log.Errorf("query clusterRole: %+v failed: %+v", crb.RoleRef.Name, err)
					continue
				}
				if len(cr.Rules) != 2 {
					continue
				}
				log.Debugf("needInspect: username: %+v, namespace: %+v, clusterRole: %+v->%v",
					username, namespace, cr.Name, cr.Rules)
				if isClusterAdmin(cr.Rules) {
					return false
				}
			}
		}
	}
	return true
}

func (i *clusterInspector) Inspect(handler http.Handler, c *genericapiserver.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if !i.needInspect(ctx, i.privilegedUsername) {
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
		username, tenantID := authentication.UsernameAndTenantID(ctx)
		log.Infof(" clusterNames: %+v, username: %+v, tenant: %+v, "+
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
			if AllowClusterNotFoundActions.Has(verb) && k8serrors.IsNotFound(err) {
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
