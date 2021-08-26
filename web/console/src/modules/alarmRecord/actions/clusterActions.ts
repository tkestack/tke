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
import { FFReduxActionName } from '../../cluster/constants/Config';
import * as WebAPI from '../../cluster/WebAPI';
import { Cluster, ClusterFilter } from '../../common/';
import { RootState } from '../models';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query, query.filter.regionId);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  }
});

export const clusterActions = extend({}, FFModelClusterActions);
