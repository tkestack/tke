import React from 'react';
import { Upload } from '@tea/component';
import * as yaml from 'js-yaml';

export function KubeconfigFileParse() {
  function beforeUpload(file: File) {
    fileParse(file);

    return false;
  }

  function fileParse(file: File) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = ({ target }) => {
        if (target.readyState === 2) {
          const rsp = yaml.load(target.result as string);

          console.log(rsp);
          resolve(rsp);
        } else {
          reject(new Error());
        }
      };

      reader.readAsText(file);
    });
  }

  return (
    <Upload beforeUpload={beforeUpload}>
      <a style={{ fontSize: '12px' }}>点击上传kubeconfig文件</a>
    </Upload>
  );
}
