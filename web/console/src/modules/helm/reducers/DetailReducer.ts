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
import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { HelmHistory } from '../models';

const TempReducer = combineReducers({
  helm: reduceToPayload(ActionType.FetchHelm, null),
  isRefresh: reduceToPayload(ActionType.IsRefresh, false),
  historyQuery: generateQueryReducer({
    actionType: ActionType.QueryHistory
  }),
  histories: generateFetcherReducer<RecordSet<HelmHistory>>({
    actionType: ActionType.FetchHistory,
    initialData: {
      recordCount: 0,
      records: [] as HelmHistory[]
    }
  })
});

export const DetailReducer = (inputState, action) => {
  let state = inputState;
  // 销毁详情页面
  if (action.type === ActionType.ClearDetail) {
    state = undefined;
  }
  return TempReducer(state, action);
};
