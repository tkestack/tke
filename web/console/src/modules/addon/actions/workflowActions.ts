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
    generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';

import { CreateResource } from '../../common';
import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { clusterActions } from './clusterActions';

type GetState = () => RootState;

export const workflowActions = {
  /** 创建、编辑resourceIns */
  modifyResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyResource,
    workflowStateLocator: (state: RootState) => state.modifyResourceFlow,
    operationExecutor: CommonAPI.modifyResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { modifyResourceFlow, route } = getState();
        if (isSuccessWorkflow(modifyResourceFlow)) {
          // 跳转到列表页面
          router.navigate({}, route.queries);
          // reset modifyResourcflow
          dispatch(workflowActions.modifyResource.reset());
        }
      }
    }
  }),

  /** 创建多种资源的flow */
  applyResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ApplyResource,
    workflowStateLocator: (state: RootState) => state.applyResourceFlow,
    operationExecutor: CommonAPI.applyResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { applyResourceFlow, route } = getState();

        if (isSuccessWorkflow(applyResourceFlow)) {
          // 跳转到列表页
          router.navigate({}, route.queries);
        }
      }
    }
  }),

  /** 删除resourceIns */
  deleteResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeleteResource,
    workflowStateLocator: (state: RootState) => state.deleteResourceFlow,
    operationExecutor: CommonAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { route, deleteResourceFlow } = getState();

        if (isSuccessWorkflow(deleteResourceFlow)) {
          dispatch(workflowActions.deleteResource.reset());
          dispatch(clusterActions.addon.poll());
        }
      }
    }
  })
};
