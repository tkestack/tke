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
    username: string;

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
