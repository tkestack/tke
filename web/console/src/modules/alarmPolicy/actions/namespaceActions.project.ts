import { extend, FetchOptions, generateFetcherActionCreator, uuid } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { alarmPolicyActions } from './alarmPolicyActions';
import { clusterActions } from './clusterActions';
import { projectNamespaceActions } from './projectNamespaceActions.project';

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
          clusterVersion: item.spec.clusterVersion,
          clusterId: item.spec.clusterName,
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
      let { route, cluster, projectNamespaceList } = getState(),
        urlParams = router.resolve(route);

      router.navigate(
        urlParams,
        Object.assign(route.queries, {
          np: namespace
        })
      );
      dispatch({
        type: ActionType.SelectNamespace,
        payload: namespace
      });

      if (namespace) {
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

          let clusterId = finder.spec.clusterName;
          let clusterFinder = cluster.list.data.records.find(cluster => cluster.metadata.name === clusterId);
          if (clusterFinder) {
            dispatch(projectNamespaceActions.selectCluster(clusterFinder));
          } else {
            dispatch(
              clusterActions.selectCluster({
                id: clusterId,
                metadata: { name: clusterId },
                spec: { dispalyName: '-' },
                status: { version: '1.16.6' }
              })
            );
          }
        }
        dispatch(alarmPolicyActions.selectsWorkLoadNamespace(finder.metadata.name));
      } else {
        dispatch(clusterActions.selectCluster(undefined));
        dispatch(clusterActions.select(null));
      }
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
