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

import { FFListModel, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { CreateResource, Region, RegionFilter, Resource, ResourceFilter } from '../../common';
import { Addon } from './';
import { AddonEdit } from './AddonEdit';

type ResourceModifyWorkflow = WorkflowState<CreateResource, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

  /** 集群的版本 */
  clusterVersion?: string;

  /** 集群的下的addon列表 */
  openAddon?: FFListModel<Resource, ResourceFilter>;

  /** 所有的add的列表 */
  addon?: FFListModel<Addon, ResourceFilter>;

  /** 开通addon组件 */
  editAddon?: AddonEdit;

  /** 创建resource资源的操作流程 */
  modifyResourceFlow?: ResourceModifyWorkflow;

  /** 创建多种resource资源的操作流程 */
  applyResourceFlow?: ResourceModifyWorkflow;

  /** 删除resource资源的操作流程 */
  deleteResourceFlow?: ResourceModifyWorkflow;
}
