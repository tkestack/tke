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

import { includes } from '../../common';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';

type GetState = () => RootState;

export const peEditActions = {
  /** 输入 es的地址 */
  inputEsAddress: (address: string): ReduxAction<string> => {
    return {
      type: ActionType.EsAddress,
      payload: address
    };
  },

  /** 输入当前的索引 */
  inputIndexName: (indexName: string): ReduxAction<string> => {
    return {
      type: ActionType.IndexName,
      payload: indexName
    };
  },

  /** 输入当前的索引 */
  inputEsUsername: (username: string): ReduxAction<string> => {
    return {
      type: ActionType.EsUsername,
      payload: username
    };
  },

  /** 输入当前的索引 */
  inputEsPassword: (password: string): ReduxAction<string> => {
    return {
      type: ActionType.EsPassword,
      payload: password
    };
  },
};
