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

package auth

import (
	"fmt"
	"github.com/docker/distribution/registry/auth/token"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/util/log"
)

type userRequest struct {
	username      string
	userTenantID  string
	authenticated bool
}

// An accessFilter will filter access based on userinfo
type accessFilter interface {
	filter(a *token.ResourceActions, u *userRequest) error
}

type registryFilter struct {
}

func (reg *registryFilter) filter(a *token.ResourceActions, _ *userRequest) error {
	// Do not filter if the request is to access registry catalog
	if a.Name != "catalog" {
		return fmt.Errorf("unable to handle, type: %s, name: %s", a.Type, a.Name)
	}
	a.Actions = []string{}
	return nil
}

// repositoryFilter filters the access based on Harbor's permission model
type repositoryFilter struct {
	parser         imageParser
	registryClient *registryinternalclient.RegistryClient
	adminUsername  string
}

func (r *repositoryFilter) filter(a *token.ResourceActions, u *userRequest) error {
	// clear action list to assign to new access element after perm check.
	img, err := r.parser.parse(a.Name)
	if err != nil {
		return err
	}
	permission := ""

	if img.namespace == "" || img.repo == "" {
		log.Debugf("Namespace `%s` OR repo `%s` does not exist, set empty permission", img.namespace, img.repo)
		a.Actions = []string{}
		return nil
	}

	if img.tenantID == "" {
		if u.authenticated && u.username == r.adminUsername && u.userTenantID == "" {
			permission = "RWM"
		} else {
			permission = "R"
		}
		a.Actions = permToActions(permission)
		return nil
	}

	namespaceList, err := r.registryClient.Namespaces().List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", img.tenantID, img.namespace),
	})
	if err != nil {
		return err
	}
	if len(namespaceList.Items) == 0 {
		log.Debugf("Namespace %s in tenant %s does not exist, set empty permission", img.namespace, img.tenantID)
		a.Actions = []string{}
		return nil
	}
	namespace := namespaceList.Items[0]
	if namespace.Status.Locked == nil || !*namespace.Status.Locked {
		if u.authenticated && u.userTenantID == img.tenantID {
			permission = "RWM"
		} else if namespace.Spec.Visibility == registry.VisibilityPublic {
			permission = "R"
		}
	}

	log.Debug("Filtered repository authorization", log.Any("resourceAction", a), log.String("requestTenantID", img.tenantID), log.String("userTenantID", u.userTenantID), log.String("username", u.username), log.String("permission", permission))

	a.Actions = permToActions(permission)
	return nil
}

func permToActions(p string) []string {
	var res []string
	if strings.Contains(p, "W") {
		res = append(res, "push")
	}
	if strings.Contains(p, "M") {
		res = append(res, "*")
	}
	if strings.Contains(p, "R") {
		res = append(res, "pull")
	}
	return res
}
