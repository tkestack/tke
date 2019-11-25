import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** 获取控制台日志列表Action */
const fetchLogActions = generateFetcherActionCreator({
  actionType: ActionType.FetchLogList,
  fetcher: (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    return WebAPI.fetchLogList(getState().logQuery);
  }
});

/** 查询控制台日志列表Action */
const queryLogActions = generateQueryActionCreator({
  actionType: ActionType.QueryLogList,
  bindFetcher: fetchLogActions
});

export const logActions = extend({}, queryLogActions, fetchLogActions);
