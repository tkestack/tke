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
import { t, Trans } from '@tencent/tea-app/lib/i18n';

type GetState = () => RootState;

export const validateNamespaceActions = {
  /**
   * 校验namespace名称是否正确
   */
  _validateNamespaceName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证ingress名称
    if (!name) {
      status = 2;
      message = t('Namespace名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('Namespace名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Namespace名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNamespaceName() {
    return async (dispatch, getState: GetState) => {
      let { namespaceEdit } = getState().subRoot,
        { name } = namespaceEdit;

      const result = await validateNamespaceActions._validateNamespaceName(name);

      dispatch({
        type: ActionType.NV_Name,
        payload: result
      });
    };
  },

  /** 校验描述是否正确 */
  _validateNamespaceDesp(desp: string) {
    let status = 0,
      message = '';

    // 验证ingress描述
    if (desp && desp.length > 1000) {
      status = 2;
      message = t('Namespace描述不能超过1000个字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNamespaceDesp() {
    return async (dispatch, getState: GetState) => {
      let { description } = getState().subRoot.namespaceEdit;

      const result = await validateNamespaceActions._validateNamespaceDesp(description);

      dispatch({
        type: ActionType.NV_Description,
        payload: result
      });
    };
  },

  /** 校验namespaceEdit的正确性 */
  _validateNamespaceEdit(namespaceEdit: NamespaceEdit) {
    let { name, description } = namespaceEdit;

    let result = true;

    result =
      result &&
      validateNamespaceActions._validateNamespaceName(name).status === 1 &&
      validateNamespaceActions._validateNamespaceDesp(description).status === 1;

    return result;
  },

  validateNamespaceEdit() {
    return async (dispatch, getState: GetState) => {
      dispatch(validateNamespaceActions.validateNamespaceName());
      dispatch(validateNamespaceActions.validateNamespaceDesp());
    };
  }
};
