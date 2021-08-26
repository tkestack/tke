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
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { ClusterHelmStatus, FFReduxActionName, OtherType } from '../constants/Config';
import { Helm, InstallingHelm, TencenthubChartVersion } from '../models';

const TempReducer = combineReducers({
  region: createFFListReducer(FFReduxActionName.REGION),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

  helmList: generateFetcherReducer<RecordSet<Helm>>({
    actionType: ActionType.FetchHelmList,
    initialData: {
      recordCount: 0,
      records: [] as Helm[]
    }
  }),

  clusterHelmStatus: reduceToPayload(ActionType.ClusterHelmStatus, {
    code: ClusterHelmStatus.NONE,
    reason: ''
  }),

  helmQuery: generateQueryReducer({
    actionType: ActionType.QueryHelmList,
    initialState: {
      paging: {
        pageSize: 10
      }
    }
  }),

  helmSelection: reduceToPayload(ActionType.SelectHelm, null),

  installingHelmList: generateFetcherReducer<RecordSet<InstallingHelm>>({
    actionType: ActionType.FetchInstallingHelmList,
    initialData: {
      recordCount: 0,
      records: [] as InstallingHelm[]
    }
  }),
  installingHelmSelection: reduceToPayload(ActionType.SelectInstallingHelm, null),
  installingHelmDetail: reduceToPayload(ActionType.FetchInstallingHelm, null),
  tencenthubChartVersionList: generateFetcherReducer<RecordSet<TencenthubChartVersion>>({
    actionType: ActionType.TableFetchTencenthubChartVersionList,
    initialData: {
      recordCount: 0,
      records: [] as TencenthubChartVersion[]
    }
  }),
  tencenthubChartVersionSelection: reduceToPayload(ActionType.TableTencenthubChartVersionSelection, null),
  token: reduceToPayload(ActionType.TableTencenthubToken, ''),

  otherChartUrl: reduceToPayload(ActionType.ListOtherChartUrl, ''),
  otherTypeSelection: reduceToPayload(ActionType.ListOtherType, OtherType.Public),
  otherUserName: reduceToPayload(ActionType.ListOtherUserName, ''),
  otherPassword: reduceToPayload(ActionType.ListOtherPassword, ''),

  isValid: reduceToPayload(ActionType.ListUpdateIsValid, {
    otherChartUrl: '',
    otherUserName: '',
    otherPassword: ''
  }),
  kvs: reduceToPayload(ActionType.ListKeyValue, [])
});

export const ListReducer = (inputState, action) => {
  let state = inputState;
  // 销毁详情页面
  if (action.type === ActionType.ClearList) {
    state = undefined;
  }
  //切换集群时清理掉部分数据
  if (action.type === ActionType.ClearListOnClusterChange) {
    delete state.helmList;
    delete state.helmQuery;
    delete state.helmSelection;
    delete state.installingHelmList;
    delete state.installingHelmSelection;
    delete state.installingHelmDetail;
  }
  return TempReducer(state, action);
};
