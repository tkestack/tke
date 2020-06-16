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
    projectNamespaceList.data.records.forEach(item => {
      namespaceList.push({ id: uuid(), name: item.spec.namespace });
    });

    return { recordCount: namespaceList.length, records: namespaceList };
  },
  finish: (dispatch, getState: GetState) => {
    let { namespaceList, route } = getState();
    let defauleNamespace = (namespaceList.data.recordCount && namespaceList.data.records[0].name) || '';
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

      let finder = projectNamespaceList.data.records.find(item => item.spec.namespace === namespace);
      if (!finder) {
        finder = projectNamespaceList.data.records.length ? projectNamespaceList.data.records[0] : null;
      }

      if (finder) {
        router.navigate(
          urlParams,
          Object.assign(route.queries, {
            np: finder.spec.namespace
          })
        );
        dispatch({
          type: ActionType.SelectNamespace,
          payload: finder.spec.namespace
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
