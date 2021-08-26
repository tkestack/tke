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


import { ReduxAction, extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/ff-redux';
import { Role, RootState } from '../../models/index';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initRoleCreationState } from '../../constants/initState';
import { RoleValidateSchema } from '../../constants/RoleValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
type GetState = () => RootState;

/**
 * 增加角色
 */
const addRoleWorkflow = generateWorkflowActionCreator<Role, void>({
  actionType: ActionTypes.AddRole,
  workflowStateLocator: (state: RootState) => state.roleAddWorkflow,
  operationExecutor: WebAPI.addRole,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { roleAddWorkflow, route } = getState();
      if (isSuccessWorkflow(roleAddWorkflow)) {
        router.navigate({ module: 'role', sub: '' }, route.queries);
        //进入列表时自动加载
        //退出状态页面时自动清理状态
      }
      /** 结束工作流 */
      dispatch(createActions.addRoleWorkflow.reset());
    }
  }
});

const restActions = {
  addRoleWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: RoleValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.roleCreation;
    },
    validatorStateLocation: (store: RootState) => {
      return store.roleValidator;
    }
  }),

  /** 更新状态 */
  updateCreationState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { roleCreation } = getState();
      dispatch({
        type: ActionTypes.UpdateRoleCreationState,
        payload: Object.assign({}, roleCreation, obj)
      });
    };
  },

  /** 离开创建页面，清除Creation当中的内容 */
  clearCreationState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateRoleCreationState,
      payload: initRoleCreationState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(RoleValidateSchema.formKey),
      payload: {}
    };
  }
};

export const createActions = extend({}, restActions);
