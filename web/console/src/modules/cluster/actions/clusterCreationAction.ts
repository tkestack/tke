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

import { deepClone, ReduxAction, uuid } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { initClusterCreationState } from '../constants/initState';
import { RootState } from '../models';

type GetState = () => RootState;

export const clusterCreationAction = {
  /** 更新cluser的名称 */
  updateClusterCreationState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterCreationState } = getState();
      dispatch({
        type: ActionType.UpdateclusterCreationState,
        payload: Object.assign({}, clusterCreationState, obj)
      });
    };
  },

  /** 离开创建页面，清除 Creation当中的内容 */
  clearClusterCreationState: (): ReduxAction<any> => {
    return {
      type: ActionType.UpdateclusterCreationState,
      payload: initClusterCreationState
    };
  }
};
