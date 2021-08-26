/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { deepClone, ReduxAction, uuid } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { ContainerRuntimeEnum } from '../constants/Config';
import { ICComponter, LabelsKeyValue, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

export const createICAction = {
  /** 更新cluser的名称 */
  inputClusterName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.IC_Name,
      payload: name
    };
  },

  fetchK8sVersion: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const response = await WebAPI.fetchCreateICK8sVersion();
      dispatch({
        type: ActionType.IC_FetchK8SVersion,
        payload: response
      });

      if (response.length) {
        dispatch({
          type: ActionType.IC_K8SVersion,
          payload: response[0].value
        });
      }
    };
  },

  selectK8SVersion: (k8sVersion: string): ReduxAction<string> => {
    return {
      type: ActionType.IC_K8SVersion,
      payload: k8sVersion
    };
  },

  inputNetworkDevice: (networkDevice: string): ReduxAction<string> => {
    return {
      type: ActionType.IC_NetworkDevice,
      payload: networkDevice
    };
  },

  inputVipAddress: (vipAddress: string): ReduxAction<string> => {
    return {
      type: ActionType.IC_VipAddress,
      payload: vipAddress
    };
  },

  inputVipPort: (vipPort: string): ReduxAction<string> => {
    return {
      type: ActionType.IC_VipPort,
      payload: vipPort
    };
  },

  selectVipType: (vipType: string): ReduxAction<string> => {
    return {
      type: ActionType.v_IC_Vip,
      payload: vipType
    };
  },

  useGPU: (gpu: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.v_IC_Gpu,
      payload: gpu
    };
  },

  useMerticsServer: (merticsServer: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.v_IC_Mertics_server,
      payload: merticsServer
    };
  },

  useCilium: (cilium: string): ReduxAction<string> => {
    return {
      type: ActionType.v_IC_Cilium,
      payload: cilium
    };
  },

  setNetWorkMode: (networkMode: string) => {
    return {
      type: ActionType.v_IC_NetworkMode,
      payload: networkMode
    };
  },

  setAsNumber: (asNumber: string) => {
    return {
      type: ActionType.IC_AS,
      payload: asNumber
    };
  },

  setSwitchIp: (ip: string) => {
    return {
      type: ActionType.IC_SwitchIp,
      payload: ip
    };
  },

  inputGPUType: (type: string): ReduxAction<string> => {
    return {
      type: ActionType.v_IC_GpuType,
      payload: type
    };
  },

  setCidr: (cidr: string, maxClusterServiceNum: string, maxNodePodNum: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.IC_MaxClusterServiceNum,
        payload: maxClusterServiceNum
      });
      dispatch({
        type: ActionType.IC_MaxNodePodNum,
        payload: maxNodePodNum
      });
      dispatch({
        type: ActionType.IC_Cidr,
        payload: cidr
      });
    };
  },

  updateComputerList: (computerList: ICComponter[]): ReduxAction<ICComponter[]> => {
    return {
      type: ActionType.IC_ComputerList,
      payload: computerList
    };
  },
  updateComputerAction: (computerEdit: ICComponter): ReduxAction<ICComponter> => {
    return {
      type: ActionType.IC_ComputerEdit,
      payload: computerEdit
    };
  },

  setEnableContainerRuntime: (runtime: ContainerRuntimeEnum) => {
    return {
      type: ActionType.IC_EnableContainerRuntime,
      payload: runtime
    };
  },

  /** 离开创建页面，清除 Creation当中的内容 */
  clear: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.IC_Clear,
        payload: null
      });
      dispatch(createICAction.updateComputerList([]));
    };
  }
};
