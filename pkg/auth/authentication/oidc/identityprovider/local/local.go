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

	"k8s.io/apimachinery/pkg/labels"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dexidp/dex/connector"
	dexlog "github.com/dexidp/dex/pkg/log"
	dexstorage "github.com/dexidp/dex/storage"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"tkestack.io/tke/api/auth"
	authv1 "tkestack.io/tke/api/auth/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	authv1lister "tkestack.io/tke/api/client/listers/auth/v1"
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

// DefaultIdentityProvider is the default idp for tke local identity login.
type DefaultIdentityProvider struct {
	tenantID            string
	localIdentityLister authv1lister.LocalIdentityLister
	localGroupLister    authv1lister.LocalGroupLister
}

func NewDefaultIdentityProvider(tenantID string, versionInformers versionedinformers.SharedInformerFactory) identityprovider.IdentityProvider {
	return &DefaultIdentityProvider{
		tenantID:            tenantID,
		localIdentityLister: versionInformers.Auth().V1().LocalIdentities().Lister(),
		localGroupLister:    versionInformers.Auth().V1().LocalGroups().Lister(),
	}
}

// Open returns a strategy for logging in through TKE
func (c *DefaultIdentityProvider) Open(id string, logger dexlog.Logger) (
	connector.Connector, error) {

	if authClient == nil {
		return nil, fmt.Errorf("kubernetes client config is nil")
	}

	return &localConnector{authClient: authClient, tenantID: id}, nil
}

func (c *DefaultIdentityProvider) Connector() (*dexstorage.Connector, error) {
	if c.tenantID == "" {
		return nil, fmt.Errorf("must specify tenantID")
	}

	return &dexstorage.Connector{
		Type:   ConnectorType,
		ID:     c.tenantID,
		Name:   c.tenantID,
		Config: []byte("{}"),
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
	localIdentity, err := util.GetLocalIdentity(authClient, p.tenantID, username)
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
	groups, err := util.GetGroupsForUser(authClient, localIdentity.ObjectMeta.Name)
	if err == nil {
		for _, g := range groups.Items {
			ident.Groups = append(ident.Groups, g.ObjectMeta.Name)
		}
	}

	ident.Email = localIdentity.Spec.Email

	if emailVerified, ok := localIdentity.Spec.Extra["emailVerified"]; ok {
		ident.EmailVerified, _ = strconv.ParseBool(emailVerified)
	}

	log.Info("Check user login success", log.Any("User info", ident))

	return ident, true, nil
}

func (p *localConnector) Refresh(ctx context.Context, s connector.Scopes, identity connector.Identity) (connector.Identity, error) {
	// If the user has been deleted, the refresh token will be rejected.
	ident, err := util.GetLocalIdentity(p.authClient, p.tenantID, identity.Username)
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
func (c *DefaultIdentityProvider) GetUser(ctx context.Context, name string, options *metav1.GetOptions) (*auth.User, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	localIdentity, err := c.localIdentityLister.Get(name)
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
func (c *DefaultIdentityProvider) ListUsers(ctx context.Context, options *metainternal.ListOptions) (*auth.UserList, error) {
	keyword, limit := util.ParseQueryKeywordAndLimit(options)

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	allList, err := c.localIdentityLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var localIdentityList []*authv1.LocalIdentity
	for i, item := range allList {
		if item.Spec.TenantID == c.tenantID {
			localIdentityList = append(localIdentityList, allList[i])
		}
	}

	if keyword != "" {
		var newList []*authv1.LocalIdentity
		for i, val := range localIdentityList {
			if strings.Contains(val.Name, keyword) || strings.Contains(val.Spec.Username, keyword) || strings.Contains(val.Spec.DisplayName, keyword) {
				newList = append(newList, localIdentityList[i])
			}
		}
		localIdentityList = newList
	}

	items := localIdentityList[0:min(len(localIdentityList), limit)]

	userList := auth.UserList{}
	for _, item := range items {
		user := convertToUser(item)
		userList.Items = append(userList.Items, user)
	}

	return &userList, nil
}

// Get is an object that can get the user that match the provided field and label criteria.
func (c *DefaultIdentityProvider) GetGroup(ctx context.Context, name string, options *metav1.GetOptions) (*auth.Group, error) {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	localGroup, err := c.localGroupLister.Get(name)
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
func (c *DefaultIdentityProvider) ListGroups(ctx context.Context, options *metainternal.ListOptions) (*auth.GroupList, error) {

	keyword, limit := util.ParseQueryKeywordAndLimit(options)

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID != "" && tenantID != c.tenantID {
		return nil, apierrors.NewBadRequest("must in the same tenant")
	}

	allList, err := c.localGroupLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var localGroupList []*authv1.LocalGroup
	for i, item := range allList {
		if item.Spec.TenantID == c.tenantID {
			localGroupList = append(localGroupList, allList[i])
		}
	}

	if keyword != "" {
		var newList []*authv1.LocalGroup
		for i, val := range localGroupList {
			if strings.Contains(val.Name, keyword) || strings.Contains(val.Spec.DisplayName, keyword) {
				newList = append(newList, localGroupList[i])
			}
		}
		localGroupList = newList
	}

	if limit > 0 {
		localGroupList = localGroupList[0:min(len(localGroupList), limit)]
	}

	groupList := auth.GroupList{}
	for _, item := range localGroupList {
		group := convertToGroup(item)
		groupList.Items = append(groupList.Items, group)
	}

	return &groupList, nil
}

func convertToUser(localIdentity *authv1.LocalIdentity) auth.User {
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

func convertToGroup(localGroup *authv1.LocalGroup) auth.Group {
	return auth.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: localGroup.ObjectMeta.Name,
		},
		Spec: auth.GroupSpec{
			ID:          localGroup.ObjectMeta.Name,
			DisplayName: localGroup.Spec.DisplayName,
			TenantID:    localGroup.Spec.TenantID,
			Description: localGroup.Spec.TenantID,
		},
		Status: auth.GroupStatus{
			Users: fromV1Subject(localGroup.Status.Users),
		},
	}
}

func fromV1Subject(v1Subjects []authv1.Subject) []auth.Subject {
	var subjects []auth.Subject

	for _, sub := range v1Subjects {
		subjects = append(subjects, auth.Subject{
			ID:   sub.ID,
			Name: sub.Name,
		})
	}

	return subjects
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
