import { FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { extend } from '@tencent/qcloud-lib';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch subRouter list */
const fetchSubRouterActions = generateFetcherActionCreator({
  actionType: ActionType.FetchSubRouterList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot } = getState(),
      { subRouterQuery } = subRoot;
    let response = await WebAPI.fetchSubRouterList(subRouterQuery);
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
  }
});

/** query subRouter list */
const querySubRouterActions = generateQueryActionCreator({
  actionType: ActionType.QuerySubRouterList,
  bindFetcher: fetchSubRouterActions
});

const restActions = {};

export const subRouterActions = extend(fetchSubRouterActions, querySubRouterActions, restActions);
