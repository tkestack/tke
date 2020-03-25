import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.N_Name, ''),

  v_name: reduceToPayload(ActionType.NV_Name, initValidator),

  description: reduceToPayload(ActionType.N_Description, ''),

  v_description: reduceToPayload(ActionType.NV_Description, initValidator)
});

export const NamespaceEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建namespace 页面
  if (action.type === ActionType.ClearNamespaceEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
