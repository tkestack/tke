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

import * as ActionType from '../constants/ActionType';
import { RootState, NamespaceEdit } from '../models';
import validatorJs from 'validator';

type GetState = () => RootState;

export const validateCreateICAction = {
  /**
   * 校验cluster名称是否正确
   */
  _validateClusterName(name: string) {
    let status = 0,
      message = '';

    //验证集群名称
    if (!name) {
      status = 2;
      message = '集群名称不能为空';
    } else if (name.length > 60) {
      status = 2;
      message = '集群名称不能超过60个字符';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateClusterName() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { name } = getState().createIC;
      const result = await validateCreateICAction._validateClusterName(name);
      dispatch({
        type: ActionType.v_IC_Name,
        payload: result
      });
    };
  },

  _validateNetworkDevice(networkDevice: string) {
    let status = 0,
      message = '';
    const reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/;
    //验证集群网卡模式
    if (!networkDevice) {
      status = 2;
      message = '集群网卡名称不能为空';
    } else if (networkDevice.length > 60) {
      status = 2;
      message = '集群网卡不能超过60个字符';
    } else if (!reg.test(networkDevice)) {
      status = 2;
      message = '集群网卡名称格式不对';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  validateNetworkDevice() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { networkDevice } = getState().createIC;
      const result = await validateCreateICAction._validateNetworkDevice(networkDevice);
      dispatch({
        type: ActionType.v_IC_NetworkDevice,
        payload: result
      });
    };
  },
  /** 校验port是否正确 */
  _validatePort(port: string) {
    let status = 0,
      message = '';

    if (!port) {
      status = 2;
      message = '端口不能为空';
    } else if (isNaN(+port)) {
      status = 2;
      message = '端口格式错误';
    } else if (+port < 1 || +port > 65535) {
      status = 2;
      message = '端口范围为1～65535';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validatePort() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { vipPort } = getState().createIC;
      const result = validateCreateICAction._validatePort(vipPort);
      dispatch({
        type: ActionType.v_IC_VipPort,
        payload: result
      });
    };
  },

  /**
   * 校验apiserver地址是否正确
   */
  _validateVIPServer(name: string) {
    let status = 0,
      message = '',
      ipReg =
        /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
    // hostReg = /^([\w-]+\.)+[\w-]+(\/[\w- .\/?%&=]*)?$/;
    //验证集群名称

    if (!name) {
      status = 2;
      message = 'VIP不能为空';
    } else if (!ipReg.test(name)) {
      status = 2;
      message = 'VIP格式不正确';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateVIPServer() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { vipAddress } = getState().createIC;
      const result = await validateCreateICAction._validateVIPServer(vipAddress);

      dispatch({
        type: ActionType.v_IC_VipAddress,
        payload: result
      });
    };
  },

  _validateAsNumber(asNumber: string) {
    let status = 0,
      message = '';

    if (!asNumber) {
      status = 2;
      message = '自治系统号不能为空';
    } else if (!validatorJs.isInt(asNumber, { min: 0 })) {
      status = 2;
      message = '自治系统号必须为正整数';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  validateAsNumber() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { cilium, networkMode, asNumber } = getState().createIC;
      let result = {
        status: 0,
        message: ''
      };

      if (cilium === 'Cilium' && networkMode === 'underlay') {
        result = await validateCreateICAction._validateAsNumber(asNumber);
      }

      dispatch({
        type: ActionType.v_IC_AS,
        payload: result
      });
    };
  },

  _validateSwitchIp(ip: string) {
    let status = 0,
      message = '';

    if (!ip) {
      status = 2;
      message = '交换机IP不能为空';
    } else if (!validatorJs.isIP(ip)) {
      status = 2;
      message = '交换机IP格式不正确';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  validateSwitchIp() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { cilium, networkMode, switchIp } = getState().createIC;
      let result = {
        status: 0,
        message: ''
      };

      if (cilium === 'Cilium' && networkMode === 'underlay') {
        result = await validateCreateICAction._validateSwitchIp(switchIp);
      }

      dispatch({
        type: ActionType.v_IC_SwitchIp,
        payload: result
      });
    };
  }
};
