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

const manager = seajs.require('manager');

type Callback<T> = (data: T) => any;

function promisify<T>(receive: (cb: Callback<T>) => any): Promise<T> {
  return new Promise<T>(resolve => receive(resolve));
}

export function queryWhitelist(whiteKey: string[]) {
  return new Promise<nmc.WhitelistMap>(resolve => manager.queryWhiteList({ whiteKey }, resolve));
}

export async function isInWhitelist(key: string) {
  const commonData = await promisify(manager.getComData.bind(manager) as typeof manager.getComData);
  if (!commonData.userInfo) return false;
  const uin = String(commonData.userInfo.ownerUin);

  const cached = await promisify(manager.getAllWhiteList.bind(manager) as typeof manager.getAllWhiteList);
  if (cached[key] && cached[key].indexOf(uin) > -1) {
    return true;
  }

  const query = await queryWhitelist([key]);
  return query[key] && query[key].indexOf(uin) > -1;
}
