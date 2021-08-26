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
 * 深度拷贝值
 */
export function deepClone<T>(source: T): T {
  if (source instanceof Array) {
    return deepCloneArray(source as any) as any;
  }

  if (source && typeof source === 'object') {
    return deepCloneObject(source);
  }

  return source;
}

function deepCloneObject(source: any): any {
  const target: any = {};

  Object.keys(source).forEach(key => {
    target[key] = deepClone(source[key]);
  });

  return target;
}

function deepCloneArray(source: any[]) {
  return source.map(x => deepClone(x));
}
