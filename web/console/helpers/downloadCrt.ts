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

import { Base64 } from 'js-base64';

// export function downloadCrt(crtText, filename = 'cluster-ca.crt') {
//   let crtFile = '\ufeff';
//   let userAgent = navigator.userAgent;

//   if (navigator.msSaveBlob) {
//     let blob = new Blob([crtText], { type: 'text/plain;charset=utf-8;' });
//     navigator.msSaveBlob(blob, filename);
//   } else if (userAgent.indexOf('MSIE 9.0') > 0) {
//     tips.error('该浏览器暂不支持导出功能');
//   } else {
//     let blob = new Blob([crtText], { type: 'text/plain;charset=utf-8;' });
//     let link = document.createElement('a') as any;

//     if (link.download !== undefined) {
//       let url = URL.createObjectURL(blob);
//       link.setAttribute('href', url);
//       link.setAttribute('download', filename);
//       link.style.visibility = 'hidden';
//       document.body.appendChild(link);
//       link.click();
//       document.body.removeChild(link);
//     }
//   }
// }

const tips = seajs.require('tips');

export function downloadText(crtText, filename, contentType = 'text/plain;charset=utf-8;') {
  const crtFile = '\ufeff';
  const userAgent = navigator.userAgent;

  if (navigator.msSaveBlob) {
    const blob = new Blob([crtText], { type: contentType });
    navigator.msSaveBlob(blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    tips.error('该浏览器暂不支持导出功能');
  } else {
    const blob = new Blob([crtText], { type: contentType });
    const link = document.createElement('a') as any;

    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', filename);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }
}

export function downloadCrt(crtText, filename = 'cluster-ca.crt') {
  downloadText(crtText, filename, 'text/plain;charset=utf-8;');
}

export function downloadKubeconfig(crtText, filename = 'kubeconfig') {
  downloadText(crtText, filename, 'applicatoin/octet-stream;charset=utf-8;');
}

export function getKubectlConfig({ caCert, token, host, clusterId, clientKey, clientCert }) {
  const user = Object.entries({
    token: token,
    'client-key-data': clientKey,
    'client-certificate-data': clientCert
  })
    .filter(([key, value]) => value)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n    ');

  return `
apiVersion: v1
clusters:
- cluster:
    ${caCert ? `certificate-authority-data: ${caCert}` : `insecure-skip-tls-verify: true`}
    server: ${host}
  name: ${clusterId}
contexts:
- context:
    cluster: ${clusterId}
    user: ${clusterId}-admin
  name: ${clusterId}-context-default
current-context: ${clusterId}-context-default
kind: Config
preferences: {}
users:
- name: ${clusterId}-admin
  user:
    ${user}
    `.trim();
}
