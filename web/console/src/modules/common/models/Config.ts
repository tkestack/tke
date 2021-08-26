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

interface Base extends Identifiable {
  /**配置文件Id */
  configId?: string;

  /**配置文件名称 */
  name?: string;

  /**创建时间 */
  createdAt?: string;
}

export interface Config extends Base {
  /**版本数量 */
  totalCount?: number;

  /**修改时间 */
  updatedAt?: string;
}

export interface ConfigFilter {
  /**搜索字段 */
  search?: string;
}

export interface Version extends Base {
  /**版本名称 */
  version?: string;

  /**配置数据 */
  data?: string;

  /**描述 */
  description?: string;
}

export interface VersionFilter {
  /**配置文件Id */
  configId?: string;

  /**地域 */
  regionId?: number | string;
}

export interface Variable extends Identifiable {
  /**变量名称 */
  key?: string;

  /**变量值 */
  value?: string;

  /**变量类型 */
  type?: string;

  /**变量名是否是规范化变量 */
  isLegal?: boolean;
}

export interface VariableFilter {
  /**配置文件Id */
  configId?: string;

  /**版本名称 */
  version?: string;
}
