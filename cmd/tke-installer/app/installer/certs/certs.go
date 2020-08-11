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

package certs

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"net"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/pkg/util/files"
	"tkestack.io/tke/pkg/util/pkiutil"
)

var (
	components = []string{
		"tke-platform-api",
		"tke-business-api",
		"tke-notify-api",
		"tke-auth-api",
		"tke-gateway",
		"tke-monitor-api",
		"tke-registry-api",
		"tke-logagent-api",
		"tke-audit-api",
		"tke-application-api",
	}
)

func Generate(dnsNames []string, ips []net.IP, dir string) error {
	caCert, caKey, err := generateRootCA()
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.CACrtFileBaseName, pkiutil.EncodeCertPEM(caCert), 0644)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.CAKeyFileBaseName, pkiutil.EncodePrivateKeyPEM(caKey), 0644)
	if err != nil {
		return err
	}

	serverCert, serverKey, err := generateServerCertKey(caCert, caKey, dnsNames, ips)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.ServerCrtFileBaseName, pkiutil.EncodeCertPEM(serverCert), 0644)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.ServerKeyFileBaseName, pkiutil.EncodePrivateKeyPEM(serverKey), 0644)
	if err != nil {
		return err
	}

	adminCert, adminKey, err := generateAdminCertKey(caCert, caKey)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.AdminCrtFileBaseName, pkiutil.EncodeCertPEM(adminCert), 0644)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, constants.AdminKeyFileBaseName, pkiutil.EncodePrivateKeyPEM(adminKey), 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateRootCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	config := &certutil.Config{
		CommonName:   "TKE",
		Organization: []string{"Tencent"},
	}
	cert, key, err := pkiutil.NewCertificateAuthority(config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create self-signed certificate")
	}

	return cert, key, nil
}

func generateServerCertKey(caCert *x509.Certificate, caKey crypto.Signer, dnsNames []string, ips []net.IP) (*x509.Certificate, *rsa.PrivateKey, error) {
	config := &certutil.Config{
		CommonName:   "TKE-SERVER",
		Organization: []string{"Tencent"},
		AltNames: certutil.AltNames{
			IPs:      ips,
			DNSNames: append(dnsNames, AlternateDNS()...),
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to sign certificate")
	}

	return cert, key, nil
}

func generateAdminCertKey(caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, *rsa.PrivateKey, error) {
	config := &certutil.Config{
		CommonName:   "admin",
		Organization: []string{"system:masters"},
		AltNames:     certutil.AltNames{},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to sign certificate")
	}

	return cert, key, nil
}

// AlternateDNS return TKE alternateDNS
func AlternateDNS() []string {
	result := []string{
		"localhost",
	}
	for _, one := range components {
		result = append(result, one)            // service in same namespace
		result = append(result, one+".tke.svc") // for apiservice
		result = append(result, one[4:])        // strip tke- for same namespace viste
	}

	return result
}
