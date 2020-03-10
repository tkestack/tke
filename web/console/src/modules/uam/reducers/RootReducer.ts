import { combineReducers } from 'redux';

import { createFFListReducer, generateFetcherReducer, generateWorkflowReducer, reduceToPayload } from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { router } from '../router';
import {
  initRoleCreationState,
  initRoleEditorState,
  initRoleFilterState,
  initRoleAssociationState,
  initGroupCreationState,
  initGroupEditorState,
  initGroupAssociationState,
  initGroupFilterState,
  initCommonUserAssociationState,
  initCommonUserFilterState,
  initPolicyAssociationState,
  initPolicyEditorState,
  initPolicyFilterState,
} from '../constants/initState';
import { createValidatorReducer } from '@tencent/ff-validator';
import { GroupValidateSchema } from '../constants/GroupValidateConfig';
import { RoleValidateSchema } from '../constants/RoleValidateConfig';

export const RootReducer = combineReducers({
  route: router.getReducer(),
  userList: createFFListReducer(ActionTypes.UserList),
  addUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddUser
  }),
  removeUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveUser
  }),
  filterUsers: reduceToPayload(ActionTypes.FetchUserByName, []),
  getUser: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetUser,
    initialData: {}
  }),
  updateUser: generateFetcherReducer<Object>({
    actionType: ActionTypes.UpdateUser,
    initialData: {}
  }),

  userStrategyList: createFFListReducer(ActionTypes.UserStrategyList),

  strategyList: createFFListReducer(ActionTypes.StrategyList),
  addStrategyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddStrategy
  }),
  removeStrategyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveStrategy
  }),
  getStrategy: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetStrategy,
    initialData: {}
  }),
  updateStrategy: generateFetcherReducer<Object>({
    actionType: ActionTypes.UpdateStrategy,
    initialData: {}
  }),

  categoryList: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetCategories,
    initialData: {}
  }),
  associatedUsersList: createFFListReducer(ActionTypes.GetStrategyAssociatedUsers),
  removeAssociatedUser: generateWorkflowReducer({
    actionType: ActionTypes.RemoveAssociatedUser
  }),
  addAssociatedUser: generateWorkflowReducer({
    actionType: ActionTypes.AddAssociatedUser
  }),

  /** 角色相关 */
  roleList: createFFListReducer(ActionTypes.RoleList),
  roleCreation: reduceToPayload(ActionTypes.UpdateRoleCreationState, initRoleCreationState),
  roleEditor: reduceToPayload(ActionTypes.UpdateRoleEditorState, initRoleEditorState),
  roleValidator: createValidatorReducer(RoleValidateSchema),
  roleAddWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddRole
  }),
  roleUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.UpdateRole
  }),
  roleRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveRole
  }),
  /** 关联角色相关 */
  rolePlainList: createFFListReducer(ActionTypes.RolePlainList),
  roleAssociatedList: createFFListReducer(ActionTypes.RoleAssociatedList),
  associateRoleWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AssociateRole
  }),
  disassociateRoleWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.DisassociateRole
  }),
  roleAssociation: reduceToPayload(ActionTypes.UpdateRoleAssociation, initRoleAssociationState),
  roleFilter: reduceToPayload(ActionTypes.UpdateRoleFilter, initRoleFilterState),
  /** 用户组相关 */
  groupList: createFFListReducer(ActionTypes.GroupList),
  groupCreation: reduceToPayload(ActionTypes.UpdateGroupCreationState, initGroupCreationState),
  groupEditor: reduceToPayload(ActionTypes.UpdateGroupEditorState, initGroupEditorState),
  groupValidator: createValidatorReducer(GroupValidateSchema),
  groupAddWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddGroup
  }),
  groupUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.UpdateGroup
  }),
  groupRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveGroup
  }),
  /** 关联用户组相关 */
  groupPlainList: createFFListReducer(ActionTypes.GroupPlainList),
  groupAssociatedList: createFFListReducer(ActionTypes.GroupAssociatedList),
  associateGroupWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AssociateGroup
  }),
  disassociateGroupWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.DisassociateGroup
  }),
  groupAssociation: reduceToPayload(ActionTypes.UpdateGroupAssociation, initGroupAssociationState),
  groupFilter: reduceToPayload(ActionTypes.UpdateGroupFilter, initGroupFilterState),
  /** 关联用户相关 */
  userPlainList: createFFListReducer(ActionTypes.UserPlainList),
  commonUserAssociatedList: createFFListReducer(ActionTypes.CommonUserAssociatedList),
  commonAssociateUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.CommonAssociateUser
  }),
  commonDisassociateUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.CommonDisassociateUser
  }),
  commonUserAssociation: reduceToPayload(ActionTypes.UpdateCommonUserAssociation, initCommonUserAssociationState),
  commonUserFilter: reduceToPayload(ActionTypes.UpdateCommonUserFilter, initCommonUserFilterState),
  /** 策略相关 */
  policyEditor: reduceToPayload(ActionTypes.UpdatePolicyEditorState, initPolicyEditorState),
  /** 关联策略相关 */
  policyPlainList: createFFListReducer(ActionTypes.PolicyPlainList),
  policyAssociatedList: createFFListReducer(ActionTypes.PolicyAssociatedList),
  associatePolicyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AssociatePolicy
  }),
  disassociatePolicyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.DisassociatePolicy
  }),
  policyAssociation: reduceToPayload(ActionTypes.UpdatePolicyAssociation, initPolicyAssociationState),
  policyFilter: reduceToPayload(ActionTypes.UpdatePolicyFilter, initPolicyFilterState),
});
