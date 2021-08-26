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
 * 从 source 抓取 target 中定义的值
 *
 * @example
 *
   ```js
   let foo = { a: 1, b: 2, c: null };
   let bar = { a: 3, b: 4, c: 5, d: 6 };

   fetch(foo, bar); // { a: 3, b: 4, c: 5 }
   fetch(foo, bar, true); // { a: 1, b: 2, c: 5 }
   ```
 */
export function fetch<T extends Object, S extends Object>(target: T, source: S, fetchNullOnly: boolean = false): T {
  for (let prop in target) {
    if (target.hasOwnProperty(prop)) {
      if (fetchNullOnly && target[prop] !== null) break;
      if (source.hasOwnProperty(prop)) {
        target[prop] = source[prop as any];
      }
    }
  }
  return target;
}
