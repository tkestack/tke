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

export const isEmpty = (value: any): boolean => {
  if (Array.isArray(value)) {
    //value为数组
    return !value.length;
  } else if (typeof value === 'object') {
    //value为对象
    if (value === null) {
      //value为null
      return true;
    } else {
      //value是否没有key
      return !Object.keys(value).length;
    }
  } else if (typeof value === 'undefined') {
    //value为undefinded
    return true;
  } else if (Number.isFinite(value)) {
    //value为数值
    return false;
  } else {
    //value为默认值
    return !value;
  }
};
