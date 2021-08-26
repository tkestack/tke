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
  ReduxAction,
  extend,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  getWorkflowStatistics
} from '@tencent/ff-redux';
import {
  RootState,
  ChartGroupFilter,
  ChartGroupEditor,
  ChartGroup,
  ChartGroupDetailFilter,
  UserPlain
} from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { initChartGroupEditorState } from '../../constants/initState';
import { ChartGroupValidateSchema } from '../../constants/ChartGroupValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
type GetState = () => RootState;

/**
 * 修改仓库
 */
const updateChartGroupWorkflow = generateWorkflowActionCreator<ChartGroup, void>({
  actionType: ActionTypes.UpdateChartGroup,
  workflowStateLocator: (state: RootState) => state.chartGroupUpdateWorkflow,
  operationExecutor: WebAPI.updateChartGroup,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      if (isSuccessWorkflow(getState().chartGroupUpdateWorkflow)) {
        //表示编辑模式结束
        let { chartGroupEditor } = getState();
        dispatch({
          type: ActionTypes.UpdateChartGroupEditorState,
          payload: Object.assign({}, chartGroupEditor, { v_editing: false })
        });
        /** 重新获取最新数据，从而Detail可以连续编辑且使用到最新的resourceVersion */
        dispatch(
          detailActions.fetchChartGroup({
            name: chartGroupEditor.metadata.name,
            projectID:
              chartGroupEditor.spec.projects && chartGroupEditor.spec.projects.length > 0
                ? chartGroupEditor.spec.projects[0]
                : ''
          })
        );
      }
      /** 结束工作流 */
      dispatch(detailActions.updateChartGroupWorkflow.reset());
    }
  }
});

const restActions = {
  updateChartGroupWorkflow,

  validator: createValidatorActions({
    userDefinedSchema: ChartGroupValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.chartGroupEditor;
    },
    validatorStateLocation: (store: RootState) => {
      return store.chartGroupValidator;
    },
    // used in extraStore, i.t. customFunc: (value, store, extraStore)
    extraValidateStateLocatorPath: ['userInfo']
  }),

  fetchChartGroup: (filter: ChartGroupDetailFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchChartGroup(filter);
      let editor: ChartGroupEditor = response;
      dispatch({
        type: ActionTypes.UpdateChartGroupEditorState,
        payload: editor
      });

      // 拉取关联列表之后，更新关联用户面板WorkflowDialog会用到的CommonUserAssociation状态数据
      let users: UserPlain[] = editor.spec.users
        ? editor.spec.users.map(name => {
            return {
              id: name,
              name: name,
              displayName: name
            };
          })
        : [];
      dispatch({
        type: ActionTypes.UpdateCommonUserAssociation,
        payload: Object.assign({}, getState().commonUserAssociation, {
          users: users,
          originUsers: users,
          addUsers: [],
          removeUsers: []
        })
      });
    };
  },

  /** 更新状态 */
  updateEditorState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { chartGroupEditor } = getState();
      dispatch({
        type: ActionTypes.UpdateChartGroupEditorState,
        payload: Object.assign({}, chartGroupEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateChartGroupEditorState,
      payload: initChartGroupEditorState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(ChartGroupValidateSchema.formKey),
      payload: {}
    };
  }
};
export const detailActions = extend({}, restActions);
