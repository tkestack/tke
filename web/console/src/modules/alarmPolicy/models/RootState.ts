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

//import { RouteState } from "@tencent/qcloud-nmc";
import { RouteState } from '../../../../helpers/Router';
import { Cluster } from '../../common';
import { User, UserFilter } from '../../uam/models';
import {
  AlarmPolicy,
  AlarmPolicyEdition,
  AlarmPolicyFilter,
  AlarmPolicyOperator,
  ClusterFilter,
  Group,
  Namespace,
  NamespaceFilter,
  Region,
  RegionFilter,
  Resource,
  ResourceFilter
} from './';
import { GroupFilter } from './Group';
import { AddonStatus } from './Addon';

type AlarmPolicyOpWorkflow = WorkflowState<AlarmPolicy, AlarmPolicyOperator>;
type AlarmPolicyCreateWorkflow = WorkflowState<AlarmPolicyEdition, AlarmPolicyOperator>;

export interface RootState {
  /**
   * 路由
   */
  route?: RouteState;

  /**
   * 地域查询
   */
  regionQuery?: QueryState<RegionFilter>;

  /**
   * 地域列表
   */
  regionList?: FetcherState<RecordSet<Region>>;

  /**
   * 选择的地域
   */
  regionSelection?: Region;

  cluster?: FFListModel<Cluster, ClusterFilter>;

  addons?: AddonStatus;

  /**当前集群命名空间 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  namespaceQuery?: QueryState<NamespaceFilter>;

  /**当前命名空间下pod列表 */
  workloadList?: FetcherState<RecordSet<Resource>>;

  workloadQuery?: QueryState<ResourceFilter>;

  clusterVersion?: string;

  alarmPolicy?: FFListModel<AlarmPolicy, AlarmPolicyFilter>;

  userList?: FFListModel<User, UserFilter>;

  /** 当前新建告警 */
  alarmPolicyEdition?: AlarmPolicyEdition;

  /** 创建告警workflow */
  alarmPolicyCreateWorkflow?: AlarmPolicyCreateWorkflow;

  /** 更新告警workflow */
  alarmPolicyUpdateWorkflow?: AlarmPolicyCreateWorkflow;

  /** 删除告警workflow */
  alarmPolicyDeleteWorkflow?: AlarmPolicyOpWorkflow;

  /**详情 */
  alarmPolicyDetail?: AlarmPolicy;

  /**组列表 */
  // groupList?: FetcherState<RecordSet<Group>>;

  channel?: FFListModel<Resource, ResourceFilter>;
  template?: FFListModel<Resource, ResourceFilter>;
  receiver?: FFListModel<Resource, ResourceFilter>;
  receiverGroup?: FFListModel<Resource, ResourceFilter>;

  groupQuery?: QueryState<GroupFilter>;

  /** 是否为国际版 */
  isI18n?: boolean;

  // /** namespace列表 */
  // namespaceList?: FetcherState<RecordSet<Namespace>>;

  // /** namespace查询条件 */
  // namespaceQuery?: QueryState<ResourceFilter>;

  /** namespace selection */
  namespaceSelection?: string;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;
}
