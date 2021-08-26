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

/**
 * @returns 删除数据后的新数组，改变原数组
 */
export const remove = (arr: Array<any>, func: (value: any, index?: number, array?: any[]) => any) => {
  let targets = arr.filter(func);

  targets.forEach(elem => {
    arr.forEach((value, index) => {
      if (JSON.stringify(value) === JSON.stringify(elem)) {
        arr.splice(index, 1);
      }
    });
  });

  return arr;
};
