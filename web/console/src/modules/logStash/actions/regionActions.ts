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
import { extend, getRegionId } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { assureRegion, setRegionId } from '../../../../helpers';
import { CommonAPI } from '../../../modules/common/webapi';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
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
