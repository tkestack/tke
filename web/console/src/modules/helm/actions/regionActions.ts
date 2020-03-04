import { Region, RegionFilter } from '../../common/models';
import { RootState } from '../models';
import * as WebAPI from '../WebAPI';
import { assureRegion } from '../../../../helpers';
import { getRegionId } from '../../../../helpers/appUtil';
import { router } from '../router';
import { clusterActions } from './clusterActions';
import { createFFListActions } from '@tencent/ff-redux';
import { FFReduxActionName } from '../constants/Config';
import { CommonAPI } from '../../common/webapi';

type GetState = () => RootState;

/** 地域列表的Actions */
export const regionActions = createFFListActions<Region, RegionFilter>({
  actionName: FFReduxActionName.REGION,
  fetcher: async query => {
    let response = await CommonAPI.fetchRegionList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().listState.region;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    if (record.data.recordCount) {
      let defaultRegion = route.queries['rid'] || getRegionId();
      defaultRegion = assureRegion(record.data.records, defaultRegion, 1);

      let r = record.data.records.find(r => r.value + '' === defaultRegion + '');
      dispatch(regionActions.select(r));
    }
  },
  onSelect: (region: Region, dispatch: Redux.Dispatch, getState: GetState) => {
    let {
        listState: { cluster },
        route
      } = getState(),
      urlParams = router.resolve(route);
    router.navigate(urlParams, Object.assign({}, route.queries, { rid: region ? region.value : 1 }));
    /// #if tke
    dispatch(
      clusterActions.applyFilter(Object.assign({}, cluster.query.filter, { regionId: region ? region.value : 1 }))
    );
    /// #endif
  }
});
