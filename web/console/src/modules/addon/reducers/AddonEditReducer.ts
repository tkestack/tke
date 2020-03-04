import { combineReducers } from 'redux';

import { reduceToPayload } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { PeEditReducer } from './PeEditReducer';

const TempReducer = combineReducers({
  addonName: reduceToPayload(ActionType.AddonName, ''),

  peEdit: PeEditReducer
});

export const AddonEditReducer = (state, action) => {
  let newState = state;
  // 销毁新建扩展组件页面
  if (action.type === ActionType.ClearAddonEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
