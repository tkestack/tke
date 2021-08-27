/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { createFFListActions } from '@tencent/ff-redux';

import { assureRegion } from '../../../../helpers';
import { getRegionId } from '../../../../helpers/appUtil';
import { Region, RegionFilter } from '../../common/models';
import { CommonAPI } from '../../common/webapi';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { clusterActions } from './clusterActions';

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
