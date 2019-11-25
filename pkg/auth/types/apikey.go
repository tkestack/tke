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

package types

import (
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"time"
)

// Duration implements json marshal func for time.Duration.
//+k8s:openapi-gen=true
type Duration struct {
	time.Duration
}

// UnmarshalJSON parses struct by unmarshaling from either the numeric or string representations.
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	if err != nil {
		return
	}
	d.Duration = time.Duration(id) * time.Second
	return
}

// MarshalJSON supports marshaling to the string representation of the duration.
func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// SignKeys hold encryption and signing keys.
type SignKeys struct {
	// Key for creating and verifying signatures. These may be nil.
	SigningKey    *jose.JSONWebKey
	SigningKeyPub *jose.JSONWebKey
}

// APIKeyReq contains expiration time used to apply the api key.
//+k8s:openapi-gen=true
type APIKeyReq struct {
	//Exipre is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire Duration `json:"expire,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`
}

// APIKeyReqPassword contains userinfo and expiration time used to apply the api key.
//+k8s:openapi-gen=true
type APIKeyReqPassword struct {
	// TenantID for user
	TenantID string `json:"tenantID,omitempty"`

	// UserName
	UserName string `json:"username,omitempty"`

	// Password (encoded by base64)
	Password string `json:"password,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`

	// Expire holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire Duration `json:"expire,omitempty"`
}

// APIKeyData contains expiration time used to apply the api key.
//+k8s:openapi-gen=true
type APIKeyData struct {
	// APIkey is the jwt token used to authenticate user, and contains user info and sign.
	APIkey string `json:"apiKey,omitempty"`

	// Disabled represents whether the apikey has been disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// Deleted represents whether the apikey has been deleted.
	Deleted *bool `json:"deleted,omitempty"`

	// Expired represents whether the apikey has been expired.
	Expired bool `json:"expired,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description,omitempty"`

	// IssueAt is the created time for api key
	IssueAt time.Time `json:"issue_at,omitempty"`

	// ExpireAt is the expire time for api key
	ExpireAt time.Time `json:"expire_at,omitempty"`
}

// APIKeyList is the whole list of APIKeyData.
//+k8s:openapi-gen=true
type APIKeyList struct {
	Items []*APIKeyData `json:"items,omitempty"`
}
