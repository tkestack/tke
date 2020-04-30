import {
    createFFListActions, extend, FetchOptions, generateFetcherActionCreator,
    generateWorkflowActionCreator, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, Audit, AuditFilter } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
    noCache: false
};

const FFModelAuditActions = createFFListActions<Audit, AuditFilter>({
    actionName: ActionTypes.FetchAuditList,
    fetcher: async (query, getState: GetState) => {
        let response = await WebAPI.fetchAuditList(query);
        return response;
    },
    getRecord: (getState: GetState) => {
        return getState().auditList;
    },
    onFinish: (record, dispatch: Redux.Dispatch) => {
        if (record.data.recordCount) {
            let isNotNeedPoll = record.data.records.filter(item => item.status['phase'] === 'Terminating').length === 0;

            if (isNotNeedPoll) {
                dispatch(FFModelAuditActions.clearPolling());
            }
        }
    }
});

/**
 * 获取策略
 */
const getAuditFilterCondition = generateFetcherActionCreator({
    actionType: ActionTypes.FetchAuditFilterCondition,
    fetcher: async (getState: GetState, options: FetchOptions, dispatch) => {
        let result = await WebAPI.fetchAuditFilterCondition();
        return result;
    }
});

const restActions = {
    poll: () => {
        return async (dispatch: Redux.Dispatch, getState: GetState) => {
            dispatch(
                auditActions.polling({
                    delayTime: 5000
                })
            );
        };
    },
    getAuditFilterCondition,
};

export const auditActions = extend({}, FFModelAuditActions, restActions);
