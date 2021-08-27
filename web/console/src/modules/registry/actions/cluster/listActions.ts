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
import { RootState, Cluster, ChartInfoFilter, ClusterFilter } from '../../models';
import * as ActionType from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { namespaceActions } from '../namespace';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchClusterActions = createFFListActions<Cluster, ClusterFilter, ChartInfoFilter>({
  actionName: ActionType.ClusterList,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().clusterList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(listActions.selectCluster(record.data.records[0].metadata.name, record.data.data));
    }
  }
});

const restActions = {
  selectCluster: (cluster: string, chartInfoFilter?: ChartInfoFilter) => {
    return async (dispatch, getState: GetState) => {
      dispatch(listActions.selectByValue(cluster));
      dispatch(
        namespaceActions.list.applyFilter({
          cluster: cluster,
          chartInfoFilter: chartInfoFilter
        })
      );
    };
  }
};

export const listActions = extend({}, fetchClusterActions, restActions);
