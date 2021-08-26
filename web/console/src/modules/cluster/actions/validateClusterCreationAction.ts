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

import { ClusterCreationState } from './../models/ClusterCreationState';
import { RootState, NamespaceEdit } from '../models';
import { clusterCreationAction } from './clusterCreationAction';

type GetState = () => RootState;

export const validateClusterCreationAction = {
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
      const { name } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateClusterName(name);
      dispatch(clusterCreationAction.updateClusterCreationState({ v_name: result }));
    };
  },

  _validateApiServer(name: string) {
    let status = 0,
      message = '',
      numberReg = /^\d+$/,
      ipReg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,
      hostReg = /^([\w-]+\.)+[\w-]+(\/[\w- .\/?%&=]*)?$/;
    //验证集群名称

    if (!name) {
      status = 2;
      message = 'API Server地址不能为空';
    } else if (name.startsWith('https://')) {
      const tempName = name.substring(8);
      const tempSplit = tempName.split(':');
      let host = tempSplit[0];
      let path = '',
        port = '';
      if (host.indexOf('/') !== -1) {
        const index = host.indexOf('/');
        path = host.substring(index);
        host = host.substring(0, index);
        port = '443';
      } else {
        port = tempSplit[1] ? tempSplit[1].split('/')[0] : '443';
        if (tempSplit[1] && tempSplit[1].indexOf('/') !== -1) {
          path = tempSplit[1] ? tempSplit[1].substring(tempSplit[1].indexOf('/')) : '';
        }
      }
      if (!host) {
        status = 2;
        message = 'API Server访问地址域名不能为空';
      } else if (!ipReg.test(host) && !hostReg.test(host)) {
        status = 2;
        message = 'API Server格式不正确';
      } else {
        status = 1;
        message = '';
      }
      if (!numberReg.test(port)) {
        status = 2;
        message = '端口格式错误';
      } else if (+port < 1 || +port > 65535) {
        status = 2;
        message = '端口范围为1～65535';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 2;
      message = 'API Server访问地址，必须是https';
    }

    return { status, message };
  },

  validateApiServer() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { apiServer } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateApiServer(apiServer);

      dispatch(clusterCreationAction.updateClusterCreationState({ v_apiServer: result }));
    };
  },

  _validateCertfile(certFile: string) {
    return {
      status: 1,
      message: ''
    };
  },

  validateCertfile() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { certFile } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateCertfile(certFile);

      dispatch(clusterCreationAction.updateClusterCreationState({ v_certFile: result }));
    };
  },

  _validateToken(token: string) {
    return {
      status: 1,
      message: ''
    };
  },

  validateToken() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { token } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateToken(token);
      dispatch(clusterCreationAction.updateClusterCreationState({ v_token: result }));
    };
  },
  /** 校验clusterconnection的正确性 */
  _validateclusterCreationState(clusterCreationState: ClusterCreationState) {
    const { name, apiServer, certFile, token } = clusterCreationState;

    let result = true;

    result =
      result &&
      validateClusterCreationAction._validateClusterName(name).status === 1 &&
      validateClusterCreationAction._validateApiServer(apiServer).status === 1 &&
      validateClusterCreationAction._validateCertfile(certFile).status === 1 &&
      // validateClusterCreationAction._validatePort(port).status === 1 &&
      validateClusterCreationAction._validateToken(token).status === 1;

    return result;
  },

  validateclusterCreationState() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(validateClusterCreationAction.validateClusterName());
      dispatch(validateClusterCreationAction.validateCertfile());
      dispatch(validateClusterCreationAction.validateApiServer());
      // dispatch(validateClusterCreationAction.validatePort());
      dispatch(validateClusterCreationAction.validateToken());
    };
  }
};
