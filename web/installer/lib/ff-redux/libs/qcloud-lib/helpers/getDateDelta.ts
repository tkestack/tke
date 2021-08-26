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

/**Parses string formatted as YYYY-MM-DD to a Date object.
 * If the supplied string does not match the format, an
 * invalid Date (value NaN) is returned.
 * @param {string} dateStringInRange format YYYY-MM-DD, with year in
 * range of 0000-9999, inclusive.
 * @return {Date} Date object representing the string.
 */
function parseISO8601(dateStringInRange) {
  let isoExp = /^\s*(\d{4})-(\d\d)-(\d\d)\s*/,
    date = new Date(NaN),
    month,
    parts = isoExp.exec(dateStringInRange);

  if (parts) {
    month = +parts[2];
    date.setFullYear(+parts[1], month - 1, +parts[3]);
    if (month !== date.getMonth() + 1) {
      date.setTime(NaN);
    }
  }
  return date;
}

const ONE_DAY = 24 * 60 * 60 * 1000;

export function getDateDelta(date1, date2) {
  return Math.floor(+parseISO8601(date2) - +parseISO8601(date1)) / ONE_DAY;
}
