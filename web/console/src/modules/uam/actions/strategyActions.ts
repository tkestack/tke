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
    createFFListActions, extend, FetchOptions, generateFetcherActionCreator,
    generateWorkflowActionCreator, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, Strategy, StrategyFilter } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const FFModelStrategyActions = createFFListActions<Strategy, StrategyFilter>({
  actionName: ActionTypes.StrategyList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchStrategyList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().strategyList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    if (record.data.recordCount) {
      let isNotNeedPoll = record.data.records.filter(item => item.status['phase'] === 'Terminating').length === 0;

      if (isNotNeedPoll) {
        dispatch(FFModelStrategyActions.clearPolling());
      }
    }
  }
});

/**
 * 增加策略
 */
const addStrategy = generateWorkflowActionCreator<Strategy, void>({
  actionType: ActionTypes.AddStrategy,
  workflowStateLocator: (state: RootState) => state.addStrategyWorkflow,
  operationExecutor: WebAPI.addStrategy,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      const { route } = getState();
      let { sub } = router.resolve(route);
      dispatch(strategyActions.poll({ type: sub }));
    }
  }
});

/**
 * 删除策略
 */
const removeStrategy = generateWorkflowActionCreator<any, void>({
  actionType: ActionTypes.RemoveStrategy,
  workflowStateLocator: (state: RootState) => state.removeStrategyWorkflow,
  operationExecutor: WebAPI.removeStrategy,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      const { route } = getState();
      let { sub } = router.resolve(route);
      dispatch(strategyActions.poll({ type: sub }));
    }
  }
});

/**
 * 获取策略
 */
const getStrategy = generateFetcherActionCreator({
  actionType: ActionTypes.GetStrategy,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    // const { id, userNames } = options.data;
    let result = await WebAPI.getStrategy(options.data.id);
    return result;
  }
});

/**
 * 更新策略
 */
const updateStrategy = generateFetcherActionCreator({
  actionType: ActionTypes.UpdateStrategy,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    // const { id, userNames } = options.data;
    let result = await WebAPI.updateStrategy(options.data);
    return result;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    let urlParams = router.resolve(route);
    dispatch(
      strategyActions.getStrategy.fetch({
        noCache: true,
        data: {
          id: route.queries['id']
        }
      })
    );
  }
});

/**
 * 获取服务
 */
const getCategories = generateFetcherActionCreator({
  actionType: ActionTypes.GetCategories,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    let result = await WebAPI.fetchCategoryList();
    return result;
  }
});

// /**
//  * 获取平台策略
//  */
// const getPlatformCategories = generateFetcherActionCreator({
//   actionType: ActionTypes.GetPlatformCategories,
//   fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
//     let result = await WebAPI.getPlatformCategories();
//     return result;
//   }
// });

const restActions = {
  poll: (type) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        strategyActions.polling({
          filter: type,
          delayTime: 5000
        })
      );
    };
  },

  addStrategy,
  removeStrategy,
  getCategories,
  getStrategy,
  updateStrategy,
  // getPlatformCategories
};

export const strategyActions = extend({}, FFModelStrategyActions, restActions);
