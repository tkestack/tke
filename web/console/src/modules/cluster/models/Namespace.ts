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

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name: string;

  displayName: string;

  //业务侧使用
  clusterVersion?: string;

  clusterId?: string;

  clusterDisplayName?: string;

  clusterName?: string;

  namespace?: string;
}

/** 可视化创建的namespace的相关配置 */
export interface NamespaceEdit {
  /** name */
  name?: string;
  v_name?: Validation;

  /** 描述 */
  description?: string;
  v_description?: Validation;
}

/** 创建Namespace的时候，提交的jasonSchema */
export interface NamespaceEditJSONYaml {
  /** 资源的类型 */
  kind: string;

  /** api的版本 */
  apiVersion: string;

  /** metadata */
  metadata: NamespaceMetadata;

  /** spec */
  spec?: {};

  /** status */
  status?: {};
}

/** metadata的配置，非全部配置项 */
interface NamespaceMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** namespace的名称 */
  name: string;
}
