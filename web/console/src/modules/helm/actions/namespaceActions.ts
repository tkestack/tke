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
