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
	"strconv"
	"strings"

	"github.com/dexidp/dex/connector"
	dexlog "github.com/dexidp/dex/pkg/log"
	dexserver "github.com/dexidp/dex/server"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// Local connectorType type
	ConnectorType = "tke"
)

var (
	authClient authinternalclient.AuthInterface
)

func init() {
	// create dex local identity provider for tke connector.
	dexserver.ConnectorsConfig[ConnectorType] = func() dexserver.ConnectorConfig {
		return new(identityProvider)
	}
}

// identityProvider is the default idp for tke local identity login.
type identityProvider struct {
	tenantID       string
	administrators []string

	authClient authinternalclient.AuthInterface
}

func NewDefaultIdentityProvider(tenantID string, administrators []string, authClient authinternalclient.AuthInterface) (identityprovider.IdentityProvider, error) {
	return &identityProvider{
		tenantID:       tenantID,
		administrators: administrators,
		authClient:     authClient,
	}, nil
}

// Open returns a strategy for logging in through TKE
func (c *identityProvider) Open(id string, logger dexlog.Logger) (
	connector.Connector, error) {

	if authClient == nil {
		return nil, fmt.Errorf("kubernetes client config is nil")
	}

	return &localConnector{authClient: authClient, tenantID: id}, nil
}

func (c *identityProvider) Store() (*auth.IdentityProvider, error) {
	if c.tenantID == "" {
		return nil, fmt.Errorf("must specify tenantID")
	}
	return &auth.IdentityProvider{
		ObjectMeta: v1.ObjectMeta{Name: c.tenantID},
		Spec: auth.IdentityProviderSpec{
			Name:           c.tenantID,
			Type:           ConnectorType,
			Administrators: c.administrators,
			Config:         "{}",
		},
	}, nil

}

func SetupRestClient(authInterface authinternalclient.AuthInterface) {
	authClient = authInterface
}

type localConnector struct {
	tenantID   string
	authClient authinternalclient.AuthInterface
}

func (p *localConnector) Prompt() string {
	return "Username"
}

func (p *localConnector) Login(ctx context.Context, scopes connector.Scopes, username, password string) (connector.Identity, bool, error) {
	ident := connector.Identity{}
	if len(username) == 0 {
		return ident, false, nil
	}

	log.Debug("Check user login", log.String("tenantID", p.tenantID), log.String("username", username), log.String("password", password))
	localIdentity, err := util.GetLocalIdentity(ctx, authClient, p.tenantID, username)
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

	extra["status"] = strconv.FormatBool(localIdentity.Status.Locked)

	if ident.ConnectorData, err = json.Marshal(extra); err != nil {
		log.Error("Marshal extra data failed", log.Err(err))
		return ident, false, nil
	}

	ident.UserID = localIdentity.ObjectMeta.Name
	ident.Username = localIdentity.Spec.Username
	groups, err := util.GetGroupsForUser(ctx, authClient, localIdentity.ObjectMeta.Name)
	if err == nil {
		for _, g := range groups.Items {
			ident.Groups = append(ident.Groups, g.ObjectMeta.Name)
		}
	}

	ident.Email = localIdentity.Spec.Email
	ident.PreferredUsername = localIdentity.Spec.DisplayName
	if emailVerified, ok := localIdentity.Spec.Extra["emailVerified"]; ok {
		ident.EmailVerified, _ = strconv.ParseBool(emailVerified)
	}

	log.Info("Check user login success", log.Any("User info", ident))

	return ident, true, nil
}

func (p *localConnector) Refresh(ctx context.Context, s connector.Scopes, identity connector.Identity) (connector.Identity, error) {
	// If the user has been deleted, the refresh token will be rejected.
	ident, err := util.GetLocalIdentity(ctx, p.authClient, p.tenantID, identity.Username)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return connector.Identity{}, errors.New("user not found")
		}
		return connector.Identity{}, fmt.Errorf("get user faild: %v", err)
	}

	// User removed but a new user with the same name exists.
	if ident.ObjectMeta.Name != identity.UserID {
		return connector.Identity{}, errors.New("user not found")
	}

	return identity, nil
}

