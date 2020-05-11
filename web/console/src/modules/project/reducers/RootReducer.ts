import { FFReduxActionName } from './../constants/Config';
import { combineReducers } from 'redux';

import {
  createFFListReducer,
  generateFetcherReducer,
  generateWorkflowReducer,
  reduceToPayload,
  createFFObjectReducer
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { initNamespaceEdition, initProjectEdition } from '../constants/Config';
import { router } from '../router';

export const RootReducer = combineReducers({
  route: router.getReducer(),

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
