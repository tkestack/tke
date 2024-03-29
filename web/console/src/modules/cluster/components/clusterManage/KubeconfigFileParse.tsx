/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

import React, { useState } from 'react';
import { Upload, Text } from '@tea/component';
import * as yaml from 'js-yaml';
import { v4 as uuidv4 } from 'uuid';

export interface KubeConfig {
  apiVersion: string;
  clusters: Array<{
    cluster: {
      'certificate-authority-data': string;
      server: string;
    };
    name: string;
  }>;
  'current-context': string;
  kind: string;
  preferences: {};
  users: Array<{
    name: string;
    user: {
      token?: string;
      'client-certificate-data'?: string;
      'client-key-data'?: string;
      username?: string;
      as?: string;
      'as-groups'?: string[];
      'as-user-extra'?: {
        [key: string]: string[];
      };
    };
  }>;
}

export interface KubeconfigFileParseProps {
  onSuccess: (targetConfig: {
    apiServer: string;
    certFile: string;
    token: string;
    clientCert: string;
    clientKey: string;
    username: string;
    as: string;
    asGroups: string;
    asUserExtra: Array<{
      id: string;
      key: string;
      value: string;
    }>;
  }) => any;
}

export function KubeconfigFileParse({ onSuccess }: KubeconfigFileParseProps) {
  const [errorMessage, setErrorMessage] = useState('');

  function beforeUpload(file: File) {
    setErrorMessage('');
    fileParse(file)
      .then((rsp: KubeConfig) => {
        onSuccess(yamlToConfig(rsp));
      })
      .catch(error => {
        onFiled(error);
      });

    return false;
  }

  function onFiled(error: Error) {
    console.log(error);
    setErrorMessage(error.message);
  }

  function fileParse(file: File) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = ({ target }: any) => {
        if (target.readyState === 2) {
          let rsp;
          try {
            rsp = yaml.load(target.result as string);
          } catch (error) {
            return reject(new Error('文件解析失败'));
          }

          if (rsp.kind !== 'Config') {
            return reject(new Error('这不是一个kubeconfig文件'));
          }

          return resolve(rsp);
        } else {
          return reject(new Error('文件解析失败'));
        }
      };

      reader.readAsText(file);
    });
  }

  function yamlToConfig({
    clusters: [
      {
        cluster: { server = '', 'certificate-authority-data': certFile = '' }
      }
    ],
    users: [
      {
        user: {
          token = '',
          'client-certificate-data': clientCert = '',
          'client-key-data': clientKey = '',
          username = '',
          as = '',
          'as-groups': asGroups = [],
          'as-user-extra': asUserExtra = {}
        }
      }
    ]
  }: KubeConfig) {
    return {
      apiServer: server,
      certFile,
      token,
      clientCert,
      clientKey,
      username,
      as,
      asGroups: asGroups.join(','),
      asUserExtra: Object.entries(asUserExtra).reduce(
        (all, [key, values]) => [...all, ...values.map(v => ({ key, value: v, id: uuidv4() }))],
        []
      )
    };
  }

  return (
    <Upload beforeUpload={beforeUpload}>
      <a style={{ fontSize: '12px' }}>点击上传kubeconfig文件</a>
      <br />
      {errorMessage && (
        <Text style={{ fontSize: '12px' }} theme="danger">
          {errorMessage}
        </Text>
      )}
    </Upload>
  );
}
