const tips = seajs.require('tips');

export function downloadCrt(crtText, filename = 'cluster-ca.crt') {
  let crtFile = '\ufeff';
  let userAgent = navigator.userAgent;

  if (navigator.msSaveBlob) {
    let blob = new Blob([crtText], { type: 'text/plain;charset=utf-8;' });
    navigator.msSaveBlob(blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    tips.error('该浏览器暂不支持导出功能');
  } else {
    let blob = new Blob([crtText], { type: 'text/plain;charset=utf-8;' });
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
