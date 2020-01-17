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
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/ldap"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider/local"
	"tkestack.io/tke/pkg/util/log"

	dexldap "github.com/dexidp/dex/connector/ldap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	oidcidp "tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/registry/identityprovider"
)

// Storage includes storage for signing keys and all sub resources.
type Storage struct {
	*REST
}

// NewStorage returns a Storage object that will work against signing key.
func NewStorage(optsGetter generic.RESTOptionsGetter, authClient authinternalclient.AuthInterface, versionedInformers versionedinformers.SharedInformerFactory) *Storage {
	strategy := identityprovider.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.IdentityProvider{} },
		NewListFunc:              func() runtime.Object { return &auth.IdentityProviderList{} },
		DefaultQualifiedResource: auth.Resource("identityproviders"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		PredicateFunc: identityprovider.MatchAPIKey,
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    identityprovider.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create identityprovider etcd rest storage", log.Err(err))
	}

	return &Storage{&REST{store, authClient, versionedInformers}}
}

// REST implements a RESTStorage for signing keys against etcd.
type REST struct {
	*registry.Store

	authClient         authinternalclient.AuthInterface
	versionedInformers versionedinformers.SharedInformerFactory
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	var idp oidcidp.IdentityProvider
	var err error

	idpObj := obj.(*auth.IdentityProvider)

	log.Info("Create a new identity provider", log.String("type", idpObj.Spec.Type), log.String("name(tenant)", idpObj.Name))
	switch idpObj.Spec.Type {
	case local.ConnectorType:
		idp, err = local.NewDefaultIdentityProvider(idpObj.Name, idpObj.Spec.Administrators, r.versionedInformers)
		if err != nil {
			return nil, errors.NewInternalError(err)
		}
	case ldap.ConnectorType:
		var ldapConfig dexldap.Config
		if err = json.Unmarshal([]byte(idpObj.Spec.Config), &ldapConfig); err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		idp, err = ldap.NewLDAPIdentityProvider(ldapConfig, idpObj.Spec.Administrators, idpObj.Name)
		if err != nil {
			return nil, errors.NewInternalError(err)
		}
	default:
		log.Warn("Identity provider type has not implemented users or groups api", log.String("type", idpObj.Spec.Type))
	}

	result, err := r.Store.Create(ctx, obj, createValidation, options)
	if err == nil && idp != nil {
		oidcidp.IdentityProvidersStore[idpObj.Name] = idp
	}

	return result, err
}
