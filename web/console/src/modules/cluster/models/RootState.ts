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
import { FetcherState, FFListModel, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { Cluster, ClusterFilter, Region, RegionFilter } from '../../common/models';
import { CreateResource } from './';
import { ClusterCreationState } from './ClusterCreationState';
import { Clustercredential } from './Clustercredential';
import { CreateIC } from './CreateIC';
import { DialogState } from './DialogState';
import { Namespace } from './Namespace';
import { Resource, ResourceFilter } from './ResourceOption';
import { SubRootState } from './SubRoot';

type ResourceModifyFlow = WorkflowState<CreateResource, number>;
type CreateICFlow = WorkflowState<CreateIC, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  region?: FFListModel<Region, RegionFilter>;

  /** 集群的版本 */
  clusterVersion?: string;

  /** 集群的相关配置 */
  cluster?: FFListModel<Cluster, ClusterFilter>;

  clustercredential?: Clustercredential;

  /** 修改Token */

  updateClusterToken?: ResourceModifyFlow;

  /** 集群详情的yaml，通过tke-apiserver进行拉取 */
  clusterInfoQuery?: QueryState<ClusterFilter>;

  clusterInfoList?: FetcherState<RecordSet<Cluster>>;

  /** namespace列表 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  /** namespace查询条件 */
  namespaceQuery?: QueryState<ResourceFilter>;

  /** namespace selection */
  namespaceSelection?: string;

  /** 二级菜单数据结构 */
  subRoot?: SubRootState;

  /** 当前模式 create: 创建集群模式；expand：扩展节点模式*/
  mode?: string;

  /** 弹窗状态的 状态集合 */
  dialogState?: DialogState;

  /** 删除集群的工作流 */
  deleteClusterFlow?: ResourceModifyFlow;

  /** 创建集群状态 */
  clusterCreationState?: ClusterCreationState;

  /**创建集群工作流 */
  createClusterFlow?: ResourceModifyFlow;

  createIC?: CreateIC;

  createICWorkflow?: CreateICFlow;

  modifyClusterName?: ResourceModifyFlow;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;
}
