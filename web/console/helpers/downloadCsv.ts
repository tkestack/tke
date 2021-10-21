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
// const tips = seajs.require('tips');

export function downloadCsv(rows, header, filename) {
  let processRow = function(row) {
    let finalVal = '';
    for (let j = 0; j < row.length; j++) {
      let innerValue = typeof row[j] === 'undefined' || row[j] === null ? '' : row[j].toString();
      if (row[j] instanceof Date) {
        innerValue = row[j].toLocaleString();
      }
      let result = innerValue.replace(/"/g, '""');
      if (result.search(/("|,|\n)/g) >= 0) result = '"' + result + '"';
      if (j > 0) finalVal += ',';
      finalVal += result;
    }
    return finalVal + '\n';
  };

  let csvFile = '\ufeff';
  let csv_array = [header.join(',') + '\n'];
  let userAgent = navigator.userAgent;
  for (let i = 0; i < rows.length; i++) {
    csv_array.push(processRow(rows[i]));
  }
  csvFile += csv_array.join('');

  if (navigator.msSaveBlob) {
    let blob = new Blob([csvFile], { type: 'text/csv;charset=utf-8;' });
    navigator.msSaveBlob(blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    // tips.error('该浏览器暂不支持导出功能');
  } else {
    let blob = new Blob([csvFile], { type: 'text/csv;charset=utf-8;' });
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
