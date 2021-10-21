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
import { CreateResource } from 'src/modules/cluster/models';

import { FetcherState, FFListModel, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers/Router';
import {
    Cluster, ClusterFilter, Namespace, NamespaceFilter, Region, RegionFilter, Resource,
    ResourceFilter
} from '../../common/models';
import { LogDaemonset, LogDaemonSetFliter, LogDaemonSetStatus } from './LogDaemonset';
import { LogStashEdit } from './LogStashEdit';
import { Log, LogFilter } from './LogStatsh';

type LogOpenDeployWorkflow = WorkflowState<CreateResource, number>;
type ModifyLogStashWorkflow = WorkflowState<CreateResource, number>;
type InlineDeleteLog = WorkflowState<CreateResource, number>;
export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域的查询 */
  regionQuery?: QueryState<RegionFilter>;

  /** 地域的列表 */
  regionList?: FetcherState<RecordSet<Region>>;

  /** 地域的选择 */
  regionSelection?: Region;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;

  /** 集群列表的查询的查询 */
  clusterQuery?: QueryState<ClusterFilter>;

  /** 集群列表 */
  clusterList?: FetcherState<RecordSet<Cluster>>;

  /** 集群的选择 */
  clusterSelection?: Cluster[];

  /** 集群的版本 */
  clusterVersion?: string;

  /**namespace搜索过滤*/
  namespaceSelection?: string;

  /** namespace列表 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  /** namespace查询条件 */
  namespaceQuery?: QueryState<NamespaceFilter>;

  /** log日志收集规则的查询 */
  logQuery?: QueryState<LogFilter>;

  /** log日志收集器列表 */
  logList?: FetcherState<RecordSet<Log>>;

  /** 当前选中的log日志收集器 */
  logSelection?: Log[];

  /**日志采集器 */
  logDaemonset?: FetcherState<RecordSet<LogDaemonset>>;

  /**日志采集器查询条件 */
  logDaemonsetQuery?: QueryState<LogDaemonSetFliter>;

  /** 是否已经开通了日志采集器 */
  isOpenLogStash?: boolean;

  /** 判断当前的日志采集器规则是否正常 */
  isDaemonsetNormal?: LogDaemonSetStatus;

  /** 开通日志服务操作的工作流 */
  authorizeOpenLogFlow?: LogOpenDeployWorkflow;

  /** 创建日志采集规则的相关编辑项 */
  logStashEdit?: LogStashEdit;

  /** 创建、修改日志采集规则的工作流 */
  modifyLogStashFlow?: ModifyLogStashWorkflow;

  /** 删除日志的工作流 */
  inlineDeleteLog?: InlineDeleteLog;

  /** 是否正在拉取日志采集规则详情 */
  isFetchDoneSpecificLog?: boolean;

  /** 集群下的addon列表 */
  openAddon?: FFListModel<Resource, ResourceFilter>;
}
