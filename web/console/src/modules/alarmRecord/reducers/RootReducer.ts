import { combineReducers } from 'redux';
import { router } from '../router';
import  * as ActionTypes from '../constants/ActionTypes';
import { Cluster } from '../../common';
import { ClusterFilter } from '../models';
import { FFReduxActionName } from '../../cluster/constants/Config';
import {
    createFFListReducer,
} from '@tencent/ff-redux';

export const RootReducer = combineReducers({
    route: router.getReducer(),
    alarmRecord: createFFListReducer(ActionTypes.FetchAlarmRecord),
    cluster: createFFListReducer<Cluster, ClusterFilter>(FFReduxActionName.CLUSTER),
});
