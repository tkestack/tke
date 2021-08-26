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

import { FetchOptions, ReduxAction } from '@tencent/ff-redux';

import { Resource } from '../../common';
import * as ActionType from '../constants/ActionType';
import { EsInfo, RootState } from '../models';

import { Base64 } from 'js-base64';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

export const peEditActions = {
  /** 是否开启持久化存储 */
  isOpenPE: (isOpen: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.IsOpenPE,
      payload: isOpen
    };
  },

  /** 输入es的地址 */
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

  /** 更新页面初始化数据 */
  initPeEditInfoForUpdate: (resource: Resource) => {
    return async (dispatch, getState: GetState) => {
      let storeType = Object.keys(resource.spec.persistentBackEnd)[0];

      let esInfo: EsInfo = resource.spec.persistentBackEnd[storeType];
      dispatch(peEditActions.inputEsAddress((esInfo.scheme ? esInfo.scheme : 'http') + '://' + esInfo.ip + ':' + esInfo.port));
      dispatch(peEditActions.inputIndexName(esInfo.indexName || 'fluentd'));
      dispatch(peEditActions.inputEsUsername(esInfo.user || ''));
      dispatch(peEditActions.inputEsPassword(Base64.decode(esInfo.password) || ''));
    };
  },

  /** 离开设置页面，清除peEdit当中的内容 */
  clearPeEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearPeEdit
    };
  }
};
