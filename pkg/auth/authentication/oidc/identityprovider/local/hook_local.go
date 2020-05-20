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

package local

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	genericapiserver "k8s.io/apiserver/pkg/server"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/util/log"
)

type localHookHandler struct {
	authClient authinternalclient.AuthInterface
}

// NewLocalHookHandler creates a new localHookHandler object.
func NewLocalHookHandler(authClient authinternalclient.AuthInterface) genericapiserver.PostStartHookProvider {
	return &localHookHandler{
		authClient: authClient,
	}
}

func (d *localHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "wait-local-sync", func(ctx genericapiserver.PostStartHookContext) error {
		go wait.JitterUntil(func() {
			tenantUserSelector := fields.AndSelectors(
				fields.OneTermEqualSelector("spec.type", ConnectorType),
			)
			conns, err := d.authClient.IdentityProviders().List(context.Background(), v1.ListOptions{FieldSelector: tenantUserSelector.String()})
			if err != nil {
				log.Error("List default idp from registry failed", log.Err(err))
				return
			}

			for _, conn := range conns.Items {
				if _, ok := identityprovider.GetIdentityProvider(conn.Name); ok {
					continue
				}

				idp, err := NewDefaultIdentityProvider(conn.Name, conn.Spec.Administrators, d.authClient)
				if err != nil {
					log.Error("NewDefaultIdentityProvider failed", log.String("idp", conn.Spec.Name), log.Err(err))
					continue
				}

				identityprovider.SetIdentityProvider(conn.Name, idp)
				log.Info("load local identity provider successfully", log.String("idp", conn.Name))
			}
		}, 30*time.Second, 0.0, false, ctx.StopCh)
		return nil
	}, nil
}
