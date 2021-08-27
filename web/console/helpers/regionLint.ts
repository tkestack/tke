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
type RegionId = number | string;
type LooseRegionId = number | string;
interface Region {
  value: RegionId;
}
/**
 * 根据地域列表、外界传入的regionId 以及 当前的regionId，来判断应该切换到哪个regionId
 * @author cluezhang
 * @param list 支持的地域列表，每项的value对应数字的regionId
 * @param regionId 想要切换的地域ID，可以为0、空串
 * @param otherRegionIds 其它备选地域ID
 */
export function assureRegion(list: Region[], regionId: LooseRegionId, ...otherRegionIds: LooseRegionId[]): RegionId {
  let id: RegionId = +regionId;
  // 只在地域列表加载完成后才进行处理
  if (!list.length) {
    return id;
  }
  // 如果地域不支持
  if (!list.some(region => region.value === id)) {
    // 可以tips提示之
    seajs.require('tips').error('所选地域不支持该功能，已自动为您切换到支持的地域');
    // 则先尝试当前的值
    if (otherRegionIds.length) {
      /* eslint-disable */
      id = assureRegion.apply(null, [list, ...otherRegionIds]);
      /* eslint-enable */
    } else {
      // 如果没有设置地域，用第一个
      id = list[0].value;
    }
  }
  return id;
}
