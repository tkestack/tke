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
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { initSecretData } from '../constants/initState';
import { Namespace } from '../models';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.Sec_Name, ''),

  v_name: reduceToPayload(ActionType.SecV_Name, initValidator),

  nsList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.Sec_FetchNsList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),

  nsQuery: generateQueryReducer({
    actionType: ActionType.Sec_QueryNsList
  }),

  secretType: reduceToPayload(ActionType.Sec_SecretType, 'Opaque'),

  data: reduceToPayload(ActionType.Sec_UpdateData, [initSecretData]),

  nsType: reduceToPayload(ActionType.Sec_NsType, 'specific'),

  nsListSelection: reduceToPayload(ActionType.Sec_NamespaceSelection, []),

  domain: reduceToPayload(ActionType.Sec_Domain, ''),

  v_domain: reduceToPayload(ActionType.SecV_Domain, initValidator),

  username: reduceToPayload(ActionType.Sec_Username, ''),

  v_username: reduceToPayload(ActionType.SecV_Username, initValidator),

  password: reduceToPayload(ActionType.Sec_Password, ''),

  v_password: reduceToPayload(ActionType.SecV_Password, initValidator)
});

export const SecretEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建 Secret 界面
  if (action.type === ActionType.ClearSecretEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
