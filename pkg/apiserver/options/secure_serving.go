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

package options

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"path"
	"strings"
	"time"

	"tkestack.io/tke/pkg/util/log"

	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	cliflag "k8s.io/component-base/cli/flag"
)

const (
	flagSecureServingBindAddress                  = "bind-address"
	flagSecureServingBindPort                     = "secure-port"
	flagSecureServingCertDir                      = "cert-dir"
	flagSecureServingTLSCertFile                  = "tls-cert-file"
	flagSecureServingTLSKeyFile                   = "tls-private-key-file"
	flagSecureServingTLSCipherSuites              = "tls-cipher-suites"
	flagSecureServingTLSMinVersion                = "tls-min-version"
	flagSecureServingTLSSNICertKeys               = "tls-sni-cert-key"
	flagSecureServingHTTP2MaxStreamsPerConnection = "http2-max-streams-per-connection"
)

const (
	configSecureServingBindAddress                  = "secure_serving.bind_address"
	configSecureServingBindPort                     = "secure_serving.port"
	configSecureServingCertDir                      = "secure_serving.cert_dir"
	configSecureServingTLSCertFile                  = "secure_serving.tls_cert_file"
	configSecureServingTLSKeyFile                   = "secure_serving.tls_private_key_file"
	configSecureServingTLSCipherSuites              = "secure_serving.tls_cipher_suites"
	configSecureServingTLSMinVersion                = "secure_serving.tls_min_version"
	configSecureServingTLSSNICertKeys               = "secure_serving.tls_sni_cert_key"
	configSecureServingHTTP2MaxStreamsPerConnection = "secure_serving.http2_max_streams_per_connection"
)

// SecureServingOptions contains the options that serve HTTPS.
type SecureServingOptions struct {
	*apiserveroptions.SecureServingOptionsWithLoopback
}

