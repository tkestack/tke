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

import { extend, FetchOptions, generateFetcherActionCreator, uuid } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

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
    let { projectNamespaceList, namespaceQuery } = getState();
    // 获取当前的资源的配置
    let namespaceList = [];
    projectNamespaceList.data.records
      .filter(item => item.status.phase === 'Available')
      .forEach(item => {
        namespaceList.push({
          id: uuid(),
          name: item.metadata.name,
          displayName: `${item.spec.namespace}(${item.spec.clusterName})`,
          clusterDisplayName: item.spec.clusterDisplayName,
          clusterName: item.spec.clusterName,
          namespace: item.spec.namespace
        });
      });

    return { recordCount: namespaceList.length, records: namespaceList };
  },
  finish: (dispatch, getState: GetState) => {
    let { namespaceList, route } = getState();
    let defauleNamespace =
      route.queries['np'] || (namespaceList.data.recordCount && namespaceList.data.records[0].name) || '';
    dispatch(namespaceActions.selectNamespace(defauleNamespace));
  }
});

/** query namespace list action */
const queryNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.QueryNamespaceList,
  bindFetcher: fetchNamespaceActions
});

const restActions = {
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let {
          route,
          listState: { cluster },
          projectNamespaceList
        } = getState(),
        urlParams = router.resolve(route);

      let finder = projectNamespaceList.data.records.find(item => item.metadata.name === namespace);
      if (!finder) {
        finder = projectNamespaceList.data.records.length ? projectNamespaceList.data.records[0] : null;
      }

      if (finder) {
        router.navigate(
          urlParams,
          Object.assign(route.queries, {
            np: finder.metadata.name
          })
        );
        dispatch({
          type: ActionType.SelectNamespace,
          payload: finder.metadata.name
        });

        dispatch(
          clusterActions.selectCluster({
            id: finder.spec.clusterName,
            metadata: { name: finder.spec.clusterName },
            spec: { dispalyName: '-' },
            status: { version: '1.12.4' }
          })
        );
      } else {
        dispatch({
          type: ActionType.FetchHelmList + 'Done',
          payload: {
            data: {
              recordCount: 0,
              records: []
            },
            trigger: 'Done'
          }
        });
        dispatch({
          type: ActionType.FetchInstallingHelmList + 'Done',
          payload: {
            data: {
              recordCount: 0,
              records: []
            },
            trigger: 'Done'
          }
        });
      }
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
