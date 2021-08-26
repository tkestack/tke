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

import { RecordSet, reduceToPayload, ReduxAction } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { initVariable } from '../models';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.CM_Name, ''),

  v_name: reduceToPayload(ActionType.V_CM_Name, initValidator),

  namespace: reduceToPayload(ActionType.CM_Namespace, 'default'),

  variables: (state = [initVariable], action: ReduxAction<any>) => {
    switch (action.type) {
      case ActionType.CM_AddVariable:
      case ActionType.CM_EditVariable:
      case ActionType.CM_DeleteVariable:
      case ActionType.CM_ValidateVariable:
        return action.payload;
      default:
        return state;
    }
  }
});

export const ConfigMapEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建pv的界面
  if (action.type === ActionType.ClearConfigMapEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
