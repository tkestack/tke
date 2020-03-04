import { extend } from '@tencent/qcloud-lib';
import { RootState, Region, RegionFilter } from '../models';
import * as WebAPI from '../WebAPI';
import { clusterActions } from './clusterActions';
import { createFFListActions } from '@tencent/ff-redux';

type GetState = () => RootState;

const FFModelRegionActions = createFFListActions<Region, RegionFilter>({
  actionName: 'region',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchRegionList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().region;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    if (record.data.recordCount) {
      dispatch(regionActions.selectRegion(1));
    }
  }
});

const restActions = {
  selectRegion: (regionId: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        FFModelRegionActions.select(getState().region.list.data.recordCount && getState().region.list.data.records[0])
      );

      dispatch(clusterActions.applyFilter({ regionId: 1 }));
    };
  }
};

export const regionActions = extend({}, FFModelRegionActions, restActions);
