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
import { RootState, ChartGroup, ChartGroupFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchChartGroupActions = createFFListActions<ChartGroup, ChartGroupFilter>({
  actionName: ActionTypes.ChartGroupList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartGroupList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartGroupList;
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
      dispatch(fetchChartGroupActions.clearPolling());
    }
  }
});

const restActions = {};

export const listActions = extend({}, fetchChartGroupActions, restActions);
