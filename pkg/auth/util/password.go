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

package util

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"tkestack.io/tke/pkg/util/log"
)

// BcryptPassword decodes base64 string and bcrypts password.
func BcryptPassword(password string) (string, error) {
	decodedPasswd, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		log.Error("Base64 Decode password failed", log.Err(err))
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword(decodedPasswd, 10)
	if err != nil {
		log.Error("Bcrypt hash password failed", log.Err(err))
		return "", err
	}
	log.Debug("password", log.ByteString("pwd", decodedPasswd), log.ByteString("hased", hashed))
	return base64.StdEncoding.EncodeToString(hashed), nil
}

// VerifyDecodedPassword verifies password.
func VerifyDecodedPassword(decodedPasswd string, bcryptedPasswd string) error {
	if decodedPasswd == "" && bcryptedPasswd == "" {
		return nil
	}

	if decodedPasswd == "" {
		return fmt.Errorf("input original password is empty")
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(decodedPasswd)
	if err != nil {
		return err
	}

	hashBytes, err := base64.StdEncoding.DecodeString(bcryptedPasswd)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hashBytes, decodedBytes)
}
