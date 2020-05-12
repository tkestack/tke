import { combineReducers } from 'redux';
import { router } from '../router';
import  * as ActionTypes from '../constants/ActionTypes';
import {
    createFFListReducer,
    generateWorkflowReducer,
    reduceToPayload,
    generateFetcherReducer
} from '@tencent/ff-redux';


export const RootReducer = combineReducers({
    route: router.getReducer(),
    auditList: createFFListReducer(ActionTypes.FetchAuditList),
    auditFilterCondition: generateFetcherReducer<Object>({
        actionType: ActionTypes.FetchAuditFilterCondition,
        initialData: {}
    }),
});
