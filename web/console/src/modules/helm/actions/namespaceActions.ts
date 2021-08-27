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
import { helmActions } from './helmActions';
import { resourceConfig } from '@config';
import { extend, FetchOptions, generateFetcherActionCreator, uuid } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { CommonAPI, Resource, ResourceFilter, ResourceInfo } from '../../common';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { clusterActions } from './clusterActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch namespace list */
const fetchNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchNamespaceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let {
      listState: { cluster },
      namespaceQuery
    } = getState();
    let namespaceInfo: ResourceInfo = resourceConfig()['np'];
    let response = await CommonAPI.fetchResourceList({ query: namespaceQuery, resourceInfo: namespaceInfo });
    return {
      recordCount: response.recordCount,
      records: response.records.map(r => ({ id: r.metadata.name, name: r.metadata.name, displayName: r.metadata.name }))
    };
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { namespaceList, route, namespaceQuery } = getState();

    if (namespaceQuery.filter) {
      let defauleNamespace =
        route.queries['np'] || (namespaceList.data.recordCount ? namespaceList.data.records[0].id + '' : 'default');

      dispatch(namespaceActions.selectNamespace(defauleNamespace));
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
      let { route } = getState(),
        urlParams = router.resolve(route);

      dispatch({
        type: ActionType.SelectNamespace,
        payload: namespace
      });
      router.navigate(
        urlParams,
        Object.assign(route.queries, {
          np: namespace
        })
      );
      dispatch(helmActions.fetch());
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
