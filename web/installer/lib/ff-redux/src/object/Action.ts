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
import { FetcherState, FetchOptions, QueryState } from '../../';
import { extend } from '../../libs/qcloud-lib';
import { BaseAction, createBaseAction, DataType, RecordSet } from '../base';
import { createFFObjectActionType } from './ActionType';
import { FFObjectModel } from './Model';

const ifvisible = require('ifvisible.js');

interface PollingOptions<T, TFilter> {
  /** 拉取数据的filter */
  filter?: TFilter;

  /** 数据拉取错误后的回调 */
  onError?: (dispatch: Redux.Dispatch) => void;

  /** 重试次数的限制，默认为3 */
  retryTimes?: number;

  /** timer的延迟时间，默认为3000 */
  delayTime?: number;
}

export interface FFObjectAction<T = any, TFilter = any> extends BaseAction<T, TFilter> {
  /** 清楚选择项：包含 单项数组 和 列表数组 */
  clear?: () => void;

  polling?: ({ filter, onError, retryTimes, delayTime }: PollingOptions<T, TFilter>, visibleCheck?: boolean) => void;
  clearPolling?: () => void;
}

export interface FFObjectActionParams<T, TFilter> {
  id?: string;
  actionName: string;
  fetcher: (query: QueryState<TFilter>, getState, fetchOptions: FetchOptions, dispatch: Redux.Dispatch) => Promise<T>;
  getRecord: (getState) => FFObjectModel<T, TFilter>;
  onFinish?: (record: FetcherState<T>, dispatch: Redux.Dispatch, getState) => void;
  onClear?: (dispatch: Redux.Dispatch, getState) => void;
}

export function createFFObjectActions<T, TFilter>({
  id,
  actionName,
  fetcher,
  getRecord,
  onFinish,
  onClear
}: FFObjectActionParams<T, TFilter>): FFObjectAction<T, TFilter> {
  const ActionType = createFFObjectActionType(actionName, id);

  const baseAction = createBaseAction<T, TFilter>({
    dataType: DataType.Object,
    actionType: ActionType.Base,
    fetcher: async (getState, fetchOptions, dispatch) => {
      let response = await fetcher(getRecord(getState).query, getState, fetchOptions, dispatch);
      return response;
    },
    finish: (dispatch, getState) => {
      let model = getRecord(getState);
      if (model.object.error) {
        onPollingErr(dispatch);
      } else {
        onFinish && onFinish.apply(actions, [model.object, dispatch, getState]);
        onPollingSucc();
      }
    }
  });

  let _polling = {
    timer: null,
    errorTimes: 0,
    retryTimes: 3,
    delayTime: 3000,
    onError: null
  };

  let clearPolling = () => {
    clearTimeout(_polling.timer);
    _polling = {
      timer: null,
      errorTimes: 0,
      retryTimes: 3,
      delayTime: 3000,
      onError: null
    };
  };

  let onPollingSucc = () => {
    //如果返回成功，就重置错误次数
    _polling.errorTimes = 0;
  };

  let onPollingErr = (dispatch: Redux.Dispatch) => {
    _polling.errorTimes++;
    if (_polling.errorTimes >= _polling.retryTimes) {
      _polling.onError && _polling.onError(dispatch);
      clearPolling();
    }
  };

  let doPolling = (filter: TFilter, dispatch, visibleCheck: boolean) => {
    if (visibleCheck) {
      if (ifvisible.now()) {
        dispatch(actions.applyPolling(filter));
      } else {
        // do not polling
      }
    } else {
      dispatch(actions.applyPolling(filter));
    }
    _polling.timer = setTimeout(() => {
      doPolling(filter, dispatch, visibleCheck);
    }, _polling.delayTime);
  };

  const restActions = {
    polling: ({ filter, onError, retryTimes = 3, delayTime = 3000 }, visibleCheck = false) => {
      return async (dispatch, getState) => {
        // 这里在多次调用polling的时候，需要先清除之前的timer
        if (_polling.timer) {
          clearPolling();
        }
        _polling = {
          timer: null,
          onError,
          retryTimes,
          delayTime,
          errorTimes: 0
        };
        doPolling(filter, dispatch, visibleCheck);
      };
    },
    clearPolling: () => {
      return async (dispatch: Redux.Dispatch, getState) => {
        clearPolling();
      };
    },
    clear: () => {
      return async (dispatch: Redux.Dispatch, getState) => {
        dispatch({
          type: ActionType.Clear
        });
        if (actions.clearFetch) {
          dispatch(actions.clearFetch());
        }
        onClear && onClear.apply(actions, [dispatch, getState]);
      };
    }
  };

  const actions = extend({}, baseAction, restActions);
  return actions;
}
