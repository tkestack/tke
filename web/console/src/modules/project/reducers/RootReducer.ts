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

import { combineReducers } from 'redux';

import {
  createFFListReducer,
  createFFObjectReducer,
  generateFetcherReducer,
  generateWorkflowReducer,
  reduceToPayload
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName, initNamespaceEdition, initProjectEdition } from '../constants/Config';
import { router } from '../router';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  platformType: reduceToPayload(ActionType.PlatformType, 'init'),

  userInfo: createFFObjectReducer(FFReduxActionName.UserInfo),

  userManagedProjects: createFFListReducer(FFReduxActionName.UserManagedProjects),

  project: createFFListReducer('project'),

  projectEdition: reduceToPayload(ActionType.UpdateProjectEdition, initProjectEdition),

  createProject: generateWorkflowReducer({
    actionType: ActionType.CreateProject
  }),

  editProjectName: generateWorkflowReducer({
    actionType: ActionType.EditProjectName
  }),

  editProjectManager: generateWorkflowReducer({
    actionType: ActionType.EditProjectManager
  }),

  editProjecResourceLimit: generateWorkflowReducer({
    actionType: ActionType.EditProjecResourceLimit
  }),

  deleteProject: generateWorkflowReducer({
    actionType: ActionType.DeleteProject
  }),

  manager: createFFListReducer('manager'),

  namespace: createFFListReducer('namespace'),

  namespaceEdition: reduceToPayload(ActionType.UpdateNamespaceEdition, initNamespaceEdition),

  createNamespace: generateWorkflowReducer({
    actionType: ActionType.CreateNamespace
  }),

  editNamespaceResourceLimit: generateWorkflowReducer({
    actionType: ActionType.EditNamespaceResourceLimit
  }),

  deleteNamespace: generateWorkflowReducer({
    actionType: ActionType.DeleteNamespace
  }),

  region: createFFListReducer('region'),

  cluster: createFFListReducer('cluster'),

  /** 设置管理员*/
  modifyAdminstrator: generateWorkflowReducer({
    actionType: ActionType.ModifyAdminstrator
  }),

  /**当前管理员 */
  adminstratorInfo: reduceToPayload(ActionType.FetchAdminstratorInfo, {}),

  /** 用户相关*/
  userList: createFFListReducer(ActionType.UserList),
  addUserWorkflow: generateWorkflowReducer({
    actionType: ActionType.AddUser
  }),

  /** 关联策略相关 */
  policyPlainList: createFFListReducer(ActionType.PolicyPlainList),

  projectUserInfo: createFFObjectReducer(FFReduxActionName.ProjectUserInfo),

  detailProject: createFFListReducer('detailProject'),

  addExistMultiProject: generateWorkflowReducer({
    actionType: ActionType.AddExistMultiProject
  }),
  deleteParentProject: generateWorkflowReducer({
    actionType: ActionType.DeleteParentProject
  }),

  projectDetail: reduceToPayload(ActionType.ProjectDetail, null),

  namespaceKubectlConfig: createFFObjectReducer(FFReduxActionName.NamespaceKubectlConfig),

  migrateNamesapce: generateWorkflowReducer({
    actionType: ActionType.MigrateNamesapce
  })
});
