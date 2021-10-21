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
 * 在指定函数后面追加一个函数
 */
export function appendFunction<T extends Function, R>(origin: (...args) => R, append: (r: R, ...args) => R) {
  return function appended<T>(...args) {
    let result = origin.apply(this, args);
    return append.apply(this, [result].concat(args));
  };
}
