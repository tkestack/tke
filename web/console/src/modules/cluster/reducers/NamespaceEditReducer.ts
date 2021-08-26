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

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.N_Name, ''),

  v_name: reduceToPayload(ActionType.NV_Name, initValidator),

  description: reduceToPayload(ActionType.N_Description, ''),

  v_description: reduceToPayload(ActionType.NV_Description, initValidator)
});

export const NamespaceEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建namespace 页面
  if (action.type === ActionType.ClearNamespaceEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
