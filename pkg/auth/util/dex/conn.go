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

package dex

import (
	"context"
	"fmt"

	dexstorage "github.com/dexidp/dex/storage"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

type conn struct {
	dexstorage.Storage

	authClient authinternalclient.AuthInterface
}

func (c *conn) CreateConnector(connector dexstorage.Connector) error {
	return fmt.Errorf("not support CreateConnector, please use identityproviders api")
}

func (c *conn) GetConnector(id string) (conn dexstorage.Connector, err error) {
	idp, err := c.authClient.IdentityProviders().Get(context.Background(), id, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return conn, dexstorage.ErrNotFound
		}
		return conn, err
	}

	return toDexConnector(idp), nil
}

func (c *conn) UpdateConnector(id string, updater func(s dexstorage.Connector) (dexstorage.Connector, error)) error {
	current, err := c.authClient.IdentityProviders().Get(context.Background(), id, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return dexstorage.ErrNotFound
		}
		return err
	}

	administrators := current.Spec.Administrators
	currConn := toDexConnector(current)
	updatedConn, err := updater(currConn)
	if err != nil {
		return err
	}

	updated := fromDexConnector(updatedConn)

	current.Spec = updated.Spec
	current.Spec.Administrators = administrators

	_, err = c.authClient.IdentityProviders().Update(context.Background(), current, metav1.UpdateOptions{})
	return err
}

func (c *conn) DeleteConnector(id string) error {
	err := c.authClient.IdentityProviders().Delete(context.Background(), id, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return dexstorage.ErrNotFound
		}
		return err
	}

	return nil
}

func (c *conn) ListConnectors() (connectors []dexstorage.Connector, err error) {
	idpList, err := c.authClient.IdentityProviders().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, v := range idpList.Items {
		connectors = append(connectors, toDexConnector(&v))
	}
	return connectors, nil
}

func toDexConnector(idp *auth.IdentityProvider) dexstorage.Connector {
	return dexstorage.Connector{
		ID:              idp.Name,
		Name:            idp.Spec.Name,
		Type:            idp.Spec.Type,
		ResourceVersion: idp.ResourceVersion,
		Config:          []byte(idp.Spec.Config),
	}
}

func fromDexConnector(conn dexstorage.Connector) *auth.IdentityProvider {
	return &auth.IdentityProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:            conn.ID,
			ResourceVersion: conn.ResourceVersion,
		},
		Spec: auth.IdentityProviderSpec{
			Name:   conn.Name,
			Type:   conn.Type,
			Config: string(conn.Config),
		},
	}
}
