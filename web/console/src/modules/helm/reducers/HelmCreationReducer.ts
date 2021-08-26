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

import { createFFListReducer, RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName, HelmResource, OtherType, TencentHubType } from '../constants/Config';
import { TencenthubChart, TencenthubChartVersion, TencenthubNamespace } from '../models';

const TempReducer = combineReducers({
  region: createFFListReducer(FFReduxActionName.REGION, 'HelmCreate'),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER, 'HelmCreate'),

  name: reduceToPayload(ActionType.C_CreateionName, ''),

  isValid: reduceToPayload(ActionType.IsValid, {
    name: '',
    otherChartUrl: '',
    otherUserName: '',
    otherPassword: ''
  }),

  resourceSelection: reduceToPayload(ActionType.ResourceSelection, HelmResource.Other),

  token: reduceToPayload(ActionType.TencenthubToken, ''),

  tencenthubTypeSelection: reduceToPayload(ActionType.TencenthubTypeSelection, TencentHubType.Public),
  tencenthubNamespaceList: generateFetcherReducer<RecordSet<TencenthubNamespace>>({
    actionType: ActionType.FetchTencenthubNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as TencenthubNamespace[]
    }
  }),
  tencenthubNamespaceSelection: reduceToPayload(ActionType.TencenthubNamespaceSelection, ''),
  tencenthubChartList: generateFetcherReducer<RecordSet<TencenthubChart>>({
    actionType: ActionType.FetchTencenthubChartList,
    initialData: {
      recordCount: 0,
      records: [] as TencenthubChart[]
    }
  }),
  tencenthubChartSelection: reduceToPayload(ActionType.TencenthubChartSelection, null),
  tencenthubChartVersionList: generateFetcherReducer<RecordSet<TencenthubChartVersion>>({
    actionType: ActionType.FetchTencenthubChartVersionList,
    initialData: {
      recordCount: 0,
      records: [] as TencenthubChartVersion[]
    }
  }),
  tencenthubChartVersionSelection: reduceToPayload(ActionType.TencenthubChartVersionSelection, null),
  tencenthubChartReadMe: reduceToPayload(ActionType.TencenthubChartReadMe, null),

  otherChartUrl: reduceToPayload(ActionType.OtherChartUrl, ''),
  otherTypeSelection: reduceToPayload(ActionType.OtherType, OtherType.Public),
  otherUserName: reduceToPayload(ActionType.OtherUserName, ''),
  otherPassword: reduceToPayload(ActionType.OtherPassword, ''),
  kvs: reduceToPayload(ActionType.KeyValue, [])
});

export const HelmCreationReducer = (inputState, action) => {
  let state = inputState;
  // 销毁页面
  if (action.type === ActionType.ClearCreation) {
    state = undefined;
  }
  return TempReducer(state, action);
};
