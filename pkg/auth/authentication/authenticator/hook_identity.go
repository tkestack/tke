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

package authenticator

import (
	"encoding/base64"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util/log"
)

type adminIdentityHookHandler struct {
	authClient authinternalclient.AuthInterface

	tenantID string
	userName string
	password string
}

// NewAPISigningKeyHookHandler creates a new authnHookHandler object.
func NewAdminIdentityHookHandler(authClient authinternalclient.AuthInterface, tenantID, userName, password string) genericapiserver.PostStartHookProvider {
	return &adminIdentityHookHandler{
		authClient: authClient,
		tenantID:   tenantID,
		userName:   userName,
		password:   password,
	}
}

func (d *adminIdentityHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "generate-default-admin-identity", func(context genericapiserver.PostStartHookContext) error {
		tenantUserSelector := fields.AndSelectors(
			fields.OneTermEqualSelector("spec.tenantID", d.tenantID),
			fields.OneTermEqualSelector("spec.username", d.userName))

		localIdentityList, err := d.authClient.LocalIdentities().List(metav1.ListOptions{FieldSelector: tenantUserSelector.String()})
		if err != nil {
			return err
		}

		if len(localIdentityList.Items) != 0 {
			return nil
		}

		_, err = d.authClient.LocalIdentities().Create(&auth.LocalIdentity{
			Spec: auth.LocalIdentitySpec{
				HashedPassword: base64.StdEncoding.EncodeToString([]byte(d.password)),
				TenantID:       d.tenantID,
				Username:       d.userName,
				DisplayName:    "Administrator",
				Extra: map[string]string{
					"platformadmin": "true",
				},
			},
		})
		if err != nil {
			log.Error("Failed to create the default admin identity", log.Err(err))
			return err
		}
		return nil
	}, nil
}
