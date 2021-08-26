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

import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { Pod, Resource } from '../models';

const TempReducer = combineReducers({
  workloadType: reduceToPayload(ActionType.L_WorkloadType, 'deployment'),

  namespaceSelection: reduceToPayload(ActionType.L_NamespaceSelection, ''),

  workloadSelection: reduceToPayload(ActionType.L_WorkloadSelection, ''),

  workloadQuery: generateQueryReducer({
    actionType: ActionType.L_QueryWorkloadList
  }),

  workloadList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.L_FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  podQuery: generateQueryReducer({
    actionType: ActionType.L_QueryPodList
  }),

  podList: generateFetcherReducer<RecordSet<Pod>>({
    actionType: ActionType.L_FetchPodList,
    initialData: {
      recordCount: 0,
      records: [] as Pod[]
    }
  }),

  podSelection: reduceToPayload(ActionType.L_PodSelection, ''),

  containerSelection: reduceToPayload(ActionType.L_ContainerSelection, ''),

  logQuery: generateQueryReducer({
    actionType: ActionType.L_QueryLogList
  }),

  logList: generateFetcherReducer<RecordSet<string>>({
    actionType: ActionType.L_FetchLogList,
    initialData: {
      recordCount: 0,
      records: [] as string[]
    }
  }),

  tailLines: reduceToPayload(ActionType.L_TailLines, '100'),

  isAutoRenew: reduceToPayload(ActionType.L_IsAutoRenew, false)
});

export const ResourceLogReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResourceLog) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
