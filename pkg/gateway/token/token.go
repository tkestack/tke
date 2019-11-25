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

package token

import (
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"net/http"
	"time"
	"tkestack.io/tke/pkg/util/log"
)

const (
	cookieName = "tke"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Token struct {
	ID      string    `json:"i"`
	Refresh string    `json:"r"`
	Expire  time.Time `json:"e"`
}

// RetrieveToken gets the idToken and related information from the cookie
// requested by HTTP.
func RetrieveToken(request *http.Request) (*Token, error) {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	tokenJSON, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Error("Failed to base64 decode cookie value", log.Err(err))
		return nil, err
	}
	var t Token
	if err := json.Unmarshal(tokenJSON, &t); err != nil {
		log.Error("Failed to json decode cookie value", log.Err(err))
		return nil, err
	}
	// validate token
	if t.ID == "" || t.Expire.Before(time.Now()) {
		return nil, fmt.Errorf("invalid token")
	}
	return &t, nil
}

// ResponseToken writes a cookie in HTTP response return according to the given
// OAuth2 token.
func ResponseToken(t *oauth2.Token, writer http.ResponseWriter) error {
	idToken, ok := t.Extra("id_token").(string)
	if !ok {
		log.Error("Failed to extra oauth2 token to id token", log.Any("token", t))
		return fmt.Errorf("failed to extra oauth2 token to id token")
	}

	tokenJSON, err := json.Marshal(Token{
		ID:      idToken,
		Refresh: t.RefreshToken,
		Expire:  t.Expiry,
	})
	if err != nil {
		log.Error("Failed to marshal oauth2 token", log.Err(err))
		return fmt.Errorf("failed to mashal oauth2 token")
	}
	tokenStr := base64.StdEncoding.EncodeToString(tokenJSON)

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   int(time.Until(t.Expiry).Seconds()),
	}
	http.SetCookie(writer, cookie)
	return nil
}

// DeleteCookie to delete cookie in HTTP response. It used to logout.
func DeleteCookie(writer http.ResponseWriter) {
	cookie := http.Cookie{Name: cookieName, Path: "/", MaxAge: -1}
	http.SetCookie(writer, &cookie)
}
