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
import { isEmpty } from './isEmpty';

/** 判断数组是否存在某个值 */
export const includes = (arr: Array<any> | string, value: any): boolean => {
  if (isEmpty(arr)) {
    /**如果数组为空 */
    return false;
  } else if (arr instanceof Array) {
    //判断数组是否包含元素
    return arr.findIndex(item => item === value) > -1;
  } else {
    //判断字符串是否包含子串
    return arr.indexOf(value) > -1;
  }
};
