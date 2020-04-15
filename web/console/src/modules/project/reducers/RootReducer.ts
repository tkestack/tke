import { combineReducers } from 'redux';

import { createFFListReducer, generateFetcherReducer, generateWorkflowReducer, reduceToPayload } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { initNamespaceEdition, initProjectEdition } from '../constants/Config';
import { router } from '../router';
import {
  initPolicyAssociationState,
  initPolicyEditorState,
  initPolicyFilterState,
  initRoleAssociationState,
  initRoleFilterState,
  initGroupAssociationState,
  initGroupFilterState
} from '../constants/initState';

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
  removeUserWorkflow: generateWorkflowReducer({
    actionType: ActionType.RemoveUser
  }),
  filterUsers: reduceToPayload(ActionType.FetchUserByName, []),
  getUser: generateFetcherReducer<Object>({
    actionType: ActionType.GetUser,
    initialData: {}
  }),
  updateUser: generateFetcherReducer<Object>({
    actionType: ActionType.UpdateUser,
    initialData: {}
  }),

  /** 策略相关 */
  policyEditor: reduceToPayload(ActionType.UpdatePolicyEditorState, initPolicyEditorState),
  /** 关联策略相关 */
  policyPlainList: createFFListReducer(ActionType.PolicyPlainList),
  policyAssociatedList: createFFListReducer(ActionType.PolicyAssociatedList),
  associatePolicyWorkflow: generateWorkflowReducer({
    actionType: ActionType.AssociatePolicy
  }),
  disassociatePolicyWorkflow: generateWorkflowReducer({
    actionType: ActionType.DisassociatePolicy
  }),
  policyAssociation: reduceToPayload(ActionType.UpdatePolicyAssociation, initPolicyAssociationState),
  policyFilter: reduceToPayload(ActionType.UpdatePolicyFilter, initPolicyFilterState),

  /** 关联角色相关 */
  rolePlainList: createFFListReducer(ActionType.RolePlainList),
  roleAssociatedList: createFFListReducer(ActionType.RoleAssociatedList),
  associateRoleWorkflow: generateWorkflowReducer({
    actionType: ActionType.AssociateRole
  }),
  disassociateRoleWorkflow: generateWorkflowReducer({
    actionType: ActionType.DisassociateRole
  }),
  roleAssociation: reduceToPayload(ActionType.UpdateRoleAssociation, initRoleAssociationState),
  roleFilter: reduceToPayload(ActionType.UpdateRoleFilter, initRoleFilterState),

  /** 关联用户组相关 */
  groupPlainList: createFFListReducer(ActionType.GroupPlainList),
  groupAssociatedList: createFFListReducer(ActionType.GroupAssociatedList),
  associateGroupWorkflow: generateWorkflowReducer({
    actionType: ActionType.AssociateGroup
  }),
  disassociateGroupWorkflow: generateWorkflowReducer({
    actionType: ActionType.DisassociateGroup
  }),
  groupAssociation: reduceToPayload(ActionType.UpdateGroupAssociation, initGroupAssociationState),
  groupFilter: reduceToPayload(ActionType.UpdateGroupFilter, initGroupFilterState),
});
