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
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"net"
	"os"
	"path/filepath"

	"tkestack.io/tke/pkg/util/files"

	"tkestack.io/tke/cmd/tke-installer/app/installer/certs"

	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/segmentio/ksuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/kubeconfig"
)

type TkeCert struct {
	tmpDir string
}

func (c *TkeCert) InitTmpDir(namespace string) {
	pattern := fmt.Sprintf("tkestack.%s", namespace)
	c.tmpDir, _ = ioutil.TempDir("", pattern)
	_ = os.MkdirAll(filepath.Join(c.tmpDir, "data"), 0755)
}

func (c *TkeCert) ClearTmpDir() {
	os.RemoveAll(c.tmpDir)
}

func (c *TkeCert) CreateCertMap(ctx context.Context, client kubernetes.Interface, dnsNames []string, ips []net.IP,
	namespace string) error {
	err := certs.Generate(dnsNames, ips, c.tmpDir)
	if err != nil {
		return err
	}

	caCert, err := files.ReadFileWithDir(c.tmpDir, constants.CACrtFileBaseName)
	if err != nil {
		return err
	}
	caKey, err := files.ReadFileWithDir(c.tmpDir, constants.CAKeyFileBaseName)
	if err != nil {
		return err
	}
	serverCert, err := files.ReadFileWithDir(c.tmpDir, constants.ServerCrtFileBaseName)
	if err != nil {
		return err
	}
	serverKey, err := files.ReadFileWithDir(c.tmpDir, constants.ServerKeyFileBaseName)
	if err != nil {
		return err
	}
	adminCert, err := files.ReadFileWithDir(c.tmpDir, constants.AdminCrtFileBaseName)
	if err != nil {
		return err
	}
	adminKey, err := files.ReadFileWithDir(c.tmpDir, constants.AdminKeyFileBaseName)
	if err != nil {
		return err
	}
	webhookCert, err := files.ReadFileWithDir(c.tmpDir, constants.WebhookCrtFileBaseName)
	if err != nil {
		return err
	}
	webhookKey, err := files.ReadFileWithDir(c.tmpDir, constants.WebhookKeyFileBaseName)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: namespace,
		},
		Data: map[string]string{
			"etcd-ca.crt": string(caCert),
			"etcd.crt":    string(adminCert),
			"etcd.key":    string(adminKey),
			"ca.crt":      string(caCert),
			"ca.key":      string(caKey),
			"server.crt":  string(serverCert),
			"server.key":  string(serverKey),
			"admin.crt":   string(adminCert),
			"admin.key":   string(adminKey),
			"webhook.crt": string(webhookCert),
			"webhook.key": string(webhookKey),
		},
	}

	token := ksuid.New().String()
	cm.Data["token.csv"] = fmt.Sprintf("%s,admin,1,administrator", token)

	return apiclient.CreateOrUpdateConfigMap(ctx, client, cm)
}

func (c *TkeCert) WriteKubeConfig(host string, port int, namespace string) error {
	fmt.Println("write kube config")
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	fmt.Println("kubeconfig addr:", addr)

	caCert, err := files.ReadFileWithDir(c.tmpDir, constants.CACrtFileBaseName)
	if err != nil {
		return err
	}

	adminCert, err := files.ReadFileWithDir(c.tmpDir, constants.AdminCrtFileBaseName)
	if err != nil {
		return err
	}
	adminKey, err := files.ReadFileWithDir(c.tmpDir, constants.AdminKeyFileBaseName)
	if err != nil {
		return err
	}

	cfg := kubeconfig.CreateWithCerts(addr, namespace, "admin", caCert, adminKey, adminCert)
	data, err := runtime.Encode(clientcmdlatest.Codec, cfg)
	if err != nil {
		return err
	}

	klog.Info("Kubeconfig file path: ", c.GetKubeConfigFile())
	return files.WriteFileWithDir(c.tmpDir, constants.KubeconfigFileBaseName, data, 0644)
}

func (c *TkeCert) GetKubeConfig() (*restclient.Config, error) {
	kubeConfig, err := files.ReadFileWithDir(c.tmpDir, constants.KubeconfigFileBaseName)
	if err != nil {
		return nil, err
	}

	return clientcmd.RESTConfigFromKubeConfig(kubeConfig)
}

func (c *TkeCert) GetKubeConfigFile() string {
	return c.tmpDir + "/" + constants.KubeconfigFileBaseName
}
