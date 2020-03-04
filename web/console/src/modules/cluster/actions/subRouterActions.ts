import { extend, FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import * as WebAPI from '../WebAPI';
import { clusterActions } from './clusterActions';

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
    let { clusterId, rid } = route.queries;
    // 获取当前集群所有开启的Addon

    dispatch(clusterActions.fetchClusterAddon(clusterId, +rid));
  }
});

/** query subRouter list */
const querySubRouterActions = generateQueryActionCreator({
  actionType: ActionType.QuerySubRouterList,
  bindFetcher: fetchSubRouterActions
});

const restActions = {};

export const subRouterActions = extend(fetchSubRouterActions, querySubRouterActions, restActions);
