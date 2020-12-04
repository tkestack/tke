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

package installer

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	certutil "k8s.io/client-go/util/cert"
	"tkestack.io/tke/pkg/util/pkiutil"
)

var (
	tke = newTKE()
)

func newTKE() *TKE {
	return &TKE{
		namespace: namespace,
	}
}

func TestTKE_validateCertAndKey(t *testing.T) {
	caCert, caKey, _ := pkiutil.NewCertificateAuthority(&certutil.Config{
		CommonName:   "TKE",
		Organization: []string{"Tencent"},
	})

	cert, key, _ := pkiutil.NewCertAndKey(caCert, caKey, &certutil.Config{
		CommonName:   "TKE-SERVER",
		Organization: []string{"Tencent"},
		AltNames: certutil.AltNames{
			DNSNames: []string{"*.tke.com", "*.registry.tke.com"},
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	})

	err := tke.validateCertAndKey(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}),
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}),
		"console.tke.com",
	)
	assert.True(t, err == nil)
}
