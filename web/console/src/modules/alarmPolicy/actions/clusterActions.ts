import { createListAction } from '@tencent/redux-list';
import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../../cluster/WebAPI';
import { setClusterId } from '../../../../helpers';
import { router } from '../router';
import { TellIsNeedFetchNS, TellIsNotNeedFetchResource } from '../components/resource/ResourceSidebarPanel';
import { resourceConfig } from '../../../../config';
import { Cluster, ClusterFilter, ResourceInfo } from '../../common/';
import { FFReduxActionName } from '../../cluster/constants/Config';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createListAction<Cluster, ClusterFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query, query.filter.regionId);
    let ps = await WebAPI.fetchPrometheuses();
    let clusterHasPs = {};
    for (let p of ps.records) {
      clusterHasPs[p.spec.clusterName] = true;
    }
    for (let record of response.records) {
      record.spec.hasPrometheus = clusterHasPs[record.metadata.name];
    }
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    let { cluster, route } = getState();

    if (cluster.list.data.recordCount) {
      let routeClusterId = route.queries['clusterId'];
      let finder = routeClusterId ? cluster.list.data.records.find(c => c.metadata.name === routeClusterId) : undefined;
      let defaultCluster = finder ? finder : cluster.list.data.records[0];
      dispatch(clusterActions.selectCluster(defaultCluster));
    }
  }
});

const restActions = {
  selectCluster(cluster: Cluster, isNeedInitClusterVersion: boolean = false) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { regionSelection, route } = getState(),
        urlParams = router.resolve(route);
      if (cluster) {
        dispatch(FFModelClusterActions.select(cluster));

        dispatch(clusterActions.initClusterVersion(cluster.status.version));

        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
        dispatch(
          alarmPolicyActions.applyFilter({ regionId: +regionSelection.value, clusterId: cluster.metadata.name })
        );
      } else {
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: '' }));
        dispatch(alarmPolicyActions.clear());
      }
    };
  },
  /** 初始化集群的版本 */
  initClusterVersion: (clusterVersion: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let k8sVersion = clusterVersion
        .split('.')
        .slice(0, 2)
        .join('.');
      dispatch({
        type: ActionType.InitClusterVersion,
        payload: k8sVersion
      });
    };
  }
};

export const clusterActions = extend({}, FFModelClusterActions, restActions);
