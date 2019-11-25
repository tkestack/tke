import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import { CommonAPI } from '../../../modules/common/webapi';
import { getRegionId } from '@tencent/qcloud-nmc';
import { assureRegion, setRegionId } from '../../../../helpers';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { router } from '../router';
import { clusterActions } from './clusterActions';

type GetState = () => RootState;

/** 获取地域列表的action */
const fetchRegionActions = generateFetcherActionCreator({
  actionType: ActionType.FetchRegion,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await CommonAPI.fetchRegionList(getState().regionQuery);
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { regionList, route } = getState();

    if (regionList.data.recordCount) {
      let defaultRegion = route.queries['rid'] || getRegionId();
      defaultRegion = assureRegion(regionList.data.records, defaultRegion, 1);
      dispatch(regionActions.selectRegion(+defaultRegion));
    }
  }
});

/** 查询地域列表Action */
const queryRegionActions = generateQueryActionCreator({
  actionType: ActionType.QueryRegion,
  bindFetcher: fetchRegionActions
});

const restActions = {
  /** 选择地域 */
  selectRegion: (regionId: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { regionList, route } = getState(),
        urlParams = router.resolve(route);

      let regionSelection = regionList.data.records.find(r => r.value === regionId);
      dispatch({
        type: ActionType.SelectRegion,
        payload: regionSelection
      });
      setRegionId(regionId);

      // 进行路由的更新
      router.navigate(urlParams, Object.assign({}, route.queries, { rid: regionId + '' }));

      /**
       * 集群列表的拉取
       */
      dispatch(clusterActions.applyFilter({ regionId }));
    };
  }
};

export const regionActions = extend({}, queryRegionActions, fetchRegionActions, restActions);
