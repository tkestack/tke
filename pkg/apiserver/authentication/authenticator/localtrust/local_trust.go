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

package localtrust

import (
	"context"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"tkestack.io/tke/pkg/apiserver/filter"
)

type localTrustAuthenticator struct {
}

// NewLocalTrustAuthenticator returns an authenticator.Request.
func NewLocalTrustAuthenticator() authenticator.Token {
	return &localTrustAuthenticator{}
}

func (a *localTrustAuthenticator) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	local := filter.LocalFrom(ctx)
	if !local {
		return nil, false, nil
	}

	user, err := decodeToken(token)
	if err != nil {
		return nil, false, nil
	}

	return &authenticator.Response{
		User: user,
	}, true, nil
}
