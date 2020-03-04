import { extend } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { CommonAPI } from '../../common';
import { ResourceInfo } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { PodListFilter, RootState } from '../models';

type GetState = () => RootState;

/** 获取PodList */
const fetchPodList = generateFetcherActionCreator({
  actionType: ActionType.FetchPodList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { logStashEdit, clusterVersion } = getState(),
      { podListQuery } = logStashEdit;
    let workloadType = logStashEdit.containerFileWorkloadType;
    let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    let response = await CommonAPI.fetchExtraResourceList({
      query: podListQuery,
      resourceInfo,
      isClearData,
      extraResource: 'pods'
    });
    return response;
  }
});

/** 获取PodList的查询 */
const queryPodList = generateQueryActionCreator<PodListFilter>({
  actionType: ActionType.QueryPodList,
  bindFetcher: fetchPodList
});

const restActions = {};

export const podActions = extend({}, fetchPodList, queryPodList, restActions);
