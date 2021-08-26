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

export interface Repo extends Identifiable {
  apiVersion?: string;

  kind?: boolean;

  /** 元数据 */
  metadata?: {
    annotations?: any;
    clusterName?: string;
    creationTimestamp?: string;
    deletionGracePeriodSeconds?: number;
    deletionTimestamp?: string;
    finalizers?: string[];
    generateName?: string;
    generation?: string;
    labels?: any;
    managedFields?: any[];
    /** 键值 name */
    name?: string;
    namespace?: string;
    ownerReferences?: any;
    resourceVersion?: string;
    selfLink?: string;
    uid?: string;
  };

  spec?: {
    /** 描述 */
    displayName?: string;
    /** 仓库名称 */
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    locked?: boolean;
    /** 仓库数量 */
    repoCount?: number;
  };
}

export interface RepoFilter {
  /** 仓库名称 */
  name?: string;
  /** 描述 */
  displayName?: string;
}

export interface RepoCreation extends Identifiable {
  /** 描述 */
  displayName?: 'string';
  /** 仓库名称 */
  name?: 'string';
  v_name?: Validation;
  /** 公开或私有 */
  visibility?: 'Public' | 'Private';
}
