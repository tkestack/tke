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
import { extend } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import * as WebAPI from '../../cluster/WebAPI';
import { ResourceInfo } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/** fetch namespace list */
const fetchNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchNamespaceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    const { clusterVersion } = getState();
    const namespaceInfo: ResourceInfo = resourceConfig(clusterVersion)['np'];
    const response = await WebAPI.fetchNamespaceList(getState().namespaceQuery, namespaceInfo);
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    const { namespaceList, route, namespaceQuery } = getState();

    if (namespaceQuery.filter.default) {
      const defauleNamespace =
        (namespaceList.data.recordCount && namespaceList.data.records.find(n => n.name === 'default').name) ||
        'default';

      dispatch(alarmPolicyActions.selectsWorkLoadNamespace(defauleNamespace));
    }
  }
});

/** query namespace list action */
const queryNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.QueryNamespaceList,
  bindFetcher: fetchNamespaceActions
});

//占位
const restActions = {
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch(alarmPolicyActions.selectsWorkLoadNamespace(namespace));
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
