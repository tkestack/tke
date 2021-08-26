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

import { FetcherState, FFListModel, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { Region, RegionFilter, Resource, ResourceFilter, ResourceInfo } from '../../common';
import { CreateResource, PeEdit } from './';

type PeModifyWorkflow = WorkflowState<CreateResource, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

  /** PersistentEvent的列表 */
  peList?: FetcherState<RecordSet<Resource>>;

  /** PersistentEvent的查询 */
  peQuery?: QueryState<ResourceFilter>;

  /** PersistentEvent的选择 */
  peSelection?: Resource[];

  /** peEdit */
  peEdit?: PeEdit;

  /** resourceInfo */
  resourceInfo?: ResourceInfo;

  /** 设置持久化事件 创建的操作流 */
  modifyPeFlow?: PeModifyWorkflow;

  /** 删除持久化事件 Delete的操作流 */
  deletePeFlow?: PeModifyWorkflow;
}

export interface FetchPeList {
  /** clusterId */
  clusterId?: string;

  /** peName */
  peName?: string;
}
