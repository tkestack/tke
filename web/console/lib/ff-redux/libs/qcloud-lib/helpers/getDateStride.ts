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
import { getDateDelta } from './getDateDelta';

function rangeLength(from: string, to: string) {
  return Math.abs(getDateDelta(from, to)) + 1;
}

export function getDateStride(dateFrom: string, dateTo: string) {
  let length = rangeLength(dateFrom, dateTo);

  /*
        小于等于1天，5分钟粒度数据。
        大于1天小于等于7天，1小时粒度数据。
        大于7天，1天粒度数据。
    */
  if (length <= 1) {
    return 5;
  } else if (length <= 7) {
    return 60;
  } else {
    return 24 * 60;
  }
}
