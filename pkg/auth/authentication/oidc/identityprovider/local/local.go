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

package local

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dexidp/dex/connector"
	dexlog "github.com/dexidp/dex/pkg/log"
	dexstorage "github.com/dexidp/dex/storage"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	storageoptions "tkestack.io/tke/pkg/apiserver/storage/options"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/util/log"
)

var (
	// TkeConnectorType type and id
	TkeConnectorType = "tke"
)

// Config holds the configuration parameters for tke local connector login.
type Config struct {
	EtcdOpts storageoptions.ETCDClientOptions
}

// Open returns a strategy for logging in through TKE
func (c *Config) Open(id string, logger dexlog.Logger) (
	connector.Connector, error) {
	client, err := c.EtcdOpts.NewClient()
	if err != nil {
		return nil, err
	}

	return &localIdentityProvider{identityStore: localidentity.NewLocalIdentity(client), tenantID: id}, nil
}

// NewLocalConnector creates a demo tke connector when there is no connector in backend.
func NewLocalConnector(etcdOpts *storageoptions.ETCDClientOptions, tenantID string) (*dexstorage.Connector, error) {
	bytes, err := json.Marshal(Config{
		*etcdOpts,
	})
	if err != nil {
		return nil, err
	}

	return &dexstorage.Connector{
		Type:   TkeConnectorType,
		ID:     tenantID,
		Name:   tenantID,
		Config: bytes,
	}, nil
}

type localIdentityProvider struct {
	tenantID      string
	identityStore *localidentity.Storage
}

func (p *localIdentityProvider) Prompt() string {
	return "UserName"
}

func (p *localIdentityProvider) Login(ctx context.Context, scopes connector.Scopes, username, password string) (connector.Identity, bool, error) {
	ident := connector.Identity{}
	if len(username) == 0 {
		return ident, false, nil
	}

	log.Debug("Check user login", log.String("tenantID", p.tenantID), log.String("username", username), log.String("password", password))
	localIdentity, err := p.identityStore.Get(p.tenantID, username)
	if err != nil {
		log.Error("Get user failed", log.String("user", username), log.Err(err))
		return ident, false, nil
	}

	hashBytes, err := base64.StdEncoding.DecodeString(localIdentity.Spec.HashedPassword)
	if err != nil {
		log.Error("Parse hash password failed", log.String("hashedPassword", localIdentity.Spec.HashedPassword), log.Err(err))
		return ident, false, nil
	}
	if err := bcrypt.CompareHashAndPassword(hashBytes, []byte(password)); err != nil {
		log.Error("Invalid password", log.ByteString("input password", []byte(password)), log.ByteString("store password", hashBytes))
		return ident, false, nil
	}

	extra := map[string]string{
		oidc.TenantIDKey: localIdentity.Spec.TenantID,
	}
	if localIdentity.Status != nil {
		extra["status"] = strconv.FormatBool(localIdentity.Status.Locked)
	}

	if ident.ConnectorData, err = json.Marshal(extra); err != nil {
		log.Error("Marshal extra data failed", log.Err(err))
		return ident, false, nil
	}

	ident.UserID = localIdentity.UID
	ident.Username = localIdentity.Name
	ident.Groups = localIdentity.Spec.Groups

	if email, ok := localIdentity.Spec.Extra["email"]; ok {
		ident.Email = email
	}

	if emailVerified, ok := localIdentity.Spec.Extra["emailVerified"]; ok {
		ident.EmailVerified, _ = strconv.ParseBool(emailVerified)
	}

	log.Info("Check user login success", log.Any("User info", ident))

	return ident, true, nil
}

func (p *localIdentityProvider) Refresh(ctx context.Context, s connector.Scopes, identity connector.Identity) (connector.Identity, error) {
	// If the user has been deleted, the refresh token will be rejected.
	ident, err := p.identityStore.Get(p.tenantID, identity.Username)
	if err != nil {
		if err == dexstorage.ErrNotFound {
			return connector.Identity{}, errors.New("user not found")
		}
		return connector.Identity{}, fmt.Errorf("get user faild: %v", err)
	}

	// User removed but a new user with the same name exists.
	if ident.UID != identity.UserID {
		return connector.Identity{}, errors.New("user not found")
	}

	return identity, nil
}
