import { extend, deepClone, uuid } from '@tencent/qcloud-lib';
import {
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  createFFListActions
} from '@tencent/ff-redux';
import { RootState, ApiKey, ApiKeyFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import { InitApiKey } from '../constants/Config';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ApiKeyCreation } from '../models/ApiKey';

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
