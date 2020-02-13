import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common/models';

import { k8sVersionList, GPUTYPE } from '../constants/Config';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.IC_Name, ''),

  v_name: reduceToPayload(ActionType.v_IC_Name, initValidator),

  networkDevice: reduceToPayload(ActionType.IC_NetworkDevice, 'eth0'),

  v_networkDevice: reduceToPayload(ActionType.v_IC_NetworkDevice, initValidator),

  maxClusterServiceNum: reduceToPayload(ActionType.IC_MaxClusterServiceNum, 256),

  maxNodePodNum: reduceToPayload(ActionType.IC_MaxNodePodNum, 256),

  k8sVersion: reduceToPayload(ActionType.IC_K8SVersion, '1.16.6'),

  k8sVersionList: reduceToPayload(ActionType.IC_FetchK8SVersion, []),

  cidr: reduceToPayload(ActionType.IC_Cidr, '10.244.0.0/16'),

  computerList: reduceToPayload(ActionType.IC_ComputerList, []),
  computerEdit: reduceToPayload(ActionType.IC_ComputerEdit, null),
  vipAddress: reduceToPayload(ActionType.IC_VipAddress, ''),
  vipPort: reduceToPayload(ActionType.IC_VipPort, ''),

  v_vipAddress: reduceToPayload(ActionType.v_IC_VipAddress, initValidator),

  v_vipPort: reduceToPayload(ActionType.v_IC_VipPort, initValidator),

  vip: reduceToPayload(ActionType.v_IC_Vip, false),

  gpu: reduceToPayload(ActionType.v_IC_Gpu, false),

  gpuType: reduceToPayload(ActionType.v_IC_GpuType, GPUTYPE.PGPU)
});

export const CreateICReducer = (state, action) => {
  let newState = state;
  // 销毁创建namespace 页面
  if (action.type === ActionType.IC_Clear) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