// NewSecureServingOptions gives default values for the HTTPS server which are
// not the options wanted by "normal" servers running on the platform.
func NewSecureServingOptions(serverName string, defaultPort int) *SecureServingOptions {
	o := apiserveroptions.SecureServingOptions{
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    defaultPort,
		Required:    true,
		ServerCert: apiserveroptions.GeneratableKeyCert{
			PairName:      serverName,
			CertDirectory: "_output/certificates",
		},
	}
	return &SecureServingOptions{
		o.WithLoopback(),
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	o.SecureServingOptionsWithLoopback.AddFlags(fs)

	_ = viper.BindPFlag(configSecureServingBindAddress, fs.Lookup(flagSecureServingBindAddress))
	_ = viper.BindPFlag(configSecureServingBindPort, fs.Lookup(flagSecureServingBindPort))
	_ = viper.BindPFlag(configSecureServingCertDir, fs.Lookup(flagSecureServingCertDir))
	_ = viper.BindPFlag(configSecureServingTLSCertFile, fs.Lookup(flagSecureServingTLSCertFile))
	_ = viper.BindPFlag(configSecureServingTLSKeyFile, fs.Lookup(flagSecureServingTLSKeyFile))
	_ = viper.BindPFlag(configSecureServingTLSCipherSuites, fs.Lookup(flagSecureServingTLSCipherSuites))
	_ = viper.BindPFlag(configSecureServingTLSMinVersion, fs.Lookup(flagSecureServingTLSMinVersion))
	_ = viper.BindPFlag(configSecureServingTLSSNICertKeys, fs.Lookup(flagSecureServingTLSSNICertKeys))
	_ = viper.BindPFlag(configSecureServingHTTP2MaxStreamsPerConnection, fs.Lookup(flagSecureServingHTTP2MaxStreamsPerConnection))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *SecureServingOptions) ApplyFlags() []error {
	var errs []error

	o.BindAddress = net.ParseIP(viper.GetString(configSecureServingBindAddress))
	o.BindPort = viper.GetInt(configSecureServingBindPort)
	o.ServerCert.CertDirectory = viper.GetString(configSecureServingCertDir)
	o.ServerCert.CertKey.CertFile = viper.GetString(configSecureServingTLSCertFile)
	o.ServerCert.CertKey.KeyFile = viper.GetString(configSecureServingTLSKeyFile)
	o.CipherSuites = viper.GetStringSlice(configSecureServingTLSCipherSuites)
	o.MinTLSVersion = viper.GetString(configSecureServingTLSMinVersion)

	nck := cliflag.NewNamedCertKeyArray(&o.SNICertKeys)
	sniCertKeysString := viper.GetString(configSecureServingTLSSNICertKeys)
	sniCertKeysStringRaw := strings.TrimPrefix(strings.TrimSuffix(sniCertKeysString, "]"), "[")
	sniCertKeysStringArray := strings.Split(sniCertKeysStringRaw, ";")
	if len(sniCertKeysStringArray) > 0 {
		for _, sniCertKey := range sniCertKeysStringArray {
			if sniCertKey != "" {
				if err := nck.Set(sniCertKey); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	o.HTTP2MaxStreamsPerConnection = viper.GetInt(configSecureServingHTTP2MaxStreamsPerConnection)

	if validateErrs := o.Validate(); len(validateErrs) > 0 {
		errs = append(errs, validateErrs...)
	}

	return errs
}

// rewrite apiserveroptions.SecureServingOptionsWithLoopback ApplyTo
// ApplyTo fills up serving information in the server configuration.
func (o *SecureServingOptions) ApplyTo(secureServingInfo **server.SecureServingInfo, loopbackClientConfig **rest.Config) error {
	if o == nil || o.SecureServingOptions == nil || secureServingInfo == nil {
		return nil
	}

	if err := o.SecureServingOptions.ApplyTo(secureServingInfo); err != nil {
		return err
	}

	if *secureServingInfo == nil || loopbackClientConfig == nil {
		return nil
	}

	// create self-signed cert+key with the fake server.LoopbackClientServerNameOverride and
	// let the server return it when the loopback client connects.
	certPem, keyPem, err := GenerateSelfSignedCertKey(server.LoopbackClientServerNameOverride, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
	}
	log.Infof("secureLoopbackClient cert: %s", string(certPem))
	log.Infof("secureLoopbackClient key: %s", string(keyPem))
	certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", certPem, keyPem, server.LoopbackClientServerNameOverride)
	if err != nil {
		return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
	}

	secureLoopbackClientConfig, err := (*secureServingInfo).NewLoopbackClientConfig(uuid.New().String(), certPem)
	switch {
	// if we failed and there's no fallback loopback client config, we need to fail
	case err != nil && *loopbackClientConfig == nil:
		return err

	// if we failed, but we already have a fallback loopback client config (usually insecure), allow it
	case err != nil && *loopbackClientConfig != nil:

	default:
		*loopbackClientConfig = secureLoopbackClientConfig
		// Write to the front of SNICerts so that this overrides any other certs with the same name
		(*secureServingInfo).SNICerts = append([]dynamiccertificates.SNICertKeyContentProvider{certProvider}, (*secureServingInfo).SNICerts...)
	}

	return nil
}

// copy from k8s.io/client-go/util/cert
// GenerateSelfSignedCertKey creates a self-signed certificate and key for the given host.
// Host may be an IP or a DNS name
// You may also specify additional subject alt names (either ip or dns names) for the certificate.
func GenerateSelfSignedCertKey(host string, alternateIPs []net.IP, alternateDNS []string) ([]byte, []byte, error) {
	return GenerateSelfSignedCertKeyWithFixtures(host, alternateIPs, alternateDNS, "")
}

// copy from k8s.io/client-go/util/cert
// GenerateSelfSignedCertKeyWithFixtures creates a self-signed certificate and key for the given host.
// Host may be an IP or a DNS name. You may also specify additional subject alt names (either ip or dns names)
// for the certificate.
//
// If fixtureDirectory is non-empty, it is a directory path which can contain pre-generated certs. The format is:
// <host>_<ip>-<ip>_<alternateDNS>-<alternateDNS>.crt
// <host>_<ip>-<ip>_<alternateDNS>-<alternateDNS>.key
// Certs/keys not existing in that directory are created.
func GenerateSelfSignedCertKeyWithFixtures(host string, alternateIPs []net.IP, alternateDNS []string, fixtureDirectory string) ([]byte, []byte, error) {
	validFrom := time.Now().Add(-time.Hour) // valid an hour earlier to avoid flakes due to clock skew
	maxAge := time.Hour * 24 * 365 * 10     // 10 year self-signed certs

	baseName := fmt.Sprintf("%s_%s_%s", host, strings.Join(ipsToStrings(alternateIPs), "-"), strings.Join(alternateDNS, "-"))
	certFixturePath := path.Join(fixtureDirectory, baseName+".crt")
	keyFixturePath := path.Join(fixtureDirectory, baseName+".key")
	if len(fixtureDirectory) > 0 {
		cert, err := ioutil.ReadFile(certFixturePath)
		if err == nil {
			key, err := ioutil.ReadFile(keyFixturePath)
			if err == nil {
				return cert, key, nil
			}
			return nil, nil, fmt.Errorf("cert %s can be read, but key %s cannot: %v", certFixturePath, keyFixturePath, err)
		}
		maxAge = 100 * time.Hour * 24 * 365 // 100 years fixtures
	}

	caKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s-ca@%d", host, time.Now().Unix()),
		},
		NotBefore: validFrom,
		NotAfter:  validFrom.Add(maxAge),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	caCertificate, err := x509.ParseCertificate(caDERBytes)
	if err != nil {
		return nil, nil, err
	}

	priv, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, time.Now().Unix()),
		},
		NotBefore: validFrom,
		NotAfter:  validFrom.Add(maxAge),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, caCertificate, &priv.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate cert, followed by ca
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: certutil.CertificateBlockType, Bytes: derBytes}); err != nil {
		return nil, nil, err
	}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: certutil.CertificateBlockType, Bytes: caDERBytes}); err != nil {
		return nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: keyutil.RSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	if len(fixtureDirectory) > 0 {
		if err := ioutil.WriteFile(certFixturePath, certBuffer.Bytes(), 0644); err != nil {
			return nil, nil, fmt.Errorf("failed to write cert fixture to %s: %v", certFixturePath, err)
		}
		if err := ioutil.WriteFile(keyFixturePath, keyBuffer.Bytes(), 0644); err != nil {
			return nil, nil, fmt.Errorf("failed to write key fixture to %s: %v", certFixturePath, err)
		}
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

func ipsToStrings(ips []net.IP) []string {
	ss := make([]string, 0, len(ips))
	for _, ip := range ips {
		ss = append(ss, ip.String())
	}
	return ss
}
