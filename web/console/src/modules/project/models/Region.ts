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

export interface Region extends Identifiable {
  /** 地域的值 */
  value: number;

  /** 地域的名称 */
  name?: string;

  /** 地域是否可用 */
  disabled?: boolean;

  /** 地域所属大区 */
  area?: string;

  /** 是否新版控制台当中的地域 */
  Remark?: string;

  //[props: string]: any;
}

export interface RegionFilter {
  /** 地域id */
  regionId?: number;
}
