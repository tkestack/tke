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
import { Event, Resource } from '../models';

let defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, 'default');
/// #if tke
defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, 'default');
/// #endif
/// #if project
defaultNamespace = reduceToPayload(ActionType.E_NamespaceSelection, '');
/// #endif

const TempReducer = combineReducers({
  workloadType: reduceToPayload(ActionType.E_WorkloadType, ''),

  namespaceSelection: defaultNamespace,

  workloadQuery: generateQueryReducer({
    actionType: ActionType.E_QueryWorkloadList
  }),

  workloadList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.E_FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  workloadSelection: reduceToPayload(ActionType.E_WorkloadSelection, ''),

  eventQuery: generateQueryReducer({
    actionType: ActionType.E_QueryEventList
  }),

  eventList: generateFetcherReducer<RecordSet<Event>>({
    actionType: ActionType.E_FetchEventList,
    initialData: {
      recordCount: 0,
      records: [] as Event[]
    }
  }),

  isAutoRenew: reduceToPayload(ActionType.E_IsAutoRenew, true)
});

export const ResourceEventReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResourceEvent) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
