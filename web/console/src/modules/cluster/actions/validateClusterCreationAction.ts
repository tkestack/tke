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
      let { name } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateClusterName(name);
      dispatch(clusterCreationAction.updateClusterCreationState({ v_name: result }));
    };
  },

  /** 校验port是否正确 */
  _validatePort(port: string) {
    let status = 0,
      message = '';

    if (!port) {
      status = 2;
      message = 'port端口不能为空';
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
      let { port } = getState().clusterCreationState;
      let result = validateClusterCreationAction._validatePort(port);
      dispatch(clusterCreationAction.updateClusterCreationState({ v_port: result }));
    };
  },

  /**
   * 校验apiserver地址是否正确
   */
  _validateApiServer(name: string) {
    let status = 0,
      message = '',
      ipReg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/,
      hostReg = /^([\w-]+\.)+[\w-]+(\/[\w- .\/?%&=]*)?$/;
    //验证集群名称

    if (!name) {
      status = 2;
      message = 'API Server地址不能为空';
    } else if (!ipReg.test(name) && !hostReg.test(name)) {
      status = 2;
      message = 'API Server格式不正确';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateApiServer() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { apiServer } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateApiServer(apiServer);

      dispatch(clusterCreationAction.updateClusterCreationState({ v_apiServer: result }));
    };
  },

  _validateCertfile(certFile: string) {
    let status = 0,
      message = '';

    //验证集群名称
    if (!certFile) {
      status = 2;
      message = '证书不能为空';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateCertfile() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { certFile } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateCertfile(certFile);

      dispatch(clusterCreationAction.updateClusterCreationState({ v_certFile: result }));
    };
  },

  _validateToken(token: string) {
    let status = 0,
      message = '';

    //验证集群名称
    if (!token) {
      status = 2;
      message = 'token不能为空';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateToken() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { token } = getState().clusterCreationState;
      const result = await validateClusterCreationAction._validateToken(token);
      dispatch(clusterCreationAction.updateClusterCreationState({ v_token: result }));
    };
  },
  /** 校验clusterconnection的正确性 */
  _validateclusterCreationState(clusterCreationState: ClusterCreationState) {
    let { name, apiServer, certFile, port, token } = clusterCreationState;

    let result = true;

    result =
      result &&
      validateClusterCreationAction._validateClusterName(name).status === 1 &&
      validateClusterCreationAction._validateApiServer(apiServer).status === 1 &&
      validateClusterCreationAction._validateCertfile(certFile).status === 1 &&
      validateClusterCreationAction._validatePort(port).status === 1 &&
      validateClusterCreationAction._validateToken(token).status === 1;

    return result;
  },

  validateclusterCreationState() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(validateClusterCreationAction.validateClusterName());
      dispatch(validateClusterCreationAction.validateCertfile());
      dispatch(validateClusterCreationAction.validateApiServer());
      dispatch(validateClusterCreationAction.validatePort());
      dispatch(validateClusterCreationAction.validateToken());
    };
  }
};
