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
/**
 * 计算字符串所占的内存字节数，默认使用 utf-8的形式进行编码
 * @param userScript:string
 * @param charset:编码方式
 * @returns length  length = bytes
 */
export const stringSizeOf = (userScript: string, encode: 'utf-8' | 'utf-16' = 'utf-8') => {
  let totalSize = 0,
    charset = encode ? encode.toLowerCase() : '';

  if (charset === 'utf-16') {
    userScript.split('').forEach((item, index) => {
      let charCode = userScript.charCodeAt(index);
      if (charCode <= 0xffff) {
        totalSize += 2;
      } else {
        totalSize += 4;
      }
    });
  } else {
    userScript.split('').forEach((item, index) => {
      let charCode = userScript.charCodeAt(index);
      if (charCode <= 0x007f) {
        totalSize += 1;
      } else if (charCode <= 0x07ff) {
        totalSize += 2;
      } else if (charCode <= 0xffff) {
        // utf-8 在范围 0xD800 - 0xDFFF不存在任何字符
        totalSize += 3;
      } else {
        totalSize += 4;
      }
    });
  }

  return totalSize;
};
