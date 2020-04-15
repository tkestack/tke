import { combineReducers } from 'redux';

import { reduceToPayload, createFFListReducer } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';

const TempReducer = combineReducers({
  ffResourceList: createFFListReducer(FFReduxActionName.Resource_Workload),

  resourceMultipleSelection: reduceToPayload(ActionType.SelectMultipleResource, []),

  resourceDeleteSelection: reduceToPayload(ActionType.SelectDeleteResource, [])
});

export const ResourceReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResource) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
