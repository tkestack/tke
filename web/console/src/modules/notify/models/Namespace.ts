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

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name: string;
}

export interface NamespaceFilter {
  /** clusterId */
  clusterId?: string;

  /** regionId */
  regionId?: string;

  /**是否选择默认值 */

  default?: boolean;
}
