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

package authzwebhook

import (
	"bytes"
	"io/ioutil"

	installerconstants "tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	utilfile "tkestack.io/tke/pkg/util/file"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"

	"github.com/pkg/errors"
)

const (
	authzWebhookConfig = `
apiVersion: v1
kind: Config
clusters:
  - name: tke
    cluster:
      server: {{.AuthzEndpoint}}
      insecure-skip-tls-verify: true
users:
  - name: admin-cert
    user:
      client-certificate: {{.WebhookCertFile}}
      client-key: {{.WebhookKeyFile}}
current-context: tke
contexts:
- context:
    cluster: tke
    user: admin-cert
  name: tke
`
)

type Option struct {
	AuthzWebhookEndpoint string
	IsGlobalCluster      bool
	IsClusterUpscaling   bool
}

// WebhookCertAndKeyExist checks whether the certificate and private key exist,
// for compatibility with old version clusters' webhook certificates and private keys which version are before 1.5,
// and we will completely replace webhook certificates and private keys' file name in 1.6 or future release.
func WebhookCertAndKeyExist(basePath string) bool {
	return utilfile.Exists(basePath+constants.WebhookCertName) &&
		utilfile.Exists(basePath+constants.WebhookKeyName)
}

func Install(s ssh.Interface, option *Option) error {
	var webhookCertFile = constants.WebhookCertFile
	var webhookKeyFile = constants.WebhookKeyFile
	var webhookCertName = constants.WebhookCertName
	var webhookKeyName = constants.WebhookKeyName

	basePath := constants.AppCertDir
	if option.IsGlobalCluster && !option.IsClusterUpscaling {
		basePath = installerconstants.DataDir
	}
	// For compatibility with old version clusters' webhook certificates and private keys.
	if !WebhookCertAndKeyExist(basePath) {
		webhookCertFile = constants.AdminCertFile
		webhookKeyFile = constants.AdminKeyFile
		webhookCertName = constants.AdminCertName
		webhookKeyName = constants.AdminKeyName
	}

	authzWebhookConfig, err := template.ParseString(authzWebhookConfig, map[string]interface{}{
		"AuthzEndpoint":   option.AuthzWebhookEndpoint,
		"WebhookCertFile": webhookCertFile,
		"WebhookKeyFile":  webhookKeyFile,
	})
	if err != nil {
		return errors.Wrap(err, "parse authzWebhookConfig error")
	}

	err = s.WriteFile(bytes.NewReader(authzWebhookConfig), constants.KubernetesAuthzWebhookConfigFile)
	if err != nil {
		return err
	}
	webhookCertData, err := ioutil.ReadFile(basePath + webhookCertName)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(webhookCertData), webhookCertFile)
	if err != nil {
		return err
	}
	webhookKeyData, err := ioutil.ReadFile(basePath + webhookKeyName)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(webhookKeyData), webhookKeyFile)
	if err != nil {
		return err
	}

	return nil
}
