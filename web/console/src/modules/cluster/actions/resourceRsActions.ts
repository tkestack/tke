import { extend, FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import * as ActionType from '../constants/ActionType';
import { Replicaset, ResourceFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** 获取rs列表的action */
const fetchRsActions = generateFetcherActionCreator({
  actionType: ActionType.FetchRsList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { resourceDetailState, resourceOption } = subRoot,
      { ffResourceList } = resourceOption,
      { rsQuery } = resourceDetailState;

    let labelInfo = ffResourceList.selection.metadata.labels;
    let labelKeys = Object.keys(labelInfo);

    let isClearData = (fetchOptions && fetchOptions.noCache) || labelKeys.length === 0 ? true : false;
    let replicasetResourceInfo = resourceConfig(clusterVersion)['rs'];
    let k8sQuery = labelKeys.length
      ? {
          labelSelector: `${labelKeys[0]}=${labelInfo[labelKeys[0]]}`
        }
      : {};

    let response = await WebAPI.fetchSpecificResourceList(rsQuery, replicasetResourceInfo, isClearData, true, k8sQuery);
    // 这里主要是根据时间进行排序，时间最新的，排在最前面，即不可修改
    response.records.sort((pre: Replicaset, next: Replicaset) => {
      // return new Date(pre.metadata.creationTimestamp).getTime() - new Date(next.metadata.creationTimestamp).getTime() <
      //   0
      //   ? 1
      //   : -1;
      return +pre.metadata.annotations['deployment.kubernetes.io/revision'] <
        +next.metadata.annotations['deployment.kubernetes.io/revision']
        ? 1
        : -1;
    });
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { rsList } = getState().subRoot.resourceDetailState;

    rsList.data.recordCount && dispatch(resourceRsActions.selectRs([rsList.data.records[0]]));
  }
});

/** 查询rs列表的Action */
const queryRsActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.QueryRsList,
  bindFetcher: fetchRsActions
});

const restActions = {
  /** 选择rs */
  selectRs: (rs: Replicaset[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.RsSelection,
        payload: rs
      });
    };
  }
};

export const resourceRsActions = extend(fetchRsActions, queryRsActions, restActions);
