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

/** 原有User在localidentities和user概念之间混用了，anyway，用于关联角色等 */
export interface UserPlain extends Identifiable {
  /** 名称 */
  name?: string;
  /** 展示名 */
  displayName?: string;
}

export interface CommonUserFilter {
  /** 目标资源，如localgroup/role/policy */
  resource: string;
  /** 资源id */
  resourceID: string;
  /** 关联/解关联操作后的回调函数 */
  callback?: () => void;
}

export interface CommonUserAssociation extends Identifiable{
  /** 后端绑定接口不支持同时绑定和解绑，因此，这里设计灵活点，存储原始数据和即将增删的数据 */
  /** 最新数据 */
  users?: UserPlain[];
  /** 原来数据 */
  originUsers?: UserPlain[];
  /** 新增数据 */
  addUsers?: UserPlain[];
  /** 删除数据 */
  removeUsers?: UserPlain[];
}
