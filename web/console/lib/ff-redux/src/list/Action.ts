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
import { extend } from '../../libs/qcloud-lib';
import { BaseAction, createBaseAction, FetcherState, FetchOptions, FetchState, QueryState, RecordSet } from '../base';
import { createFFListActionType } from './ActionType';
import { FFListModel } from './Model';

const ifvisible = require('ifvisible.js');
function getFieldValue(record, field: String | Function) {
  if (typeof field === 'function') {
    return field(record);
  } else {
    return record[field as string];
  }
}
interface PollingOptions<T, TFilter> {
  /** 拉取数据的filter */
  filter?: TFilter;
  /** 数据拉取错误后的回调 */
  onError?: (dispatch: Redux.Dispatch) => void;

  /** 重试次数的限制，默认为3 */
  retryTimes?: number;

  /** timer的延迟时间，默认为3000 */
  delayTime?: number;

  /** 是否需要开启可视范围内 不进行数据的拉取 */
  visibleCheck?: boolean;
}

export interface FFListAction<T = any, TFilter = any, TSFilter = any> extends BaseAction<T, TFilter, TSFilter> {
  /** 选择相关操作，列表单项的值 */
  select?: (t: T) => void;
  setInitValue?: (value: string | number) => void;
  selectByValue?: (value: string | number) => void;

  /** 选择相关操作，数组项 */
  selects?: (ts: T[]) => void;
  setInitValues?: (values: string[] | number[]) => void;
  selectsByValue?: (values: string[] | number[]) => void;

  /** 清楚选择项：包含 单项数组 和 列表数组 */
  clear?: () => void;
  clearSelection?: () => void;

  polling?: ({ filter, onError, retryTimes, delayTime, visibleCheck }: PollingOptions<T, TFilter>) => void;
  clearPolling?: () => void;
  startPolling?: ({ onError, retryTimes, delayTime, visibleCheck }: PollingOptions<T, TFilter>) => {};
}

export interface FFListActionParams<T, TFilter, ExtendParamsT = any, TSFilter = any> {
  id?: string;
  actionName: string;
  fetcher: (
    query: QueryState<TFilter, TSFilter>,
    getState,
    fetchOptions: FetchOptions,
    dispatch: Redux.Dispatch
  ) => Promise<RecordSet<T, ExtendParamsT>>;
  getRecord: (getState) => FFListModel<T, TFilter>;
  selectFirst?: boolean;
  keepLastSelection?: boolean;
  onFinish?: (record: FetcherState<RecordSet<T, ExtendParamsT>>, dispatch: Redux.Dispatch, getState) => void;
  onSelect?: (t: T, dispatch: Redux.Dispatch, getState) => void;
  onSelects?: (ts: T[], dispatch: Redux.Dispatch, getState) => void;
  onClear?: (dispatch: Redux.Dispatch, getState) => void;
  onClearSelection?: (dispatch: Redux.Dispatch, getState) => void;
}

