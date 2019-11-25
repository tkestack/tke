import { extend } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { clusterActions } from './clusterActions';
import { assureRegion } from '../../../../helpers';
import { setRegionId, getRegionId } from '../../../../helpers/appUtil';
import { router } from '../router';
import { CommonAPI } from '../../common/webapi';
import { Region, RegionFilter } from '../../common';
import { FFReduxActionName } from '../constants/Config';
import { createListAction } from '@tencent/redux-list';

type GetState = () => RootState;

/** 地域列表的Actions */
const FFModelRegionActions = createListAction<Region, RegionFilter>({
  actionName: FFReduxActionName.REGION,
  fetcher: async query => {
    let response = await CommonAPI.fetchRegionList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().region;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { region, route } = getState();
    if (region.list.data.recordCount) {
      let defaulRegion = route.queries['rid'] || getRegionId();
      defaulRegion = assureRegion(region.list.data.records, defaulRegion, 1);
      dispatch(regionActions.selectRegion(+defaulRegion, false));
    }
  }
});

const restActions = {
  selectRegion: (regionId: number, isNeedFetchZoneList: boolean = true) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { cluster, region, route } = getState(),
        urlParams = router.resolve(route);

      let finalRegion = assureRegion(region.list.data.records, regionId, 1);

      let regionInfo = region.list.data.records.find(r => r.value === regionId);
      dispatch(FFModelRegionActions.select(regionInfo));
      setRegionId(finalRegion);
      router.navigate(urlParams, Object.assign({}, route.queries, { rid: finalRegion }));

      dispatch(clusterActions.selectCluster([]));
      dispatch(clusterActions.poll());
    };
  }
};

export const regionActions = extend({}, FFModelRegionActions, restActions);
