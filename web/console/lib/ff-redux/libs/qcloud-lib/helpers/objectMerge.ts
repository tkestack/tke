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
const deepExtend = require('deep-extend');

export function merge<T, S>(target: T, source: S): T & S;
export function merge<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function merge<T, S1, S2>(target: T, source1: S1, source2: S2): T & S1 & S2;
export function merge<T, S1, S2, S3>(target: T, source1: S1, source2: S2, source3: S3): T & S1 & S2 & S3;
export function merge<T, S1, S2, S3, S4>(
  target: T,
  source1: S1,
  source2: S2,
  source3: S3,
  source4,
  S4
): T & S1 & S2 & S3 & S4;

export function merge(target: any, ...sources: any[]): any {
  return deepExtend(target, ...sources);
}
