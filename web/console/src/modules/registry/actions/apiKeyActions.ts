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

import {
    createFFListActions, extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import { InitApiKey } from '../constants/Config';
import { ApiKey, ApiKeyFilter, RootState } from '../models';
import { ApiKeyCreation } from '../models/ApiKey';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

const FFModelApiKeyActions = createFFListActions<ApiKey, ApiKeyFilter>({
  actionName: 'apiKey',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchApiKeyList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().apiKey;
  }
});

const restActions = {
  /** 创建 ApiKey */
  createApiKey: generateWorkflowActionCreator<ApiKeyCreation, void>({
    actionType: ActionType.CreateApiKey,
    workflowStateLocator: (state: RootState) => state.createApiKey,
    operationExecutor: WebAPI.createApiKey,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createApiKey, route } = getState();
        if (isSuccessWorkflow(createApiKey)) {
          dispatch(restActions.createApiKey.reset());
          dispatch(restActions.clearEdition());
          dispatch(apiKeyActions.fetch());
          let urlParams = router.resolve(route);
          router.navigate(Object.assign({}, urlParams, { sub: 'apikey', mode: 'list' }), {});
        }
      }
    }
  }),

  /** 删除 ApiKey */
  deleteApiKey: generateWorkflowActionCreator<ApiKey, void>({
    actionType: ActionType.DeleteApiKey,
    workflowStateLocator: (state: RootState) => state.deleteApiKey,
    operationExecutor: WebAPI.deleteApiKey,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteApiKey, route } = getState();
        if (isSuccessWorkflow(deleteApiKey)) {
          dispatch(restActions.deleteApiKey.reset());
          dispatch(apiKeyActions.fetch());
        }
      }
    }
  }),

  /** enable/disable ApiKey */
  toggleKeyStatus: generateWorkflowActionCreator<ApiKey, void>({
    actionType: ActionType.ToggleKeyStatus,
    workflowStateLocator: (state: RootState) => state.toggleKeyStatus,
    operationExecutor: WebAPI.toggleKeyStatus,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { toggleKeyStatus, route } = getState();
        if (isSuccessWorkflow(toggleKeyStatus)) {
          dispatch(restActions.toggleKeyStatus.reset());
          dispatch(apiKeyActions.fetch());
        }
      }
    }
  }),

  /** --begin编辑action */
  inputApiKeyDesc: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateApiKeyCreation,
        payload: Object.assign({}, getState().apiKeyCreation, { description: value })
      });
    };
  },

  inputApiKeyExpire: (value: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateApiKeyCreation,
        payload: Object.assign({}, getState().apiKeyCreation, { expire: value })
      });
      dispatch(apiKeyActions.validateApiKeyExpire(value));
    };
  },

  selectApiKeyUnit: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateApiKeyCreation,
        payload: Object.assign({}, getState().apiKeyCreation, { unit: value })
      });
    };
  },

  validateApiKeyExpire(value: number) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = apiKeyActions._validateApiKeyExpires(value);
      dispatch({
        type: ActionType.UpdateApiKeyCreation,
        payload: Object.assign({}, getState().apiKeyCreation, { v_expire: result })
      });
    };
  },

  _validateApiKeyExpires(expires: number) {
    let reg = /^\d+?$/,
      status = 0,
      message = '';

    if (isNaN(expires)) {
      status = 2;
      message = t('只能输入正整数');
    } else if (!expires) {
      status = 1;
      message = '';
    } else if (!reg.test(expires + '')) {
      status = 2;
      message = t('只能输入正整数');
    } else if (expires <= 0) {
      status = 2;
      message = t('过期时间必须大于 0');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateApiKeyCreation,
        payload: InitApiKey
      });
    };
  }
  /** --end编辑action */
};

export const apiKeyActions = extend({}, FFModelApiKeyActions, restActions);
