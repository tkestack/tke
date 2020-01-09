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

package ldap

import (
	"encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	dexldap "github.com/dexidp/dex/connector/ldap"
	genericapiserver "k8s.io/apiserver/pkg/server"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/util/log"
)

type ldapHookHandler struct {
	authClient authinternalclient.AuthInterface
}

// NewLdapHookHandler creates a new ldapHookHandler object.
func NewLdapHookHandler(authClient authinternalclient.AuthInterface) genericapiserver.PostStartHookProvider {
	return &ldapHookHandler{
		authClient: authClient,
	}
}

func (d *ldapHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "load-ldap-idp", func(context genericapiserver.PostStartHookContext) error {
		tenantUserSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("spec.type", ConnectorType),
		)
		conns, err := d.authClient.IdentityProviders().List(v1.ListOptions{FieldSelector: tenantUserSelector.String()})
		if err != nil {
			return err
		}

		for _, conn := range conns.Items {
			var ldapConfig dexldap.Config
			err = json.Unmarshal([]byte(conn.Spec.Config), &ldapConfig)
			if err != nil {
				log.Error("Unmarshal idp config failed", log.String("idp", conn.Spec.Name), log.Err(err))
				continue
			}

			idp, err := NewLDAPIdentityProvider(ldapConfig, conn.Name)
			if err != nil {
				log.Error("NewLDAPIdentityProvider failed", log.String("idp", conn.Spec.Name), log.Err(err))
				continue
			}

			identityprovider.IdentityProvidersStore[conn.Name] = idp
			log.Info("load ldap identity provider successfully", log.String("idp", conn.Name))
		}

		return nil
	}, nil
}
