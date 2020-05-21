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
    alarmRecord: createFFListReducer(ActionTypes.FetchAlarmRecord),
});
