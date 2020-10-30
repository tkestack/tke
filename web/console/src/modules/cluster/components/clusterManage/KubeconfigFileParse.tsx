import React from 'react';
import { Upload } from '@tea/component';
import * as yaml from 'js-yaml';

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
      token: string;
    };
  }>;
}

export interface KubeconfigFileParseProps {
  onSuccess: (targetConfig: { apiServer: string; certFile: string; token: string }) => any;
  onFaild?: () => any;
}

export function KubeconfigFileParse({ onSuccess, onFaild }: KubeconfigFileParseProps) {
  function beforeUpload(file: File) {
    fileParse(file)
      .then((rsp: KubeConfig) => {
        onSuccess(yamlToConfig(rsp));
      })
      .catch(error => {
        onFaild && onFaild();
      });

    return false;
  }

  function fileParse(file: File) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = ({ target }) => {
        if (target.readyState === 2) {
          let rsp;
          try {
            rsp = yaml.load(target.result as string);
          } catch (error) {
            return reject(error);
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
        cluster: { server, 'certificate-authority-data': certFile }
      }
    ],
    users: [
      {
        user: { token }
      }
    ]
  }: KubeConfig) {
    return {
      apiServer: server,
      certFile,
      token
    };
  }

  return (
    <Upload beforeUpload={beforeUpload}>
      <a style={{ fontSize: '12px' }}>点击上传kubeconfig文件</a>
    </Upload>
  );
}
