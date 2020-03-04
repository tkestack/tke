import { createFFListReducer, generateWorkflowReducer } from '@tencent/ff-redux';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { combineReducers } from 'redux';
import * as ActionType from '../constants/ActionType';
import { router } from '../router';
import { initNamespaceEdition, initProjectEdition } from './../constants/Config';

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
  adminstratorInfo: reduceToPayload(ActionType.FetchAdminstratorInfo, {})
});
