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
// import { extend } from '@tencent/ff-redux';
// import { router } from '../router';
// import { getRegionId, setRegionId, assureRegion } from '../../../../helpers';
// import { RootState, RegionFilter, Region } from '../models';
// import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
// import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
// import * as ActionType from '../constants/ActionType';
// import * as WebAPI from '../WebAPI';
// import { clusterActions } from './clusterActions';
// import { includes } from '../../common/utils/includes';

// type GetState = () => RootState;
// const fetchOptions: FetchOptions = {
//   noCache: false
// };

// /**
//  * 获取地域列表Action
//  */
// const fetchRegionActions = generateFetcherActionCreator({
//   actionType: ActionType.FetchRegion,
//   fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
//     let { regionQuery } = getState();
//     return WebAPI.fetchRegionList(getState().regionQuery);
//   },
//   finish: (dispatch, getState: GetState) => {
//     let { route, regionList } = getState(),
//       urlParams = router.resolve(route),
//       regionId = route.queries['rid'];
//     let rid = regionId || getRegionId() || regionList.data.records[0].value;
//     dispatch(regionActions.select(+rid));
//   }
// });

// /**
//  * 查询Region列表Action
//  */
// const queryRegionActions = generateQueryActionCreator<RegionFilter>({
//   actionType: ActionType.QueryRegion,
//   bindFetcher: fetchRegionActions
// });

// /**
//  * 其他地域Action
//  */
// const restActions = {
//   select: (selectRegion: number) => {
//     return (dispatch, getState: GetState) => {
//       const { regionList, regionSelection, route } = getState(),
//         urlParams = router.resolve(route);
//       let regionId = assureRegion(regionList.data.records, selectRegion, regionSelection.value);
//       let regionInfo = regionList.data.records.find(r => r.value === regionId);
//       setRegionId(regionId);
//       dispatch({
//         type: ActionType.SelectRegion,
//         payload: regionInfo
//       });
//       router.navigate(urlParams, Object.assign({}, route.queries, { rid: regionId }));
//       dispatch(clusterActions.applyFilter({ regionId: +regionId }));
//     };
//   }
// };

// export const regionActions = extend({}, queryRegionActions, fetchRegionActions, restActions);
