import { FFReduxActionName } from './../constants/Config';
import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload, createFFListReducer } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { LbcfConfig } from '../constants/Config';
import { initLbcfBackGroupEdition, initLbcfBGPort, initSelector } from '../constants/initState';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.Gate_Name, ''),

  v_name: reduceToPayload(ActionType.V_Gate_Name, initValidator),

  namespace: reduceToPayload(ActionType.Gate_Namespace, ''),

  v_namespace: reduceToPayload(ActionType.V_Gate_Namespace, initValidator),

  config: reduceToPayload(ActionType.Lbcf_Config, [
    {
      key: '',
      value: ''
    }
  ]),
  args: reduceToPayload(ActionType.Lbcf_Args, [
    {
      key: '',
      value: ''
    }
  ]),

  v_config: reduceToPayload(ActionType.V_Lbcf_Config, initValidator),

  v_args: reduceToPayload(ActionType.V_Lbcf_Args, initValidator),
  /** LBReducer*/
  // vpcSelection: reduceToPayload(ActionType.GLB_VpcSelection, ''),

  // clbList: reduceToPayload(ActionType.GLB_FecthClb, []),

  // clbSelection: reduceToPayload(ActionType.GLB_SelectClb, ''),
  // v_clbSelection: reduceToPayload(ActionType.V_GLB_SelectClb, initValidator),

  // createLbWay: reduceToPayload(ActionType.GLB_SwitchCreateLbWay, 'new'),
  /** LBReducer*/

  /** backGroupReducer*/

  lbcfBackGroupEditions: reduceToPayload(ActionType.GBG_UpdateLbcfBackGroup, [initLbcfBackGroupEdition]),

  driver: createFFListReducer(
    FFReduxActionName.LBCF_DRIVER,
    null,
    x => x.metadata.name,
    x => x.metadata.name
  )

  // gameAppList: reduceToPayload(ActionType.GBG_FetchGameApp, []),

  // gameAppSelection: reduceToPayload(ActionType.GBG_SelectGameApp, ''),

  // isShowGameAppDialog: reduceToPayload(ActionType.GBG_ShowGameAppDialog, false)
  /** backGroupReducer*/
});

export const LbcfEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建 Lbcf 界面
  if (action.type === ActionType.ClearLbcfEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
