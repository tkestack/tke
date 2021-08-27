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
/** 从枚举类型创建 key-value 一样的常量 */
export function createConstant<T>(enumType: T): T {
  let constant = {};
  for (let key in enumType) {
    if (enumType.hasOwnProperty(key) && isNaN(+key)) {
      constant[key as any] = key;
    }
  }
  return constant as T;
}
