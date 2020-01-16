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

package identityprovider

import (
	"context"

	"github.com/dexidp/dex/connector"
	dexlog "github.com/dexidp/dex/pkg/log"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/auth"
)

// IdentityProvider defines a object that generate a dex connector.
type IdentityProvider interface {

	// Open is used to open a dex connector instance.
	Open(id string, logger dexlog.Logger) (connector.Connector, error)

	// Store generates a identity provider object into storage.
	Store() (*auth.IdentityProvider, error)
}

// IdentityProvidersStore represents identity providers for every tenantID.
var IdentityProvidersStore = make(map[string]IdentityProvider)

// UserGetter is an object that can get the user that match the provided field and label criteria.
type UserGetter interface {
	GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*auth.User, error)
}

// UserLister is an object that can list users that match the provided field and label criteria.
type UserLister interface {
	ListUsers(ctx context.Context, options *metainternal.ListOptions) (*auth.UserList, error)
}

// GroupGetter is an object that can get the group that match the provided field and label criteria.
type GroupGetter interface {
	GetGroup(ctx context.Context, name string, options *metav1.GetOptions) (*auth.Group, error)
}

// GroupLister is an object that can list groups that match the provided field and label criteria.
type GroupLister interface {
	ListGroups(ctx context.Context, options *metainternal.ListOptions) (*auth.GroupList, error)
}
