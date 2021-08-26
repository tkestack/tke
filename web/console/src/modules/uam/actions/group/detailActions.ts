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
import { RootState, GroupInfoFilter, GroupEditor, Group } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initGroupEditorState } from '../../constants/initState';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
import { GroupValidateSchema } from '../../constants/GroupValidateConfig';
import { listActions } from './listActions';
type GetState = () => RootState;

/**
 * 修改用户组
 */
const updateGroupWorkflow = generateWorkflowActionCreator<Group, void>({
  actionType: ActionTypes.UpdateGroup,
  workflowStateLocator: (state: RootState) => state.groupUpdateWorkflow,
  operationExecutor: WebAPI.updateGroup,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      if (isSuccessWorkflow(getState().groupUpdateWorkflow)) {
        //表示编辑模式结束
        let { groupEditor } = getState();
        dispatch({
          type: ActionTypes.UpdateGroupEditorState,
          payload: Object.assign({}, groupEditor, { v_editing: false })
        });
      }
      /** 结束工作流 */
      dispatch(detailActions.updateGroupWorkflow.reset());
      let count = 0;
      const timer = setInterval(() => {
        dispatch(listActions.poll());
        ++count;
        if (count >= 3) {
          clearInterval(timer);
        }
      }, 1500);
    }
  }
});

const restActions = {
  updateGroupWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: GroupValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.groupEditor;
    },
    validatorStateLocation: (store: RootState) => {
      return store.groupValidator;
    }
  }),

  fetchGroup: (filter: GroupInfoFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchGroup(filter);
      let editor: GroupEditor = response;
      dispatch({
        type: ActionTypes.UpdateGroupEditorState,
        payload: editor
      });
    };
  },

  /** 更新状态 */
  updateEditorState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { groupEditor } = getState();
      dispatch({
        type: ActionTypes.UpdateGroupEditorState,
        payload: Object.assign({}, groupEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateGroupEditorState,
      payload: initGroupEditorState
    };
  },

  /** 离开更新页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(GroupValidateSchema.formKey),
      payload: {}
    };
  }
};
export const detailActions = extend({}, restActions);
