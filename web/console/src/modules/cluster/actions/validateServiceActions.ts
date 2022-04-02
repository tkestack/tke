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
import { SessionAffinity } from './../constants/Config';
import * as ActionType from '../constants/ActionType';
import { RootState, PortMap, ServiceEdit, Selector } from '../models';
import { cloneDeep } from '../../common/utils';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

type GetState = () => RootState;

export const validateServiceActions = {
  /** 校验服务名称是否正确 */
  _validateServiceName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证服务名称
    if (!name) {
      status = 2;
      message = t('服务名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('服务名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('服务名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateServiceName() {
    return async (dispatch, getState: GetState) => {
      const { serviceEdit } = getState().subRoot,
        { serviceName } = serviceEdit;

      const result = await validateServiceActions._validateServiceName(serviceName);

      dispatch({
        type: ActionType.SV_ServiceName,
        payload: result
      });
    };
  },

  /**校验描述是否合法 */
  _validateServiceDesp(desp: string) {
    let status = 0,
      message = '';

    //验证服务描述
    if (desp && desp.length > 1000) {
      status = 2;
      message = t('服务描述不能超过1000个字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateServiceDesp() {
    return async (dispatch, getState: GetState) => {
      const { description } = getState().subRoot.serviceEdit;

      const result = await validateServiceActions._validateServiceDesp(description);

      dispatch({
        type: ActionType.SV_Description,
        payload: result
      });
    };
  },

  /**校验命名空间是否合法 */
  _validateNamespace(namespace: string) {
    let status = 0,
      message = '';

    // 验证命名空间的选择
    if (!namespace) {
      status = 2;
      message = t('命名空间不能为空');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNamespace() {
    return async (dispatch, getState: GetState) => {
      const { namespace } = getState().subRoot.serviceEdit;
      const result = validateServiceActions._validateNamespace(namespace);
      dispatch({
        type: ActionType.SV_Namespace,
        payload: result
      });
    };
  },

  /** 校验端口映射协议 */
  _validatePortProtocol(protocol: string, portsMap: PortMap[], communicationType: string) {
    let status = 0,
      message = '';

    if (communicationType === 'LoadBalancer' || communicationType === 'SvcLBTypeInner') {
      if (portsMap.length > 1 && portsMap.filter(p => p.protocol === protocol).length !== portsMap.length) {
        status = 2;
        message = t('协议必须相同');
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validatePortProtocol(protocol: string, portMapId: string) {
    return async (dispatch, getState: GetState) => {
      const { portsMap, communicationType } = getState().subRoot.serviceEdit;
      const newPortsMap: PortMap[] = cloneDeep(portsMap),
        portIndex = newPortsMap.findIndex(p => p.id === portMapId),
        result = validateServiceActions._validatePortProtocol(protocol, newPortsMap, communicationType);

      newPortsMap[portIndex]['v_protocol'] = result;
      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: newPortsMap
      });
    };
  },

  _validateAllProtocol(ports: PortMap[], communicationType: string) {
    let result = true;
    ports.forEach(item => {
      const temp = validateServiceActions._validatePortProtocol(item.protocol, ports, communicationType);
      result = result && temp.status === 1;
    });
    return result;
  },

  validateAllProtocol() {
    return async (dispatch, getState: GetState) => {
      const ports = getState().subRoot.serviceEdit.portsMap;
      ports.forEach(item => {
        dispatch(validateServiceActions.validatePortProtocol(item.protocol, item.id + ''));
      });
    };
  },

  /** 校验容器端口 */
  _validateTargetPort(targetPort: string, ports: PortMap[], protocol: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!targetPort) {
      status = 2;
      message = t('请输入容器端口');
    } else if (!reg.test(targetPort)) {
      status = 2;
      message = t('端口格式不正确');
    } else if (+targetPort < 1 || +targetPort > 65535) {
      status = 2;
      message = t('端口大小必须为1~65535');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateTargetPort(targetPort: string, tId: string) {
    return async (dispatch, getState: GetState) => {
      const portMaps: PortMap[] = cloneDeep(getState().subRoot.serviceEdit.portsMap),
        portIndex = portMaps.findIndex(p => p.id === tId),
        result = validateServiceActions._validateTargetPort(targetPort, portMaps, portMaps[portIndex].protocol);

      portMaps[portIndex]['v_targetPort'] = result;

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: portMaps
      });
    };
  },

  _validateAllTargetPort(ports: PortMap[]) {
    let result = true;
    ports.forEach(item => {
      const temp = validateServiceActions._validateTargetPort(item.targetPort, ports, item.protocol);
      result = result && temp.status === 1;
    });
    return result;
  },

  validateAllTargetPort() {
    return async (dispatch, getState: GetState) => {
      const ports = getState().subRoot.serviceEdit.portsMap;
      ports.forEach(p => {
        dispatch(validateServiceActions.validateTargetPort(p.targetPort, p.id + ''));
      });
    };
  },

  /**
   * 校验端口映射 - 主机端口
   */
  _validateNodePort(nodePort: string, ports: PortMap[]) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    // nodeport可以不填，不填则默认自动分配
    if (nodePort === '') {
      status = 1;
      message = '';
    } else if (!reg.test(nodePort)) {
      status = 2;
      message = t('端口格式不正确');
    } else if (+nodePort < 30000 || +nodePort > 32767) {
      status = 2;
      message = t('端口大小必须为30000~32767');
    } else if (ports.filter(p => p.nodePort === nodePort).length > 1) {
      status = 2;
      message = t('端口不可重复映射');
    } else {
      status = 1;
      message = '';
    }
    return { message, status };
  },

  validateNodePort(nodePort: string, nId: string) {
    return async (dispatch, getState: GetState) => {
      const portMaps: PortMap[] = cloneDeep(getState().subRoot.serviceEdit.portsMap),
        portIndex = portMaps.findIndex(p => p.id === nId),
        result = validateServiceActions._validateNodePort(nodePort, portMaps);

      portMaps[portIndex]['v_nodePort'] = result;
      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: portMaps
      });
    };
  },

  _validateAllNodePort(ports: PortMap[]) {
    let result = true;
    ports.forEach(item => {
      const temp = validateServiceActions._validateNodePort(item.nodePort, ports);
      result = result && temp.status === 1;
    });
    return result;
  },

  validateAllNodePort() {
    return async (dispatch, getState: GetState) => {
      const ports = getState().subRoot.serviceEdit.portsMap;
      ports.forEach(item => {
        dispatch(validateServiceActions.validateNodePort(item.nodePort, item.id + ''));
      });
    };
  },

  /**
   * 端口映射 - 服务端口
   */
  _validateServicePort(port: string, ports: PortMap[], protocol: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!port) {
      status = 2;
      message = t('请输入服务端口');
    } else if (!reg.test(port)) {
      status = 2;
      message = t('端口格式不正确');
    } else if (+port < 1 || +port > 65535) {
      status = 2;
      message = t('端口大小必须为1~65535');
    } else if (ports.filter(p => p.port === port && p.protocol === protocol).length > 1) {
      status = 2;
      message = t('端口不可重复映射');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateServicePort(port: string, sId: string) {
    return async (dispatch, getState: GetState) => {
      const portMaps: PortMap[] = cloneDeep(getState().subRoot.serviceEdit.portsMap),
        portIndex = portMaps.findIndex(p => p.id === sId),
        result = validateServiceActions._validateServicePort(port, portMaps, portMaps[portIndex].protocol);

      portMaps[portIndex]['v_port'] = result;

      dispatch({
        type: ActionType.S_UpdatePortsMap,
        payload: portMaps
      });
    };
  },

  _validateAllServicePort(ports: PortMap[]) {
    let result = true;
    ports.forEach(item => {
      const temp = validateServiceActions._validateServicePort(item.port, ports, item.protocol);
      result = result && temp.status === 1;
    });
    return result;
  },

  validateAllServicePort() {
    return async (dispatch, getState: GetState) => {
      const ports = getState().subRoot.serviceEdit.portsMap;
      ports.forEach(p => {
        dispatch(validateServiceActions.validateServicePort(p.port, p.id + ''));
      });
    };
  },

  /**
   * 校验selectors填入值的正确性
   */
  _validateSelectorContent(data: string, isKey = false) {
    let status = 0,
      message = '',
      reg = /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/;

    if (!data) {
      status = 2;
      message = t('值不能为空');
    } else if (!reg.test(data)) {
      status = 2;
      message = isKey ? t('格式不正确，只能包含小写字母、数字及 "-" "_" "/" "."') : t('格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateSelectorContent(obj: any, sId: string) {
    return async (dispatch, getState: GetState) => {
      const selectors: Selector[] = cloneDeep(getState().subRoot.serviceEdit.selector),
        sIndex = selectors.findIndex(s => s.id === sId),
        keyName = Object.keys(obj)[0],
        result = validateServiceActions._validateSelectorContent(obj[keyName], keyName === 'key');

      selectors[sIndex]['v_' + keyName] = result;
      dispatch({
        type: ActionType.S_Selector,
        payload: selectors
      });
    };
  },

  _validateAllSelectorContent(selector: Selector[]) {
    let result = true;
    selector.forEach(s => {
      const keyResult = validateServiceActions._validateSelectorContent(s.key),
        valueResult = validateServiceActions._validateSelectorContent(s.value);

      result = result && keyResult.status === 1 && valueResult.status === 1;
    });
    return result;
  },

  validateAllSelectorContent() {
    return async (dispatch, getState: GetState) => {
      const selectors = getState().subRoot.serviceEdit.selector;
      selectors.forEach(s => {
        dispatch(validateServiceActions.validateSelectorContent({ key: s.key }, s.id + ''));
        dispatch(validateServiceActions.validateSelectorContent({ value: s.value }, s.id + ''));
      });
    };
  },

  /**
   * 校验整个表单是否正确
   */
  _validateServiceEdit(serviceEdit: ServiceEdit) {
    const {
      serviceName,
      description,
      namespace,
      portsMap,
      communicationType,
      selector,
      sessionAffinity,
      sessionAffinityTimeout
    } = serviceEdit;

    let result = true;
    result =
      result &&
      validateServiceActions._validateServiceName(serviceName).status === 1 &&
      validateServiceActions._validateServiceDesp(description).status === 1 &&
      validateServiceActions._validateNamespace(namespace).status === 1 &&
      validateServiceActions._validateAllProtocol(portsMap, communicationType) &&
      validateServiceActions._validateAllServicePort(portsMap) &&
      validateServiceActions._validateAllTargetPort(portsMap);

    // 这里是自己创建Service，关联相关的selector
    if (selector.length) {
      result = result && validateServiceActions._validateAllSelectorContent(selector);
    }

    const isNodePort = communicationType === 'NodePort';

    // 只有非 clusterIP的访问类型，才有 nodePort的设置
    if (isNodePort) {
      result = result && validateServiceActions._validateAllNodePort(portsMap);
    }

    //**只有开启sessionAffinity才验证最大会话保持时间 */
    if (sessionAffinity === SessionAffinity.ClientIP) {
      result =
        result &&
        validateServiceActions._validatesessionAffinityTimeout(sessionAffinityTimeout, communicationType).status === 1;
    }

    return result;
  },

  validateServiceEdit() {
    return async (dispatch, getState: GetState) => {
      let { subRoot } = getState(),
        { serviceEdit } = subRoot,
        { communicationType, selector, sessionAffinity } = serviceEdit;

      dispatch(validateServiceActions.validateServiceName());
      dispatch(validateServiceActions.validateServiceDesp());
      dispatch(validateServiceActions.validateNamespace());
      dispatch(validateServiceActions.validateAllProtocol());
      dispatch(validateServiceActions.validateAllTargetPort());
      dispatch(validateServiceActions.validateAllServicePort());

      // 只有 selector 数组存在时，才校验
      selector.length >= 1 && dispatch(validateServiceActions.validateAllSelectorContent());

      const isNodePort = communicationType === 'NodePort';

      // 只有访问类型为 非 ClusterIP，才需要校验nodePort
      if (isNodePort) {
        dispatch(validateServiceActions.validateAllNodePort());
      }

      //**只有开启sessionAffinity才验证最大会话保持时间 */
      if (sessionAffinity === SessionAffinity.ClientIP) {
        dispatch(validateServiceActions.validatesessionAffinityTimeout());
      }
    };
  },

  /** 校验更新访问方式 */
  _validateUpdateServiceAccessEdit(serviceEdit: ServiceEdit) {
    const { portsMap, communicationType, sessionAffinity, sessionAffinityTimeout } = serviceEdit;

    let result = true;
    result =
      result &&
      validateServiceActions._validateAllProtocol(portsMap, communicationType) &&
      validateServiceActions._validateAllServicePort(portsMap) &&
      validateServiceActions._validateAllTargetPort(portsMap);

    const isNodePort = communicationType === 'NodePort';

    // 只有非 clusterIP的访问类型，才有 nodePort的设置
    if (isNodePort) {
      result = result && validateServiceActions._validateAllNodePort(portsMap);
    }

    //**只有开启sessionAffinity才验证最大会话保持时间 */
    if (sessionAffinity === SessionAffinity.ClientIP) {
      result =
        result &&
        validateServiceActions._validatesessionAffinityTimeout(sessionAffinityTimeout, communicationType).status === 1;
    }
    return result;
  },

  validateUpdateServiceAccessEdit() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { communicationType, sessionAffinity } = getState().subRoot.serviceEdit;

      dispatch(validateServiceActions.validateAllProtocol());
      dispatch(validateServiceActions.validateAllTargetPort());
      dispatch(validateServiceActions.validateAllServicePort());

      const isNodePort = communicationType === 'NodePort';

      // 只有访问类型为 非 ClusterIP，才需要校验nodePort
      if (isNodePort) {
        dispatch(validateServiceActions.validateAllNodePort());
      }

      //**只有开启sessionAffinity才验证最大会话保持时间 */
      if (sessionAffinity === SessionAffinity.ClientIP) {
        dispatch(validateServiceActions.validatesessionAffinityTimeout());
      }
    };
  },

  _validatesessionAffinityTimeout(item, communicationType) {
    let reg = /^\d+$/,
      message = '',
      status = 0;
    if (!item) {
      status = 2;
      message = t('会话保持时间不能为空');
    } else if (!reg.test(item)) {
      status = 2;
      message = t('会话保持时间格式错误');
    } else {
      if (communicationType !== 'ClusterIP' && communicationType !== 'NodePort') {
        if (item < 30 || item > 3600) {
          status = 2;
          message = t('会话保持时间范围错误');
        } else {
          status = 1;
          message = '';
        }
      } else {
        if (item <= 0 || item > 86400) {
          status = 2;
          message = t('会话保持时间范围错误');
        } else {
          status = 1;
          message = '';
        }
      }
    }
    return {
      status,
      message
    };
  },

  validatesessionAffinityTimeout() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { sessionAffinityTimeout, communicationType } = getState().subRoot.serviceEdit;
      const result = validateServiceActions._validatesessionAffinityTimeout(sessionAffinityTimeout, communicationType);
      dispatch({
        type: ActionType.SV_sessionAffinityTimeout,
        payload: result
      });
    };
  }
};
