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

import { Dispatch } from 'redux';

import {
  ActionTypesEnum,
  FetcherPayload,
  FetcherTrigger,
  FetchOptions,
  PagingQuery,
  ReduxAction,
  FetcherState,
  QueryState,
  DataType
} from './Model';

/* eslint-disable */
export interface BaseActionParams<TData, TFilter, TSilter = any> {
  actionType: number | string;
  fetcher: (getState: () => any, options: FetchOptions, dispatch: Redux.Dispatch) => Promise<TData>;
  finish?: (dispatch: Redux.Dispatch, getState: () => any) => any;
  dataType?: DataType;

  getFetchReducer?: (getState) => FetcherState<TData>;
  getQueryReducer?: (getState) => QueryState<TFilter, TSilter>;
}
export interface BaseAction<TData, TFilter, TSilter = any> {
  fetch(options?: FetchOptions): any;
  update?: (data) => void;
  clearFetch?: () => void;

  changeKeyword(nextKeyword: string): any;
  performSearch(keyword: string): any;

  changeFilter?: (nextFilter: TFilter) => any;
  applyFilter(nextFilter: TFilter): any;
  applyPolling(nextFilter?: TFilter): any;

  changeSearchFilter?: (nextFilter: TSilter) => any;
  applySearchFilter?: (nextFilter: TSilter) => any;

  changePaging(nextPaging: PagingQuery): any;
  changePagingIndex?: (pageIndex: number) => any;
  resetPaging(): any;
  next?: () => any;

  sortBy(by: string, desc?: boolean): any;
  reset(): any;
}

