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
import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, Cluster } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { namespaceActions } from '../namespace';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchClusterActions = createFFListActions<Cluster, void>({
  actionName: ActionTypes.ClusterList,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().clusterList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      let { route } = getState();
      const cluster = route.queries['cluster'];
      let exist = record.data.records.find(r => {
        return r.metadata.name === cluster;
      });
      if (exist) {
        dispatch(listActions.selectCluster(exist.metadata.name));
      } else {
        dispatch(listActions.selectCluster(record.data.records[0].metadata.name));
      }
    }
  }
});

const restActions = {
  selectCluster: (cluster: string) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { cluster: cluster }));

      dispatch(listActions.selectByValue(cluster));
      dispatch(
        namespaceActions.list.applyFilter({
          cluster: cluster
        })
      );
    };
  }
};

export const listActions = extend({}, fetchClusterActions, restActions);
