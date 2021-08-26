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

import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, ProjectNamespace, ProjectNamespaceFilter, ChartInfoFilter } from '../../models';
import * as ActionType from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { appActions } from '../app';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchProjectNamespaceActions = createFFListActions<ProjectNamespace, ProjectNamespaceFilter>({
  actionName: ActionType.ProjectNamespaceList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().projectNamespaceList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(
        projectNamespaceListActions.selectProjectNamespace(
          record.data.records[0].metadata.namespace,
          record.data.records[0].spec.clusterName,
          record.data.records[0].spec.namespace,
          record.data.data
        )
      );
    }
  }
});

const restActions = {
  selectProjectNamespace: (
    projectID: string,
    cluster: string,
    namespace: string,
    chartInfoFilter?: ChartInfoFilter
  ) => {
    return async (dispatch, getState: GetState) => {
      let { route, appCreation, chartEditor } = getState(),
        urlParams = router.resolve(route);

      dispatch(projectNamespaceListActions.selectByValue(cluster + '/' + namespace));

      if (!urlParams['sub'] || urlParams['sub'] === 'chart') {
        if (!urlParams['mode'] || urlParams['mode'] === 'detail') {
          dispatch(
            appActions.create.updateCreationState({
              metadata: Object.assign({}, appCreation.metadata, {
                namespace: namespace
              }),
              spec: Object.assign({}, appCreation.spec, {
                targetCluster: cluster
              })
            })
          );

          //加载values.yaml
          dispatch(
            appActions.create.chart.applyFilter({
              cluster: cluster,
              namespace: namespace,
              metadata: { ...chartInfoFilter.metadata },
              chartVersion: chartInfoFilter.chartVersion,
              projectID: projectID
            })
          );
        }
      }
    };
  }
};

export const projectNamespaceListActions = extend({}, fetchProjectNamespaceActions, restActions);
