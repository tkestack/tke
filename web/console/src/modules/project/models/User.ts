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
import { Identifiable } from '@tencent/ff-redux';

export interface User extends Identifiable {
  metadata?: {
    /** 用户的资源id */
    name: string;
  };

  /** 用户名（唯一） */
  spec: {
    /** 用户名 */
    name: string;

    /** 展示名 */
    displayName: string;

    /** 邮箱 */
    email?: string;

    /** 手机号 */
    phoneNumber?: string;

    /** 密码 */
    hashedPassword: string;

    /** 额外属性 */
    extra?: {
      /** 是否是管理员 */
      platformadmin?: boolean;
      [props: string]: any;
    };
  };

  status?: any;
}

export interface UserFilter {
  /** 业务Id */
  projectId?: string;

  /** 用户名(唯一) */
  name?: string;

  /** 展示名 */
  displayName?: string;

  /** 相关参数 */
  search?: string;

  ifAll?: boolean;

  /** 是否只拉取策策略需要绑定的用户 */
  isPolicyUser?: boolean;
}

/** 原有User在localidentities和user概念之间混用了，anyway，用于关联角色等 */
export interface UserPlain extends Identifiable {
  /** 名称 */
  name?: string;
  /** 展示名 */
  displayName?: string;
}

export interface Member extends Identifiable {
  projectId: string;
  users: any;
  policies: string[];
}

export interface UserInfo {
  extra: any;
  uid: string;
  groups: string[];
  name: string;
}
