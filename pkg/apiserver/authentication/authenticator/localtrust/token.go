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
	"encoding/base64"
	"encoding/json"
	"k8s.io/apiserver/pkg/authentication/user"
	"tkestack.io/tke/pkg/util/log"
)

// User defines the user info.
type User struct {
	Name   string              `json:"name"`
	UID    string              `json:"uid"`
	Groups []string            `json:"groups"`
	Extra  map[string][]string `json:"extra"`
}

// GenerateToken to generate token for local trust authenticator.
func GenerateToken(info user.Info) (string, error) {
	u := User{
		Name:   info.GetName(),
		UID:    info.GetUID(),
		Groups: info.GetGroups(),
		Extra:  info.GetExtra(),
	}
	buf, err := json.Marshal(u)
	if err != nil {
		log.Error("Failed to marshal user info", log.Any("user", info), log.Err(err))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func decodeToken(token string) (user.Info, error) {
	buf, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(buf, &u); err != nil {
		return nil, err
	}

	return &user.DefaultInfo{
		Name:   u.Name,
		UID:    u.UID,
		Groups: u.Groups,
		Extra:  u.Extra,
	}, nil
}