export function createBaseAction<TData, TFilter, TSFilter = any>({
  actionType,
  fetcher,
  finish,
  dataType = DataType.List,
  getFetchReducer,
  getQueryReducer
}: BaseActionParams<TData, TFilter, TSFilter>): BaseAction<TData, TFilter, TSFilter> {
  type ActionType = ReduxAction<FetcherPayload<TData>>;

  let syncId = 0;
  let lastLoadingTimeout = 0;

  function fetch(options?: FetchOptions) {
    return (dispatch: Dispatch, getState: () => any) => {
      const fetchAction: ActionType = {
        type: actionType + (FetcherTrigger.Start as any),
        payload: {
          trigger: FetcherTrigger.Start
        }
      };
      dispatch(fetchAction);

      if (dataType === DataType.List && getFetchReducer && getQueryReducer) {
        const fetchReducer = getFetchReducer(getState);
        const queryReducer = getQueryReducer(getState);

        let pagingIndex = queryReducer.paging.pageIndex;
        let page = fetchReducer.pages[pagingIndex - 2];
        let continueToken = page ? (page.data as any).continueToken : '';

        const continueAction: ActionType = {
          type: actionType + ActionTypesEnum.ChangeContinueToken,
          payload: continueToken
        };
        dispatch(continueAction);
      }

      // keep the action is always dispatch with the latest promise result by `start()`
      const currentSyncId = ++syncId;
      const dispatchOnSync = (action: any) => {
        if (syncId === currentSyncId) {
          dispatch(action);
        }
        clearTimeout(lastLoadingTimeout);
      };

      dispatch(loading());

      const fetched = fetcher(getState, options, dispatch).then(
        data => dispatchOnSync(done(getState, data)),
        error => dispatchOnSync(fail(getState, error))
      );

      if (typeof finish === 'function') {
        fetched.then(() => finish(dispatch, getState));
      }
    };
  }

  function loading(): ActionType {
    return {
      type: actionType + (FetcherTrigger.Loading as any),
      payload: {
        trigger: FetcherTrigger.Loading
      }
    };
  }

  function done(getState, data: TData): ActionType {
    if (getQueryReducer) {
      const query = getQueryReducer(getState);
      return {
        type: actionType + (FetcherTrigger.Done as any),
        payload: {
          trigger: FetcherTrigger.Done,
          data,
          pageIndex: query.paging.pageIndex,
          append: query.paging.append,
          clear: query.paging.clear
        }
      };
    } else {
      return {
        type: actionType + (FetcherTrigger.Done as any),
        payload: {
          trigger: FetcherTrigger.Done,
          data
        }
      };
    }
  }

  function fail(getState, error: Error): ActionType {
    if (getQueryReducer) {
      const query = getQueryReducer(getState);
      return {
        type: actionType + (FetcherTrigger.Fail as any),
        payload: {
          trigger: FetcherTrigger.Fail,
          error,
          pageIndex: query.paging.pageIndex,
          append: query.paging.append,
          clear: query.paging.clear
        }
      };
    } else {
      return {
        type: actionType + (FetcherTrigger.Fail as any),
        payload: {
          trigger: FetcherTrigger.Fail,
          error
        }
      };
    }
  }

  function update(data: TData): ActionType {
    return {
      type: actionType + (FetcherTrigger.Update as any),
      payload: {
        trigger: FetcherTrigger.Update,
        data
      }
    };
  }

  function clearFetch(): ActionType {
    syncId = -1;
    return {
      type: actionType + (FetcherTrigger.Clear as any),
      payload: {
        trigger: FetcherTrigger.Clear
      }
    };
  }

  const prefix = actionType as string;

  function changeKeyword(nextKeyword: string): ReduxAction<string> {
    return {
      type: prefix + ActionTypesEnum.ChangeKeyword,
      payload: nextKeyword
    };
  }
  function performSearch(nextSearch: string): ReduxAction<string> {
    return {
      type: prefix + ActionTypesEnum.PerformSearch,
      payload: nextSearch
    };
  }

  function changeFilter(nextFilter: TFilter): ReduxAction<TFilter> {
    return {
      type: prefix + ActionTypesEnum.ChangeFilter,
      payload: nextFilter
    };
  }
  function applyFilter(nextFilter: TFilter): ReduxAction<TFilter> {
    return {
      type: prefix + ActionTypesEnum.ApplyFilter,
      payload: nextFilter
    };
  }
  function applyPolling(nextFilter: TFilter): ReduxAction<TFilter> {
    return {
      type: prefix + ActionTypesEnum.ApplyPolling,
      payload: nextFilter
    };
  }

  function changeSearchFilter(nextFilter: TSFilter): ReduxAction<TSFilter> {
    return {
      type: prefix + ActionTypesEnum.ChangeSearchFilter,
      payload: nextFilter
    };
  }
  function applySearchFilter(nextFilter: TSFilter): ReduxAction<TSFilter> {
    return {
      type: prefix + ActionTypesEnum.ApplySearchFilter,
      payload: nextFilter
    };
  }

  function changePaging(nextPaging: PagingQuery): ReduxAction<PagingQuery> {
    return {
      type: prefix + ActionTypesEnum.ChangePaging,
      payload: nextPaging
    };
  }

  function changePagingIndex(pageIndex: number): ReduxAction<number> {
    return {
      type: prefix + ActionTypesEnum.ChangePagingIndex,
      payload: pageIndex
    };
  }

  function resetPaging(): ReduxAction<void> {
    return {
      type: prefix + ActionTypesEnum.RestPaging
    };
  }
  function next(): ReduxAction<void> {
    return {
      type: prefix + ActionTypesEnum.Next
    };
  }

  function sortBy(by: string, desc?: boolean) {
    return {
      type: prefix + ActionTypesEnum.SortBy,
      payload: { by, desc }
    };
  }

  function reset() {
    return {
      type: prefix + ActionTypesEnum.Reset,
      payload: ''
    };
  }

  const wrap = (action: Function) => (...args: any[]) => (dispatch: Dispatch, getState: () => any) => {
    const queryAction: ReduxAction<any> = action.apply(null, args);
    dispatch(queryAction);
    dispatch(fetch());
  };
  return {
    fetch,
    update,
    clearFetch,
    changeKeyword: changeKeyword,
    performSearch: wrap(performSearch),

    changeFilter: changeFilter,
    applyFilter: wrap(applyFilter),
    applyPolling: wrap(applyPolling),

    changeSearchFilter: changeSearchFilter,
    applySearchFilter: wrap(applySearchFilter),

    changePaging: wrap(changePaging),
    changePagingIndex: wrap(changePagingIndex),
    resetPaging: wrap(resetPaging),
    next: wrap(next),

    sortBy: wrap(sortBy),
    reset: reset
  };
}
