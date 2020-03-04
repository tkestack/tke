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
    projectNamespaceList.data.records.forEach(item => {
      namespaceList.push({
        id: uuid(),
        name: item.metadata.name,
        clusterVersion: item.spec.clusterVersion,
        clusterId: item.spec.clusterVersion,
        clusterDisplayName: item.spec.clusterDisplayName
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
        if (finder) {
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
        dispatch(alarmPolicyActions.selectsWorkLoadNamespace(namespace));
      } else {
        dispatch(clusterActions.selectCluster(undefined));
        dispatch(clusterActions.select(null));
      }
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
