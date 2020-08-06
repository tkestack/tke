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

package certs

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"

	"tkestack.io/tke/pkg/util/apiclient"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	certutil "k8s.io/client-go/util/cert"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/files"
	"tkestack.io/tke/pkg/util/pkiutil"
)

const (
	CACrtFileBaseName     = "ca.crt"
	CAKeyFileBaseName     = "ca.key"
	ServerCrtFileBaseName = "server.crt"
	ServerKeyFileBaseName = "server.key"
)

func GenerateInDir(dir string, namespace string, dnsNames []string, ips []net.IP) error {
	caCert, caKey, err := NewRootCA()
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, CACrtFileBaseName, pkiutil.EncodeCertPEM(caCert), 0644)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, CAKeyFileBaseName, pkiutil.EncodePrivateKeyPEM(caKey), 0644)
	if err != nil {
		return err
	}

	serverCert, serverKey, err := NewServerCertKey(caCert, caKey, namespace, dnsNames, ips)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, ServerCrtFileBaseName, pkiutil.EncodeCertPEM(serverCert), 0644)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(dir, ServerKeyFileBaseName, pkiutil.EncodePrivateKeyPEM(serverKey), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GenerateInK8s(clientset kubernetes.Interface, name string, namespace string, dnsNames []string, ips []net.IP) error {
	caCert, caKey, err := NewRootCA()
	if err != nil {
		return err
	}

	serverCert, serverKey, err := NewServerCertKey(caCert, caKey, namespace, dnsNames, ips)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Data: map[string][]byte{
			CACrtFileBaseName:     pkiutil.EncodeCertPEM(caCert),
			CAKeyFileBaseName:     pkiutil.EncodePrivateKeyPEM(caKey),
			ServerCrtFileBaseName: pkiutil.EncodeCertPEM(serverCert),
			ServerKeyFileBaseName: pkiutil.EncodePrivateKeyPEM(serverKey),
		},
	}

	return apiclient.CreateOrUpdateSecret(context.Background(), clientset, secret)
}

func NewRootCA() (*x509.Certificate, *rsa.PrivateKey, error) {
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

func NewServerCertKey(caCert *x509.Certificate, caKey crypto.Signer, namespace string, dnsNames []string, ips []net.IP) (*x509.Certificate, *rsa.PrivateKey, error) {
	config := &certutil.Config{
		CommonName:   "TKE-SERVER",
		Organization: []string{"Tencent"},
		AltNames: certutil.AltNames{
			IPs:      ips,
			DNSNames: append(dnsNames, alternateDNS(namespace)...),
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to sign certificate")
	}

	return cert, key, nil
}

func NewAdminCertKey(caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, *rsa.PrivateKey, error) {
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

// alternateDNS return TKE alternateDNS
func alternateDNS(namespace string) []string {
	result := []string{
		"localhost",
	}
	for _, one := range spec.Components {
		result = append(result, one)                                      // service in same namespace
		result = append(result, fmt.Sprintf("%s.%s.svc", one, namespace)) // for apiservice
	}

	return result
}
