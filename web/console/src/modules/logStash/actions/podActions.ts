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

import { extend } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { CommonAPI } from '../../common';
import { ResourceInfo } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { PodListFilter, RootState } from '../models';

type GetState = () => RootState;

/** 获取PodList */
const fetchPodList = generateFetcherActionCreator({
  actionType: ActionType.FetchPodList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { logStashEdit, clusterVersion } = getState(),
      { podListQuery } = logStashEdit;
    let workloadType = logStashEdit.containerFileWorkloadType;
    let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    let response = await CommonAPI.fetchExtraResourceList({
      query: podListQuery,
      resourceInfo,
      isClearData,
      extraResource: 'pods'
    });
    return response;
  }
});

/** 获取PodList的查询 */
const queryPodList = generateQueryActionCreator<PodListFilter>({
  actionType: ActionType.QueryPodList,
  bindFetcher: fetchPodList
});

const restActions = {};

export const podActions = extend({}, fetchPodList, queryPodList, restActions);
