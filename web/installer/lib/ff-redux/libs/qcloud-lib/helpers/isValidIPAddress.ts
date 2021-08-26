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
 * 检查一个 IP 地址是否为合法的 IP 地址
 * */
export function isValidIPAddress(ipAddress: string) {
  if (!ipAddress) return false;

  let newIpAddress = ipAddress.trim();

  const segments = newIpAddress.split('.');
  if (segments.length !== 4) return false;

  return segments.reduce((prev, curr) => {
    const value = parseInt(curr, 10);
    return prev && value >= 0 && value <= 255 && String(value) === curr;
  }, true);
}
