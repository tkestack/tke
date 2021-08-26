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

import { RootState, PeEdit, AddonEdit } from '../models';
import * as ActionType from '../constants//ActionType';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { AddonNameEnum } from '../constants/Config';

type GetState = () => RootState;

export const validatorActions = {
  /** ================================ 事件持久化的相关校验 ============================= */
  /** 校验当前的es地址是否正确 */
  _validateEsAddress(address: string) {
    let status = 0,
      message = '',
      hostReg = /^((http|https):\/\/)((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
      domainReg = /^((http|https):\/\/)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]:([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$/;

    if (!address) {
      status = 2;
      message = t('Elasticsearch地址不能为空');
    } else if (!hostReg.test(address) && !domainReg.test(address)) {
      status = 2;
      message = t('Elasticsearch地址格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateEsAddress() {
    return async (dispatch, getState: GetState) => {
      const { esAddress } = getState().editAddon.peEdit;

      const result = validatorActions._validateEsAddress(esAddress);
      dispatch({
        type: ActionType.V_EsAddress,
        payload: result
      });
    };
  },

  /** 校验当前的索引名是否正确 */
  _validateIndexName(indexName: string) {
    let status = 0,
      message = '',
      reg = /^[a-z][0-9a-z_+-]*$/;
    if (!indexName) {
      status = 2;
      message = t('索引名不能为空');
    } else if (!reg.test(indexName)) {
      status = 2;
      message = t('索引名格式不正确');
    } else if (indexName.length > 60) {
      status = 2;
      message = t('索引名不能超过60个字符');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateIndexName() {
    return async (dispatch, getState: GetState) => {
      const { indexName } = getState().editAddon.peEdit;
      const result = validatorActions._validateIndexName(indexName);
      dispatch({
        type: ActionType.V_IndexName,
        payload: result
      });
    };
  },
  /** ================================ 事件持久化的相关校验 ============================= */

  /** 创建addon的校验 */
  _validateAddonEdit(addonEdit: AddonEdit) {
    let result = true;
    const { addonName, peEdit } = addonEdit;

    result = result && addonName !== '';

    if (addonName === AddonNameEnum.PersistentEvent) {
      result =
        result &&
        validatorActions._validateEsAddress(peEdit.esAddress).status === 1 &&
        validatorActions._validateIndexName(peEdit.indexName).status === 1;
    }

    return result;
  },

  validateAddonEdit() {
    return async (dispatch, getState: GetState) => {
      let { editAddon } = getState(),
        { peEdit, addonName } = editAddon;

      if (addonName === AddonNameEnum.PersistentEvent) {
        dispatch(validatorActions.validateEsAddress());
        dispatch(validatorActions.validateIndexName());
      }
    };
  }
};
