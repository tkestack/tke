import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import { getRegionId } from '@tencent/qcloud-nmc';
import { assureRegion, setRegionId } from '../../../../helpers';
import { router } from '../router';
import { clusterActions } from './clusterActions';
import { CommonAPI } from '../../common/webapi';
import { createListAction } from '@tencent/redux-list';
import { Region, RegionFilter } from '../../common';
import { FFReduxActionName } from '../constants/Config';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** 地域列表的Actions */
const ListRegionActions = createListAction<Region, RegionFilter>({
  actionName: FFReduxActionName.REGION,
  fetcher: async query => {
    let response = await CommonAPI.fetchRegionList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().region;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    if (record.data.recordCount) {
      let defaulRegion = route.queries['rid'] || getRegionId();
      defaulRegion = assureRegion(record.data.records, defaulRegion, 1);
      dispatch(regionActions.selectRegion(+defaulRegion));
    }
  }
});

const restActions = {
  selectRegion: (regionId: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { region, route } = getState(),
        urlParams = router.resolve(route);

      let regionInfo = region.list.data.records.find(r => r.value === regionId);
      dispatch(ListRegionActions.select(regionInfo));
      setRegionId(regionId);
      router.navigate(urlParams, Object.assign({}, route.queries, { rid: regionId }));

      // 拉取相对的集群列表
      dispatch(clusterActions.applyFilter({ regionId }));
    };
  }
};

export const regionActions = extend({}, ListRegionActions, restActions);
