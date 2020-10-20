const tips = seajs.require('tips');
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

export function downloadText(crtText, filename, contentType = 'text/plain;charset=utf-8;') {
  let crtFile = '\ufeff';
  let userAgent = navigator.userAgent;

  if (navigator.msSaveBlob) {
    let blob = new Blob([crtText], { type: contentType });
    navigator.msSaveBlob(blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    tips.error('该浏览器暂不支持导出功能');
  } else {
    let blob = new Blob([crtText], { type: contentType });
    let link = document.createElement('a') as any;

    if (link.download !== undefined) {
      let url = URL.createObjectURL(blob);
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

export function getKubectlConfig({ caCert, token, host, clusterId }) {
  let config = `apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: ${caCert}\n    server: ${host}\n  name: ${clusterId}\ncontexts:\n- context:\n    cluster: ${clusterId}\n    user: ${clusterId}-admin\n  name: ${clusterId}-context-default\ncurrent-context: ${clusterId}-context-default\nkind: Config\npreferences: {}\nusers:\n- name: ${clusterId}-admin\n  user:\n    token: ${token}\n`;
  return config;
}
