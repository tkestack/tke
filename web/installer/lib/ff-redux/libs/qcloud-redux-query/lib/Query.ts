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
import { Dispatch } from 'redux';

import { ActionTypesEnum, PagingQuery, QueryState, ReduxAction } from '../../../src/base';
import { extend, merge } from '../../qcloud-lib';
import { FetcherActionCreator } from '../../qcloud-redux-fetcher';

// export interface QueryState<TFilter> {
//   date?: DateQuery;
//   paging?: PagingQuery;
//   search?: string;
//   keyword?: string;
//   filter?: TFilter;
//   sort?: SortQuery;
// }

// enum ActionTypesEnum {
//   ChangeDate = "ChangeDate",
//   ChangePaging = "ChangePaging",
//   RestPaging = "RestPaging",
//   ChangeKeyword = "ChangeKeyword",
//   ApplyFilter = "ApplyFilter",
//   ApplyPolling = "ApplyPolling",
//   SortBy = "SortBy",
//   PerformSearch = "PerformSearch",
//   Reset = "Reset"
// }

/* eslint-disable */
const defaultInitialState: QueryState<any> = {
  paging: { pageIndex: 1, pageSize: 20 },
  keyword: '',
  search: '',
  filter: {},
  sort: {}
};

export function generateQueryReducer<TFilter>({
  actionType,
  initialState
}: {
  actionType: string | number;
  initialState?: QueryState<TFilter>;
}) {
  const prefix = actionType as string;
  initialState = merge({}, defaultInitialState, initialState);

  function resetPaging(paging: PagingQuery): PagingQuery {
    return { pageIndex: 1, pageSize: paging.pageSize };
  }

  function reducer(state: QueryState<TFilter> = initialState, action: ReduxAction<any>) {
    let { paging, keyword, search, filter, sort } = state;

    switch (action.type) {
      case prefix + ActionTypesEnum.ChangePaging:
        paging = action.payload;
        break;

      case prefix + ActionTypesEnum.RestPaging:
        paging = { pageIndex: 1, pageSize: paging.pageSize };
        break;

      case prefix + ActionTypesEnum.ChangeKeyword:
        keyword = action.payload;
        break;

      case prefix + ActionTypesEnum.PerformSearch:
        search = action.payload;
        keyword = search;
        paging = resetPaging(paging);
        break;

      case prefix + ActionTypesEnum.ApplyFilter:
        filter = extend({}, filter, action.payload);
        paging = resetPaging(paging);
        break;

      case prefix + ActionTypesEnum.ApplyPolling:
        filter = extend({}, filter, action.payload);
        break;

      case prefix + ActionTypesEnum.SortBy:
        sort = action.payload;
        break;

      case prefix + ActionTypesEnum.Reset:
        return initialState;

      default:
        return state;
    }

    return { paging, keyword, search, filter, sort };
  }

  return reducer;
}

export interface QueryActionCreator<TFilter> {
  changePaging(nextPaging: PagingQuery): any;
  resetPaging(): any;
  changeKeyword(nextKeyword: string): any;
  performSearch(keyword: string): any;
  applyFilter(nextFilter: TFilter): any;
  applyPolling(nextFilter: TFilter): any;
  sortBy(by: string, desc?: boolean): any;
  reset(): any;
}

export function generateQueryActionCreator<TFilter>({
  actionType,
  bindFetcher,
  bindFetchers
}: {
  actionType: number | string;
  bindFetcher?: FetcherActionCreator;
  bindFetchers?: Array<{
    actions: ActionTypesEnum[] | 'all';
    fetcher: FetcherActionCreator;
  }>;
}): QueryActionCreator<TFilter> {
  const prefix = actionType as string;

  function changePaging(nextPaging: PagingQuery): ReduxAction<PagingQuery> {
    return {
      type: prefix + ActionTypesEnum.ChangePaging,
      payload: nextPaging
    };
  }

  function resetPaging(): ReduxAction<void> {
    return {
      type: prefix + ActionTypesEnum.RestPaging
    };
  }

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

  const creator: QueryActionCreator<TFilter> = {
    changePaging,
    applyFilter,
    applyPolling,
    resetPaging,
    changeKeyword,
    performSearch,
    sortBy,
    reset
  };

  if (bindFetcher || bindFetchers) {
    const wrap = (creator: Function) => (...args: any[]) => (dispatch: Dispatch, getState: () => any) => {
      const queryAction: ReduxAction<any> = creator.apply(null, args);
      dispatch(queryAction);
      if (bindFetcher) {
        dispatch(bindFetcher.fetch());
      }
      if (bindFetchers) {
        bindFetchers.forEach(bind => {
          if (bind.actions === 'all') {
            dispatch(bind.fetcher.fetch());
          } else if (bind.actions instanceof Array) {
            const bindActions = bind.actions as ActionTypesEnum[];
            if (bindActions.map(x => (actionType as string) + x).indexOf(queryAction.type as string) > -1) {
              dispatch(bind.fetcher.fetch());
            }
          }
        });
      }
    };
    return {
      changeKeyword,
      changePaging: wrap(creator.changePaging),
      resetPaging: wrap(creator.resetPaging),
      applyFilter: wrap(creator.applyFilter),
      applyPolling: wrap(creator.applyPolling),
      performSearch: wrap(creator.performSearch),
      sortBy: wrap(creator.sortBy),
      reset
    };
  }
  return creator;
}
