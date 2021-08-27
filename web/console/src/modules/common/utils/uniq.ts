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
 * 数组去重函数
 * @param Array<any> arr
 * @param string by 数组去重字段，非必需
 * @returns 返回去重后的数组，原数组保持不变
 */
export const uniq = (arr: Array<any>, by?: string) => {
  return arr.filter((elem, pos, arr) => {
    if (by) {
      return pos === arr.findIndex(item => JSON.stringify(item[by]) === JSON.stringify(elem[by]));
    } else {
      return pos === arr.findIndex(item => JSON.stringify(item) === JSON.stringify(elem));
    }
  });
};
