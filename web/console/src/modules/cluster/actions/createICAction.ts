import { deepClone, ReduxAction, uuid } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
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
      let response = await WebAPI.fetchCreateICK8sVersion();
      dispatch({
        type: ActionType.IC_FetchK8SVersion,
        payload: response
      });
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
