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

import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { CreateResource, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/**
 * 操作流Action
 */
export const workflowActions = {
  /** 设置持久化事件 */
  modifyPeFlow: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyPeFlow,
    workflowStateLocator: (state: RootState) => state.modifyPeFlow,
    operationExecutor: WebAPI.modifyPeConfig,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { route, modifyPeFlow } = getState(),
          urlParams = router.resolve(route);

        if (isSuccessWorkflow(modifyPeFlow)) {
          // 跳转到 持久化事件的列表页
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { clusterId: undefined })));
          let newUrlParams = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: undefined })));
          router.navigate(newUrlParams, routeQueries);
          // dispatch(clusterActions.applyFilter({ regionId: route.queries['rid'] }));
          // dispatch(peActions.applyFilter({ regionId: route.queries['rid'] }));
          dispatch(workflowActions.modifyPeFlow.reset());
        }
      }
    }
  }),

  /** 删除PersistentEvent */
  deletePeFlow: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeletePeFlow,
    workflowStateLocator: (state: RootState) => state.deletePeFlow,
    operationExecutor: CommonAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { deletePeFlow, route } = getState(),
          urlParams = router.resolve(route);

        if (isSuccessWorkflow(deletePeFlow)) {
          // 跳转到 持久化事件的列表页
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { clusterId: undefined })));
          let newUrlParams = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: undefined })));
          router.navigate(newUrlParams, routeQueries);
          dispatch(workflowActions.deletePeFlow.reset());
        }
      }
    }
  })
};
