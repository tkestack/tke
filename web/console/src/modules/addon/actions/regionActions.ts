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
import { createFFListActions, extend } from '@tencent/ff-redux';

import { assureRegion, getRegionId, setRegionId } from '../../../../helpers';
import { Region, RegionFilter } from '../../common';
import { CommonAPI } from '../../common/webapi';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import { clusterActions } from './clusterActions';

type GetState = () => RootState;

/** 地域列表的Actions */
const ListRegionActions = createFFListActions<Region, RegionFilter>({
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
      let defaultRegion = route.queries['rid'] || getRegionId();
      defaultRegion = assureRegion(region.list.data.records, defaultRegion, 1);
      dispatch(regionActions.selectRegion(+defaultRegion));
    }
  }
});

const restActions = {
  selectRegion: (regionId: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { region, route, cluster } = getState(),
        urlParams = router.resolve(route);

      let regionInfo = region.list.data.records.find(r => r.value === regionId);
      dispatch(ListRegionActions.select(regionInfo));
      setRegionId(regionId);
      router.navigate(urlParams, Object.assign({}, route.queries, { rid: regionId }));

      // 进行集群列表的获取
      dispatch(clusterActions.applyFilter(Object.assign({}, cluster.query.filter, { regionId })));
    };
  }
};

export const regionActions = extend({}, ListRegionActions, restActions);
