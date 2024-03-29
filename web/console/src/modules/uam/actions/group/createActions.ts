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
import { Group, RootState } from '../../models/index';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initGroupCreationState } from '../../constants/initState';
import { GroupValidateSchema } from '../../constants/GroupValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
import { listActions } from './listActions';
type GetState = () => RootState;

/**
 * 增加用户组
 */
const addGroupWorkflow = generateWorkflowActionCreator<Group, void>({
  actionType: ActionTypes.AddGroup,
  workflowStateLocator: (state: RootState) => state.groupAddWorkflow,
  operationExecutor: WebAPI.addGroup,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      let { groupAddWorkflow, route } = getState();
      if (isSuccessWorkflow(groupAddWorkflow)) {
        router.navigate({ module: 'user', sub: 'group' }, route.queries);
        //进入列表时自动加载
        //退出状态页面时自动清理状态
      }
      /** 结束工作流 */
      dispatch(createActions.addGroupWorkflow.reset());
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
  addGroupWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: GroupValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.groupCreation;
    },
    validatorStateLocation: (store: RootState) => {
      return store.groupValidator;
    }
  }),

  /** 更新状态 */
  updateCreationState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { groupCreation } = getState();
      dispatch({
        type: ActionTypes.UpdateGroupCreationState,
        payload: Object.assign({}, groupCreation, obj)
      });
    };
  },

  /** 离开创建页面，清除Creation当中的内容 */
  clearCreationState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateGroupCreationState,
      payload: initGroupCreationState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(GroupValidateSchema.formKey),
      payload: {}
    };
  }
};

export const createActions = extend({}, restActions);
