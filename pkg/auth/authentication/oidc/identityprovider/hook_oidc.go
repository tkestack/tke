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
	"crypto/tls"
	"fmt"
	"net/http"

	gooidc "github.com/coreos/go-oidc"
	dexserver "github.com/dexidp/dex/server"
	dexstorage "github.com/dexidp/dex/storage"
	"gopkg.in/square/go-jose.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"

	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/authentication/authenticator"
	"tkestack.io/tke/pkg/util/log"
)

type dexHookHandler struct {
	handler    *DexHander
	dexConfig  *dexserver.Config
	dexStorage dexstorage.Storage
	authClient authinternalclient.AuthInterface

	publicAddress   string
	internalAddress string
	tokenAuthn      *authenticator.TokenAuthenticator
	ctx             context.Context
}

func NewDexHookHandler(ctx context.Context, authClient authinternalclient.AuthInterface, config *dexserver.Config, storage dexstorage.Storage, handler *DexHander,
	publicAddress string, internalAddress string, tokenAuthn *authenticator.TokenAuthenticator) genericapiserver.PostStartHookProvider {
	return &dexHookHandler{
		handler:         handler,
		dexConfig:       config,
		dexStorage:      storage,
		authClient:      authClient,
		publicAddress:   publicAddress,
		internalAddress: internalAddress,
		tokenAuthn:      tokenAuthn,
		ctx:             ctx,
	}
}

// PostStartHook provides a function that is called after the server has started.
func (d *dexHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "create-dex-server", func(_ genericapiserver.PostStartHookContext) error {
		log.Info("start create dex server")
		// Ensure all identity providers defined exists in dex.
		for tenantID, idp := range IdentityProvidersStore {
			idp, err := idp.Store()
			if err != nil {
				log.Errorf("Get connector for tenant failed", log.String("tenantID", tenantID), log.Err(err))
				continue
			}

			// if conn is nil, not create into dexStorage
			if idp == nil {
				continue
			}

			if _, err := d.authClient.IdentityProviders().Create(idp); err != nil && !errors.IsAlreadyExists(err) {
				log.Error("Create idp for tenant failed", log.String("tenantID", tenantID), log.Any("idp", *idp), log.Err(err))
			}
		}

		// Ensure there is at least one connector available for dex
		conns, err := d.dexStorage.ListConnectors()
		if err != nil {
			return err
		}

		if len(conns) == 0 {
			return fmt.Errorf("create connectors failed")
		}
		server, err := dexserver.NewServer(d.ctx, *d.dexConfig)
		if err != nil {
			return err
		}
		d.handler.handler = server
		return d.registerIDTokenVerifier()
	}, nil
}

func (d *dexHookHandler) registerIDTokenVerifier() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	ctx := gooidc.ClientContext(d.ctx, client)
	idTokenVerifier, err := oidc.NewIDTokenVerifier(ctx, d.internalAddress, d.publicAddress, &gooidc.Config{
		SkipClientIDCheck: true,
		SupportedSigningAlgs: []string{
			string(jose.RS256),
		},
	})

	if err != nil {
		log.Error("Failed to create the oidc verifier", log.Err(err))
		return err
	}

	d.tokenAuthn.IDTokenVerifier = idTokenVerifier
	return nil
}
