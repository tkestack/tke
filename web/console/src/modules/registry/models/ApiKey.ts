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

export interface ApiKey extends Identifiable {
  metadata?: {
    name: string;
  };

  spec: {
    /** key 内容 */
    apiKey?: string;

    /** 描述 */
    description?: string;

    /** 过期时间 */
    expire_at?: string;

    /** 创建时间 */
    issue_at?: string;
  };

  status: {
    /** 启用|禁用 */
    disabled?: boolean;

    /** 是否过期 */
    expired?: boolean;
  };

  /** 软删除标记? */
  deleted?: boolean;
}

export interface ApiKeyFilter {
  /** 描述字段 */
  desc?: string;
}

export interface ApiKeyCreation extends Identifiable {
  /** key 描述 */
  description?: string;
  /** key 过期时间，单位 h */
  expire?: number;
  v_expire: Validation;
  unit?: string;
}
