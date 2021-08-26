/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { createFFListActions, extend } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { CommonAPI, Resource, ResourceFilter, ResourceInfo } from '../../common';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import { helmActions } from './helmActions';
import { namespaceActions } from './namespaceActions';

type GetState = () => RootState;

/** 集群列表的Actions */
const ListClusterActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async query => {
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let response = await CommonAPI.fetchResourceList({ query, resourceInfo: clusterInfo });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().listState.cluster;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    if (record.data.recordCount) {
      let defaultClusterId = route.queries['clusterId'];
      let defaultCluster = record.data.records.find(item => item.metadata.name === defaultClusterId);
      dispatch(clusterActions.selectCluster(defaultCluster || record.data.records[0]));
    }
  }
});

const restActions = {
  selectCluster: (cluster: Resource) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
      dispatch(ListClusterActions.select(cluster));

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
      dispatch(helmActions.checkClusterHelmStatus());
      /// #if tke
      dispatch(namespaceActions.applyFilter({ clusterId: cluster.metadata.name }));

      /// #endif
    };
  }
};

export const clusterActions = extend({}, ListClusterActions, restActions);
