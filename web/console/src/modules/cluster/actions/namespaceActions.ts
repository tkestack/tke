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
import { extend, FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { resourceActions } from './resourceActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch namespace list */
const fetchNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchNamespaceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    const { clusterVersion } = getState();
    // 获取当前的资源的配置
    const namespaceInfo = resourceConfig(clusterVersion)['ns'];
    const response = await WebAPI.fetchNamespaceList(getState().namespaceQuery, namespaceInfo);
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    const { namespaceList, route } = getState();

    const defauleNamespace =
      route.queries['np'] ||
      (namespaceList.data.recordCount && namespaceList.data.records.find(n => n.name === 'default').name) ||
      'default';

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
      const { subRoot, route } = getState(),
        urlParams = router.resolve(route),
        { isNeedFetchNamespace, mode } = subRoot;

      dispatch({
        type: ActionType.SelectNamespace,
        payload: namespace
      });

      // 这里进行路由的更新，如果不需要命名空间的话，路由就不需要有np的信息
      if (isNeedFetchNamespace) {
        router.navigate(urlParams, Object.assign({}, route.queries, { np: namespace }));
      } else {
        const routeQueries = Object.assign({}, route.queries, { np: undefined });
        router.navigate(urlParams, JSON.parse(JSON.stringify(routeQueries)));
      }

      // 初始化或者变更Resource的信息，在创建页面当中，变更ns，不需要拉取resource

      if (mode !== 'create') {
        dispatch(resourceActions.resetPaging());
        dispatch(resourceActions.poll());
      }
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
