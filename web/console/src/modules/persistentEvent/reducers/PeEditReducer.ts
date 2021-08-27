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

import { initValidator } from '../../common/models/Validation';
import * as ActionType from '../constants/ActionType';

const TempReducer = combineReducers({
  isOpen: reduceToPayload(ActionType.IsOpenPE, true),

  esAddress: reduceToPayload(ActionType.EsAddress, ''),

  v_esAddress: reduceToPayload(ActionType.V_EsAddress, initValidator),

  indexName: reduceToPayload(ActionType.IndexName, ''),

  v_indexName: reduceToPayload(ActionType.V_IndexName, initValidator),

  esUsername: reduceToPayload(ActionType.EsUsername, ''),

  esPassword: reduceToPayload(ActionType.EsPassword, '')
});

export const PeEditReducer = (state, action) => {
  let newState = state;
  // 销毁设置事件持久化的界面
  if (action.type === ActionType.ClearPeEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
