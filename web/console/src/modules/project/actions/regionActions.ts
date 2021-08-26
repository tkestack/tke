/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import { createFFListActions, extend } from '@tencent/ff-redux';

import { Region, RegionFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';
import { clusterActions } from './clusterActions';

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
