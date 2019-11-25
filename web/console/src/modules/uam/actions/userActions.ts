import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { generateWorkflowActionCreator, OperationResult, OperationTrigger } from '@tencent/qcloud-redux-workflow';
import { RootState, Strategy, StrategyFilter } from '../models';
import * as ActionTypes from '../constants/ActionTypes';
import * as WebAPI from '../WebAPI';
import { User, UserFilter } from '../models/index';
import { createListAction } from '@tencent/redux-list';
import { CommonAPI, ResourceFilter, ResourceInfo } from '@src/modules/common';
import { resourceConfig } from '@config';
import { AssociatedUserParams } from './associatedActions';
type GetState = () => RootState;

/**
 * 增加用户
 */
const addUser = generateWorkflowActionCreator<User, void>({
  actionType: ActionTypes.AddUser,
  workflowStateLocator: (state: RootState) => state.addUserWorkflow,
  operationExecutor: WebAPI.addUser,
  after: {
    [OperationTrigger.Done]: dispatch => {
      dispatch(FFModelUserActions.applyFilter({}));
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
    [OperationTrigger.Done]: dispatch => {
      dispatch(FFModelUserActions.applyFilter({}));
    }
  }
});

/**
 * 获取用户
 */
const getUser = generateFetcherActionCreator({
  actionType: ActionTypes.GetUser,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    // const { id, userNames } = options.data;
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
    // const { id, userNames } = options.data;
    let result = await WebAPI.updateUser(options.data.user);
    return result;
  }
});

/**
 * 用户列表操作
 */
const FFModelUserActions = createListAction<User, UserFilter>({
  actionName: ActionTypes.UserList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchUserList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userList;
  }
});

/* ================================ start 权限列表相关的 ================================ */
const StrategyListActions = createListAction<Strategy, ResourceFilter>({
  actionName: ActionTypes.UserStrategyList,
  fetcher: async (query, getState: GetState) => {
    let resourceInfo: ResourceInfo = resourceConfig()['localidentity'];

    let response = await CommonAPI.fetchExtraResourceList<Strategy>({
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

  /** 初始化集群的版本 */
  getUsersByName: (username: string) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let filterUsers = [];
      getState().userList.list.data.records.forEach(user => {
        if (user.name === username || user.Spec.extra.displayName === username) {
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
