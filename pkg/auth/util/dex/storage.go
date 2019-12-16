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
	"github.com/dexidp/dex/pkg/log"
	dexstorage "github.com/dexidp/dex/storage"
	"github.com/dexidp/dex/storage/etcd"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

// WrapEtcdConfig options for connecting to etcd databases and override connector crud by auth client.
type WrapEtcdConfig struct {
	etcd etcd.Etcd

	authClient authinternalclient.AuthInterface
}

func NewWrapEtcdStorage(etcd etcd.Etcd, authClient authinternalclient.AuthInterface) *WrapEtcdConfig {
	return &WrapEtcdConfig{etcd: etcd, authClient: authClient}
}

// Open creates a new storage implementation backed by Etcd
func (p *WrapEtcdConfig) Open(logger log.Logger) (dexstorage.Storage, error) {
	store, err := p.etcd.Open(logger)
	if err != nil {
		return nil, err
	}
	return &conn{store, p.authClient}, nil
}
