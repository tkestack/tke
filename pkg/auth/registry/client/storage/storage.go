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

	dexstorage "github.com/dexidp/dex/storage"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	Client *REST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(_ genericregistry.RESTOptionsGetter, storage dexstorage.Storage) *Storage {
	return &Storage{
		Client: &REST{dexStorage: storage},
	}
}

// REST implements a RESTStorage for configmap against etcd.
type REST struct {
	rest.Storage
	dexStorage dexstorage.Storage
}

func (r *REST) NamespaceScoped() bool {
	return false
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Creater = &REST{}
var _ rest.Scoper = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cli"}
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &auth.Client{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &auth.ClientList{}
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	idp := obj.(*auth.Client)

	cli := toDexClient(idp)
	if err := r.dexStorage.CreateClient(cli); err != nil {
		if err == dexstorage.ErrAlreadyExists {
			return nil, apierrors.NewConflict(auth.Resource("Client"), idp.Name, err)
		}

		return nil, apierrors.NewInternalError(err)
	}

	return idp, nil
}

// Delete enforces life-cycle rules for policy termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	cli, err := r.dexStorage.GetClient(name)
	if err != nil {
		if err == dexstorage.ErrNotFound {
			return nil, false, apierrors.NewNotFound(auth.Resource("client"), name)
		}

		return nil, false, apierrors.NewInternalError(err)
	}

	err = r.dexStorage.DeleteClient(name)
	if err != nil {
		if err == dexstorage.ErrNotFound {
			return nil, false, apierrors.NewNotFound(auth.Resource("client"), name)
		}

		return nil, false, apierrors.NewInternalError(err)
	}
	return fromDexClient(cli), true, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	cli, err := r.dexStorage.GetClient(name)
	if err != nil {
		if err == dexstorage.ErrNotFound {
			return nil, apierrors.NewNotFound(auth.Resource("client"), name)
		}

		return nil, apierrors.NewInternalError(err)
	}

	return fromDexClient(cli), nil
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	clis, err := r.dexStorage.ListClients()
	if err != nil {
		return nil, apierrors.NewInternalError(err)
	}

	var clientList auth.ClientList

	for _, cli := range clis {
		idp := fromDexClient(cli)
		clientList.Items = append(clientList.Items, *idp)
	}

	return &clientList, nil
}

func toDexClient(client *auth.Client) dexstorage.Client {
	return dexstorage.Client{
		ID:           client.Spec.ID,
		Secret:       client.Spec.Secret,
		RedirectURIs: client.Spec.RedirectUris,
		TrustedPeers: client.Spec.TrustedPeers,
		Public:       client.Spec.Public,
		Name:         client.Spec.Name,
		LogoURL:      client.Spec.LogoURL,
	}
}

func fromDexClient(client dexstorage.Client) *auth.Client {
	return &auth.Client{
		ObjectMeta: metav1.ObjectMeta{
			Name: client.ID,
		},
		Spec: auth.ClientSpec{
			ID:           client.ID,
			Secret:       client.Secret,
			RedirectUris: client.RedirectURIs,
			TrustedPeers: client.TrustedPeers,
			Public:       client.Public,
			Name:         client.Name,
			LogoURL:      client.LogoURL,
		},
	}
}
