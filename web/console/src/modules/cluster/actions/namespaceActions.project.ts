import { FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { extend, uuid } from '@tencent/qcloud-lib';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { projectNamespaceActions } from './projectNamespaceActions.project';
import { resourceActions } from './resourceActions';

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
    let {
      filter: { clusterId }
    } = namespaceQuery;
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
      let { subRoot, route, cluster, projectNamespaceList } = getState(),
        urlParams = router.resolve(route),
        { isNeedFetchNamespace, mode } = subRoot;

      let finder = projectNamespaceList.data.records.find(item => item.metadata.name === namespace);
      if (finder) {
        let clusterFinder = cluster.list.data.records.find(item => item.metadata.name === finder.spec.clusterName);
        if (clusterFinder) {
          dispatch(projectNamespaceActions.selectCluster(clusterFinder));
        }
      }

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
      mode !== 'create' && finder
        ? dispatch(
            resourceActions.poll({
              namespace,
              clusterId: finder.spec.clusterName,
              regionId: +route.queries['rid']
            })
          )
        : dispatch({
            type: ActionType.FetchResourceList + 'Done',
            payload: {
              data: {
                recordCount: 0,
                records: []
              },
              trigger: 'Done'
            }
          });
    };
  }
};

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions, restActions);
