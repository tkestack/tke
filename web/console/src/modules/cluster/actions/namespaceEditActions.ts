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
import { RootState } from '../models';

type GetState = () => RootState;

export const namespaceEditActions = {
  /** 更新namespace的名称 */
  inputNamespaceName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.N_Name,
      payload: name
    };
  },

  /** 更新namespace的描述 */
  inputNamespaceDesp: (desp: string): ReduxAction<string> => {
    return {
      type: ActionType.N_Description,
      payload: desp
    };
  },

  /** 离开创建页面，清除 namespaceEdit当中的内容 */
  clearNamespaceEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearNamespaceEdit
    };
  }
};
