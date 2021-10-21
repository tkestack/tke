/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { GPUTYPE, k8sVersionList, CreateICVipType, ContainerRuntimeEnum } from '../constants/Config';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.IC_Name, ''),

  v_name: reduceToPayload(ActionType.v_IC_Name, initValidator),

  networkDevice: reduceToPayload(ActionType.IC_NetworkDevice, 'eth0'),

  v_networkDevice: reduceToPayload(ActionType.v_IC_NetworkDevice, initValidator),

  maxClusterServiceNum: reduceToPayload(ActionType.IC_MaxClusterServiceNum, 256),

  maxNodePodNum: reduceToPayload(ActionType.IC_MaxNodePodNum, 256),

  k8sVersion: reduceToPayload(ActionType.IC_K8SVersion, ''),

  k8sVersionList: reduceToPayload(ActionType.IC_FetchK8SVersion, []),

  cidr: reduceToPayload(ActionType.IC_Cidr, '10.244.0.0/16'),

  computerList: reduceToPayload(ActionType.IC_ComputerList, []),
  computerEdit: reduceToPayload(ActionType.IC_ComputerEdit, null),
  vipAddress: reduceToPayload(ActionType.IC_VipAddress, ''),
  vipPort: reduceToPayload(ActionType.IC_VipPort, '6443'),

  v_vipAddress: reduceToPayload(ActionType.v_IC_VipAddress, initValidator),

  v_vipPort: reduceToPayload(ActionType.v_IC_VipPort, initValidator),

  vipType: reduceToPayload(ActionType.v_IC_Vip, CreateICVipType.unuse),

  gpu: reduceToPayload(ActionType.v_IC_Gpu, false),

  merticsServer: reduceToPayload(ActionType.v_IC_Mertics_server, true),

  cilium: reduceToPayload(ActionType.v_IC_Cilium, 'Galaxy'),

  networkMode: reduceToPayload(ActionType.v_IC_NetworkMode, 'overlay'),

  asNumber: reduceToPayload(ActionType.IC_AS, ''),
  v_asNumber: reduceToPayload(ActionType.v_IC_AS, initValidator),

  switchIp: reduceToPayload(ActionType.IC_SwitchIp, ''),
  v_switchIp: reduceToPayload(ActionType.v_IC_SwitchIp, initValidator),

  gpuType: reduceToPayload(ActionType.v_IC_GpuType, GPUTYPE.PGPU),

  containerRuntime: reduceToPayload(ActionType.IC_EnableContainerRuntime, ContainerRuntimeEnum.CONTAINERD)
});

export const CreateICReducer = (state, action) => {
  let newState = state;
  // 销毁创建namespace 页面
  if (action.type === ActionType.IC_Clear) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
