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
import { ReduxAction, extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/ff-redux';
import { RootState, RoleInfoFilter, RoleEditor, Role } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initRoleEditorState } from '../../constants/initState';
import { RoleValidateSchema } from '../../constants/RoleValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
type GetState = () => RootState;

/**
 * 修改角色
 */
const updateRoleWorkflow = generateWorkflowActionCreator<Role, void>({
  actionType: ActionTypes.UpdateRole,
  workflowStateLocator: (state: RootState) => state.roleUpdateWorkflow,
  operationExecutor: WebAPI.updateRole,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      if (isSuccessWorkflow(getState().roleUpdateWorkflow)) {
        //表示编辑模式结束
        let { roleEditor } = getState();
        dispatch({
          type: ActionTypes.UpdateRoleEditorState,
          payload: Object.assign({}, roleEditor, { v_editing: false })
        });
      }
      /** 结束工作流 */
      dispatch(detailActions.updateRoleWorkflow.reset());
    }
  }
});

const restActions = {
  updateRoleWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: RoleValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.roleEditor;
    },
    validatorStateLocation: (store: RootState) => {
      return store.roleValidator;
    }
  }),

  fetchRole: (filter: RoleInfoFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchRole(filter);
      let editor: RoleEditor = response;
      dispatch({
        type: ActionTypes.UpdateRoleEditorState,
        payload: editor
      });
    };
  },

  /** 更新状态 */
  updateEditorState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { roleEditor } = getState();
      dispatch({
        type: ActionTypes.UpdateRoleEditorState,
        payload: Object.assign({}, roleEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateRoleEditorState,
      payload: initRoleEditorState
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
export const detailActions = extend({}, restActions);