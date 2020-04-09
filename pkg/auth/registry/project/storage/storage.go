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

	"github.com/casbin/casbin/v2"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	Project *REST
	User    *UserREST
	Group   *GroupREST
	Policy  *PolicyREST

	Binding   *BindingREST
	UnBinding *UnBindingREST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(_ genericregistry.RESTOptionsGetter, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer) *Storage {
	return &Storage{
		Project:   &REST{},
		User:      &UserREST{&BindingREST{authClient}, &UnBindingREST{authClient}, authClient},
		Group:     &GroupREST{&BindingREST{authClient}, &UnBindingREST{authClient}, authClient},
		Policy:    &PolicyREST{authClient},
		Binding:   &BindingREST{authClient},
		UnBinding: &UnBindingREST{authClient},
	}
}

// REST implements a RESTStorage for configmap against etcd.
type REST struct {
	rest.Storage
}

func (r *REST) NamespaceScoped() bool {
	return false
}

var _ rest.Scoper = &REST{}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &auth.ProjectBelongs{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &auth.ProjectBelongs{}
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return &auth.ProjectBelongs{}, nil
}

func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {

	return nil, nil
}
