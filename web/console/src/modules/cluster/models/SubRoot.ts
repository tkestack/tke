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

import { FetcherState, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { ResourceInfo } from '../../common/models';
import { AddonStatus } from './Addon';
import { AllocationRatioEdition } from './AllocationRatioEdition';
import { DetailResourceOption } from './DetailResourceOption';
import {
    ComputerState, ConfigMapEdit, CreateResource, LbcfEdit, NamespaceEdit, ResourceDetailState,
    ResourceEventOption, ResourceLogOption, ResourceOption, SecretEdit, ServiceEdit, SubRouter,
    SubRouterFilter, WorkloadEdit
} from './index';

type ResourceModifyWorkflow = WorkflowState<CreateResource, number | any>;

export interface SubRootState {
  /** 节点列表 */
  computerState?: ComputerState;

  /** 超售比 */
  clusterAllocationRatioEdition?: AllocationRatioEdition;

  updateClusterAllocationRatio?: ResourceModifyWorkflow;

  /** 二级菜单栏配置列表查询 */
  subRouterQuery?: QueryState<SubRouterFilter>;

  /** 二级菜单栏配置列表 */
  subRouterList?: FetcherState<RecordSet<SubRouter>>;

  /** 当前的模式 create | update | resource */
  mode?: string;

  /** 创建多种resource资源的操作流程 */
  applyResourceFlow?: ResourceModifyWorkflow;

  /**创建多种resource (使用不同的接口)的操作流程 */
  applyDifferentInterfaceResourceFlow?: ResourceModifyWorkflow;
  /** 创建resource资源的操作流 */
  modifyResourceFlow?: ResourceModifyWorkflow;

  modifyMultiResourceWorkflow?: ResourceModifyWorkflow;

  /** 删除resource资源的操作流 */
  deleteResourceFlow?: ResourceModifyWorkflow;

  /** 更新Service的访问方式、更新Ingress的转发配置等的操作流 */
  updateResourcePart?: ResourceModifyWorkflow;

  updateMultiResource?: ResourceModifyWorkflow;

  /** 当前的请求资源名称 */
  resourceName?: string;

  /** resourcrInfo */
  resourceInfo?: ResourceInfo;

  detailResourceOption?: DetailResourceOption;

  /** 通用resource 数据结构 */
  resourceOption?: ResourceOption;

  /** resource detail详情 */
  resourceDetailState?: ResourceDetailState;

  /** editService */
  serviceEdit?: ServiceEdit;

  /** editNamespace */
  namespaceEdit?: NamespaceEdit;

  /** editResource */
  workloadEdit?: WorkloadEdit;

  /** editSecret */
  secretEdit?: SecretEdit;

  lbcfEdit?: LbcfEdit;

  /** editConfigMap */
  cmEdit?: ConfigMapEdit;

  /** resourcelog的相关配置 */
  resourceLogOption?: ResourceLogOption;

  /** resourceEvent的相关配置 */
  resourceEventOption?: ResourceEventOption;

  /** 是否需要进行命名空间的拉取 */
  isNeedFetchNamespace?: boolean;

  /** 使用已有lb白名单 */
  isNeedExistedLb?: boolean;

  addons?: AddonStatus;
}