// Get is an object that can get the user that match the provided field and label criteria.
func (c *identityProvider) GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*auth.User, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	localIdentity, err := c.authClient.LocalIdentities().Get(ctx, name, *options)
	if err != nil {
		return nil, err
	}

	if localIdentity.Spec.TenantID != c.tenantID {
		return nil, apierrors.NewNotFound(auth.Resource("user"), name)
	}

	user := convertToUser(localIdentity)
	return &user, nil
}

// List is an object that can list users that match the provided field and label criteria.
func (c *identityProvider) ListUsers(ctx context.Context, options *metainternal.ListOptions) (*auth.UserList, error) {
	keyword, limit := util.ParseQueryKeywordAndLimit(options)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	v1Opt := util.PredicateV1ListOptions(c.tenantID, options)
	localIdentityList, err := c.authClient.LocalIdentities().List(ctx, *v1Opt)
	if err != nil {
		return nil, err
	}

	if keyword != "" {
		var newList []auth.LocalIdentity
		for _, val := range localIdentityList.Items {
			if strings.Contains(val.Name, keyword) || strings.Contains(val.Spec.Username, keyword) || strings.Contains(val.Spec.DisplayName, keyword) {
				newList = append(newList, val)
			}
		}

		localIdentityList.Items = newList
	}

	if limit > 0 {
		localIdentityList.Items = localIdentityList.Items[0:min(len(localIdentityList.Items), limit)]
	}

	userList := auth.UserList{}
	for _, item := range localIdentityList.Items {
		user := convertToUser(&item)
		userList.Items = append(userList.Items, user)
	}

	return &userList, nil
}

// Get is an object that can get the user that match the provided field and label criteria.
func (c *identityProvider) GetGroup(ctx context.Context, name string, options *metav1.GetOptions) (*auth.Group, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	localGroup, err := c.authClient.LocalGroups().Get(ctx, name, *options)
	if err != nil {
		return nil, err
	}

	if localGroup.Spec.TenantID != c.tenantID {
		return nil, apierrors.NewNotFound(auth.Resource("group"), name)
	}

	group := convertToGroup(localGroup)
	return &group, nil
}

// List is an object that can list users that match the provided field and label criteria.
func (c *identityProvider) ListGroups(ctx context.Context, options *metainternal.ListOptions) (*auth.GroupList, error) {
	keyword, limit := util.ParseQueryKeywordAndLimit(options)
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	v1Opt := util.PredicateV1ListOptions(c.tenantID, options)

	localGroupList, err := c.authClient.LocalGroups().List(ctx, *v1Opt)
	if err != nil {
		return nil, err
	}
	if keyword != "" {
		var newList []auth.LocalGroup
		for _, val := range localGroupList.Items {
			if strings.Contains(val.Name, keyword) || strings.Contains(val.Spec.DisplayName, keyword) {
				newList = append(newList, val)
			}
		}
		localGroupList.Items = newList
	}

	if limit > 0 {
		localGroupList.Items = localGroupList.Items[0:min(len(localGroupList.Items), limit)]
	}

	groupList := auth.GroupList{}
	for _, item := range localGroupList.Items {
		group := convertToGroup(&item)
		groupList.Items = append(groupList.Items, group)
	}

	return &groupList, nil
}

func convertToUser(localIdentity *auth.LocalIdentity) auth.User {
	return auth.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: localIdentity.ObjectMeta.Name,
		},
		Spec: auth.UserSpec{
			ID:          localIdentity.ObjectMeta.Name,
			Name:        localIdentity.Spec.Username,
			DisplayName: localIdentity.Spec.DisplayName,
			Email:       localIdentity.Spec.Email,
			PhoneNumber: localIdentity.Spec.PhoneNumber,
			TenantID:    localIdentity.Spec.TenantID,
			Extra:       localIdentity.Spec.Extra,
		},
	}
}

func convertToGroup(localGroup *auth.LocalGroup) auth.Group {
	return auth.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: localGroup.ObjectMeta.Name,
		},
		Spec: auth.GroupSpec{
			ID:          localGroup.ObjectMeta.Name,
			DisplayName: localGroup.Spec.DisplayName,
			TenantID:    localGroup.Spec.TenantID,
			Description: localGroup.Spec.Description,
		},
		Status: auth.GroupStatus{
			Users: localGroup.Status.Users,
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
