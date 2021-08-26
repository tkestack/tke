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
    extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationHooks, OperationTrigger,
    RecordSet, ReduxAction
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { AlarmPolicy, AlarmPolicyEdition, AlarmPolicyOperator } from '../models/AlarmPolicy';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/**
 * 操作流Action
 */
export const workflowActions = {
  editAlarmPolicy: generateWorkflowActionCreator<AlarmPolicyEdition, AlarmPolicyOperator>({
    actionType: ActionType.CreateAlarmPolicy,
    workflowStateLocator: (state: RootState) => state.alarmPolicyCreateWorkflow,
    operationExecutor: (
      targets: AlarmPolicyEdition[],
      params: AlarmPolicyOperator,
      dispatch: Redux.Dispatch,
      getState: GetState
    ) => {
      return WebAPI.editAlarmPolicy(targets, params, getState().receiverGroup);
    },
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { alarmPolicyCreateWorkflow, route } = getState();
        if (isSuccessWorkflow(alarmPolicyCreateWorkflow)) {
          dispatch(workflowActions.editAlarmPolicy.reset());
          dispatch(alarmPolicyActions.clearAlarmPolicyEdit());
          router.navigate({}, { rid: route.queries['rid'], clusterId: route.queries['clusterId'] });
        }
      }
    }
  }),
  deleteAlarmPolicy: generateWorkflowActionCreator<AlarmPolicy, AlarmPolicyOperator>({
    actionType: ActionType.DeleteAlarmPolicy,
    workflowStateLocator: (state: RootState) => state.alarmPolicyDeleteWorkflow,
    operationExecutor: WebAPI.deleteAlarmPolicy,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { alarmPolicyDeleteWorkflow, route } = getState();
        if (isSuccessWorkflow(alarmPolicyDeleteWorkflow)) {
          dispatch(workflowActions.deleteAlarmPolicy.reset());
          dispatch(alarmPolicyActions.fetch());
          dispatch(alarmPolicyActions.selects([]));
        }
      }
    }
  })
};
