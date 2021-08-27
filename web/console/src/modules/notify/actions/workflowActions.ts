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
import {
    generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { Resource, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { resourceActions } from './resourceActions';

type GetState = () => RootState;

/**
 * 操作流Action
 */
export const workflowActions = {
  deleteResource: generateWorkflowActionCreator<Resource, {}>({
    actionType: ActionType.DeleteResource,
    workflowStateLocator: (state: RootState) => state.resourceDeleteWorkflow,
    operationExecutor: WebAPI.deleteResource,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { resourceDeleteWorkflow, route } = getState();
        let urlParams = router.resolve(route);
        if (isSuccessWorkflow(resourceDeleteWorkflow)) {
          dispatch(workflowActions.deleteResource.reset());
          dispatch(resourceActions[urlParams.resourceName || 'channel'].fetch());
          dispatch(resourceActions[urlParams.resourceName || 'channel'].selects([]));
        }
      }
    }
  }),

  modifyResource: generateWorkflowActionCreator<Resource, {}>({
    actionType: ActionType.ModifyResource,
    workflowStateLocator: (state: RootState) => state.modifyResourceFlow,
    operationExecutor: WebAPI.modifyResource,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        // 清空更改名称的表单
        let { modifyResourceFlow, route } = getState();
        let urlParams = router.resolve(route);
        if (isSuccessWorkflow(modifyResourceFlow)) {
          dispatch(workflowActions.modifyResource.reset());
          dispatch(resourceActions[urlParams.resourceName || 'channel'].fetch());
          dispatch(resourceActions[urlParams.resourceName || 'channel'].selects([]));
          router.navigate({ ...urlParams, mode: urlParams.mode === 'create' ? 'list' : 'detail' }, route.queries);
        }
      }
    }
  })
};
