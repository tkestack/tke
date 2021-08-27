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
import { ReduxAction } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { peEditActions } from './peEditActions';

export const addonEditActions = {
  pe: peEditActions,

  /** 需要开通的扩展组件的名称 */
  selectAddonName: (name: string) => {
    return async (dispatch: Redux.Dispatch) => {
      dispatch({
        type: ActionType.AddonName,
        payload: name
      });
    };
  },

  /** 清除开通addon的相关信息 */
  clearCreateAddon: (): ReduxAction<void> => {
    return { type: ActionType.ClearAddonEdit };
  }
};
