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
import { resourceConfig } from '@config';
import { CommonAPI, ResourceFilter, ResourceInfo } from '@src/modules/common';
import {
  createFFListActions,
  extend,
  FetchOptions,
  generateFetcherActionCreator,
  generateWorkflowActionCreator,
  isSuccessWorkflow,
  OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, Strategy } from '../models';
import { User, UserFilter } from '../models/index';
import * as WebAPI from '../WebAPI';
import { router } from '../router';

type GetState = () => RootState;

/**
 * 增加用户
 */
const addUser = generateWorkflowActionCreator<User, void>({
  actionType: ActionTypes.AddUser,
  workflowStateLocator: (state: RootState) => state.addUserWorkflow,
  operationExecutor: WebAPI.addUser,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      let { addUserWorkflow, route } = getState();
      if (isSuccessWorkflow(addUserWorkflow)) {
        router.navigate({ module: 'user' }, route.queries);
      }
      /** 结束工作流 */
      // dispatch(userActions.poll());
      // dispatch(FFModelUserActions.applyFilter({}));
      let count = 0;
      const timer = setInterval(() => {
        dispatch(FFModelUserActions.applyFilter({}));
        // dispatch(userActions.poll());
        ++count;
        if (count >= 3) {
          clearInterval(timer);
        }
      }, 1500);
    }
  }
});

/**
 * 删除用户
 */
const removeUser = generateWorkflowActionCreator<any, void>({
  actionType: ActionTypes.RemoveUser,
  workflowStateLocator: (state: RootState) => state.removeUserWorkflow,
  operationExecutor: WebAPI.removeUser,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch) => {
      dispatch(userActions.poll());
    }
  }
});

/**
 * 获取用户
 */
const getUser = generateFetcherActionCreator({
  actionType: ActionTypes.GetUser,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    let result = await WebAPI.getUser(options.data.name);
    return result;
  }
});

/**
 * 更新用户
 */
const updateUser = generateFetcherActionCreator({
  actionType: ActionTypes.UpdateUser,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    let result = await WebAPI.updateUser(options.data.user);
    return result;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let count = 0;
    const timer = setInterval(() => {
      dispatch(FFModelUserActions.applyFilter({}));
      // dispatch(userActions.poll());
      ++count;
      if (count >= 3) {
        clearInterval(timer);
      }
    }, 1500);
  }
});

/**
 * 用户列表操作
 */
const FFModelUserActions = createFFListActions<User, UserFilter>({
  actionName: ActionTypes.UserList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchUserList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    if (record.data.recordCount) {
      let isNotNeedPoll =
        record.data.records.filter(item => item.status && item.status['phase'] && item.status['phase'] === 'Deleting')
          .length === 0;
      if (isNotNeedPoll) {
        dispatch(FFModelUserActions.clearPolling());
      }
    }
  }
});

/* ================================ start 权限列表相关的 ================================ */
const StrategyListActions = createFFListActions<Strategy, ResourceFilter>({
  actionName: ActionTypes.UserStrategyList,
  fetcher: async (query, getState: GetState) => {
    let resourceInfo: ResourceInfo = resourceConfig()['localidentity'];

    let response = await CommonAPI.fetchExtraResourceList({
      query,
      resourceInfo,
      extraResource: 'policies'
    });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userStrategyList;
  }
});

const strategyRestActions = {};

const strategyActions = extend({}, StrategyListActions, strategyRestActions);
/* ================================ end 权限列表相关的 ================================ */

const restActions = {
  addUser,
  removeUser,
  getUser,
  updateUser,
  strategy: strategyActions,

  /** 轮训操作 */
  poll: () => {
    return async (dispatch: Redux.Dispatch) => {
      dispatch(
        userActions.polling({
          filter: {
            ifAll: true
          },
          delayTime: 5000
        })
      );
    };
  },

  /** 初始化集群的版本 */
  getUsersByName: (username: string) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let filterUsers = [];
      getState().userList.list.data.records.forEach(user => {
        if (user.spec.username === username || user.spec.displayName === username) {
          filterUsers.push(user);
        }
      });
      dispatch({
        type: ActionTypes.FetchUserByName,
        payload: filterUsers
      });
    };
  }
};

export const userActions = extend({}, FFModelUserActions, restActions);
