import { alarmPolicyActions } from './alarmPolicyActions';
import { clusterActions } from './clusterActions';
import { extend, RecordSet, ReduxAction, uuid } from '@tencent/qcloud-lib';
import { FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { FFReduxActionName } from '../../cluster/constants/Config';
import { Cluster, ClusterFilter, ResourceInfo } from '../../common/';
import { uniq } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

import { namespaceActions } from './namespaceActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch namespacesetlist */
const fetchProjectNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchProjectNamespace,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { projectNamespaceQuery } = getState();
    let response = await WebAPI.fetchProjectNamespaceList(projectNamespaceQuery);

    let clusterList = uniq(
      response.records.map(namespace => ({
        clusterId: namespace.spec.clusterName,
        clusterDisplayName: namespace.spec.clusterDisplayName,
        clusterVersion: namespace.spec.clusterVersion
      })),
      'clusterId'
    );
    let clusterListRecord = clusterList.map(item => {
      return {
        metadata: { name: item.clusterId },
        spec: { displayName: item.clusterDisplayName, hasPrometheus: false },
        status: { version: item.clusterVersion }
      };
    });
    let ps = await WebAPI.fetchPrometheuses();
    let clusterHasPs = {};
    for (let p of ps.records) {
      clusterHasPs[p.spec.clusterName] = true;
    }
    for (let record of clusterListRecord) {
      record.spec.hasPrometheus = clusterHasPs[record.metadata.name];
    }

    dispatch(projectNamespaceActions.initClusterList(clusterListRecord));
    return response;
  },
  finish: async (dispatch: Redux.Dispatch, getState: GetState) => {
    dispatch(namespaceActions.fetch());
  }
});

/** query namespace list action */
const queryProjectNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.QueryProjectNamespace,
  bindFetcher: fetchProjectNamespaceActions
});

const restActions = {
  /** 初始化 NamespaceList列表 */
  initProjectList: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, projectSelection } = getState();
      let portalResourceInfo = resourceConfig().portal;
      let portal = await WebAPI.fetchUserPortal(portalResourceInfo);
      let userProjectList = Object.keys(portal.projects).map(key => {
        return {
          name: key,
          displayName: portal.projects[key]
        };
      });
      dispatch({
        type: ActionType.InitProjectList,
        payload: userProjectList
      });
      let defaultProjectName = projectSelection
        ? projectSelection
        : route.queries['projectName']
        ? route.queries['projectName']
        : userProjectList.length
        ? userProjectList[0].name
        : '';
      defaultProjectName && dispatch(projectNamespaceActions.selectProject(defaultProjectName));
    };
  },

  /** 选择业务 */
  selectProject: (project: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      dispatch({
        type: ActionType.ProjectSelection,
        payload: project
      });
      dispatch(projectNamespaceActions.applyFilter({ specificName: project }));
      router.navigate(
        urlParams,
        Object.assign(route.queries, {
          projectName: project
        })
      );
    };
  },

  /** 初始化集群列表 */
  initClusterList: clusterList => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result: RecordSet<Cluster> = {
        recordCount: clusterList.length,
        records: clusterList
      };

      dispatch({
        type: FFReduxActionName.CLUSTER + '_FetchDone',
        payload: {
          data: result,
          trigger: 'Done'
        }
      });
    };
  },

  selectCluster(cluster: Cluster, isNeedInitClusterVersion: boolean = false) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { regionSelection, route } = getState(),
        urlParams = router.resolve(route);
      if (cluster) {
        dispatch(clusterActions.select(cluster));
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
        dispatch(
          alarmPolicyActions.applyFilter({ regionId: +regionSelection.value, clusterId: cluster.metadata.name })
        );
      } else {
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: '' }));
        dispatch(alarmPolicyActions.clear());
        dispatch({
          type: 'AlarmPolicy_FetchDone',
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

export const projectNamespaceActions = extend(fetchProjectNamespaceActions, queryProjectNamespaceActions, restActions);