export function createFFListActions<T, TFilter, ExtendParamsT = any, TSFilter = any>({
  id,
  actionName,
  fetcher,
  getRecord,
  keepLastSelection = true,
  selectFirst,
  onFinish,
  onSelect,
  onSelects,
  onClear,
  onClearSelection
}: FFListActionParams<T, TFilter, ExtendParamsT, TSFilter>): FFListAction<T, TFilter, TSFilter> {
  const ActionType = createFFListActionType(actionName, id);

  const baseActions = createBaseAction<RecordSet<T, ExtendParamsT>, TFilter, TSFilter>({
    actionType: ActionType.Base,
    fetcher: async (getState, fetchOptions, dispatch) => {
      let response = await fetcher(getRecord(getState).query, getState, fetchOptions, dispatch);

      if (fetchOptions && fetchOptions.fetchAll) {
        let query = Object.assign({}, getRecord(getState).query);
        let times = 0,
          maxTimes = fetchOptions.maxFetchTimes || 5;
        while (response.records.length < response.recordCount && times < maxTimes) {
          query.paging.pageIndex++;
          times++;
          let pResponse = await fetcher(query, getState, fetchOptions, dispatch);
          if (pResponse.records.length) {
            response.records = response.records.concat(pResponse.records);
          }
        }
      }

      return response;
    },
    getFetchReducer: getState => {
      return getRecord(getState).list;
    },
    getQueryReducer: getState => {
      return getRecord(getState).query;
    },
    finish: (dispatch, getState) => {
      let model = getRecord(getState);
      if (model.list.error) {
        onPollingErr(dispatch);
      } else {
        if (model.initValue !== null) {
          dispatch(actions.setInitValue(model.initValue));
        } else if (model.initValues && model.initValues.length) {
          dispatch(actions.setInitValues(model.initValues));
        } else {
          if (
            keepLastSelection &&
            ((model.selection &&
              model.list.data.recordCount &&
              model.list.data.records.filter(record => {
                return getFieldValue(record, model.valueField) === getFieldValue(model.selection, model.valueField);
              }).length) ||
              (model.selections && model.selections.length))
          ) {
            //如果想要保留上一次的值，并且上一次的值还在列表里，就停留原来的selection
          } else if (selectFirst && model.list.data.recordCount) {
            //如果要求选第一项
            dispatch(actions.select(model.list.data.records[0]));
          } else if (model.selection !== null || (model.selections && model.selections.length !== 0)) {
            //否则把之前的selection清除
            dispatch(actions.clearSelection());
          }
        }

        onFinish && onFinish.apply(actions, [model.list, dispatch, getState]);
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
    setInitValue: (value: string | number) => {
      return async (dispatch, getState) => {
        let model = getRecord(getState);
        if (model.list.fetchState === FetchState.Ready && model.list.fetched) {
          //如果已经加载完成
          if (model.valueField) {
            let record = model.list.data.records.find(r => {
              return getFieldValue(r, model.valueField) === value;
            });
            if (record) {
              dispatch(restActions.select(record));
            } else if (selectFirst) {
              dispatch(restActions.select(model.list.data.records[0]));
            }
          } else {
            if (selectFirst) {
              dispatch(restActions.select(model.list.data.records[0]));
            }
          }
          dispatch({
            type: ActionType.InitValue,
            payload: null
          });
        } else {
          dispatch({
            type: ActionType.InitValue,
            payload: value
          });
        }
      };
    },
    select: (record: T) => {
      return async (dispatch, getState) => {
        dispatch({
          type: ActionType.Selection,
          payload: record
        });

        onSelect && onSelect.apply(actions, [record, dispatch, getState]);
      };
    },
    selectByValue: (value: string | number) => {
      return async (dispatch, getState) => {
        let model = getRecord(getState);
        if (model.valueField) {
          let record = model.list.data.records.find(r => {
            return getFieldValue(r, model.valueField) === value;
          });
          if (record) {
            dispatch(actions.select(record));
          }
        }
      };
    },
    setInitValues: (values: string[] | number[]) => {
      return async (dispatch, getState) => {
        let model = getRecord(getState);
        if (model.list.fetchState === FetchState.Ready && model.list.fetched) {
          //如果已经加载完成
          if (model.valueField) {
            let record = model.list.data.records.filter(r => {
              return (values as any[]).indexOf(getFieldValue(r, model.valueField)) !== -1;
            });
            if (record && record.length) {
              dispatch(restActions.selects(record));
            }
          }
          dispatch({
            type: ActionType.InitValues,
            payload: []
          });
        } else {
          dispatch({
            type: ActionType.InitValues,
            payload: values
          });
        }
      };
    },
    selects: (ts: T[]) => {
      return async (dispatch, getState) => {
        dispatch({
          type: ActionType.Selections,
          payload: ts
        });

        onSelects && onSelects.apply(actions, [ts, dispatch, getState]);
      };
    },
    selectsByValue: (values: string[] | number[]) => {
      return async (dispatch, getState) => {
        let model = getRecord(getState);
        if (model.valueField) {
          let records = model.list.data.records.filter(r => {
            return (values as any[]).indexOf(getFieldValue(r, model.valueField)) !== -1;
          });
          if (records && records.length) {
            dispatch(actions.selects(records));
          }
        }
      };
    },
    polling: ({ filter, onError, retryTimes = 3, delayTime = 3000, visibleCheck = true }) => {
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
    startPolling: ({ onError = null, retryTimes = 3, delayTime = 3000, visibleCheck = true } = {}) => {
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
        _polling.timer = setTimeout(() => {
          doPolling(undefined, dispatch, visibleCheck);
        }, _polling.delayTime);
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
    },
    clearSelection: () => {
      return async (dispatch: Redux.Dispatch, getState) => {
        dispatch({
          type: ActionType.ClearSelection
        });

        onClearSelection && onClearSelection.apply(actions, [dispatch, getState]);
      };
    }
  };

  const actions = extend({}, baseActions, restActions);
  return actions;
}

export interface FFListActionFactoryParams<T, TFilter, ExtendParamsT = any, TSFilter = any>
  extends FFListActionParams<T, TFilter, ExtendParamsT, TSFilter> {}

export function createFFListActionsFactory<T, TFilter, ExtendParamsT = any, TSFilter = any>({
  ...props
}: FFListActionFactoryParams<T, TFilter, ExtendParamsT, TSFilter>): (id: string) => FFListAction<T, TFilter, TSFilter> {
  return id =>
    createFFListActions<T, TFilter, ExtendParamsT, TSFilter>({
      id,
      ...props
    });
}
