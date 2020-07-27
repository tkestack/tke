/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package ssh

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func MakePrivateKeySignerFromFile(key string) (ssh.Signer, error) {
	// Create an actual signer.
	buffer, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, fmt.Errorf("error reading SSH key %s: '%v'", key, err)
	}
	return MakePrivateKeySigner(buffer, nil)
}

func MakePrivateKeySigner(privateKey []byte, passPhrase []byte) (ssh.Signer, error) {
	var signer ssh.Signer
	var err error
	if passPhrase == nil {
		signer, err = ssh.ParsePrivateKey(privateKey)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, passPhrase)
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing SSH key: '%v'", err)
	}
	return signer, nil
}

func ParsePublicKeyFromFile(keyFile string) (*rsa.PublicKey, error) {
	buffer, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading SSH key %s: '%v'", keyFile, err)
	}
	keyBlock, _ := pem.Decode(buffer)
	if keyBlock == nil {
		return nil, fmt.Errorf("error parsing SSH key %s: 'invalid PEM format'", keyFile)
	}
	key, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing SSH key %s: '%v'", keyFile, err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("SSH key could not be parsed as rsa public key")
	}
	return rsaKey, nil
}
