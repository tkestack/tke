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

import { Identifiable } from '@tencent/ff-redux';

export interface Zone extends Identifiable {
  /**
   * 可用区ID
   */
  id: string;

  /**
   * 可用区名称
   * */
  name?: string;

  /**
   * 可用区状态
   * */
  status?: number;

  /**
   * cbs
   */
  cbs?: number;

  /**
   * 是否默认选中
   */
  isdefault?: boolean;

  /**
   * 是否可用
   */
  disable?: boolean;

  /**
   * 提示
   */
  tip?: string;

  /**
   * 是否下线
   */
  offline?: boolean;

  /**
   * 白名单
   */
  whiteList?: string;
}

export interface ZoneInfo {
  default?: number;

  id?: number;

  name?: string;

  payMode?: string[];

  zoneId?: string;
}

export interface ZoneFilter {
  regionId?: number | string;

  devPayMode?: string;
}

export interface ZoneQuotaFilter {
  cvmPayMode?: number;
}
