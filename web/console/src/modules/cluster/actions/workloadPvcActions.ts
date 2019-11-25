import { extend } from '@tencent/qcloud-lib';
import { FetchOptions, generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState, ResourceFilter, Resource } from '../models';
import * as WebAPI from '../WebAPI';
import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const fetchPvcActions = generateFetcherActionCreator({
  actionType: ActionType.W_FetchPvcList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadEdit } = subRoot,
      { pvcQuery } = workloadEdit;

    let pvcResourceInfo = resourceConfig(clusterVersion)['pvc'];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchSpecificResourceList(pvcQuery, pvcResourceInfo, isClearData, true);
    return response;
  },
  finish: (dispatch, getState: GetState) => {}
});

const queryPvcActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.W_QueryPvcList,
  bindFetcher: fetchPvcActions
});

const restActions = {};

export const workloadPvcActions = extend({}, fetchPvcActions, queryPvcActions, restActions);
