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
function parseDate(dateStr) {
  if (/^(\d{4})[-\s\.,\/]*(\d+)[-\s\.,\/]*(\d+)(?:\s*(\d+)\:(\d+):(\d+))?$/.test(dateStr)) {
    let year = +RegExp.$1;
    let month = +RegExp.$2;
    let day = +RegExp.$3;
    let hour = +RegExp.$4 || 0;
    let minute = +RegExp.$5 || 0;
    let second = +RegExp.$6 || 0;
    return [year, +month - 1, day, hour, minute, second];
  }
  return null;
}

export function getDateUTC(str) {
  return Date.UTC.apply(Date, parseDate(str));
}
