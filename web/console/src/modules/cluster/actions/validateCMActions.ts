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
import * as ActionType from '../constants/ActionType';
import { RootState, ConfigMapEdit } from '../models';
import { cloneDeep } from '../../common/utils';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

type GetState = () => RootState;

export const validateCMActions = {
  /** 校验名称是否正确 */
  _validateCMName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证服务名称
    if (!name) {
      status = 2;
      message = t('ConfigMap名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('ConfigMap不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('ConfigMap名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateCMName() {
    return async (dispatch, getState: GetState) => {
      let { name } = getState().subRoot.cmEdit;
      const result = validateCMActions._validateCMName(name);
      dispatch({
        type: ActionType.V_CM_Name,
        payload: result
      });
    };
  },

  /** 校验变量名称是否正确 */
  _validateVariableKey(name: string) {
    let reg = /^([A-Za-z0-9][-A-Za-z0-9_\.]*)?[A-Za-z0-9]$/,
      status = 0,
      message = '';

    // 验证服务名称
    if (!name) {
      status = 2;
      message = t('变量名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('变量名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('变量名称名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateVariableKey() {
    return async (dispatch, getState: GetState) => {
      let variables = cloneDeep(getState().subRoot.cmEdit.variables);

      getState().subRoot.cmEdit.variables.forEach((v, index) => {
        variables[index].v_key = validateCMActions._validateVariableKey(v.key);
      });

      dispatch({
        type: ActionType.CM_ValidateVariable,
        payload: variables
      });
    };
  },

  /** 校验整个表单 */
  _validateCMEdit(cmEdit: ConfigMapEdit) {
    let result = true;
    result = result && validateCMActions._validateCMName(cmEdit.name).status === 1;

    cmEdit.variables.forEach((v, index) => {
      result = result && validateCMActions._validateVariableKey(v.key).status === 1;
    });

    return result;
  },

  validateCMEdit() {
    return async (dispatch, getState: GetState) => {
      dispatch(validateCMActions.validateCMName());
      dispatch(validateCMActions.validateVariableKey());
    };
  }
};
