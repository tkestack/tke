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
import { Validation } from '../../common/models';

export interface Role extends Identifiable {
  metadata?: {
    /** 名称 */
    name: string;
    /** 创建时间 */
    creationTimestamp?: string;
  };
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
    /** 策略 */
    policies?: string[];
  };
  [props: string]: any;
}

export interface RolePlain extends Identifiable {
  /** id */
  name?: string;
  /** 名称 */
  displayName?: string;
  /** 描述 */
  description: string;
}

/** 用于列表查询 */
export interface RoleFilter {
  /** 目标资源，如localgroup,policy,localidentity */
  resource: string;
  /** 资源id */
  resourceID: string;
  /** 关联/解关联操作后的回调函数 */
  callback?: () => void;
}

/** 用于单个查询 */
export interface RoleInfoFilter {
  name: string;
}

export interface RoleCreation extends Identifiable {
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
    /** 策略 */
    policies?: [];
  };
  status?: {
    /** 用户 */
    users?: {
      id: string;
    }[];
    /** 用户组 */
    groups?: {
      id: string;
    }[];
  };
}

export interface RoleEditor extends Identifiable {
  metadata?: {
    /** 名称 */
    name: string;
    /** 创建时间 */
    creationTimestamp?: string;
  };
  spec: {
    /** 展示名 */
    displayName: string;
    /** 描述 */
    description: string;
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
}

export interface RoleAssociation extends Identifiable{
  /** 后端绑定接口不支持同时绑定和解绑，因此，这里设计灵活点，存储原始数据和即将增删的数据 */
  /** 最新数据 */
  roles?: RolePlain[];
  /** 原来数据 */
  originRoles?: RolePlain[];
  /** 新增数据 */
  addRoles?: RolePlain[];
  /** 删除数据 */
  removeRoles?: RolePlain[];
}
