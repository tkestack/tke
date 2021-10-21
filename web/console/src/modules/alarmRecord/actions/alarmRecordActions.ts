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
  createFFListActions,
  extend,
  FetchOptions,
  generateFetcherActionCreator,
  generateWorkflowActionCreator,
  OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, AlarmRecord, AlarmRecordFilter } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/**
 * 获取告警记录列表
 */
const FFModelAlarmActions = createFFListActions<AlarmRecord, AlarmRecordFilter>({
  actionName: ActionTypes.FetchAlarmRecord,
  fetcher: async (query, getState: GetState) => {
    const { alarmRecord } = getState();
    const alarmRecordData = alarmRecord.list.data;
    const response = await WebAPI.fetchAlarmRecord(query, { continueToken: alarmRecordData.continueToken });
    return {
      ...response,
      records: [...response.records].sort(
        (pre, current) => current?.metadata?.creationTimestamp - pre?.metadata?.creationTimestamp
      )
    };
  },
  getRecord: (getState: GetState) => {
    return getState().alarmRecord;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    if (record.data.recordCount) {
      const isNotNeedPoll = record.data.records.filter(item => item.status['phase'] === 'Terminating').length === 0;

      if (isNotNeedPoll) {
        dispatch(FFModelAlarmActions.clearPolling());
      }
    }
  }
});

const restActions = {
  // poll: () => {
  //     return async (dispatch: Redux.Dispatch, getState: GetState) => {
  //         dispatch(
  //             alarmRecordActions.polling({
  //                 delayTime: 5000
  //             })
  //         );
  //     };
  // },
};

export const alarmRecordActions = extend({}, FFModelAlarmActions, restActions);
