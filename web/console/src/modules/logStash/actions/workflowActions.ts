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

import { CreateResource } from 'src/modules/cluster/models';

import {
    extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';

import { CommonAPI } from '../../../../src/modules/common';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { logActions } from './logActions';
import { logDaemonsetActions } from './logDaemonsetActions';
import { namespaceActions } from './namespaceActions';

type GetState = () => RootState;

/** 操作流actions */
export const workflowActions = {
  /**
   * 单行删除日志采集规则
   */
  inlineDeleteLog: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.InlineDeleteLog,
    workflowStateLocator: (state: RootState) => state.inlineDeleteLog,
    operationExecutor: CommonAPI.deleteResourceIns,
    after: extend(
      {},
      {
        [OperationTrigger.Done]: (dispatch, getState) => {
          let deleteLog = getState().inlineDeleteLog;
          let { route, namespaceSelection } = getState();
          if (isSuccessWorkflow(deleteLog)) {
            dispatch(
              logActions.applyFilter({
                clusterId: route.queries['clusterId'],
                isClear: false,
                namespace: namespaceSelection
              })
            );
            dispatch(workflowActions.inlineDeleteLog.reset());
          }
        }
      }
    )
  }),
  /** 开通日志采集规则的工作流 */
  authorizeOpenLog: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.AuthorizeOpenLog,
    workflowStateLocator: (state: RootState) => state.authorizeOpenLogFlow,
    operationExecutor: CommonAPI.modifyResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { authorizeOpenLogFlow, route } = getState(),
          urlParams = router.resolve(route);
        let mode = urlParams['mode'];
        if (isSuccessWorkflow(authorizeOpenLogFlow)) {
          dispatch(workflowActions.authorizeOpenLog.reset());
          dispatch(logDaemonsetActions.fetch());
          // 进行路由的跳转，如果没有开通的，并且在列表页，则默认跳转到创建日志采集规则的页面
          if (!mode) {
            router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
          }
        }
      }
    }
  }),

  /** 创建、修改日志采集规则 */
  modifyLogStash: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyLogStashFlow,
    workflowStateLocator: (state: RootState) => state.modifyLogStashFlow,
    operationExecutor: WebAPI.modifyLogStash,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { modifyLogStashFlow, route } = getState();

        if (isSuccessWorkflow(modifyLogStashFlow)) {
          // 初始化flow
          dispatch(workflowActions.modifyLogStash.reset());

          // 进行路由的跳转，回列表页
          let newRouteQueies = JSON.parse(
            JSON.stringify(Object.assign({}, route.queries, { stashName: undefined, namespace: undefined }))
          );
          router.navigate({}, newRouteQueies);
        }
      }
    }
  })
};
