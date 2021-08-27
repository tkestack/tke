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

import { reduceToPayload, createFFListReducer } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';

const TempReducer = combineReducers({
  ffResourceList: createFFListReducer(FFReduxActionName.Resource_Workload),

  resourceMultipleSelection: reduceToPayload(ActionType.SelectMultipleResource, []),

  resourceDeleteSelection: reduceToPayload(ActionType.SelectDeleteResource, [])
});

export const ResourceReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResource) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
