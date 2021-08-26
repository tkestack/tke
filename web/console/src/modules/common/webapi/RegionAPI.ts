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

import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { Region, RegionFilter } from '../models';

/**获取地域列表 */
export async function fetchRegionList(query?: QueryState<RegionFilter>) {
  // 目前是hardcode，后面换成接口获取
  let regionList = [
    {
      id: uuid(),
      Remark: 'SUITABLE_TKE',
      area: '华南地区',
      name: '广州',
      value: 1
    },
    {
      id: uuid(),
      Remark: 'SUITABLE_TKE',
      area: '华东地区',
      name: '上海',
      value: 4
    }
  ];

  const result: RecordSet<Region> = {
    recordCount: regionList.length,
    records: regionList
  };

  return result;
}
