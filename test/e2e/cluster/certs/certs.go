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
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	certutil "k8s.io/client-go/util/cert"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/kubeconfig"
	"tkestack.io/tke/pkg/util/pkiutil"
)

var (
	tmpDir     = ""
	components = []string{
		"tke-platform-api",
		"etcd",
		"etcd-client",
	}
)

func InitTmpDir() {
	tmpDir, _ = ioutil.TempDir("", "tkestack")
	_ = os.MkdirAll(prefixWithTmpDir("data"), 0755)

}

func ClearTmpDir() {
	os.RemoveAll(tmpDir)
}

func CreateCertMap(ctx context.Context, client kubernetes.Interface, ips []net.IP, namespace string) error {
	err := generateCerts(ips)
	if err != nil {
		return err
	}

	caCrt, err := ioutil.ReadFile(prefixWithTmpDir(constants.CACrtFile))
	if err != nil {
		return err
	}
	caKey, err := ioutil.ReadFile(prefixWithTmpDir(constants.CAKeyFile))
	if err != nil {
		return err
	}
	serverCrt, err := ioutil.ReadFile(prefixWithTmpDir(constants.ServerCrtFile))
	if err != nil {
		return err
	}
	serverKey, err := ioutil.ReadFile(prefixWithTmpDir(constants.ServerKeyFile))
	if err != nil {
		return err
	}
	adminCrt, err := ioutil.ReadFile(prefixWithTmpDir(constants.AdminCrtFile))
	if err != nil {
		return err
	}
	adminKey, err := ioutil.ReadFile(prefixWithTmpDir(constants.AdminKeyFile))
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: namespace,
		},
		Data: map[string]string{
			"etcd-ca.crt": string(caCrt),
			"etcd.crt":    string(adminCrt),
			"etcd.key":    string(adminKey),
			"ca.crt":      string(caCrt),
			"ca.key":      string(caKey),
			"server.crt":  string(serverCrt),
			"server.key":  string(serverKey),
			"admin.crt":   string(adminCrt),
			"admin.key":   string(adminKey),
		},
	}

	cm.Data["password.csv"] = fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())
	token := ksuid.New().String()
	cm.Data["token.csv"] = fmt.Sprintf("%s,admin,1,administrator", token)

	return apiclient.CreateOrUpdateConfigMap(ctx, client, cm)
}

func WriteKubeConfig(host string, port int, namespace string) error {
	fmt.Println("write kube config")
	addr := fmt.Sprintf("%s:%d", host, port)

	fmt.Println("kubeconfig addr:", addr)

	caCrt, err := ioutil.ReadFile(prefixWithTmpDir(constants.CACrtFile))
	if err != nil {
		return err
	}

	adminCrt, err := ioutil.ReadFile(prefixWithTmpDir(constants.AdminCrtFile))
	if err != nil {
		return err
	}
	adminKey, err := ioutil.ReadFile(prefixWithTmpDir(constants.AdminKeyFile))
	if err != nil {
		return err
	}

	cfg := kubeconfig.CreateWithCerts(addr, namespace, "admin", caCrt, adminKey, adminCrt)
	data, err := runtime.Encode(clientcmdlatest.Codec, cfg)
	if err != nil {
		return err
	}

	_ = os.MkdirAll("/root/.kube", 0755)
	ioutil.WriteFile("/root/.kube/config", data, 0644)
	return ioutil.WriteFile(prefixWithTmpDir(constants.KubeconfigFile), data, 0644)
}

func GetKubeConfig() (*restclient.Config, error) {
	kubeconfig, err := ioutil.ReadFile(prefixWithTmpDir(constants.KubeconfigFile))
	if err != nil {
		return nil, err
	}

	return clientcmd.RESTConfigFromKubeConfig(kubeconfig)
}

func prefixWithTmpDir(path string) string {
	return tmpDir + "/" + path
}

func generateCerts(ips []net.IP) error {
	caCert, caKey, err := generateRootCA()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.CACrtFile), pkiutil.EncodeCertPEM(caCert), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.CAKeyFile), pkiutil.EncodePrivateKeyPEM(caKey), 0644)
	if err != nil {
		return err
	}
	serverCert, serverKey, err := generateServerCertKey(caCert, caKey, nil, ips)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.ServerCrtFile), pkiutil.EncodeCertPEM(serverCert), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.ServerKeyFile), pkiutil.EncodePrivateKeyPEM(serverKey), 0644)
	if err != nil {
		return err
	}
	adminCert, adminKey, err := generateAdminCertKey(caCert, caKey)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.AdminCrtFile), pkiutil.EncodeCertPEM(adminCert), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(prefixWithTmpDir(constants.AdminKeyFile), pkiutil.EncodePrivateKeyPEM(adminKey), 0644)
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
