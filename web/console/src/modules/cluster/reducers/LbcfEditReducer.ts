import { initLbcfBGPort, initSelector, initLbcfBackGroupEdition } from './../constants/initState';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common/models';
import { LbcfConfig } from '../constants/Config';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.Gate_Name, ''),

  v_name: reduceToPayload(ActionType.V_Gate_Name, initValidator),

  namespace: reduceToPayload(ActionType.Gate_Namespace, ''),

  v_namespace: reduceToPayload(ActionType.V_Gate_Namespace, initValidator),

  config: reduceToPayload(ActionType.Lbcf_Config, [
    {
      key: 'loadBalancerID',
      value: LbcfConfig.find(o => o.value === 'loadBalancerID').defaultValue || ''
    },
    {
      key: 'loadBalancerType',
      value: LbcfConfig.find(o => o.value === 'loadBalancerType').defaultValue || ''
    },
    {
      key: 'vpcID',
      value: LbcfConfig.find(o => o.value === 'vpcID').defaultValue || ''
    },
    {
      key: 'listenerProtocol',
      value: LbcfConfig.find(o => o.value === 'listenerProtocol').defaultValue || ''
    },
    {
      key: 'listenerPort',
      value: LbcfConfig.find(o => o.value === 'listenerPort').defaultValue || ''
    }
  ]),
  args: reduceToPayload(ActionType.Lbcf_Args, []),

  /** LBReducer*/
  // vpcSelection: reduceToPayload(ActionType.GLB_VpcSelection, ''),

  // clbList: reduceToPayload(ActionType.GLB_FecthClb, []),

  // clbSelection: reduceToPayload(ActionType.GLB_SelectClb, ''),
  // v_clbSelection: reduceToPayload(ActionType.V_GLB_SelectClb, initValidator),

  // createLbWay: reduceToPayload(ActionType.GLB_SwitchCreateLbWay, 'new'),
  /** LBReducer*/

  /** backGroupReducer*/

  lbcfBackGroupEditions: reduceToPayload(ActionType.GBG_UpdateLbcfBackGroup, [initLbcfBackGroupEdition])

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
