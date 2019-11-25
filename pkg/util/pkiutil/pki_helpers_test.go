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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pkiutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	certutil "k8s.io/client-go/util/cert"
	"os"
	"testing"
)

func TestNewCertificateAuthority(t *testing.T) {
	cert, key, err := NewCertificateAuthority(&certutil.Config{CommonName: "kubernetes"})

	if cert == nil {
		t.Error("failed NewCertificateAuthority, cert == nil")
	} else if !cert.IsCA {
		t.Error("cert is not a valida CA")
	}

	if key == nil {
		t.Error("failed NewCertificateAuthority, key == nil")
	}

	if err != nil {
		t.Errorf("failed NewCertificateAuthority with an error: %+v", err)
	}
}

func TestNewCertAndKey(t *testing.T) {
	var tests = []struct {
		name       string
		keyGenFunc func() (crypto.Signer, error)
		expected   bool
	}{
		{
			name: "RSA key too small",
			keyGenFunc: func() (crypto.Signer, error) {
				return rsa.GenerateKey(rand.Reader, 128)
			},
			expected: false,
		},
		{
			name: "RSA should succeed",
			keyGenFunc: func() (crypto.Signer, error) {
				return rsa.GenerateKey(rand.Reader, 2048)
			},
			expected: true,
		},
		{
			name: "ECDSA should succeed",
			keyGenFunc: func() (crypto.Signer, error) {
				return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			},
			expected: true,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			caKey, err := rt.keyGenFunc()
			if err != nil {
				t.Fatalf("Couldn't create Private Key")
			}
			caCert := &x509.Certificate{}
			config := &certutil.Config{
				CommonName:   "test",
				Organization: []string{"test"},
				Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}
			_, _, actual := NewCertAndKey(caCert, caKey, config)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed NewCertAndKey:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					actual == nil,
				)
			}
		})
	}
}

func TestHasServerAuth(t *testing.T) {
	caCert, caKey, _ := NewCertificateAuthority(&certutil.Config{CommonName: "kubernetes"})

	var tests = []struct {
		name     string
		config   certutil.Config
		expected bool
	}{
		{
			name: "has ServerAuth",
			config: certutil.Config{
				CommonName: "test",
				Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			},
			expected: true,
		},
		{
			name: "doesn't have ServerAuth",
			config: certutil.Config{
				CommonName: "test",
				Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			},
			expected: false,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			cert, _, err := NewCertAndKey(caCert, caKey, &rt.config)
			if err != nil {
				t.Fatalf("Couldn't create cert: %v", err)
			}
			actual := HasServerAuth(cert)
			if actual != rt.expected {
				t.Errorf(
					"failed HasServerAuth:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					actual,
				)
			}
		})
	}
}

func TestWriteCertAndKey(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Couldn't create rsa Private Key")
	}
	caCert := &x509.Certificate{}
	actual := WriteCertAndKey(tmpdir, "foo", caCert, caKey)
	if actual != nil {
		t.Errorf(
			"failed WriteCertAndKey with an error: %v",
			actual,
		)
	}
}

func TestWriteCert(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caCert := &x509.Certificate{}
	actual := WriteCert(tmpdir, "foo", caCert)
	if actual != nil {
		t.Errorf(
			"failed WriteCertAndKey with an error: %v",
			actual,
		)
	}
}

func TestWriteKey(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Couldn't create rsa Private Key")
	}
	actual := WriteKey(tmpdir, "foo", caKey)
	if actual != nil {
		t.Errorf(
			"failed WriteCertAndKey with an error: %v",
			actual,
		)
	}
}

func TestWritePublicKey(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Couldn't create rsa Private Key")
	}
	actual := WritePublicKey(tmpdir, "foo", &caKey.PublicKey)
	if actual != nil {
		t.Errorf(
			"failed WriteCertAndKey with an error: %v",
			actual,
		)
	}
}

func TestCertOrKeyExist(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Couldn't create rsa Private Key")
	}
	caCert := &x509.Certificate{}
	actual := WriteCertAndKey(tmpdir, "foo", caCert, caKey)
	if actual != nil {
		t.Errorf(
			"failed WriteCertAndKey with an error: %v",
			actual,
		)
	}

	var tests = []struct {
		desc     string
		path     string
		name     string
		expected bool
	}{
		{
			desc:     "empty path and name",
			path:     "",
			name:     "",
			expected: false,
		},
		{
			desc:     "valid path and name",
			path:     tmpdir,
			name:     "foo",
			expected: true,
		},
	}
	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			actual := CertOrKeyExist(rt.path, rt.name)
			if actual != rt.expected {
				t.Errorf(
					"failed CertOrKeyExist:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					actual,
				)
			}
		})
	}
}

func TestTryLoadCertAndKeyFromDisk(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caCert, caKey, err := NewCertificateAuthority(&certutil.Config{CommonName: "kubernetes"})
	if err != nil {
		t.Errorf(
			"failed to create cert and key with an error: %v",
			err,
		)
	}
	err = WriteCertAndKey(tmpdir, "foo", caCert, caKey)
	if err != nil {
		t.Errorf(
			"failed to write cert and key with an error: %v",
			err,
		)
	}

	var tests = []struct {
		desc     string
		path     string
		name     string
		expected bool
	}{
		{
			desc:     "empty path and name",
			path:     "",
			name:     "",
			expected: false,
		},
		{
			desc:     "valid path and name",
			path:     tmpdir,
			name:     "foo",
			expected: true,
		},
	}
	for _, rt := range tests {
		t.Run(rt.desc, func(t *testing.T) {
			_, _, actual := TryLoadCertAndKeyFromDisk(rt.path, rt.name)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed TryLoadCertAndKeyFromDisk:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					(actual == nil),
				)
			}
		})
	}
}

