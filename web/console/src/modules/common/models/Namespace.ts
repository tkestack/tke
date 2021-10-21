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
import { Cluster } from './Cluster';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name?: string;

  /**命名空间 */
  namespace?: string;

  /** 用在业务侧的命名空间全名 */
  namespaceValue?: string;

  /**描述 */
  description?: string;

  /**状态 */
  status?: string;

  /**创建时间 */
  createdAt?: string;

  metadata?;

  cluster?: Cluster;
}

export interface NamespaceFilter {
  /**业务 */
  projectName?: string;

  /**集群Id */
  clusterId?: string;

  /**地域Id */
  regionId?: number;
}

export interface NamespaceOperator {
  /**集群Id */
  clusterId?: string;

  /**地域Id */
  regionId?: number;
}
