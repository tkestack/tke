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
import { KeyValue } from 'src/modules/common';
import { deepClone, FFListModel } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { LbcfEdit, RootState } from '../models';
import { Selector } from '../models/ServiceEdit';
import { BackendType } from '../constants/Config';

type GetState = () => RootState;

export const validateLbcfActions = {
  /** 校验名称是否正确  复用*/
  _validateLbcfName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证服务名称
    if (!name) {
      status = 2;
      message = t('名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateLbcfName() {
    return async (dispatch, getState: GetState) => {
      let { name } = getState().subRoot.lbcfEdit;
      const result = validateLbcfActions._validateLbcfName(name);
      dispatch({
        type: ActionType.V_Gate_Name,
        payload: result
      });
    };
  },
  _validateLbcfNamespace(namespace: string) {
    let status = 0,
      message = '';
    if (!namespace) {
      status = 2;
      message = t('负载均衡名称不能为空');
    } else {
      status = 1;
    }
    return {
      status,
      message
    };
  },
  validateLbcfNamespace() {
    return async (dispatch, getState: GetState) => {
      let { namespace } = getState().subRoot.lbcfEdit;
      const result = validateLbcfActions._validateLbcfName(namespace);
      dispatch({
        type: ActionType.V_Gate_Namespace,
        payload: result
      });
    };
  },

  _validateLbcfClbSelection(clbSelection: string) {
    let status = 0,
      message = '';
    if (!clbSelection) {
      status = 2;
      message = t('已有LB不能为空');
    } else {
      status = 1;
    }
    return {
      status,
      message
    };
  },

  _validateLbcfConfig(kvs: KeyValue[]) {
    let status = 1,
      message = '';
    kvs.forEach(kv => {
      if (kv.key === '' || kv.value === '') {
        status = 2;
        message = t('该配置未填写完');
      }
    });
    return {
      status,
      message
    };
  },

  validateLbcfConfig() {
    return async (dispatch, getState: GetState) => {
      let { config } = getState().subRoot.lbcfEdit;
      const result = validateLbcfActions._validateLbcfConfig(config);
      dispatch({
        type: ActionType.V_Lbcf_Config,
        payload: result
      });
    };
  },

  validateLbcfArgs() {
    return async (dispatch, getState: GetState) => {
      let { args } = getState().subRoot.lbcfEdit;
      const result = validateLbcfActions._validateLbcfConfig(args);
      dispatch({
        type: ActionType.V_Lbcf_Args,
        payload: result
      });
    };
  },

  _validateLbcfDriver(driver: FFListModel) {
    let status = 0,
      message = '';
    if (!driver.selection) {
      status = 2;
      message = 'Driver不能为空';
    } else {
      status = 1;
      message = '';
    }
    return {
      status,
      message
    };
  },
  validateLbcfDriver() {
    return async (dispatch, getState: GetState) => {
      let { driver } = getState().subRoot.lbcfEdit;
      const result = validateLbcfActions._validateLbcfDriver(driver);
      dispatch({
        type: ActionType.V_Lbcf_Driver,
        payload: result
      });
    };
  },

  // validateLbcfClbSelection() {
  //   return async (dispatch, getState: GetState) => {
  //     let { clbSelection } = getState().subRoot.lbcfEdit;
  //     const result = validateLbcfActions._validateLbcfClbSelection(clbSelection);
  //     dispatch({
  //       type: ActionType.V_GLB_SelectClb,
  //       payload: result
  //     });
  //   };
  // },

  /** 校验整个表单 */
  _validateLbcfEdit(lbcfEdit: LbcfEdit) {
    let result = true;
    result =
      result &&
      validateLbcfActions._validateLbcfName(lbcfEdit.name).status === 1 &&
      validateLbcfActions._validateLbcfNamespace(lbcfEdit.namespace).status === 1 &&
      validateLbcfActions._validateLbcfConfig(lbcfEdit.config).status === 1 &&
      validateLbcfActions._validateLbcfConfig(lbcfEdit.args).status === 1 &&
      validateLbcfActions._validateLbcfDriver(lbcfEdit.driver).status === 1;

    // if (lbcfEdit.createLbWay === 'existed') {
    //   result = result && validateLbcfActions._validateLbcfClbSelection(lbcfEdit.clbSelection).status === 1;
    // }

    return result;
  },

  validateLbcfEdit() {
    return async (dispatch, getState: GetState) => {
      dispatch(validateLbcfActions.validateLbcfName());
      dispatch(validateLbcfActions.validateLbcfNamespace());
      dispatch(validateLbcfActions.validateLbcfConfig());
      dispatch(validateLbcfActions.validateLbcfArgs());
      dispatch(validateLbcfActions.validateLbcfDriver());
      // if (getState().subRoot.lbcfEdit.createLbWay === 'existed') {
      //   dispatch(validateLbcfActions.validateLbcfClbSelection());
      // }
    };
  },

  /**backGroup */

  validateLbcfBackGroupName(backGroupId: string, value: string) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.v_name = validateLbcfActions._validateLbcfName(value);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  validateLbcfBackGroupServiceName(backGroupId: string, value: string) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.v_serviceName = validateLbcfActions._validateLbcfName(value);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  validatePort(backGroupId: string, id: string, value: string) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();
      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { ports } = backGroupEdition;
      let index = ports.findIndex(item => item.id === id);
      ports[index].v_portNumber = validateLbcfActions._validatePort(value);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  _validatePort(port: string) {
    let reg = /^\d+?$/,
      status = 0,
      message = '';

    // 验证内存限制
    if (isNaN(+port)) {
      status = 2;
      message = t('只能输入正整数');
    } else if (port === '') {
      status = 2;
      message = t('端口不能为空');
    } else if (!reg.test(port + '')) {
      status = 2;
      message = t('只能输入正整数');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  /**
   * 校验selectors填入值的正确性
   */
  _validateLabelContent(data: string, isKey: boolean = false) {
    let status = 0,
      message = '',
      reg = /^([A-Za-z0-9][-A-Za-z0-9_]*)?[A-Za-z0-9]$/;

    if (!data) {
      status = 2;
      message = t('值不能为空');
    } else if (!reg.test(data)) {
      status = 2;
      message = isKey
        ? t('格式不正确，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')
        : t('格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateLabelContent(backGroupId: string, id: string, obj: any) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { labels } = backGroupEdition;
      let index = labels.findIndex(item => item.id === id),
        keyName = Object.keys(obj)[0];
      labels[index]['v_' + keyName] = validateLbcfActions._validateLabelContent(obj[keyName], keyName === 'key');

      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  validateAddress(backGroupId: string, id: string) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { staticAddress } = backGroupEdition;
      let index = staticAddress.findIndex(item => item.id === id);
      staticAddress[index].v_value = validateLbcfActions._validateStaticAddress(staticAddress[index].value);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  _validateStaticAddress(data: string) {
    let status = 0,
      message = '';

    if (!data) {
      status = 2;
      message = t('值不能为空');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validatePodName(backGroupId: string, id: string) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { byName } = backGroupEdition;
      let index = byName.findIndex(item => item.id === id);
      byName[index].v_value = validateLbcfActions._validateLbcfName(byName[index].value);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  /** 校验整个表单 */
  _validateGameBGEdit(lbcfEdit: LbcfEdit) {
    let result = true;
    lbcfEdit.lbcfBackGroupEditions.forEach(item => {
      let { ports, labels, name, backgroupType, staticAddress, serviceName, byName } = item;
      result = result && validateLbcfActions._validateLbcfName(name).status === 1;
      if (backgroupType === BackendType.Static) {
        staticAddress.forEach(address => {
          result = result && validateLbcfActions._validateStaticAddress(address.value).status === 1;
        });
      } else if (backgroupType === BackendType.Service) {
        result = result && validateLbcfActions._validateLbcfName(serviceName).status === 1;
        ports.forEach(port => {
          result = result && validateLbcfActions._validatePort(port.portNumber).status === 1;
        });
        labels.forEach(label => {
          result =
            result &&
            validateLbcfActions._validateLabelContent(label.key, true).status === 1 &&
            validateLbcfActions._validateLabelContent(label.value, false).status === 1;
        });
      } else {
        byName.forEach(name => {
          result = result && validateLbcfActions._validateLbcfName(name.value).status === 1;
        });
        ports.forEach(port => {
          result = result && validateLbcfActions._validatePort(port.portNumber).status === 1;
        });
        labels.forEach(label => {
          result =
            result &&
            validateLbcfActions._validateLabelContent(label.key, true).status === 1 &&
            validateLbcfActions._validateLabelContent(label.value, false).status === 1;
        });
      }
    });
    return result;
  },

  validateGameBGEdit() {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();
      lbcfBackGroupEditions.forEach(item => {
        let { ports, labels, id, name, backgroupType, serviceName, staticAddress, byName } = item;
        dispatch(validateLbcfActions.validateLbcfBackGroupName(id + '', name));

        if (backgroupType === BackendType.Static) {
          staticAddress.forEach(address => {
            dispatch(validateLbcfActions.validateAddress(id + '', address.id + ''));
          });
        } else if (backgroupType === BackendType.Service) {
          dispatch(validateLbcfActions.validateLbcfBackGroupServiceName(id + '', serviceName));
          ports.forEach(port => {
            dispatch(validateLbcfActions.validatePort(id + '', port.id + '', port.portNumber));
          });
          labels.forEach(label => {
            dispatch(validateLbcfActions.validateLabelContent(id + '', label.id + '', { key: label.key }));
            dispatch(validateLbcfActions.validateLabelContent(id + '', label.id + '', { value: label.value }));
          });
        } else {
          byName.forEach(name => {
            dispatch(validateLbcfActions.validatePodName(id + '', name.id + ''));
          });
          ports.forEach(port => {
            dispatch(validateLbcfActions.validatePort(id + '', port.id + '', port.portNumber));
          });
          labels.forEach(label => {
            dispatch(validateLbcfActions.validateLabelContent(id + '', label.id + '', { key: label.key }));
            dispatch(validateLbcfActions.validateLabelContent(id + '', label.id + '', { value: label.value }));
          });
        }
      });
    };
  }
};
