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
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"tkestack.io/tke/pkg/auth/util"

	dexldap "github.com/dexidp/dex/connector/ldap"
	dexstorage "github.com/dexidp/dex/storage"
	"gopkg.in/ldap.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/auth/authentication/oidc/identityprovider"
	"tkestack.io/tke/pkg/util/log"
)

// Referenced and modified from github.com/dexidp/dex/blob/master/connector/ldap/ldap.go

const (
	ConnectorType = "ldap"
)

// identityProvider is the third-party idp that support LDAP.
type identityProvider struct {
	dexldap.Config

	userSearchScope  int
	groupSearchScope int

	tlsConfig *tls.Config

	tenantID string
}

// DefaultIdentityProvider is the default idp for tke local identity login.
func NewLDAPIdentityProvider(c dexldap.Config, tenantID string) (identityprovider.IdentityProvider, error) {
	requiredFields := []struct {
		name string
		val  string
	}{
		{"host", c.Host},
		{"userSearch.baseDN", c.UserSearch.BaseDN},
		{"userSearch.username", c.UserSearch.Username},
		{"userSearch.filter", c.UserSearch.Filter},
	}

	for _, field := range requiredFields {
		if field.val == "" {
			return nil, fmt.Errorf("ldap: missing required field %q", field.name)
		}
	}

	var (
		host string
		err  error
	)
	if host, _, err = net.SplitHostPort(c.Host); err != nil {
		host = c.Host
		if c.InsecureNoSSL {
			c.Host = c.Host + ":389"
		} else {
			c.Host = c.Host + ":636"
		}
	}

	tlsConfig := &tls.Config{ServerName: host, InsecureSkipVerify: c.InsecureSkipVerify}
	if c.RootCA != "" || len(c.RootCAData) != 0 {
		data := c.RootCAData
		if len(data) == 0 {
			var err error
			if data, err = ioutil.ReadFile(c.RootCA); err != nil {
				return nil, fmt.Errorf("ldap: read ca file: %v", err)
			}
		}
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("ldap: no certs found in ca file")
		}
		tlsConfig.RootCAs = rootCAs
	}

	if c.ClientKey != "" && c.ClientCert != "" {
		cert, err := tls.LoadX509KeyPair(c.ClientCert, c.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("ldap: load client cert failed: %v", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	userSearchScope, ok := parseScope(c.UserSearch.Scope)
	if !ok {
		return nil, fmt.Errorf("userSearch.Scope unknown value %q", c.UserSearch.Scope)
	}
	groupSearchScope, ok := parseScope(c.GroupSearch.Scope)
	if !ok {
		return nil, fmt.Errorf("groupSearch.Scope unknown value %q", c.GroupSearch.Scope)
	}

	return &identityProvider{Config: c, userSearchScope: userSearchScope, groupSearchScope: groupSearchScope, tlsConfig: tlsConfig, tenantID: tenantID}, nil
}

func (c *identityProvider) Connector() (*dexstorage.Connector, error) {
	if c.tenantID == "" {
		return nil, fmt.Errorf("must specify tenantID")
	}

	bytes, err := json.Marshal(c.Config)
	if err != nil {
		return nil, fmt.Errorf("mashal ldap config failed: %+v", err)
	}

	return &dexstorage.Connector{
		ID:     c.tenantID,
		Type:   ConnectorType,
		Name:   c.tenantID,
		Config: bytes,
	}, nil
}

func (c *identityProvider) GetUser(ctx context.Context, name string, options *v1.GetOptions) (*auth.User, error) {
	ldapUser := ldap.Entry{}
	if err := c.do(ctx, func(conn *ldap.Conn) error {
		entry, found, err := c.userEntry(conn, name)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		if !found {
			return apierrors.NewNotFound(auth.Resource("user"), name)
		}
		ldapUser = entry
		return nil
	}); err != nil {
		return nil, err
	}

	user, err := c.userFromEntry(ldapUser)
	if err != nil {
		return nil, apierrors.NewInternalError(err)
	}

	return user, nil
}

func (c *identityProvider) ListUsers(ctx context.Context, options *internalversion.ListOptions) (*auth.UserList, error) {
	var ldapUsers []*ldap.Entry

	keyword, limit := util.ParseQueryKeywordAndLimit(options)
	if err := c.do(ctx, func(conn *ldap.Conn) error {
		entries, err := c.usersEntry(conn, keyword, limit)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		ldapUsers = entries
		return nil
	}); err != nil {
		return nil, err
	}

	userList := auth.UserList{}
	for _, entry := range ldapUsers {
		user, err := c.userFromEntry(*entry)
		if err != nil {
			continue
		}
		userList.Items = append(userList.Items, *user)

	}

	return &userList, nil
}

func (c *identityProvider) GetGroup(ctx context.Context, name string, options *v1.GetOptions) (*auth.Group, error) {
	ldapGroup := ldap.Entry{}
	if err := c.do(ctx, func(conn *ldap.Conn) error {
		entry, found, err := c.groupEntry(conn, name)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		if !found {
			return apierrors.NewNotFound(auth.Resource("user"), name)
		}
		ldapGroup = entry
		return nil
	}); err != nil {
		return nil, err
	}

	group, err := c.groupFromEntry(ldapGroup)
	if err != nil {
		return nil, apierrors.NewInternalError(err)
	}

	return group, nil
}

func (c *identityProvider) ListGroups(ctx context.Context, options *internalversion.ListOptions) (*auth.GroupList, error) {
	var ldapGroups []*ldap.Entry

	keyword, limit := util.ParseQueryKeywordAndLimit(options)
	if err := c.do(ctx, func(conn *ldap.Conn) error {
		entries, err := c.groupsEntry(conn, keyword, limit)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		ldapGroups = entries
		return nil
	}); err != nil {
		return nil, err
	}

	groupList := auth.GroupList{}
	for _, entry := range ldapGroups {
		grp, err := c.groupFromEntry(*entry)
		if err != nil {
			continue
		}
		groupList.Items = append(groupList.Items, *grp)

	}

	return &groupList, nil
}

func (c *identityProvider) userEntry(conn *ldap.Conn, username string) (user ldap.Entry, found bool, err error) {
	filter := fmt.Sprintf("(%s=%s)", c.UserSearch.Username, ldap.EscapeFilter(username))
	if c.UserSearch.Filter != "" {
		filter = fmt.Sprintf("(&%s%s)", c.UserSearch.Filter, filter)
	}

	// Initial search.
	req := &ldap.SearchRequest{
		BaseDN: c.UserSearch.BaseDN,
		Filter: filter,
		Scope:  c.userSearchScope,
		// We only need to search for these specific requests.
		Attributes: []string{
			c.UserSearch.Username,
			c.UserSearch.IDAttr,
			c.UserSearch.EmailAttr,
			c.GroupSearch.UserAttr,
		},
	}

	if c.UserSearch.NameAttr != "" {
		req.Attributes = append(req.Attributes, c.UserSearch.NameAttr)
	}

	if c.UserSearch.PreferredUsernameAttrAttr != "" {
		req.Attributes = append(req.Attributes, c.UserSearch.PreferredUsernameAttrAttr)
	}

	log.Infof("performing ldap search %s %s %s",
		req.BaseDN, scopeString(req.Scope), req.Filter)
	resp, err := conn.Search(req)
	if err != nil {
		return ldap.Entry{}, false, fmt.Errorf("ldap: search with filter %q failed: %v", req.Filter, err)
	}

	switch n := len(resp.Entries); n {
	case 0:
		log.Errorf("ldap: no results returned for filter: %q", filter)
		return ldap.Entry{}, false, nil
	case 1:
		user = *resp.Entries[0]
		log.Infof("username %q mapped to entry %s", username, user.DN)
		return user, true, nil
	default:
		return ldap.Entry{}, false, fmt.Errorf("ldap: filter returned multiple (%d) results: %q", n, filter)
	}
}

func (c *identityProvider) usersEntry(conn *ldap.Conn, keyword string, limit int) (user []*ldap.Entry, err error) {
	filter := ""
	if keyword != "" {
		filter = fmt.Sprintf("(%s=*%s*)", c.UserSearch.Username, ldap.EscapeFilter(keyword))
	}
	if c.UserSearch.Filter != "" {
		filter = fmt.Sprintf("(&%s%s)", c.UserSearch.Filter, filter)
	}

	// Initial search.
	req := &ldap.SearchRequest{
		BaseDN: c.UserSearch.BaseDN,
		Filter: filter,
		Scope:  c.userSearchScope,
		// We only need to search for these specific requests.
		Attributes: []string{
			c.UserSearch.Username,
			c.UserSearch.IDAttr,
			c.UserSearch.EmailAttr,
			c.GroupSearch.UserAttr,
		},
	}

	if c.UserSearch.NameAttr != "" {
		req.Attributes = append(req.Attributes, c.UserSearch.NameAttr)
	}

	if c.UserSearch.PreferredUsernameAttrAttr != "" {
		req.Attributes = append(req.Attributes, c.UserSearch.PreferredUsernameAttrAttr)
	}

	log.Infof("performing ldap search %s %s %s",
		req.BaseDN, scopeString(req.Scope), req.Filter)
	resp, err := conn.SearchWithPaging(req, uint32(limit))
	if err != nil {
		return nil, fmt.Errorf("ldap: search with filter %q failed: %v", req.Filter, err)
	}

	return resp.Entries, nil
}

func (c *identityProvider) userFromEntry(user ldap.Entry) (authUser *auth.User, err error) {
	authUser = &auth.User{}
	// If we're missing any attributes, such as email or ID, we want to report
	// an error rather than continuing.
	var missing []string

	if c.UserSearch.Username != "" {
		if authUser.Spec.Name = getAttr(user, c.UserSearch.Username); authUser.Spec.Name == "" {
			missing = append(missing, c.UserSearch.Username)
		} else {
			// ldap id and name is same
			authUser.Spec.ID = authUser.Spec.Name
			authUser.ObjectMeta.Name = authUser.Spec.Name
		}
	}

	if c.UserSearch.PreferredUsernameAttrAttr != "" {
		if authUser.Spec.DisplayName = getAttr(user, c.UserSearch.PreferredUsernameAttrAttr); authUser.Spec.DisplayName == "" {
			authUser.Spec.DisplayName = getAttr(user, c.UserSearch.NameAttr)
		}
	}

	if c.UserSearch.EmailSuffix != "" {
		authUser.Spec.Email = authUser.Spec.Name + "@" + c.UserSearch.EmailSuffix
	} else {
		authUser.Spec.Email = getAttr(user, c.UserSearch.EmailAttr)
	}

	authUser.Spec.TenantID = c.tenantID
	if len(missing) != 0 {
		err := fmt.Errorf("ldap: entry %q missing following required attribute(s): %q", user.DN, missing)
		return nil, err
	}
	return authUser, nil
}

func (c *identityProvider) groupEntry(conn *ldap.Conn, name string) (user ldap.Entry, found bool, err error) {
	filter := fmt.Sprintf("(%s=%s)", c.GroupSearch.NameAttr, ldap.EscapeFilter(name))
	if c.UserSearch.Filter != "" {
		filter = fmt.Sprintf("(&%s%s)", c.GroupSearch.Filter, filter)
	}

	// Initial search.
	req := &ldap.SearchRequest{
		BaseDN: c.GroupSearch.BaseDN,
		Filter: filter,
		Scope:  c.groupSearchScope,
		// We only need to search for these specific requests.
		Attributes: []string{
			c.GroupSearch.GroupAttr,
			c.GroupSearch.NameAttr,
		},
	}

	log.Infof("performing ldap search %s %s %s",
		req.BaseDN, scopeString(req.Scope), req.Filter)
	resp, err := conn.Search(req)
	if err != nil {
		return ldap.Entry{}, false, fmt.Errorf("ldap: search with filter %q failed: %v", req.Filter, err)
	}

	switch n := len(resp.Entries); n {
	case 0:
		log.Errorf("ldap: no results returned for filter: %q", filter)
		return ldap.Entry{}, false, nil
	case 1:
		user = *resp.Entries[0]
		log.Infof("grouop name %q mapped to entry %s", name, user.DN)
		return user, true, nil
	default:
		return ldap.Entry{}, false, fmt.Errorf("ldap: filter returned multiple (%d) results: %q", n, filter)
	}
}

func (c *identityProvider) groupsEntry(conn *ldap.Conn, keyword string, limit int) (user []*ldap.Entry, err error) {
	filter := ""
	if keyword != "" {
		filter = fmt.Sprintf("(%s=*%s*)", c.GroupSearch.NameAttr, ldap.EscapeFilter(keyword))
	}
	if c.UserSearch.Filter != "" {
		filter = fmt.Sprintf("(&%s%s)", c.GroupSearch.Filter, filter)
	}

	// Initial search.
	req := &ldap.SearchRequest{
		BaseDN: c.GroupSearch.BaseDN,
		Filter: filter,
		Scope:  c.groupSearchScope,
		// We only need to search for these specific requests.
		Attributes: []string{
			c.GroupSearch.GroupAttr,
			c.GroupSearch.NameAttr,
		},
	}

	log.Infof("performing ldap search %s %s %s %d",
		req.BaseDN, scopeString(req.Scope), req.Filter, limit)

	var resp *ldap.SearchResult
	if limit != 0 {
		resp, err = conn.SearchWithPaging(req, uint32(limit))
		if err != nil {
			return nil, fmt.Errorf("ldap: search with filter %q failed: %v", req.Filter, err)
		}
	} else {
		resp, err = conn.Search(req)
		if err != nil {
			return nil, fmt.Errorf("ldap: search with filter %q failed: %v", req.Filter, err)
		}
	}

	log.Info("resp", log.Any("resp", resp))

	return resp.Entries, nil
}

func (c *identityProvider) groupFromEntry(group ldap.Entry) (authGroup *auth.Group, err error) {
	authGroup = &auth.Group{}
	// If we're missing any attributes, such as email or ID, we want to report
	// an error rather than continuing.
	var missing []string

	// Fill the identity struct using the attributes from the user entry.
	if authGroup.Spec.ID = getAttr(group, c.GroupSearch.NameAttr); authGroup.Spec.ID == "" {
		missing = append(missing, c.UserSearch.IDAttr)
	} else {
		authGroup.ObjectMeta.Name = authGroup.Spec.ID
		authGroup.Spec.DisplayName = authGroup.Spec.ID
	}

	if c.GroupSearch.GroupAttr != "" {
		members := getAttrs(group, c.GroupSearch.GroupAttr)
		for _, dn := range members {
			name := parseNameFromDN(dn, c.UserSearch.Username)
			if name != "" {
				authGroup.Status.Users = append(authGroup.Status.Users, auth.Subject{
					ID:   name,
					Name: name,
				})
			}
		}
	}
	authGroup.Spec.TenantID = c.tenantID
	if len(missing) != 0 {
		err := fmt.Errorf("ldap: entry %q missing following required attribute(s): %q", group.DN, missing)
		return nil, err
	}

	return authGroup, nil
}

var _ identityprovider.UserGetter = &identityProvider{}

var _ identityprovider.UserLister = &identityProvider{}

var _ identityprovider.GroupGetter = &identityProvider{}
var _ identityprovider.GroupLister = &identityProvider{}

// do initializes a connection to the LDAP directory and passes it to the
// provided function. It then performs appropriate teardown or reuse before
// returning.
func (c *identityProvider) do(ctx context.Context, f func(c *ldap.Conn) error) error {
	var (
		conn *ldap.Conn
		err  error
	)
	switch {
	case c.InsecureNoSSL:
		conn, err = ldap.Dial("tcp", c.Host)
	case c.StartTLS:
		conn, err = ldap.Dial("tcp", c.Host)
		if err != nil {
			return fmt.Errorf("failed to connect: %v", err)
		}
		if err := conn.StartTLS(c.tlsConfig); err != nil {
			return fmt.Errorf("start TLS failed: %v", err)
		}
	default:
		conn, err = ldap.DialTLS("tcp", c.Host, c.tlsConfig)
	}
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// If bindDN and bindPW are empty this will default to an anonymous bind.
	if err := conn.Bind(c.BindDN, c.BindPW); err != nil {
		if c.BindDN == "" && c.BindPW == "" {
			return fmt.Errorf("ldap: initial anonymous bind failed: %v", err)
		}
		return fmt.Errorf("ldap: initial bind for user %q failed: %v", c.BindDN, err)
	}

	return f(conn)
}

func parseScope(s string) (int, bool) {
	// NOTE(ericchiang): ScopeBaseObject doesn't really make sense for us because we
	// never know the user's or group's DN.
	switch s {
	case "", "sub":
		return ldap.ScopeWholeSubtree, true
	case "one":
		return ldap.ScopeSingleLevel, true
	}
	return 0, false
}

func scopeString(i int) string {
	switch i {
	case ldap.ScopeBaseObject:
		return "base"
	case ldap.ScopeSingleLevel:
		return "one"
	case ldap.ScopeWholeSubtree:
		return "sub"
	default:
		return ""
	}
}

func getAttr(e ldap.Entry, name string) string {
	if a := getAttrs(e, name); len(a) > 0 {
		return a[0]
	}
	return ""
}

func getAttrs(e ldap.Entry, name string) []string {
	for _, a := range e.Attributes {
		if a.Name != name {
			continue
		}
		return a.Values
	}
	if name == "DN" {
		return []string{e.DN}
	}
	return nil
}

func parseNameFromDN(s string, nameAttr string) string {
	dn, err := ldap.ParseDN(s)
	if err != nil {
		return ""
	}
	for _, dn := range dn.RDNs {
		for _, at := range dn.Attributes {
			if at.Type == nameAttr {
				return at.Value
			}
		}
	}
	return ""
}
