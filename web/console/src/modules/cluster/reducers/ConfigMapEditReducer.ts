import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload, ReduxAction } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { initVariable } from '../models';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.CM_Name, ''),

  v_name: reduceToPayload(ActionType.V_CM_Name, initValidator),

  namespace: reduceToPayload(ActionType.CM_Namespace, 'default'),

  variables: (state = [initVariable], action: ReduxAction<any>) => {
    switch (action.type) {
      case ActionType.CM_AddVariable:
      case ActionType.CM_EditVariable:
      case ActionType.CM_DeleteVariable:
      case ActionType.CM_ValidateVariable:
        return action.payload;
      default:
        return state;
    }
  }
});

export const ConfigMapEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建pv的界面
  if (action.type === ActionType.ClearConfigMapEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
