import { t, Trans } from '@tencent/tea-app/lib/i18n';
const tips = seajs.require('tips');

export function downloadText(crtText, filename, contentType = 'text/plain;charset=utf-8;') {
  const crtFile = '\ufeff';
  const userAgent = navigator.userAgent;

  if (navigator.msSaveBlob) {
    const blob = new Blob([crtText], { type: contentType });
    navigator.msSaveBlob(blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    tips.error(t('该浏览器暂不支持导出功能'));
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
