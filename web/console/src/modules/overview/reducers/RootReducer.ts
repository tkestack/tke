import { ClusterOverview, ClusterOverviewFilter } from './../models/RootState';
import { combineReducers } from 'redux';
import * as ActionType from '../constants/ActionType';
import { router } from '../router';
import { createFFObjectReducer } from '@tencent/ff-redux';

export const RootReducer = combineReducers({
  route: router.getReducer(),
  clusterOverview: createFFObjectReducer<ClusterOverview, ClusterOverviewFilter>(ActionType.ClusterOverview)
});