func TestTryLoadCertFromDisk(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Couldn't create tmpdir")
	}
	defer os.RemoveAll(tmpdir)

	caCert, _, err := NewCertificateAuthority(&certutil.Config{CommonName: "kubernetes"})
	if err != nil {
		t.Errorf(
			"failed to create cert and key with an error: %v",
			err,
		)
	}
	err = WriteCert(tmpdir, "foo", caCert)
	if err != nil {
		t.Errorf(
			"failed to write cert and key with an error: %v",
			err,
		)
	}

	var tests = []struct {
		desc     string
		path     string
		name     string
		expected bool
	}{
		{
			desc:     "empty path and name",
			path:     "",
			name:     "",
			expected: false,
		},
		{
			desc:     "valid path and name",
			path:     tmpdir,
			name:     "foo",
			expected: true,
		},
	}
	for _, rt := range tests {
		t.Run(rt.desc, func(t *testing.T) {
			_, actual := TryLoadCertFromDisk(rt.path, rt.name)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed TryLoadCertAndKeyFromDisk:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					(actual == nil),
				)
			}
		})
	}
}

func TestTryLoadKeyFromDisk(t *testing.T) {

	var tests = []struct {
		desc       string
		pathSuffix string
		name       string
		keyGenFunc func() (crypto.Signer, error)
		expected   bool
	}{
		{
			desc:       "empty path and name",
			pathSuffix: "somegarbage",
			name:       "",
			keyGenFunc: func() (crypto.Signer, error) {
				return rsa.GenerateKey(rand.Reader, 2048)
			},
			expected: false,
		},
		{
			desc:       "RSA valid path and name",
			pathSuffix: "",
			name:       "foo",
			keyGenFunc: func() (crypto.Signer, error) {
				return rsa.GenerateKey(rand.Reader, 2048)
			},
			expected: true,
		},
		{
			desc:       "ECDSA valid path and name",
			pathSuffix: "",
			name:       "foo",
			keyGenFunc: func() (crypto.Signer, error) {
				return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			},
			expected: true,
		},
	}
	for _, rt := range tests {
		t.Run(rt.desc, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatalf("Couldn't create tmpdir")
			}
			defer os.RemoveAll(tmpdir)

			caKey, err := rt.keyGenFunc()
			if err != nil {
				t.Errorf(
					"failed to create key with an error: %v",
					err,
				)
			}

			err = WriteKey(tmpdir, "foo", caKey)
			if err != nil {
				t.Errorf(
					"failed to write key with an error: %v",
					err,
				)
			}
			_, actual := TryLoadKeyFromDisk(tmpdir+rt.pathSuffix, rt.name)
			if (actual == nil) != rt.expected {
				t.Errorf(
					"failed TryLoadCertAndKeyFromDisk:\n\texpected: %t\n\t  actual: %t",
					rt.expected,
					(actual == nil),
				)
			}
		})
	}
}

func TestPathsForCertAndKey(t *testing.T) {
	crtPath, keyPath := PathsForCertAndKey("/foo", "bar")
	if crtPath != "/foo/bar.crt" {
		t.Errorf("unexpected certificate path: %s", crtPath)
	}
	if keyPath != "/foo/bar.key" {
		t.Errorf("unexpected key path: %s", keyPath)
	}
}

func TestPathForCert(t *testing.T) {
	crtPath := pathForCert("/foo", "bar")
	if crtPath != "/foo/bar.crt" {
		t.Errorf("unexpected certificate path: %s", crtPath)
	}
}

func TestPathForKey(t *testing.T) {
	keyPath := pathForKey("/foo", "bar")
	if keyPath != "/foo/bar.key" {
		t.Errorf("unexpected certificate path: %s", keyPath)
	}
}

func TestPathForPublicKey(t *testing.T) {
	pubPath := pathForPublicKey("/foo", "bar")
	if pubPath != "/foo/bar.pub" {
		t.Errorf("unexpected certificate path: %s", pubPath)
	}
}

func TestPathForCSR(t *testing.T) {
	csrPath := pathForCSR("/foo", "bar")
	if csrPath != "/foo/bar.csr" {
		t.Errorf("unexpected certificate path: %s", csrPath)
	}
}

func TestAppendSANsToAltNames(t *testing.T) {
	var tests = []struct {
		sans     []string
		expected int
	}{
		{[]string{}, 0},
		{[]string{"abc"}, 1},
		{[]string{"*.abc"}, 1},
		{[]string{"**.abc"}, 0},
		{[]string{"a.*.bc"}, 0},
		{[]string{"a.*.bc", "abc.def"}, 1},
		{[]string{"a*.bc", "abc.def"}, 1},
	}
	for _, rt := range tests {
		altNames := certutil.AltNames{}
		appendSANsToAltNames(&altNames, rt.sans, "foo")
		actual := len(altNames.DNSNames)
		if actual != rt.expected {
			t.Errorf(
				"failed AppendSANsToAltNames Numbers:\n\texpected: %d\n\t  actual: %d",
				rt.expected,
				actual,
			)
		}
	}

}
