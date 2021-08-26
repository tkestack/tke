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

import { FetcherState, FFListModel, OperationResult, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { FFObjectModel } from '../../../../lib/ff-redux/src/object/Model';
import { Region, RegionFilter, ResourceFilter } from '../../common/models';
import { Resource } from '../../common/models/Resource';
import {
  Cluster,
  ClusterFilter,
  Manager,
  ManagerFilter,
  Member,
  Namespace,
  NamespaceEdition,
  NamespaceFilter,
  NamespaceOperator,
  PolicyFilter,
  PolicyPlain,
  Project,
  ProjectEdition,
  ProjectFilter,
  User,
  UserFilter
} from './index';
import { NamespaceCert } from './Namespace';
import { ProjectUserMap, UserManagedProject, UserManagedProjectFilter } from './Project';
import { UserInfo } from './User';

type ProjectWorkflow = WorkflowState<Project, string>;
type ProjectEditWorkflow = WorkflowState<ProjectEdition, void>;
type NamespaceWorkflow = WorkflowState<Namespace, NamespaceOperator>;
type NamespaceEditWorkflow = WorkflowState<NamespaceEdition, NamespaceOperator>;
type userWorkflow = WorkflowState<Member, any>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  platformType?: string;

  userInfo?: FFObjectModel<UserInfo, any>;

  userManagedProjects?: FFListModel<UserManagedProject, UserManagedProjectFilter>;

  project?: FFListModel<Project, ProjectFilter>;

  /** 业务编辑参数 */
  projectEdition?: ProjectEdition;

  /** 创建业务工作流 */
  createProject?: ProjectEditWorkflow;

  /** 编辑业务名称工作流 */
  editProjectName?: ProjectEditWorkflow;

  /** 编辑业务负责人工作流 */
  editProjectManager?: ProjectEditWorkflow;

  /** 编辑业务描述工作流 */
  editProjecResourceLimit?: ProjectEditWorkflow;

  /** 删除业务工作流 */
  deleteProject?: ProjectWorkflow;

  namespace?: FFListModel<Namespace, NamespaceFilter>;

  /** Namespace编辑参数 */
  namespaceEdition?: NamespaceEdition;

  /** 创建业务工作流 */
  createNamespace?: NamespaceEditWorkflow;

  /** 创建业务工作流 */
  editNamespaceResourceLimit?: NamespaceEditWorkflow;

  /** 删除业务工作流 */
  deleteNamespace?: NamespaceWorkflow;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表*/
  cluster?: FFListModel<Cluster, ClusterFilter>;

  /** 负责人列表 */
  manager?: FFListModel<Manager, ManagerFilter>;

  /** 设置管理员*/
  modifyAdminstrator?: ProjectEditWorkflow;

  /**当前管理员 */
  adminstratorInfo?: Resource;

  /** 用户信息 */
  userList?: FFListModel<User, UserFilter>;
  addUserWorkflow?: userWorkflow;

  /** 关联策略相关，单独设置，不赋予任何场景相关的命名 */
  policyPlainList?: FFListModel<PolicyPlain, PolicyFilter>;

  /**project和用户信息的映射 */
  projectUserInfo?: FFObjectModel<ProjectUserMap, ProjectFilter>;

  detailProject?: FFListModel<Project, ProjectFilter>;

  addExistMultiProject?: ProjectWorkflow;

  deleteParentProject?: ProjectWorkflow;

  projectDetail?: Project;

  /**namespaceTable */
  namespaceKubectlConfig?: FFObjectModel<NamespaceCert, NamespaceFilter>;

  migrateNamesapce?: NamespaceWorkflow;
}
