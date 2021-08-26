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

import {
  User,
  UserFilter,
  Strategy,
  StrategyFilter,
  Category,
  Role,
  RoleCreation,
  RoleEditor,
  RoleFilter,
  RolePlain,
  RoleAssociation,
  Group,
  GroupCreation,
  GroupEditor,
  GroupFilter,
  GroupPlain,
  GroupAssociation,
  UserPlain,
  CommonUserFilter,
  CommonUserAssociation,
  PolicyEditor,
  PolicyPlain,
  PolicyFilter,
  PolicyAssociation
} from './index';
import { ResourceFilter } from '@src/modules/common';
import { FetcherState, FFListModel, OperationResult, RecordSet, WorkflowState } from '@tencent/ff-redux';
import { Validation, ValidatorModel } from '@tencent/ff-validator';
import { RouteState } from '../../../../helpers';

type userWorkflow = WorkflowState<User, any>;
type strategyWorkflow = WorkflowState<Strategy, any>;
type associateWorkflow = WorkflowState<{ id: string; userNames: [] }, any>;

export interface RootState {
  /** 用户信息 */
  userList?: FFListModel<User, UserFilter>;
  addUserWorkflow?: userWorkflow;
  removeUserWorkflow?: userWorkflow;
  user?: User;
  filterUsers?: User[];
  getUser?: OperationResult<User>;
  updateUser?: FetcherState<RecordSet<any>>;
  userStrategyList?: FFListModel<Strategy, ResourceFilter>;

  /** 策略相关 */
  strategyList?: FFListModel<Strategy, StrategyFilter>;
  addStrategyWorkflow?: strategyWorkflow;
  removeStrategyWorkflow?: strategyWorkflow;
  associatedUsersList?: FFListModel<User, UserFilter>;
  removeAssociatedUser?: associateWorkflow;
  addAssociatedUser?: associateWorkflow;
  getStrategy?: OperationResult<Strategy>;
  updateStrategy?: FetcherState<RecordSet<any>>;

  /** 类别 */
  categoryList?: FetcherState<RecordSet<Category>>;

  /** 角色相关 */
  roleList?: FFListModel<Role, RoleFilter>;
  roleCreation?: RoleCreation;
  roleEditor?: RoleEditor;
  roleValidator?: ValidatorModel;
  roleAddWorkflow?: WorkflowState<Role, any>;
  roleUpdateWorkflow?: WorkflowState<Role, any>;
  roleRemoveWorkflow?: WorkflowState<Role, any>;

  /** 关联角色相关，单独设置，不赋予任何场景相关的命名 */
  rolePlainList?: FFListModel<RolePlain, RoleFilter>;
  roleAssociatedList?: FFListModel<RolePlain, RoleFilter>;
  associateRoleWorkflow?: WorkflowState<RoleAssociation, any>;
  disassociateRoleWorkflow?: WorkflowState<RoleAssociation, any>;
  roleAssociation?: RoleAssociation;
  roleFilter?: RoleFilter;

  /** 用户组相关 */
  groupList?: FFListModel<Group, GroupFilter>;
  groupCreation?: GroupCreation;
  groupEditor?: GroupEditor;
  groupValidator?: ValidatorModel;
  groupAddWorkflow?: WorkflowState<Group, any>;
  groupUpdateWorkflow?: WorkflowState<Group, any>;
  groupRemoveWorkflow?: WorkflowState<Group, any>;

  /** 关联用户组相关，单独设置，不赋予任何场景相关的命名 */
  groupPlainList?: FFListModel<GroupPlain, GroupFilter>;
  groupAssociatedList?: FFListModel<GroupPlain, GroupFilter>;
  associateGroupWorkflow?: WorkflowState<GroupAssociation, any>;
  disassociateGroupWorkflow?: WorkflowState<GroupAssociation, any>;
  groupAssociation?: GroupAssociation;
  groupFilter?: GroupFilter;

  /** 关联用户相关，单独设置，不赋予任何场景相关的命名 */
  userPlainList?: FFListModel<UserPlain, CommonUserFilter>;
  commonUserAssociatedList?: FFListModel<UserPlain, CommonUserFilter>;
  commonAssociateUserWorkflow?: WorkflowState<CommonUserAssociation, any>;
  commonDisassociateUserWorkflow?: WorkflowState<CommonUserAssociation, any>;
  commonUserAssociation?: CommonUserAssociation;
  commonUserFilter?: CommonUserFilter;

  /** 策略相关 */
  policyEditor?: PolicyEditor;

  /** 关联策略相关，单独设置，不赋予任何场景相关的命名 */
  policyPlainList?: FFListModel<PolicyPlain, PolicyFilter>;
  policyAssociatedList?: FFListModel<PolicyPlain, PolicyFilter>;
  associatePolicyWorkflow?: WorkflowState<PolicyAssociation, any>;
  disassociatePolicyWorkflow?: WorkflowState<PolicyAssociation, any>;
  policyAssociation?: PolicyAssociation;
  policyFilter?: PolicyFilter;

  /** 路由 */
  route?: RouteState;
}
