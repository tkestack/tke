import { extend, ReduxAction } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions, FetchState } from '@tencent/qcloud-redux-fetcher';
import { generateWorkflowActionCreator, OperationResult, OperationTrigger } from '@tencent/qcloud-redux-workflow';
import { generateQueryActionCreator, QueryState } from '@tencent/qcloud-redux-query';
import { ComputerFilter } from '../../cluster/models';
import * as ActionTypes from '../constants/ActionTypes';
import * as WebAPI from '../WebAPI';
import { RootState, Strategy, StrategyFilter } from '../models';
import { createListAction } from '@tencent/redux-list';
type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const FFModelStrategyActions = createListAction<Strategy, StrategyFilter>({
  actionName: ActionTypes.StrategyList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchStrategyList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().strategyList;
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
    [OperationTrigger.Done]: dispatch => {
      dispatch(FFModelStrategyActions.applyFilter({}));
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
    [OperationTrigger.Done]: dispatch => {
      dispatch(FFModelStrategyActions.applyFilter({}));
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
  }
});

/**
 * 获取服务
 */
const getCategories = generateFetcherActionCreator({
  actionType: ActionTypes.GetCategories,
  fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
    // const { id, userNames } = options.data;
    let result = await WebAPI.fetchCategoryList();
    return result;
  }
});

const restActions = {
  addStrategy,
  removeStrategy,
  getCategories,
  getStrategy,
  updateStrategy
};

export const strategyActions = extend({}, FFModelStrategyActions, restActions);
