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
import { reduceNs } from '@helper';
import { resourceConfig } from '@config';
import { ResourceInfo } from '@src/modules/common';
import { createFFListActions, extend } from '@tencent/ff-redux';

import { FFReduxActionName } from '../../cluster/constants/Config';
import * as WebAPI from '../../cluster/WebAPI';
import { Cluster, ClusterFilter } from '../../common/';
import * as ActionType from '../constants/ActionType';
import { RootState, AddonStatus } from '../models';
import { router } from '../router';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    const response = await WebAPI.fetchClusterList(query, query.filter.regionId);
    const ps = await WebAPI.fetchPrometheuses();
    const clusterHasPs = {};
    for (const p of ps.records) {
      clusterHasPs[p.spec.clusterName] = true;
    }
    for (const record of response.records) {
      record.spec.hasPrometheus = clusterHasPs[record.metadata.name];
    }
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    const { cluster, route } = getState();

    if (cluster.list.data.recordCount) {
      const routeClusterId = route.queries['clusterId'];
      const finder = routeClusterId
        ? cluster.list.data.records.find(c => c.metadata.name === routeClusterId)
        : undefined;
      const defaultCluster = finder ? finder : cluster.list.data.records[0];
      dispatch(clusterActions.selectCluster(defaultCluster));
    }
  }
});

const restActions = {
  selectCluster(cluster: Cluster, isNeedInitClusterVersion = false) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { regionSelection, route } = getState(),
        urlParams = router.resolve(route);
      if (cluster) {
        dispatch(FFModelClusterActions.select(cluster));

        dispatch(clusterActions.initClusterVersion(cluster.status.version));

        dispatch(clusterActions.fetchClusterAddon(cluster.metadata.name, 1));

        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
        /// #if project
        dispatch(
          alarmPolicyActions.applyFilter({
            regionId: +regionSelection.value,
            clusterId: cluster.metadata.name,
            namespace: reduceNs(route.queries['np']),
            alarmPolicyType: 'pod'
          })
        );
        /// #endif
        /// #if tke
        dispatch(
          alarmPolicyActions.applyFilter({ regionId: +regionSelection.value, clusterId: cluster.metadata.name })
        );
        /// #endif
      } else {
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: '' }));
        dispatch(alarmPolicyActions.clear());
      }
    };
  },
  /** 初始化集群的版本 */
  initClusterVersion: (clusterVersion: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const k8sVersion = clusterVersion.split('.').slice(0, 2).join('.');
      dispatch({
        type: ActionType.InitClusterVersion,
        payload: k8sVersion
      });
    };
  },
  /** 获取当前集群开启的Addon */
  fetchClusterAddon: (clusterId, regionId) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { clusterVersion } = getState();
      const clusterInfo: ResourceInfo = resourceConfig(clusterVersion).cluster;
      const response = await WebAPI.fetchExtraResourceList(
        {
          filter: {
            regionId: regionId,
            specificName: clusterId
          }
        },
        clusterInfo,
        false,
        'addons'
      );
      const addons: AddonStatus = {};
      response.records.forEach(item => {
        addons[item.spec.type] = {
          status: item.status.phase,
          name: item.metadata.name
        };
      });
      dispatch({
        type: ActionType.FetchClusterAddons,
        payload: addons
      });
    };
  }
};

export const clusterActions = extend({}, FFModelClusterActions, restActions);
