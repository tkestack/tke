/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package storage

import (
	"context"
	"strings"

	"github.com/casbin/casbin/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// ProjectREST implements the REST endpoint, list policies bound to the user.
type ProjectREST struct {
	userRest   *REST
	authClient authinternalclient.AuthInterface
	enforcer   *casbin.SyncedEnforcer
}

var _ = rest.Lister(&ProjectREST{})

// NewList returns an empty object that can be used with the List call.
func (r *ProjectREST) NewList() runtime.Object {
	return &auth.ProjectBelongs{}
}

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *ProjectREST) New() runtime.Object {
	return &auth.ProjectBelongs{}
}

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (r *ProjectREST) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	// TODO: convert role list to table
	tableConvertor := rest.NewDefaultTableConvertor(auth.Resource("projects"))
	return tableConvertor.ConvertToTable(ctx, object, tableOptions)
}

func (r *ProjectREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("unable to get request info from context")
	}

	userID := requestInfo.Name
	// todo: confirm this
	//var tenantID string
	//if _, t := authentication.GetUsernameAndTenantID(ctx); t == "" {
	//	tenantID = t
	//} else {
	//	tenantID = filter.TenantIDFrom(ctx)
	//}

	obj, err := r.userRest.Get(ctx, userID, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	user := obj.(*auth.User)

	managed := make(map[string]auth.ExtraValue)
	memberd := make(map[string]auth.ExtraValue)
	projectOwner := util.ProjectOwnerPolicyID(user.Spec.TenantID)

	rules := r.enforcer.GetFilteredGroupingPolicy(0, util.UserKey(user.Spec.TenantID, user.Spec.Name))
	for _, r := range rules {
		if len(r) != 3 {
			log.Warn("invalid rule", log.Strings("rule", r))
			continue
		}
		prj := r[2]
		role := r[1]

		if strings.HasPrefix(prj, "prj-") {
			if role == projectOwner {
				managed[prj] = append(managed[prj], role)
			}
			memberd[prj] = append(memberd[prj], role)
		}
	}

	return &auth.ProjectBelongs{
		TenantID:        user.Spec.TenantID,
		ManagedProjects: managed,
		MemberdProjects: memberd,
	}, nil
}
