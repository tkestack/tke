import {
    createFFListActions, extend, FetchOptions, generateFetcherActionCreator,
    generateWorkflowActionCreator, OperationTrigger
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
        let response = await WebAPI.fetchAlarmRecord(query, { continueToken: alarmRecordData.continueToken });
        return response;
    },
    getRecord: (getState: GetState) => {
        return getState().alarmRecord;
    },
    onFinish: (record, dispatch: Redux.Dispatch) => {
        if (record.data.recordCount) {
            let isNotNeedPoll = record.data.records.filter(item => item.status['phase'] === 'Terminating').length === 0;

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
