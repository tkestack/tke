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
import { extend, generateWorkflowActionCreator, OperationResult, OperationTrigger, createFFListActions } from '@tencent/ff-redux';
import { RootState, Group, GroupFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchGroupActions = createFFListActions<Group, GroupFilter>({
  actionName: ActionTypes.GroupList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchGroupList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().groupList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    let isNotNeedPoll = true;
    if (record.data.recordCount) {
      isNotNeedPoll = record.data.records.filter(item => item.status && item.status['phase'] && item.status['phase'] === 'Terminating').length === 0;
    }
    if (isNotNeedPoll) {
      dispatch(fetchGroupActions.clearPolling());
    }
  }
});

/**
 * 删除用户组
 */
const removeGroupWorkflow = generateWorkflowActionCreator<Group, void>({
  actionType: ActionTypes.RemoveGroup,
  workflowStateLocator: (state: RootState) => state.groupRemoveWorkflow,
  operationExecutor: WebAPI.deleteGroup,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      //在列表页删除的动作，因此直接重新拉取一次数据
      dispatch(listActions.poll());
      /** 结束工作流 */
      dispatch(listActions.removeGroupWorkflow.reset());
    }
  }
});

const restActions = {
  removeGroupWorkflow,

  /** 轮询操作 */
  poll: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        listActions.polling({
            delayTime: 3000
          }
        )
      );
    };
  },

};

export const listActions = extend({}, fetchGroupActions, restActions);
