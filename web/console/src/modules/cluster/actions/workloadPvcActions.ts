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
import { extend, FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';
import { Resource, ResourceFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const fetchPvcActions = generateFetcherActionCreator({
  actionType: ActionType.W_FetchPvcList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadEdit } = subRoot,
      { pvcQuery } = workloadEdit;

    let pvcResourceInfo = resourceConfig(clusterVersion)['pvc'];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchSpecificResourceList(pvcQuery, pvcResourceInfo, isClearData, true);
    return response;
  },
  finish: (dispatch, getState: GetState) => {}
});

const queryPvcActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.W_QueryPvcList,
  bindFetcher: fetchPvcActions
});

const restActions = {};

export const workloadPvcActions = extend({}, fetchPvcActions, queryPvcActions, restActions);
