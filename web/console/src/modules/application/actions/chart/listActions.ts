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
import { RootState, Chart, ChartFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchChartActions = createFFListActions<Chart, ChartFilter>({
  actionName: ActionTypes.ChartList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartList(query);
    // 此部分若放置在onFinish，会出现lastVersion异步返回的现象，且ui不会更新
    if (response.recordCount > 0) {
      let records = response.records.map(r => {
        if (r.status.versions) {
          let sorted = r.status.versions.sort((a, b) => {
            let oDate1 = new Date(a.timeCreated);
            let oDate2 = new Date(b.timeCreated);
            return oDate1.getTime() > oDate2.getTime() ? -1 : 1;
          });
          r.lastVersion = sorted[0];
          r.sortedVersions = sorted;
        }
        return r;
      });
      response.records = records;
    }
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    let isNotNeedPoll = true;
    if (record.data.recordCount) {
      isNotNeedPoll =
        record.data.records.filter(
          item => item.status && item.status['phase'] && item.status['phase'] === 'Terminating'
        ).length === 0;
    }
    if (isNotNeedPoll) {
      dispatch(fetchChartActions.clearPolling());
    }
  }
});

const restActions = {
  /** 轮询操作 */
  poll: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        listActions.polling({
          delayTime: 3000
        })
      );
    };
  }
};

export const listActions = extend({}, fetchChartActions, restActions);
