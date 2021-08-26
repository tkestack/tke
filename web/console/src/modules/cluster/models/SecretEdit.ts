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

import { FetcherState, Identifiable, QueryState, RecordSet } from '@tencent/ff-redux';

import { Validation } from '../../common/models';
import { Namespace } from './Namespace';
import { ResourceFilter } from './ResourceOption';

export interface SecretEdit extends Identifiable {
  /** secret名称 */
  name?: string;
  v_name?: Validation;

  /** namespace列表 */
  nsList?: FetcherState<RecordSet<Namespace>>;

  /** namepsace列表的查询 */
  nsQuery?: QueryState<ResourceFilter>;

  /** secret类型 */
  secretType?: string;

  /** secret的数据 */
  data?: SecretData[];

  /** ns的类型，是全部命名空间 还是 指定命名空间 */
  nsType?: string;

  /** 添加第三方镜像仓库的命名空间 */
  nsListSelection?: Namespace[];

  /** 当前填写的第三方镜像仓库的域名 */
  domain?: string;
  v_domain?: Validation;

  /** 第三方镜像仓库的用户名 */
  username?: string;
  v_username?: Validation;

  /** 第三方镜像仓库的密码 */
  password?: string;
  v_password?: Validation;
}

/** secret的数据类型 */
export interface SecretData extends Identifiable {
  /** key名称 */
  keyName?: string;
  v_keyName?: Validation;

  /** value名称 */
  value?: string;
  v_value?: Validation;
}

export interface SecretEditJSONYaml {
  /** 资源的类型 */
  kind?: string;

  /** api的版本 */
  apiVersion?: string;

  /** metadata */
  metadata?: SecretMetadata;

  /** data */
  data: {
    [props: string]: string;
  };

  type?: string;
}

interface SecretMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** pvc的名称 */
  name: string;

  /** pvc的命名空间 */
  namespace?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };
}
