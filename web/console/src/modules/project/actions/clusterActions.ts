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
import { createFFListActions, extend } from '@tencent/ff-redux';

import { Cluster, ClusterFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: 'cluster',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch, getState: GetState) => {}
});

const restActions = {
  selectCluster: (cluster: Cluster[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(FFModelClusterActions.select(cluster[0]));
    };
  }
};

export const clusterActions = extend({}, FFModelClusterActions, restActions);
