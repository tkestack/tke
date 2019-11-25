import { FetchOptions, generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { extend, ReduxAction } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models/RootState';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { ResourceInfo, Resource, ResourceFilter } from '../../common/models';
import { router } from '../router';
import { peEditActions } from './peEditActions';
import { PollEventName, isNeedPollPE } from '../constants/Config';
import { includes } from '../../common/utils';
import { CommonAPI } from '../../common/webapi';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** 获取集群的列表 */
const fetchPEActions = generateFetcherActionCreator({
  actionType: ActionType.FetchPeList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { peQuery, resourceInfo } = getState();

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await CommonAPI.fetchResourceList({ query: peQuery, resourceInfo, isClearData });
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { peList, route } = getState(),
      urlParams = router.resolve(route);

    if (peList.data.recordCount) {
      if (!urlParams['mode']) {
        dispatch(peActions.selectPe([peList.data.records[0]]));
        // 这里需要去判断是否需要轮询
        let isNeedPoll = false;
        peList.data.records.forEach(item => {
          if (includes(isNeedPollPE, (item.status.phase as string).toLowerCase())) {
            isNeedPoll = true;
          }
        });
        !isNeedPoll && dispatch(peActions.clearPollEvent());
      } else if (urlParams['mode'] === 'update') {
        let peInfo = peList.data.records.find(item => item.spec.clusterName === route.queries['clusterId']);
        dispatch(peActions.selectPe([peInfo]));
        // 初始化该集群的persistentEvent的数据
        dispatch(peEditActions.initPeEditInfoForUpdate(peInfo));

        // 如果在更新界面，直接清除原来的轮询
        dispatch(peActions.clearPollEvent());
      } else if (urlParams['mode'] === 'create') {
        dispatch(peActions.clearPollEvent());
      }
    } else {
      dispatch(peActions.clearPollEvent());
    }
  }
});

/** 查询集群列表的action */
const qeuryPEActions = generateQueryActionCreator({
  actionType: ActionType.QueryPeList,
  bindFetcher: fetchPEActions
});

const restActions = {
  /** 初始化peResourceInfo的相关信息 */
  initPeResourceInfo: (resourceInfo: ResourceInfo): ReduxAction<ResourceInfo> => {
    return {
      type: ActionType.InitResourceInfo,
      payload: resourceInfo
    };
  },

  /** 选择PersistenEvent */
  selectPe: (resource: Resource[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectPe,
        payload: resource
      });
    };
  },

  /** 轮询拉取条件 */
  poll: (queryObj: ResourceFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      // 每次轮询之前先清空之前的轮询
      dispatch(peActions.clearPollEvent());
      // 触发列表的查询
      dispatch(peActions.applyFilter(queryObj));

      window[PollEventName['peList']] = setInterval(() => {
        dispatch(peActions.poll(queryObj));
      }, 10000);
    };
  },

  /** 清空轮询条件 */
  clearPollEvent: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['peList']]);
    };
  }
};

export const peActions = extend({}, fetchPEActions, qeuryPEActions, restActions);
