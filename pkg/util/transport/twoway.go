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

package transport

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"tkestack.io/tke/pkg/util/log"
)

// NewTwoWayTLSTransport create a two-way SSL HTTP transport object by given
// certificate file and client certificate and private key file.
func NewTwoWayTLSTransport(caFile string, clientCertFile string, clientKeyFile string) (*http.Transport, error) {
	var cert *x509.CertPool
	if caFile == "" {
		cert = nil
	} else {
		cert = x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Error("Failed to read the CA certificate file", log.String("file", caFile), log.Err(err))
			return nil, err
		}
		cert.AppendCertsFromPEM(ca)
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Error("Failed to load X509 client certificate", log.String("certFile", clientCertFile), log.String("keyFile", clientKeyFile), log.Err(err))
		return nil, err
	}

	return twoWayTransport(cert, &clientCert), nil
}

func twoWayTransport(cert *x509.CertPool, clientCert *tls.Certificate) *http.Transport {
	tlsConfig := &tls.Config{}
	if cert != nil {
		tlsConfig.RootCAs = cert
	}
	if clientCert != nil {
		tlsConfig.Certificates = []tls.Certificate{*clientCert}
	}
	return Transport(tlsConfig)
}
