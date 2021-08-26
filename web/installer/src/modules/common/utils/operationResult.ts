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

import { OperationResult } from '@tencent/ff-redux';

// 返回标准操作结果
export const operationResult = <T>(target: T[] | T, error?: any): OperationResult<T>[] => {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
};
