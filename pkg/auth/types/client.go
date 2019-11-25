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

// Client represents an OAuth2 client.
//+k8s:openapi-gen=true
type Client struct {
	ID           string   `json:"id,omitempty"`
	Secret       string   `json:"secret,omitempty"`
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// TrustedPeers are a list of peers which can issue tokens on this client's behalf using the dynamic "oauth2:server:client_id:(client_id)" scope.
	TrustedPeers []string `json:"trusted_peers,omitempty"`
	// Public clients must use either use a redirectURL 127.0.0.1:X or "urn:ietf:wg:oauth:2.0:oob".
	Public  bool   `json:"public,omitempty"`
	Name    string `json:"name,omitempty"`
	LogoURL string `json:"logo_url,omitempty"`
}

// ClientList is the whole list of OAuth2 client.
//+k8s:openapi-gen=true
type ClientList struct {
	// List of policies.
	Items []*Client `json:"items,omitempty"`
}
