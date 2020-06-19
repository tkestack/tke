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
    let { clusterVersion } = getState();
    // 获取当前的资源的配置
    let namespaceInfo = resourceConfig(clusterVersion)['ns'];
    let response = await WebAPI.fetchNamespaceList(getState().namespaceQuery, namespaceInfo);
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { namespaceList, route } = getState();

    let defauleNamespace =
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
      let { subRoot, route } = getState(),
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
        let routeQueries = Object.assign({}, route.queries, { np: undefined });
        router.navigate(urlParams, JSON.parse(JSON.stringify(routeQueries)));
      }

      // 初始化或者变更Resource的信息，在创建页面当中，变更ns，不需要拉取resource
      mode !== 'create' && dispatch(resourceActions.poll());
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
