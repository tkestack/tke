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
import * as WebAPI from '../../cluster/WebAPI';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';

type GetState = () => RootState;

/** fetch workload list */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { workloadQuery, clusterVersion } = getState();
    let { filter } = workloadQuery;
    let workloadTypeMap = {
      Deployment: 'deployment',
      StatefulSet: 'statefulset',
      DaemonSet: 'daemonset',
      TApp: 'tapp'
    };
    let resourceInfo = resourceConfig(clusterVersion)[workloadTypeMap[filter.workloadType]];
    let response = await WebAPI.fetchResourceList(getState().workloadQuery, { resourceInfo });
    return response;
  }
});

/** query Pod list action */
const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

export const workloadActions = extend(fetchWorkloadActions, queryWorkloadActions);
